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
	var packet packettypes.Packet
	if err := packet.ABIDecode(msg.Packet); err != nil {
		return nil, sdkerrors.Wrapf(packettypes.ErrABIPack, "RecvPacket failed, decode packet err: %s", err)
	}

	relayer, found := k.ClientKeeper.GetRelayerAddressOnOtherChain(ctx, packet.SrcChain, msg.Signer)
	if !found {
		return nil, sdkerrors.Wrapf(packettypes.ErrRelayerNotFound, "relayer on source chain not found")
	}

	if packet.GetDstChain() == k.ClientKeeper.GetChainName(cctx) {
		// call packet onRecvPacket
		res, err := k.PacketKeeper.CallPacket(ctx, "onRecvPacket", packet)
		if err != nil {
			// Write ErrAck
			return nil, k.PacketKeeper.WriteAcknowledgement(
				ctx, &packet, packettypes.NewAcknowledgement(1, []byte{}, "receive packet callback failed", relayer, packet.FeeOption),
			)
		}
		// call onRecvPacket end then get the result to write the ack
		var result packettypes.Result
		if err := packetcontract.PacketContract.ABI.UnpackIntoInterface(&result, "onRecvPacket", res.Ret); err != nil {
			return nil, sdkerrors.Wrapf(packettypes.ErrABIPack, "recv packet failed, decode result err: %s", err)
		}
		if err := k.PacketKeeper.WriteAcknowledgement(
			ctx, &packet, packettypes.NewAcknowledgement(result.Code, result.Result, result.Message, relayer, packet.FeeOption),
		); err != nil {
			return nil, err
		}
	} else if _, found := k.ClientKeeper.GetClientState(ctx, packet.GetDstChain()); !found {
		// Write ErrAck
		return nil, k.PacketKeeper.WriteAcknowledgement(
			ctx, &packet, packettypes.NewAcknowledgement(1, []byte{}, "dstChain not found", relayer, packet.FeeOption),
		)
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
	if err := packet.ABIDecode(msg.Packet); err != nil {
		return nil, sdkerrors.Wrapf(packettypes.ErrDecodeAbi, "decode packet bytes failed: %v", err)
	}

	var ack packettypes.Acknowledgement
	if err := ack.ABIDecode(msg.Acknowledgement); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "decode acknowledgement bytes failed: %v", err)
	}
	if len(ack.String()) == 0 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "decode acknowledgement bytes failed")
	}

	if packet.GetSrcChain() == k.ClientKeeper.GetChainName(ctx) {
		success := ack.Code == 0
		if success {
			// set sequence in packet contract
			if _, err := k.PacketKeeper.CallPacket(
				ctx,
				"setAckStatus",
				packet.GetDstChain(),
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
				packet.GetDstChain(),
				packet.Sequence,
				uint8(2),
			); err != nil {
				return nil, err
			}
		}

		relayer, found := k.ClientKeeper.GetRelayerAddressOnTeleport(ctx, packet.GetDstChain(), ack.Relayer)
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
			packet.GetDstChain(),
			packet.Sequence,
			common.BytesToAddress(relayerAddr),
		); err != nil {
			return nil, err
		}

		// OnAcknowledgementPacket
		if _, err := k.PacketKeeper.CallPacket(ctx, "OnAcknowledgePacket", packet, ack); err != nil {
			return nil, err
		}
	}

	return &packettypes.MsgAcknowledgementResponse{}, nil
}
