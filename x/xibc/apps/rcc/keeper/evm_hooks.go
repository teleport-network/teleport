package keeper

import (
	"github.com/teleport-network/teleport/x/xibc/apps/rcc/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	rcc "github.com/teleport-network/teleport/syscontracts/xibc_rcc"
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
	from common.Address,
	to *common.Address,
	receipt *ethtypes.Receipt,
) error {
	rccContract := rcc.RCCContract.ABI

	for i, log := range receipt.Logs {
		if log.Address != rcc.RCCContractAddress {
			continue
		}

		if len(log.Topics) == 0 {
			continue
		}

		eventID := log.Topics[0]
		event, err := rccContract.EventByID(eventID)
		if err != nil {
			return err
		}

		if event.Name == "Ack" {
			return nil
		}
		packet, err := rccContract.Unpack(event.Name, log.Data)
		if err != nil {
			h.k.Logger(ctx).Error("failed to unpack send packet event", "error", err.Error())
			return err
		}

		var sendPacketEvent types.RCCEventSendPacketData
		err = sendPacketEvent.DecodeInterface(packet[0])
		if err != nil {
			h.k.Logger(ctx).Error("failed to decode send packet event", "error", err.Error())
			return err
		}

		// send cross chain contract call
		if err := h.k.SendRemoteContractCall(
			ctx,
			sendPacketEvent.DestChain,
			sendPacketEvent.RelayChain,
			sendPacketEvent.Sequence,
			sendPacketEvent.Sender,
			sendPacketEvent.ContractAddress,
			sendPacketEvent.Data,
		); err != nil {
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
