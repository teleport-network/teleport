package tools

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/bitdao-io/bitnetwork/tools/extract"
	"github.com/bitdao-io/bitnetwork/tools/rb"

	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewExtractLatestVersionAppCmd(storeKeys map[string]*sdk.KVStoreKey) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extract-app",
		Short: "extract the latest height application data from the original db, and store it into the target DB",
		RunE: func(cmd *cobra.Command, args []string) error {
			targetHome, _ := cmd.Flags().GetString("target-home")
			if len(targetHome) == 0 {
				return errors.New("missing target-home")
			}
			cfg := server.GetServerContextFromCmd(cmd).Config
			home := cfg.RootDir
			return extract.CopyApplication(home, targetHome, storeKeys)
		},
	}
	cmd.Flags().String("target-home", "", "The home path of the target db")
	return cmd
}

func NewExtractLatestVersionBlockCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extract-block",
		Short: "extract the latest height block store data from the original db, and store it into the target DB",
		Long: `
the extractor is performed to extract the latest height of data from the databases created in tendermint layer,
and store it into the target database. It is always used to prune the historic data from tendermint databases to recover the db performance 
which was badly affected by the large data size. 
Note: the DB will only keep the recentest height of data, thus the earlier blocks, transactions 
will not be able to be queried via rpc/p2p server built upon it.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			targetHome, _ := cmd.Flags().GetString("target-home")
			if len(targetHome) == 0 {
				return errors.New("missing target-home")
			}
			cfg := server.GetServerContextFromCmd(cmd).Config
			home := cfg.RootDir
			extract.CopyBlockStore(home, targetHome)
			return nil
		},
	}
	cmd.Flags().String("target-home", "", "The home path of the target db")
	return cmd
}

func NewRollbackAnyCmd(storeKeys map[string]*sdk.KVStoreKey) *cobra.Command {
	var height int64
	cmd := &cobra.Command{
		Use:   "rollback-any",
		Short: "rollback cosmos-sdk and tendermint state to specified height",
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
			config := ctx.Config
			home := config.RootDir

			if err := rb.RollbackBlockAndState(home, height, config); err != nil {
				return err
			}

			if err := rb.RollbackApp(home, storeKeys, height); err != nil {
				return err
			}

			return nil
		},
	}
	cmd.Flags().Int64Var(&height, "height", -1, "the target height we want to rollback")
	return cmd
}
