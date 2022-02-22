package keeper_test

import (
	"fmt"
	"time"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	xibctmtypes "github.com/teleport-network/teleport/x/xibc/clients/light-clients/tendermint/types"
	"github.com/teleport-network/teleport/x/xibc/core/client/types"
	commitmenttypes "github.com/teleport-network/teleport/x/xibc/core/commitment/types"
	xibctesting "github.com/teleport-network/teleport/x/xibc/testing"
)

func (suite *KeeperTestSuite) TestQueryClientState() {
	var (
		req            *types.QueryClientStateRequest
		expClientState *codectypes.Any
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{{
		"invalid chainName",
		func() {
			req = &types.QueryClientStateRequest{}
		},
		false,
	}, {
		"client not found",
		func() {
			req = &types.QueryClientStateRequest{
				ChainName: testChainName,
			}
		},
		false,
	}, {
		"success",
		func() {
			clientState := xibctmtypes.NewClientState(testChainID, xibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.ZeroHeight(), commitmenttypes.GetSDKSpecs(), xibctesting.Prefix, 0)
			suite.keeper.SetClientState(suite.ctx, testChainName, clientState)

			var err error
			expClientState, err = types.PackClientState(clientState)
			suite.Require().NoError(err)

			req = &types.QueryClientStateRequest{
				ChainName: testChainName,
			}
		},
		true,
	}}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset

			tc.malleate()
			ctx := sdk.WrapSDKContext(suite.ctx)
			res, err := suite.queryClient.ClientState(ctx, req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(expClientState, res.ClientState)

				// ensure UnpackInterfaces is defined
				cachedValue := res.ClientState.GetCachedValue()
				suite.Require().NotNil(cachedValue)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQueryClientStates() {
	var (
		req             *types.QueryClientStatesRequest
		expClientStates = types.IdentifiedClientStates{}
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{{
		"empty pagination",
		func() {
			req = &types.QueryClientStatesRequest{}
		},
		true,
	}, {
		"success, no results",
		func() {
			req = &types.QueryClientStatesRequest{}
		},
		true,
	}, {
		"success",
		func() {
			// setup testing conditions
			path1 := xibctesting.NewPath(suite.chainA, suite.chainB)
			path2 := xibctesting.NewPath(suite.chainA, suite.chainC)
			suite.coordinator.SetupClients(path1)
			suite.coordinator.SetupClients(path2)

			clientStateA1 := path1.EndpointA.GetClientState()
			clientStateA2 := path2.EndpointA.GetClientState()

			idcs1 := types.NewIdentifiedClientState(path1.EndpointB.ChainName, clientStateA1)
			idcs2 := types.NewIdentifiedClientState(path2.EndpointB.ChainName, clientStateA2)

			// order is sorted by client id, localhost is last
			expClientStates = types.IdentifiedClientStates{idcs1, idcs2}.Sort()
			req = &types.QueryClientStatesRequest{}
		},
		true,
	}}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset
			tc.malleate()

			ctx := sdk.WrapSDKContext(suite.chainA.GetContext())
			res, err := suite.chainA.QueryServer.ClientStates(ctx, req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(expClientStates.Sort(), res.ClientStates)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQueryConsensusState() {
	var (
		req               *types.QueryConsensusStateRequest
		expConsensusState *codectypes.Any
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{{
		"invalid chainName",
		func() {
			req = &types.QueryConsensusStateRequest{}
		},
		false,
	}, {
		"invalid height",
		func() {
			req = &types.QueryConsensusStateRequest{
				ChainName:      testChainName,
				RevisionNumber: 0,
				RevisionHeight: 0,
				LatestHeight:   false,
			}
		},
		false,
	}, {
		"consensus state not found",
		func() {
			req = &types.QueryConsensusStateRequest{
				ChainName:    testChainName,
				LatestHeight: true,
			}
		},
		false,
	}, {
		"success latest height",
		func() {
			clientState := xibctmtypes.NewClientState(testChainID, xibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeight, commitmenttypes.GetSDKSpecs(), xibctesting.Prefix, 0)
			cs := xibctmtypes.NewConsensusState(
				suite.consensusState.Timestamp, []byte("hash1"), nil,
			)
			suite.keeper.SetClientState(suite.ctx, testChainName, clientState)
			suite.keeper.SetClientConsensusState(suite.ctx, testChainName, testClientHeight, cs)

			var err error
			expConsensusState, err = types.PackConsensusState(cs)
			suite.Require().NoError(err)

			req = &types.QueryConsensusStateRequest{
				ChainName:    testChainName,
				LatestHeight: true,
			}
		},
		true,
	}, {
		"success with height",
		func() {
			cs := xibctmtypes.NewConsensusState(
				suite.consensusState.Timestamp, []byte("hash1"), nil,
			)
			suite.keeper.SetClientConsensusState(suite.ctx, testChainName, testClientHeight, cs)

			var err error
			expConsensusState, err = types.PackConsensusState(cs)
			suite.Require().NoError(err)

			req = &types.QueryConsensusStateRequest{
				ChainName:      testChainName,
				RevisionNumber: 0,
				RevisionHeight: height,
			}
		},
		true,
	}}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset

			tc.malleate()
			ctx := sdk.WrapSDKContext(suite.ctx)
			res, err := suite.queryClient.ConsensusState(ctx, req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(expConsensusState, res.ConsensusState)

				// ensure UnpackInterfaces is defined
				cachedValue := res.ConsensusState.GetCachedValue()
				suite.Require().NotNil(cachedValue)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQueryConsensusStates() {
	var (
		req                *types.QueryConsensusStatesRequest
		expConsensusStates = []types.ConsensusStateWithHeight{}
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{{
		"invalid client identifier",
		func() {
			req = &types.QueryConsensusStatesRequest{}
		},
		false,
	}, {
		"empty pagination",
		func() {
			req = &types.QueryConsensusStatesRequest{
				ChainName: testChainName,
			}
		},
		true,
	}, {
		"success, no results",
		func() {
			req = &types.QueryConsensusStatesRequest{
				ChainName: testChainName,
				Pagination: &query.PageRequest{
					Limit:      3,
					CountTotal: true,
				},
			}
		},
		true,
	}, {
		"success",
		func() {
			cs := xibctmtypes.NewConsensusState(
				suite.consensusState.Timestamp, []byte("hash1"), nil,
			)
			cs2 := xibctmtypes.NewConsensusState(
				suite.consensusState.Timestamp.Add(time.Second), []byte("hash2"), nil,
			)

			clientState := xibctmtypes.NewClientState(
				testChainID, xibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeight, commitmenttypes.GetSDKSpecs(), xibctesting.Prefix, 0,
			)

			// Use CreateClient to ensure that processedTime metadata gets stored.
			err := suite.keeper.CreateClient(suite.ctx, testChainName, clientState, cs)
			suite.Require().NoError(err)
			suite.keeper.SetClientConsensusState(suite.ctx, testChainName, testClientHeight.Increment(), cs2)

			// order is swapped because the res is sorted by client id
			expConsensusStates = []types.ConsensusStateWithHeight{
				types.NewConsensusStateWithHeight(testClientHeight, cs),
				types.NewConsensusStateWithHeight(testClientHeight.Increment().(types.Height), cs2),
			}
			req = &types.QueryConsensusStatesRequest{
				ChainName: testChainName,
				Pagination: &query.PageRequest{
					Limit:      3,
					CountTotal: true,
				},
			}
		},
		true,
	}}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset

			tc.malleate()
			ctx := sdk.WrapSDKContext(suite.ctx)

			res, err := suite.queryClient.ConsensusStates(ctx, req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(len(expConsensusStates), len(res.ConsensusStates))
				for i := range expConsensusStates {
					suite.Require().NotNil(res.ConsensusStates[i])
					suite.Require().Equal(expConsensusStates[i], res.ConsensusStates[i])

					// ensure UnpackInterfaces is defined
					cachedValue := res.ConsensusStates[i].ConsensusState.GetCachedValue()
					suite.Require().NotNil(cachedValue)
				}
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
