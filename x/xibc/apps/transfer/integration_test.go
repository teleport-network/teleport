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

	"github.com/tharsis/ethermint/server/config"
	"github.com/tharsis/ethermint/tests"
	evm "github.com/tharsis/ethermint/x/evm/types"

	agentcontract "github.com/teleport-network/teleport/syscontracts/agent"
	erc20contracts "github.com/teleport-network/teleport/syscontracts/erc20"
	wtelecontract "github.com/teleport-network/teleport/syscontracts/wtele"
	multicallcontract "github.com/teleport-network/teleport/syscontracts/xibc_multicall"
	packetcontract "github.com/teleport-network/teleport/syscontracts/xibc_packet"
	rcccontract "github.com/teleport-network/teleport/syscontracts/xibc_rcc"
	transfercontract "github.com/teleport-network/teleport/syscontracts/xibc_transfer"
	multicalltypes "github.com/teleport-network/teleport/x/xibc/apps/multicall/types"
	rcctypes "github.com/teleport-network/teleport/x/xibc/apps/rcc/types"
	"github.com/teleport-network/teleport/x/xibc/apps/transfer/types"
	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
	xibctesting "github.com/teleport-network/teleport/x/xibc/testing"
)

type TransferTestSuite struct {
	suite.Suite
	coordinator *xibctesting.Coordinator
	chainA      *xibctesting.TestChain
	chainB      *xibctesting.TestChain
	chainC      *xibctesting.TestChain
}

func TestTransferTestSuite(t *testing.T) {
	suite.Run(t, new(TransferTestSuite))
}

func (suite *TransferTestSuite) SetupTest() {
	suite.coordinator = xibctesting.NewCoordinator(suite.T(), 3)
	suite.chainA = suite.coordinator.GetChain(xibctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(xibctesting.GetChainID(1))
	suite.chainC = suite.coordinator.GetChain(xibctesting.GetChainID(2))

}

func (suite *TransferTestSuite) TestTransferBase() common.Address {
	path := xibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupClients(path)

	// prepare test data
	total := big.NewInt(100000000000000)
	amount := big.NewInt(100)

	// check balance
	balance := suite.chainA.App.EvmKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAddress)
	suite.Require().Equal(total.String(), balance.String())

	// deploy ERC20 on chainB
	erc20Address := suite.DeployERC20(suite.chainB, transfercontract.TransferContractAddress, uint8(18))

	// add erc20 trace on chainB
	err := suite.chainB.App.AggregateKeeper.RegisterERC20Trace(
		suite.chainB.GetContext(),
		erc20Address,
		common.BigToAddress(big.NewInt(0)).String(),
		suite.chainA.ChainID,
		uint8(0),
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
	suite.SendTransfer(
		suite.chainA,
		types.TransferData{
			TokenAddress: common.HexToAddress("0x0000000000000000000000000000000000000000"),
			Receiver:     strings.ToLower(suite.chainB.SenderAddress.String()),
			Amount:       amount,
			DestChain:    suite.chainB.ChainID,
			RelayChain:   "",
		},
		types.Fee{
			TokenAddress: common.HexToAddress("0x0000000000000000000000000000000000000000"),
			Amount:       big.NewInt(0),
		},
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
	sequence := uint64(1)
	// relay packet
	packetData := types.NewFungibleTokenPacketData(
		path.EndpointA.ChainName,
		path.EndpointB.ChainName,
		sequence,
		strings.ToLower(suite.chainA.SenderAddress.String()),
		strings.ToLower(suite.chainB.SenderAddress.String()),
		amount.Bytes(),
		strings.ToLower(common.BigToAddress(big.NewInt(0)).String()),
		strings.ToLower(""),
	)
	DataListBaseBz, err := packetData.GetBytes()
	suite.NoError(err)
	packet := packettypes.NewPacket(
		sequence,
		path.EndpointA.ChainName,
		path.EndpointB.ChainName,
		"",
		[]string{types.PortID},
		[][]byte{DataListBaseBz},
	)
	ack, err := packettypes.NewResultAcknowledgement(
		[][]byte{{byte(1)}},
		path.EndpointB.Chain.SenderAcc.String(),
	).GetBytes()

	suite.Require().NoError(err)
	err = path.RelayPacket(packet, ack)
	suite.Require().NoError(err)

	// check balance
	recvBalance := suite.BalanceOf(suite.chainB, erc20Address, suite.chainB.SenderAddress)
	suite.Require().Equal(amount.String(), recvBalance.String())

	return erc20Address
}

func (suite *TransferTestSuite) TestTransferBaseBack() {
	path := xibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupClients(path)

	erc20Address := suite.TestTransferBase()
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

	suite.SendTransfer(
		suite.chainB,
		types.TransferData{
			TokenAddress: erc20Address,
			Receiver:     strings.ToLower(suite.chainA.SenderAddress.String()),
			Amount:       amount,
			DestChain:    suite.chainA.ChainID,
			RelayChain:   "",
		},
		types.Fee{
			TokenAddress: common.HexToAddress("0x0000000000000000000000000000000000000000"),
			Amount:       big.NewInt(0),
		},
	)

	recvBalance = suite.BalanceOf(suite.chainB, erc20Address, suite.chainB.SenderAddress)
	suite.Require().Equal("0", recvBalance.String())
	sequence := uint64(1)

	// relay packet
	packetData := types.NewFungibleTokenPacketData(
		path.EndpointB.ChainName,
		path.EndpointA.ChainName,
		sequence,
		strings.ToLower(suite.chainB.SenderAddress.String()),
		strings.ToLower(suite.chainA.SenderAddress.String()),
		amount.Bytes(),
		strings.ToLower(erc20Address.String()),
		strings.ToLower(common.BigToAddress(big.NewInt(0)).String()),
	)
	DataListERC20Bz, err := packetData.GetBytes()
	suite.NoError(err)
	packet := packettypes.NewPacket(
		sequence,
		path.EndpointB.ChainName,
		path.EndpointA.ChainName,
		"",
		[]string{types.PortID},
		[][]byte{DataListERC20Bz},
	)

	ack, err := packettypes.NewResultAcknowledgement(
		[][]byte{{byte(1)}},
		path.EndpointB.Chain.SenderAcc.String(),
	).GetBytes()

	suite.Require().NoError(err)
	err = path.RelayPacket(packet, ack)
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

	chainAERC20Address := suite.DeployERC20(suite.chainA, suite.chainA.SenderAddress, uint8(18))
	suite.MintERC20Token(suite.chainA, suite.chainA.SenderAddress, chainAERC20Address, amount)

	// commit block
	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)

	// check balance
	recvBalance := suite.BalanceOf(suite.chainA, chainAERC20Address, suite.chainA.SenderAddress)
	suite.Require().Equal(amount.String(), recvBalance.String())

	// Approve erc20 to transfer
	suite.Approve(suite.chainA, chainAERC20Address, out)

	// deploy ERC20 on chainB
	chainBERC20Address := suite.DeployERC20(suite.chainB, transfercontract.TransferContractAddress, uint8(18))

	// add erc20 trace on chainB
	err := suite.chainB.App.AggregateKeeper.RegisterERC20Trace(
		suite.chainB.GetContext(),
		chainBERC20Address,
		strings.ToLower(chainAERC20Address.String()),
		suite.chainA.ChainID,
		uint8(0),
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
	suite.SendTransfer(
		suite.chainA,
		types.TransferData{
			TokenAddress: chainAERC20Address,
			Receiver:     suite.chainB.SenderAddress.String(),
			Amount:       out,
			DestChain:    suite.chainB.ChainID,
			RelayChain:   "",
		},
		types.Fee{
			TokenAddress: common.HexToAddress("0x0000000000000000000000000000000000000000"),
			Amount:       big.NewInt(0),
		},
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
	sequence := uint64(1)

	// relay packet
	packetData := types.NewFungibleTokenPacketData(
		path.EndpointA.ChainName,
		path.EndpointB.ChainName,
		sequence,
		strings.ToLower(suite.chainA.SenderAddress.String()),
		strings.ToLower(suite.chainB.SenderAddress.String()),
		out.Bytes(),
		strings.ToLower(chainAERC20Address.String()),
		strings.ToLower(""),
	)
	DataListERC20Bz, err := packetData.GetBytes()
	suite.NoError(err)
	packet := packettypes.NewPacket(
		1,
		path.EndpointA.ChainName,
		path.EndpointB.ChainName,
		"",
		[]string{types.PortID},
		[][]byte{DataListERC20Bz},
	)

	// commit block
	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)

	ack, err := packettypes.NewResultAcknowledgement(
		[][]byte{{byte(1)}},
		path.EndpointB.Chain.SenderAcc.String(),
	).GetBytes()

	suite.Require().NoError(err)
	err = path.RelayPacket(packet, ack)
	suite.Require().NoError(err)

	// check balance
	recvBalance = suite.BalanceOf(suite.chainB, chainBERC20Address, suite.chainB.SenderAddress)
	suite.Require().Equal(out.String(), recvBalance.String())
}

func (suite *TransferTestSuite) TestTransferERC20Back() {
	path := xibctesting.NewPath(suite.chainA, suite.chainB)
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

	suite.SendTransfer(
		suite.chainB,
		types.TransferData{
			TokenAddress: chainBERC20Address,
			Receiver:     suite.chainA.SenderAddress.String(),
			Amount:       amount,
			DestChain:    suite.chainA.ChainID,
			RelayChain:   "",
		},
		types.Fee{
			TokenAddress: common.HexToAddress("0x0000000000000000000000000000000000000000"),
			Amount:       big.NewInt(0),
		},
	)

	recvBalance = suite.BalanceOf(suite.chainB, chainBERC20Address, suite.chainB.SenderAddress)
	suite.Require().Equal("0", recvBalance.String())
	sequence := uint64(1)
	// relay packet
	packetData := types.NewFungibleTokenPacketData(
		path.EndpointB.ChainName,
		path.EndpointA.ChainName,
		sequence,
		strings.ToLower(suite.chainB.SenderAddress.String()),
		strings.ToLower(suite.chainA.SenderAddress.String()),
		amount.Bytes(),
		strings.ToLower(chainBERC20Address.String()),
		strings.ToLower(chainAERC20Address.String()),
	)
	DataListERC20Bz, err := packetData.GetBytes()
	suite.NoError(err)
	packet := packettypes.NewPacket(
		sequence,
		path.EndpointB.ChainName,
		path.EndpointA.ChainName,
		"",
		[]string{types.PortID},
		[][]byte{DataListERC20Bz},
	)

	ack, err := packettypes.NewResultAcknowledgement(
		[][]byte{{byte(1)}},
		path.EndpointB.Chain.SenderAcc.String(),
	).GetBytes()

	suite.Require().NoError(err)
	err = path.RelayPacket(packet, ack)
	suite.Require().NoError(err)

	// check chainA token out
	outAmount = suite.OutTokens(
		suite.chainA,
		chainAERC20Address,
		suite.chainB.ChainID,
	)
	suite.Require().Equal("0", outAmount.String())
}

func (suite *TransferTestSuite) TestTransferScaledERC20() {
	path := xibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupClients(path)

	amount := big.NewInt(10000000000)
	out := big.NewInt(1000)
	scale, _ := sdk.NewIntFromString("1000000000000")

	chainAERC20Address := suite.DeployERC20(suite.chainA, suite.chainA.SenderAddress, uint8(6))
	suite.MintERC20Token(suite.chainA, suite.chainA.SenderAddress, chainAERC20Address, amount)

	// commit block
	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)

	// check balance
	recvBalance := suite.BalanceOf(suite.chainA, chainAERC20Address, suite.chainA.SenderAddress)
	suite.Require().Equal(amount.String(), recvBalance.String())

	// Approve erc20 to transfer
	suite.Approve(suite.chainA, chainAERC20Address, out)

	// deploy ERC20 on chainB
	chainBERC20Address := suite.DeployERC20(suite.chainB, transfercontract.TransferContractAddress, uint8(18))

	// add erc20 trace on chainB
	err := suite.chainB.App.AggregateKeeper.RegisterERC20Trace(
		suite.chainB.GetContext(),
		chainBERC20Address,
		strings.ToLower(chainAERC20Address.String()),
		suite.chainA.ChainID,
		uint8(12),
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
	suite.SendTransfer(
		suite.chainA,
		types.TransferData{
			TokenAddress: chainAERC20Address,
			Receiver:     suite.chainB.SenderAddress.String(),
			Amount:       out,
			DestChain:    suite.chainB.ChainID,
			RelayChain:   "",
		},
		types.Fee{
			TokenAddress: common.HexToAddress("0x0000000000000000000000000000000000000000"),
			Amount:       big.NewInt(0),
		},
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
	sequence := uint64(1)
	// relay packet
	packetData := types.NewFungibleTokenPacketData(
		path.EndpointA.ChainName,
		path.EndpointB.ChainName,
		sequence,
		strings.ToLower(suite.chainA.SenderAddress.String()),
		strings.ToLower(suite.chainB.SenderAddress.String()),
		out.Bytes(),
		strings.ToLower(chainAERC20Address.String()),
		strings.ToLower(""),
	)
	DataListERC20Bz, err := packetData.GetBytes()
	suite.NoError(err)
	packet := packettypes.NewPacket(
		sequence,
		path.EndpointA.ChainName,
		path.EndpointB.ChainName,
		"",
		[]string{types.PortID},
		[][]byte{DataListERC20Bz},
	)

	// commit block
	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)

	ack, err := packettypes.NewResultAcknowledgement(
		[][]byte{{byte(1)}},
		path.EndpointB.Chain.SenderAcc.String(),
	).GetBytes()

	suite.Require().NoError(err)
	err = path.RelayPacket(packet, ack)
	suite.Require().NoError(err)

	// check balance
	recvBalance = suite.BalanceOf(suite.chainB, chainBERC20Address, suite.chainB.SenderAddress)
	suite.Require().Equal(new(big.Int).Mul(out, scale.BigInt()).String(), recvBalance.String())
}

func (suite *TransferTestSuite) TestTransferScaledERC20Back() {
	path := xibctesting.NewPath(suite.chainA, suite.chainB)
	amount := big.NewInt(1000)
	scale, _ := sdk.NewIntFromString("1000000000000")

	suite.TestTransferScaledERC20()

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
	suite.Require().Equal(new(big.Int).Mul(amount, scale.BigInt()).String(), recvBalance.String())

	// Approve erc20 to transfer
	suite.Approve(suite.chainB, chainBERC20Address, new(big.Int).Mul(amount, scale.BigInt()))

	suite.SendTransfer(
		suite.chainB,
		types.TransferData{
			TokenAddress: chainBERC20Address,
			Receiver:     suite.chainA.SenderAddress.String(),
			Amount:       amount,
			DestChain:    suite.chainA.ChainID,
			RelayChain:   "",
		},
		types.Fee{
			TokenAddress: common.HexToAddress("0x0000000000000000000000000000000000000000"),
			Amount:       big.NewInt(0),
		},
	)

	recvBalance = suite.BalanceOf(suite.chainB, chainBERC20Address, suite.chainB.SenderAddress)
	suite.Require().Equal("0", recvBalance.String())
	sequence := uint64(1)
	// relay packet
	packetData := types.NewFungibleTokenPacketData(
		path.EndpointB.ChainName,
		path.EndpointA.ChainName,
		sequence,
		strings.ToLower(suite.chainB.SenderAddress.String()),
		strings.ToLower(suite.chainA.SenderAddress.String()),
		amount.Bytes(),
		strings.ToLower(chainBERC20Address.String()),
		strings.ToLower(chainAERC20Address.String()),
	)
	DataListERC20Bz, err := packetData.GetBytes()
	suite.NoError(err)
	packet := packettypes.NewPacket(
		sequence,
		path.EndpointB.ChainName,
		path.EndpointA.ChainName,
		"",
		[]string{types.PortID},
		[][]byte{DataListERC20Bz},
	)

	ack, err := packettypes.NewResultAcknowledgement(
		[][]byte{{byte(1)}},
		path.EndpointB.Chain.SenderAcc.String(),
	).GetBytes()

	suite.Require().NoError(err)
	err = path.RelayPacket(packet, ack)
	suite.Require().NoError(err)

	// check chainA token out
	outAmount = suite.OutTokens(
		suite.chainA,
		chainAERC20Address,
		suite.chainB.ChainID,
	)
	suite.Require().Equal("0", outAmount.String())
}

func (suite *TransferTestSuite) TestTransferWTele() {
	path := xibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupClients(path)

	amount := big.NewInt(10000000000)
	out := big.NewInt(1000)

	suite.DepositWTeleToken(suite.chainA, amount)

	// check balance
	recvBalance := suite.BalanceOf(suite.chainA, wtelecontract.WTELEContractAddress, suite.chainA.SenderAddress)
	suite.Require().Equal(amount.String(), recvBalance.String())

	// Approve erc20 to transfer
	suite.Approve(suite.chainA, wtelecontract.WTELEContractAddress, out)

	// deploy ERC20 on chainB
	chainBERC20Address := suite.DeployERC20(suite.chainB, transfercontract.TransferContractAddress, uint8(18))

	// add erc20 trace on chainB
	err := suite.chainB.App.AggregateKeeper.RegisterERC20Trace(
		suite.chainB.GetContext(),
		chainBERC20Address,
		strings.ToLower(wtelecontract.WTELEContractAddress.String()),
		suite.chainA.ChainID,
		uint8(0),
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
	suite.SendTransfer(
		suite.chainA,
		types.TransferData{
			TokenAddress: wtelecontract.WTELEContractAddress,
			Receiver:     suite.chainB.SenderAddress.String(),
			Amount:       out,
			DestChain:    suite.chainB.ChainID,
			RelayChain:   "",
		},
		types.Fee{
			TokenAddress: common.HexToAddress("0x0000000000000000000000000000000000000000"),
			Amount:       big.NewInt(0),
		},
	)
	recvBalance = suite.BalanceOf(suite.chainA, wtelecontract.WTELEContractAddress, suite.chainA.SenderAddress)
	suite.Require().Equal(strconv.FormatUint(amount.Uint64()-out.Uint64(), 10), recvBalance.String())

	// check chainA token out
	outAmount := suite.OutTokens(
		suite.chainA,
		wtelecontract.WTELEContractAddress,
		suite.chainB.ChainID,
	)
	suite.Require().Equal(out.String(), outAmount.String())
	sequence := uint64(1)
	// relay packet
	packetData := types.NewFungibleTokenPacketData(
		path.EndpointA.ChainName,
		path.EndpointB.ChainName,
		sequence,
		strings.ToLower(suite.chainA.SenderAddress.String()),
		strings.ToLower(suite.chainB.SenderAddress.String()),
		out.Bytes(),
		strings.ToLower(wtelecontract.WTELEContractAddress.String()),
		strings.ToLower(""),
	)
	DataListERC20Bz, err := packetData.GetBytes()
	suite.NoError(err)
	packet := packettypes.NewPacket(
		sequence,
		path.EndpointA.ChainName,
		path.EndpointB.ChainName,
		"",
		[]string{types.PortID},
		[][]byte{DataListERC20Bz},
	)

	// commit block
	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)

	ack, err := packettypes.NewResultAcknowledgement(
		[][]byte{{byte(1)}},
		path.EndpointB.Chain.SenderAcc.String(),
	).GetBytes()

	suite.Require().NoError(err)
	err = path.RelayPacket(packet, ack)
	suite.Require().NoError(err)

	// check balance
	recvBalance = suite.BalanceOf(suite.chainB, chainBERC20Address, suite.chainB.SenderAddress)
	suite.Require().Equal(out.String(), recvBalance.String())
}

func (suite *TransferTestSuite) TestTransferWTeleBack() {
	path := xibctesting.NewPath(suite.chainA, suite.chainB)
	total := big.NewInt(10000000000)
	amount := big.NewInt(1000)

	suite.TestTransferWTele()

	chainBERC20Address := crypto.CreateAddress(transfercontract.TransferContractAddress, 0)

	// check chainA token out
	outAmount := suite.OutTokens(
		suite.chainA,
		wtelecontract.WTELEContractAddress,
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

	suite.SendTransfer(
		suite.chainB,
		types.TransferData{
			TokenAddress: chainBERC20Address,
			Receiver:     suite.chainA.SenderAddress.String(),
			Amount:       amount,
			DestChain:    suite.chainA.ChainID,
			RelayChain:   "",
		},
		types.Fee{
			TokenAddress: common.HexToAddress("0x0000000000000000000000000000000000000000"),
			Amount:       big.NewInt(0),
		},
	)

	recvBalance = suite.BalanceOf(suite.chainB, chainBERC20Address, suite.chainB.SenderAddress)
	suite.Require().Equal("0", recvBalance.String())
	sequence := uint64(1)

	// relay packet
	packetData := types.NewFungibleTokenPacketData(
		path.EndpointB.ChainName,
		path.EndpointA.ChainName,
		sequence,
		strings.ToLower(suite.chainB.SenderAddress.String()),
		strings.ToLower(suite.chainA.SenderAddress.String()),
		amount.Bytes(),
		strings.ToLower(chainBERC20Address.String()),
		strings.ToLower(wtelecontract.WTELEContractAddress.String()),
	)
	DataListERC20Bz, err := packetData.GetBytes()
	suite.NoError(err)
	packet := packettypes.NewPacket(
		sequence,
		path.EndpointB.ChainName,
		path.EndpointA.ChainName,
		"",
		[]string{types.PortID},
		[][]byte{DataListERC20Bz},
	)

	ack, err := packettypes.NewResultAcknowledgement(
		[][]byte{{byte(1)}},
		path.EndpointB.Chain.SenderAcc.String(),
	).GetBytes()

	suite.Require().NoError(err)
	err = path.RelayPacket(packet, ack)
	suite.Require().NoError(err)

	// check chainA token out
	outAmount = suite.OutTokens(
		suite.chainA,
		wtelecontract.WTELEContractAddress,
		suite.chainB.ChainID,
	)
	suite.Require().Equal("0", outAmount.String())

	// check chainA WTele balance
	balance := suite.BalanceOf(suite.chainA, wtelecontract.WTELEContractAddress, suite.chainA.SenderAddress)
	suite.Require().Equal(total.String(), balance.String())
}

// ================================================================================================================
// Functions for step
// ================================================================================================================
func (suite *TransferTestSuite) Refund(fromChain *xibctesting.TestChain, srcChain, destChain, sequence string) bool {
	cus := agentcontract.AgentContract.ABI
	res, err := fromChain.App.XIBCTransferKeeper.CallEVM(
		fromChain.GetContext(),
		cus,
		types.ModuleAddress,
		agentcontract.AgentContractAddress,
		"refunded",
		srcChain+"/"+destChain+"/"+sequence,
	)
	suite.Require().NoError(err)

	var refunded bool
	err = cus.UnpackIntoInterface(&refunded, "refunded", res.Ret)
	suite.Require().NoError(err)

	return refunded
}

func (suite *TransferTestSuite) AckStatus(fromChain *xibctesting.TestChain, srcChain, destChain string, sequence uint64) uint8 {
	cus := packetcontract.PacketContract.ABI
	res, err := fromChain.App.XIBCTransferKeeper.CallEVM(
		fromChain.GetContext(),
		cus,
		types.ModuleAddress,
		packetcontract.PacketContractAddress,
		"getAckStatus",
		srcChain,
		destChain,
		sequence,
	)
	suite.Require().NoError(err)

	var status uint8
	err = cus.UnpackIntoInterface(&status, "getAckStatus", res.Ret)
	suite.Require().NoError(err)

	return status
}

func (suite *TransferTestSuite) GetNextSequenceSend(fromChain *xibctesting.TestChain, srcChain, destChain string) uint64 {
	cus := packetcontract.PacketContract.ABI
	res, err := fromChain.App.XIBCTransferKeeper.CallEVM(
		fromChain.GetContext(),
		cus,
		types.ModuleAddress,
		packetcontract.PacketContractAddress,
		"getNextSequenceSend",
		srcChain,
		destChain,
	)
	suite.Require().NoError(err)

	var seq uint64
	err = cus.UnpackIntoInterface(&seq, "getNextSequenceSend", res.Ret)
	suite.Require().NoError(err)

	return seq
}

func (suite *TransferTestSuite) GetAgentPacketExist(fromChain *xibctesting.TestChain, srcChain, destChain, sequences string) bool {
	cus := agentcontract.AgentContract.ABI
	res, err := fromChain.App.XIBCTransferKeeper.CallEVM(
		fromChain.GetContext(),
		cus,
		types.ModuleAddress,
		agentcontract.AgentContractAddress,
		"agentData",
		srcChain+"/"+destChain+"/"+sequences,
	)
	suite.Require().NoError(err)

	var exist struct {
		Sent                    bool
		RefundAddressOnTeleport common.Address
		TokenAddress            common.Address
		Amount                  *big.Int
	}
	err = cus.UnpackIntoInterface(&exist, "agentData", res.Ret)
	suite.Require().NoError(err)

	return exist.Sent
}

func (suite *TransferTestSuite) AgentRefund(fromChain *xibctesting.TestChain, srcChain, destChain string, sequence uint64) {
	agentData, err := agentcontract.AgentContract.ABI.Pack("refund", srcChain, destChain, sequence)
	suite.Require().NoError(err)

	_ = suite.SendTx(fromChain, agentcontract.AgentContractAddress, big.NewInt(0), agentData)
	suite.coordinator.CommitBlock(suite.chainA, suite.chainB, suite.chainC)
}

func (suite *TransferTestSuite) SendMultiCall(fromChain *xibctesting.TestChain, amount *big.Int, data multicalltypes.MultiCallData) {
	multiCallData, err := multicallcontract.MultiCallContract.ABI.Pack("multiCall", data)
	suite.Require().NoError(err)

	_ = suite.SendTx(fromChain, multicallcontract.MultiCallContractAddress, amount, multiCallData)
	suite.coordinator.CommitBlock(suite.chainA, suite.chainB, suite.chainC)
}

func (suite *TransferTestSuite) SendRemoteContractCall(fromChain *xibctesting.TestChain, data rcctypes.CallRCCData) {
	rccData, err := rcccontract.RCCContract.ABI.Pack("sendRemoteContractCall", data)
	suite.Require().NoError(err)

	_ = suite.SendTx(fromChain, rcccontract.RCCContractAddress, big.NewInt(0), rccData)
	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)
}

func (suite *TransferTestSuite) RCCAcks(fromChain *xibctesting.TestChain, hash [32]byte) []byte {
	rcc := rcccontract.RCCContract.ABI

	res, err := fromChain.App.XIBCTransferKeeper.CallEVM(
		fromChain.GetContext(),
		rcc,
		rcctypes.ModuleAddress,
		rcccontract.RCCContractAddress,
		"acks",
		hash,
	)
	suite.Require().NoError(err)

	var ack struct{ Value []byte }
	err = rcc.UnpackIntoInterface(&ack, "acks", res.Ret)
	suite.Require().NoError(err)

	return ack.Value
}

func (suite *TransferTestSuite) DeployERC20(fromChain *xibctesting.TestChain, deployer common.Address, scale uint8) common.Address {
	ctorArgs, err := erc20contracts.ERC20MinterBurnerDecimalsContract.ABI.Pack("", "name", "symbol", scale)
	suite.Require().NoError(err)

	data := make([]byte, len(erc20contracts.ERC20MinterBurnerDecimalsContract.Bin)+len(ctorArgs))
	copy(data[:len(erc20contracts.ERC20MinterBurnerDecimalsContract.Bin)], erc20contracts.ERC20MinterBurnerDecimalsContract.Bin)
	copy(data[len(erc20contracts.ERC20MinterBurnerDecimalsContract.Bin):], ctorArgs)

	nonce := fromChain.App.EvmKeeper.GetNonce(fromChain.GetContext(), deployer)
	contractAddr := crypto.CreateAddress(deployer, nonce)

	res, err := fromChain.App.XIBCTransferKeeper.CallEVMWithData(fromChain.GetContext(), deployer, nil, data)
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

func (suite *TransferTestSuite) SendTransfer(fromChain *xibctesting.TestChain, data types.TransferData, fee types.Fee) {
	transferData, err := transfercontract.TransferContract.ABI.Pack("sendTransfer", data, fee)
	suite.Require().NoError(err)

	amount := big.NewInt(0)
	if data.TokenAddress == common.HexToAddress("0x0000000000000000000000000000000000000000") {
		amount = amount.Add(amount, data.Amount)
	}

	if fee.TokenAddress == common.HexToAddress("0x0000000000000000000000000000000000000000") {
		amount = amount.Add(amount, fee.Amount)
	}

	_ = suite.SendTx(fromChain, transfercontract.TransferContractAddress, amount, transferData)
	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)
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

func (suite *TransferTestSuite) DepositWTeleToken(fromChain *xibctesting.TestChain, amount *big.Int) {
	transferData, err := wtelecontract.WTELEContract.ABI.Pack("deposit")
	suite.Require().NoError(err)

	_ = suite.SendTx(fromChain, wtelecontract.WTELEContractAddress, amount, transferData)
	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)
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
		config.DefaultGasCap,
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
	suite.Require().Empty(rsp.VmError, rsp.VmError)
	return ercTransferTx
}
