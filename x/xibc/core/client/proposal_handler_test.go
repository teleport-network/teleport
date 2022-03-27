package client_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	client "github.com/teleport-network/teleport/x/xibc/core/client"
	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	xibctesting "github.com/teleport-network/teleport/x/xibc/testing"
)

type ClientTestSuite struct {
	suite.Suite
	coordinator *xibctesting.Coordinator
	chainA      *xibctesting.TestChain
	chainB      *xibctesting.TestChain
}

func (suite *ClientTestSuite) SetupTest() {
	suite.coordinator = xibctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(xibctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(xibctesting.GetChainID(1))
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

func (suite *ClientTestSuite) TestNewClientUpdateProposalHandler() {
	var (
		content govtypes.Content
		err     error
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{{
		"valid create client proposal",
		func() {
			// setup testing conditions
			path := xibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)
			clientState := path.EndpointA.GetClientState()
			consensusState := path.EndpointA.GetConsensusState(clientState.GetLatestHeight())

			content, err = clienttypes.NewCreateClientProposal(
				xibctesting.Title,
				xibctesting.Description,
				"test-chain-name",
				clientState,
				consensusState,
			)
			suite.Require().NoError(err)
		},
		true,
	}, {
		"valid create client proposal",
		func() {
			// setup testing conditions
			path := xibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)
			clientState := path.EndpointA.GetClientState()
			consensusState := path.EndpointA.GetConsensusState(clientState.GetLatestHeight())

			content, err = clienttypes.NewUpgradeClientProposal(
				xibctesting.Title,
				xibctesting.Description,
				path.EndpointB.ChainName,
				clientState,
				consensusState,
			)
			suite.Require().NoError(err)
		},
		true,
	}, {
		// TODO
		// "valid create client proposal",
		// func() {

		// 	// setup testing conditions
		// 	path := xibctesting.NewPath(suite.chainA, suite.chainB)
		// 	suite.coordinator.SetupClients(path)

		// 	relayers := []string{
		// 		suite.chainB.SenderAcc.String(),
		// 	}

		// 	content = clienttypes.NewRegisterRelayerProposal(
		// 		xibctesting.Title,
		// 		xibctesting.Description,
		// 		path.EndpointB.ChainName,
		// 		relayers,
		// 	)
		// 	suite.Require().NoError(err)
		// },
		// true,
	}, {
		"nil proposal",
		func() {
			content = nil
		},
		false,
	}, {
		"unsupported proposal type",
		func() {
			content = distributiontypes.NewCommunityPoolSpendProposal(
				xibctesting.Title,
				xibctesting.Description,
				suite.chainA.SenderAcc,
				sdk.NewCoins(sdk.NewCoin("communityfunds", sdk.NewInt(10))),
			)
		},
		false,
	}}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset
			tc.malleate()

			proposalHandler := client.NewClientProposalHandler(suite.chainA.App.XIBCKeeper.ClientKeeper)

			err = proposalHandler(suite.chainA.GetContext(), content)
			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
