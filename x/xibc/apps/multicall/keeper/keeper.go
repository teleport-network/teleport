package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/teleport-network/teleport/x/xibc/apps/multicall/types"
)

type Keeper struct {
	packetKeeper    types.PacketKeeper
	clientKeeper    types.ClientKeeper
	aggregateKeeper types.AggregateKeeper
}

// NewKeeper creates a new XIBC MultiCall Keeper instance
func NewKeeper(
	packetKeeper types.PacketKeeper,
	clientKeeper types.ClientKeeper,
	aggregateKeeper types.AggregateKeeper,
) Keeper {
	return Keeper{
		packetKeeper:    packetKeeper,
		clientKeeper:    clientKeeper,
		aggregateKeeper: aggregateKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
