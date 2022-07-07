package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitdao-io/bitnetwork/x/rvesting/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	return &types.QueryParamsResponse{
		Params: k.GetParams(ctx),
	}, nil
}

func (k Keeper) Remaining(c context.Context, _ *types.QueryRemainingRequest) (*types.QueryRemainingResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	rvestingAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	totalRemaining := k.bankKeeper.GetAllBalances(ctx, rvestingAddr)
	return &types.QueryRemainingResponse{
		Address:   rvestingAddr.String(),
		Remaining: totalRemaining,
	}, nil
}
