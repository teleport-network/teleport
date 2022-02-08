package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

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

	if err = k.ClientKeeper.UpdateClient(ctx, msg.ChainName, header); err != nil {
		return nil, err
	}

	return &clienttypes.MsgUpdateClientResponse{}, nil
}

// RecvPacket defines a rpc handler method for MsgRecvPacket.
func (k Keeper) RecvPacket(goCtx context.Context, msg *packettypes.MsgRecvPacket) (*packettypes.MsgRecvPacketResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.PacketKeeper.RecvPacket(ctx, msg.Packet, msg.ProofCommitment, msg.ProofHeight); err != nil {
		return nil, sdkerrors.Wrap(err, "receive packet verification failed")
	}

	cctx, write := ctx.CacheContext()

	if msg.Packet.GetDestChain() == k.ClientKeeper.GetChainName(cctx) {
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
				if err := k.PacketKeeper.WriteAcknowledgement(
					ctx,
					msg.Packet,
					packettypes.NewErrorAcknowledgement(result.Message).GetBytes(),
				); err != nil {
					return nil, err
				}

				return &packettypes.MsgRecvPacketResponse{}, nil
			}

			results = append(results, result.Result)
		}

		if err := k.PacketKeeper.WriteAcknowledgement(ctx, msg.Packet, packettypes.NewResultAcknowledgement(results).GetBytes()); err != nil {
			return nil, err
		}
	}

	write()

	return &packettypes.MsgRecvPacketResponse{}, nil
}

// Acknowledgement defines a rpc handler method for MsgAcknowledgement.
func (k Keeper) Acknowledgement(goCtx context.Context, msg *packettypes.MsgAcknowledgement) (*packettypes.MsgAcknowledgementResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.PacketKeeper.AcknowledgePacket(ctx, msg.Packet, msg.Acknowledgement, msg.ProofAcked, msg.ProofHeight); err != nil {
		return nil, sdkerrors.Wrap(err, "acknowledge packet verification failed")
	}

	var ack packettypes.Acknowledgement
	if err := ack.Unmarshal(msg.Acknowledgement); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal acknowledgement: %v", err)
	}

	success := ack.Results != nil && len(ack.Results) > 0
	for i, port := range msg.Packet.Ports {
		cbs, ok := k.RoutingKeeper.Router.GetRoute(port)
		if !ok {
			return nil, sdkerrors.Wrapf(routingtypes.ErrInvalidRoute, "route not found to module: %s", port)
		}

		if success {
			if _, err := cbs.OnAcknowledgementPacket(ctx, msg.Packet.GetDataList()[i], ack.Results[i]); err != nil {
				return nil, sdkerrors.Wrap(err, "acknowledge packet callback failed")
			}
			// set sequence in packet contract
			if _, err := k.PacketKeeper.CallPacket(ctx, "setAckStatus", msg.Packet.GetSourceChain(), msg.Packet.GetDestChain(), msg.Packet.Sequence, uint8(1)); err != nil {
				return nil, err
			}
		} else {
			if _, err := cbs.OnAcknowledgementPacket(ctx, msg.Packet.GetDataList()[i], []byte{}); err != nil {
				return nil, sdkerrors.Wrap(err, "acknowledge packet callback failed")
			}
			// set sequence in packet contract
			if _, err := k.PacketKeeper.CallPacket(ctx, "setAckStatus", msg.Packet.GetSourceChain(), msg.Packet.GetDestChain(), msg.Packet.Sequence, uint8(2)); err != nil {
				return nil, err
			}
		}
	}

	return &packettypes.MsgAcknowledgementResponse{}, nil
}
