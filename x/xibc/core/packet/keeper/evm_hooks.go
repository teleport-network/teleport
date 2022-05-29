package keeper

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	packetcontract "github.com/teleport-network/teleport/syscontracts/xibc_packet"
	"github.com/teleport-network/teleport/x/xibc/core/packet/types"

	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
)

// Hooks wrapper struct for packet keeper
type Hooks struct {
	k Keeper
}

var _ evmtypes.EvmHooks = Hooks{}

// Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// PostTxProcessing implements EvmHooks.PostTxProcessing
func (h Hooks) PostTxProcessing(
	ctx sdk.Context,
	msg core.Message,
	receipt *ethtypes.Receipt,
) error {
	packetContract := packetcontract.PacketContract.ABI

	for i, log := range receipt.Logs {
		if log.Address != packetcontract.PacketContractAddress {
			continue
		}
		if len(log.Topics) == 0 {
			continue
		}
		eventID := log.Topics[0]
		event, err := packetContract.EventByID(eventID)
		if err != nil {
			h.k.Logger(ctx).Error("failed to unpack send packet event", "error", err.Error())
			return err
		}

		if event.Name != types.PacketSendEvent {
			continue
		}
		sendEvent, err := packetContract.Unpack(event.Name, log.Data)
		if err != nil {
			h.k.Logger(ctx).Error("failed to unpack send packet event", "error", err.Error())
			return err
		}
		var packetBytes []byte
		bzTmp, err := json.Marshal(sendEvent[0])
		if err != nil {
			fmt.Println(err.Error())
			h.k.Logger(ctx).Error("failed to decode interface to event", "error", err.Error())
			return err
		}
		err = json.Unmarshal(bzTmp, &packetBytes)
		if err != nil {
			fmt.Println(err.Error())
			h.k.Logger(ctx).Error("failed to decode interface to event", "error", err.Error())
			return err
		}

		var packet types.Packet
		if err = packet.ABIDecode(packetBytes); err != nil {
			h.k.Logger(ctx).Error("failed to decode packet", "error", err.Error())
			return err
		}
		// todo ? add data validate?
		//var transferData types.TransferData
		//if err = transferData.ABIDecode(packet.GetTransferData()); err != nil {
		//	h.k.Logger(ctx).Error("failed to decode packet", "error", err.Error())
		//	return err
		//}
		//fmt.Println("hooks : callData : ", hex.EncodeToString(packet.GetCallData()))
		//var callData types.CallData
		//if err = callData.ABIDecode(packet.GetCallData()); err != nil {
		//	h.k.Logger(ctx).Error("failed to decode packet", "error", err.Error())
		//	return err
		//}

		// send cross chain contract packet
		if err = h.k.SendPacket(
			ctx,
			&packet,
		); err != nil {
			h.k.Logger(ctx).Debug(
				"Packet EVM hook failed ,err on send packet",
				"txhash", receipt.TxHash.Hex(),
				"log-idx", i,
				"error", err.Error(),
			)
			return err
		}
	}

	return nil
}
