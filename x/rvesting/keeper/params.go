package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitdao-io/bitnetwork/x/rvesting/types"
)

// SetParams sets the auth module's parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSubspace.SetParamSet(ctx, &params)
}

// GetParams gets the auth module's parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSubspace.GetParamSet(ctx, &params)
	return
}
