package types_test

import (
	"encoding/json"

	"io/ioutil"

	bibcethtypes "github.com/teleport-network/teleport/x/xibc/clients/light-clients/eth/types"
	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	"github.com/teleport-network/teleport/x/xibc/exported"
)

var chainName = "eth"

func (suite *ETHTestSuite) TestCheckHeaderAndUpdateState() {
	var updateHeaders []*bibcethtypes.EthHeader
	updateHeadersBz, _ := ioutil.ReadFile("testdata/update_headers.json")
	err := json.Unmarshal(updateHeadersBz, &updateHeaders)
	suite.Require().NoError(err)
	suite.GreaterOrEqual(len(updateHeaders), 1)

	header := updateHeaders[0]

	number := clienttypes.NewHeight(0, header.Number.Uint64())

	clientState := exported.ClientState(&bibcethtypes.ClientState{
		Header:          header.ToHeader(),
		ChainId:         1,
		ContractAddress: []byte("0x00"),
		TrustingPeriod:  99999999,
		TimeDelay:       0,
		BlockDelay:      1,
	})

	consensusState := exported.ConsensusState(&bibcethtypes.ConsensusState{
		Timestamp: header.Time,
		Number:    number,
		Root:      header.Root[:],
	})
	err = suite.app.XIBCKeeper.ClientKeeper.CreateClient(suite.ctx, chainName, clientState, consensusState)
	suite.Require().NoError(err)
	state, exist := suite.app.XIBCKeeper.ClientKeeper.GetClientConsensusState(suite.ctx, chainName, number)
	suite.Require().True(exist)
	suite.Equal(state.GetRoot(), consensusState.GetRoot())

	for i, updateHeader := range updateHeaders[1:5] {
		protoHeader := updateHeader.ToHeader()
		suite.Require().NoError(err)
		err = suite.app.XIBCKeeper.ClientKeeper.UpdateClient(suite.ctx, chainName, &protoHeader)
		suite.Require().NoError(err)

		number.RevisionHeight = protoHeader.Height.RevisionHeight
		getClientState, exist := suite.app.XIBCKeeper.ClientKeeper.GetClientState(suite.ctx, chainName)
		suite.Require().True(exist)
		suite.Equal(getClientState.GetLatestHeight().GetRevisionHeight(), number.RevisionHeight)
		state, exist = suite.app.XIBCKeeper.ClientKeeper.GetClientConsensusState(suite.ctx, chainName, number)
		suite.Require().True(exist)
		suite.Equal(state.GetRoot(), protoHeader.Root)

		gm := make([]exported.GenesisMetadata, 0)
		callback := func(key, val []byte) bool {
			gm = append(gm, clienttypes.NewGenesisMetadata(key, val))
			return false
		}
		bibcethtypes.IteratorEthMetaDataByPrefix(suite.app.XIBCKeeper.ClientKeeper.ClientStore(suite.ctx, chainName), bibcethtypes.KeyIndexEthHeaderPrefix, callback)
		suite.Equal(len(gm), i+2)
	}
}
