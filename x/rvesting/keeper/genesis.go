package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitdao-io/bitnetwork/x/rvesting/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, genesisState *types.GenesisState) {
	k.SetParams(ctx, genesisState.GetParams())
	if len(genesisState.From) == 0 {
		return
	}
	from, err := sdk.AccAddressFromBech32(genesisState.From)
	if err != nil {
		panic(err)
	}
	if err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, genesisState.InitReward); err != nil {
		panic(err)
	}
}

// ExportGenesis returns a GenesisState .
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	params := k.GetParams(ctx)
	return types.NewGenesisState(params)
}
