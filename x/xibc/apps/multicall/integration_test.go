package multicall_test

// import (
// 	"encoding/hex"
// 	"math/big"
// 	"strings"
// 	"testing"

// 	"github.com/stretchr/testify/suite"

// 	sdk "github.com/cosmos/cosmos-sdk/types"

// 	"github.com/ethereum/go-ethereum/accounts/abi"
// 	"github.com/ethereum/go-ethereum/common"
// 	ethtypes "github.com/ethereum/go-ethereum/core/types"
// 	"github.com/ethereum/go-ethereum/crypto"

// 	"github.com/tharsis/ethermint/server/config"
// 	"github.com/tharsis/ethermint/tests"
// 	evm "github.com/tharsis/ethermint/x/evm/types"

// 	erc20contracts "github.com/teleport-network/teleport/syscontracts/erc20"
// 	multicallcontract "github.com/teleport-network/teleport/syscontracts/xibc_multicall"
// 	rcccontract "github.com/teleport-network/teleport/syscontracts/xibc_rcc"
// 	transfercontract "github.com/teleport-network/teleport/syscontracts/xibc_transfer"
// 	"github.com/teleport-network/teleport/x/xibc/apps/multicall/types"
// 	rcctypes "github.com/teleport-network/teleport/x/xibc/apps/rcc/types"
// 	transfertypes "github.com/teleport-network/teleport/x/xibc/apps/transfer/types"
// 	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
// 	xibctesting "github.com/teleport-network/teleport/x/xibc/testing"
// )

// type MultiCallTestSuite struct {
// 	suite.Suite
// 	coordinator *xibctesting.Coordinator
// 	chainA      *xibctesting.TestChain
// 	chainB      *xibctesting.TestChain
// }

// func TestMultiCallTestSuite(t *testing.T) {
// 	suite.Run(t, new(MultiCallTestSuite))
// }

// func (suite *MultiCallTestSuite) SetupTest() {
// 	suite.coordinator = xibctesting.NewCoordinator(suite.T(), 2)
// 	suite.chainA = suite.coordinator.GetChain(xibctesting.GetChainID(0))
// 	suite.chainB = suite.coordinator.GetChain(xibctesting.GetChainID(1))
// }

// func (suite *MultiCallTestSuite) TestTransferBaseCall() common.Address {
// 	path := xibctesting.NewPath(suite.chainA, suite.chainB)
// 	suite.coordinator.SetupClients(path)

// 	// prepare test data
// 	total := big.NewInt(100000000000000)
// 	amount := big.NewInt(100)

// 	// check balance
// 	balance := suite.chainA.App.EvmKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAddress)
// 	suite.Require().Equal(total.String(), balance.String())

// 	// deploy ERC20 on chainB
// 	erc20Address := suite.DeployERC20(suite.chainB, transfercontract.TransferContractAddress, uint8(18))

// 	// add erc20 trace on chainB
// 	err := suite.chainB.App.AggregateKeeper.RegisterERC20Trace(
// 		suite.chainB.GetContext(),
// 		erc20Address,
// 		common.BigToAddress(big.NewInt(0)).String(),
// 		suite.chainA.ChainID,
// 		uint8(0),
// 	)
// 	suite.Require().NoError(err)

// 	// commit block
// 	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)

// 	// check ERC20 trace
// 	_, _, exist, err := suite.chainB.App.AggregateKeeper.QueryERC20Trace(
// 		suite.chainB.GetContext(),
// 		erc20Address,
// 		suite.chainA.ChainID,
// 	)
// 	suite.Require().NoError(err)
// 	suite.Require().True(exist)

// 	transferBaseDataBytes, err := abi.Arguments{{Type: types.TupleBaseTransferData}}.Pack(
// 		types.BaseTransferData{
// 			Receiver: suite.chainB.SenderAddress.String(),
// 			Amount:   amount,
// 		},
// 	)
// 	suite.Require().NoError(err)

// 	data := types.MultiCallData{
// 		DestChain:  suite.chainB.ChainID,
// 		RelayChain: "",
// 		Functions:  []uint8{1},
// 		Data:       [][]byte{transferBaseDataBytes},
// 	}
// 	suite.SendMultiCall(suite.chainA, amount, data)

// 	// commit block
// 	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)

// 	// check balance
// 	balance = suite.chainA.App.EvmKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAddress)
// 	suite.Require().Equal(big.NewInt(0).Sub(total, amount).String(), balance.String())

// 	// check token out
// 	outAmount := suite.OutTokens(
// 		suite.chainA,
// 		common.BigToAddress(big.NewInt(0)),
// 		suite.chainB.ChainID,
// 	)
// 	suite.Require().Equal(amount.String(), outAmount.String())

// 	// relay packet
// 	transferBasePacketData := transfertypes.NewFungibleTokenPacketData(
// 		path.EndpointA.ChainName,
// 		path.EndpointB.ChainName,
// 		1,
// 		strings.ToLower(suite.chainA.SenderAddress.String()),
// 		strings.ToLower(suite.chainB.SenderAddress.String()),
// 		amount.Bytes(),
// 		strings.ToLower(common.BigToAddress(big.NewInt(0)).String()),
// 		strings.ToLower(""),
// 	)
// 	DataListBaseBz, err := transferBasePacketData.GetBytes()
// 	suite.NoError(err)
// 	packet := packettypes.NewPacket(
// 		1,
// 		path.EndpointA.ChainName,
// 		path.EndpointB.ChainName,
// 		"",
// 		[]string{transfertypes.PortID},
// 		[][]byte{DataListBaseBz},
// 	)

// 	ack, err := packettypes.NewResultAcknowledgement([][]byte{{byte(1)}}, "").GetBytes()
// 	suite.NoError(err)
// 	err = path.RelayPacket(packet, ack)
// 	suite.Require().NoError(err)

// 	// check balance
// 	recvBalance := suite.BalanceOf(suite.chainB, erc20Address, suite.chainB.SenderAddress)
// 	suite.Require().Equal(amount.String(), recvBalance.String())

// 	return erc20Address
// }

// func (suite *MultiCallTestSuite) TestTransferBaseBackCall() {
// 	path := xibctesting.NewPath(suite.chainA, suite.chainB)
// 	suite.coordinator.SetupClients(path)

// 	erc20Address := suite.TestTransferBaseCall()
// 	amount := big.NewInt(100)

// 	// check ERC20 trace
// 	_, _, exist, err := suite.chainB.App.AggregateKeeper.QueryERC20Trace(
// 		suite.chainB.GetContext(),
// 		erc20Address,
// 		suite.chainA.ChainID,
// 	)
// 	suite.Require().NoError(err)
// 	suite.Require().True(exist)

// 	// check balance
// 	recvBalance := suite.BalanceOf(suite.chainB, erc20Address, suite.chainB.SenderAddress)
// 	suite.Require().Equal(amount.String(), recvBalance.String())

// 	// Approve erc20 to transfer
// 	suite.Approve(suite.chainB, erc20Address, amount)

// 	// send multi call
// 	transferBaseBackDataBytes, err := abi.Arguments{{Type: types.TupleERC20TransferData}}.Pack(
// 		types.ERC20TransferData{
// 			TokenAddress: erc20Address,
// 			Receiver:     suite.chainA.SenderAddress.String(),
// 			Amount:       amount,
// 		},
// 	)
// 	suite.Require().NoError(err)

// 	data := types.MultiCallData{
// 		DestChain:  suite.chainA.ChainID,
// 		RelayChain: "",
// 		Functions:  []uint8{0},
// 		Data:       [][]byte{transferBaseBackDataBytes},
// 	}
// 	suite.SendMultiCall(suite.chainB, big.NewInt(0), data)

// 	// commit block
// 	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)

// 	// check balance
// 	erc20Balance := suite.BalanceOf(suite.chainB, erc20Address, suite.chainB.SenderAddress)
// 	suite.Require().Equal("0", erc20Balance.String())

// 	// relay packet
// 	packetData := transfertypes.NewFungibleTokenPacketData(
// 		path.EndpointB.ChainName,
// 		path.EndpointA.ChainName,
// 		1,
// 		strings.ToLower(suite.chainB.SenderAddress.String()),
// 		strings.ToLower(suite.chainA.SenderAddress.String()),
// 		amount.Bytes(),
// 		strings.ToLower(erc20Address.String()),
// 		strings.ToLower(common.BigToAddress(big.NewInt(0)).String()),
// 	)
// 	DataListERC20Bz, err := packetData.GetBytes()
// 	suite.NoError(err)
// 	packet := packettypes.NewPacket(
// 		1,
// 		path.EndpointB.ChainName,
// 		path.EndpointA.ChainName,
// 		"",
// 		[]string{transfertypes.PortID},
// 		[][]byte{DataListERC20Bz},
// 	)

// 	ack, err := packettypes.NewResultAcknowledgement([][]byte{{byte(1)}}, "").GetBytes()
// 	suite.NoError(err)
// 	err = path.RelayPacket(packet, ack)
// 	suite.Require().NoError(err)

// 	// check chainA token out
// 	outAmount := suite.OutTokens(
// 		suite.chainA,
// 		common.BigToAddress(big.NewInt(0)),
// 		suite.chainB.ChainID,
// 	)
// 	suite.Require().Equal("0", outAmount.String())
// }

// func (suite *MultiCallTestSuite) TestRCCCall() {
// 	path := xibctesting.NewPath(suite.chainA, suite.chainB)
// 	suite.coordinator.SetupClients(path)

// 	// prepare test data
// 	total := big.NewInt(100000000000000)
// 	amount := big.NewInt(100)

// 	// check balance
// 	balanceA := suite.chainA.App.EvmKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAddress)
// 	suite.Require().Equal(total.String(), balanceA.String())
// 	balanceB := suite.chainB.App.EvmKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAddress)
// 	suite.Require().Equal(total.String(), balanceB.String())

// 	// deploy ERC20 on chainB
// 	erc20Address := suite.DeployERC20(suite.chainB, transfercontract.TransferContractAddress, uint8(18))

// 	// commit block
// 	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)

// 	// send multi call
// 	payload, err := erc20contracts.ERC20MinterBurnerDecimalsContract.ABI.Pack("approve", suite.chainA.SenderAddress, amount)
// 	suite.Require().NoError(err)
// 	rccDataBytes, err := abi.Arguments{{Type: types.TupleRCCData}}.Pack(
// 		types.RCCData{
// 			ContractAddress: strings.ToLower(erc20Address.String()),
// 			Data:            payload,
// 		},
// 	)
// 	suite.Require().NoError(err)

// 	data := types.MultiCallData{
// 		DestChain:  path.EndpointB.ChainName,
// 		RelayChain: "",
// 		Functions:  []uint8{2},
// 		Data:       [][]byte{rccDataBytes},
// 	}
// 	suite.SendMultiCall(suite.chainA, amount, data)

// 	// commit block
// 	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)
// 	sequence := uint64(1)
// 	// relay packet
// 	rccPacketData := rcctypes.NewRCCPacketData(
// 		path.EndpointA.ChainName,
// 		path.EndpointB.ChainName,
// 		sequence,
// 		strings.ToLower(suite.chainA.SenderAddress.String()),
// 		strings.ToLower(erc20Address.String()),
// 		payload,
// 	)
// 	bz, err := rccPacketData.GetBytes()
// 	suite.NoError(err)
// 	packet := packettypes.NewPacket(
// 		sequence,
// 		path.EndpointA.ChainName,
// 		path.EndpointB.ChainName,
// 		"",
// 		[]string{rcctypes.PortID},
// 		[][]byte{bz},
// 	)

// 	resultRcc, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000001")
// 	suite.Require().NoError(err)

// 	ack, err := packettypes.NewResultAcknowledgement([][]byte{resultRcc}, "").GetBytes()
// 	suite.Require().NoError(err)
// 	err = path.RelayPacket(packet, ack)
// 	suite.Require().NoError(err)

// 	// check allowance
// 	allowance := suite.Allowance(suite.chainB, erc20Address, rcccontract.RCCContractAddress, suite.chainA.SenderAddress)
// 	suite.Require().Equal(amount.String(), allowance.String())
// }

// func (suite *MultiCallTestSuite) TestMultiCall_VV() {
// 	path := xibctesting.NewPath(suite.chainA, suite.chainB)
// 	suite.coordinator.SetupClients(path)

// 	// prepare test data
// 	total := big.NewInt(100000000000000)
// 	amount := big.NewInt(100)

// 	// check balance
// 	balanceA := suite.chainA.App.EvmKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAddress)
// 	suite.Require().Equal(total.String(), balanceA.String())
// 	balanceB := suite.chainB.App.EvmKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAddress)
// 	suite.Require().Equal(total.String(), balanceB.String())

// 	// deploy ERC20 on chainB
// 	erc20Address := suite.DeployERC20(suite.chainB, transfercontract.TransferContractAddress, uint8(18))

// 	// add erc20 trace on chainB
// 	err := suite.chainB.App.AggregateKeeper.RegisterERC20Trace(
// 		suite.chainB.GetContext(),
// 		erc20Address,
// 		common.BigToAddress(big.NewInt(0)).String(),
// 		suite.chainA.ChainID,
// 		uint8(0),
// 	)
// 	suite.Require().NoError(err)

// 	// commit block
// 	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)

// 	// check ERC20 trace
// 	_, _, exist, err := suite.chainB.App.AggregateKeeper.QueryERC20Trace(
// 		suite.chainB.GetContext(),
// 		erc20Address,
// 		suite.chainA.ChainID,
// 	)
// 	suite.Require().NoError(err)
// 	suite.Require().True(exist)

// 	// send multi call
// 	payload, err := erc20contracts.ERC20MinterBurnerDecimalsContract.ABI.Pack("approve", suite.chainA.SenderAddress, amount)
// 	suite.Require().NoError(err)

// 	rccDataBytes, err := abi.Arguments{{Type: types.TupleRCCData}}.Pack(
// 		types.RCCData{
// 			ContractAddress: strings.ToLower(erc20Address.String()),
// 			Data:            payload,
// 		},
// 	)
// 	suite.Require().NoError(err)

// 	transferBaseDataBytes, err := abi.Arguments{{Type: types.TupleBaseTransferData}}.Pack(
// 		types.BaseTransferData{
// 			Receiver: suite.chainB.SenderAddress.String(),
// 			Amount:   amount,
// 		},
// 	)
// 	suite.Require().NoError(err)

// 	data := types.MultiCallData{
// 		DestChain:  path.EndpointB.ChainName,
// 		RelayChain: "",
// 		Functions:  []uint8{2, 1},
// 		Data:       [][]byte{rccDataBytes, transferBaseDataBytes},
// 	}
// 	suite.SendMultiCall(suite.chainA, amount, data)

// 	// commit block
// 	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)
// 	sequence := uint64(1)
// 	// relay packet
// 	rccPacketData := rcctypes.NewRCCPacketData(
// 		path.EndpointA.ChainName,
// 		path.EndpointB.ChainName,
// 		sequence,
// 		strings.ToLower(suite.chainA.SenderAddress.String()),
// 		strings.ToLower(erc20Address.String()),
// 		payload,
// 	)
// 	transferBasePacketData := transfertypes.NewFungibleTokenPacketData(
// 		path.EndpointA.ChainName,
// 		path.EndpointB.ChainName,
// 		sequence,
// 		strings.ToLower(suite.chainA.SenderAddress.String()),
// 		strings.ToLower(suite.chainB.SenderAddress.String()),
// 		amount.Bytes(),
// 		strings.ToLower(common.BigToAddress(big.NewInt(0)).String()),
// 		strings.ToLower(""),
// 	)
// 	bz, err := rccPacketData.GetBytes()
// 	suite.NoError(err)
// 	DataListERC20Bz, err := transferBasePacketData.GetBytes()
// 	suite.NoError(err)
// 	packet := packettypes.NewPacket(
// 		1,
// 		path.EndpointA.ChainName,
// 		path.EndpointB.ChainName,
// 		"",
// 		[]string{
// 			rcctypes.PortID,
// 			transfertypes.PortID,
// 		},
// 		[][]byte{
// 			bz,
// 			DataListERC20Bz,
// 		},
// 	)

// 	resultRcc, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000001")
// 	suite.Require().NoError(err)
// 	resultTransferBase := []byte{byte(1)}

// 	ack, err := packettypes.NewResultAcknowledgement([][]byte{resultRcc, resultTransferBase}, "").GetBytes()
// 	suite.Require().NoError(err)
// 	err = path.RelayPacket(packet, ack)
// 	suite.Require().NoError(err)

// 	// check allowance
// 	allowance := suite.Allowance(suite.chainB, erc20Address, rcccontract.RCCContractAddress, suite.chainA.SenderAddress)
// 	suite.Require().Equal(amount.String(), allowance.String())

// 	// check balance
// 	recvBalance := suite.BalanceOf(suite.chainB, erc20Address, suite.chainB.SenderAddress)
// 	suite.Require().Equal(amount.String(), recvBalance.String())
// }

// func (suite *MultiCallTestSuite) TestMultiCall_VX() {
// 	path := xibctesting.NewPath(suite.chainA, suite.chainB)
// 	suite.coordinator.SetupClients(path)

// 	// prepare test data
// 	total := big.NewInt(100000000000000)
// 	amount := big.NewInt(100)

// 	// check balance
// 	balanceA := suite.chainA.App.EvmKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAddress)
// 	suite.Require().Equal(total.String(), balanceA.String())
// 	balanceB := suite.chainB.App.EvmKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAddress)
// 	suite.Require().Equal(total.String(), balanceB.String())

// 	// deploy ERC20 on chainB
// 	erc20Address := suite.DeployERC20(suite.chainB, transfercontract.TransferContractAddress, uint8(18))

// 	// commit block
// 	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)

// 	// check ERC20 trace
// 	_, _, exist, err := suite.chainB.App.AggregateKeeper.QueryERC20Trace(
// 		suite.chainB.GetContext(),
// 		erc20Address,
// 		suite.chainA.ChainID,
// 	)
// 	suite.Require().NoError(err)
// 	suite.Require().False(exist)

// 	// send multi call
// 	payload, err := erc20contracts.ERC20MinterBurnerDecimalsContract.ABI.Pack("approve", suite.chainA.SenderAddress, amount)
// 	suite.Require().NoError(err)

// 	rccDataBytes, err := abi.Arguments{{Type: types.TupleRCCData}}.Pack(
// 		types.RCCData{
// 			ContractAddress: strings.ToLower(erc20Address.String()),
// 			Data:            payload,
// 		},
// 	)
// 	suite.Require().NoError(err)

// 	transferBaseDataBytes, err := abi.Arguments{{Type: types.TupleBaseTransferData}}.Pack(
// 		types.BaseTransferData{
// 			Receiver: suite.chainB.SenderAddress.String(),
// 			Amount:   amount,
// 		},
// 	)
// 	suite.Require().NoError(err)

// 	data := types.MultiCallData{
// 		DestChain:  path.EndpointB.ChainName,
// 		RelayChain: "",
// 		Functions:  []uint8{2, 1},
// 		Data:       [][]byte{rccDataBytes, transferBaseDataBytes},
// 	}
// 	suite.SendMultiCall(suite.chainA, amount, data)

// 	// commit block
// 	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)
// 	sequence := uint64(1)

// 	// relay packet
// 	rccPacketData := rcctypes.NewRCCPacketData(
// 		path.EndpointA.ChainName,
// 		path.EndpointB.ChainName,
// 		sequence,
// 		strings.ToLower(suite.chainA.SenderAddress.String()),
// 		strings.ToLower(erc20Address.String()),
// 		payload,
// 	)
// 	transferBasePacketData := transfertypes.NewFungibleTokenPacketData(
// 		path.EndpointA.ChainName,
// 		path.EndpointB.ChainName,
// 		sequence,
// 		strings.ToLower(suite.chainA.SenderAddress.String()),
// 		strings.ToLower(suite.chainB.SenderAddress.String()),
// 		amount.Bytes(),
// 		strings.ToLower(common.BigToAddress(big.NewInt(0)).String()),
// 		strings.ToLower(""),
// 	)
// 	bz, err := rccPacketData.GetBytes()
// 	suite.NoError(err)
// 	DataListERC20Bz, err := transferBasePacketData.GetBytes()
// 	suite.NoError(err)
// 	packet := packettypes.NewPacket(
// 		1,
// 		path.EndpointA.ChainName,
// 		path.EndpointB.ChainName,
// 		"",
// 		[]string{
// 			rcctypes.PortID,
// 			transfertypes.PortID,
// 		},
// 		[][]byte{
// 			bz,
// 			DataListERC20Bz,
// 		},
// 	)

// 	ack, err := packettypes.NewErrorAcknowledgement("onRecvPackt: binding is not exist", "").GetBytes()
// 	suite.Require().NoError(err)

// 	err = path.RelayPacket(packet, ack)
// 	suite.Require().NoError(err)

// 	// check allowance
// 	allowance := suite.Allowance(suite.chainB, erc20Address, rcccontract.RCCContractAddress, suite.chainA.SenderAddress)
// 	suite.Require().Equal("0", allowance.String())

// 	// check balance
// 	recvBalance := suite.BalanceOf(suite.chainB, erc20Address, suite.chainB.SenderAddress)
// 	suite.Require().Equal("0", recvBalance.String())
// }

// // ================================================================================================================
// // Functions for step
// // ================================================================================================================

// func (suite *MultiCallTestSuite) SendMultiCall(fromChain *xibctesting.TestChain, amount *big.Int, data types.MultiCallData) {
// 	multiCallData, err := multicallcontract.MultiCallContract.ABI.Pack("multiCall", data)
// 	suite.Require().NoError(err)

// 	_ = suite.SendTx(fromChain, multicallcontract.MultiCallContractAddress, amount, multiCallData)
// }

// func (suite *MultiCallTestSuite) DeployERC20(fromChain *xibctesting.TestChain, deployer common.Address, scale uint8) common.Address {
// 	ctorArgs, err := erc20contracts.ERC20MinterBurnerDecimalsContract.ABI.Pack("", "name", "symbol", scale)
// 	suite.Require().NoError(err)

// 	data := make([]byte, len(erc20contracts.ERC20MinterBurnerDecimalsContract.Bin)+len(ctorArgs))
// 	copy(data[:len(erc20contracts.ERC20MinterBurnerDecimalsContract.Bin)], erc20contracts.ERC20MinterBurnerDecimalsContract.Bin)
// 	copy(data[len(erc20contracts.ERC20MinterBurnerDecimalsContract.Bin):], ctorArgs)

// 	nonce := fromChain.App.EvmKeeper.GetNonce(fromChain.GetContext(), deployer)
// 	contractAddr := crypto.CreateAddress(deployer, nonce)

// 	res, err := fromChain.App.XIBCTransferKeeper.CallEVMWithData(fromChain.GetContext(), deployer, nil, data)
// 	suite.Require().NoError(err)
// 	suite.Require().False(res.Failed(), res.VmError)

// 	return contractAddr
// }

// func (suite *MultiCallTestSuite) Allowance(
// 	fromChain *xibctesting.TestChain,
// 	contract common.Address,
// 	owner common.Address,
// 	spender common.Address,
// ) *big.Int {
// 	erc20 := erc20contracts.ERC20MinterBurnerDecimalsContract.ABI

// 	res, err := fromChain.App.XIBCTransferKeeper.CallEVM(
// 		fromChain.GetContext(),
// 		erc20,
// 		rcctypes.ModuleAddress,
// 		contract,
// 		"allowance",
// 		owner,
// 		spender,
// 	)
// 	suite.Require().NoError(err)

// 	var amount types.Amount
// 	err = erc20.UnpackIntoInterface(&amount, "allowance", res.Ret)
// 	suite.Require().NoError(err)

// 	return amount.Value
// }

// func (suite *MultiCallTestSuite) RCCAcks(fromChain *xibctesting.TestChain, hash [32]byte) []byte {
// 	rcc := multicallcontract.MultiCallContract.ABI

// 	res, err := fromChain.App.XIBCTransferKeeper.CallEVM(
// 		fromChain.GetContext(),
// 		rcc,
// 		rcctypes.ModuleAddress,
// 		multicallcontract.MultiCallContractAddress,
// 		"acks",
// 		hash,
// 	)
// 	suite.Require().NoError(err)

// 	var ack struct{ Value []byte }
// 	err = rcc.UnpackIntoInterface(&ack, "acks", res.Ret)
// 	suite.Require().NoError(err)

// 	return ack.Value
// }

// func (suite *MultiCallTestSuite) BalanceOf(fromChain *xibctesting.TestChain, contract common.Address, account common.Address) *big.Int {
// 	erc20 := erc20contracts.ERC20MinterBurnerDecimalsContract.ABI

// 	res, err := fromChain.App.XIBCTransferKeeper.CallEVM(
// 		fromChain.GetContext(),
// 		erc20,
// 		transfertypes.ModuleAddress,
// 		contract,
// 		"balanceOf",
// 		account,
// 	)
// 	suite.Require().NoError(err)

// 	var balance types.Amount
// 	err = erc20.UnpackIntoInterface(&balance, "balanceOf", res.Ret)
// 	suite.Require().NoError(err)

// 	return balance.Value
// }

// func (suite *MultiCallTestSuite) OutTokens(fromChain *xibctesting.TestChain, tokenAddress common.Address, destChain string) *big.Int {
// 	res, err := fromChain.App.XIBCTransferKeeper.CallEVM(
// 		fromChain.GetContext(),
// 		transfercontract.TransferContract.ABI,
// 		transfertypes.ModuleAddress,
// 		transfercontract.TransferContractAddress,
// 		"outTokens",
// 		tokenAddress,
// 		destChain,
// 	)
// 	suite.Require().NoError(err)

// 	var amount types.Amount
// 	err = transfercontract.TransferContract.ABI.UnpackIntoInterface(&amount, "outTokens", res.Ret)
// 	suite.Require().NoError(err)

// 	return amount.Value
// }

// func (suite *MultiCallTestSuite) Approve(fromChain *xibctesting.TestChain, erc20Address common.Address, amount *big.Int) {
// 	transferData, err := erc20contracts.ERC20MinterBurnerDecimalsContract.ABI.Pack("approve", transfercontract.TransferContractAddress, amount)
// 	suite.Require().NoError(err)

// 	_ = suite.SendTx(fromChain, erc20Address, big.NewInt(0), transferData)
// 	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)
// }

// // ================================================================================================================
// // EVM transaction (return events)
// // ================================================================================================================

// func (suite *MultiCallTestSuite) SendTx(fromChain *xibctesting.TestChain, contractAddr common.Address, amount *big.Int, transferData []byte) *evm.MsgEthereumTx {
// 	ctx := sdk.WrapSDKContext(fromChain.GetContext())
// 	chainID := fromChain.App.EvmKeeper.ChainID()
// 	signer := tests.NewSigner(fromChain.SenderPrivKey)

// 	nonce := fromChain.App.EvmKeeper.GetNonce(fromChain.GetContext(), fromChain.SenderAddress)
// 	ercTransferTx := evm.NewTx(
// 		chainID,
// 		nonce,
// 		&contractAddr,
// 		amount,
// 		config.DefaultGasCap,
// 		big.NewInt(0),
// 		big.NewInt(0),
// 		big.NewInt(0),
// 		transferData,
// 		&ethtypes.AccessList{}, // accesses
// 	)

// 	ercTransferTx.From = fromChain.SenderAddress.Hex()
// 	err := ercTransferTx.Sign(ethtypes.LatestSignerForChainID(chainID), signer)
// 	suite.Require().NoError(err)
// 	rsp, err := fromChain.App.EvmKeeper.EthereumTx(ctx, ercTransferTx)
// 	suite.Require().NoError(err)
// 	suite.Require().Empty(rsp.VmError, rsp.VmError)
// 	return ercTransferTx
// }
