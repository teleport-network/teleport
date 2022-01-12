package packet

import (
	"github.com/gogo/protobuf/grpc"
	"github.com/spf13/cobra"

	"github.com/teleport-network/teleport/x/xibc/core/packet/client/cli"
	"github.com/teleport-network/teleport/x/xibc/core/packet/types"
)

// Name returns the XIBC packet name.
func Name() string {
	return types.SubModuleName
}

// GetTxCmd returns the root tx command for XIBC packet.
func GetTxCmd() *cobra.Command {
	return nil
}

// GetQueryCmd returns the root query command for XIBC packet.
func GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// RegisterQueryService registers the gRPC query service for XIBC packet.
func RegisterQueryService(server grpc.Server, queryServer types.QueryServer) {
	types.RegisterQueryServer(server, queryServer)
}
