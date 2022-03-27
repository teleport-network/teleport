package keeper_test

import (
	"fmt"

	xibctmtypes "github.com/teleport-network/teleport/x/xibc/clients/light-clients/tendermint/types"
	"github.com/teleport-network/teleport/x/xibc/core/client/types"
	commitmenttypes "github.com/teleport-network/teleport/x/xibc/core/commitment/types"
	xibctesting "github.com/teleport-network/teleport/x/xibc/testing"
)

func (suite KeeperTestSuite) TestHandleCreateClientProposal() {
	header := suite.chainA.CreateTMClientHeader(
		suite.chainA.ChainID, suite.chainA.CurrentHeader.Height,
		types.NewHeight(0, uint64(suite.chainA.CurrentHeader.Height-1)),
		suite.chainA.CurrentHeader.Time, suite.chainA.Vals,
		suite.chainA.Vals, suite.chainA.Signers,
	)
	testCases := []struct {
		msg      string
		malleate func()
	}{
		{
			"success, create new client",
			func() {
				clientState := xibctmtypes.NewClientState("test", xibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.ZeroHeight(), commitmenttypes.GetSDKSpecs(), xibctesting.Prefix, 0)
				consensusState := xibctmtypes.NewConsensusState(
					header.GetTime(), header.Header.AppHash, header.Header.NextValidatorsHash,
				)
				var proposal *types.CreateClientProposal
				var err error
				proposal, err = types.NewCreateClientProposal("test", "test", "test", clientState, consensusState)
				suite.Require().NoError(err)
				_, err = suite.chainA.App.XIBCKeeper.ClientKeeper.HandleCreateClient(suite.chainA.GetContext(), proposal)
				suite.Require().NoError(err)
			},
		},
		{
			"fail, A client for this chainname already exists",
			func() {
				clientState := xibctmtypes.NewClientState("test", xibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.ZeroHeight(), commitmenttypes.GetSDKSpecs(), xibctesting.Prefix, 0)
				consensusState := xibctmtypes.NewConsensusState(
					header.GetTime(), header.Header.AppHash, header.Header.NextValidatorsHash,
				)
				var proposal *types.CreateClientProposal
				var err error
				proposal, err = types.NewCreateClientProposal("test", "test", "test", clientState, consensusState)
				suite.Require().NoError(err)
				_, err = suite.chainA.App.XIBCKeeper.ClientKeeper.HandleCreateClient(suite.chainA.GetContext(), proposal)
				suite.Require().NoError(err)
				_, err = suite.chainA.App.XIBCKeeper.ClientKeeper.HandleCreateClient(suite.chainA.GetContext(), proposal)
				suite.Require().Error(err)
			},
		},
		{
			"success, get client and compare",
			func() {
				clientState := xibctmtypes.NewClientState(
					"test", xibctmtypes.DefaultTrustLevel,
					trustingPeriod, ubdPeriod, maxClockDrift, types.NewHeight(0, 5),
					commitmenttypes.GetSDKSpecs(), xibctesting.Prefix, 0,
				)
				consensusState := xibctmtypes.NewConsensusState(
					header.GetTime(), header.Header.AppHash, header.Header.NextValidatorsHash,
				)
				var proposal *types.CreateClientProposal
				var err error
				proposal, err = types.NewCreateClientProposal("test", "test", "test", clientState, consensusState)
				suite.Require().NoError(err)
				_, err = suite.chainA.App.XIBCKeeper.ClientKeeper.HandleCreateClient(suite.chainA.GetContext(), proposal)
				suite.Require().NoError(err)
				client, _ := suite.chainA.App.XIBCKeeper.ClientKeeper.GetClientState(suite.chainA.GetContext(), "test")
				suite.Require().Equal(clientState, client, "clientState not equal")
				consensus, _ := suite.chainA.App.XIBCKeeper.ClientKeeper.GetClientConsensusState(suite.chainA.GetContext(), "test", types.NewHeight(0, 5))
				suite.Require().Equal(consensusState, consensus, "consensusState not equal")
			},
		},
	}
	for i, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s, %d/%d tests", tc.msg, i+1, len(testCases)), func() {
			suite.SetupTest() // reset the context
			tc.malleate()
		})
	}
}

func (suite KeeperTestSuite) TestHandleUpgradeClientProposal() {
	header := suite.chainA.CreateTMClientHeader(
		suite.chainA.ChainID, suite.chainA.CurrentHeader.Height,
		types.NewHeight(0, uint64(suite.chainA.CurrentHeader.Height-1)),
		suite.chainA.CurrentHeader.Time, suite.chainA.Vals,
		suite.chainA.Vals, suite.chainA.Signers,
	)
	testCases := []struct {
		msg      string
		malleate func()
	}{
		{
			"fail, client and consensus are not existing",
			func() {
				clientState := xibctmtypes.NewClientState("test", xibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.ZeroHeight(), commitmenttypes.GetSDKSpecs(), xibctesting.Prefix, 0)
				consensusState := xibctmtypes.NewConsensusState(
					header.GetTime(), header.Header.AppHash, header.Header.NextValidatorsHash,
				)
				var proposal *types.UpgradeClientProposal
				var err error
				proposal, err = types.NewUpgradeClientProposal("test", "test", "test", clientState, consensusState)
				suite.Require().NoError(err)
				_, err = suite.chainA.App.XIBCKeeper.ClientKeeper.HandleUpgradeClient(suite.chainA.GetContext(), proposal)
				suite.Require().Error(err)
			},
		},
		{
			"success, client and consensus are existing",
			func() {
				clientState := xibctmtypes.NewClientState(
					"test", xibctmtypes.DefaultTrustLevel,
					trustingPeriod, ubdPeriod, maxClockDrift, types.NewHeight(0, 5),
					commitmenttypes.GetSDKSpecs(), xibctesting.Prefix, 0,
				)
				consensusState := xibctmtypes.NewConsensusState(
					header.GetTime(), header.Header.AppHash, header.Header.NextValidatorsHash,
				)
				var proposal *types.CreateClientProposal
				var err error
				proposal, err = types.NewCreateClientProposal("test", "test", "test", clientState, consensusState)
				suite.Require().NoError(err)
				_, err = suite.chainA.App.XIBCKeeper.ClientKeeper.HandleCreateClient(suite.chainA.GetContext(), proposal)
				suite.Require().NoError(err)
				clientState2 := xibctmtypes.NewClientState(
					"test", xibctmtypes.DefaultTrustLevel,
					trustingPeriod*2, ubdPeriod, maxClockDrift, types.NewHeight(0, 6),
					commitmenttypes.GetSDKSpecs(), xibctesting.Prefix, 0,
				)
				consensusState2 := xibctmtypes.NewConsensusState(
					header.GetTime().Add(1), header.Header.AppHash, header.Header.NextValidatorsHash,
				)
				var proposal2 *types.UpgradeClientProposal
				proposal2, err = types.NewUpgradeClientProposal("test", "test", "test", clientState2, consensusState2)
				suite.Require().NoError(err)
				_, err = suite.chainA.App.XIBCKeeper.ClientKeeper.HandleUpgradeClient(suite.chainA.GetContext(), proposal2)
				suite.Require().NoError(err)

				// Check the consistency of clientState and consensusState
				client, _ := suite.chainA.App.XIBCKeeper.ClientKeeper.GetClientState(suite.chainA.GetContext(), "test")
				suite.Require().Equal(clientState2, client, "clientState not equal")
				consensus, _ := suite.chainA.App.XIBCKeeper.ClientKeeper.GetClientConsensusState(suite.chainA.GetContext(), "test", types.NewHeight(0, 6))
				suite.Require().Equal(consensusState2, consensus, "consensusState not equal")
			},
		},
	}
	for i, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s, %d/%d tests", tc.msg, i+1, len(testCases)), func() {
			suite.SetupTest() // reset the context
			tc.malleate()
		})
	}
}

func (suite KeeperTestSuite) TestHandleRegisterRelayerProposal() {
	header := suite.chainA.CreateTMClientHeader(
		suite.chainA.ChainID, suite.chainA.CurrentHeader.Height,
		types.NewHeight(0, uint64(suite.chainA.CurrentHeader.Height-1)),
		suite.chainA.CurrentHeader.Time, suite.chainA.Vals,
		suite.chainA.Vals, suite.chainA.Signers,
	)
	testCases := []struct {
		msg      string
		malleate func()
	}{
		{
			"success, exist client",
			func() {
				// TODO
				// clientState := xibctmtypes.NewClientState(
				// 	"test", xibctmtypes.DefaultTrustLevel,
				// 	trustingPeriod, ubdPeriod, maxClockDrift, types.NewHeight(0, 5),
				// 	commitmenttypes.GetSDKSpecs(), xibctesting.Prefix, 0,
				// )
				// consensusState := xibctmtypes.NewConsensusState(
				// 	header.GetTime(), header.Header.AppHash, header.Header.NextValidatorsHash,
				// )
				// var proposalCreate *types.CreateClientProposal
				// var err error
				// proposalCreate, err = types.NewCreateClientProposal("test", "test", "test", clientState, consensusState)
				// suite.Require().NoError(err)
				// _, err = suite.chainA.App.XIBCKeeper.ClientKeeper.HandleCreateClient(suite.chainA.GetContext(), proposalCreate)
				// suite.Require().NoError(err)

				// // set relayers
				// address := "xxx"
				// chians := []string{"xxx", "yyy"}
				// addresses := []string{"xxx", "yyy"}
				// relayerProposal := types.NewRegisterRelayerProposal("test", "test", address, chians, addresses)
				// err = suite.chainA.App.XIBCKeeper.ClientKeeper.HandleRegisterRelayer(suite.chainA.GetContext(), relayerProposal)
				// suite.Require().NoError(err)

				// // get relayers and compare
				// relayers2 := suite.chainA.App.XIBCKeeper.ClientKeeper.GetRelayers(suite.chainA.GetContext(), "test")
				// suite.Require().Equal(relayers, relayers2)
			},
		},
		{
			"success, no client",
			func() {
				address := "xxx"
				chians := []string{"xxx", "yyy"}
				addresses := []string{"xxx", "yyy"}
				relayerProposal := types.NewRegisterRelayerProposal("test", "test", address, chians, addresses)
				err := suite.chainA.App.XIBCKeeper.ClientKeeper.HandleRegisterRelayer(suite.chainA.GetContext(), relayerProposal)
				suite.Require().NoError(err)
			},
		},
		{
			"fail, no-existing client",
			func() {
				// set client "test"
				clientState := xibctmtypes.NewClientState(
					"test", xibctmtypes.DefaultTrustLevel,
					trustingPeriod, ubdPeriod, maxClockDrift, types.NewHeight(0, 5),
					commitmenttypes.GetSDKSpecs(), xibctesting.Prefix, 0,
				)
				consensusState := xibctmtypes.NewConsensusState(
					header.GetTime(), header.Header.AppHash, header.Header.NextValidatorsHash,
				)
				var proposalCreate *types.CreateClientProposal
				var err error
				proposalCreate, err = types.NewCreateClientProposal("test", "test", "test", clientState, consensusState)
				suite.Require().NoError(err)
				_, err = suite.chainA.App.XIBCKeeper.ClientKeeper.HandleCreateClient(suite.chainA.GetContext(), proposalCreate)
				suite.Require().NoError(err)
			},
		},
	}
	for i, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s, %d/%d tests", tc.msg, i+1, len(testCases)), func() {
			suite.SetupTest() // reset the context
			tc.malleate()
		})
	}
}
