package keeper_test

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func (suite *KeeperTestSuite) TestERC20Trace() {
	testCases := []struct {
		name     string
		malleate func()
	}{
		{
			"register one trace",
			func() {
				tokenAddress := common.BigToAddress(big.NewInt(200))
				originToken := "token0"
				originChain := "chain0"

				_, _, exist, err := suite.app.AggregateKeeper.QueryERC20Trace(suite.ctx, tokenAddress, originChain)
				suite.Require().NoError(err)
				suite.Require().False(exist)

				_, err = suite.app.AggregateKeeper.AddERC20TraceToTransferContract(suite.ctx, tokenAddress, originToken, originChain)
				suite.Require().NoError(err)

				_, _, exist, err = suite.app.AggregateKeeper.QueryERC20Trace(suite.ctx, tokenAddress, originChain)
				suite.Require().NoError(err)
				suite.Require().True(exist)
			},
		},
		{
			"register two trace",
			func() {
				tokenAddress := common.BigToAddress(big.NewInt(200))
				originToken0 := "token0"
				originChain0 := "chain0"
				originToken1 := "token1"
				originChain1 := "chain1"

				// ====================================

				_, _, exist, err := suite.app.AggregateKeeper.QueryERC20Trace(suite.ctx, tokenAddress, originChain0)
				suite.Require().NoError(err)
				suite.Require().False(exist)

				_, err = suite.app.AggregateKeeper.AddERC20TraceToTransferContract(suite.ctx, tokenAddress, originToken0, originChain0)
				suite.Require().NoError(err)

				_, _, exist, err = suite.app.AggregateKeeper.QueryERC20Trace(suite.ctx, tokenAddress, originChain0)
				suite.Require().NoError(err)
				suite.Require().True(exist)

				// ====================================

				_, _, exist, err = suite.app.AggregateKeeper.QueryERC20Trace(suite.ctx, tokenAddress, originChain1)
				suite.Require().NoError(err)
				suite.Require().False(exist)

				_, err = suite.app.AggregateKeeper.AddERC20TraceToTransferContract(suite.ctx, tokenAddress, originToken1, originChain1)
				suite.Require().NoError(err)

				_, _, exist, err = suite.app.AggregateKeeper.QueryERC20Trace(suite.ctx, tokenAddress, originChain1)
				suite.Require().NoError(err)
				suite.Require().True(exist)
			},
		},
		{
			"repeat registration",
			func() {
				tokenAddress := common.BigToAddress(big.NewInt(200))
				originToken := "token0"
				originChain := "chain0"

				_, _, exist, err := suite.app.AggregateKeeper.QueryERC20Trace(suite.ctx, tokenAddress, originChain)
				suite.Require().NoError(err)
				suite.Require().False(exist)

				_, err = suite.app.AggregateKeeper.AddERC20TraceToTransferContract(suite.ctx, tokenAddress, originToken, originChain)
				suite.Require().NoError(err)

				_, err = suite.app.AggregateKeeper.AddERC20TraceToTransferContract(suite.ctx, tokenAddress, originToken, originChain)
				suite.Require().NoError(err)
			},
		},
		{
			"register two traces with one origin chain",
			func() {
				tokenAddress := common.BigToAddress(big.NewInt(200))
				originToken0 := "token0"
				originToken1 := "token1"
				originChain := "chain0"

				_, _, exist, err := suite.app.AggregateKeeper.QueryERC20Trace(suite.ctx, tokenAddress, originChain)
				suite.Require().NoError(err)
				suite.Require().False(exist)

				_, err = suite.app.AggregateKeeper.AddERC20TraceToTransferContract(suite.ctx, tokenAddress, originToken0, originChain)
				suite.Require().NoError(err)

				_, err = suite.app.AggregateKeeper.AddERC20TraceToTransferContract(suite.ctx, tokenAddress, originToken1, originChain)
				suite.Require().NoError(err)
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest()
			tc.malleate()
		})
	}
}
