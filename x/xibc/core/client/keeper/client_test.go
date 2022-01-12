package keeper_test

import (
	"fmt"
	"time"

	tmtypes "github.com/tendermint/tendermint/types"

	xibctmtypes "github.com/teleport-network/teleport/x/xibc/clients/light-clients/tendermint/types"
	"github.com/teleport-network/teleport/x/xibc/core/client/types"
	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	commitmenttypes "github.com/teleport-network/teleport/x/xibc/core/commitment/types"
	xibctesting "github.com/teleport-network/teleport/x/xibc/testing"
)

func (suite *KeeperTestSuite) TestUpdateClientTendermint() {
	// Must create header creation functions since suite.header gets recreated on each test case
	createFutureUpdateFn := func(s *KeeperTestSuite) *xibctmtypes.Header {
		heightPlus3 := clienttypes.NewHeight(suite.header.GetHeight().GetRevisionNumber(), suite.header.GetHeight().GetRevisionHeight()+3)
		height := suite.header.GetHeight().(clienttypes.Height)

		return suite.chainA.CreateTMClientHeader(
			testChainID, int64(heightPlus3.RevisionHeight), height, suite.header.Header.Time.Add(time.Hour),
			suite.valSet, suite.valSet, []tmtypes.PrivValidator{suite.privVal},
		)
	}
	createPastUpdateFn := func(s *KeeperTestSuite) *xibctmtypes.Header {
		heightMinus2 := clienttypes.NewHeight(suite.header.GetHeight().GetRevisionNumber(), suite.header.GetHeight().GetRevisionHeight()-2)
		heightMinus4 := clienttypes.NewHeight(suite.header.GetHeight().GetRevisionNumber(), suite.header.GetHeight().GetRevisionHeight()-4)

		return suite.chainA.CreateTMClientHeader(
			testChainID, int64(heightMinus2.RevisionHeight), heightMinus4, suite.header.Header.Time,
			suite.valSet, suite.valSet, []tmtypes.PrivValidator{suite.privVal},
		)
	}

	var (
		updateHeader *xibctmtypes.Header
		clientState  *xibctmtypes.ClientState
		chainName    string
		err          error
	)

	cases := []struct {
		name     string
		malleate func() error
		expPass  bool
	}{{
		"valid update",
		func() error {
			clientState = xibctmtypes.NewClientState(testChainID, xibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeight, commitmenttypes.GetSDKSpecs(), xibctesting.Prefix, 0)
			err = suite.keeper.CreateClient(suite.ctx, chainName, clientState, suite.consensusState)

			// store intermediate consensus state to check that trustedHeight does not need to be highest consensus state before header height
			incrementedClientHeight := testClientHeight.Increment().(types.Height)
			intermediateConsState := &xibctmtypes.ConsensusState{
				Timestamp:          suite.now.Add(time.Minute),
				NextValidatorsHash: suite.valSetHash,
			}
			suite.keeper.SetClientConsensusState(suite.ctx, chainName, incrementedClientHeight, intermediateConsState)

			clientState.LatestHeight = incrementedClientHeight
			suite.keeper.SetClientState(suite.ctx, chainName, clientState)

			updateHeader = createFutureUpdateFn(suite)
			return err
		},
		true,
	}, {
		"valid past update",
		func() error {
			clientState = xibctmtypes.NewClientState(testChainID, xibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeight, commitmenttypes.GetSDKSpecs(), xibctesting.Prefix, 0)
			err = suite.keeper.CreateClient(suite.ctx, chainName, clientState, suite.consensusState)
			suite.Require().NoError(err)

			height1 := types.NewHeight(0, 1)

			// store previous consensus state
			prevConsState := &xibctmtypes.ConsensusState{
				Timestamp:          suite.past,
				NextValidatorsHash: suite.valSetHash,
			}
			suite.keeper.SetClientConsensusState(suite.ctx, chainName, height1, prevConsState)

			height2 := types.NewHeight(0, 2)

			// store intermediate consensus state to check that trustedHeight does not need to be hightest consensus state before header height
			intermediateConsState := &xibctmtypes.ConsensusState{
				Timestamp:          suite.past.Add(time.Minute),
				NextValidatorsHash: suite.valSetHash,
			}
			suite.keeper.SetClientConsensusState(suite.ctx, chainName, height2, intermediateConsState)

			// updateHeader will fill in consensus state between prevConsState and suite.consState
			// clientState should not be updated
			updateHeader = createPastUpdateFn(suite)
			return nil
		},
		true,
	}, {
		"client state not found",
		func() error {
			updateHeader = createFutureUpdateFn(suite)
			return nil
		},
		false,
	}, {
		"consensus state not found",
		func() error {
			clientState = xibctmtypes.NewClientState(testChainID, xibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeight, commitmenttypes.GetSDKSpecs(), xibctesting.Prefix, 0)
			suite.keeper.SetClientState(suite.ctx, testChainName, clientState)
			updateHeader = createFutureUpdateFn(suite)

			return nil
		},
		false,
	}, {
		"valid past update before client was frozen",
		func() error {
			clientState = xibctmtypes.NewClientState(testChainID, xibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeight, commitmenttypes.GetSDKSpecs(), xibctesting.Prefix, 0)
			err = suite.keeper.CreateClient(suite.ctx, chainName, clientState, suite.consensusState)
			suite.Require().NoError(err)

			height1 := types.NewHeight(0, 1)

			// store previous consensus state
			prevConsState := &xibctmtypes.ConsensusState{
				Timestamp:          suite.past,
				NextValidatorsHash: suite.valSetHash,
			}
			suite.keeper.SetClientConsensusState(suite.ctx, chainName, height1, prevConsState)

			// updateHeader will fill in consensus state between prevConsState and suite.consState
			// clientState should not be updated
			updateHeader = createPastUpdateFn(suite)
			return nil
		}, true,
	}, {
		"invalid header", func() error {
			clientState = xibctmtypes.NewClientState(
				testChainID, xibctmtypes.DefaultTrustLevel, trustingPeriod,
				ubdPeriod, maxClockDrift, testClientHeight,
				commitmenttypes.GetSDKSpecs(), xibctesting.Prefix, 0,
			)
			err := suite.keeper.CreateClient(suite.ctx, chainName, clientState, suite.consensusState)
			suite.Require().NoError(err)
			updateHeader = createPastUpdateFn(suite)

			return nil
		},
		false,
	}}

	for i, tc := range cases {
		suite.Run(
			fmt.Sprintf("Case %s", tc.name),
			func() {
				suite.SetupTest()
				chainName = testChainName // must be explicitly changed
				suite.Require().NoError(tc.malleate())

				suite.ctx = suite.ctx.WithBlockTime(updateHeader.Header.Time.Add(time.Minute))
				err = suite.keeper.UpdateClient(suite.ctx, chainName, updateHeader)

				if tc.expPass {
					suite.Require().NoError(err, tc.name)

					expConsensusState := &xibctmtypes.ConsensusState{
						Timestamp:          updateHeader.GetTime(),
						Root:               updateHeader.Header.GetAppHash(),
						NextValidatorsHash: updateHeader.Header.NextValidatorsHash,
					}

					newClientState, found := suite.keeper.GetClientState(suite.ctx, chainName)
					suite.Require().True(found, "valid test case %d failed: %s", i, tc.name)

					consensusState, found := suite.keeper.GetClientConsensusState(suite.ctx, chainName, updateHeader.GetHeight())
					suite.Require().True(found, "valid test case %d failed: %s", i, tc.name)

					// Determine if clientState should be updated or not
					if updateHeader.GetHeight().GT(clientState.GetLatestHeight()) {
						// Header Height is greater than clientState latest Height, clientState should be updated with header.GetHeight()
						suite.Require().Equal(updateHeader.GetHeight(), newClientState.GetLatestHeight(), "clientstate height did not update")
					} else {
						// Update will add past consensus state, clientState should not be updated at all
						suite.Require().Equal(clientState.GetLatestHeight(), newClientState.GetLatestHeight(), "client state height updated for past header")
					}

					suite.Require().NoError(err, "valid test case %d failed: %s", i, tc.name)
					suite.Require().Equal(expConsensusState, consensusState, "consensus state should have been updated on case %s", tc.name)
				} else {
					suite.Require().Error(err, "invalid test case %d passed: %s", i, tc.name)
				}
			},
		)
	}
}
