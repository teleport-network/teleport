package keeper_test

import (
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmversion "github.com/tendermint/tendermint/proto/tendermint/version"
	"github.com/tendermint/tendermint/version"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"

	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
	"github.com/tharsis/ethermint/encoding"
	"github.com/tharsis/ethermint/server/config"
	"github.com/tharsis/ethermint/tests"
	ethermint "github.com/tharsis/ethermint/types"
	evm "github.com/tharsis/ethermint/x/evm/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
	feemarkettypes "github.com/tharsis/ethermint/x/feemarket/types"

	"github.com/teleport-network/teleport/app"
	erc20contracts "github.com/teleport-network/teleport/syscontracts/erc20"
	transfer "github.com/teleport-network/teleport/syscontracts/xibc_transfer"
	"github.com/teleport-network/teleport/x/xibc/apps/transfer/types"
	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx              sdk.Context
	app              *app.Teleport
	queryClientEvm   evm.QueryClient
	address          common.Address
	consAddress      sdk.ConsAddress
	clientCtx        client.Context
	ethSigner        ethtypes.Signer
	signer           keyring.Signer
	mintFeeCollector bool
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.DoSetupTest(suite.T())

	println(common.Address(common.BytesToAddress(authtypes.NewModuleAddress("FT").Bytes())).String())
	println(common.Address(common.BytesToAddress(authtypes.NewModuleAddress("packet").Bytes())).String())
}

func (suite *KeeperTestSuite) DoSetupTest(t require.TestingT) {
	checkTx := false

	// account key
	priv, err := ethsecp256k1.GenerateKey()
	require.NoError(t, err)
	suite.address = common.BytesToAddress(priv.PubKey().Address().Bytes())
	suite.signer = tests.NewSigner(priv)

	// consensus key
	priv, err = ethsecp256k1.GenerateKey()
	require.NoError(t, err)
	suite.consAddress = sdk.ConsAddress(priv.PubKey().Address())

	// setup feemarketGenesis params
	feemarketGenesis := feemarkettypes.DefaultGenesisState()
	feemarketGenesis.Params.EnableHeight = 1
	feemarketGenesis.Params.NoBaseFee = false
	feemarketGenesis.BaseFee = sdk.NewInt(feemarketGenesis.Params.InitialBaseFee)
	suite.app = app.Setup(checkTx, feemarketGenesis)

	if suite.mintFeeCollector {
		// mint some coin to fee collector
		coins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(int64(params.TxGas)-1)))
		genesisState := app.ModuleBasics.DefaultGenesis(suite.app.AppCodec())
		balances := []banktypes.Balance{
			{
				Address: suite.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName).String(),
				Coins:   coins,
			},
		}
		// update total supply
		bankGenesis := banktypes.NewGenesisState(
			banktypes.DefaultGenesisState().Params,
			balances,
			sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt((int64(params.TxGas)-1)))),
			[]banktypes.Metadata{},
		)
		bz := suite.app.AppCodec().MustMarshalJSON(bankGenesis)
		require.NotNil(t, bz)
		genesisState[banktypes.ModuleName] = suite.app.AppCodec().MustMarshalJSON(bankGenesis)

		// we marshal the genesisState of all module to a byte array
		stateBytes, err := tmjson.MarshalIndent(genesisState, "", " ")
		require.NoError(t, err)

		//Initialize the chain
		suite.app.InitChain(
			abci.RequestInitChain{
				ChainId:         "teleport_9000-1",
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simapp.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	suite.ctx = suite.app.BaseApp.NewContext(checkTx, tmproto.Header{
		Height:          1,
		ChainID:         "teleport_9000-1",
		Time:            time.Now().UTC(),
		ProposerAddress: suite.consAddress.Bytes(),

		Version: tmversion.Consensus{
			Block: version.BlockProtocol,
		},
		LastBlockId: tmproto.BlockID{
			Hash: tmhash.Sum([]byte("block_id")),
			PartSetHeader: tmproto.PartSetHeader{
				Total: 11,
				Hash:  tmhash.Sum([]byte("partset_header")),
			},
		},
		AppHash:            tmhash.Sum([]byte("app")),
		DataHash:           tmhash.Sum([]byte("data")),
		EvidenceHash:       tmhash.Sum([]byte("evidence")),
		ValidatorsHash:     tmhash.Sum([]byte("validators")),
		NextValidatorsHash: tmhash.Sum([]byte("next_validators")),
		ConsensusHash:      tmhash.Sum([]byte("consensus")),
		LastResultsHash:    tmhash.Sum([]byte("last_result")),
	})

	queryHelperEvm := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	evm.RegisterQueryServer(queryHelperEvm, suite.app.EvmKeeper)
	suite.queryClientEvm = evm.NewQueryClient(queryHelperEvm)

	acc := &ethermint.EthAccount{
		BaseAccount: authtypes.NewBaseAccount(sdk.AccAddress(suite.address.Bytes()), nil, 0, 0),
		CodeHash:    common.BytesToHash(crypto.Keccak256(nil)).String(),
	}

	suite.app.AccountKeeper.SetAccount(suite.ctx, acc)

	valAddr := sdk.ValAddress(suite.address.Bytes())
	validator, err := stakingtypes.NewValidator(valAddr, priv.PubKey(), stakingtypes.Description{})
	require.NoError(t, err)
	err = suite.app.StakingKeeper.SetValidatorByConsAddr(suite.ctx, validator)
	require.NoError(t, err)
	suite.app.StakingKeeper.SetValidator(suite.ctx, validator)

	encodingConfig := encoding.MakeConfig(app.ModuleBasics)
	suite.clientCtx = client.Context{}.WithTxConfig(encodingConfig.TxConfig)
	suite.ethSigner = ethtypes.LatestSignerForChainID(suite.app.EvmKeeper.ChainID())
}

func (suite *KeeperTestSuite) MintFeeCollector(coins sdk.Coins) {
	err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToModule(suite.ctx, types.ModuleName, authtypes.FeeCollectorName, coins)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) Commit() {
	_ = suite.app.Commit()
	header := suite.ctx.BlockHeader()
	header.Height += 1
	suite.app.BeginBlock(abci.RequestBeginBlock{
		Header: header,
	})

	// update ctx
	suite.ctx = suite.app.BaseApp.NewContext(false, header)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	evm.RegisterQueryServer(queryHelper, suite.app.EvmKeeper)
	suite.queryClientEvm = evm.NewQueryClient(queryHelper)
}

// ================================================================================================================
// EVM transaction (return events)
// ================================================================================================================

func (suite *KeeperTestSuite) SendTx(contractAddr common.Address, transferData []byte) *evm.MsgEthereumTx {
	ctx := sdk.WrapSDKContext(suite.ctx)
	chainID := suite.app.EvmKeeper.ChainID()

	args, err := json.Marshal(&evm.TransactionArgs{To: &contractAddr, From: &suite.address, Data: (*hexutil.Bytes)(&transferData)})
	suite.Require().NoError(err)
	res, err := suite.queryClientEvm.EstimateGas(ctx, &evm.EthCallRequest{
		Args:   args,
		GasCap: uint64(config.DefaultGasCap),
	})
	suite.Require().NoError(err)

	nonce := suite.app.EvmKeeper.GetNonce(suite.ctx, suite.address)

	// Mint the max gas to the FeeCollector to ensure balance in case of refund
	suite.MintFeeCollector(sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(suite.app.FeeMarketKeeper.GetBaseFee(suite.ctx).Int64()*int64(res.Gas)))))

	ercTransferTx := evm.NewTx(
		chainID,
		nonce,
		&contractAddr,
		big.NewInt(0),
		res.Gas,
		big.NewInt(0),
		big.NewInt(0),
		big.NewInt(0),
		transferData,
		&ethtypes.AccessList{}, // accesses
	)

	ercTransferTx.From = suite.address.Hex()
	err = ercTransferTx.Sign(ethtypes.LatestSignerForChainID(chainID), suite.signer)
	suite.Require().NoError(err)
	rsp, err := suite.app.EvmKeeper.EthereumTx(ctx, ercTransferTx)
	suite.Require().NoError(err)
	suite.Require().Empty(rsp.VmError)
	return ercTransferTx
}

// ================================================================================================================
// ERC20 contract
// ================================================================================================================

func (suite *KeeperTestSuite) DeployERC20MintableContract(sender common.Address, name string, symbol string, decimal uint8) common.Address {
	ctorArgs, err := erc20contracts.ERC20MinterBurnerDecimalsContract.ABI.Pack("", name, symbol, decimal)
	suite.Require().NoError(err)

	data := make([]byte, len(erc20contracts.ERC20MinterBurnerDecimalsContract.Bin)+len(ctorArgs))
	copy(data[:len(erc20contracts.ERC20MinterBurnerDecimalsContract.Bin)], erc20contracts.ERC20MinterBurnerDecimalsContract.Bin)
	copy(data[len(erc20contracts.ERC20MinterBurnerDecimalsContract.Bin):], ctorArgs)

	nonce := suite.app.EvmKeeper.GetNonce(suite.ctx, sender)
	contractAddr := crypto.CreateAddress(sender, nonce)

	res, err := suite.app.XIBCTransferKeeper.CallEVMWithData(suite.ctx, sender, nil, data)
	suite.Require().NoError(err)
	suite.Require().False(res.Failed(), res.VmError)

	return contractAddr
}

func (suite *KeeperTestSuite) TotalSupply(contract common.Address) *big.Int {
	erc20 := erc20contracts.ERC20MinterBurnerDecimalsContract.ABI

	res, err := suite.app.XIBCTransferKeeper.CallEVM(
		suite.ctx,
		erc20,
		types.ModuleAddress,
		contract,
		"totalSupply",
	)
	suite.Require().NoError(err)

	var supply types.Amount
	err = erc20.UnpackIntoInterface(&supply, "totalSupply", res.Ret)
	suite.Require().NoError(err)

	return supply.Value
}

func (suite *KeeperTestSuite) BalanceOf(contract common.Address, account common.Address) *big.Int {
	erc20 := erc20contracts.ERC20MinterBurnerDecimalsContract.ABI

	res, err := suite.app.XIBCTransferKeeper.CallEVM(
		suite.ctx,
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

func (suite *KeeperTestSuite) Mint(contract common.Address, sender common.Address, to common.Address, amount *big.Int) {
	erc20 := erc20contracts.ERC20MinterBurnerDecimalsContract.ABI

	res, err := suite.app.XIBCTransferKeeper.CallEVM(
		suite.ctx,
		erc20,
		sender,
		contract,
		"mint",
		to,
		amount,
	)
	suite.Require().NoError(err)
	suite.Require().False(res.Failed(), res.VmError)
}

func (suite *KeeperTestSuite) Approve(contract common.Address, sender common.Address, spender common.Address, amount *big.Int) {
	erc20 := erc20contracts.ERC20MinterBurnerDecimalsContract.ABI

	res, err := suite.app.XIBCTransferKeeper.CallEVM(
		suite.ctx,
		erc20,
		sender,
		contract,
		"approve",
		spender,
		amount,
	)
	suite.Require().NoError(err)
	suite.Require().False(res.Failed(), res.VmError)
}

func (suite *KeeperTestSuite) Allowance(contract common.Address, owner common.Address, spender common.Address) *big.Int {
	erc20 := erc20contracts.ERC20MinterBurnerDecimalsContract.ABI

	res, err := suite.app.XIBCTransferKeeper.CallEVM(
		suite.ctx,
		erc20,
		types.ModuleAddress,
		contract,
		"allowance",
		owner,
		spender,
	)
	suite.Require().NoError(err)

	var amount types.Amount
	err = erc20.UnpackIntoInterface(&amount, "allowance", res.Ret)
	suite.Require().NoError(err)

	return amount.Value
}

// ================================================================================================================
// XIBC transfer contract
// ================================================================================================================

func (suite *KeeperTestSuite) SendTransferERC20(sender common.Address, data types.ERC20TransferData) {
	res, err := suite.app.XIBCTransferKeeper.CallEVM(
		suite.ctx,
		transfer.TransferContract.ABI,
		sender,
		transfer.TransferContractAddress,
		"sendTransferERC20",
		data,
	)
	suite.Require().NoError(err)
	suite.Require().False(res.Failed(), res.VmError)
}

func (suite *KeeperTestSuite) SendTransferBase(sender common.Address, data types.BaseTransferData, amount *big.Int) {
	transferData, err := transfer.TransferContract.ABI.Pack("sendTransferBase", data)
	suite.Require().NoError(err)

	nonce, err := suite.app.AccountKeeper.GetSequence(suite.ctx, sender.Bytes())
	suite.Require().NoError(err)

	msg := ethtypes.NewMessage(
		sender,
		&transfer.TransferContractAddress,
		nonce,
		amount,                // amount
		config.DefaultGasCap,  // gasLimit
		big.NewInt(0),         // gasFeeCap
		big.NewInt(0),         // gasTipCap
		big.NewInt(0),         // gasPrice
		transferData,          // tx data
		ethtypes.AccessList{}, // accessList
		true,                  // checkNonce
	)

	res, err := suite.app.EvmKeeper.ApplyMessage(suite.ctx, msg, evmtypes.NewNoOpTracer(), true)
	suite.Require().NoError(err)
	suite.Require().False(res.Failed(), res.VmError)
}

func (suite *KeeperTestSuite) OutTokens(erc20Address common.Address, destChain string) *big.Int {
	res, err := suite.app.XIBCTransferKeeper.CallEVM(
		suite.ctx,
		transfer.TransferContract.ABI,
		types.ModuleAddress,
		transfer.TransferContractAddress,
		"outTokens",
		erc20Address,
		destChain,
	)
	suite.Require().NoError(err)

	var amount types.Amount
	err = transfer.TransferContract.ABI.UnpackIntoInterface(&amount, "outTokens", res.Ret)
	suite.Require().NoError(err)

	return amount.Value
}

func (suite *KeeperTestSuite) RecvPacket(data types.FungibleTokenPacketData) {
	res, err := suite.app.XIBCTransferKeeper.RecvPacket(suite.ctx, data)
	suite.Require().NoError(err)
	var r struct{ Result packettypes.Result }
	err = transfer.TransferContract.ABI.UnpackIntoInterface(&r, "onRecvPacket", res.Ret)
	suite.Require().NoError(err)
	suite.Require().Equal("", r.Result.Message)
	suite.Require().Equal([]byte{01}, r.Result.Result)
}

func (suite *KeeperTestSuite) AcknowledgementPacket(data types.FungibleTokenPacketData, result []byte) {
	res, err := suite.app.XIBCTransferKeeper.AcknowledgementPacket(suite.ctx, data, result)
	suite.Require().NoError(err)
	suite.Require().False(res.Failed(), res.VmError)
}
