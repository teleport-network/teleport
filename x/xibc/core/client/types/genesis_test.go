package types_test

import (
	"time"

	tmtypes "github.com/tendermint/tendermint/types"

	ibctmtypes "github.com/teleport-network/teleport/x/xibc/clients/light-clients/tendermint/types"
	client "github.com/teleport-network/teleport/x/xibc/core/client"
	"github.com/teleport-network/teleport/x/xibc/core/client/types"
	commitmenttypes "github.com/teleport-network/teleport/x/xibc/core/commitment/types"
	xibctesting "github.com/teleport-network/teleport/x/xibc/testing"
	xibctestingmock "github.com/teleport-network/teleport/x/xibc/testing/mock"
)

const (
	chainID          = "chainID"
	tmChainName0     = "tendermint-0"
	tmChainName1     = "tendermint-1"
	invalidChainName = "myclient/0"
	chainName        = tmChainName0
	height           = 10
)

var clientHeight = types.NewHeight(0, 10)

func (suite *TypesTestSuite) TestMarshalGenesisState() {
	cdc := suite.chainA.App.AppCodec()
	path := xibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupClients(path)
	err := path.EndpointA.UpdateClient()
	suite.Require().NoError(err)

	genesis := client.ExportGenesis(suite.chainA.GetContext(), suite.chainA.App.XIBCKeeper.ClientKeeper)

	bz, err := cdc.MarshalJSON(&genesis)
	suite.Require().NoError(err)
	suite.Require().NotNil(bz)

	var gs types.GenesisState
	err = cdc.UnmarshalJSON(bz, &gs)
	suite.Require().NoError(err)
}

func (suite *TypesTestSuite) TestValidateGenesis() {
	privVal := xibctestingmock.NewPV()
	pubKey, err := privVal.GetPubKey()
	suite.Require().NoError(err)

	now := time.Now().UTC()

	val := tmtypes.NewValidator(pubKey, 10)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{val})

	heightMinus1 := types.NewHeight(0, height-1)
	header := suite.chainA.CreateTMClientHeader(
		chainID,
		int64(clientHeight.RevisionHeight),
		heightMinus1,
		now,
		valSet,
		valSet,
		[]tmtypes.PrivValidator{privVal},
	)

	testCases := []struct {
		name     string
		genState types.GenesisState
		expPass  bool
	}{{
		name:     "default",
		genState: types.DefaultGenesisState(),
		expPass:  true,
	}, {
		name: "valid custom genesis",
		genState: types.NewGenesisState(
			[]types.IdentifiedClientState{
				types.NewIdentifiedClientState(
					tmChainName0,
					ibctmtypes.NewClientState(
						chainID,
						xibctesting.DefaultTrustLevel,
						xibctesting.TrustingPeriod,
						xibctesting.UnbondingPeriod,
						xibctesting.MaxClockDrift,
						clientHeight,
						commitmenttypes.GetSDKSpecs(),
						xibctesting.Prefix,
						0,
					),
				),
			},
			[]types.ClientConsensusStates{
				types.NewClientConsensusStates(
					tmChainName0,
					[]types.ConsensusStateWithHeight{
						types.NewConsensusStateWithHeight(
							header.GetHeight().(types.Height),
							ibctmtypes.NewConsensusState(
								header.GetTime(),
								header.Header.GetAppHash(),
								header.Header.NextValidatorsHash,
							),
						),
					},
				),
			},
			[]types.IdentifiedGenesisMetadata{
				types.NewIdentifiedGenesisMetadata(
					chainName,
					[]types.GenesisMetadata{
						types.NewGenesisMetadata([]byte("key1"), []byte("val1")),
						types.NewGenesisMetadata([]byte("key2"), []byte("val2")),
					},
				),
			},
			tmChainName1,
		),
		expPass: true,
	}, {
		name: "invalid chain-name",
		genState: types.NewGenesisState(
			[]types.IdentifiedClientState{
				types.NewIdentifiedClientState(
					invalidChainName,
					ibctmtypes.NewClientState(
						chainID,
						ibctmtypes.DefaultTrustLevel,
						xibctesting.TrustingPeriod,
						xibctesting.UnbondingPeriod,
						xibctesting.MaxClockDrift,
						clientHeight,
						commitmenttypes.GetSDKSpecs(),
						xibctesting.Prefix,
						0,
					),
				),
			},
			[]types.ClientConsensusStates{
				types.NewClientConsensusStates(
					invalidChainName,
					[]types.ConsensusStateWithHeight{
						types.NewConsensusStateWithHeight(
							header.GetHeight().(types.Height),
							ibctmtypes.NewConsensusState(
								header.GetTime(),
								header.Header.GetAppHash(),
								header.Header.NextValidatorsHash,
							),
						),
					},
				),
			},
			nil,
			tmChainName1,
		),
		expPass: false,
	}, {
		name: "consensus state chain name does not match chain name in genesis clients",
		genState: types.NewGenesisState(
			[]types.IdentifiedClientState{
				types.NewIdentifiedClientState(
					tmChainName0,
					ibctmtypes.NewClientState(
						chainID,
						ibctmtypes.DefaultTrustLevel,
						xibctesting.TrustingPeriod,
						xibctesting.UnbondingPeriod,
						xibctesting.MaxClockDrift,
						clientHeight,
						commitmenttypes.GetSDKSpecs(),
						xibctesting.Prefix,
						0,
					),
				),
			},
			[]types.ClientConsensusStates{
				types.NewClientConsensusStates(
					tmChainName1,
					[]types.ConsensusStateWithHeight{
						types.NewConsensusStateWithHeight(
							types.NewHeight(0, 1),
							ibctmtypes.NewConsensusState(
								header.GetTime(),
								header.Header.GetAppHash(),
								header.Header.NextValidatorsHash,
							),
						),
					},
				),
			},
			nil,
			tmChainName1,
		),
		expPass: false,
	}, {
		name: "invalid consensus state height",
		genState: types.NewGenesisState(
			[]types.IdentifiedClientState{
				types.NewIdentifiedClientState(
					tmChainName0, ibctmtypes.NewClientState(
						chainID,
						ibctmtypes.DefaultTrustLevel,
						xibctesting.TrustingPeriod,
						xibctesting.UnbondingPeriod,
						xibctesting.MaxClockDrift,
						clientHeight,
						commitmenttypes.GetSDKSpecs(),
						xibctesting.Prefix,
						0,
					),
				),
			},
			[]types.ClientConsensusStates{
				types.NewClientConsensusStates(
					tmChainName0,
					[]types.ConsensusStateWithHeight{
						types.NewConsensusStateWithHeight(
							types.ZeroHeight(),
							ibctmtypes.NewConsensusState(
								header.GetTime(),
								header.Header.GetAppHash(),
								header.Header.NextValidatorsHash,
							),
						),
					},
				),
			},
			nil,
			tmChainName1,
		),
		expPass: false,
	}, {
		name: "invalid consensus state",
		genState: types.NewGenesisState(
			[]types.IdentifiedClientState{
				types.NewIdentifiedClientState(
					tmChainName0, ibctmtypes.NewClientState(
						chainID,
						ibctmtypes.DefaultTrustLevel,
						xibctesting.TrustingPeriod,
						xibctesting.UnbondingPeriod,
						xibctesting.MaxClockDrift,
						clientHeight,
						commitmenttypes.GetSDKSpecs(),
						xibctesting.Prefix,
						0,
					),
				),
			},
			[]types.ClientConsensusStates{
				types.NewClientConsensusStates(
					tmChainName0,
					[]types.ConsensusStateWithHeight{
						types.NewConsensusStateWithHeight(
							types.NewHeight(0, 1),
							ibctmtypes.NewConsensusState(
								time.Time{},
								header.Header.GetAppHash(),
								header.Header.NextValidatorsHash,
							),
						),
					},
				),
			},
			nil,
			tmChainName1,
		),
		expPass: false,
	}, {
		name: "metadata chain-name does not match a genesis client",
		genState: types.NewGenesisState(
			[]types.IdentifiedClientState{
				types.NewIdentifiedClientState(
					chainName, ibctmtypes.NewClientState(
						chainID,
						xibctesting.DefaultTrustLevel,
						xibctesting.TrustingPeriod,
						xibctesting.UnbondingPeriod,
						xibctesting.MaxClockDrift,
						clientHeight,
						commitmenttypes.GetSDKSpecs(),
						xibctesting.Prefix,
						0,
					),
				),
			},
			[]types.ClientConsensusStates{
				types.NewClientConsensusStates(
					chainName,
					[]types.ConsensusStateWithHeight{
						types.NewConsensusStateWithHeight(
							header.GetHeight().(types.Height),
							ibctmtypes.NewConsensusState(
								header.GetTime(),
								header.Header.GetAppHash(),
								header.Header.NextValidatorsHash,
							),
						),
					},
				),
			},
			[]types.IdentifiedGenesisMetadata{
				types.NewIdentifiedGenesisMetadata(
					"wrongclientid",
					[]types.GenesisMetadata{
						types.NewGenesisMetadata([]byte("key1"), []byte("val1")),
						types.NewGenesisMetadata([]byte("key2"), []byte("val2")),
					},
				),
			},
			tmChainName1,
		),
		expPass: false,
	}, {
		name: "invalid metadata",
		genState: types.NewGenesisState(
			[]types.IdentifiedClientState{
				types.NewIdentifiedClientState(
					chainName, ibctmtypes.NewClientState(
						chainID,
						ibctmtypes.DefaultTrustLevel,
						xibctesting.TrustingPeriod,
						xibctesting.UnbondingPeriod,
						xibctesting.MaxClockDrift,
						clientHeight,
						commitmenttypes.GetSDKSpecs(),
						xibctesting.Prefix,
						0,
					),
				),
			},
			[]types.ClientConsensusStates{
				types.NewClientConsensusStates(
					chainName,
					[]types.ConsensusStateWithHeight{
						types.NewConsensusStateWithHeight(
							header.GetHeight().(types.Height),
							ibctmtypes.NewConsensusState(
								header.GetTime(),
								header.Header.GetAppHash(),
								header.Header.NextValidatorsHash,
							),
						),
					},
				),
			},
			[]types.IdentifiedGenesisMetadata{
				types.NewIdentifiedGenesisMetadata(
					chainName,
					[]types.GenesisMetadata{
						types.NewGenesisMetadata([]byte(""), []byte("val1")),
						types.NewGenesisMetadata([]byte("key2"), []byte("val2")),
					},
				),
			},
			tmChainName1,
		),
	}, {
		name: "failed to parse client identifier in client state loop",
		genState: types.NewGenesisState(
			[]types.IdentifiedClientState{
				types.NewIdentifiedClientState(
					"my-client",
					ibctmtypes.NewClientState(
						chainID,
						xibctesting.DefaultTrustLevel,
						xibctesting.TrustingPeriod,
						xibctesting.UnbondingPeriod,
						xibctesting.MaxClockDrift,
						clientHeight,
						commitmenttypes.GetSDKSpecs(),
						xibctesting.Prefix,
						0,
					),
				),
			},
			[]types.ClientConsensusStates{
				types.NewClientConsensusStates(
					tmChainName0,
					[]types.ConsensusStateWithHeight{
						types.NewConsensusStateWithHeight(
							header.GetHeight().(types.Height),
							ibctmtypes.NewConsensusState(
								header.GetTime(),
								header.Header.GetAppHash(),
								header.Header.NextValidatorsHash,
							),
						),
					},
				),
			},
			nil,
			tmChainName1,
		),
		expPass: false,
	}}

	for _, tc := range testCases {
		if tc.expPass {
			suite.Require().NoError(tc.genState.Validate(), tc.name)
		} else {
			suite.Require().Error(tc.genState.Validate(), tc.name)
		}
	}
}
