package transfer_test

import (
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	agentcontract "github.com/teleport-network/teleport/syscontracts/agent"
	wtelecontract "github.com/teleport-network/teleport/syscontracts/wtele"
	transfercontract "github.com/teleport-network/teleport/syscontracts/xibc_transfer"
	multicalltypes "github.com/teleport-network/teleport/x/xibc/apps/multicall/types"
	rcctypes "github.com/teleport-network/teleport/x/xibc/apps/rcc/types"
	"github.com/teleport-network/teleport/x/xibc/apps/transfer/types"
	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
	xibctesting "github.com/teleport-network/teleport/x/xibc/testing"
)

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
	suite.Require().NoError(err)

	rccDataBytes, err := abi.Arguments{{Type: multicalltypes.TupleRCCData}}.Pack(
		multicalltypes.RCCData{
			ContractAddress: strings.ToLower(agentcontract.AgentContractAddress.String()),
			Data:            agentPayload,
		},
	)
	suite.Require().NoError(err)

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
	suite.Require().NoError(err)

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
