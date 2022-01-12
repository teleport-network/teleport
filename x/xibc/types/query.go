package types

import (
	"github.com/gogo/protobuf/grpc"

	client "github.com/teleport-network/teleport/x/xibc/core/client/module"
	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	packet "github.com/teleport-network/teleport/x/xibc/core/packet/module"
	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
)

// QueryServer defines the XIBC interfaces that the gRPC query server must implement
type QueryServer interface {
	clienttypes.QueryServer
	packettypes.QueryServer
}

// RegisterQueryService registers each individual XIBC submodule query service
func RegisterQueryService(server grpc.Server, queryService QueryServer) {
	client.RegisterQueryService(server, queryService)
	packet.RegisterQueryService(server, queryService)
}
