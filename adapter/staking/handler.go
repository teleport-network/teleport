package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	distypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"

	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/teleport-network/teleport/adapter/common"
	"github.com/teleport-network/teleport/syscontracts"
	"github.com/teleport-network/teleport/syscontracts/staking"
)

func (h *HookAdapter) HandleDelegated(ctx sdk.Context, log *ethtypes.Log) error {
	event := new(staking.StakingDelegated)
	if err := syscontracts.ParseLog(event, h.abi, log, "Delegated"); err != nil {
		return err
	}

	delegator, err := bech32.ConvertAndEncode(sdk.GetConfig().GetBech32AccountAddrPrefix(), event.Delegator.Bytes())
	if err != nil {
		return err
	}
	bondDenom := h.stakingKeeper.BondDenom(ctx)
	msg := &types.MsgDelegate{
		DelegatorAddress: delegator,
		ValidatorAddress: event.Validator,
		Amount:           sdk.NewCoin(bondDenom, sdk.NewIntFromBigInt(event.Amount)),
	}
	return common.ExecuteMsg(ctx, h.router, msg)
}

func (h *HookAdapter) HandleUndelegated(ctx sdk.Context, log *ethtypes.Log) error {
	event := new(staking.StakingUndelegated)
	if err := syscontracts.ParseLog(event, h.abi, log, "Undelegated"); err != nil {
		return err
	}

	delegator, err := bech32.ConvertAndEncode(sdk.GetConfig().GetBech32AccountAddrPrefix(), event.Delegator.Bytes())
	if err != nil {
		return err
	}
	bondDenom := h.stakingKeeper.BondDenom(ctx)
	msg := &types.MsgUndelegate{
		DelegatorAddress: delegator,
		ValidatorAddress: event.Validator,
		Amount:           sdk.NewCoin(bondDenom, sdk.NewIntFromBigInt(event.Amount)),
	}
	return common.ExecuteMsg(ctx, h.router, msg)
}

func (h *HookAdapter) HandleRedelegated(ctx sdk.Context, log *ethtypes.Log) error {
	event := new(staking.StakingRedelegated)
	if err := syscontracts.ParseLog(event, h.abi, log, "Redelegated"); err != nil {
		return err
	}

	delegator, err := bech32.ConvertAndEncode(sdk.GetConfig().GetBech32AccountAddrPrefix(), event.Delegator.Bytes())
	if err != nil {
		return err
	}
	bondDenom := h.stakingKeeper.BondDenom(ctx)
	msg := &types.MsgBeginRedelegate{
		DelegatorAddress:    delegator,
		ValidatorSrcAddress: event.ValidatorSrc,
		ValidatorDstAddress: event.ValidatorDest,
		Amount:              sdk.NewCoin(bondDenom, sdk.NewIntFromBigInt(event.Amount)),
	}
	return common.ExecuteMsg(ctx, h.router, msg)
}

func (h *HookAdapter) HandleWithdrew(ctx sdk.Context, log *ethtypes.Log) error {
	event := new(staking.StakingWithdrew)
	if err := syscontracts.ParseLog(event, h.abi, log, "Withdrew"); err != nil {
		return err
	}

	delegator, err := bech32.ConvertAndEncode(sdk.GetConfig().GetBech32AccountAddrPrefix(), event.Delegator.Bytes())
	if err != nil {
		return err
	}
	msg := &distypes.MsgWithdrawDelegatorReward{
		DelegatorAddress: delegator,
		ValidatorAddress: event.Validator,
	}
	return common.ExecuteMsg(ctx, h.router, msg)
}
