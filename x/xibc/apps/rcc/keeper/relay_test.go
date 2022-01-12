package keeper_test

import (
	"fmt"
)

func (suite *KeeperTestSuite) TestRecvPacket() {
	testCases := []struct {
		name     string
		malleate func()
	}{
		{
			"",
			func() {
				// TODO
			},
		},
		// TODO
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest()
			tc.malleate()
		})
	}
}

func (suite *KeeperTestSuite) TestAcknowledgementPacket() {
	testCases := []struct {
		name     string
		malleate func()
	}{
		{
			"ack success",
			func() {
				// TODO
			},
		},
		// TODO
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest()
			tc.malleate()
		})
	}
}
