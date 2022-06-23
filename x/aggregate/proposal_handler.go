package aggregate

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/teleport-network/teleport/x/aggregate/keeper"
	"github.com/teleport-network/teleport/x/aggregate/types"
)

// NewAggregateProposalHandler creates a governance handler to manage new proposal types.
// It enables RegisterTokenPairProposal to propose a registration of token mapping
func NewAggregateProposalHandler(k *keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.RegisterCoinProposal:
			return handleRegisterCoinProposal(ctx, k, c)
		case *types.AddCoinProposal:
			return handleAddCoinProposal(ctx, k, c)
		case *types.RegisterERC20Proposal:
			return handleRegisterERC20Proposal(ctx, k, c)
		case *types.ToggleTokenConversionProposal:
			return handleToggleConversionProposal(ctx, k, c)
		case *types.RegisterERC20TraceProposal:
			return handleRegisterERC20TraceProposal(ctx, k, c)
		case *types.EnableTimeBasedSupplyLimitProposal:
			return handleEnableTimeBasedSupplyLimitProposal(ctx, k, c)
		case *types.DisableTimeBasedSupplyLimitProposal:
			return handleDisableTimeBasedSupplyLimitProposal(ctx, k, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s proposal content type: %T", types.ModuleName, c)
		}
	}
}

func handleRegisterCoinProposal(ctx sdk.Context, k *keeper.Keeper, p *types.RegisterCoinProposal) error {
	pair, err := k.RegisterCoin(ctx, p.Metadata)
	if err != nil {
		return err
	}
	return ctx.EventManager().EmitTypedEvent(
		&types.EventRegisterTokens{
			Denom:      pair.Denoms,
			Erc20Token: pair.ERC20Address,
		},
	)
}

func handleAddCoinProposal(ctx sdk.Context, k *keeper.Keeper, p *types.AddCoinProposal) error {
	pair, err := k.AddCoin(ctx, p.Metadata, p.ContractAddress)
	if err != nil {
		return err
	}
	return ctx.EventManager().EmitTypedEvent(
		&types.EventRegisterTokens{
			Denom:      pair.Denoms,
			Erc20Token: pair.ERC20Address,
		},
	)
}

func handleRegisterERC20Proposal(ctx sdk.Context, k *keeper.Keeper, p *types.RegisterERC20Proposal) error {
	pair, err := k.RegisterERC20(ctx, common.HexToAddress(p.ERC20Address))
	if err != nil {
		return err
	}
	return ctx.EventManager().EmitTypedEvent(
		&types.EventRegisterTokens{
			Denom:      pair.Denoms,
			Erc20Token: pair.ERC20Address,
		},
	)
}

func handleToggleConversionProposal(ctx sdk.Context, k *keeper.Keeper, p *types.ToggleTokenConversionProposal) error {
	pair, err := k.ToggleConversion(ctx, p.Token)
	if err != nil {
		return err
	}
	return ctx.EventManager().EmitTypedEvent(
		&types.EventRegisterTokens{
			Denom:      pair.Denoms,
			Erc20Token: pair.ERC20Address,
		},
	)
}

func handleRegisterERC20TraceProposal(ctx sdk.Context, k *keeper.Keeper, p *types.RegisterERC20TraceProposal) error {
	if err := k.RegisterERC20Trace(ctx, common.HexToAddress(p.ERC20Address), p.OriginToken, p.OriginChain, uint8(p.Scale)); err != nil {
		return err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRegisterERC20Trace,
			sdk.NewAttribute(types.AttributeKeyERC20Token, p.ERC20Address),
			sdk.NewAttribute(types.AttributeKeyOriginToken, p.OriginToken),
			sdk.NewAttribute(types.AttributeKeyOriginChain, p.OriginChain),
		),
	)
	return nil
}

func handleEnableTimeBasedSupplyLimitProposal(ctx sdk.Context, k *keeper.Keeper, p *types.EnableTimeBasedSupplyLimitProposal) error {
	timePeriod, _ := new(big.Int).SetString(p.TimePeriod, 10)
	timeBasedLimit, _ := new(big.Int).SetString(p.TimeBasedLimit, 10)
	maxAmount, _ := new(big.Int).SetString(p.MaxAmount, 10)
	minAmount, _ := new(big.Int).SetString(p.MinAmount, 10)

	if err := k.EnableTimeBasedSupplyLimit(
		ctx,
		common.HexToAddress(p.ERC20Address),
		timePeriod,
		timeBasedLimit,
		maxAmount,
		minAmount,
	); err != nil {
		return err
	}
	_ = ctx.EventManager().EmitTypedEvent(
		&types.EnableTimeBasedSupplyLimitProposal{
			ERC20Address:   p.ERC20Address,
			TimePeriod:     p.TimePeriod,
			TimeBasedLimit: p.TimeBasedLimit,
			MaxAmount:      p.MaxAmount,
			MinAmount:      p.MinAmount,
		},
	)
	return nil
}

func handleDisableTimeBasedSupplyLimitProposal(ctx sdk.Context, k *keeper.Keeper, p *types.DisableTimeBasedSupplyLimitProposal) error {
	if err := k.DisableTimeBasedSupplyLimit(ctx, common.HexToAddress(p.ERC20Address)); err != nil {
		return err
	}
	_ = ctx.EventManager().EmitTypedEvent(
		&types.DisableTimeBasedSupplyLimitProposal{
			ERC20Address: p.ERC20Address,
		},
	)
	return nil
}
