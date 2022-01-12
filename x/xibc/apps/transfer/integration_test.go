package transfer_test

import (
	"math/big"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/tharsis/ethermint/tests"
	evm "github.com/tharsis/ethermint/x/evm/types"

	transfercontract "github.com/teleport-network/teleport/syscontracts/xibc_transfer"
	erc20contracts "github.com/teleport-network/teleport/x/aggregate/types/contracts"
	"github.com/teleport-network/teleport/x/xibc/apps/transfer/types"
	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
	xibctesting "github.com/teleport-network/teleport/x/xibc/testing"
)

type TransferTestSuite struct {
	suite.Suite
	coordinator *xibctesting.Coordinator
	chainA      *xibctesting.TestChain
	chainB      *xibctesting.TestChain
}

func TestTransferTestSuite(t *testing.T) {
	suite.Run(t, new(TransferTestSuite))
}

func (suite *TransferTestSuite) SetupTest() {
	suite.coordinator = xibctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(xibctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(xibctesting.GetChainID(1))
}

func (suite *TransferTestSuite) TestTransferBase() {
	path := xibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupClients(path)

	// prepare test data
	total := big.NewInt(100000000000000)
	amount := big.NewInt(100)

	// check balance
	balance := suite.chainA.App.EvmKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAddress)
	suite.Require().Equal(total.String(), balance.String())

	// deploy ERC20 on chainB
	erc20Address := suite.DeployERC20ByTransfer(suite.chainB)

	// add erc20 trace on chainB
	err := suite.chainB.App.AggregateKeeper.RegisterERC20Trace(
		suite.chainB.GetContext(),
		erc20Address,
		common.BigToAddress(big.NewInt(0)).String(),
		suite.chainA.ChainID,
	)
	suite.Require().NoError(err)

	// check ERC20 trace
	_, _, exist, err := suite.chainB.App.AggregateKeeper.QueryERC20Trace(
		suite.chainB.GetContext(),
		erc20Address,
		suite.chainA.ChainID,
	)
	suite.Require().NoError(err)
	suite.Require().True(exist)

	// send transferBase on chainA
	suite.SendTransferBase(
		suite.chainA,
		types.BaseTransferData{
			Receiver:   suite.chainB.SenderAddress.String(),
			DestChain:  suite.chainB.ChainID,
			RelayChain: "",
		},
		amount,
	)

	// commit block
	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)

	// check balance
	balance = suite.chainA.App.EvmKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAddress)
	suite.Require().Equal(big.NewInt(0).Sub(total, amount).String(), balance.String())

	// check token out
	outAmount := suite.OutTokens(
		suite.chainA,
		common.BigToAddress(big.NewInt(0)),
		suite.chainB.ChainID,
	)
	suite.Require().Equal(amount.String(), outAmount.String())

	// relay packet
	packetData := types.NewFungibleTokenPacketData(
		path.EndpointA.ChainName,
		path.EndpointB.ChainName,
		strings.ToLower(suite.chainA.SenderAddress.String()),
		strings.ToLower(suite.chainB.SenderAddress.String()),
		amount.Bytes(),
		strings.ToLower(common.BigToAddress(big.NewInt(0)).String()),
		strings.ToLower(""),
	)
	packet := packettypes.NewPacket(
		1,
		path.EndpointA.ChainName,
		path.EndpointB.ChainName,
		"",
		[]string{types.PortID},
		[][]byte{packetData.GetBytes()},
	)

	ack := packettypes.NewResultAcknowledgement([][]byte{{byte(1)}})
	err = path.RelayPacket(packet, ack.GetBytes())
	suite.Require().NoError(err)

	// check balance
	recvBalance := suite.BalanceOf(suite.chainB, erc20Address, suite.chainB.SenderAddress)
	suite.Require().Equal(amount.String(), recvBalance.String())

}

func (suite *TransferTestSuite) TestTransferBaseBack() {
	path := xibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupClients(path)
	suite.TestTransferBase()

	erc20Address := common.HexToAddress("0x13efE3b42ca903c6E96Dc300DCf3bdC32C5A1aD1")
	amount := big.NewInt(100)

	// check ERC20 trace
	_, _, exist, err := suite.chainB.App.AggregateKeeper.QueryERC20Trace(
		suite.chainB.GetContext(),
		erc20Address,
		suite.chainA.ChainID,
	)
	suite.Require().NoError(err)
	suite.Require().True(exist)

	// check balance
	recvBalance := suite.BalanceOf(suite.chainB, erc20Address, suite.chainB.SenderAddress)
	suite.Require().Equal(amount.String(), recvBalance.String())

	// Approve erc20 to transfer
	suite.Approve(suite.chainB, erc20Address, amount)

	suite.SendTransferERC20(
		suite.chainB,
		types.ERC20TransferData{
			TokenAddress: erc20Address,
			Receiver:     suite.chainA.SenderAddress.String(),
			Amount:       amount,
			DestChain:    suite.chainA.ChainID,
			RelayChain:   "",
		},
		big.NewInt(0),
	)

	recvBalance = suite.BalanceOf(suite.chainB, erc20Address, suite.chainB.SenderAddress)
	suite.Require().Equal("0", recvBalance.String())

	// relay packet
	packetData := types.NewFungibleTokenPacketData(
		path.EndpointB.ChainName,
		path.EndpointA.ChainName,
		strings.ToLower(suite.chainB.SenderAddress.String()),
		strings.ToLower(suite.chainA.SenderAddress.String()),
		amount.Bytes(),
		strings.ToLower(erc20Address.String()),
		strings.ToLower(common.BigToAddress(big.NewInt(0)).String()),
	)

	packet := packettypes.NewPacket(
		1,
		path.EndpointB.ChainName,
		path.EndpointA.ChainName,
		"",
		[]string{types.PortID},
		[][]byte{packetData.GetBytes()},
	)

	ack := packettypes.NewResultAcknowledgement([][]byte{{byte(1)}})
	err = path.RelayPacket(packet, ack.GetBytes())
	suite.Require().NoError(err)

	// check balance
	total := big.NewInt(100000000000000)
	balance := suite.chainA.App.EvmKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAddress)
	suite.Require().Equal(total.String(), balance.String())

	// check chainA token out
	outAmount := suite.OutTokens(
		suite.chainA,
		common.BigToAddress(big.NewInt(0)),
		suite.chainB.ChainID,
	)
	suite.Require().Equal("0", outAmount.String())
}

func (suite *TransferTestSuite) TestTransferERC20() {
	path := xibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupClients(path)

	amount := big.NewInt(10000000000)
	out := big.NewInt(1000)

	chainAERC20Address := suite.DeployERC20ByAccount(suite.chainA)
	suite.MintERC20Token(suite.chainA, suite.chainA.SenderAddress, chainAERC20Address, amount)

	// commit block
	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)

	// check balance
	recvBalance := suite.BalanceOf(suite.chainA, chainAERC20Address, suite.chainA.SenderAddress)
	suite.Require().Equal(amount.String(), recvBalance.String())

	// Approve erc20 to transfer
	suite.Approve(suite.chainA, chainAERC20Address, out)

	// deploy ERC20 on chainB
	chainBERC20Address := suite.DeployERC20ByTransfer(suite.chainB)

	// add erc20 trace on chainB
	err := suite.chainB.App.AggregateKeeper.RegisterERC20Trace(
		suite.chainB.GetContext(),
		chainBERC20Address,
		strings.ToLower(chainAERC20Address.String()),
		suite.chainA.ChainID,
	)
	suite.Require().NoError(err)

	// check ERC20 trace
	_, _, exist, err := suite.chainB.App.AggregateKeeper.QueryERC20Trace(
		suite.chainB.GetContext(),
		chainBERC20Address,
		suite.chainA.ChainID,
	)
	suite.Require().NoError(err)
	suite.Require().True(exist)

	// send transferBase on chainA
	suite.SendTransferERC20(
		suite.chainA,
		types.ERC20TransferData{
			TokenAddress: chainAERC20Address,
			Receiver:     suite.chainB.SenderAddress.String(),
			Amount:       out,
			DestChain:    suite.chainB.ChainID,
			RelayChain:   "",
		},
		big.NewInt(0),
	)
	recvBalance = suite.BalanceOf(suite.chainA, chainAERC20Address, suite.chainA.SenderAddress)
	suite.Require().Equal(strconv.FormatUint(amount.Uint64()-out.Uint64(), 10), recvBalance.String())

	// check chainA token out
	outAmount := suite.OutTokens(
		suite.chainA,
		chainAERC20Address,
		suite.chainB.ChainID,
	)
	suite.Require().Equal(out.String(), outAmount.String())

	// relay packet
	packetData := types.NewFungibleTokenPacketData(
		path.EndpointA.ChainName,
		path.EndpointB.ChainName,
		strings.ToLower(suite.chainA.SenderAddress.String()),
		strings.ToLower(suite.chainB.SenderAddress.String()),
		out.Bytes(),
		strings.ToLower(chainAERC20Address.String()),
		strings.ToLower(""),
	)
	packet := packettypes.NewPacket(
		1,
		path.EndpointA.ChainName,
		path.EndpointB.ChainName,
		"",
		[]string{types.PortID},
		[][]byte{packetData.GetBytes()},
	)

	// commit block
	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)

	ack := packettypes.NewResultAcknowledgement([][]byte{{byte(1)}})
	err = path.RelayPacket(packet, ack.GetBytes())
	suite.Require().NoError(err)

	// check balance
	recvBalance = suite.BalanceOf(suite.chainB, chainBERC20Address, suite.chainB.SenderAddress)
	suite.Require().Equal(out.String(), recvBalance.String())
}

func (suite *TransferTestSuite) TestTransferERC20Back() {
	path := xibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupClients(path)

	amount := big.NewInt(1000)

	suite.TestTransferERC20()

	chainAERC20Address := crypto.CreateAddress(suite.chainA.SenderAddress, 0)
	chainBERC20Address := crypto.CreateAddress(transfercontract.TransferContractAddress, 0)

	// check chainA token out
	outAmount := suite.OutTokens(
		suite.chainA,
		chainAERC20Address,
		suite.chainB.ChainID,
	)
	suite.Require().Equal(amount.String(), outAmount.String())

	// check ERC20 trace
	_, _, exist, err := suite.chainB.App.AggregateKeeper.QueryERC20Trace(
		suite.chainB.GetContext(),
		chainBERC20Address,
		suite.chainA.ChainID,
	)
	suite.Require().NoError(err)
	suite.Require().True(exist)

	// check balance
	recvBalance := suite.BalanceOf(suite.chainB, chainBERC20Address, suite.chainB.SenderAddress)
	suite.Require().Equal(amount.String(), recvBalance.String())

	// Approve erc20 to transfer
	suite.Approve(suite.chainB, chainBERC20Address, amount)

	suite.SendTransferERC20(
		suite.chainB,
		types.ERC20TransferData{
			TokenAddress: chainBERC20Address,
			Receiver:     suite.chainA.SenderAddress.String(),
			Amount:       amount,
			DestChain:    suite.chainA.ChainID,
			RelayChain:   "",
		},
		big.NewInt(0),
	)

	recvBalance = suite.BalanceOf(suite.chainB, chainBERC20Address, suite.chainB.SenderAddress)
	suite.Require().Equal("0", recvBalance.String())

	// relay packet
	packetData := types.NewFungibleTokenPacketData(
		path.EndpointB.ChainName,
		path.EndpointA.ChainName,
		strings.ToLower(suite.chainB.SenderAddress.String()),
		strings.ToLower(suite.chainA.SenderAddress.String()),
		amount.Bytes(),
		strings.ToLower(chainBERC20Address.String()),
		strings.ToLower(chainAERC20Address.String()),
	)

	packet := packettypes.NewPacket(
		1,
		path.EndpointB.ChainName,
		path.EndpointA.ChainName,
		"",
		[]string{types.PortID},
		[][]byte{packetData.GetBytes()},
	)

	ack := packettypes.NewResultAcknowledgement([][]byte{{byte(1)}})
	err = path.RelayPacket(packet, ack.GetBytes())
	suite.Require().NoError(err)

	// check chainA token out
	outAmount = suite.OutTokens(
		suite.chainA,
		chainAERC20Address,
		suite.chainB.ChainID,
	)
	suite.Require().Equal("0", outAmount.String())
}

// ================================================================================================================
// Functions for step
// ================================================================================================================

func (suite *TransferTestSuite) DeployERC20ByTransfer(fromChain *xibctesting.TestChain) common.Address {
	ctorArgs, err := erc20contracts.ERC20MinterBurnerDecimalsContract.ABI.Pack("", "name", "symbol", uint8(18))
	suite.Require().NoError(err)

	data := make([]byte, len(erc20contracts.ERC20MinterBurnerDecimalsContract.Bin)+len(ctorArgs))
	copy(data[:len(erc20contracts.ERC20MinterBurnerDecimalsContract.Bin)], erc20contracts.ERC20MinterBurnerDecimalsContract.Bin)
	copy(data[len(erc20contracts.ERC20MinterBurnerDecimalsContract.Bin):], ctorArgs)

	nonce := fromChain.App.EvmKeeper.GetNonce(fromChain.GetContext(), transfercontract.TransferContractAddress)
	contractAddr := crypto.CreateAddress(transfercontract.TransferContractAddress, nonce)

	res, err := fromChain.App.XIBCTransferKeeper.CallEVMWithPayload(fromChain.GetContext(), transfercontract.TransferContractAddress, nil, data)
	suite.Require().NoError(err)
	suite.Require().False(res.Failed(), res.VmError)

	return contractAddr
}

func (suite *TransferTestSuite) DeployERC20ByAccount(fromChain *xibctesting.TestChain) common.Address {
	ctorArgs, err := erc20contracts.ERC20MinterBurnerDecimalsContract.ABI.Pack("", "name", "symbol", uint8(18))
	suite.Require().NoError(err)

	data := make([]byte, len(erc20contracts.ERC20MinterBurnerDecimalsContract.Bin)+len(ctorArgs))
	copy(data[:len(erc20contracts.ERC20MinterBurnerDecimalsContract.Bin)], erc20contracts.ERC20MinterBurnerDecimalsContract.Bin)
	copy(data[len(erc20contracts.ERC20MinterBurnerDecimalsContract.Bin):], ctorArgs)

	nonce := fromChain.App.EvmKeeper.GetNonce(fromChain.GetContext(), fromChain.SenderAddress)
	contractAddr := crypto.CreateAddress(fromChain.SenderAddress, nonce)

	res, err := fromChain.App.XIBCTransferKeeper.CallEVMWithPayload(fromChain.GetContext(), fromChain.SenderAddress, nil, data)
	suite.Require().NoError(err)
	suite.Require().False(res.Failed(), res.VmError)

	return contractAddr
}

func (suite *TransferTestSuite) BalanceOf(fromChain *xibctesting.TestChain, contract common.Address, account common.Address) *big.Int {
	erc20 := erc20contracts.ERC20MinterBurnerDecimalsContract.ABI

	res, err := fromChain.App.XIBCTransferKeeper.CallEVM(
		fromChain.GetContext(),
		erc20,
		types.ModuleAddress,
		contract,
		"balanceOf",
		account,
	)
	suite.Require().NoError(err)

	var balance types.Amount
	err = erc20.UnpackIntoInterface(&balance, "balanceOf", res.Ret)
	suite.Require().NoError(err)

	return balance.Value
}

func (suite *TransferTestSuite) SendTransferBase(fromChain *xibctesting.TestChain, data types.BaseTransferData, amount *big.Int) {
	transferData, err := transfercontract.TransferContract.ABI.Pack("sendTransferBase", data)
	suite.Require().NoError(err)

	_ = suite.SendTx(fromChain, transfercontract.TransferContractAddress, amount, transferData)
}

func (suite *TransferTestSuite) MintERC20Token(fromChain *xibctesting.TestChain, to, erc20Address common.Address, amount *big.Int) {
	ctorArgs, err := erc20contracts.ERC20MinterBurnerDecimalsContract.ABI.Pack("mint", to, amount)
	suite.Require().NoError(err)

	_ = suite.SendTx(fromChain, erc20Address, big.NewInt(0), ctorArgs)
}

func (suite *TransferTestSuite) Approve(fromChain *xibctesting.TestChain, erc20Address common.Address, amount *big.Int) {
	transferData, err := erc20contracts.ERC20MinterBurnerDecimalsContract.ABI.Pack("approve", transfercontract.TransferContractAddress, amount)
	suite.Require().NoError(err)

	_ = suite.SendTx(fromChain, erc20Address, big.NewInt(0), transferData)
	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)
}

func (suite *TransferTestSuite) SendTransferERC20(fromChain *xibctesting.TestChain, data types.ERC20TransferData, amount *big.Int) {
	transferData, err := transfercontract.TransferContract.ABI.Pack("sendTransferERC20", data)
	suite.Require().NoError(err)

	_ = suite.SendTx(fromChain, transfercontract.TransferContractAddress, amount, transferData)
	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)
}

func (suite *TransferTestSuite) OutTokens(fromChain *xibctesting.TestChain, tokenAddress common.Address, destChain string) *big.Int {
	res, err := fromChain.App.XIBCTransferKeeper.CallEVM(
		fromChain.GetContext(),
		transfercontract.TransferContract.ABI,
		types.ModuleAddress,
		transfercontract.TransferContractAddress,
		"outTokens",
		tokenAddress,
		destChain,
	)
	suite.Require().NoError(err)

	var amount types.Amount
	err = transfercontract.TransferContract.ABI.UnpackIntoInterface(&amount, "outTokens", res.Ret)
	suite.Require().NoError(err)

	return amount.Value
}

// ================================================================================================================
// EVM transaction (return events)
// ================================================================================================================

func (suite *TransferTestSuite) SendTx(fromChain *xibctesting.TestChain, contractAddr common.Address, amount *big.Int, transferData []byte) *evm.MsgEthereumTx {
	ctx := sdk.WrapSDKContext(fromChain.GetContext())
	chainID := fromChain.App.EvmKeeper.ChainID()
	signer := tests.NewSigner(fromChain.SenderPrivKey)

	nonce := fromChain.App.EvmKeeper.GetNonce(fromChain.GetContext(), fromChain.SenderAddress)
	ercTransferTx := evm.NewTx(
		chainID,
		nonce,
		&contractAddr,
		amount,
		25000000,
		big.NewInt(0),
		big.NewInt(0),
		big.NewInt(0),
		transferData,
		&ethtypes.AccessList{}, // accesses
	)

	ercTransferTx.From = fromChain.SenderAddress.Hex()
	err := ercTransferTx.Sign(ethtypes.LatestSignerForChainID(chainID), signer)
	suite.Require().NoError(err)
	rsp, err := fromChain.App.EvmKeeper.EthereumTx(ctx, ercTransferTx)
	suite.Require().NoError(err)
	suite.Require().Empty(rsp.VmError)
	return ercTransferTx
}
