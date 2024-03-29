package keeper_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/teleport-network/teleport/app"
	xibctmtypes "github.com/teleport-network/teleport/x/xibc/clients/light-clients/tendermint/types"
	"github.com/teleport-network/teleport/x/xibc/core/client/keeper"
	"github.com/teleport-network/teleport/x/xibc/core/client/types"
	commitmenttypes "github.com/teleport-network/teleport/x/xibc/core/commitment/types"
	"github.com/teleport-network/teleport/x/xibc/exported"
	xibctesting "github.com/teleport-network/teleport/x/xibc/testing"
	xibctestingmock "github.com/teleport-network/teleport/x/xibc/testing/mock"
)

const (
	testChainID    = "gaiahub-0"
	testChainName  = "tendermint-0"
	testChainName2 = "tendermint-1"
	testChainName3 = "tendermint-2"

	height = 5

	trustingPeriod time.Duration = time.Hour * 24 * 7 * 2
	ubdPeriod      time.Duration = time.Hour * 24 * 7 * 3
	maxClockDrift  time.Duration = time.Second * 10
)

var (
	testClientHeight = types.NewHeight(0, 5)
)

type KeeperTestSuite struct {
	suite.Suite

	coordinator *xibctesting.Coordinator

	chainA *xibctesting.TestChain
	chainB *xibctesting.TestChain
	chainC *xibctesting.TestChain

	cdc            codec.Codec
	ctx            sdk.Context
	keeper         *keeper.Keeper
	consensusState *xibctmtypes.ConsensusState
	header         *xibctmtypes.Header
	valSet         *tmtypes.ValidatorSet
	valSetHash     tmbytes.HexBytes
	privVal        tmtypes.PrivValidator
	now            time.Time
	past           time.Time

	queryClient types.QueryClient
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.coordinator = xibctesting.NewCoordinator(suite.T(), 3)

	suite.chainA = suite.coordinator.GetChain(xibctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(xibctesting.GetChainID(1))
	suite.chainC = suite.coordinator.GetChain(xibctesting.GetChainID(2))

	isCheckTx := false
	suite.now = time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
	suite.past = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	now := suite.now.Add(time.Hour)
	app := app.Setup(isCheckTx, nil)

	suite.cdc = app.AppCodec()
	suite.ctx = app.BaseApp.NewContext(
		isCheckTx,
		tmproto.Header{
			Height:  height,
			ChainID: testChainName,
			Time:    now,
		},
	)
	suite.keeper = &app.XIBCKeeper.ClientKeeper
	suite.privVal = xibctestingmock.NewPV()

	pubKey, err := suite.privVal.GetPubKey()
	suite.Require().NoError(err)

	testClientHeightMinus1 := types.NewHeight(0, height-1)

	validator := tmtypes.NewValidator(pubKey, 1)
	suite.valSet = tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})
	suite.valSetHash = suite.valSet.Hash()
	suite.header = suite.chainA.CreateTMClientHeader(
		testChainID, int64(testClientHeight.RevisionHeight),
		testClientHeightMinus1, now, suite.valSet, suite.valSet,
		[]tmtypes.PrivValidator{suite.privVal},
	)
	suite.consensusState = xibctmtypes.NewConsensusState(suite.now, []byte("hash"), suite.valSetHash)

	var validators stakingtypes.Validators
	for i := 1; i < 11; i++ {
		privVal := xibctestingmock.NewPV()
		tmPk, err := privVal.GetPubKey()
		suite.Require().NoError(err)
		pk, err := cryptocodec.FromTmPubKeyInterface(tmPk)
		suite.Require().NoError(err)
		val, err := stakingtypes.NewValidator(sdk.ValAddress(pk.Address()), pk, stakingtypes.Description{})
		suite.Require().NoError(err)

		val.Status = stakingtypes.Bonded
		val.Tokens = sdk.NewInt(rand.Int63())
		validators = append(validators, val)

		hi := stakingtypes.NewHistoricalInfo(suite.ctx.BlockHeader(), validators, sdk.DefaultPowerReduction)
		app.StakingKeeper.SetHistoricalInfo(suite.ctx, int64(i), &hi)
	}

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.XIBCKeeper.ClientKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestSetClientState() {
	clientState := xibctmtypes.NewClientState(
		testChainID, xibctmtypes.DefaultTrustLevel,
		trustingPeriod, ubdPeriod, maxClockDrift,
		types.ZeroHeight(), commitmenttypes.GetSDKSpecs(),
		xibctesting.Prefix, 0,
	)
	suite.keeper.SetClientState(suite.ctx, testChainName, clientState)

	retrievedState, found := suite.keeper.GetClientState(suite.ctx, testChainName)
	suite.Require().True(found, "GetClientState failed")
	suite.Require().Equal(clientState, retrievedState, "Client states are not equal")
}

func (suite *KeeperTestSuite) TestSetClientConsensusState() {
	suite.keeper.SetClientConsensusState(suite.ctx, testChainName, testClientHeight, suite.consensusState)

	retrievedConsState, found := suite.keeper.GetClientConsensusState(suite.ctx, testChainName, testClientHeight)
	suite.Require().True(found, "GetConsensusState failed")

	tmConsState, ok := retrievedConsState.(*xibctmtypes.ConsensusState)
	suite.Require().True(ok)
	suite.Require().Equal(suite.consensusState, tmConsState, "ConsensusState not stored correctly")
}

func (suite KeeperTestSuite) TestGetAllGenesisClients() {
	chainNames := []string{
		testChainName2, testChainName3, testChainName,
	}
	expClients := []exported.ClientState{
		xibctmtypes.NewClientState(testChainID, xibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.ZeroHeight(), commitmenttypes.GetSDKSpecs(), xibctesting.Prefix, 0),
		xibctmtypes.NewClientState(testChainID, xibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.ZeroHeight(), commitmenttypes.GetSDKSpecs(), xibctesting.Prefix, 0),
		xibctmtypes.NewClientState(testChainID, xibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.ZeroHeight(), commitmenttypes.GetSDKSpecs(), xibctesting.Prefix, 0),
	}

	expGenClients := make(types.IdentifiedClientStates, len(expClients))

	for i := range expClients {
		suite.chainA.App.XIBCKeeper.ClientKeeper.SetClientState(suite.chainA.GetContext(), chainNames[i], expClients[i])
		expGenClients[i] = types.NewIdentifiedClientState(chainNames[i], expClients[i])
	}

	genClients := suite.chainA.App.XIBCKeeper.ClientKeeper.GetAllGenesisClients(suite.chainA.GetContext())
	suite.Require().Equal(expGenClients.Sort(), genClients)
}

func (suite KeeperTestSuite) TestGetAllGenesisMetadata() {
	expectedGenMetadata := []types.IdentifiedGenesisMetadata{
		types.NewIdentifiedGenesisMetadata(
			"clientA",
			[]types.GenesisMetadata{
				types.NewGenesisMetadata(xibctmtypes.ProcessedTimeKey(types.NewHeight(0, 1)), []byte("foo")),
				types.NewGenesisMetadata(xibctmtypes.ProcessedTimeKey(types.NewHeight(0, 2)), []byte("bar")),
				types.NewGenesisMetadata(xibctmtypes.ProcessedTimeKey(types.NewHeight(0, 3)), []byte("baz")),
			},
		),
		types.NewIdentifiedGenesisMetadata(
			"clientB",
			[]types.GenesisMetadata{
				types.NewGenesisMetadata(xibctmtypes.ProcessedTimeKey(types.NewHeight(1, 100)), []byte("val1")),
				types.NewGenesisMetadata(xibctmtypes.ProcessedTimeKey(types.NewHeight(2, 300)), []byte("val2")),
			},
		),
	}

	genClients := []types.IdentifiedClientState{
		types.NewIdentifiedClientState("clientA", &xibctmtypes.ClientState{}), types.NewIdentifiedClientState("clientB", &xibctmtypes.ClientState{}),
	}

	suite.chainA.App.XIBCKeeper.ClientKeeper.SetAllClientMetadata(suite.chainA.GetContext(), expectedGenMetadata)

	actualGenMetadata, err := suite.chainA.App.XIBCKeeper.ClientKeeper.GetAllClientMetadata(suite.chainA.GetContext(), genClients)
	suite.Require().NoError(err, "get client metadata returned error unexpectedly")
	suite.Require().Equal(expectedGenMetadata, actualGenMetadata, "retrieved metadata is unexpected")
}

func (suite KeeperTestSuite) TestConsensusStateHelpers() {
	// initial setup
	clientState := xibctmtypes.NewClientState(
		testChainID, xibctmtypes.DefaultTrustLevel,
		trustingPeriod, ubdPeriod, maxClockDrift,
		testClientHeight, commitmenttypes.GetSDKSpecs(),
		xibctesting.Prefix, 0,
	)

	suite.keeper.SetClientState(suite.ctx, testChainName, clientState)
	suite.keeper.SetClientConsensusState(suite.ctx, testChainName, testClientHeight, suite.consensusState)

	nextState := xibctmtypes.NewConsensusState(suite.now, []byte("next"), suite.valSetHash)

	testClientHeightPlus5 := types.NewHeight(0, height+5)

	header := suite.chainA.CreateTMClientHeader(
		testChainName, int64(testClientHeightPlus5.RevisionHeight),
		testClientHeight, suite.header.Header.Time.Add(time.Minute),
		suite.valSet, suite.valSet, []tmtypes.PrivValidator{suite.privVal},
	)

	// mock update functionality
	clientState.LatestHeight = header.GetHeight().(types.Height)
	suite.keeper.SetClientConsensusState(suite.ctx, testChainName, header.GetHeight(), nextState)
	suite.keeper.SetClientState(suite.ctx, testChainName, clientState)

	latest, ok := suite.keeper.GetLatestClientConsensusState(suite.ctx, testChainName)
	suite.Require().True(ok)
	suite.Require().Equal(nextState, latest, "Latest client not returned correctly")
}

// 2 clients in total are created on chainA. The first client is updated so it contains an initial consensus state
// and a consensus state at the update height.
func (suite KeeperTestSuite) TestGetAllConsensusStates() {
	// setup testing conditions
	path := xibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupClients(path)

	clientState := suite.chainA.GetClientState(path.EndpointB.ChainName)
	expConsensusHeight0 := clientState.GetLatestHeight()
	consensusState0, ok := suite.chainA.GetConsensusState(path.EndpointB.ChainName, expConsensusHeight0)
	suite.Require().True(ok)

	// update client to create a second consensus state
	err := path.EndpointA.UpdateClient()
	suite.Require().NoError(err)

	clientState = suite.chainA.GetClientState(path.EndpointB.ChainName)
	expConsensusHeight1 := clientState.GetLatestHeight()
	suite.Require().True(expConsensusHeight1.GT(expConsensusHeight0))
	consensusState1, ok := suite.chainA.GetConsensusState(path.EndpointB.ChainName, expConsensusHeight1)
	suite.Require().True(ok)

	consensusStateAny0, err := types.PackConsensusState(consensusState0)
	suite.Require().NoError(err)

	consensusStateAny1, err := types.PackConsensusState(consensusState1)
	suite.Require().NoError(err)
	expConsensus := []types.ConsensusStateWithHeight{
		{Height: expConsensusHeight0.(types.Height), ConsensusState: consensusStateAny0},
		{Height: expConsensusHeight1.(types.Height), ConsensusState: consensusStateAny1},
	}

	consStates := path.EndpointA.Chain.App.XIBCKeeper.ClientKeeper.GetAllConsensusStates(suite.chainA.GetContext())
	suite.Require().Len(consStates, 1)
	suite.Require().Equal(path.EndpointB.ChainName, consStates[0].ChainName)
	suite.Require().Equal(expConsensus, consStates[0].ConsensusStates)
}
