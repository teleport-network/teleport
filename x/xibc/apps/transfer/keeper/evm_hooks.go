package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	transfer "github.com/teleport-network/teleport/syscontracts/xibc_transfer"
	"github.com/teleport-network/teleport/x/xibc/apps/transfer/types"
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
	transferContract := transfer.TransferContract.ABI

	for i, log := range receipt.Logs {
		if log.Address != transfer.TransferContractAddress {
			continue
		}

		if len(log.Topics) == 0 {
			continue
		}

		eventID := log.Topics[0]
		event, err := transferContract.EventByID(eventID)
		if err != nil {
			return err
		}
		packet, err := transferContract.Unpack(event.Name, log.Data)
		if err != nil {
			h.k.Logger(ctx).Error("failed to unpack send packet event", "error", err.Error())
			return err
		}

		var sendPacketEvent types.TransferEventSendPacketData
		if err = sendPacketEvent.DecodeInterface(packet[0]); err != nil {
			h.k.Logger(ctx).Error("failed to decode send packet event", "error", err.Error())
			return err
		}

		// send cross chain transfer
		if err := h.k.SendTransfer(
			ctx,
			sendPacketEvent.DestChain,
			sendPacketEvent.RelayChain,
			sendPacketEvent.Sequence,
			sendPacketEvent.Sender,
			sendPacketEvent.Receiver,
			sendPacketEvent.Amount.Bytes(),
			sendPacketEvent.Token,
			sendPacketEvent.OriToken,
		); err != nil {
			h.k.Logger(ctx).Debug(
				"failed to process EVM hook for XIBC transfer",
				"tx-hash", receipt.TxHash.Hex(),
				"log-idx", i,
				"error", err.Error(),
			)
			return err
		}
	}

	return nil
}
