package rcc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"

	"github.com/teleport-network/teleport/x/xibc/apps/rcc/types"
)

// InitGenesis import module genesis
func InitGenesis(ctx sdk.Context, accountKeeper authkeeper.AccountKeeper) {
	// ensure xibc rcc module account is set on genesis
	if acc := accountKeeper.GetModuleAccount(ctx, types.ModuleName); acc == nil {
		panic("the xibc rcc module account has not been set")
	}
}
