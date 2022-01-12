package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	transfer "github.com/teleport-network/teleport/syscontracts/xibc_transfer"
)

var _ evmtypes.EvmHooks = (*Keeper)(nil)

func (k Keeper) PostTxProcessing(
	ctx sdk.Context,
	from common.Address,
	to *common.Address,
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

		sendPacketEvent, err := transferContract.Unpack(event.Name, log.Data)
		if err != nil {
			k.Logger(ctx).Error("failed to unpack send packet event", "error", err.Error())
			return err
		}

		// srcChain := sendPacketEvent[0].(string)
		destChain := sendPacketEvent[1].(string)
		relayChain := sendPacketEvent[2].(string)
		sender := sendPacketEvent[3].(string)
		receiver := sendPacketEvent[4].(string)
		amount := sendPacketEvent[5].(*big.Int)
		token := sendPacketEvent[6].(string)
		oriToken := sendPacketEvent[7].(string)

		// send cross chain transfer
		if err := k.SendTransfer(
			ctx,
			destChain,
			relayChain,
			sender,
			receiver,
			amount.Bytes(),
			token,
			oriToken,
		); err != nil {
			k.Logger(ctx).Debug(
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
