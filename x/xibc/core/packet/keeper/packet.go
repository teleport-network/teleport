package keeper

import (
	"bytes"
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	"github.com/teleport-network/teleport/x/xibc/core/packet/types"
	"github.com/teleport-network/teleport/x/xibc/exported"
)

// SendPacket is called by a module to send an XIBC packet on a port owned
// by the calling module to the corresponding module on the counterparty chain.
func (k Keeper) SendPacket(ctx sdk.Context, packet exported.PacketI) error {
	if err := packet.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "packet failed basic validation")
	}
	if packet.GetSrcChain() != k.clientKeeper.GetChainName(ctx) {
		return sdkerrors.Wrap(types.ErrInvalidPacket, "source chain of packet is not this chain")
	}

	if _, found := k.clientKeeper.GetClientState(ctx, packet.GetDstChain()); !found {
		return clienttypes.ErrClientNotFound
	}

	nextSequenceSend := k.GetNextSequenceSend(ctx, packet.GetSrcChain(), packet.GetDstChain())

	if packet.GetSequence() != nextSequenceSend {
		return sdkerrors.Wrapf(
			types.ErrInvalidPacket,
			"packet sequence ≠ next send sequence (%d ≠ %d)", packet.GetSequence(), nextSequenceSend,
		)
	}

	commitment, err := types.CommitPacket(packet)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrABIPack, "SendPacket error ,err : %s", err)
	}
	nextSequenceSend++
	k.SetNextSequenceSend(ctx, packet.GetSrcChain(), packet.GetDstChain(), nextSequenceSend)

	// set sequence in packet contract
	if _, err := k.CallPacket(ctx, "setSequence", packet.GetDstChain(), nextSequenceSend); err != nil {
		return err
	}

	k.SetPacketCommitment(ctx, packet.GetSrcChain(), packet.GetDstChain(), packet.GetSequence(), commitment)

	packetBytes, err := packet.ABIPack()
	if err != nil {
		return sdkerrors.Wrapf(types.ErrABIPack, "SendPacket error , err: %s", err)
	}
	_ = ctx.EventManager().EmitTypedEvent(
		&types.EventSendPacket{
			SrcChain: packet.GetSrcChain(),
			DstChain: packet.GetDstChain(),
			Sequence: strconv.FormatUint(packet.GetSequence(), 10),
			Packet:   packetBytes,
		},
	)

	k.Logger(ctx).Info("packet sent", "packet", fmt.Sprintf("%v : %v : %v", packet.GetSrcChain(), packet.GetDstChain(), packet.GetSequence()))
	return nil
}

// RecvPacket is called by a module to receive & process an XIBC packet
// sent on the corresponding port on the counterparty chain.
func (k Keeper) RecvPacket(ctx sdk.Context, msg *types.MsgRecvPacket) error {
	var packet types.Packet
	if err := packet.ABIDecode(msg.Packet); err != nil && packet.Sequence == 0 {
		return sdkerrors.Wrap(err, "packet failed basic validation")
	}

	if err := k.ValidatePacket(ctx, &packet); err != nil {
		return sdkerrors.Wrap(err, "packet failed basic validation")
	}
	// check if the packet receipt has been received already
	if _, found := k.GetPacketReceipt(ctx, packet.GetSrcChain(), packet.GetDstChain(), packet.GetSequence()); found {
		return sdkerrors.Wrapf(
			types.ErrInvalidPacket,
			"packet sequence (%d) already has been received", packet.GetSequence(),
		)
	}
	fromChain := packet.GetSrcChain()

	clientState, found := k.clientKeeper.GetClientState(ctx, fromChain)
	if !found {
		return sdkerrors.Wrap(clienttypes.ErrClientNotFound, fromChain)
	}

	commitment, err := types.CommitPacket(&packet)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrABIPack, "SendPacket error ,err : %s", err)
	}
	// use signer as tss client proof
	proof := msg.ProofCommitment
	if clientState.ClientType() == exported.TSS {
		proof = []byte(msg.Signer)
	}
	// verify that the counterparty did commit to sending this packet
	if err := clientState.VerifyPacketCommitment(
		ctx,
		k.clientKeeper.ClientStore(ctx, fromChain),
		k.cdc,
		msg.ProofHeight,
		proof,
		packet.GetSrcChain(),
		packet.GetDstChain(),
		packet.GetSequence(),
		commitment,
	); err != nil {
		return sdkerrors.Wrapf(err, "failed packet commitment verification for client (%s)", fromChain)
	}

	// All verification complete, update state
	k.SetPacketReceipt(ctx, packet.GetSrcChain(), packet.GetDstChain(), packet.GetSequence())

	// log that a packet has been received & executed
	k.Logger(ctx).Info("packet received", "packet", fmt.Sprintf("%v", packet))

	_ = ctx.EventManager().EmitTypedEvent(
		&types.EventRecvPacket{
			SrcChain: packet.GetSrcChain(),
			DstChain: packet.GetDstChain(),
			Sequence: strconv.FormatUint(packet.GetSequence(), 10),
			Packet:   msg.Packet,
		},
	)

	chainName := k.clientKeeper.GetChainName(ctx)

	// teleport is relayChain if found dest client then emit event to send packet
	_, found = k.clientKeeper.GetClientState(ctx, packet.GetDstChain())
	if packet.GetDstChain() != chainName && found {
		k.SetPacketCommitment(ctx, packet.GetSrcChain(), packet.GetDstChain(), packet.GetSequence(), commitment)

		_ = ctx.EventManager().EmitTypedEvent(&types.EventSendPacket{
			Sequence: fmt.Sprintf("%d", packet.GetSequence()),
			SrcChain: packet.GetSrcChain(),
			DstChain: packet.GetDstChain(),
			Packet:   msg.Packet,
		})
	}

	k.Logger(ctx).Info("packet recv", "packet", fmt.Sprintf("%v : %v : %v", packet.GetSrcChain(), packet.GetDstChain(), packet.GetSequence()))
	return nil
}

// WriteAcknowledgement writes the packet execution acknowledgement to the state,
// which will be verified by the counterparty chain using AcknowledgePacket.
//
// CONTRACT:
//
// 1) For synchronous execution, this function is be called in the XIBC handler .
// For async handling, it needs to be called directly by the module which originally
// processed the packet.
//
// 2) Assumes that packet receipt has been written.
// previously by RecvPacket.
func (k Keeper) WriteAcknowledgement(
	ctx sdk.Context,
	packet exported.PacketI,
	ack types.Acknowledgement,
) error {
	ackBz, err := ack.ABIPack()
	if err != nil {
		return sdkerrors.Wrapf(types.ErrInvalidAcknowledgement, "pack ack failed")
	}
	if len(ackBz) == 0 {
		return sdkerrors.Wrap(types.ErrInvalidAcknowledgement, "acknowledgement cannot be empty")
	}

	// NOTE: XIBC app modules might have written the acknowledgement synchronously on
	// the OnRecvPacket callback so we need to check if the acknowledgement is already
	// set on the store and return an error if so.
	if k.HasPacketAcknowledgement(
		ctx,
		packet.GetSrcChain(),
		packet.GetDstChain(),
		packet.GetSequence(),
	) {
		return types.ErrAcknowledgementExists
	}

	if _, found := k.clientKeeper.GetClientState(ctx, packet.GetSrcChain()); !found {
		return clienttypes.ErrClientNotFound
	}

	// set the acknowledgement so that it can be verified on the other side
	k.SetPacketAcknowledgement(
		ctx,
		packet.GetSrcChain(),
		packet.GetDstChain(),
		packet.GetSequence(),
		types.CommitAcknowledgement(ackBz),
	)

	// log that a packet acknowledgement has been written
	k.Logger(ctx).Info("acknowledged written", "packet", fmt.Sprintf("%v", packet))

	packetBytes, err := packet.ABIPack()
	if err != nil {
		return sdkerrors.Wrapf(types.ErrABIPack, "SendPacket error , err: %s", err)
	}

	_ = ctx.EventManager().EmitTypedEvent(
		&types.EventWriteAck{
			SrcChain: packet.GetSrcChain(),
			DstChain: packet.GetDstChain(),
			Sequence: strconv.FormatUint(packet.GetSequence(), 10),
			Packet:   packetBytes,
			Ack:      ackBz,
			Code:     ack.Code,
		},
	)

	return nil
}

// AcknowledgePacket is called by a module to process the acknowledgement of a
// packet previously sent by the calling module on a port to a counterparty
// module on the counterparty chain. Its intended usage is within the ante
// handler. AcknowledgePacket will clean up the packet commitment,
// which is no longer necessary since the packet has been received and acted upon.
func (k Keeper) AcknowledgePacket(
	ctx sdk.Context,
	msg *types.MsgAcknowledgement,
) error {
	var packet types.Packet
	if err := packet.ABIDecode(msg.Packet); err != nil {
		return sdkerrors.Wrapf(types.ErrDecodeAbi, "AcknowledgePacket failed ,err %s", err)
	}
	if err := k.ValidatePacket(ctx, &packet); err != nil {
		return sdkerrors.Wrap(err, "AcknowledgePacket failed basic validation")
	}
	commitment := k.GetPacketCommitment(
		ctx,
		packet.GetSrcChain(),
		packet.GetDstChain(),
		packet.GetSequence(),
	)

	packetCommitment, err := types.CommitPacket(&packet)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrABIPack, "AcknowledgePacket failed ,err %s", err)
	}
	// verify we sent the packet and haven't cleared it out yet
	if !bytes.Equal(commitment, packetCommitment) {
		return sdkerrors.Wrapf(
			types.ErrInvalidPacket,
			"commitment bytes are not equal: got (%v), expected (%v)",
			packetCommitment, commitment,
		)
	}

	fromChain := packet.GetDstChain()
	clientState, found := k.clientKeeper.GetClientState(ctx, fromChain)
	if !found {
		return sdkerrors.Wrap(clienttypes.ErrClientNotFound, fromChain)
	}

	ackCommitment := types.CommitAcknowledgement(msg.Acknowledgement)

	// use signer as tss client proof
	proof := msg.ProofAcked
	if clientState.ClientType() == exported.TSS {
		proof = []byte(msg.Signer)
	}

	if err := clientState.VerifyPacketAcknowledgement(
		ctx,
		k.clientKeeper.ClientStore(ctx, fromChain),
		k.cdc,
		msg.ProofHeight,
		proof,
		packet.GetSrcChain(),
		packet.GetDstChain(),
		packet.GetSequence(),
		ackCommitment,
	); err != nil {
		return sdkerrors.Wrapf(err, "failed packet acknowledgement verification for client (%s)", fromChain)
	}

	// Delete packet commitment, since the packet has been acknowledged, the commitement is no longer necessary
	k.deletePacketCommitment(ctx, packet.GetSrcChain(), packet.GetDstChain(), packet.GetSequence())

	// log that a packet has been acknowledged
	k.Logger(ctx).Info("packet acknowledged", "packet", fmt.Sprintf("%v", packet))

	_ = ctx.EventManager().EmitTypedEvent(
		&types.EventAcknowledgePacket{
			SrcChain: packet.GetSrcChain(),
			DstChain: packet.GetDstChain(),
			Sequence: strconv.FormatUint(packet.GetSequence(), 10),
			Packet:   msg.Packet,
			Ack:      msg.Acknowledgement,
		},
	)

	if packet.GetSrcChain() != k.clientKeeper.GetChainName(ctx) {
		if _, found = k.clientKeeper.GetClientState(ctx, packet.GetSrcChain()); !found {
			return sdkerrors.Wrap(clienttypes.ErrClientNotFound, fromChain)
		}
		// set the acknowledgement so that it can be verified on the other side
		k.SetPacketAcknowledgement(
			ctx, packet.GetSrcChain(), packet.GetDstChain(), packet.GetSequence(),
			ackCommitment,
		)

		_ = ctx.EventManager().EmitTypedEvent(&types.EventWriteAck{
			Sequence: fmt.Sprintf("%d", packet.GetSequence()),
			SrcChain: packet.GetSrcChain(),
			DstChain: packet.GetDstChain(),
			Packet:   msg.Packet,
			Ack:      msg.Acknowledgement,
		})
	}

	k.Logger(ctx).Info("acknowledge packet ", "packet", fmt.Sprintf("%v : %v : %v", packet.GetSrcChain(), packet.GetDstChain(), packet.GetSequence()))
	return nil
}
