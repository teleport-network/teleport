package types_test

import (
	"time"

	"github.com/teleport-network/teleport/x/xibc/clients/light-clients/tendermint/types"
	"github.com/teleport-network/teleport/x/xibc/exported"
)

func (suite *TendermintTestSuite) TestConsensusStateValidateBasic() {
	testCases := []struct {
		msg            string
		consensusState *types.ConsensusState
		expectPass     bool
	}{{
		"success",
		&types.ConsensusState{
			Timestamp:          suite.now,
			Root:               []byte("app_hash"),
			NextValidatorsHash: suite.valsHash,
		},
		true,
	}, {
		"root is nil",
		&types.ConsensusState{
			Timestamp:          suite.now,
			Root:               []byte{},
			NextValidatorsHash: suite.valsHash,
		},
		false,
	}, {
		"root is empty",
		&types.ConsensusState{
			Timestamp:          suite.now,
			Root:               []byte{},
			NextValidatorsHash: suite.valsHash,
		},
		false,
	}, {
		"nextvalshash is invalid",
		&types.ConsensusState{
			Timestamp:          suite.now,
			Root:               []byte("app_hash"),
			NextValidatorsHash: []byte("hi"),
		},
		false,
	}, {
		"timestamp is zero",
		&types.ConsensusState{
			Timestamp:          time.Time{},
			Root:               []byte("app_hash"),
			NextValidatorsHash: suite.valsHash,
		},
		false,
	}}

	for i, tc := range testCases {
		// check just to increase coverage
		suite.Require().Equal(exported.Tendermint, tc.consensusState.ClientType())
		suite.Require().Equal(tc.consensusState.GetRoot(), tc.consensusState.Root)

		if tc.expectPass {
			suite.Require().NoError(tc.consensusState.ValidateBasic(), "valid test case %d failed: %s", i, tc.msg)
		} else {
			suite.Require().Error(tc.consensusState.ValidateBasic(), "invalid test case %d passed: %s", i, tc.msg)
		}
	}
}
