package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/tharsis/ethermint/tests"

	"github.com/teleport-network/teleport/x/aggregate/types"
)

func (suite *KeeperTestSuite) TestMintingEnabled() {
	sender := sdk.AccAddress(tests.GenerateAddress().Bytes())
	receiver := sdk.AccAddress(tests.GenerateAddress().Bytes())
	expPair := types.NewTokenPair(tests.GenerateAddress(), []string{"coin"}, true, types.OWNER_MODULE)
	id := expPair.GetID()

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"intrarelaying is disabled globally",
			func() {
				params := types.DefaultParams()
				params.EnableAggregate = false
				suite.app.AggregateKeeper.SetParams(suite.ctx, params)
			},
			false,
		},
		{
			"token pair not found",
			func() {},
			false,
		},
		{
			"intrarelaying is disabled for the given pair",
			func() {
				expPair.Enabled = false
				suite.app.AggregateKeeper.SetTokenPair(suite.ctx, expPair)
				suite.app.AggregateKeeper.SetDenomsMap(suite.ctx, expPair.Denoms, id)
				suite.app.AggregateKeeper.SetERC20Map(suite.ctx, expPair.GetERC20Contract(), id)
			},
			false,
		},
		{
			"token transfers are disabled",
			func() {
				expPair.Enabled = true
				suite.app.AggregateKeeper.SetTokenPair(suite.ctx, expPair)
				suite.app.AggregateKeeper.SetDenomsMap(suite.ctx, expPair.Denoms, id)
				suite.app.AggregateKeeper.SetERC20Map(suite.ctx, expPair.GetERC20Contract(), id)

				params := banktypes.DefaultParams()
				params.SendEnabled = []*banktypes.SendEnabled{
					{Denom: expPair.Denoms[0], Enabled: false},
				}
				suite.app.BankKeeper.SetParams(suite.ctx, params)
			},
			false,
		},
		{
			"ok",
			func() {
				suite.app.AggregateKeeper.SetTokenPair(suite.ctx, expPair)
				suite.app.AggregateKeeper.SetDenomsMap(suite.ctx, expPair.Denoms, id)
				suite.app.AggregateKeeper.SetERC20Map(suite.ctx, expPair.GetERC20Contract(), id)
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			tc.malleate()

			pair, err := suite.app.AggregateKeeper.MintingEnabled(suite.ctx, sender, receiver, expPair.ERC20Address, expPair.Denoms[0])
			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expPair, pair)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
