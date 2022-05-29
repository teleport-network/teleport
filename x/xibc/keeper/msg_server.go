package keeper

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	packetcontract "github.com/teleport-network/teleport/syscontracts/xibc_packet"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
)

var _ clienttypes.MsgServer = Keeper{}
var _ packettypes.MsgServer = Keeper{}

// UpdateClient defines a rpc handler method for MsgUpdateClient.
func (k Keeper) UpdateClient(goCtx context.Context, msg *clienttypes.MsgUpdateClient) (*clienttypes.MsgUpdateClientResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	header, err := clienttypes.UnpackHeader(msg.Header)
	if err != nil {
		return nil, err
	}

	// Verify that the account has permission to update the client
	if !k.ClientKeeper.AuthRelayer(ctx, msg.ChainName, msg.Signer) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "relayer: %s", msg.Signer)
	}

	clientState, found := k.ClientKeeper.GetClientState(ctx, msg.ChainName)
	if !found {
		return nil, sdkerrors.Wrapf(clienttypes.ErrClientNotFound, "client state not found %s", msg.ChainName)
	}
	if err = clientState.CheckMsg(msg); err != nil {
		return nil, err
	}

	if err = k.ClientKeeper.UpdateClient(ctx, msg.ChainName, header); err != nil {
		return nil, err
	}

	return &clienttypes.MsgUpdateClientResponse{}, nil
}

// RecvPacket defines a rpc handler method for MsgRecvPacket.
func (k Keeper) RecvPacket(goCtx context.Context, msg *packettypes.MsgRecvPacket) (*packettypes.MsgRecvPacketResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.PacketKeeper.RecvPacket(ctx, msg); err != nil {
		return nil, sdkerrors.Wrap(err, "receive packet verification failed")
	}

	cctx, write := ctx.CacheContext()
	// todo relay chain check
	var packet packettypes.Packet
	err := packet.ABIDecode(msg.Packet)
	if err != nil {
		return nil, sdkerrors.Wrapf(packettypes.ErrAbiPack, "RecvPacket failed,decode packet err : %s", err)
	}

	relayer, found := k.ClientKeeper.GetRelayerAddressOnOtherChain(ctx, packet.SourceChain, msg.Signer)
	if !found {
		return nil, sdkerrors.Wrapf(packettypes.ErrRelayerNotFound, "relayer on source chain not found")
	}

	res, err := k.PacketKeeper.CallPacket(ctx, "onRecvPacket", packet.ToWPacket())
	// call packet to onRecvPacket
	if err != nil {
		// Write ErrAck
		errAckBz, err := packettypes.NewAcknowledgement(1, []byte{}, "receive packet callback failed", relayer, packet.FeeOption).ABIPack()
		if err != nil {
			return nil, sdkerrors.Wrapf(packettypes.ErrInvalidAcknowledgement, "pack ack failed")
		}
		if err := k.PacketKeeper.WriteAcknowledgement(ctx, &packet, errAckBz); err != nil {
			return nil, err
		}
		return &packettypes.MsgRecvPacketResponse{}, nil
	}

	var result packettypes.Result
	errDecodeResult := packetcontract.PacketContract.ABI.UnpackIntoInterface(&result, "onRecvPacket", res.Ret)
	if errDecodeResult != nil {
		return nil, sdkerrors.Wrapf(packettypes.ErrAbiPack, "RecvPacket failed,decode result err : %s", errDecodeResult)
	}
	ackBz, err := packettypes.NewAcknowledgement(result.Code, result.Result, result.Message, relayer, packet.FeeOption).ABIPack()
	if err != nil {
		return nil, sdkerrors.Wrapf(packettypes.ErrInvalidAcknowledgement, "pack ack failed")
	}
	if err := k.PacketKeeper.WriteAcknowledgement(ctx, &packet, ackBz); err != nil {
		return nil, err
	}

	write()
	ctx.EventManager().EmitEvents(cctx.EventManager().Events())

	return &packettypes.MsgRecvPacketResponse{}, nil
}

// Acknowledgement defines a rpc handler method for MsgAcknowledgement.
func (k Keeper) Acknowledgement(goCtx context.Context, msg *packettypes.MsgAcknowledgement) (*packettypes.MsgAcknowledgementResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.PacketKeeper.AcknowledgePacket(ctx, msg); err != nil {
		return nil, sdkerrors.Wrap(err, "acknowledge packet verification failed")
	}

	var packet packettypes.Packet
	err := packet.ABIDecode(msg.Packet)
	if err != nil {
		return nil, sdkerrors.Wrapf(packettypes.ErrDecodeAbi, "Acknowledgement failed,decode packet err : %v", err)
	}

	var ack packettypes.Acknowledgement
	if err := ack.ABIDecode(msg.Acknowledgement); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "decode acknowledgement bytes failed: %v", err)
	}
	if len(ack.String()) == 0 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "decode acknowledgement bytes failed")
	}

	// todo ?
	//success := ack.Result != nil && len(ack.Result) > 0
	success := ack.Code == 0
	if success {
		// set sequence in packet contract
		if _, err := k.PacketKeeper.CallPacket(
			ctx,
			"setAckStatus",
			packet.GetDestChain(),
			packet.Sequence,
			uint8(1),
		); err != nil {
			return nil, err
		}
	} else {
		// set sequence in packet contract
		if _, err := k.PacketKeeper.CallPacket(
			ctx,
			"setAckStatus",
			packet.GetDestChain(),
			packet.Sequence,
			uint8(2),
		); err != nil {
			return nil, err
		}
	}

	relayer, found := k.ClientKeeper.GetRelayerAddressOnTeleport(ctx, packet.GetDestChain(), ack.Relayer)
	if !found {
		return nil, sdkerrors.Wrapf(packettypes.ErrRelayerNotFound, "relayer on source chain not found")
	}

	relayerAddr, err := sdk.AccAddressFromBech32(relayer)
	if err != nil {
		return nil, sdkerrors.Wrapf(packettypes.ErrInvalidRelayer, "convert relayer address error")
	}

	if _, err := k.PacketKeeper.CallPacket(
		ctx,
		"sendPacketFeeToRelayer",
		packet.GetDestChain(),
		packet.Sequence,
		common.BytesToAddress(relayerAddr),
	); err != nil {
		return nil, err
	}

	// OnAcknowledgementPacket
	if _, err := k.PacketKeeper.CallPacket(
		ctx,
		"OnAcknowledgePacket",
		packet.ToWPacket(),
		ack,
	); err != nil {
		return nil, err
	}

	return &packettypes.MsgAcknowledgementResponse{}, nil
}
