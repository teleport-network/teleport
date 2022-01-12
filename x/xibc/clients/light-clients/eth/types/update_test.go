package types_test

import (
	"encoding/json"

	"io/ioutil"

	xibcethtypes "github.com/teleport-network/teleport/x/xibc/clients/light-clients/eth/types"
	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	"github.com/teleport-network/teleport/x/xibc/exported"
)

var chainName = "eth"

func (suite *ETHTestSuite) TestCheckHeaderAndUpdateState() {
	var updateHeaders []*xibcethtypes.EthHeader
	updateHeadersBz, _ := ioutil.ReadFile("testdata/update_headers.json")
	err := json.Unmarshal(updateHeadersBz, &updateHeaders)
	suite.Require().NoError(err)
	suite.GreaterOrEqual(len(updateHeaders), 1)

	header := updateHeaders[0]

	number := clienttypes.NewHeight(0, header.Number.Uint64())

	clientState := exported.ClientState(&xibcethtypes.ClientState{
		Header:          header.ToHeader(),
		ChainId:         1,
		ContractAddress: []byte("0x00"),
		TrustingPeriod:  200000,
		TimeDelay:       0,
		BlockDelay:      1,
	})

	consensusState := exported.ConsensusState(&xibcethtypes.ConsensusState{
		Timestamp: header.Time,
		Number:    number,
		Root:      header.Root[:],
	})

	suite.app.XIBCKeeper.ClientKeeper.SetClientConsensusState(suite.ctx, chainName, number, consensusState)
	protoHeader := header.ToHeader()
	store := suite.app.XIBCKeeper.ClientKeeper.ClientStore(suite.ctx, chainName)
	headerBytes, err := suite.app.AppCodec().MarshalInterface(&protoHeader)
	suite.Require().NoError(err)

	xibcethtypes.SetEthHeaderIndex(store, protoHeader, headerBytes)
	xibcethtypes.SetEthConsensusRoot(store, protoHeader.Height.RevisionHeight, protoHeader.ToEthHeader().Root, header.Hash())

	for _, updateHeader := range updateHeaders[1:5] {
		protoHeader := updateHeader.ToHeader()
		suite.Require().NoError(err)

		clientState, consensusState, err = clientState.CheckHeaderAndUpdateState(
			suite.ctx,
			suite.app.AppCodec(),
			suite.app.XIBCKeeper.ClientKeeper.ClientStore(suite.ctx, chainName), // pass in chainName prefixed clientStore
			&protoHeader,
		)

		suite.Require().NoError(err)

		number.RevisionHeight = protoHeader.Height.RevisionHeight
		suite.app.XIBCKeeper.ClientKeeper.SetClientConsensusState(suite.ctx, chainName, number, consensusState)

		suite.Require().Equal(updateHeader.Number.Uint64(), clientState.GetLatestHeight().GetRevisionHeight())
	}
}
