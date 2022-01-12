package xibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	client "github.com/teleport-network/teleport/x/xibc/core/client"
	packet "github.com/teleport-network/teleport/x/xibc/core/packet"
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
