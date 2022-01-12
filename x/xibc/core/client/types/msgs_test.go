package types_test

import (
	"testing"

	"github.com/golang/protobuf/proto" // nolint
	"github.com/stretchr/testify/suite"

	xibctmtypes "github.com/teleport-network/teleport/x/xibc/clients/light-clients/tendermint/types"
	"github.com/teleport-network/teleport/x/xibc/core/client/types"
	xibctesting "github.com/teleport-network/teleport/x/xibc/testing"
)

type TypesTestSuite struct {
	suite.Suite
	coordinator *xibctesting.Coordinator
	chainA      *xibctesting.TestChain
	chainB      *xibctesting.TestChain
}

func (suite *TypesTestSuite) SetupTest() {
	suite.coordinator = xibctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(xibctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(xibctesting.GetChainID(1))
}

func TestTypesTestSuite(t *testing.T) {
	suite.Run(t, new(TypesTestSuite))
}

// tests that different header within MsgUpdateClient can be marshaled  and unmarshaled
func (suite *TypesTestSuite) TestMarshalMsgUpdateClient() {
	var (
		msg *types.MsgUpdateClient
		err error
	)

	testCases := []struct {
		name     string
		malleate func()
	}{{
		"tendermint client",
		func() {
			msg, err = types.NewMsgUpdateClient("tendermint", suite.chainA.CurrentTMClientHeader(), suite.chainA.SenderAcc)
			suite.Require().NoError(err)
		},
	}}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			tc.malleate()

			cdc := suite.chainA.App.AppCodec()

			// marshal message
			bz, err := cdc.MarshalJSON(msg)
			suite.Require().NoError(err)

			// unmarshal message
			newMsg := &types.MsgUpdateClient{}
			err = cdc.UnmarshalJSON(bz, newMsg)
			suite.Require().NoError(err)

			suite.Require().True(proto.Equal(msg, newMsg))
		})
	}
}

func (suite *TypesTestSuite) TestMsgUpdateClient_ValidateBasic() {
	var (
		msg = &types.MsgUpdateClient{}
		err error
	)

	cases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{{
		"invalid chain-name",
		func() {
			msg.ChainName = ""
		},
		false,
	}, {
		"valid - tendermint header",
		func() {
			msg, err = types.NewMsgUpdateClient("tendermint", suite.chainA.CurrentTMClientHeader(), suite.chainA.SenderAcc)
			suite.Require().NoError(err)
		},
		true,
	}, {
		"invalid tendermint header",
		func() {
			msg, err = types.NewMsgUpdateClient("tendermint", &xibctmtypes.Header{}, suite.chainA.SenderAcc)
			suite.Require().NoError(err)
		},
		false,
	}, {
		"failed to unpack header",
		func() {
			msg.Header = nil
		},
		false,
	}, {
		"invalid signer",
		func() {
			msg.Signer = ""
		},
		false,
	}}

	for _, tc := range cases {
		tc.malleate()
		err = msg.ValidateBasic()
		if tc.expPass {
			suite.Require().NoError(err, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
	}
}
