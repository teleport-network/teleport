package rvesting

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/teleport-network/teleport/x/rvesting/keeper"
	"github.com/teleport-network/teleport/x/rvesting/types"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	rvestingAddr := k.AccountKeeper.GetModuleAddress(types.ModuleName)
	params := k.GetParams(ctx)

	vestedCoins := sdk.NewCoins()
	for _, reward := range params.PerBlockReward {
		remainingCoin := k.BankKeeper.GetBalance(ctx, rvestingAddr, reward.GetDenom())
		if remainingCoin.IsZero() {
			continue
		}
		if remainingCoin.Amount.LT(reward.Amount) {
			vestedCoins = vestedCoins.Add(remainingCoin)
		} else {
			vestedCoins = vestedCoins.Add(reward)
		}
	}

	if vestedCoins.IsZero() {
		return
	}
	if err := k.BankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.FeeCollectorName, vestedCoins); err != nil {
		panic(err)
	}
}
