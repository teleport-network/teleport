package keeper_test

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/tharsis/ethermint/tests"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	erc20contracts "github.com/teleport-network/teleport/syscontracts/erc20"
	"github.com/teleport-network/teleport/x/aggregate/types"
)

func (suite *KeeperTestSuite) TestQueryERC20() {
	var contract common.Address
	testCases := []struct {
		name     string
		malleate func()
		res      bool
	}{
		{
			"erc20 not deployed",
			func() { contract = common.Address{} },
			false,
		},
		{
			"ok",
			func() { contract = suite.DeployContract("coin", "token", erc20Decimals) },
			true,
		},
	}
	for _, tc := range testCases {
		suite.SetupTest() // reset

		tc.malleate()

		res, err := suite.app.AggregateKeeper.QueryERC20(suite.ctx, contract)
		if tc.res {
			suite.Require().NoError(err)
			suite.Require().Equal(
				types.ERC20Data{Name: "coin", Symbol: "token", Decimals: erc20Decimals},
				res,
			)
		} else {
			suite.Require().Error(err)
		}
	}
}

func (suite *KeeperTestSuite) TestCallEVM() {
	testCases := []struct {
		name    string
		method  string
		expPass bool
	}{
		{
			"unknown method",
			"",
			false,
		},
		{
			"pass",
			"balanceOf",
			true,
		},
	}
	for _, tc := range testCases {
		suite.SetupTest() // reset

		erc20 := erc20contracts.ERC20MinterBurnerDecimalsContract.ABI
		contract := suite.DeployContract("coin", "token", erc20Decimals)
		account := tests.GenerateAddress()

		res, err := suite.app.AggregateKeeper.CallEVM(suite.ctx, erc20, types.ModuleAddress, contract, tc.method, account)
		if tc.expPass {
			suite.Require().IsTypef(&evmtypes.MsgEthereumTxResponse{}, res, tc.name)
			suite.Require().NoError(err)
		} else {
			suite.Require().Error(err)
		}
	}
}
