package keeper

import (
	"bytes"
	"fmt"

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
	if packet.GetSourceChain() != k.clientKeeper.GetChainName(ctx) {
		return sdkerrors.Wrap(types.ErrInvalidPacket, "source chain of packet is not this chain")
	}

	targetChain := packet.GetDestChain()
	if len(packet.GetRelayChain()) > 0 {
		targetChain = packet.GetRelayChain()
	}

	if _, found := k.clientKeeper.GetClientState(ctx, targetChain); !found {
		return clienttypes.ErrClientNotFound
	}

	nextSequenceSend := k.GetNextSequenceSend(ctx, packet.GetSourceChain(), packet.GetDestChain())

	if packet.GetSequence() != nextSequenceSend {
		return sdkerrors.Wrapf(
			types.ErrInvalidPacket,
			"packet sequence ≠ next send sequence (%d ≠ %d)", packet.GetSequence(), nextSequenceSend,
		)
	}

	commitment, err := types.CommitPacket(packet)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrAbiPack, "SendPacket error ,err : %s", err)
	}
	nextSequenceSend++
	k.SetNextSequenceSend(ctx, packet.GetSourceChain(), packet.GetDestChain(), nextSequenceSend)

	// set sequence in packet contract
	if _, err := k.CallPacket(ctx, "setSequence", packet.GetSourceChain(), packet.GetDestChain(), nextSequenceSend); err != nil {
		return err
	}

	k.SetPacketCommitment(ctx, packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence(), commitment)

	packetData, err := packet.AbiPack()
	if err != nil {
		return sdkerrors.Wrapf(types.ErrAbiPack, "SendPacket error , err: %s", err)
	}
	_ = ctx.EventManager().EmitTypedEvent(
		&types.EventSendPacket{
			Packet: packetData,
		},
	)

	k.Logger(ctx).Info("packet sent", "packet", fmt.Sprintf("%v", packet))
	return nil
}

// RecvPacket is called by a module to receive & process an XIBC packet
// sent on the corresponding port on the counterparty chain.
func (k Keeper) RecvPacket(ctx sdk.Context, msg *types.MsgRecvPacket) error {
	var packet types.Packet
	err := packet.DecodeAbiBytes(msg.Packet)
	if err != nil && packet.Sequence == 0 {
		return sdkerrors.Wrap(err, "packet failed basic validation")
	}

	if err := k.ValidatePacket(ctx, &packet); err != nil {
		return sdkerrors.Wrap(err, "packet failed basic validation")
	}
	// check if the packet receipt has been received already
	if _, found := k.GetPacketReceipt(ctx, packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence()); found {
		return sdkerrors.Wrapf(
			types.ErrInvalidPacket,
			"packet sequence (%d) already has been received", packet.GetSequence(),
		)
	}
	chainName := k.clientKeeper.GetChainName(ctx)
	fromChain := packet.GetSourceChain()
	if packet.GetDestChain() == chainName && len(packet.GetRelayChain()) > 0 {
		fromChain = packet.GetRelayChain()
	}

	targetClient, found := k.clientKeeper.GetClientState(ctx, fromChain)
	if !found {
		return sdkerrors.Wrap(clienttypes.ErrClientNotFound, fromChain)
	}

	commitment, err := types.CommitPacket(&packet)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrAbiPack, "SendPacket error ,err : %s", err)
	}
	// use signer as tss client proof
	proof := msg.ProofCommitment
	if targetClient.ClientType() == exported.TSS {
		proof = []byte(msg.Signer)
	}
	// verify that the counterparty did commit to sending this packet
	if err := targetClient.VerifyPacketCommitment(
		ctx,
		k.clientKeeper.ClientStore(ctx, fromChain),
		k.cdc,
		msg.ProofHeight,
		proof,
		packet.GetSourceChain(),
		packet.GetDestChain(),
		packet.GetSequence(),
		commitment,
	); err != nil {
		return sdkerrors.Wrapf(err, "failed packet commitment verification for client (%s)", fromChain)
	}

	// All verification complete, update state
	k.SetPacketReceipt(ctx, packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence())

	// log that a packet has been received & executed
	k.Logger(ctx).Info("packet received", "packet", fmt.Sprintf("%v", packet))

	_ = ctx.EventManager().EmitTypedEvent(
		&types.EventRecvPacket{
			Packet: msg.Packet,
		},
	)

	if packet.GetRelayChain() == chainName {
		k.SetPacketCommitment(ctx, packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence(), commitment)
		k.SetPacketRelayer(ctx, packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence(), msg.Signer)

		_ = ctx.EventManager().EmitTypedEvent(
			&types.EventSendPacket{
				Packet: msg.Packet,
			},
		)

		// if destChain not exist, return error ack to source chain
		if _, found := k.clientKeeper.GetClientState(ctx, packet.GetDestChain()); !found {
			errAckBz, err := types.NewErrorAcknowledgement(1, []byte(""), "invalid destChain", msg.Signer).AbiPack()
			if err != nil {
				return sdkerrors.Wrapf(types.ErrInvalidAcknowledgement, "pack ack failed")
			}
			if err := k.WriteAcknowledgement(ctx, &packet, errAckBz); err != nil {
				return err
			}
		}
	}

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
	acknowledgement []byte,
) error {
	if len(acknowledgement) == 0 {
		return sdkerrors.Wrap(types.ErrInvalidAcknowledgement, "acknowledgement cannot be empty")
	}

	// NOTE: XIBC app modules might have written the acknowledgement synchronously on
	// the OnRecvPacket callback so we need to check if the acknowledgement is already
	// set on the store and return an error if so.
	if k.HasPacketAcknowledgement(
		ctx,
		packet.GetSourceChain(),
		packet.GetDestChain(),
		packet.GetSequence(),
	) {
		return types.ErrAcknowledgementExists
	}

	targetChain := packet.GetSourceChain()
	if len(packet.GetRelayChain()) > 0 && packet.GetDestChain() == k.clientKeeper.GetChainName(ctx) {
		targetChain = packet.GetRelayChain()
	}

	if _, found := k.clientKeeper.GetClientState(ctx, targetChain); !found {
		return clienttypes.ErrClientNotFound
	}

	// set the acknowledgement so that it can be verified on the other side
	k.SetPacketAcknowledgement(
		ctx,
		packet.GetSourceChain(),
		packet.GetDestChain(),
		packet.GetSequence(),
		types.CommitAcknowledgement(acknowledgement),
	)

	// log that a packet acknowledgement has been written
	k.Logger(ctx).Info("acknowledged written", "packet", fmt.Sprintf("%v", packet))

	packetData, err := packet.AbiPack()
	if err != nil {
		return sdkerrors.Wrapf(types.ErrAbiPack, "SendPacket error , err: %s", err)
	}
	_ = ctx.EventManager().EmitTypedEvent(
		&types.EventWriteAck{
			Packet: packetData,
			Ack:    acknowledgement,
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
	err := packet.DecodeAbiBytes(msg.Packet)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrDecodeAbi, "AcknowledgePacket failed ,err %s", err)
	}
	if err := k.ValidatePacket(ctx, &packet); err != nil {
		return sdkerrors.Wrap(err, "AcknowledgePacket failed basic validation")
	}
	commitment := k.GetPacketCommitment(
		ctx,
		packet.GetSourceChain(),
		packet.GetDestChain(),
		packet.GetSequence(),
	)

	packetCommitment, err := types.CommitPacket(&packet)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrAbiPack, "AcknowledgePacket failed ,err %s", err)
	}
	// verify we sent the packet and haven't cleared it out yet
	if !bytes.Equal(commitment, packetCommitment) {
		return sdkerrors.Wrapf(
			types.ErrInvalidPacket,
			"commitment bytes are not equal: got (%v), expected (%v)",
			packetCommitment, commitment,
		)
	}

	chainName := k.clientKeeper.GetChainName(ctx)
	fromChain := packet.GetDestChain()
	if packet.GetSourceChain() == chainName && len(packet.GetRelayChain()) > 0 {
		fromChain = packet.GetRelayChain()
	}

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
		packet.GetSourceChain(),
		packet.GetDestChain(),
		packet.GetSequence(),
		ackCommitment,
	); err != nil {
		return sdkerrors.Wrapf(err, "failed packet acknowledgement verification for client (%s)", fromChain)
	}

	// Delete packet commitment, since the packet has been acknowledged, the commitement is no longer necessary
	k.deletePacketCommitment(ctx, packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence())

	// log that a packet has been acknowledged
	k.Logger(ctx).Info("packet acknowledged", "packet", fmt.Sprintf("%v", packet))

	_ = ctx.EventManager().EmitTypedEvent(
		&types.EventAcknowledgePacket{
			Packet: msg.Packet,
			Ack:    msg.Acknowledgement,
		},
	)

	if packet.GetRelayChain() == chainName {
		if _, found = k.clientKeeper.GetClientState(ctx, packet.GetSourceChain()); !found {
			return sdkerrors.Wrap(clienttypes.ErrClientNotFound, fromChain)
		}
		// set the acknowledgement so that it can be verified on the other side
		k.SetPacketAcknowledgement(
			ctx,
			packet.GetSourceChain(),
			packet.GetDestChain(),
			packet.GetSequence(),
			ackCommitment,
		)

		var ack types.Acknowledgement
		if err := ack.DecodeAbiBytes(msg.Acknowledgement); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "decode acknowledgement bytes failed: %v", err)
		}

		relayerOnTeleport := k.GetPacketRelayer(
			ctx,
			packet.GetSourceChain(),
			packet.GetDestChain(),
			packet.GetSequence(),
		)

		relayer, found := k.clientKeeper.GetRelayerAddressOnOtherChain(ctx, packet.SourceChain, relayerOnTeleport)
		if !found {
			return sdkerrors.Wrapf(types.ErrRelayerNotFound, "relayer on source chain not found")
		}

		ack.Relayer = relayer
		ackBz, err := ack.AbiPack()
		if err != nil {
			return sdkerrors.Wrapf(types.ErrInvalidAcknowledgement, "pack ack failed")
		}

		k.deletePacketRelayer(ctx,
			packet.GetSourceChain(),
			packet.GetDestChain(),
			packet.GetSequence(),
		)

		_ = ctx.EventManager().EmitTypedEvent(
			&types.EventWriteAck{
				Packet: msg.Packet,
				Ack:    ackBz,
			},
		)
	}

	return nil
}
