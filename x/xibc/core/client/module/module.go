package client

import (
	"github.com/gogo/protobuf/grpc"
	"github.com/spf13/cobra"

	"github.com/teleport-network/teleport/x/xibc/core/client/client/cli"
	"github.com/teleport-network/teleport/x/xibc/core/client/types"
)

// Name returns the XIBC client name.
func Name() string {
	return types.SubModuleName
}

// GetQueryCmd returns no root query command for the XIBC client.
func GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// GetTxCmd returns the root tx command for client.
func GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}

// RegisterQueryService registers the gRPC query service for XIBC client.
func RegisterQueryService(server grpc.Server, queryServer types.QueryServer) {
	types.RegisterQueryServer(server, queryServer)
}
