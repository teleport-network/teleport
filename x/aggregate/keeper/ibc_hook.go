package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	"github.com/cosmos/ibc-go/v3/modules/core/exported"
	"github.com/ethereum/go-ethereum/common"

	"github.com/teleport-network/teleport/x/aggregate/types"

	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
)

// OnRecvPacket will get the denom name from ibc ,generate by port/channel/denom
func (k Keeper) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	ack exported.Acknowledgement,
) exported.Acknowledgement {
	event := &types.EventIBCAggregate{
		Status:             types.STATUS_UNKNOWN,
		Message:            "",
		Sequence:           packet.Sequence,
		SourceChannel:      packet.SourceChannel,
		DestinationChannel: packet.DestinationChannel,
	}
	cctx, write := ctx.CacheContext()

	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		event.Status = types.STATUS_FAILED
		event.Message = err.Error()
		_ = ctx.EventManager().EmitTypedEvent(event)
		return nil
	}
	transferAmount, ok := sdk.NewIntFromString(data.Amount)
	if !ok {
		event.Status = types.STATUS_FAILED
		event.Message = "Change data.Amount type to int error"
		_ = ctx.EventManager().EmitTypedEvent(event)
		return nil
	}
	receiver, _ := sdk.AccAddressFromBech32(data.Receiver)
	denom, err := types.IBCDenom(packet.GetDestPort(), packet.GetDestChannel(), data.Denom)
	if err != nil {
		event.Status = types.STATUS_FAILED
		event.Message = err.Error()
		_ = ctx.EventManager().EmitTypedEvent(event)
		return nil
	}

	if !k.IsDenomRegistered(ctx, denom) {
		event.Status = types.STATUS_FAILED
		event.Message = fmt.Sprintf("denom %s not registered", denom)
		_ = ctx.EventManager().EmitTypedEvent(event)
		return nil
	}
	msg := types.NewMsgConvertCoin(
		sdk.NewCoin(denom, transferAmount),
		common.BytesToAddress(receiver.Bytes()),
		receiver,
	)
	// use cctx to ConvertCoin
	context := sdk.WrapSDKContext(cctx)
	_, err = k.ConvertCoin(context, msg)
	if err != nil {
		event.Status = types.STATUS_FAILED
		event.Message = err.Error()
		_ = ctx.EventManager().EmitTypedEvent(event)
		return nil
	}

	write()
	ctx.EventManager().EmitEvents(cctx.EventManager().Events())
	event.Status = types.STATUS_SUCCESS
	_ = ctx.EventManager().EmitTypedEvent(event)
	return nil
}

func (k Keeper) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
) error {
	// nothing to do
	return nil
}

func (k Keeper) SendPacket(ctx sdk.Context, channelCap *capabilitytypes.Capability, packet exported.PacketI) error {
	return k.ics4Wrapper.SendPacket(ctx, channelCap, packet)
}

func (k Keeper) WriteAcknowledgement(ctx sdk.Context, channelCap *capabilitytypes.Capability, packet exported.PacketI, ack exported.Acknowledgement) error {
	return k.ics4Wrapper.WriteAcknowledgement(ctx, channelCap, packet, ack)
}
