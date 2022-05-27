package xibc_test

import (
	"math/big"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tharsis/ethermint/server/config"
	"github.com/tharsis/ethermint/tests"
	evm "github.com/tharsis/ethermint/x/evm/types"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	crosschaincontract "github.com/teleport-network/teleport/syscontracts/cross_chain"
	erc20contracts "github.com/teleport-network/teleport/syscontracts/erc20"
	packetcontract "github.com/teleport-network/teleport/syscontracts/xibc_packet"
	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
	xibctesting "github.com/teleport-network/teleport/x/xibc/testing"
)

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
	//suite.SetupTest()
	//
	//// setup testing conditions
	//pathAToB := xibctesting.NewPath(suite.chainA, suite.chainB)
	//suite.coordinator.SetupClients(pathAToB)
	//
	//// deploy ERC20
	//chainAERC20 := suite.DeployERC20ByCrossChain(suite.chainA)
	//chainBERC20 := suite.DeployERC20ByCrossChain(suite.chainB)
	//suite.GrantERC20MintRoleByCrossChain(suite.chainA, chainAERC20, suite.chainA.SenderAddress)
	//suite.MintERC20Token(suite.chainA, suite.chainA.SenderAddress, chainAERC20, big.NewInt(10000))
	//balance := suite.ERC20Balance(suite.chainA, chainAERC20, suite.chainA.SenderAddress)
	//suite.Equal(balance.Int64(), int64(10000))
	//
	//// add erc20 trace
	//err := suite.chainB.App.AggregateKeeper.RegisterERC20Trace(
	//	suite.chainB.GetContext(),
	//	chainBERC20,
	//	strings.ToLower(chainAERC20.String()),
	//	suite.chainA.ChainID,
	//	uint8(0),
	//)
	//suite.Require().NoError(err)
	//// check ERC20 trace
	//token, amount, exist, err := suite.chainB.App.AggregateKeeper.QueryERC20Trace(
	//	suite.chainB.GetContext(),
	//	chainBERC20,
	//	suite.chainA.ChainID,
	//)
	//suite.Require().NoError(err)
	//suite.Require().True(exist)
	//suite.Equal(token, strings.ToLower(chainAERC20.String()))
	//suite.Equal(amount.Int64(), int64(0))
	//
	//crossChainData := packettypes.CrossChainData{
	//	DestChain:       suite.chainB.ChainID,
	//	TokenAddress:    chainAERC20,
	//	Receiver:        strings.ToLower(suite.chainB.SenderAddress.String()),
	//	Amount:          big.NewInt(1000),
	//	ContractAddress: "",
	//	CallData:        []byte(""),
	//	CallbackAddress: common.BigToAddress(big.NewInt(0)),
	//	FeeOption:       0,
	//}
	//fee := packettypes.Fee{
	//	TokenAddress: chainAERC20,
	//	Amount:       big.NewInt(1000),
	//}
	//
	//// send CrossChainCAll Tx
	//suite.Approve(suite.chainA, chainAERC20, crosschaincontract.CrossChainAddress, big.NewInt(2000))
	//suite.CrossChainCAll(suite.chainA, crossChainData, fee)
	//balance = suite.ERC20Balance(suite.chainA, chainAERC20, suite.chainA.SenderAddress)
	//suite.Equal(balance.Int64(), int64(8000))
	//balance = suite.ERC20Balance(suite.chainA, chainAERC20, crosschaincontract.CrossChainAddress)
	//suite.Equal(balance.Int64(), int64(1000))
	//balance = suite.ERC20Balance(suite.chainA, chainAERC20, packetcontract.PacketContractAddress)
	//suite.Equal(balance.Int64(), int64(1000))
	//// check packet fees
	//
	//fees := suite.GetPacketFees(suite.chainA, suite.chainA.ChainID, suite.chainB.ChainID, 1)
	//suite.Require().Equal(chainAERC20, fees.TokenAddress)
	//suite.Require().Equal(big.NewInt(1000).Int64(), fees.Amount.Int64())
	//
	//// packet and ack
	//decodeString, err := hex.DecodeString("00000000000000000000000000000000000000000000000000000000000003e8")
	//suite.Require().NoError(err)
	//transferData := packettypes.TransferData{
	//	Receiver: strings.ToLower(crossChainData.Receiver),
	//	Amount:   decodeString,
	//	Token:    strings.ToLower(crossChainData.TokenAddress.String()),
	//	OriToken: "",
	//}
	//transferDataAbi, err := transferData.AbiPack()
	//suite.Require().NoError(err)
	//
	//packet := packettypes.Packet{
	//	SourceChain:      suite.chainA.ChainID,
	//	DestinationChain: suite.chainB.ChainID,
	//	Sequence:         1,
	//	Sender:           strings.ToLower(suite.chainA.SenderAddress.String()),
	//	TransferData:     transferDataAbi,
	//	CallData:         []byte(""),
	//	CallbackAddress:  common.BigToAddress(big.NewInt(0)).String(),
	//	FeeOption:        0,
	//}
	//ack := packettypes.NewResultAcknowledgement(
	//	0,
	//	[]byte(""),
	//	"",
	//	strings.ToLower(suite.chainB.SenderAcc.String()),
	//)
	//ackData, err := ack.AbiPack()
	//suite.Require().NoError(err)
	//
	//// relay
	//err = pathAToB.RelayPacket(packet, ackData)
	//suite.Require().NoError(err)
	//balance = suite.ERC20Balance(suite.chainB, chainBERC20, suite.chainB.SenderAddress)
	//suite.Require().Equal(int64(1000), balance.Int64())
	//balance = suite.ERC20Balance(suite.chainA, chainAERC20, packetcontract.PacketContractAddress)
	//suite.Equal(balance.Int64(), int64(0))
	//balance = suite.ERC20Balance(suite.chainA, chainAERC20, suite.chainA.SenderAddress)
	//suite.Equal(balance.Int64(), int64(9000))
	//
	//status := suite.GetAckStatus(suite.chainA, suite.chainA.ChainID, suite.chainB.ChainID, 1)
	//suite.Require().Equal(status, uint8(1))
	//ack = suite.GetAck(suite.chainA, suite.chainA.ChainID, suite.chainB.ChainID, 1)
	//suite.Require().Equal(uint64(0), ack.Code)
	//suite.Require().Equal([]byte(""), ack.Result)
	//suite.Require().Equal("", ack.Message)
	//suite.Require().Equal(suite.chainB.SenderAcc.String(), ack.Relayer)
}

// ================================================================================================================
// CrossChain functions
// ================================================================================================================
func (suite *XIBCTestSuite) CrossChainCAll(fromChain *xibctesting.TestChain, data packettypes.CrossChainData, fee packettypes.Fee) {
	crossChainCallData, err := crosschaincontract.CrossChainContract.ABI.Pack("crossChainCall", data, fee)
	suite.Require().NoError(err)
	amount := big.NewInt(0)
	if data.TokenAddress == common.HexToAddress("0x0000000000000000000000000000000000000000") {
		amount = amount.Add(amount, data.Amount)
	}

	if fee.TokenAddress == common.HexToAddress("0x0000000000000000000000000000000000000000") {
		amount = amount.Add(amount, fee.Amount)
	}

	_ = suite.SendTx(fromChain, crosschaincontract.CrossChainAddress, amount, crossChainCallData)
	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)
}

// ================================================================================================================
// Packet functions
// ================================================================================================================
func (suite *XIBCTestSuite) GetPacketFees(fromChain *xibctesting.TestChain, srcChain, destChain string, sequence uint64) packettypes.Fee {
	data := []byte(srcChain + "/" + destChain + "/" + strconv.FormatUint(sequence, 10))
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

func (suite *XIBCTestSuite) GetAck(fromChain *xibctesting.TestChain, srcChain, destChain string, sequence uint64) packettypes.Acknowledgement {
	data := []byte(srcChain + "/" + destChain + "/" + strconv.FormatUint(sequence, 10))
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

func (suite *XIBCTestSuite) GetAckStatus(fromChain *xibctesting.TestChain, srcChain, destChain string, sequence uint64) uint8 {
	packet := packetcontract.PacketContract.ABI

	res, err := fromChain.App.AggregateKeeper.CallEVM(
		fromChain.GetContext(),
		packet,
		packettypes.ModuleAddress,
		packetcontract.PacketContractAddress,
		"getAckStatus",
		srcChain,
		destChain,
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

	nonce := fromChain.App.EvmKeeper.GetNonce(fromChain.GetContext(), crosschaincontract.CrossChainAddress)
	contractAddr := crypto.CreateAddress(crosschaincontract.CrossChainAddress, nonce)

	res, err := fromChain.App.AggregateKeeper.CallEVMWithData(fromChain.GetContext(), crosschaincontract.CrossChainAddress, nil, data)
	suite.Require().NoError(err)
	suite.Require().False(res.Failed(), res.VmError)

	return contractAddr
}

func (suite *XIBCTestSuite) GrantERC20MintRoleByCrossChain(fromChain *xibctesting.TestChain, erc20, address common.Address) {
	ctorArgs, err := erc20contracts.ERC20MinterBurnerDecimalsContract.ABI.Pack("grantRole", common.BytesToHash(crypto.Keccak256([]byte("MINTER_ROLE"))), address)
	suite.Require().NoError(err)

	res, err := fromChain.App.AggregateKeeper.CallEVMWithData(fromChain.GetContext(), crosschaincontract.CrossChainAddress, &erc20, ctorArgs)
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
