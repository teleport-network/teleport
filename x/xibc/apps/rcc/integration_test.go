package rcc_test

// import (
// 	"crypto/sha256"
// 	"encoding/hex"
// 	"math/big"
// 	"strings"
// 	"testing"

// 	"github.com/stretchr/testify/suite"

// 	sdk "github.com/cosmos/cosmos-sdk/types"

// 	"github.com/ethereum/go-ethereum/common"
// 	ethtypes "github.com/ethereum/go-ethereum/core/types"
// 	"github.com/ethereum/go-ethereum/crypto"

// 	"github.com/tharsis/ethermint/server/config"
// 	"github.com/tharsis/ethermint/tests"
// 	evm "github.com/tharsis/ethermint/x/evm/types"

// 	erc20contracts "github.com/teleport-network/teleport/syscontracts/erc20"
// 	rcccontract "github.com/teleport-network/teleport/syscontracts/xibc_rcc"
// 	"github.com/teleport-network/teleport/x/xibc/apps/rcc/types"
// 	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
// 	xibctesting "github.com/teleport-network/teleport/x/xibc/testing"
// )

// type RCCTestSuite struct {
// 	suite.Suite
// 	coordinator *xibctesting.Coordinator
// 	chainA      *xibctesting.TestChain
// 	chainB      *xibctesting.TestChain
// }

// func TestRCCTestSuite(t *testing.T) {
// 	suite.Run(t, new(RCCTestSuite))
// }

// func (suite *RCCTestSuite) SetupTest() {
// 	suite.coordinator = xibctesting.NewCoordinator(suite.T(), 2)
// 	suite.chainA = suite.coordinator.GetChain(xibctesting.GetChainID(0))
// 	suite.chainB = suite.coordinator.GetChain(xibctesting.GetChainID(1))
// }

// func (suite *RCCTestSuite) TestRemoteContractCall() {
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
// 	erc20Adress := suite.DeployERC20(suite.chainB, rcccontract.RCCContractAddress, uint8(18))

// 	// commit block
// 	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)

// 	// send remote contract call
// 	payload, err := erc20contracts.ERC20MinterBurnerDecimalsContract.ABI.Pack("approve", suite.chainA.SenderAddress, amount)
// 	suite.Require().NoError(err)

// 	data := types.CallRCCData{
// 		ContractAddress: strings.ToLower(erc20Adress.String()),
// 		Data:            payload,
// 		DestChain:       path.EndpointB.ChainName,
// 		RelayChain:      "",
// 	}
// 	suite.SendRemoteContractCall(suite.chainA, data)

// 	// commit block
// 	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)
// 	sequence := uint64(1)
// 	// relay packet
// 	packetData := types.NewRCCPacketData(
// 		path.EndpointA.ChainName,
// 		path.EndpointB.ChainName,
// 		sequence,
// 		strings.ToLower(suite.chainA.SenderAddress.String()),
// 		strings.ToLower(erc20Adress.String()),
// 		payload,
// 	)
// 	bz, err := packetData.GetBytes()
// 	suite.NoError(err)
// 	packet := packettypes.NewPacket(
// 		sequence,
// 		path.EndpointA.ChainName,
// 		path.EndpointB.ChainName,
// 		"",
// 		[]string{types.PortID},
// 		[][]byte{bz},
// 	)

// 	ackBZ := suite.RCCAcks(suite.chainA, sha256.Sum256(bz))
// 	suite.Require().Equal([]byte{}, ackBZ)

// 	result, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000001")
// 	suite.Require().NoError(err)

// 	ack, err := packettypes.NewResultAcknowledgement([][]byte{result}, "").GetBytes()
// 	suite.Require().NoError(err)
// 	err = path.RelayPacket(packet, ack)
// 	suite.Require().NoError(err)

// 	// check ack
// 	ackBZ = suite.RCCAcks(suite.chainA, sha256.Sum256(bz))
// 	suite.Require().Equal(result, ackBZ)

// 	// check allowance
// 	allowance := suite.Allowance(suite.chainB, erc20Adress, rcccontract.RCCContractAddress, suite.chainA.SenderAddress)
// 	suite.Require().Equal(amount.String(), allowance.String())
// }

// // ================================================================================================================
// // Functions for step
// // ================================================================================================================

// func (suite *RCCTestSuite) SendRemoteContractCall(fromChain *xibctesting.TestChain, data types.CallRCCData) {
// 	rccData, err := rcccontract.RCCContract.ABI.Pack("sendRemoteContractCall", data)
// 	suite.Require().NoError(err)

// 	_ = suite.SendTx(fromChain, rcccontract.RCCContractAddress, big.NewInt(0), rccData)
// }

// func (suite *RCCTestSuite) DeployERC20(fromChain *xibctesting.TestChain, deployer common.Address, scale uint8) common.Address {
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

// func (suite *RCCTestSuite) Allowance(
// 	fromChain *xibctesting.TestChain,
// 	contract common.Address,
// 	owner common.Address,
// 	spender common.Address,
// ) *big.Int {
// 	erc20 := erc20contracts.ERC20MinterBurnerDecimalsContract.ABI

// 	res, err := fromChain.App.XIBCTransferKeeper.CallEVM(
// 		fromChain.GetContext(),
// 		erc20,
// 		types.ModuleAddress,
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

// func (suite *RCCTestSuite) RCCAcks(fromChain *xibctesting.TestChain, hash [32]byte) []byte {
// 	rcc := rcccontract.RCCContract.ABI

// 	res, err := fromChain.App.XIBCTransferKeeper.CallEVM(
// 		fromChain.GetContext(),
// 		rcc,
// 		types.ModuleAddress,
// 		rcccontract.RCCContractAddress,
// 		"acks",
// 		hash,
// 	)
// 	suite.Require().NoError(err)

// 	var ack struct{ Value []byte }
// 	err = rcc.UnpackIntoInterface(&ack, "acks", res.Ret)
// 	suite.Require().NoError(err)

// 	return ack.Value
// }

// // ================================================================================================================
// // EVM transaction (return events)
// // ================================================================================================================

// func (suite *RCCTestSuite) SendTx(fromChain *xibctesting.TestChain, contractAddr common.Address, amount *big.Int, transferData []byte) *evm.MsgEthereumTx {
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
