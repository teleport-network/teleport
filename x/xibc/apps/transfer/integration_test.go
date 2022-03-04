package transfer_test

import (
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/accounts/abi"
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
	suite.SendTransferBase(
		suite.chainA,
		types.BaseTransferData{
			Receiver:   strings.ToLower(suite.chainB.SenderAddress.String()),
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
	ack, err := packettypes.NewResultAcknowledgement([][]byte{{byte(1)}}).GetBytes()
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

	suite.SendTransferERC20(
		suite.chainB,
		types.ERC20TransferData{
			TokenAddress: erc20Address,
			Receiver:     strings.ToLower(suite.chainA.SenderAddress.String()),
			Amount:       amount,
			DestChain:    suite.chainA.ChainID,
			RelayChain:   "",
		},
		big.NewInt(0),
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

	ack, err := packettypes.NewResultAcknowledgement([][]byte{{byte(1)}}).GetBytes()
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

	ack, err := packettypes.NewResultAcknowledgement([][]byte{{byte(1)}}).GetBytes()
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

	ack, err := packettypes.NewResultAcknowledgement([][]byte{{byte(1)}}).GetBytes()
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

	ack, err := packettypes.NewResultAcknowledgement([][]byte{{byte(1)}}).GetBytes()
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

	ack, err := packettypes.NewResultAcknowledgement([][]byte{{byte(1)}}).GetBytes()
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
	suite.SendTransferERC20(
		suite.chainA,
		types.ERC20TransferData{
			TokenAddress: wtelecontract.WTELEContractAddress,
			Receiver:     suite.chainB.SenderAddress.String(),
			Amount:       out,
			DestChain:    suite.chainB.ChainID,
			RelayChain:   "",
		},
		big.NewInt(0),
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

	ack, err := packettypes.NewResultAcknowledgement([][]byte{{byte(1)}}).GetBytes()
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

	ack, err := packettypes.NewResultAcknowledgement([][]byte{{byte(1)}}).GetBytes()
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

func (suite *TransferTestSuite) TestRemoteContractCallAgent() {
	pathAtoB := xibctesting.NewPath(suite.chainA, suite.chainB)
	pathBtoC := xibctesting.NewPath(suite.chainB, suite.chainC)

	suite.coordinator.SetupClients(pathAtoB)
	suite.coordinator.SetupClients(pathBtoC)

	chainBErc20 := suite.DeployERC20(suite.chainB, transfercontract.TransferContractAddress, uint8(18))
	chainCErc20 := suite.DeployERC20(suite.chainC, transfercontract.TransferContractAddress, uint8(18))

	amount := big.NewInt(10000000000)
	out := big.NewInt(10000)

	suite.DepositWTeleToken(suite.chainA, amount)

	// check balance
	recvBalance := suite.BalanceOf(suite.chainA, wtelecontract.WTELEContractAddress, suite.chainA.SenderAddress)
	suite.Require().Equal(amount.String(), recvBalance.String())

	// Approve erc20 to transfer
	suite.Approve(suite.chainA, wtelecontract.WTELEContractAddress, out)

	// register erc20 trace chainA to chainB
	err := suite.chainB.App.AggregateKeeper.RegisterERC20Trace(
		suite.chainB.GetContext(),
		chainBErc20,
		strings.ToLower(wtelecontract.WTELEContractAddress.String()),
		suite.chainA.ChainID,
		uint8(0),
	)
	suite.Require().NoError(err)
	// check ERC20 trace
	_, _, exist, err := suite.chainB.App.AggregateKeeper.QueryERC20Trace(
		suite.chainB.GetContext(),
		chainBErc20,
		suite.chainA.ChainID,
	)
	suite.Require().NoError(err)
	suite.Require().True(exist)

	// register erc20 trace chainB to chainC
	err = suite.chainC.App.AggregateKeeper.RegisterERC20Trace(
		suite.chainC.GetContext(),
		chainCErc20,
		strings.ToLower(chainBErc20.String()),
		suite.chainB.ChainID,
		uint8(0),
	)
	suite.Require().NoError(err)
	// check ERC20 trace
	_, _, exist, err = suite.chainC.App.AggregateKeeper.QueryERC20Trace(
		suite.chainC.GetContext(),
		chainCErc20,
		suite.chainB.ChainID,
	)
	suite.Require().NoError(err)
	suite.Require().True(exist)
	// send remote contract call
	agentData := types.ERC20TransferData{
		TokenAddress: chainBErc20,
		Receiver:     strings.ToLower(suite.chainC.SenderAddress.String()),
		DestChain:    suite.chainC.ChainID,
		RelayChain:   "",
	}
	id := sha256.Sum256([]byte(suite.chainA.ChainID + "/" + suite.chainB.ChainID + "/" + "2"))
	agentPayload, err := agentcontract.AgentContract.ABI.Pack(
		"send",
		id[:],
		agentData.TokenAddress,
		suite.chainA.SenderAddress,
		agentData.Receiver,
		agentData.DestChain,
		agentData.RelayChain,
	)
	suite.Require().NoError(err)

	sequence := uint64(1)
	// send multi call
	transferBaseBackDataBytes, err := abi.Arguments{{Type: multicalltypes.TupleERC20TransferData}}.Pack(
		multicalltypes.ERC20TransferData{
			TokenAddress: wtelecontract.WTELEContractAddress,
			Receiver:     strings.ToLower(agentcontract.AgentContractAddress.String()),
			Amount:       out,
		},
	)
	rccDataBytes, err := abi.Arguments{{Type: multicalltypes.TupleRCCData}}.Pack(
		multicalltypes.RCCData{
			ContractAddress: strings.ToLower(agentcontract.AgentContractAddress.String()),
			Data:            agentPayload,
		},
	)
	// transfer Erc20 chainC to chainB
	MultiCallData := multicalltypes.MultiCallData{
		DestChain:  suite.chainB.ChainID,
		RelayChain: "",
		Functions:  []uint8{0, 2},
		Data:       [][]byte{transferBaseBackDataBytes, rccDataBytes},
	}
	// Approve erc20 to transfer
	suite.Approve(suite.chainA, wtelecontract.WTELEContractAddress, amount)
	suite.SendMultiCall(suite.chainA, big.NewInt(0), MultiCallData)
	suite.coordinator.CommitBlock(suite.chainA, suite.chainB, suite.chainC)

	// relay packet
	AtoBTransferErc20PacketData := types.NewFungibleTokenPacketData(
		pathAtoB.EndpointA.ChainName,
		pathAtoB.EndpointB.ChainName,
		sequence,
		strings.ToLower(suite.chainA.SenderAddress.String()),
		strings.ToLower(agentcontract.AgentContractAddress.String()),
		out.Bytes(),
		strings.ToLower(wtelecontract.WTELEContractAddress.String()),
		strings.ToLower(""),
	)
	DataListERC20Bz, err := AtoBTransferErc20PacketData.GetBytes()
	suite.NoError(err)

	agentPacketData := rcctypes.NewRCCPacketData(
		pathAtoB.EndpointA.ChainName,
		pathAtoB.EndpointB.ChainName,
		sequence,
		strings.ToLower(suite.chainA.SenderAddress.String()),
		strings.ToLower(agentcontract.AgentContractAddress.String()),
		agentPayload,
	)
	agentBz, err := agentPacketData.GetBytes()

	packet := packettypes.NewPacket(
		sequence,
		pathAtoB.EndpointA.ChainName,
		pathAtoB.EndpointB.ChainName,
		"",
		[]string{types.PortID, rcctypes.PortID},
		[][]byte{DataListERC20Bz, agentBz},
	)
	// relay packet

	ackBZ := suite.RCCAcks(suite.chainA, sha256.Sum256(agentBz))
	suite.Require().Equal([]byte{}, ackBZ)

	resultTransferErc20 := []byte{byte(1)}
	resultRcc, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000001")
	suite.Require().NoError(err)
	ack, err := packettypes.NewResultAcknowledgement([][]byte{resultTransferErc20, resultRcc}).GetBytes()
	suite.Require().NoError(err)
	err = pathAtoB.RelayPacket(packet, ack)
	suite.Require().NoError(err)

	// check chainB balance
	chainBbalances := suite.BalanceOf(suite.chainB, chainBErc20, agentcontract.AgentContractAddress)
	suite.Require().Equal(chainBbalances.String(), "0")

	// relay packet
	BtoCTransferErc20PacketData := types.NewFungibleTokenPacketData(
		pathBtoC.EndpointA.ChainName,
		pathBtoC.EndpointB.ChainName,
		sequence,
		strings.ToLower(agentcontract.AgentContractAddress.String()),
		strings.ToLower(suite.chainC.SenderAddress.String()),
		out.Bytes(),
		strings.ToLower(chainBErc20.String()),
		strings.ToLower(""),
	)
	DataListERC20Bz, err = BtoCTransferErc20PacketData.GetBytes()
	suite.NoError(err)
	BtoCTransferErc20packet := packettypes.NewPacket(
		sequence,
		pathBtoC.EndpointA.ChainName,
		pathBtoC.EndpointB.ChainName,
		"",
		[]string{types.PortID},
		[][]byte{DataListERC20Bz},
	)
	// commit block
	suite.coordinator.CommitBlock(suite.chainA, suite.chainB, suite.chainC)

	BtoCTransferErc20ack, err := packettypes.NewResultAcknowledgement([][]byte{{byte(1)}}).GetBytes()
	suite.Require().NoError(err)
	err = pathBtoC.RelayPacket(BtoCTransferErc20packet, BtoCTransferErc20ack)
	suite.Require().NoError(err)

	// check balance
	recvBalance = suite.BalanceOf(suite.chainC, chainCErc20, suite.chainC.SenderAddress)
	suite.Require().Equal(out.String(), recvBalance.String())
}

func (suite *TransferTestSuite) TestRemoteContractCallAgentBack() {
	pathBtoA := xibctesting.NewPath(suite.chainB, suite.chainA)
	pathCtoB := xibctesting.NewPath(suite.chainC, suite.chainB)
	//out := big.NewInt(10000)
	agentOut := big.NewInt(10000)
	suite.TestRemoteContractCallAgent()
	chainBERC20 := crypto.CreateAddress(transfercontract.TransferContractAddress, 0)
	chainCERC20 := chainBERC20
	// check ERC20 trace
	_, _, exist, err := suite.chainB.App.AggregateKeeper.QueryERC20Trace(
		suite.chainB.GetContext(),
		chainBERC20,
		suite.chainA.ChainID,
	)
	suite.Require().NoError(err)
	suite.Require().True(exist)
	// check ERC20 trace
	_, _, exist, err = suite.chainC.App.AggregateKeeper.QueryERC20Trace(
		suite.chainC.GetContext(),
		chainCERC20,
		suite.chainB.ChainID,
	)
	suite.Require().NoError(err)
	suite.Require().True(exist)
	// send remote contract call
	agentData := types.ERC20TransferData{
		TokenAddress: chainBERC20,
		Receiver:     strings.ToLower(suite.chainA.SenderAddress.String()),
		DestChain:    suite.chainA.ChainID,
		RelayChain:   "",
	}
	id := sha256.Sum256([]byte(suite.chainC.ChainID + "/" + suite.chainB.ChainID + "/" + "1"))
	agentPayload, err := agentcontract.AgentContract.ABI.Pack(
		"send",
		id[:],
		agentData.TokenAddress,
		suite.chainC.SenderAddress,
		agentData.Receiver,
		agentData.DestChain,
		agentData.RelayChain,
	)
	suite.Require().NoError(err)

	rccDataBytes, err := abi.Arguments{{Type: multicalltypes.TupleRCCData}}.Pack(
		multicalltypes.RCCData{
			ContractAddress: strings.ToLower(agentcontract.AgentContractAddress.String()),
			Data:            agentPayload,
		},
	)
	suite.Require().NoError(err)
	// send multi call
	transferBaseDataBytes, err := abi.Arguments{{Type: multicalltypes.TupleERC20TransferData}}.Pack(
		multicalltypes.ERC20TransferData{
			TokenAddress: chainCERC20,
			Receiver:     strings.ToLower(agentcontract.AgentContractAddress.String()),
			Amount:       agentOut,
		},
	)
	suite.Require().NoError(err)

	// transfer Erc20 chainC to chainB
	MultiCallData := multicalltypes.MultiCallData{
		DestChain:  suite.chainB.ChainID,
		RelayChain: "",
		Functions:  []uint8{0, 2},
		Data:       [][]byte{transferBaseDataBytes, rccDataBytes},
	}

	// Approve erc20 to transfer
	suite.Approve(suite.chainC, chainCERC20, agentOut)
	suite.SendMultiCall(suite.chainC, big.NewInt(0), MultiCallData)
	suite.coordinator.CommitBlock(suite.chainA, suite.chainB, suite.chainC)

	balances := suite.BalanceOf(suite.chainC, chainCERC20, suite.chainC.SenderAddress)
	suite.Require().Equal("0", balances.String())

	sequence := uint64(1)

	// relay packet
	ERC20PacketData := types.NewFungibleTokenPacketData(
		pathCtoB.EndpointA.ChainName,
		pathCtoB.EndpointB.ChainName,
		sequence,
		strings.ToLower(suite.chainC.SenderAddress.String()),
		strings.ToLower(agentcontract.AgentContractAddress.String()),
		agentOut.Bytes(),
		strings.ToLower(chainCERC20.String()),
		strings.ToLower(chainBERC20.String()),
	)
	// rcc packet data
	rccPacketData := rcctypes.NewRCCPacketData(
		pathCtoB.EndpointA.ChainName,
		pathCtoB.EndpointB.ChainName,
		sequence,
		strings.ToLower(suite.chainC.SenderAddress.String()),
		strings.ToLower(agentcontract.AgentContractAddress.String()),
		agentPayload,
	)
	rccBz, err := rccPacketData.GetBytes()
	suite.NoError(err)
	DataListERC20Bz, err := ERC20PacketData.GetBytes()
	suite.NoError(err)
	multiCallPacket := packettypes.NewPacket(
		sequence,
		pathCtoB.EndpointA.ChainName,
		pathCtoB.EndpointB.ChainName,
		"",
		[]string{types.PortID, rcctypes.PortID},
		[][]byte{DataListERC20Bz, rccBz},
	)

	resultTransferErc20 := []byte{byte(1)}
	resultRcc, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000001")
	suite.Require().NoError(err)
	ack, err := packettypes.NewResultAcknowledgement([][]byte{resultTransferErc20, resultRcc}).GetBytes()
	suite.Require().NoError(err)
	err = pathCtoB.RelayPacket(multiCallPacket, ack)
	suite.Require().NoError(err)

	// relay packet
	BToATransferErc20PacketData := types.NewFungibleTokenPacketData(
		pathBtoA.EndpointA.ChainName,
		pathBtoA.EndpointB.ChainName,
		sequence,
		strings.ToLower(agentcontract.AgentContractAddress.String()),
		strings.ToLower(suite.chainA.SenderAddress.String()),
		agentOut.Bytes(),
		strings.ToLower(chainBERC20.String()),
		strings.ToLower(wtelecontract.WTELEContractAddress.String()),
	)
	DataListERC20Bz, err = BToATransferErc20PacketData.GetBytes()
	suite.NoError(err)
	BToATransferErc20packet := packettypes.NewPacket(
		sequence,
		pathBtoA.EndpointA.ChainName,
		pathBtoA.EndpointB.ChainName,
		"",
		[]string{types.PortID},
		[][]byte{DataListERC20Bz},
	)
	// commit block
	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)

	BtoATransferErc20ack, err := packettypes.NewResultAcknowledgement([][]byte{{byte(1)}}).GetBytes()
	suite.Require().NoError(err)
	err = pathBtoA.RelayPacket(BToATransferErc20packet, BtoATransferErc20ack)
	suite.Require().NoError(err)
	// check chainA token out
	outAmount := suite.OutTokens(
		suite.chainA,
		wtelecontract.WTELEContractAddress,
		suite.chainB.ChainID,
	)

	suite.Require().Equal("0", outAmount.String())
}

func (suite *TransferTestSuite) TestAgentSendBase() {
	pathBtoA := xibctesting.NewPath(suite.chainB, suite.chainA)
	pathAtoC := xibctesting.NewPath(suite.chainA, suite.chainC)
	suite.coordinator.SetupClients(pathBtoA)
	suite.coordinator.SetupClients(pathAtoC)

	amount := big.NewInt(100)
	// chainA transferBase to chainB
	suite.TestTransferBase()
	suite.Equal(suite.OutTokens(suite.chainA, common.BigToAddress(big.NewInt(0)), suite.chainB.ChainID), amount)

	chainBERC20 := crypto.CreateAddress(transfercontract.TransferContractAddress, 0)
	balances := suite.BalanceOf(suite.chainB, chainBERC20, suite.chainB.SenderAddress)
	suite.Require().Equal(amount.String(), balances.String())
	// check ERC20 trace
	_, _, exist, err := suite.chainB.App.AggregateKeeper.QueryERC20Trace(
		suite.chainB.GetContext(),
		chainBERC20,
		suite.chainA.ChainID,
	)
	suite.Require().NoError(err)
	suite.Require().True(exist)
	// deploy ERC20 on chainC
	chainCERC20 := suite.DeployERC20(suite.chainC, transfercontract.TransferContractAddress, uint8(18))
	// add erc20 trace on chainC
	err = suite.chainC.App.AggregateKeeper.RegisterERC20Trace(
		suite.chainC.GetContext(),
		chainCERC20,
		common.BigToAddress(big.NewInt(0)).String(),
		suite.chainA.ChainID,
		uint8(0),
	)
	suite.Require().NoError(err)
	// check ERC20 trace
	_, _, exist, err = suite.chainC.App.AggregateKeeper.QueryERC20Trace(
		suite.chainC.GetContext(),
		chainCERC20,
		suite.chainA.ChainID,
	)
	suite.Require().NoError(err)
	suite.Require().True(exist)

	// chainA transferBase to ChainC
	agentData := types.ERC20TransferData{
		TokenAddress: common.BigToAddress(big.NewInt(0)),
		Receiver:     strings.ToLower(suite.chainC.SenderAddress.String()),
		DestChain:    suite.chainC.ChainID,
		RelayChain:   "",
	}
	// id = sha256(srcChain/destChain/seq)
	id := sha256.Sum256([]byte(suite.chainB.ChainID + "/" + suite.chainA.ChainID + "/" + "1"))
	agentPayload, err := agentcontract.AgentContract.ABI.Pack(
		"send",
		id[:],
		agentData.TokenAddress,
		suite.chainB.SenderAddress,
		agentData.Receiver,
		agentData.DestChain,
		agentData.RelayChain,
	)
	suite.Require().NoError(err)
	rccDataBytes, err := abi.Arguments{{Type: multicalltypes.TupleRCCData}}.Pack(
		multicalltypes.RCCData{
			ContractAddress: strings.ToLower(agentcontract.AgentContractAddress.String()),
			Data:            agentPayload,
		},
	)
	suite.Require().NoError(err)
	// send multi call
	transferBaseBackDataBytes, err := abi.Arguments{{Type: multicalltypes.TupleERC20TransferData}}.Pack(
		multicalltypes.ERC20TransferData{
			TokenAddress: chainBERC20,
			Receiver:     strings.ToLower(agentcontract.AgentContractAddress.String()),
			Amount:       amount,
		},
	)
	suite.Require().NoError(err)

	// transfer Erc20 chainC to chainB
	MultiCallData := multicalltypes.MultiCallData{
		DestChain:  suite.chainA.ChainID,
		RelayChain: "",
		Functions:  []uint8{0, 2},
		Data:       [][]byte{transferBaseBackDataBytes, rccDataBytes},
	}
	// Approve erc20 to transfer
	suite.Approve(suite.chainB, chainBERC20, amount)
	suite.SendMultiCall(suite.chainB, big.NewInt(0), MultiCallData)
	suite.coordinator.CommitBlock(suite.chainA, suite.chainB, suite.chainC)

	balances = suite.BalanceOf(suite.chainB, chainBERC20, suite.chainB.SenderAddress)
	suite.Require().Equal("0", balances.String())
	sequence := uint64(1)
	// relay packet
	ERC20PacketData := types.NewFungibleTokenPacketData(
		pathBtoA.EndpointA.ChainName,
		pathBtoA.EndpointB.ChainName,
		sequence,
		strings.ToLower(suite.chainB.SenderAddress.String()),
		strings.ToLower(agentcontract.AgentContractAddress.String()),
		amount.Bytes(),
		strings.ToLower(chainBERC20.String()),
		strings.ToLower(common.BigToAddress(big.NewInt(0)).String()),
	)
	// rcc packet data
	rccPacketData := rcctypes.NewRCCPacketData(
		pathBtoA.EndpointA.ChainName,
		pathBtoA.EndpointB.ChainName,
		sequence,
		strings.ToLower(suite.chainB.SenderAddress.String()),
		strings.ToLower(agentcontract.AgentContractAddress.String()),
		agentPayload,
	)
	DataListRccBz, err := rccPacketData.GetBytes()
	suite.NoError(err)
	DataListERC20Bz, err := ERC20PacketData.GetBytes()
	suite.NoError(err)
	multiCallPacket := packettypes.NewPacket(
		sequence,
		pathBtoA.EndpointA.ChainName,
		pathBtoA.EndpointB.ChainName,
		"",
		[]string{types.PortID, rcctypes.PortID},
		[][]byte{DataListERC20Bz, DataListRccBz},
	)

	resultTransferErc20 := []byte{byte(1)}
	resultRcc, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000001")
	suite.Require().NoError(err)
	ack, err := packettypes.NewResultAcknowledgement([][]byte{resultTransferErc20, resultRcc}).GetBytes()
	suite.Require().NoError(err)
	err = pathBtoA.RelayPacket(multiCallPacket, ack)
	suite.Require().NoError(err)
	suite.Equal(suite.OutTokens(suite.chainA, common.BigToAddress(big.NewInt(0)), suite.chainB.ChainID).String(), "0")
	suite.Equal(suite.OutTokens(suite.chainA, common.BigToAddress(big.NewInt(0)), suite.chainC.ChainID).String(), amount.String())

	AToCTransferBasePacketData := types.NewFungibleTokenPacketData(
		pathAtoC.EndpointA.ChainName,
		pathAtoC.EndpointB.ChainName,
		sequence,
		strings.ToLower(agentcontract.AgentContractAddress.String()),
		strings.ToLower(suite.chainC.SenderAddress.String()),
		amount.Bytes(),
		strings.ToLower(common.BigToAddress(big.NewInt(0)).String()),
		strings.ToLower(""),
	)
	DataListBaseBz, err := AToCTransferBasePacketData.GetBytes()
	suite.NoError(err)
	AToCTransferBasePacket := packettypes.NewPacket(
		sequence,
		pathAtoC.EndpointA.ChainName,
		pathAtoC.EndpointB.ChainName,
		"",
		[]string{types.PortID},
		[][]byte{DataListBaseBz},
	)
	// commit block
	suite.coordinator.CommitBlock(suite.chainA, suite.chainB, suite.chainC)
	AToCTransferBaseAck, err := packettypes.NewResultAcknowledgement([][]byte{{byte(1)}}).GetBytes()
	suite.Require().NoError(err)
	err = pathAtoC.RelayPacket(AToCTransferBasePacket, AToCTransferBaseAck)
	suite.Require().NoError(err)
	suite.Equal(suite.BalanceOf(suite.chainC, chainCERC20, suite.chainC.SenderAddress), amount)
}

func (suite *TransferTestSuite) TestAgentRefund() {
	pathAtoB := xibctesting.NewPath(suite.chainA, suite.chainB)
	pathBtoC := xibctesting.NewPath(suite.chainB, suite.chainC)

	suite.coordinator.SetupClients(pathAtoB)
	suite.coordinator.SetupClients(pathBtoC)

	chainBErc20 := suite.DeployERC20(suite.chainB, transfercontract.TransferContractAddress, uint8(18))

	amount := big.NewInt(10000000000)
	out := big.NewInt(10000)

	suite.DepositWTeleToken(suite.chainA, amount)

	// check balance
	recvBalance := suite.BalanceOf(suite.chainA, wtelecontract.WTELEContractAddress, suite.chainA.SenderAddress)
	suite.Require().Equal(amount.String(), recvBalance.String())

	// Approve erc20 to transfer
	suite.Approve(suite.chainA, wtelecontract.WTELEContractAddress, out)

	// register erc20 trace chainA to chainB
	err := suite.chainB.App.AggregateKeeper.RegisterERC20Trace(
		suite.chainB.GetContext(),
		chainBErc20,
		strings.ToLower(wtelecontract.WTELEContractAddress.String()),
		suite.chainA.ChainID,
		uint8(0),
	)
	suite.Require().NoError(err)
	// check ERC20 trace
	_, _, exist, err := suite.chainB.App.AggregateKeeper.QueryERC20Trace(
		suite.chainB.GetContext(),
		chainBErc20,
		suite.chainA.ChainID,
	)
	suite.Require().NoError(err)
	suite.Require().True(exist)

	// send remote contract call
	agentData := types.ERC20TransferData{
		TokenAddress: chainBErc20,
		Receiver:     strings.ToLower(suite.chainC.SenderAddress.String()),
		DestChain:    suite.chainC.ChainID,
		RelayChain:   "",
	}
	id := sha256.Sum256([]byte(suite.chainA.ChainID + "/" + suite.chainB.ChainID + "/" + "2"))
	agentPayload, err := agentcontract.AgentContract.ABI.Pack(
		"send",
		id[:],
		agentData.TokenAddress,
		suite.chainA.SenderAddress,
		agentData.Receiver,
		agentData.DestChain,
		agentData.RelayChain,
	)
	suite.Require().NoError(err)

	rccDataBytes, err := abi.Arguments{{Type: multicalltypes.TupleRCCData}}.Pack(
		multicalltypes.RCCData{
			ContractAddress: strings.ToLower(agentcontract.AgentContractAddress.String()),
			Data:            agentPayload,
		},
	)
	suite.Require().NoError(err)
	suite.Require().NotNil(rccDataBytes)
	// send multi call
	dataBz := multicalltypes.ERC20TransferData{
		TokenAddress: wtelecontract.WTELEContractAddress,
		Receiver:     strings.ToLower(agentcontract.AgentContractAddress.String()),
		Amount:       out,
	}
	transferERC20DataBytes, err := abi.Arguments{{Type: multicalltypes.TupleERC20TransferData}}.Pack(
		dataBz,
	)
	suite.Require().NoError(err)

	// transfer Erc20 chainC to chainB
	MultiCallData := multicalltypes.MultiCallData{
		DestChain:  suite.chainB.ChainID,
		RelayChain: "",
		Functions:  []uint8{0, 2},
		Data:       [][]byte{transferERC20DataBytes, rccDataBytes},
	}

	// Approve erc20 to transfer
	suite.Approve(suite.chainA, wtelecontract.WTELEContractAddress, out)

	suite.SendMultiCall(suite.chainA, big.NewInt(0), MultiCallData)
	suite.coordinator.CommitBlock(suite.chainA, suite.chainB, suite.chainC)
	sequence := uint64(1)
	// relay packet
	ERC20PacketData := types.NewFungibleTokenPacketData(
		pathAtoB.EndpointA.ChainName,
		pathAtoB.EndpointB.ChainName,
		sequence,
		strings.ToLower(suite.chainA.SenderAddress.String()),
		strings.ToLower(agentcontract.AgentContractAddress.String()),
		out.Bytes(),
		strings.ToLower(wtelecontract.WTELEContractAddress.String()),
		"",
	)
	// rcc packet data
	rccPacketData := rcctypes.NewRCCPacketData(
		pathAtoB.EndpointA.ChainName,
		pathAtoB.EndpointB.ChainName,
		sequence,
		strings.ToLower(suite.chainA.SenderAddress.String()),
		strings.ToLower(agentcontract.AgentContractAddress.String()),
		agentPayload,
	)
	DataListRccBz, err := rccPacketData.GetBytes()
	suite.NoError(err)
	DataListERC20Bz, err := ERC20PacketData.GetBytes()
	suite.NoError(err)
	multiCallPacket := packettypes.NewPacket(
		sequence,
		pathAtoB.EndpointA.ChainName,
		pathAtoB.EndpointB.ChainName,
		"",
		[]string{types.PortID, rcctypes.PortID},
		[][]byte{DataListERC20Bz, DataListRccBz},
	)
	resultTransferErc20 := []byte{byte(1)}
	resultRcc, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000001")
	suite.Require().NoError(err)
	ack, err := packettypes.NewResultAcknowledgement([][]byte{resultTransferErc20, resultRcc}).GetBytes()
	suite.Require().NoError(err)
	err = pathAtoB.RelayPacket(multiCallPacket, ack)
	suite.Require().NoError(err)
	suite.coordinator.CommitBlock(suite.chainA, suite.chainB, suite.chainC)
	suite.Equal(suite.OutTokens(suite.chainB, chainBErc20, suite.chainC.ChainID).String(), out.String())

	// relay packet
	BtoCERC20PacketData := types.NewFungibleTokenPacketData(
		pathBtoC.EndpointA.ChainName,
		pathBtoC.EndpointB.ChainName,
		sequence,
		strings.ToLower(agentcontract.AgentContractAddress.String()),
		strings.ToLower(suite.chainC.SenderAddress.String()),
		out.Bytes(),
		strings.ToLower(chainBErc20.String()),
		"",
	)
	DataListERC20Bz, err = BtoCERC20PacketData.GetBytes()
	suite.NoError(err)
	BtoCTransferErc20packet := packettypes.NewPacket(
		sequence,
		pathBtoC.EndpointA.ChainName,
		pathBtoC.EndpointB.ChainName,
		"",
		[]string{types.PortID},
		[][]byte{DataListERC20Bz},
	)
	// commit block
	errAck, err := packettypes.NewErrorAcknowledgement("onRecvPackt: binding is not exist").GetBytes()
	suite.NoError(err)
	err = pathBtoC.RelayPacket(BtoCTransferErc20packet, errAck)
	suite.NoError(err)
	suite.coordinator.CommitBlock(suite.chainA, suite.chainB, suite.chainC)

	suite.Equal(suite.OutTokens(suite.chainB, chainBErc20, suite.chainC.ChainID).String(), "0")
	suite.Equal(suite.AckStatus(suite.chainB, pathBtoC.EndpointA.ChainName, pathBtoC.EndpointB.ChainName, 1), uint8(2))
	suite.Equal(suite.GetAgentPacketExist(suite.chainB, pathBtoC.EndpointA.ChainName, pathBtoC.EndpointB.ChainName, "1"), true)

	suite.AgentRefund(suite.chainB, pathBtoC.EndpointA.ChainName, pathBtoC.EndpointB.ChainName, 1)
	suite.Equal(suite.BalanceOf(suite.chainB, chainBErc20, suite.chainA.SenderAddress), dataBz.Amount)

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
		Sent         bool
		Refunder     common.Address
		TokenAddress common.Address
		Amount       *big.Int
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
