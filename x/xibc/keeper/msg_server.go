package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/ethereum/go-ethereum/common"

	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
	routingtypes "github.com/teleport-network/teleport/x/xibc/core/routing/types"
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

	if msg.Packet.GetDestChain() == k.ClientKeeper.GetChainName(cctx) {
		relayer, found := k.ClientKeeper.GetRelayerAddressOnOtherChain(ctx, msg.Packet.SourceChain, msg.Signer)
		if !found {
			return nil, sdkerrors.Wrapf(packettypes.ErrRelayerNotFound, "relayer on source chain not found")
		}

		var results [][]byte
		for i, port := range msg.Packet.Ports {
			// Retrieve callbacks from router
			cbs, ok := k.RoutingKeeper.Router.GetRoute(port)
			if !ok {
				return nil, sdkerrors.Wrapf(routingtypes.ErrInvalidRoute, "route not found to module: %s", port)
			}

			// Perform application logic callback
			_, result, err := cbs.OnRecvPacket(cctx, msg.Packet.GetDataList()[i])
			if err != nil {
				return nil, sdkerrors.Wrap(err, "receive packet callback failed")
			}

			if len(result.Result) == 0 {
				errAckBz, err := packettypes.NewErrorAcknowledgement(result.Message, relayer).GetBytes()
				if err != nil {
					return nil, sdkerrors.Wrapf(packettypes.ErrInvalidAcknowledgement, "pack ack failed")
				}
				if err := k.PacketKeeper.WriteAcknowledgement(
					ctx,
					msg.Packet,
					errAckBz,
				); err != nil {
					return nil, err
				}

				return &packettypes.MsgRecvPacketResponse{}, nil
			}

			results = append(results, result.Result)
		}
		ackBz, err := packettypes.NewResultAcknowledgement(results, relayer).GetBytes()
		if err != nil {
			return nil, sdkerrors.Wrapf(packettypes.ErrInvalidAcknowledgement, "pack ack failed")
		}
		if err := k.PacketKeeper.WriteAcknowledgement(ctx, msg.Packet, ackBz); err != nil {
			return nil, err
		}
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

	chainName := k.ClientKeeper.GetChainName(ctx)
	if msg.Packet.GetRelayChain() == chainName {
		return &packettypes.MsgAcknowledgementResponse{}, nil
	}

	var ack packettypes.Acknowledgement
	if err := ack.DecodeBytes(msg.Acknowledgement); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "decode acknowledgement bytes failed: %v", err)
	}
	if len(ack.String()) == 0 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "decode acknowledgement bytes failed")
	}

	success := ack.Results != nil && len(ack.Results) > 0
	if success {
		// set sequence in packet contract
		if _, err := k.PacketKeeper.CallPacket(
			ctx,
			"setAckStatus",
			msg.Packet.GetSourceChain(),
			msg.Packet.GetDestChain(),
			msg.Packet.Sequence,
			uint8(1),
		); err != nil {
			return nil, err
		}
	} else {
		// set sequence in packet contract
		if _, err := k.PacketKeeper.CallPacket(
			ctx,
			"setAckStatus",
			msg.Packet.GetSourceChain(),
			msg.Packet.GetDestChain(),
			msg.Packet.Sequence,
			uint8(2),
		); err != nil {
			return nil, err
		}
	}

	relayer, found := k.ClientKeeper.GetRelayerAddressOnTeleport(ctx, msg.Packet.GetDestChain(), ack.Relayer)
	if !found {
		return nil, sdkerrors.Wrapf(packettypes.ErrRelayerNotFound, "relayer on source chain not found")
	}

	if _, err := k.PacketKeeper.CallPacket(
		ctx,
		"sendPacketFeeToRelayer",
		msg.Packet.GetSourceChain(),
		msg.Packet.GetDestChain(),
		msg.Packet.Sequence,
		common.HexToAddress(relayer),
	); err != nil {
		return nil, err
	}

	for i, port := range msg.Packet.Ports {
		cbs, ok := k.RoutingKeeper.Router.GetRoute(port)
		if !ok {
			return nil, sdkerrors.Wrapf(routingtypes.ErrInvalidRoute, "route not found to module: %s", port)
		}
		if success {
			if _, err := cbs.OnAcknowledgementPacket(ctx, msg.Packet.GetDataList()[i], ack.Results[i]); err != nil {
				return nil, sdkerrors.Wrap(err, "acknowledge packet callback failed")
			}
		} else {
			if _, err := cbs.OnAcknowledgementPacket(ctx, msg.Packet.GetDataList()[i], []byte{}); err != nil {
				return nil, sdkerrors.Wrap(err, "acknowledge packet callback failed")
			}
		}
	}

	return &packettypes.MsgAcknowledgementResponse{}, nil
}
