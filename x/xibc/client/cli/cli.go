package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	xibcclient "github.com/teleport-network/teleport/x/xibc/core/client/module"
	xibchost "github.com/teleport-network/teleport/x/xibc/core/host"
	xibcpacket "github.com/teleport-network/teleport/x/xibc/core/packet/module"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	xibcTxCmd := &cobra.Command{
		Use:                        xibchost.ModuleName,
		Short:                      "XIBC transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	xibcTxCmd.AddCommand(
		xibcclient.GetTxCmd(),
	)

	return xibcTxCmd
}

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	// Group xibc queries under a subcommand
	xibcQueryCmd := &cobra.Command{
		Use:                        xibchost.ModuleName,
		Short:                      "Querying commands for the XIBC module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	xibcQueryCmd.AddCommand(
		xibcclient.GetQueryCmd(),
		xibcpacket.GetQueryCmd(),
	)

	return xibcQueryCmd
}
