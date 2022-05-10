package tmservice

import (
	"context"
	gogogrpc "github.com/gogo/protobuf/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/tendermint/tendermint/abci/types"
)

// This is the struct that we will implement all the handlers on.
type queryServer struct {
	clientCtx         client.Context
	interfaceRegistry codectypes.InterfaceRegistry
}

var _ ServiceServer = queryServer{}

// NewQueryServer creates a new tendermint query server. It is extended from cosmos base tmservice
func NewQueryServer(clientCtx client.Context, interfaceRegistry codectypes.InterfaceRegistry) ServiceServer {
	return queryServer{
		clientCtx:         clientCtx,
		interfaceRegistry: interfaceRegistry,
	}
}

func (s queryServer) GetBlockResults(ctx context.Context, req *GetBlockResultsRequest) (*GetBlockResultsResponse, error) {
	node, err := s.clientCtx.GetNode()
	if err != nil {
		return nil, err
	}
	blockResults, err := node.BlockResults(ctx, &req.Height)
	if err != nil {
		return nil, err
	}
	resp := GetBlockResultsResponse{
		Height:                blockResults.Height,
		TxsResults:            blockResults.TxsResults,
		BeginBlockEvents:      make([]*types.Event, len(blockResults.BeginBlockEvents)),
		EndBlockEvents:        make([]*types.Event, len(blockResults.EndBlockEvents)),
		ValidatorUpdates:      make([]*types.ValidatorUpdate, len(blockResults.ValidatorUpdates)),
		ConsensusParamUpdates: blockResults.ConsensusParamUpdates,
	}
	for i, bbe := range blockResults.BeginBlockEvents {
		resp.BeginBlockEvents[i] = &bbe
	}
	for i, ebe := range blockResults.EndBlockEvents {
		resp.EndBlockEvents[i] = &ebe
	}
	for i, vu := range blockResults.ValidatorUpdates {
		resp.ValidatorUpdates[i] = &vu
	}
	return &resp, nil
}

// RegisterTendermintService registers the tendermint queries on the gRPC router.
func RegisterTendermintService(
	qrt gogogrpc.Server,
	clientCtx client.Context,
	interfaceRegistry codectypes.InterfaceRegistry,
) {
	RegisterServiceServer(
		qrt,
		NewQueryServer(clientCtx, interfaceRegistry),
	)
}

// RegisterGRPCGatewayRoutes mounts the tendermint service's GRPC-gateway routes on the
// given Mux.
func RegisterGRPCGatewayRoutes(clientConn gogogrpc.ClientConn, mux *runtime.ServeMux) {
	RegisterServiceHandlerClient(context.Background(), mux, NewServiceClient(clientConn))
}
