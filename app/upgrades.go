package app

import (
	"fmt"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
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

			// Delete deprecated xibc contracts account, code and state

			// Set new code

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
