package types_test

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/common"

	bibcethtypes "github.com/teleport-network/teleport/x/xibc/clients/light-clients/eth/types"
	"github.com/teleport-network/teleport/x/xibc/exported"
)

func (suite *ETHTestSuite) TestExportMetadata() {
	var updateHeaders []*bibcethtypes.EthHeader
	updateHeadersBz, _ := ioutil.ReadFile("testdata/update_headers.json")
	err := json.Unmarshal(updateHeadersBz, &updateHeaders)
	suite.Require().NoError(err)
	suite.GreaterOrEqual(len(updateHeaders), 1)
	header := updateHeaders[0]
	clientState := exported.ClientState(&bibcethtypes.ClientState{
		Header:          header.ToHeader(),
		ChainId:         1,
		ContractAddress: []byte("0x00"),
		TrustingPeriod:  200000,
		TimeDelay:       0,
		BlockDelay:      1,
	})
	suite.app.XIBCKeeper.ClientKeeper.SetClientState(suite.ctx, chainName, clientState)
	gm := clientState.ExportMetadata(suite.app.XIBCKeeper.ClientKeeper.ClientStore(suite.ctx, chainName))
	suite.Require().Nil(gm, "client with no metadata returned non-nil exported metadata")

	protoHeader := header.ToHeader()
	store := suite.app.XIBCKeeper.ClientKeeper.ClientStore(suite.ctx, chainName)
	headerBytes, err := suite.app.AppCodec().MarshalInterface(&protoHeader)
	suite.Require().NoError(err)

	bibcethtypes.SetEthHeaderIndex(store, protoHeader, headerBytes)
	bibcethtypes.SetEthConsensusRoot(store, protoHeader.Height.RevisionHeight, protoHeader.ToEthHeader().Root, header.Hash())

	gm = clientState.ExportMetadata(suite.app.XIBCKeeper.ClientKeeper.ClientStore(suite.ctx, chainName))
	suite.Require().NotNil(gm, "client with metadata returned nil exported metadata")
	suite.Require().Len(gm, 2, "exported metadata has unexpected length")

	suite.Require().Equal(bibcethtypes.EthHeaderIndexKey(protoHeader.Hash(), protoHeader.Height.RevisionHeight), gm[0].GetKey(), "metadata has unexpected key")
	suite.Require().Equal(headerBytes, gm[0].GetValue(), "metadata has unexpected value")

	suite.Require().Equal(bibcethtypes.EthRootMainKey(common.BytesToHash(protoHeader.Root), protoHeader.Height.RevisionHeight), gm[1].GetKey(), "metadata has unexpected key")
	suite.Require().Equal(gm[0].GetKey(), gm[1].GetValue(), "metadata has unexpected value")
}
