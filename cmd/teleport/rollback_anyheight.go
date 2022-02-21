package main

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/server"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	"github.com/spf13/cobra"
	tmcmd "github.com/tendermint/tendermint/cmd/tendermint/commands"
)

// NewRollbackAnyCmd creates a command to rollback tendermint and multistore state by one height.
func NewRollbackAnyCmd(defaultNodeHome string) *cobra.Command {
	var height int64
	cmd := &cobra.Command{
		Use:   "rollback-any",
		Short: "rollback cosmos-sdk and tendermint state by one height",
		Long: `
A state rollback is performed to recover from an incorrect application state transition,
when Tendermint has persisted an incorrect app hash and is thus unable to make
progress. Rollback overwrites a state at height n with the state at height n - 1.
The application also roll back to height n - 1. No blocks are removed, so upon
restarting Tendermint the transactions in block n will be re-executed against the
application.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := server.GetServerContextFromCmd(cmd)
			cfg := ctx.Config
			home := cfg.RootDir
			db, err := openDB(home)
			if err != nil {
				return err
			}
			// rollback tendermint state
			err = tmcmd.RollbackAnyState(ctx.Config, height)
			if err != nil {
				return fmt.Errorf("failed to rollback tendermint state: %w", err)
			}
			// rollback the multistore
			cms := rootmulti.NewStore(db)
			cms.RollbackToPreVersion(height)

			return nil
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	cmd.Flags().Int64Var(&height, "height", -1, "the height we want to rollback")
	return cmd
}
