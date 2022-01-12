package types_test

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	xibctmtypes "github.com/teleport-network/teleport/x/xibc/clients/light-clients/tendermint/types"
	"github.com/teleport-network/teleport/x/xibc/core/client/types"
	commitmenttypes "github.com/teleport-network/teleport/x/xibc/core/commitment/types"
	"github.com/teleport-network/teleport/x/xibc/exported"
	xibctesting "github.com/teleport-network/teleport/x/xibc/testing"
)

type caseAny struct {
	name    string
	any     *codectypes.Any
	expPass bool
}

func (suite *TypesTestSuite) TestPackClientState() {

	testCases := []struct {
		name        string
		clientState exported.ClientState
		expPass     bool
	}{{
		"tendermint client",
		xibctmtypes.NewClientState(
			chainID,
			xibctesting.DefaultTrustLevel,
			xibctesting.TrustingPeriod,
			xibctesting.UnbondingPeriod,
			xibctesting.MaxClockDrift,
			clientHeight,
			commitmenttypes.GetSDKSpecs(),
			xibctesting.Prefix,
			0,
		),
		true,
	}, {
		"nil", nil, false,
	}}

	testCasesAny := []caseAny{}

	for _, tc := range testCases {
		clientAny, err := types.PackClientState(tc.clientState)
		if tc.expPass {
			suite.Require().NoError(err, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}

		testCasesAny = append(testCasesAny, caseAny{tc.name, clientAny, tc.expPass})
	}

	for i, tc := range testCasesAny {
		cs, err := types.UnpackClientState(tc.any)
		if tc.expPass {
			suite.Require().NoError(err, tc.name)
			suite.Require().Equal(testCases[i].clientState, cs, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
	}
}

func (suite *TypesTestSuite) TestPackConsensusState() {
	testCases := []struct {
		name           string
		consensusState exported.ConsensusState
		expPass        bool
	}{{
		"tendermint consensus",
		suite.chainA.LastHeader.ConsensusState(),
		true,
	}, {
		"nil", nil, false,
	}}

	testCasesAny := []caseAny{}

	for _, tc := range testCases {
		clientAny, err := types.PackConsensusState(tc.consensusState)
		if tc.expPass {
			suite.Require().NoError(err, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
		testCasesAny = append(testCasesAny, caseAny{tc.name, clientAny, tc.expPass})
	}

	for i, tc := range testCasesAny {
		cs, err := types.UnpackConsensusState(tc.any)
		if tc.expPass {
			suite.Require().NoError(err, tc.name)
			suite.Require().Equal(testCases[i].consensusState, cs, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
	}
}

func (suite *TypesTestSuite) TestPackHeader() {
	testCases := []struct {
		name    string
		header  exported.Header
		expPass bool
	}{{
		"tendermint header",
		suite.chainA.LastHeader,
		true,
	}, {
		"nil", nil, false,
	}}

	testCasesAny := []caseAny{}
	for _, tc := range testCases {
		clientAny, err := types.PackHeader(tc.header)
		if tc.expPass {
			suite.Require().NoError(err, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
		testCasesAny = append(testCasesAny, caseAny{tc.name, clientAny, tc.expPass})
	}

	for i, tc := range testCasesAny {
		cs, err := types.UnpackHeader(tc.any)
		if tc.expPass {
			suite.Require().NoError(err, tc.name)
			suite.Require().Equal(testCases[i].header, cs, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
	}
}
