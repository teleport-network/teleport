package keeper

import (
	"context"

	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
)

// ClientState implements the XIBC QueryServer interface
func (q Keeper) ClientState(c context.Context, req *clienttypes.QueryClientStateRequest) (*clienttypes.QueryClientStateResponse, error) {
	return q.ClientKeeper.ClientState(c, req)
}

// ClientStates implements the XIBC QueryServer interface
func (q Keeper) ClientStates(c context.Context, req *clienttypes.QueryClientStatesRequest) (*clienttypes.QueryClientStatesResponse, error) {
	return q.ClientKeeper.ClientStates(c, req)
}

// ConsensusState implements the XIBC QueryServer interface
func (q Keeper) ConsensusState(c context.Context, req *clienttypes.QueryConsensusStateRequest) (*clienttypes.QueryConsensusStateResponse, error) {
	return q.ClientKeeper.ConsensusState(c, req)
}

// ConsensusStates implements the XIBC QueryServer interface
func (q Keeper) ConsensusStates(c context.Context, req *clienttypes.QueryConsensusStatesRequest) (*clienttypes.QueryConsensusStatesResponse, error) {
	return q.ClientKeeper.ConsensusStates(c, req)
}

// PacketCommitment implements the XIBC QueryServer interface
func (q Keeper) PacketCommitment(c context.Context, req *packettypes.QueryPacketCommitmentRequest) (*packettypes.QueryPacketCommitmentResponse, error) {
	return q.PacketKeeper.PacketCommitment(c, req)
}

// PacketCommitments implements the XIBC QueryServer interface
func (q Keeper) PacketCommitments(c context.Context, req *packettypes.QueryPacketCommitmentsRequest) (*packettypes.QueryPacketCommitmentsResponse, error) {
	return q.PacketKeeper.PacketCommitments(c, req)
}

// PacketReceipt implements the XIBC QueryServer interface
func (q Keeper) PacketReceipt(c context.Context, req *packettypes.QueryPacketReceiptRequest) (*packettypes.QueryPacketReceiptResponse, error) {
	return q.PacketKeeper.PacketReceipt(c, req)
}

// PacketAcknowledgement implements the XIBC QueryServer interface
func (q Keeper) PacketAcknowledgement(c context.Context, req *packettypes.QueryPacketAcknowledgementRequest) (*packettypes.QueryPacketAcknowledgementResponse, error) {
	return q.PacketKeeper.PacketAcknowledgement(c, req)
}

// PacketAcknowledgements implements the XIBC QueryServer interface
func (q Keeper) PacketAcknowledgements(c context.Context, req *packettypes.QueryPacketAcknowledgementsRequest) (*packettypes.QueryPacketAcknowledgementsResponse, error) {
	return q.PacketKeeper.PacketAcknowledgements(c, req)
}

// UnreceivedPackets implements the XIBC QueryServer interface
func (q Keeper) UnreceivedPackets(c context.Context, req *packettypes.QueryUnreceivedPacketsRequest) (*packettypes.QueryUnreceivedPacketsResponse, error) {
	return q.PacketKeeper.UnreceivedPackets(c, req)
}

// UnreceivedAcks implements the XIBC QueryServer interface
func (q Keeper) UnreceivedAcks(c context.Context, req *packettypes.QueryUnreceivedAcksRequest) (*packettypes.QueryUnreceivedAcksResponse, error) {
	return q.PacketKeeper.UnreceivedAcks(c, req)
}

func (q Keeper) Relayers(c context.Context, req *clienttypes.QueryRelayersRequest) (*clienttypes.QueryRelayersResponse, error) {
	return q.ClientKeeper.Relayers(c, req)
}
