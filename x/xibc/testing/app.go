package xibctesting

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/tharsis/ethermint/encoding"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
	feemarkettypes "github.com/tharsis/ethermint/x/feemarket/types"

	"github.com/teleport-network/teleport/app"
	xibcrcctypes "github.com/teleport-network/teleport/x/xibc/apps/rcc/types"
	xibctransfertypes "github.com/teleport-network/teleport/x/xibc/apps/transfer/types"
	xibcroutingtypes "github.com/teleport-network/teleport/x/xibc/core/routing/types"
	xibcmock "github.com/teleport-network/teleport/x/xibc/testing/mock"
)

var DefaultTestingAppInit func() (*app.Teleport, map[string]json.RawMessage) = SetupTestingApp

func SetupTestingApp() (*app.Teleport, map[string]json.RawMessage) {
	db := dbm.NewMemDB()

	encCdc := encoding.MakeConfig(app.ModuleBasics)
	teleport := app.NewTeleport(log.NewNopLogger(), db, nil, true, map[int64]bool{}, app.DefaultNodeHome, 5, encCdc, simapp.EmptyAppOptions{})

	xibcTransferModule, _ := teleport.XIBCKeeper.RoutingKeeper.Router.GetRoute(xibctransfertypes.PortID)
	xibcRCCModule, _ := teleport.XIBCKeeper.RoutingKeeper.Router.GetRoute(xibcrcctypes.PortID)

	mockModule := xibcmock.NewAppModule()

	xibcRouter := xibcroutingtypes.NewRouter()
	xibcRouter.AddRoute(xibcmock.ModuleName, mockModule)
	xibcRouter.AddRoute(xibctransfertypes.PortID, xibcTransferModule)
	xibcRouter.AddRoute(xibcrcctypes.PortID, xibcRCCModule)
	teleport.XIBCKeeper.RoutingKeeper.Router = xibcRouter
	teleport.XIBCKeeper.RoutingKeeper.Router.Seal()

	return teleport, app.NewDefaultGenesisState()
}

// SetupWithGenesisValSet initializes a new Teleport with a validator set and genesis accounts that also act as delegators.
func SetupWithGenesisValSet(
	t *testing.T,
	valSet *tmtypes.ValidatorSet,
	genAccs []authtypes.GenesisAccount,
	balances ...banktypes.Balance,
) *app.Teleport {
	teleport, genesisState := DefaultTestingAppInit()
	// set genesis accounts
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = teleport.AppCodec().MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	bondAmt := sdk.NewInt(1e16)

	for _, val := range valSet.Validators {
		pk, err := cryptocodec.FromTmPubKeyInterface(val.PubKey)
		require.NoError(t, err)
		pkAny, err := codectypes.NewAnyWithValue(pk)
		require.NoError(t, err)
		validator := stakingtypes.Validator{
			OperatorAddress:   sdk.ValAddress(val.Address).String(),
			ConsensusPubkey:   pkAny,
			Jailed:            false,
			Status:            stakingtypes.Bonded,
			Tokens:            bondAmt,
			DelegatorShares:   sdk.OneDec(),
			Description:       stakingtypes.Description{},
			UnbondingHeight:   int64(0),
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
			MinSelfDelegation: sdk.ZeroInt(),
		}
		validators = append(validators, validator)
		delegations = append(delegations, stakingtypes.NewDelegation(genAccs[0].GetAddress(), val.Address.Bytes(), sdk.OneDec()))
	}

	// set validators and delegations
	stakingGenesis := stakingtypes.NewGenesisState(stakingtypes.DefaultParams(), validators, delegations)
	genesisState[stakingtypes.ModuleName] = teleport.AppCodec().MustMarshalJSON(stakingGenesis)

	// set EVM genesis
	evmGenesis := evmtypes.DefaultGenesisState()
	evmGenesis.Params.EvmDenom = sdk.DefaultBondDenom
	genesisState[evmtypes.ModuleName] = teleport.AppCodec().MustMarshalJSON(evmGenesis)

	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		// add genesis acc tokens and delegated tokens to total supply
		totalSupply = totalSupply.Add(b.Coins.Add(sdk.NewCoin(sdk.DefaultBondDenom, bondAmt))...)
	}

	// add bonded amount to bonded pool module account
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, bondAmt)},
	})

	// update total supply
	bankGenesis := banktypes.NewGenesisState(
		banktypes.DefaultGenesisState().Params,
		balances,
		totalSupply,
		[]banktypes.Metadata{},
	)

	genesisState[banktypes.ModuleName] = teleport.AppCodec().MustMarshalJSON(bankGenesis)

	// setup feemarketGenesis params
	feemarketGenesis := feemarkettypes.DefaultGenesisState()
	if feemarketGenesis == nil {
		panic("nil feemarketGenesis")
	}
	feemarketGenesis.Params.EnableHeight = 1
	feemarketGenesis.Params.NoBaseFee = false

	if err := feemarketGenesis.Validate(); err != nil {
		panic(err)
	}
	genesisState[feemarkettypes.ModuleName] = teleport.AppCodec().MustMarshalJSON(feemarketGenesis)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	require.NoError(t, err)

	// init chain will set the validator set and initialize the genesis accounts
	teleport.InitChain(
		abci.RequestInitChain{
			ChainId:         "teleport_9000-1",
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: app.DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)

	// commit genesis changes
	teleport.Commit()
	teleport.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{
		ChainID:            "teleport_9000-1",
		Height:             teleport.LastBlockHeight() + 1,
		AppHash:            teleport.LastCommitID().Hash,
		ValidatorsHash:     valSet.Hash(),
		NextValidatorsHash: valSet.Hash(),
		ProposerAddress:    valSet.Proposer.Address,
	}})

	return teleport
}
