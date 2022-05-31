package xibc_test

import (
	"encoding/hex"
	"math/big"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tharsis/ethermint/server/config"
	"github.com/tharsis/ethermint/tests"
	evm "github.com/tharsis/ethermint/x/evm/types"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/teleport-network/teleport/syscontracts"
	erc20contracts "github.com/teleport-network/teleport/syscontracts/erc20"
	agentcontract "github.com/teleport-network/teleport/syscontracts/xibc_agent"
	endpointcontract "github.com/teleport-network/teleport/syscontracts/xibc_endpoint"
	packetcontract "github.com/teleport-network/teleport/syscontracts/xibc_packet"
	aggregatetypes "github.com/teleport-network/teleport/x/aggregate/types"
	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
	xibctesting "github.com/teleport-network/teleport/x/xibc/testing"
)

const (
	chainName  = "tendermint-0"
	chainName2 = "tendermin-1"
)

var clientHeight = clienttypes.NewHeight(0, 10)

type XIBCTestSuite struct {
	suite.Suite

	coordinator *xibctesting.Coordinator

	chainA *xibctesting.TestChain
	chainB *xibctesting.TestChain
	chainC *xibctesting.TestChain
}

// SetupTest creates a coordinator with 2 test chains.
func (suite *XIBCTestSuite) SetupTest() {
	suite.coordinator = xibctesting.NewCoordinator(suite.T(), 3)

	suite.chainA = suite.coordinator.GetChain(xibctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(xibctesting.GetChainID(1))
	suite.chainC = suite.coordinator.GetChain(xibctesting.GetChainID(2))

}

func TestXIBCTestSuite(t *testing.T) {
	suite.Run(t, new(XIBCTestSuite))
}

func (suite *XIBCTestSuite) TestBindToken() {
	// deploy ERC20
	erc20Address := suite.DeployERC20ByCrossChain(suite.chainA)

	// register erc20 trace
	err := suite.chainA.App.AggregateKeeper.RegisterERC20Trace(
		suite.chainA.GetContext(),
		erc20Address,
		common.BigToAddress(big.NewInt(0)).String(),
		suite.chainB.ChainID,
		uint8(0),
	)
	suite.Require().NoError(err)
	// check ERC20 trace
	token, amount, exist, err := suite.chainA.App.AggregateKeeper.QueryERC20Trace(
		suite.chainA.GetContext(),
		erc20Address,
		suite.chainB.ChainID,
	)
	suite.Require().NoError(err)
	suite.Require().True(exist)
	suite.Equal(token, common.BigToAddress(big.NewInt(0)).String())
	suite.Equal(amount.Int64(), int64(0))
}

func (suite *XIBCTestSuite) TestCrossChainTransferERC20() {
	suite.SetupTest()

	// setup testing conditions
	pathAToB := xibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupClients(pathAToB)

	// deploy ERC20
	chainAERC20 := suite.DeployERC20ByCrossChain(suite.chainA)
	chainBERC20 := suite.DeployERC20ByCrossChain(suite.chainB)
	suite.GrantERC20MintRoleByCrossChain(suite.chainA, chainAERC20, suite.chainA.SenderAddress)
	suite.MintERC20Token(suite.chainA, suite.chainA.SenderAddress, chainAERC20, big.NewInt(10000))
	balance := suite.ERC20Balance(suite.chainA, chainAERC20, suite.chainA.SenderAddress)
	suite.Equal(balance.Int64(), int64(10000))

	// add erc20 trace
	err := suite.chainB.App.AggregateKeeper.RegisterERC20Trace(
		suite.chainB.GetContext(),
		chainBERC20,
		strings.ToLower(chainAERC20.String()),
		suite.chainA.ChainID,
		uint8(0),
	)
	suite.Require().NoError(err)
	// check ERC20 trace
	token, amount, exist, err := suite.chainB.App.AggregateKeeper.QueryERC20Trace(
		suite.chainB.GetContext(),
		chainBERC20,
		suite.chainA.ChainID,
	)
	suite.Require().NoError(err)
	suite.Require().True(exist)
	suite.Equal(token, strings.ToLower(chainAERC20.String()))
	suite.Equal(amount.Int64(), int64(0))

	crossChainData := packettypes.CrossChainData{
		DstChain:        suite.chainB.ChainID,
		TokenAddress:    chainAERC20,
		Receiver:        strings.ToLower(suite.chainB.SenderAddress.String()),
		Amount:          big.NewInt(1000),
		ContractAddress: "",
		CallData:        []byte(""),
		CallbackAddress: common.BigToAddress(big.NewInt(0)),
		FeeOption:       0,
	}
	fee := packettypes.Fee{
		TokenAddress: chainAERC20,
		Amount:       big.NewInt(1000),
	}

	// send CrossChainCall Tx
	suite.Approve(suite.chainA, chainAERC20, endpointcontract.EndpointContractAddress, big.NewInt(2000))
	suite.CrossChainCall(suite.chainA, crossChainData, fee)
	balance = suite.ERC20Balance(suite.chainA, chainAERC20, suite.chainA.SenderAddress)
	suite.Equal(balance.Int64(), int64(8000))
	balance = suite.ERC20Balance(suite.chainA, chainAERC20, endpointcontract.EndpointContractAddress)
	suite.Equal(balance.Int64(), int64(1000))
	balance = suite.ERC20Balance(suite.chainA, chainAERC20, packetcontract.PacketContractAddress)
	suite.Equal(balance.Int64(), int64(1000))
	outToken := suite.OutTokens(suite.chainA, chainAERC20, suite.chainB.ChainID)
	suite.Require().Equal(int64(1000), outToken.Int64())

	// check packet fees
	fees := suite.GetPacketFees(suite.chainA, suite.chainA.ChainID, suite.chainB.ChainID, 1)
	suite.Require().Equal(chainAERC20, fees.TokenAddress)
	suite.Require().Equal(big.NewInt(1000).Int64(), fees.Amount.Int64())

	// packet and ack
	decodeString, err := hex.DecodeString("00000000000000000000000000000000000000000000000000000000000003e8")
	suite.Require().NoError(err)
	transferData := packettypes.TransferData{
		Receiver: strings.ToLower(crossChainData.Receiver),
		Amount:   decodeString,
		Token:    strings.ToLower(crossChainData.TokenAddress.String()),
		OriToken: "",
	}
	transferDataAbi, err := transferData.ABIPack()
	suite.Require().NoError(err)

	packet := packettypes.Packet{
		SrcChain:        suite.chainA.ChainID,
		DstChain:        suite.chainB.ChainID,
		Sequence:        1,
		Sender:          strings.ToLower(suite.chainA.SenderAddress.String()),
		TransferData:    transferDataAbi,
		CallData:        []byte(""),
		CallbackAddress: common.BigToAddress(big.NewInt(0)).String(),
		FeeOption:       0,
	}
	ack := packettypes.NewAcknowledgement(
		0,
		[]byte(""),
		"",
		strings.ToLower(suite.chainB.SenderAcc.String()),
		0,
	)
	ackData, err := ack.ABIPack()
	suite.Require().NoError(err)

	// relay
	err = pathAToB.RelayPacket(packet, ackData)
	suite.Require().NoError(err)
	balance = suite.ERC20Balance(suite.chainB, chainBERC20, suite.chainB.SenderAddress)
	suite.Require().Equal(int64(1000), balance.Int64())
	balance = suite.ERC20Balance(suite.chainA, chainAERC20, packetcontract.PacketContractAddress)
	suite.Equal(balance.Int64(), int64(0))
	balance = suite.ERC20Balance(suite.chainA, chainAERC20, suite.chainA.SenderAddress)
	suite.Equal(balance.Int64(), int64(9000))
	bindings := suite.Bindings(suite.chainB, chainBERC20, suite.chainA.ChainID)
	suite.Equal(bindings.OriToken, strings.ToLower(chainAERC20.String()))
	suite.Equal(bindings.Amount.Int64(), int64(1000))
	suite.Equal(bindings.Bound, true)
	suite.Equal(bindings.Scale, uint8(0))
	latestPacket := suite.GetLatestPacket(suite.chainB)
	suite.Equal(latestPacket.DstChain, packet.DstChain)
	suite.Equal(latestPacket.SrcChain, packet.SrcChain)
	suite.Equal(latestPacket.Sequence, packet.Sequence)
	suite.Equal(latestPacket.Sender, packet.Sender)
	suite.Equal(latestPacket.TransferData, packet.TransferData)
	suite.Equal(latestPacket.CallData, packet.CallData)
	suite.Equal(latestPacket.FeeOption, packet.FeeOption)
	suite.Equal(latestPacket.CallbackAddress, packet.CallbackAddress)

	status := suite.GetAckStatus(suite.chainA, suite.chainB.ChainID, 1)
	suite.Require().Equal(status, uint8(1))
	ack = suite.GetAck(suite.chainA, suite.chainB.ChainID, 1)
	suite.Require().Equal(uint64(0), ack.Code)
	suite.Require().Equal([]byte(""), ack.Result)
	suite.Require().Equal("", ack.Message)
	suite.Require().Equal(suite.chainB.SenderAcc.String(), ack.Relayer)
}

func (suite *XIBCTestSuite) TestCrossChainTransferBaseAndTransferBack() {
	suite.SetupTest()

	// setup testing conditions
	pathAToB := xibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupClients(pathAToB)

	// deploy ERC20
	chainABase := common.HexToAddress("0x0000000000000000000000000000000000000000")
	chainBERC20 := suite.DeployERC20ByCrossChain(suite.chainB)

	// add erc20 trace
	err := suite.chainB.App.AggregateKeeper.RegisterERC20Trace(
		suite.chainB.GetContext(),
		chainBERC20,
		strings.ToLower(chainABase.String()),
		suite.chainA.ChainID,
		uint8(0),
	)
	suite.Require().NoError(err)
	// check ERC20 trace
	token, amount, exist, err := suite.chainB.App.AggregateKeeper.QueryERC20Trace(
		suite.chainB.GetContext(),
		chainBERC20,
		suite.chainA.ChainID,
	)
	suite.Require().NoError(err)
	suite.Require().True(exist)
	suite.Equal(token, strings.ToLower(chainABase.String()))
	suite.Equal(amount.Int64(), int64(0))

	crossChainData := packettypes.CrossChainData{
		DstChain:        suite.chainB.ChainID,
		TokenAddress:    chainABase,
		Receiver:        strings.ToLower(suite.chainB.SenderAddress.String()),
		Amount:          big.NewInt(100),
		ContractAddress: "",
		CallData:        []byte(""),
		CallbackAddress: common.BigToAddress(big.NewInt(0)),
		FeeOption:       0,
	}
	fee := packettypes.Fee{
		TokenAddress: chainABase,
		Amount:       big.NewInt(100),
	}

	balances := suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAcc, "stake")
	suite.Require().Equal(balances.String(), "100000000000000stake")
	suite.CrossChainCall(suite.chainA, crossChainData, fee)
	balances = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAcc, "stake")
	suite.Require().Equal(balances.String(), "99999999999800stake")

	crossChainAddr, err := sdk.AccAddressFromHex(hex.EncodeToString(endpointcontract.EndpointContractAddress.Bytes()))
	suite.Require().NoError(err)
	packetAddr, err := sdk.AccAddressFromHex(hex.EncodeToString(packetcontract.PacketContractAddress.Bytes()))
	suite.Require().NoError(err)

	balances = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), crossChainAddr, "stake")
	suite.Require().Equal(balances.String(), "100stake")
	balances = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), packetAddr, "stake")
	suite.Require().Equal(balances.String(), "100stake")
	outToken := suite.OutTokens(suite.chainA, chainABase, suite.chainB.ChainID)
	suite.Require().Equal(int64(100), outToken.Int64())

	// check packet fees
	fees := suite.GetPacketFees(suite.chainA, suite.chainA.ChainID, suite.chainB.ChainID, 1)
	suite.Require().Equal(chainABase, fees.TokenAddress)
	suite.Require().Equal(big.NewInt(100).Int64(), fees.Amount.Int64())

	// packet and ack
	decodeString, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000064")
	suite.Require().NoError(err)
	transferData := packettypes.TransferData{
		Receiver: strings.ToLower(crossChainData.Receiver),
		Amount:   decodeString,
		Token:    strings.ToLower(crossChainData.TokenAddress.String()),
		OriToken: "",
	}
	transferDataAbi, err := transferData.ABIPack()
	suite.Require().NoError(err)

	packet := packettypes.Packet{
		SrcChain:        suite.chainA.ChainID,
		DstChain:        suite.chainB.ChainID,
		Sequence:        1,
		Sender:          strings.ToLower(suite.chainA.SenderAddress.String()),
		TransferData:    transferDataAbi,
		CallData:        []byte(""),
		CallbackAddress: common.BigToAddress(big.NewInt(0)).String(),
		FeeOption:       0,
	}
	ack := packettypes.NewAcknowledgement(
		0,
		[]byte(""),
		"",
		strings.ToLower(suite.chainB.SenderAcc.String()),
		0,
	)
	ackData, err := ack.ABIPack()
	suite.Require().NoError(err)

	// relay
	err = pathAToB.RelayPacket(packet, ackData)
	suite.Require().NoError(err)

	// check balance
	balance := suite.ERC20Balance(suite.chainB, chainBERC20, suite.chainB.SenderAddress)
	suite.Require().Equal(int64(100), balance.Int64())
	balances = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), packetAddr, "stake")
	suite.Require().Equal(balances.String(), "0stake")
	balances = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAcc, "stake")
	suite.Require().Equal(balances.String(), "99999999999900stake")

	// check bindings
	bindings := suite.Bindings(suite.chainB, chainBERC20, suite.chainA.ChainID)
	suite.Equal(bindings.OriToken, strings.ToLower(chainABase.String()))
	suite.Equal(bindings.Amount.Int64(), int64(100))
	suite.Equal(bindings.Bound, true)
	suite.Equal(bindings.Scale, uint8(0))

	// check ack
	status := suite.GetAckStatus(suite.chainA, suite.chainB.ChainID, 1)
	suite.Require().Equal(status, uint8(1))
	ack = suite.GetAck(suite.chainA, suite.chainB.ChainID, 1)
	suite.Require().Equal(uint64(0), ack.Code)
	suite.Require().Equal([]byte(""), ack.Result)
	suite.Require().Equal("", ack.Message)
	suite.Require().Equal(suite.chainB.SenderAcc.String(), ack.Relayer)

	crossChainData = packettypes.CrossChainData{
		DstChain:        suite.chainA.ChainID,
		TokenAddress:    chainBERC20,
		Receiver:        strings.ToLower(suite.chainA.SenderAddress.String()),
		Amount:          big.NewInt(100),
		ContractAddress: "",
		CallData:        []byte(""),
		CallbackAddress: common.BigToAddress(big.NewInt(0)),
		FeeOption:       0,
	}
	fee = packettypes.Fee{
		TokenAddress: chainBERC20,
		Amount:       big.NewInt(0),
	}
	suite.Approve(suite.chainB, chainBERC20, endpointcontract.EndpointContractAddress, big.NewInt(100))
	suite.CrossChainCall(suite.chainB, crossChainData, fee)
	balance = suite.ERC20Balance(suite.chainB, chainBERC20, suite.chainB.SenderAddress)
	suite.Require().Equal(int64(0), balance.Int64())

	// check bindings
	bindings = suite.Bindings(suite.chainB, chainBERC20, suite.chainA.ChainID)
	suite.Equal(bindings.OriToken, strings.ToLower(chainABase.String()))
	suite.Equal(bindings.Amount.Int64(), int64(0))
	suite.Equal(bindings.Bound, true)
	suite.Equal(bindings.Scale, uint8(0))

	// check packet fees
	fees = suite.GetPacketFees(suite.chainB, suite.chainB.ChainID, suite.chainA.ChainID, 1)
	suite.Require().Equal(chainBERC20, fees.TokenAddress)
	suite.Require().Equal(big.NewInt(0).Int64(), fees.Amount.Int64())

	// packet and ack
	decodeString, err = hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000064")
	suite.Require().NoError(err)
	transferData = packettypes.TransferData{
		Receiver: strings.ToLower(crossChainData.Receiver),
		Amount:   decodeString,
		Token:    strings.ToLower(crossChainData.TokenAddress.String()),
		OriToken: strings.ToLower(chainABase.String()),
	}
	transferDataAbi, err = transferData.ABIPack()
	suite.Require().NoError(err)

	packet = packettypes.Packet{
		SrcChain:        suite.chainB.ChainID,
		DstChain:        suite.chainA.ChainID,
		Sequence:        1,
		Sender:          strings.ToLower(suite.chainB.SenderAddress.String()),
		TransferData:    transferDataAbi,
		CallData:        []byte(""),
		CallbackAddress: common.BigToAddress(big.NewInt(0)).String(),
		FeeOption:       0,
	}
	ack = packettypes.NewAcknowledgement(
		0,
		[]byte(""),
		"",
		strings.ToLower(suite.chainB.SenderAcc.String()),
		0,
	)
	ackData, err = ack.ABIPack()
	suite.Require().NoError(err)

	// relay
	err = pathAToB.RelayPacket(packet, ackData)
	suite.Require().NoError(err)
	balances = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAcc, "stake")
	suite.Require().Equal(balances.String(), "100000000000000stake")

	balances = suite.chainA.App.BankKeeper.GetBalance(suite.chainA.GetContext(), crossChainAddr, "stake")
	suite.Require().Equal(balances.String(), "0stake")
	outToken = suite.OutTokens(suite.chainA, chainABase, suite.chainB.ChainID)
	suite.Require().Equal(int64(0), outToken.Int64())
}

func (suite *XIBCTestSuite) TestCrossChainCallAgent() {
	suite.SetupTest()

	// setup testing conditions
	pathAToB := xibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupClients(pathAToB)

	pathBToC := xibctesting.NewPath(suite.chainB, suite.chainC)
	suite.coordinator.SetupClientsWithoutRelayer(pathBToC)

	// deploy ERC20
	chainAERC20 := suite.DeployERC20ByCrossChain(suite.chainA)
	chainBERC20 := suite.DeployERC20ByCrossChain(suite.chainB)
	chainCERC20 := suite.DeployERC20ByCrossChain(suite.chainC)

	suite.GrantERC20MintRoleByCrossChain(suite.chainA, chainAERC20, suite.chainA.SenderAddress)
	suite.MintERC20Token(suite.chainA, suite.chainA.SenderAddress, chainAERC20, big.NewInt(10000))
	balance := suite.ERC20Balance(suite.chainA, chainAERC20, suite.chainA.SenderAddress)
	suite.Equal(balance.Int64(), int64(10000))

	// register erc20 trace A -> B
	err := suite.chainB.App.AggregateKeeper.RegisterERC20Trace(
		suite.chainB.GetContext(),
		chainBERC20,
		strings.ToLower(chainAERC20.String()),
		suite.chainA.ChainID,
		uint8(0),
	)
	suite.Require().NoError(err)
	// check ERC20 trace
	token, amount, exist, err := suite.chainB.App.AggregateKeeper.QueryERC20Trace(
		suite.chainB.GetContext(),
		chainBERC20,
		suite.chainA.ChainID,
	)
	suite.Require().NoError(err)
	suite.Require().True(exist)
	suite.Equal(token, strings.ToLower(chainAERC20.String()))
	suite.Equal(amount.Int64(), int64(0))

	// register erc20 trace B -> C
	err = suite.chainC.App.AggregateKeeper.RegisterERC20Trace(
		suite.chainC.GetContext(),
		chainCERC20,
		strings.ToLower(chainBERC20.String()),
		suite.chainB.ChainID,
		uint8(0),
	)
	suite.Require().NoError(err)
	// check ERC20 trace
	token, amount, exist, err = suite.chainC.App.AggregateKeeper.QueryERC20Trace(
		suite.chainC.GetContext(),
		chainCERC20,
		suite.chainB.ChainID,
	)
	suite.Require().NoError(err)
	suite.Require().True(exist)
	suite.Equal(token, strings.ToLower(chainBERC20.String()))
	suite.Equal(amount.Int64(), int64(0))

	callData, err := agentcontract.AgentContract.ABI.Pack(
		"send",
		chainBERC20,
		strings.ToLower(suite.chainA.SenderAddress.String()),
		suite.chainC.ChainID,
		big.NewInt(1000),
	)
	suite.Require().NoError(err)

	crossChainData := packettypes.CrossChainData{
		DstChain:        suite.chainB.ChainID,
		TokenAddress:    chainAERC20,
		Receiver:        strings.ToLower(agentcontract.AgentContractAddress.String()),
		Amount:          big.NewInt(2000),
		ContractAddress: syscontracts.AgentContractAddress,
		CallData:        callData,
		CallbackAddress: common.BigToAddress(big.NewInt(0)),
		FeeOption:       0,
	}
	fee := packettypes.Fee{
		TokenAddress: chainAERC20,
		Amount:       big.NewInt(0),
	}

	// send CrossChainCall Tx
	suite.Approve(suite.chainA, chainAERC20, endpointcontract.EndpointContractAddress, big.NewInt(2000))
	suite.CrossChainCall(suite.chainA, crossChainData, fee)
	balance = suite.ERC20Balance(suite.chainA, chainAERC20, suite.chainA.SenderAddress)
	suite.Equal(balance.Int64(), int64(8000))
	balance = suite.ERC20Balance(suite.chainA, chainAERC20, endpointcontract.EndpointContractAddress)
	suite.Equal(balance.Int64(), int64(2000))
	balance = suite.ERC20Balance(suite.chainA, chainAERC20, packetcontract.PacketContractAddress)
	suite.Equal(balance.Int64(), int64(0))
	outToken := suite.OutTokens(suite.chainA, chainAERC20, suite.chainB.ChainID)
	suite.Require().Equal(int64(2000), outToken.Int64())

	// check packet fees
	fees := suite.GetPacketFees(suite.chainA, suite.chainA.ChainID, suite.chainB.ChainID, 1)
	suite.Require().Equal(chainAERC20, fees.TokenAddress)
	suite.Require().Equal(big.NewInt(0).Int64(), fees.Amount.Int64())

	// packet and ack
	hexAmount, err := hex.DecodeString("00000000000000000000000000000000000000000000000000000000000007d0")
	suite.Require().NoError(err)
	transferData := packettypes.TransferData{
		Receiver: strings.ToLower(crossChainData.Receiver),
		Amount:   hexAmount,
		Token:    strings.ToLower(crossChainData.TokenAddress.String()),
		OriToken: "",
	}
	transferDataAbi, err := transferData.ABIPack()
	suite.Require().NoError(err)

	packetCallData := &packettypes.CallData{
		ContractAddress: syscontracts.AgentContractAddress,
		CallData:        callData,
	}
	callDataAbi, err := packetCallData.ABIPack()
	suite.Require().NoError(err)

	packet := packettypes.Packet{
		SrcChain:        suite.chainA.ChainID,
		DstChain:        suite.chainB.ChainID,
		Sequence:        1,
		Sender:          strings.ToLower(suite.chainA.SenderAddress.String()),
		TransferData:    transferDataAbi,
		CallData:        callDataAbi,
		CallbackAddress: common.BigToAddress(big.NewInt(0)).String(),
		FeeOption:       0,
	}
	result, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000001")
	suite.Require().NoError(err)
	ack := packettypes.NewAcknowledgement(
		0,
		result,
		"",
		strings.ToLower(suite.chainB.SenderAcc.String()),
		0,
	)
	ackData, err := ack.ABIPack()
	suite.Require().NoError(err)

	// relay packet
	err = pathAToB.RelayPacket(packet, ackData)
	suite.Require().NoError(err)
	packetFees := suite.GetPacketFees(suite.chainB, suite.chainB.ChainID, suite.chainC.ChainID, 1)
	suite.Equal(packetFees.TokenAddress, chainBERC20)
	suite.Equal(packetFees.Amount.Int64(), int64(1000))
	outTokens := suite.OutTokens(suite.chainB, chainBERC20, suite.chainC.ChainID)
	suite.Equal(outTokens.Int64(), int64(1000))

	balance = suite.ERC20Balance(suite.chainB, chainBERC20, agentcontract.AgentContractAddress)
	suite.Require().Equal(int64(0), balance.Int64())
	bindings := suite.Bindings(suite.chainB, chainBERC20, suite.chainA.ChainID)
	suite.Equal(bindings.OriToken, strings.ToLower(chainAERC20.String()))
	suite.Equal(bindings.Amount.Int64(), int64(2000))
	suite.Equal(bindings.Bound, true)
	suite.Equal(bindings.Scale, uint8(0))
	latestPacket := suite.GetLatestPacket(suite.chainB)
	suite.Equal(latestPacket.DstChain, packet.DstChain)
	suite.Equal(latestPacket.SrcChain, packet.SrcChain)
	suite.Equal(latestPacket.Sequence, packet.Sequence)
	suite.Equal(latestPacket.Sender, packet.Sender)
	suite.Equal(latestPacket.TransferData, packet.TransferData)
	suite.Equal(latestPacket.CallData, packet.CallData)
	suite.Equal(latestPacket.FeeOption, packet.FeeOption)
	suite.Equal(latestPacket.CallbackAddress, packet.CallbackAddress)

	status := suite.GetAckStatus(suite.chainA, suite.chainB.ChainID, 1)
	suite.Require().Equal(status, uint8(1))
	ackResult := suite.GetAck(suite.chainA, suite.chainB.ChainID, 1)
	suite.Require().Equal(ackResult.Code, ack.Code)
	suite.Require().Equal(ackResult.Result, ack.Result)
	suite.Require().Equal(ackResult.Message, ack.Message)
	suite.Require().Equal(ackResult.Relayer, ack.Relayer)

	hexAmount, err = hex.DecodeString("00000000000000000000000000000000000000000000000000000000000003e8")
	suite.Require().NoError(err)
	transferData = packettypes.TransferData{
		Receiver: strings.ToLower(suite.chainA.SenderAddress.String()),
		Amount:   hexAmount,
		Token:    strings.ToLower(chainBERC20.String()),
		OriToken: "",
	}
	transferDataAbi, err = transferData.ABIPack()
	suite.Require().NoError(err)

	packet = packettypes.Packet{
		SrcChain:        suite.chainB.ChainID,
		DstChain:        suite.chainC.ChainID,
		Sequence:        1,
		Sender:          strings.ToLower(agentcontract.AgentContractAddress.String()),
		TransferData:    transferDataAbi,
		CallData:        []byte(""),
		CallbackAddress: strings.ToLower(agentcontract.AgentContractAddress.String()),
		FeeOption:       0,
	}

	ack = packettypes.NewAcknowledgement(
		0,
		[]byte(""),
		"",
		strings.ToLower(suite.chainC.SenderAcc.String()),
		0,
	)
	ackData, err = ack.ABIPack()
	suite.Require().NoError(err)

	pathBToC.RegisterRelayers()
	// relay packet
	err = pathBToC.RelayPacket(packet, ackData)
	suite.Require().NoError(err)
	balance = suite.ERC20Balance(suite.chainC, chainCERC20, suite.chainA.SenderAddress)
	suite.Require().Equal(balance.Int64(), int64(1000))
	bindings = suite.Bindings(suite.chainC, chainCERC20, suite.chainB.ChainID)
	suite.Equal(bindings.OriToken, strings.ToLower(chainBERC20.String()))
	suite.Equal(bindings.Amount.Int64(), int64(1000))
	suite.Equal(bindings.Bound, true)
	suite.Equal(bindings.Scale, uint8(0))
	latestPacket = suite.GetLatestPacket(suite.chainC)
	suite.Equal(latestPacket.DstChain, packet.DstChain)
	suite.Equal(latestPacket.SrcChain, packet.SrcChain)
	suite.Equal(latestPacket.Sequence, packet.Sequence)
	suite.Equal(latestPacket.Sender, packet.Sender)
	suite.Equal(latestPacket.TransferData, packet.TransferData)
	suite.Equal(latestPacket.CallData, packet.CallData)
	suite.Equal(latestPacket.FeeOption, packet.FeeOption)
	suite.Equal(latestPacket.CallbackAddress, packet.CallbackAddress)
}

// ================================================================================================================
// CrossChain functions
// ================================================================================================================
func (suite *XIBCTestSuite) CrossChainCall(fromChain *xibctesting.TestChain, data packettypes.CrossChainData, fee packettypes.Fee) {
	crossChainCallData, err := endpointcontract.EndpointContract.ABI.Pack("crossChainCall", data, fee)
	suite.Require().NoError(err)
	amount := big.NewInt(0)
	if data.TokenAddress == common.HexToAddress("0x0000000000000000000000000000000000000000") {
		amount = amount.Add(amount, data.Amount)
	}

	if fee.TokenAddress == common.HexToAddress("0x0000000000000000000000000000000000000000") {
		amount = amount.Add(amount, fee.Amount)
	}

	_ = suite.SendTx(fromChain, endpointcontract.EndpointContractAddress, amount, crossChainCallData)
	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)
}

func (suite *XIBCTestSuite) OutTokens(
	fromChain *xibctesting.TestChain,
	tokenAddress common.Address,
	dstChain string,
) *big.Int {
	res, err := fromChain.App.XIBCKeeper.PacketKeeper.CallEVM(
		fromChain.GetContext(),
		endpointcontract.EndpointContract.ABI,
		aggregatetypes.ModuleAddress,
		endpointcontract.EndpointContractAddress,
		"outTokens",
		tokenAddress,
		dstChain,
	)
	suite.Require().NoError(err)

	type Amount struct {
		Value *big.Int
	}
	var amount Amount
	err = endpointcontract.EndpointContract.ABI.UnpackIntoInterface(&amount, "outTokens", res.Ret)
	suite.Require().NoError(err)

	return amount.Value
}

func (suite *XIBCTestSuite) Bindings(
	fromChain *xibctesting.TestChain,
	tokenAddress common.Address,
	oriChain string,
) packettypes.InToken {
	res, err := fromChain.App.XIBCKeeper.PacketKeeper.CallEVM(
		fromChain.GetContext(),
		endpointcontract.EndpointContract.ABI,
		aggregatetypes.ModuleAddress,
		endpointcontract.EndpointContractAddress,
		"bindings",
		strings.ToLower(tokenAddress.String())+"/"+oriChain,
	)
	suite.Require().NoError(err)

	var bind packettypes.InToken
	err = endpointcontract.EndpointContract.ABI.UnpackIntoInterface(&bind, "bindings", res.Ret)
	suite.Require().NoError(err)

	return bind
}

// ================================================================================================================
// Packet functions
// ================================================================================================================
func (suite *XIBCTestSuite) GetPacketFees(fromChain *xibctesting.TestChain, srcChain, dstChain string, sequence uint64) packettypes.Fee {
	data := []byte(dstChain + "/" + strconv.FormatUint(sequence, 10))
	packet := packetcontract.PacketContract.ABI

	res, err := fromChain.App.AggregateKeeper.CallEVM(
		fromChain.GetContext(),
		packet,
		packettypes.ModuleAddress,
		packetcontract.PacketContractAddress,
		"packetFees",
		data,
	)
	suite.Require().NoError(err)

	var fee packettypes.Fee
	err = packet.UnpackIntoInterface(&fee, "packetFees", res.Ret)
	suite.Require().NoError(err)

	return fee
}

func (suite *XIBCTestSuite) GetAck(fromChain *xibctesting.TestChain, dstChain string, sequence uint64) packettypes.Acknowledgement {
	data := []byte(dstChain + "/" + strconv.FormatUint(sequence, 10))
	packet := packetcontract.PacketContract.ABI

	res, err := fromChain.App.AggregateKeeper.CallEVM(
		fromChain.GetContext(),
		packet,
		packettypes.ModuleAddress,
		packetcontract.PacketContractAddress,
		"acks",
		data,
	)
	suite.Require().NoError(err)

	var ack packettypes.Acknowledgement
	err = packet.UnpackIntoInterface(&ack, "acks", res.Ret)
	suite.Require().NoError(err)
	return ack
}

func (suite *XIBCTestSuite) GetLatestPacket(fromChain *xibctesting.TestChain) packettypes.Packet {
	packetABI := packetcontract.PacketContract.ABI

	res, err := fromChain.App.AggregateKeeper.CallEVM(
		fromChain.GetContext(),
		packetABI,
		packettypes.ModuleAddress,
		packetcontract.PacketContractAddress,
		"latestPacket",
	)
	suite.Require().NoError(err)

	var packet packettypes.Packet
	err = packetABI.UnpackIntoInterface(&packet, "latestPacket", res.Ret)
	suite.Require().NoError(err)
	return packet
}

func (suite *XIBCTestSuite) GetAckStatus(fromChain *xibctesting.TestChain, dstChain string, sequence uint64) uint8 {
	packet := packetcontract.PacketContract.ABI

	res, err := fromChain.App.AggregateKeeper.CallEVM(
		fromChain.GetContext(),
		packet,
		packettypes.ModuleAddress,
		packetcontract.PacketContractAddress,
		"getAckStatus",
		dstChain,
		sequence,
	)
	suite.Require().NoError(err)

	var ackStatus uint8
	err = packet.UnpackIntoInterface(&ackStatus, "getAckStatus", res.Ret)
	suite.Require().NoError(err)
	return ackStatus
}

// ================================================================================================================
// ERC20 functions
// ================================================================================================================
func (suite *XIBCTestSuite) DeployERC20ByCrossChain(fromChain *xibctesting.TestChain) common.Address {
	ctorArgs, err := erc20contracts.ERC20MinterBurnerDecimalsContract.ABI.Pack("", "name", "symbol", uint8(18))
	suite.Require().NoError(err)

	data := make([]byte, len(erc20contracts.ERC20MinterBurnerDecimalsContract.Bin)+len(ctorArgs))
	copy(data[:len(erc20contracts.ERC20MinterBurnerDecimalsContract.Bin)], erc20contracts.ERC20MinterBurnerDecimalsContract.Bin)
	copy(data[len(erc20contracts.ERC20MinterBurnerDecimalsContract.Bin):], ctorArgs)

	nonce := fromChain.App.EvmKeeper.GetNonce(fromChain.GetContext(), endpointcontract.EndpointContractAddress)
	contractAddr := crypto.CreateAddress(endpointcontract.EndpointContractAddress, nonce)

	res, err := fromChain.App.AggregateKeeper.CallEVMWithData(fromChain.GetContext(), endpointcontract.EndpointContractAddress, nil, data)
	suite.Require().NoError(err)
	suite.Require().False(res.Failed(), res.VmError)

	return contractAddr
}

func (suite *XIBCTestSuite) GrantERC20MintRoleByCrossChain(fromChain *xibctesting.TestChain, erc20, address common.Address) {
	ctorArgs, err := erc20contracts.ERC20MinterBurnerDecimalsContract.ABI.Pack("grantRole", common.BytesToHash(crypto.Keccak256([]byte("MINTER_ROLE"))), address)
	suite.Require().NoError(err)

	res, err := fromChain.App.AggregateKeeper.CallEVMWithData(fromChain.GetContext(), endpointcontract.EndpointContractAddress, &erc20, ctorArgs)
	suite.Require().NoError(err)
	suite.Require().False(res.Failed(), res.VmError)
}

func (suite *XIBCTestSuite) MintERC20Token(fromChain *xibctesting.TestChain, to, erc20 common.Address, amount *big.Int) {
	ctorArgs, err := erc20contracts.ERC20MinterBurnerDecimalsContract.ABI.Pack("mint", to, amount)
	suite.Require().NoError(err)

	_ = suite.SendTx(fromChain, erc20, big.NewInt(0), ctorArgs)
	suite.coordinator.CommitBlock(fromChain)
}

func (suite *XIBCTestSuite) ERC20Balance(
	fromChain *xibctesting.TestChain,
	contract common.Address,
	account common.Address,
) *big.Int {
	erc20 := erc20contracts.ERC20MinterBurnerDecimalsContract.ABI

	res, err := fromChain.App.AggregateKeeper.CallEVM(
		fromChain.GetContext(),
		erc20,
		packettypes.ModuleAddress,
		contract,
		"balanceOf",
		account,
	)
	suite.Require().NoError(err)

	type Amount struct {
		Value *big.Int
	}
	var balance Amount
	err = erc20.UnpackIntoInterface(&balance, "balanceOf", res.Ret)
	suite.Require().NoError(err)

	return balance.Value
}

func (suite *XIBCTestSuite) Approve(
	fromChain *xibctesting.TestChain,
	erc20Address common.Address,
	spender common.Address,
	amount *big.Int,
) {
	transferData, err := erc20contracts.ERC20MinterBurnerDecimalsContract.ABI.Pack("approve", spender, amount)
	suite.Require().NoError(err)

	_ = suite.SendTx(fromChain, erc20Address, big.NewInt(0), transferData)
	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)
}

// ================================================================================================================
// EVM transaction
// ================================================================================================================
func (suite *XIBCTestSuite) SendTx(
	fromChain *xibctesting.TestChain,
	contractAddr common.Address,
	amount *big.Int,
	data []byte,
) *evm.MsgEthereumTx {
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
		data,
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
