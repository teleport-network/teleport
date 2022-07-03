package bank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bk "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

// OverwriteBankKeeper extends from the BaseBankKeeper, and overwrites function to adjust the specified features of bitchain.
type OverwriteBankKeeper struct {
	bk.BaseKeeper
}

func NewOverwriteBankKeeper(baseKeeper bk.BaseKeeper) OverwriteBankKeeper {
	return OverwriteBankKeeper{
		BaseKeeper: baseKeeper,
	}
}

// BurnCoins overwrites the function of BurnCoins implemented in bank.BaseKeeper.
// Instead of burning the coins from the module account, it transfers the coins to FeeCollector module.
// It will panic if the module account does not exist or is unauthorized.
func (k OverwriteBankKeeper) BurnCoins(ctx sdk.Context, moduleName string, amounts sdk.Coins) error {
	return k.SendCoinsFromModuleToModule(ctx, moduleName, authtypes.FeeCollectorName, amounts)
}
