package types_test

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	xibctmtypes "github.com/teleport-network/teleport/x/xibc/clients/light-clients/tendermint/types"
	"github.com/teleport-network/teleport/x/xibc/core/client/types"
	xibctesting "github.com/teleport-network/teleport/x/xibc/testing"
)

func (suite *TypesTestSuite) TestNewCreateClientProposal() {
	p, err := types.NewCreateClientProposal(xibctesting.Title, xibctesting.Description, chainName, &xibctmtypes.ClientState{}, &xibctmtypes.ConsensusState{})
	suite.Require().NoError(err)
	suite.Require().NotNil(p)

	p, err = types.NewCreateClientProposal(xibctesting.Title, xibctesting.Description, chainName, nil, nil)
	suite.Require().Error(err)
	suite.Require().Nil(p)
}

// tests a client update proposal can be marshaled and unmarshaled, and the
// client state can be unpacked
func (suite *TypesTestSuite) TestMarshalCreateClientProposalProposal() {
	path := xibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupClients(path)

	clientState := path.EndpointA.GetClientState()
	consensusState := path.EndpointA.GetConsensusState(clientState.GetLatestHeight())
	// create proposal
	proposal, err := types.NewCreateClientProposal("update XIBC client", "description", "chain-name", clientState, consensusState)
	suite.Require().NoError(err)

	// create codec
	ir := codectypes.NewInterfaceRegistry()
	types.RegisterInterfaces(ir)
	govtypes.RegisterInterfaces(ir)
	xibctmtypes.RegisterInterfaces(ir)
	cdc := codec.NewProtoCodec(ir)

	// marshal message
	bz, err := cdc.MarshalJSON(proposal)
	suite.Require().NoError(err)

	// unmarshal proposal
	newProposal := &types.CreateClientProposal{}
	err = cdc.UnmarshalJSON(bz, newProposal)
	suite.Require().NoError(err)
}
