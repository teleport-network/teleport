package keeper

import (
	"encoding/json"

	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	// nolint: typecheck
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	multicall "github.com/teleport-network/teleport/syscontracts/xibc_multicall"
	"github.com/teleport-network/teleport/x/xibc/apps/multicall/types"
)

// Hooks wrapper struct for erc20 keeper
type Hooks struct {
	k Keeper
}

var _ evmtypes.EvmHooks = Hooks{}

// Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

func (h Hooks) PostTxProcessing(
	ctx sdk.Context,
	msg core.Message,
	receipt *ethtypes.Receipt,
) error {
	multicallContract := multicall.MultiCallContract.ABI

	for i, log := range receipt.Logs {
		if log.Address != multicall.MultiCallContractAddress {
			continue
		}

		if len(log.Topics) == 0 {
			continue
		}

		eventID := log.Topics[0]
		event, err := multicallContract.EventByID(eventID)
		if err != nil {
			return err
		}

		sendPacketEvent, err := multicallContract.Unpack(event.Name, log.Data)
		if err != nil {
			h.k.Logger(ctx).Error("failed to unpack send packet event", "error", err.Error())
			return err
		}

		sender := sendPacketEvent[0].(common.Address)
		bz, err := json.Marshal(sendPacketEvent[1])
		if err != nil {
			return err
		}

		var calldata types.MultiCallData
		if err := json.Unmarshal(bz, &calldata); err != nil {
			return err
		}

		// send cross chain contract call
		if err := h.k.SendMultiCall(ctx, sender, calldata); err != nil {
			h.k.Logger(ctx).Debug(
				"failed to process EVM hook for XIBC RCC",
				"tx-hash", receipt.TxHash.Hex(),
				"log-idx", i,
				"error", err.Error(),
			)
			return err
		}
	}

	return nil
}
