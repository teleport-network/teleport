package staking

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	sk "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// SlashStakingKeeper inherits from staking keeper, and overwrite Slash
// It is designed to deal with the slashed amount by distributing it to validators, instead of burning it
type SlashStakingKeeper struct {
	sk.Keeper
	bankKeeper stypes.BankKeeper
}

func NewSlashStakingKeeper(cosmosStakingKeeper sk.Keeper, bankKeeper stypes.BankKeeper) SlashStakingKeeper {
	return SlashStakingKeeper{
		Keeper:     cosmosStakingKeeper,
		bankKeeper: bankKeeper,
	}
}

func (k SlashStakingKeeper) Slash(ctx sdk.Context, consAddr sdk.ConsAddress, infractionHeight int64, power int64, slashFactor sdk.Dec) {
	logger := k.Logger(ctx)

	if slashFactor.IsNegative() {
		panic(fmt.Errorf("attempted to slash with a negative slash factor: %v", slashFactor))
	}

	// Amount of slashing = slash slashFactor * power at time of infraction
	amount := k.TokensFromConsensusPower(ctx, power)
	slashAmountDec := amount.ToDec().Mul(slashFactor)
	slashAmount := slashAmountDec.TruncateInt()

	// ref https://github.com/cosmos/cosmos-sdk/issues/1348

	validator, found := k.GetValidatorByConsAddr(ctx, consAddr)
	if !found {
		// If not found, the validator must have been overslashed and removed - so we don't need to do anything
		// NOTE:  Correctness dependent on invariant that unbonding delegations / redelegations must also have been completely
		//        slashed in this case - which we don't explicitly check, but should be true.
		// Log the slash attempt for future reference (maybe we should tag it too)
		logger.Error(
			"WARNING: ignored attempt to slash a nonexistent validator; we recommend you investigate immediately",
			"validator", consAddr.String(),
		)
		return
	}

	// should not be slashing an unbonded validator
	if validator.IsUnbonded() {
		panic(fmt.Sprintf("should not be slashing unbonded validator: %s", validator.GetOperator()))
	}

	operatorAddress := validator.GetOperator()

	// call the before-modification hook
	k.BeforeValidatorModified(ctx, operatorAddress)

	// Track remaining slash amount for the validator
	// This will decrease when we slash unbondings and
	// redelegations, as that stake has since unbonded
	remainingSlashAmount := slashAmount

	switch {
	case infractionHeight > ctx.BlockHeight():
		// Can't slash infractions in the future
		panic(fmt.Sprintf(
			"impossible attempt to slash future infraction at height %d but we are at height %d",
			infractionHeight, ctx.BlockHeight()))

	case infractionHeight == ctx.BlockHeight():
		// Special-case slash at current height for efficiency - we don't need to
		// look through unbonding delegations or redelegations.
		logger.Info(
			"slashing at current height; not scanning unbonding delegations & redelegations",
			"height", infractionHeight,
		)

	case infractionHeight < ctx.BlockHeight():
		// Iterate through unbonding delegations from slashed validator
		unbondingDelegations := k.GetUnbondingDelegationsFromValidator(ctx, operatorAddress)
		for _, unbondingDelegation := range unbondingDelegations {
			amountSlashed := k.SlashUnbondingDelegation(ctx, unbondingDelegation, infractionHeight, slashFactor)
			if amountSlashed.IsZero() {
				continue
			}

			remainingSlashAmount = remainingSlashAmount.Sub(amountSlashed)
		}

		// Iterate through redelegations from slashed source validator
		redelegations := k.GetRedelegationsFromSrcValidator(ctx, operatorAddress)
		for _, redelegation := range redelegations {
			amountSlashed := k.SlashRedelegation(ctx, validator, redelegation, infractionHeight, slashFactor)
			if amountSlashed.IsZero() {
				continue
			}

			remainingSlashAmount = remainingSlashAmount.Sub(amountSlashed)
		}
	}

	// cannot decrease balance below zero
	tokensToSlash := sdk.MinInt(remainingSlashAmount, validator.Tokens)
	tokensToSlash = sdk.MaxInt(tokensToSlash, sdk.ZeroInt()) // defensive.

	// we need to calculate the *effective* slash fraction for distribution
	if validator.Tokens.IsPositive() {
		effectiveFraction := tokensToSlash.ToDec().QuoRoundUp(validator.Tokens.ToDec())
		// possible if power has changed
		if effectiveFraction.GT(sdk.OneDec()) {
			effectiveFraction = sdk.OneDec()
		}
		// call the before-slashed hook
		k.BeforeValidatorSlashed(ctx, operatorAddress, effectiveFraction)
	}

	// Deduct from validator's bonded tokens and update the validator.
	validator = k.RemoveValidatorTokens(ctx, validator, tokensToSlash)

	// The below is the key difference from cosmos slashing implementation.
	// Transfer the slashed tokens from the pool account to fee collector(used to distribute the coins to validators).
	if tokensToSlash.IsPositive() {
		coins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), tokensToSlash))
		switch validator.GetStatus() {
		case stypes.Bonded:
			if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, stypes.BondedPoolName, authtypes.FeeCollectorName, coins); err != nil {
				panic(err)
			}
		case stypes.Unbonding, stypes.Unbonded:
			if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, stypes.NotBondedPoolName, authtypes.FeeCollectorName, coins); err != nil {
				panic(err)
			}
		default:
			panic("invalid validator status")
		}
	}

	logger.Info(
		"validator slashed by slash factor",
		"validator", validator.GetOperator().String(),
		"slash_factor", slashFactor.String(),
		"slashed", tokensToSlash,
	)
}
