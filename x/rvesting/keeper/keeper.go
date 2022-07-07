package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/bitdao-io/bitnetwork/x/rvesting/types"
)

type Keeper struct {
	paramSubspace    paramtypes.Subspace
	bankKeeper       types.BankKeeper
	accountKeeper    types.AccountKeeper
	feeCollectorName string // name of the FeeCollector ModuleAccount
}

func NewKeeper(paramSubspace paramtypes.Subspace, bk types.BankKeeper, ak types.AccountKeeper, feeCollectorName string) Keeper {
	// set KeyTable if it has not already been set
	if !paramSubspace.HasKeyTable() {
		paramSubspace = paramSubspace.WithKeyTable(types.ParamKeyTable())
	}
	return Keeper{
		paramSubspace:    paramSubspace,
		bankKeeper:       bk,
		accountKeeper:    ak,
		feeCollectorName: feeCollectorName,
	}
}

func (k Keeper) GetRemainingCoin(ctx sdk.Context, denom string) sdk.Coin {
	rvestingAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	return k.bankKeeper.GetBalance(ctx, rvestingAddr, denom)
}

func (k Keeper) SendVestedCoins(ctx sdk.Context, vestedCoins sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, vestedCoins)
}
