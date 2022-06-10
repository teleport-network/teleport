package xibc_test

import (
	"fmt"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/teleport-network/teleport/app"
	xibc "github.com/teleport-network/teleport/x/xibc"
	xibctmtypes "github.com/teleport-network/teleport/x/xibc/clients/light-clients/tendermint/types"
	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	commitmenttypes "github.com/teleport-network/teleport/x/xibc/core/commitment/types"
	"github.com/teleport-network/teleport/x/xibc/core/host"
	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
	xibctesting "github.com/teleport-network/teleport/x/xibc/testing"
	"github.com/teleport-network/teleport/x/xibc/types"
)

func (suite *XIBCTestSuite) TestValidateGenesis() {
	header := suite.chainA.CreateTMClientHeader(
		suite.chainA.ChainID, suite.chainA.CurrentHeader.Height,
		clienttypes.NewHeight(0, uint64(suite.chainA.CurrentHeader.Height-1)),
		suite.chainA.CurrentHeader.Time, suite.chainA.Vals, suite.chainA.Vals,
		suite.chainA.Signers,
	)

	testCases := []struct {
		name     string
		genState *types.GenesisState
		expPass  bool
	}{{
		name:     "default",
		genState: types.DefaultGenesisState(),
		expPass:  true,
	}, {
		name: "valid genesis",
		genState: &types.GenesisState{
			ClientGenesis: clienttypes.NewGenesisState(
				[]clienttypes.IdentifiedClientState{
					clienttypes.NewIdentifiedClientState(
						chainName, xibctmtypes.NewClientState(
							suite.chainA.ChainID, xibctmtypes.DefaultTrustLevel,
							xibctesting.TrustingPeriod, xibctesting.UnbondingPeriod,
							xibctesting.MaxClockDrift, clientHeight,
							commitmenttypes.GetSDKSpecs(), xibctesting.Prefix, 0,
						),
					),
				},
				[]clienttypes.ClientConsensusStates{
					clienttypes.NewClientConsensusStates(
						chainName,
						[]clienttypes.ConsensusStateWithHeight{
							clienttypes.NewConsensusStateWithHeight(
								header.GetHeight().(clienttypes.Height),
								xibctmtypes.NewConsensusState(
									header.GetTime(), header.Header.AppHash, header.Header.NextValidatorsHash,
								),
							),
						},
					),
				},
				[]clienttypes.IdentifiedGenesisMetadata{
					clienttypes.NewIdentifiedGenesisMetadata(
						chainName,
						[]clienttypes.GenesisMetadata{
							clienttypes.NewGenesisMetadata([]byte("key1"), []byte("val1")),
							clienttypes.NewGenesisMetadata([]byte("key2"), []byte("val2")),
						},
					),
				},
				chainName2,
			),
		},
		expPass: true,
	}, {
		name: "invalid client genesis",
		genState: &types.GenesisState{
			ClientGenesis: clienttypes.NewGenesisState(
				[]clienttypes.IdentifiedClientState{
					clienttypes.NewIdentifiedClientState(
						chainName, xibctmtypes.NewClientState(
							suite.chainA.ChainID, xibctmtypes.DefaultTrustLevel,
							xibctesting.TrustingPeriod, xibctesting.UnbondingPeriod,
							xibctesting.MaxClockDrift, clientHeight,
							commitmenttypes.GetSDKSpecs(), xibctesting.Prefix, 0,
						),
					),
				},
				nil,
				[]clienttypes.IdentifiedGenesisMetadata{
					clienttypes.NewIdentifiedGenesisMetadata(
						chainName,
						[]clienttypes.GenesisMetadata{
							clienttypes.NewGenesisMetadata([]byte(""), []byte("val1")),
							clienttypes.NewGenesisMetadata([]byte("key2"), []byte("")),
						},
					),
				},
				chainName2,
			),
		},
		expPass: false,
	},
	}

	for _, tc := range testCases {
		if tc.expPass {
			suite.Require().NoError(tc.genState.Validate(), tc.name)
		} else {
			suite.Require().Error(tc.genState.Validate(), tc.name)
		}
	}
}

func (suite *XIBCTestSuite) TestInitGenesis() {
	header := suite.chainA.CreateTMClientHeader(
		suite.chainA.ChainID, suite.chainA.CurrentHeader.Height,
		clienttypes.NewHeight(0, uint64(suite.chainA.CurrentHeader.Height-1)),
		suite.chainA.CurrentHeader.Time, suite.chainA.Vals,
		suite.chainA.Vals, suite.chainA.Signers,
	)

	testCases := []struct {
		name     string
		genState *types.GenesisState
	}{{
		name:     "default",
		genState: types.DefaultGenesisState(),
	}, {
		name: "valid genesis",
		genState: &types.GenesisState{
			ClientGenesis: clienttypes.NewGenesisState(
				[]clienttypes.IdentifiedClientState{
					clienttypes.NewIdentifiedClientState(
						chainName, xibctmtypes.NewClientState(
							suite.chainA.ChainID, xibctmtypes.DefaultTrustLevel,
							xibctesting.TrustingPeriod, xibctesting.UnbondingPeriod,
							xibctesting.MaxClockDrift, clientHeight,
							commitmenttypes.GetSDKSpecs(), xibctesting.Prefix, 0,
						),
					),
				},
				[]clienttypes.ClientConsensusStates{
					clienttypes.NewClientConsensusStates(
						chainName,
						[]clienttypes.ConsensusStateWithHeight{
							clienttypes.NewConsensusStateWithHeight(
								header.GetHeight().(clienttypes.Height),
								xibctmtypes.NewConsensusState(
									header.GetTime(), header.Header.AppHash, header.Header.NextValidatorsHash,
								),
							),
						},
					),
				},
				[]clienttypes.IdentifiedGenesisMetadata{
					clienttypes.NewIdentifiedGenesisMetadata(
						chainName,
						[]clienttypes.GenesisMetadata{
							clienttypes.NewGenesisMetadata([]byte("key1"), []byte("val1")),
							clienttypes.NewGenesisMetadata([]byte("key2"), []byte("val2")),
						},
					),
				},
				chainName2,
			),
		},
	}}

	for _, tc := range testCases {
		teleport := app.Setup(false, nil)

		suite.Require().NotPanics(func() {
			xibc.InitGenesis(teleport.BaseApp.NewContext(false, tmproto.Header{Height: 1}), *teleport.XIBCKeeper, true, tc.genState)
		})
	}
}

func (suite *XIBCTestSuite) TestExportGenesis() {
	testCases := []struct {
		msg      string
		malleate func()
	}{{
		"success",
		func() {
			// creates clients
			suite.coordinator.SetupClients(xibctesting.NewPath(suite.chainA, suite.chainB))
			// create extra clients
			suite.coordinator.SetupClients(xibctesting.NewPath(suite.chainA, suite.chainB))
			suite.coordinator.SetupClients(xibctesting.NewPath(suite.chainA, suite.chainB))
		},
	}}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest()

			tc.malleate()

			var gs *types.GenesisState
			suite.Require().NotPanics(func() {
				gs = xibc.ExportGenesis(suite.chainA.GetContext(), *suite.chainA.App.XIBCKeeper)
			})

			// init genesis based on export
			suite.Require().NotPanics(func() {
				xibc.InitGenesis(suite.chainA.GetContext(), *suite.chainA.App.XIBCKeeper, true, gs)
			})

			suite.Require().NotPanics(func() {
				cdc := codec.NewProtoCodec(suite.chainA.App.InterfaceRegistry())
				genState := cdc.MustMarshalJSON(gs)
				cdc.MustUnmarshalJSON(genState, gs)
			})

			// init genesis based on marshal and unmarshal
			suite.Require().NotPanics(func() {
				xibc.InitGenesis(suite.chainA.GetContext(), *suite.chainA.App.XIBCKeeper, true, gs)
			})
		})
	}
}

func (suite *XIBCTestSuite) TestResetStates() {
	header := suite.chainA.CreateTMClientHeader(
		suite.chainA.ChainID, suite.chainA.CurrentHeader.Height,
		clienttypes.NewHeight(0, uint64(suite.chainA.CurrentHeader.Height-1)),
		suite.chainA.CurrentHeader.Time, suite.chainA.Vals,
		suite.chainA.Vals, suite.chainA.Signers,
	)

	packetStates := []packettypes.PacketState{{
		SrcChain: chainName,
		DstChain: chainName2,
		Sequence: 1,
		Data:     []byte{0x01},
	}}

	packetSequences := []packettypes.PacketSequence{{
		SrcChain: chainName,
		DstChain: chainName2,
		Sequence: 1,
	}}

	genState := &types.GenesisState{
		ClientGenesis: clienttypes.NewGenesisState(
			[]clienttypes.IdentifiedClientState{
				clienttypes.NewIdentifiedClientState(
					chainName, xibctmtypes.NewClientState(
						suite.chainA.ChainID, xibctmtypes.DefaultTrustLevel,
						xibctesting.TrustingPeriod, xibctesting.UnbondingPeriod,
						xibctesting.MaxClockDrift, clientHeight,
						commitmenttypes.GetSDKSpecs(), xibctesting.Prefix, 0,
					),
				),
			},
			[]clienttypes.ClientConsensusStates{
				clienttypes.NewClientConsensusStates(
					chainName,
					[]clienttypes.ConsensusStateWithHeight{
						clienttypes.NewConsensusStateWithHeight(
							header.GetHeight().(clienttypes.Height),
							xibctmtypes.NewConsensusState(
								header.GetTime(), header.Header.AppHash, header.Header.NextValidatorsHash,
							),
						),
					},
				),
			},
			[]clienttypes.IdentifiedGenesisMetadata{
				clienttypes.NewIdentifiedGenesisMetadata(
					chainName,
					[]clienttypes.GenesisMetadata{
						clienttypes.NewGenesisMetadata([]byte("consensusStates/1/processedTime"), []byte("val1")),
						clienttypes.NewGenesisMetadata([]byte("consensusStates/2/processedTime"), []byte("val2")),
					},
				),
			},
			chainName2,
		),
		PacketGenesis: packettypes.NewGenesisState(packetStates, packetStates, packetStates, packetSequences),
	}

	teleport := app.Setup(false, nil)
	ctx := teleport.BaseApp.NewContext(false, tmproto.Header{Height: 1})
	xibc.InitGenesis(ctx, *teleport.XIBCKeeper, true, genState)

	exportedGenesis := xibc.ExportGenesis(ctx, *teleport.XIBCKeeper)

	suite.Require().Len(exportedGenesis.ClientGenesis.Clients, 1)
	suite.Require().Len(exportedGenesis.ClientGenesis.ClientsConsensus, 1)
	suite.Require().Len(exportedGenesis.ClientGenesis.ClientsMetadata, 1)
	suite.Require().Equal(exportedGenesis.ClientGenesis.NativeChainName, chainName2)

	suite.Require().Len(exportedGenesis.PacketGenesis.Commitments, 1)
	suite.Require().Len(exportedGenesis.PacketGenesis.Receipts, 1)
	suite.Require().Len(exportedGenesis.PacketGenesis.Acknowledgements, 1)
	suite.Require().Len(exportedGenesis.PacketGenesis.SendSequences, 1)

	xibc.ResetStates(ctx, teleport.GetKey(host.StoreKey), *teleport.XIBCKeeper)
	exportedGenesis = xibc.ExportGenesis(ctx, *teleport.XIBCKeeper)

	suite.Require().Len(exportedGenesis.ClientGenesis.Clients, 0)
	suite.Require().Len(exportedGenesis.ClientGenesis.ClientsConsensus, 0)
	suite.Require().Len(exportedGenesis.ClientGenesis.ClientsMetadata, 0)
	suite.Require().Equal(exportedGenesis.ClientGenesis.NativeChainName, chainName2)

	suite.Require().Len(exportedGenesis.PacketGenesis.Commitments, 0)
	suite.Require().Len(exportedGenesis.PacketGenesis.Receipts, 0)
	suite.Require().Len(exportedGenesis.PacketGenesis.Acknowledgements, 0)
	suite.Require().Len(exportedGenesis.PacketGenesis.SendSequences, 0)
}
