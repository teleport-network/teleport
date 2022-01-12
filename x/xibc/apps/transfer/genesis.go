package transfer

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"

	"github.com/teleport-network/teleport/x/xibc/apps/transfer/types"
)

// InitGenesis import module genesis
func InitGenesis(ctx sdk.Context, accountKeeper authkeeper.AccountKeeper) {
	// ensure xibc transfer module account is set on genesis
	if acc := accountKeeper.GetModuleAccount(ctx, types.ModuleName); acc == nil {
		panic("the xibc transfer module account has not been set")
	}
}
