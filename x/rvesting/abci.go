package rvesting

import (
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/teleport-network/teleport/x/rvesting/keeper"
	"github.com/teleport-network/teleport/x/rvesting/types"
	"time"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	rvestingAddr := k.AccountKeeper.GetModuleAddress(types.ModuleName)
	totalRemaining := k.BankKeeper.GetAllBalances(ctx, rvestingAddr)
	params := k.GetParams(ctx)

	vestedCoins := sdk.NewCoins()
	for _, remainingCoin := range totalRemaining {
		perBlockRewardAmt := params.PerBlockReward.AmountOf(remainingCoin.Denom)
		if perBlockRewardAmt.IsZero() {
			continue
		}
		if remainingCoin.Amount.LT(perBlockRewardAmt) {
			vestedCoins = vestedCoins.Add(remainingCoin)
		} else {
			vestedCoins = vestedCoins.Add(sdk.NewCoin(remainingCoin.Denom, perBlockRewardAmt))
		}
	}

	if vestedCoins.IsZero() {
		return
	}
	if err := k.BankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.FeeCollectorName, vestedCoins); err != nil {
		panic(err)
	}
}
