package types_test

import (
	ibctmtypes "github.com/teleport-network/teleport/x/xibc/clients/light-clients/tendermint/types"
	"github.com/teleport-network/teleport/x/xibc/core/client/types"
)

func (suite *TypesTestSuite) TestMarshalHeader() {
	cdc := suite.chainA.App.AppCodec()
	h := &ibctmtypes.Header{
		TrustedHeight: types.NewHeight(4, 100),
	}

	// marshal header
	bz, err := types.MarshalHeader(cdc, h)
	suite.Require().NoError(err)

	// unmarshal header
	newHeader, err := types.UnmarshalHeader(cdc, bz)
	suite.Require().NoError(err)

	suite.Require().Equal(h, newHeader)

	// use invalid bytes
	invalidHeader, err := types.UnmarshalHeader(cdc, []byte("invalid bytes"))
	suite.Require().Error(err)
	suite.Require().Nil(invalidHeader)
}
