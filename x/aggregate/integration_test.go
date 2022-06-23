package aggregate_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	erc20contracts "github.com/teleport-network/teleport/syscontracts/erc20"
	endpointcontract "github.com/teleport-network/teleport/syscontracts/xibc_endpoint"
	xibctesting "github.com/teleport-network/teleport/x/xibc/testing"
)

type AggregateTestSuite struct {
	suite.Suite
	coordinator *xibctesting.Coordinator
	chainA      *xibctesting.TestChain
	chainB      *xibctesting.TestChain
}

func TestAggregateTestSuite(t *testing.T) {
	suite.Run(t, new(AggregateTestSuite))
}

func (suite *AggregateTestSuite) SetupTest() {
	suite.coordinator = xibctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(xibctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(xibctesting.GetChainID(1))
}

func (suite *AggregateTestSuite) TestReBindToken() {
	// deploy ERC20
	erc20Address := suite.DeployERC20(suite.chainA)
	println("================================================")
	println(erc20Address.String())
	println("================================================")

	// add erc20 trace
	err := suite.chainA.App.AggregateKeeper.RegisterERC20Trace(
		suite.chainA.GetContext(),
		erc20Address,
		common.BigToAddress(big.NewInt(0)).String(),
		suite.chainB.ChainID,
		uint8(0),
	)
	suite.Require().NoError(err)
	// check ERC20 trace
	token, _, exist, err := suite.chainA.App.AggregateKeeper.QueryERC20Trace(
		suite.chainA.GetContext(),
		erc20Address,
		suite.chainB.ChainID,
	)
	suite.Require().NoError(err)
	suite.Require().True(exist)
	suite.Equal(token, common.BigToAddress(big.NewInt(0)).String())

	// add erc20 trace
	err = suite.chainA.App.AggregateKeeper.RegisterERC20Trace(
		suite.chainA.GetContext(),
		erc20Address,
		common.BigToAddress(big.NewInt(1)).String(),
		suite.chainB.ChainID,
		uint8(0),
	)
	suite.Require().NoError(err)
	// check ERC20 trace
	token, _, exist, err = suite.chainA.App.AggregateKeeper.QueryERC20Trace(
		suite.chainA.GetContext(),
		erc20Address,
		suite.chainB.ChainID,
	)
	suite.Require().NoError(err)
	suite.Require().True(exist)
	suite.Equal(token, common.BigToAddress(big.NewInt(1)).String())
}

// ================================================================================================================
// Functions for step
// ================================================================================================================
func (suite *AggregateTestSuite) DeployERC20(fromChain *xibctesting.TestChain) common.Address {
	ctorArgs, err := erc20contracts.ERC20MinterBurnerDecimalsContract.ABI.Pack("", "name", "symbol", uint8(18))
	suite.Require().NoError(err)

	data := make([]byte, len(erc20contracts.ERC20MinterBurnerDecimalsContract.Bin)+len(ctorArgs))
	copy(data[:len(erc20contracts.ERC20MinterBurnerDecimalsContract.Bin)], erc20contracts.ERC20MinterBurnerDecimalsContract.Bin)
	copy(data[len(erc20contracts.ERC20MinterBurnerDecimalsContract.Bin):], ctorArgs)

	nonce := fromChain.App.EvmKeeper.GetNonce(fromChain.GetContext(), endpointcontract.EndpointContractAddress)
	contractAddr := crypto.CreateAddress(endpointcontract.EndpointContractAddress, nonce)

	res, err := fromChain.App.AggregateKeeper.CallEVMWithData(fromChain.GetContext(), endpointcontract.EndpointContractAddress, nil, data, false)
	suite.Require().NoError(err)
	suite.Require().False(res.Failed(), res.VmError)

	return contractAddr
}
