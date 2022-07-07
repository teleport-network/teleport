package rvesting

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitdao-io/bitnetwork/x/rvesting/keeper"
	"github.com/bitdao-io/bitnetwork/x/rvesting/types"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	params := k.GetParams(ctx)

	if !params.EnableVesting {
		return
	}

	vestedCoins := sdk.NewCoins()
	for _, reward := range params.PerBlockReward {
		remainingCoin := k.GetRemainingCoin(ctx, reward.GetDenom())
		if remainingCoin.IsZero() {
			continue
		}
		if remainingCoin.Amount.LT(reward.Amount) {
			vestedCoins = vestedCoins.Add(remainingCoin)
		} else {
			vestedCoins = vestedCoins.Add(reward)
		}
	}

	if !vestedCoins.IsZero() {
		if err := k.SendVestedCoins(ctx, vestedCoins); err != nil {
			panic(err)
		}
	}
}
