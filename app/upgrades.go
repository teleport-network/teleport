package app

import (
	"fmt"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/ethereum/go-ethereum/common"

	syscontracts "github.com/teleport-network/teleport/syscontracts"
	agentcontract "github.com/teleport-network/teleport/syscontracts/xibc_agent"
	endpointcontract "github.com/teleport-network/teleport/syscontracts/xibc_endpoint"
	packetcontract "github.com/teleport-network/teleport/syscontracts/xibc_packet"
	"github.com/teleport-network/teleport/x/xibc"
	xibchost "github.com/teleport-network/teleport/x/xibc/core/host"
)

// nolint: ignore
func (app *Teleport) registerUpgradeHandlers() {
	// v0.2 upgrade handler
	app.UpgradeKeeper.SetUpgradeHandler(
		"v0.2",
		func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
			// Refs:
			// - https://docs.cosmos.network/master/building-modules/upgrade.html#registering-migrations
			// - https://docs.cosmos.network/master/migrations/chain-upgrade-guide-044.html#chain-upgrade

			// xibc core contracts
			DeprecatedPacketContractAddress := "0x0000000000000000000000000000000020000001"

			// xibc app contracts
			DeprecatedTransferContractAddress := "0x0000000000000000000000000000000030000001"
			DeprecatedRCCContractAddress := "0x0000000000000000000000000000000030000002"
			DeprecatedMultiCallContractAddress := "0x0000000000000000000000000000000030000003"

			// app contracts
			DeprecatedAgentContractAddress := "0x0000000000000000000000000000000040000001"

			// Delete deprecated xibc contracts account, code and state
			_ = app.EvmKeeper.DeleteAccount(ctx, common.HexToAddress(DeprecatedTransferContractAddress))
			_ = app.EvmKeeper.DeleteAccount(ctx, common.HexToAddress(DeprecatedRCCContractAddress))
			_ = app.EvmKeeper.DeleteAccount(ctx, common.HexToAddress(DeprecatedMultiCallContractAddress))
			_ = app.EvmKeeper.DeleteAccount(ctx, common.HexToAddress(DeprecatedAgentContractAddress))
			_ = app.EvmKeeper.DeleteAccount(ctx, common.HexToAddress(DeprecatedPacketContractAddress))

			// Set new code
			app.SetEVMCode(ctx, common.HexToAddress(syscontracts.AgentContractAddress), agentcontract.AgentContract.Bin)
			app.SetEVMCode(ctx, common.HexToAddress(syscontracts.PacketContractAddress), packetcontract.PacketContract.Bin)
			app.SetEVMCode(ctx, common.HexToAddress(syscontracts.EndpointContractAddress), endpointcontract.EndpointContract.Bin)
			app.SetEVMCode(ctx, common.HexToAddress(syscontracts.ExecuteContractAddress), endpointcontract.ExecuteContract.Bin)

			xibc.ResetStates(ctx, app.GetKey(xibchost.StoreKey), *app.XIBCKeeper)

			return app.mm.RunMigrations(ctx, app.configurator, vm)
		})

	// When a planned update height is reached, the old binary will panic
	// writing on disk the height and name of the update that triggered it
	// This will read that value, and execute the preparations for the upgrade.
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Errorf("failed to read upgrade info from disk: %w", err))
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	var storeUpgrades *storetypes.StoreUpgrades

	switch upgradeInfo.Name {
	case "v0.2":
		// no store upgrades in v0.2
	}

	if storeUpgrades != nil {
		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}
