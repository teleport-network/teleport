package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/teleport-network/teleport/x/xibc/core/packet/types"
)

// GetQueryCmd returns a root CLI command handler for all xibc/packet query commands.
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.SubModuleName,
		Short:                      "xibc packet query subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
		GetCmdQueryPacketCommitment(),
		GetCmdQueryPacketCommitments(),
		GetCmdQueryPacketReceipt(),
		GetCmdQueryPacketAcknowledgement(),
		GetCmdQueryUnreceivedPackets(),
		GetCmdQueryUnreceivedAcks(),
	)

	return queryCmd
}
