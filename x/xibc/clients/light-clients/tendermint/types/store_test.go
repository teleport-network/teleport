package types_test

import (
	"github.com/teleport-network/teleport/x/xibc/clients/light-clients/tendermint/types"
	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	"github.com/teleport-network/teleport/x/xibc/core/host"
	"github.com/teleport-network/teleport/x/xibc/exported"
	xibctesting "github.com/teleport-network/teleport/x/xibc/testing"
)

func (suite *TendermintTestSuite) TestGetConsensusState() {
	var (
		height exported.Height
		path   *xibctesting.Path
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{{
		"success", func() {}, true,
	}, {
		"consensus state not found", func() {
			// use height with no consensus state set
			height = height.(clienttypes.Height).Increment()
		},
		false,
	}, {
		"not a consensus state interface", func() {
			// marshal an empty client state and set as consensus state
			store := suite.chainA.App.XIBCKeeper.ClientKeeper.ClientStore(suite.chainA.GetContext(), path.EndpointB.ChainName)
			clientStateBz := suite.chainA.App.XIBCKeeper.ClientKeeper.MustMarshalClientState(&types.ClientState{})
			store.Set(host.ConsensusStateKey(height), clientStateBz)
		},
		false,
	}}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			path = xibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)

			clientState := path.EndpointA.GetClientState()
			height = clientState.GetLatestHeight()

			tc.malleate() // change vars as necessary

			store := suite.chainA.App.XIBCKeeper.ClientKeeper.ClientStore(suite.chainA.GetContext(), path.EndpointB.ChainName)
			consensusState, err := types.GetConsensusState(store, suite.chainA.Codec, height)

			if tc.expPass {
				suite.Require().NoError(err)
				expConsensusState, found := suite.chainA.GetConsensusState(path.EndpointB.ChainName, height)
				suite.Require().True(found)
				suite.Require().Equal(expConsensusState, consensusState)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(consensusState)
			}
		})
	}
}

func (suite *TendermintTestSuite) TestGetProcessedTime() {
	// setup
	path := xibctesting.NewPath(suite.chainA, suite.chainB)

	suite.coordinator.UpdateTime()
	// coordinator increments time before creating client
	expectedTime := suite.chainA.CurrentHeader.Time.Add(xibctesting.TimeIncrement)

	// Verify ProcessedTime on CreateClient
	err := path.EndpointA.CreateClient()
	suite.Require().NoError(err)

	clientState := path.EndpointA.GetClientState()
	height := clientState.GetLatestHeight()

	store := path.EndpointA.ClientStore()
	actualTime, ok := types.GetProcessedTime(store, height)
	suite.Require().True(ok, "could not retrieve processed time for stored consensus state")
	suite.Require().Equal(uint64(expectedTime.UnixNano()), actualTime, "retrieved processed time is not expected value")

	suite.coordinator.UpdateTime()
	// coordinator increments time before updating client
	expectedTime = suite.chainA.CurrentHeader.Time.Add(xibctesting.TimeIncrement)

	// Verify ProcessedTime on UpdateClient
	err = path.EndpointA.UpdateClient()
	suite.Require().NoError(err)

	clientState = path.EndpointA.GetClientState()
	height = clientState.GetLatestHeight()

	store = path.EndpointA.ClientStore()
	actualTime, ok = types.GetProcessedTime(store, height)
	suite.Require().True(ok, "could not retrieve processed time for stored consensus state")
	suite.Require().Equal(uint64(expectedTime.UnixNano()), actualTime, "retrieved processed time is not expected value")

	// try to get processed time for height that doesn't exist in store
	_, ok = types.GetProcessedTime(store, clienttypes.NewHeight(1, 1))
	suite.Require().False(ok, "retrieved processed time for a non-existent consensus state")
}
