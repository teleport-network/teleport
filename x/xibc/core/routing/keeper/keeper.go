package keeper

import (
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/teleport-network/teleport/x/xibc/core/host"
	"github.com/teleport-network/teleport/x/xibc/core/routing/types"
)

// Keeper defines the XIBC routing keeper
type Keeper struct {
	Router   *types.Router
	storeKey sdk.StoreKey
}

// NewKeeper creates a new XIBC routing Keeper instance
func NewKeeper(key sdk.StoreKey) Keeper {
	return Keeper{storeKey: key}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+host.ModuleName+"/"+types.SubModuleName)
}

// SetRouter sets the Router in XIBC Keeper and seals it. The method panics if
// there is an existing router that's already sealed.
func (k *Keeper) SetRouter(rtr *types.Router) {
	if k.Router != nil && k.Router.Sealed() {
		panic("cannot reset a sealed router")
	}
	k.Router = rtr
	k.Router.Seal()
}
