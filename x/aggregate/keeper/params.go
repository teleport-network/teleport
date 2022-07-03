package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitdao-io/bitchain/x/aggregate/types"
)

// GetParams returns the total set of aggregate parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the aggregate parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
