package types_test

import (
	"time"

	tmprotocrypto "github.com/tendermint/tendermint/proto/tendermint/crypto"

	"github.com/teleport-network/teleport/x/xibc/clients/light-clients/tendermint/types"
	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	"github.com/teleport-network/teleport/x/xibc/exported"
)

func (suite *TendermintTestSuite) TestGetHeight() {
	header := suite.chainA.LastHeader
	suite.Require().NotEqual(uint64(0), header.GetHeight())
}

func (suite *TendermintTestSuite) TestGetTime() {
	header := suite.chainA.LastHeader
	suite.Require().NotEqual(time.Time{}, header.GetTime())
}

func (suite *TendermintTestSuite) TestHeaderValidateBasic() {
	var header *types.Header

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{{
		"valid header",
		func() {},
		true,
	}, {
		"header is nil", func() {
			header.Header = nil
		},
		false,
	}, {
		"signed header is nil", func() {
			header.SignedHeader = nil
		},
		false,
	}, {
		"SignedHeaderFromProto failed", func() {
			header.SignedHeader.Commit.Height = -1
		},
		false,
	}, {
		"signed header failed tendermint ValidateBasic", func() {
			header = suite.chainA.LastHeader
			header.SignedHeader.Commit = nil
		},
		false,
	}, {
		"trusted height is greater than header height", func() {
			header.TrustedHeight = header.GetHeight().(clienttypes.Height).Increment().(clienttypes.Height)
		},
		false,
	}, {
		"validator set nil", func() {
			header.ValidatorSet = nil
		},
		false,
	}, {
		"ValidatorSetFromProto failed", func() {
			header.ValidatorSet.Validators[0].PubKey = tmprotocrypto.PublicKey{}
		},
		false,
	}, {
		"header validator hash does not equal hash of validator set", func() {
			// use chainB's randomly generated validator set
			header.ValidatorSet = suite.chainB.LastHeader.ValidatorSet
		},
		false,
	}}

	suite.Require().Equal(exported.Tendermint, suite.header.ClientType())

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			header = suite.chainA.LastHeader // must be explicitly changed in malleate
			tc.malleate()

			if tc.expPass {
				suite.Require().NoError(header.ValidateBasic(), tc.name)
			} else {
				suite.Require().Error(header.ValidateBasic(), tc.name)
			}
		})
	}
}
