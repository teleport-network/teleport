package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	rcc "github.com/teleport-network/teleport/syscontracts/xibc_rcc"
)

var _ evmtypes.EvmHooks = (*Keeper)(nil)

func (k Keeper) PostTxProcessing(
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

		sendPacketEvent, err := rccContract.Unpack(event.Name, log.Data)
		if err != nil {
			k.Logger(ctx).Error("failed to unpack send packet event", "error", err.Error())
			return err
		}

		// srcChain := sendPacketEvent[0].(string)
		destChain := sendPacketEvent[1].(string)
		relayChain := sendPacketEvent[2].(string)
		sender := sendPacketEvent[3].(string)
		contractAddress := sendPacketEvent[4].(string)
		data := sendPacketEvent[5].([]byte)

		// send cross chain contract call
		if err := k.SendRemoteContractCall(
			ctx,
			destChain,
			relayChain,
			sender,
			contractAddress,
			data,
		); err != nil {
			k.Logger(ctx).Debug(
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
