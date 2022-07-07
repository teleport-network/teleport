package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	sdk "github.com/cosmos/cosmos-sdk/types"

	app "github.com/bitdao-io/bitnetwork/app"
	"github.com/bitdao-io/bitnetwork/cmd/bitnetworkd/cmd"
	cmdCfg "github.com/bitdao-io/bitnetwork/cmd/config"
)

func main() {
	setupConfig()
	rootCmd, _ := cmd.NewRootCmd()
	cmdCfg.RegisterDenoms()

	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}

func setupConfig() {
	// set the address prefixes
	config := sdk.GetConfig()
	cmdCfg.SetBech32Prefixes(config)
	cmdCfg.SetBip44CoinType(config)
	config.Seal()
}
