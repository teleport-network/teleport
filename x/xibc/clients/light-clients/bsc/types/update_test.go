package types_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	xibcbsctypes "github.com/teleport-network/teleport/x/xibc/clients/light-clients/bsc/types"
	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	"github.com/teleport-network/teleport/x/xibc/exported"
)

var (
	chainName = "bsc"
	epoch     = uint64(200)
)

func (suite *BSCTestSuite) TestCheckHeaderAndUpdateState() {
	var genesisState GenesisState
	genesisStateBz, _ := ioutil.ReadFile("testdata/genesis_state.json")
	err := json.Unmarshal(genesisStateBz, &genesisState)
	suite.NoError(err)

	header := genesisState.GenesisHeader
	genesisValidatorHeader := genesisState.GenesisValidatorHeader

	genesisValidators, err := xibcbsctypes.ParseValidators(genesisValidatorHeader.Extra)
	suite.NoError(err)

	number := clienttypes.NewHeight(0, header.Number.Uint64())

	clientState := exported.ClientState(&xibcbsctypes.ClientState{
		Header:          header.ToHeader(),
		ChainId:         56,
		Epoch:           epoch,
		BlockInteval:    3,
		Validators:      genesisValidators,
		ContractAddress: []byte("0x00"),
		TrustingPeriod:  999999999,
	})

	consensusState := exported.ConsensusState(&xibcbsctypes.ConsensusState{
		Timestamp: header.Time,
		Height:    number,
		Root:      header.Root[:],
	})
	err = suite.app.XIBCKeeper.ClientKeeper.CreateClient(suite.ctx, chainName, clientState, consensusState)
	suite.NoError(err)
	state, exist := suite.app.XIBCKeeper.ClientKeeper.GetClientConsensusState(suite.ctx, chainName, number)
	suite.True(exist)
	equal := bytes.Equal(state.GetRoot(), consensusState.GetRoot())
	suite.True(equal)

	var updateHeaders []*xibcbsctypes.BscHeader
	updateHeadersBz, _ := ioutil.ReadFile("testdata/update_headers.json")
	err = json.Unmarshal(updateHeadersBz, &updateHeaders)
	suite.NoError(err)
	suite.Equal(int(1.5*float64(epoch)), len(updateHeaders))
	for i, updateHeader := range updateHeaders {
		protoHeader := updateHeader.ToHeader()
		suite.NoError(err)

		err = suite.app.XIBCKeeper.ClientKeeper.UpdateClient(suite.ctx, chainName, &protoHeader)
		suite.NoError(err)

		number.RevisionHeight = protoHeader.Height.RevisionHeight
		state, exist = suite.app.XIBCKeeper.ClientKeeper.GetClientConsensusState(suite.ctx, chainName, number)
		suite.True(exist)
		suite.Equal(state.GetRoot(), protoHeader.Root)

		clientState, exist = suite.app.XIBCKeeper.ClientKeeper.GetClientState(suite.ctx, chainName)
		suite.True(exist)
		suite.Equal(clientState.GetLatestHeight().GetRevisionHeight(), protoHeader.Height.RevisionHeight)

		recentSigners, err2 := xibcbsctypes.GetRecentSigners(suite.app.XIBCKeeper.ClientKeeper.ClientStore(suite.ctx, chainName))
		suite.NoError(err2)

		validatorCount := len(clientState.(*xibcbsctypes.ClientState).Validators)
		if i+2 <= validatorCount/2+1 {
			suite.Equal(i+2, len(recentSigners))
		} else {
			suite.Equal(validatorCount/2+1, len(recentSigners))
		}
	}
}

type GenesisState struct {
	GenesisHeader          *xibcbsctypes.BscHeader `json:"genesis_header"`
	GenesisValidatorHeader *xibcbsctypes.BscHeader `json:"genesis_validator_header"`
}
