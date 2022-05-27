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

// SlashUnbondingDelegation slash an unbonding delegation and update the pool
// return the amount that would have been slashed assuming
// the unbonding delegation had enough stake to slash
// (the amount actually slashed may be less if there's
// insufficient stake remaining)
func (k SlashStakingKeeper) SlashUnbondingDelegation(ctx sdk.Context, unbondingDelegation stypes.UnbondingDelegation,
	infractionHeight int64, slashFactor sdk.Dec) (totalSlashAmount sdk.Int) {
	now := ctx.BlockHeader().Time
	totalSlashAmount = sdk.ZeroInt()
	slashedAmount := sdk.ZeroInt()

	// perform slashing on all entries within the unbonding delegation
	for i, entry := range unbondingDelegation.Entries {
		// If unbonding started before this height, stake didn't contribute to infraction
		if entry.CreationHeight < infractionHeight {
			continue
		}

		if entry.IsMature(now) {
			// Unbonding delegation no longer eligible for slashing, skip it
			continue
		}

		// Calculate slash amount proportional to stake contributing to infraction
		slashAmountDec := slashFactor.MulInt(entry.InitialBalance)
		slashAmount := slashAmountDec.TruncateInt()
		totalSlashAmount = totalSlashAmount.Add(slashAmount)

		// Don't slash more tokens than held
		// Possible since the unbonding delegation may already
		// have been slashed, and slash amounts are calculated
		// according to stake held at time of infraction
		unbondingSlashAmount := sdk.MinInt(slashAmount, entry.Balance)

		// Update unbonding delegation if necessary
		if unbondingSlashAmount.IsZero() {
			continue
		}

		slashedAmount = slashedAmount.Add(unbondingSlashAmount)
		entry.Balance = entry.Balance.Sub(unbondingSlashAmount)
		unbondingDelegation.Entries[i] = entry
		k.SetUnbondingDelegation(ctx, unbondingDelegation)
	}

	if slashedAmount.IsPositive() {
		coins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), slashedAmount))
		if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, stypes.NotBondedPoolName, authtypes.FeeCollectorName, coins); err != nil {
			panic(err)
		}
	}

	return totalSlashAmount
}

// SlashRedelegation slash a redelegation and update the pool
// return the amount that would have been slashed assuming
// the unbonding delegation had enough stake to slash
// (the amount actually slashed may be less if there's
// insufficient stake remaining)
// NOTE this is only slashing for prior infractions from the source validator
func (k SlashStakingKeeper) SlashRedelegation(ctx sdk.Context, srcValidator stypes.Validator, redelegation stypes.Redelegation,
	infractionHeight int64, slashFactor sdk.Dec) (totalSlashAmount sdk.Int) {
	now := ctx.BlockHeader().Time
	totalSlashAmount = sdk.ZeroInt()
	bondedSlashedAmount, notBondedSlashedAmount := sdk.ZeroInt(), sdk.ZeroInt()

	// perform slashing on all entries within the redelegation
	for _, entry := range redelegation.Entries {
		// If redelegation started before this height, stake didn't contribute to infraction
		if entry.CreationHeight < infractionHeight {
			continue
		}

		if entry.IsMature(now) {
			// Redelegation no longer eligible for slashing, skip it
			continue
		}

		// Calculate slash amount proportional to stake contributing to infraction
		slashAmountDec := slashFactor.MulInt(entry.InitialBalance)
		slashAmount := slashAmountDec.TruncateInt()
		totalSlashAmount = totalSlashAmount.Add(slashAmount)

		// Unbond from target validator
		sharesToUnbond := slashFactor.Mul(entry.SharesDst)
		if sharesToUnbond.IsZero() {
			continue
		}

		valDstAddr, err := sdk.ValAddressFromBech32(redelegation.ValidatorDstAddress)
		if err != nil {
			panic(err)
		}

		delegatorAddress, err := sdk.AccAddressFromBech32(redelegation.DelegatorAddress)
		if err != nil {
			panic(err)
		}

		delegation, found := k.GetDelegation(ctx, delegatorAddress, valDstAddr)
		if !found {
			// If deleted, delegation has zero shares, and we can't unbond any more
			continue
		}

		if sharesToUnbond.GT(delegation.Shares) {
			sharesToUnbond = delegation.Shares
		}

		tokensToBurn, err := k.Unbond(ctx, delegatorAddress, valDstAddr, sharesToUnbond)
		if err != nil {
			panic(fmt.Errorf("error unbonding delegator: %v", err))
		}

		dstValidator, found := k.GetValidator(ctx, valDstAddr)
		if !found {
			panic("destination validator not found")
		}

		// tokens of a redelegation currently live in the destination validator
		// therefor we must burn tokens from the destination-validator's bonding status
		switch {
		case dstValidator.IsBonded():
			bondedSlashedAmount = bondedSlashedAmount.Add(tokensToBurn)
		case dstValidator.IsUnbonded() || dstValidator.IsUnbonding():
			notBondedSlashedAmount = notBondedSlashedAmount.Add(tokensToBurn)
		default:
			panic("unknown validator status")
		}
	}

	if bondedSlashedAmount.IsPositive() {
		coins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), bondedSlashedAmount))
		if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, stypes.BondedPoolName, authtypes.FeeCollectorName, coins); err != nil {
			panic(err)
		}
	}

	if notBondedSlashedAmount.IsPositive() {
		coins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), notBondedSlashedAmount))
		if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, stypes.NotBondedPoolName, authtypes.FeeCollectorName, coins); err != nil {
			panic(err)
		}
	}

	return totalSlashAmount
}
