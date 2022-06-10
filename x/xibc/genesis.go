package xibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	client "github.com/teleport-network/teleport/x/xibc/core/client"
	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	"github.com/teleport-network/teleport/x/xibc/core/host"
	packet "github.com/teleport-network/teleport/x/xibc/core/packet"
	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
	"github.com/teleport-network/teleport/x/xibc/keeper"
	"github.com/teleport-network/teleport/x/xibc/types"
)

// InitGenesis initializes the xibc state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, createLocalhost bool, gs *types.GenesisState) {
	client.InitGenesis(ctx, k.ClientKeeper, gs.ClientGenesis)
	packet.InitGenesis(ctx, k.PacketKeeper, gs.PacketGenesis)
}

// ExportGenesis returns the xibc exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		ClientGenesis: client.ExportGenesis(ctx, k.ClientKeeper),
		PacketGenesis: packet.ExportGenesis(ctx, k.PacketKeeper),
	}
}

// ResetStates reset all xibc state.
// storeKey must be xibc key.
func ResetStates(ctx sdk.Context, storeKey sdk.StoreKey, k keeper.Keeper) {
	if storeKey.Name() != host.StoreKey {
		panic("storeKey must be xibc key")
	}

	nativeChainName := k.ClientKeeper.GetChainName(ctx)

	store := ctx.KVStore(storeKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}

	clientGenesis := clienttypes.DefaultGenesisState()
	clientGenesis.NativeChainName = nativeChainName

	packetGenesis := packettypes.DefaultGenesisState()

	client.InitGenesis(ctx, k.ClientKeeper, clientGenesis)
	packet.InitGenesis(ctx, k.PacketKeeper, packetGenesis)
}
