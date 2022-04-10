package types

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type GenesisTestSuite struct {
	suite.Suite
}

func (suite *GenesisTestSuite) SetupTest() {
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

func (suite *GenesisTestSuite) TestValidateGenesis() {
	newGen := NewGenesisState(DefaultParams(), []TokenPair{})

	testCases := []struct {
		name     string
		genState *GenesisState
		expPass  bool
	}{
		{
			name:     "valid genesis constructor",
			genState: &newGen,
			expPass:  true,
		},
		{
			name:     "default",
			genState: DefaultGenesisState(),
			expPass:  true,
		},
		{
			name: "valid genesis",
			genState: &GenesisState{
				Params:     DefaultParams(),
				TokenPairs: []TokenPair{},
			},
			expPass: true,
		},
		{
			name: "valid genesis - with tokens pairs",
			genState: &GenesisState{
				Params: DefaultParams(),
				TokenPairs: []TokenPair{
					{
						ERC20Address: "0xdac17f958d2ee523a2206206994597c13d831ec7",
						Denoms:       []string{"usdt"},
						Enabled:      true,
					},
				},
			},
			expPass: true,
		},
		{
			name: "invalid genesis - duplicated token pair",
			genState: &GenesisState{
				Params: DefaultParams(),
				TokenPairs: []TokenPair{
					{
						ERC20Address: "0xdac17f958d2ee523a2206206994597c13d831ec7",
						Denoms:       []string{"usdt"},
						Enabled:      true,
					},
					{
						ERC20Address: "0xdac17f958d2ee523a2206206994597c13d831ec7",
						Denoms:       []string{"usdt"},
						Enabled:      true,
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid genesis - duplicated token pair",
			genState: &GenesisState{
				Params: DefaultParams(),
				TokenPairs: []TokenPair{
					{
						ERC20Address: "0xdac17f958d2ee523a2206206994597c13d831ec7",
						Denoms:       []string{"usdt"},
						Enabled:      true,
					},
					{
						ERC20Address: "0xdac17f958d2ee523a2206206994597c13d831ec7",
						Denoms:       []string{"usdt2"},
						Enabled:      true,
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid genesis - duplicated token pair",
			genState: &GenesisState{
				Params: DefaultParams(),
				TokenPairs: []TokenPair{
					{
						ERC20Address: "0xdac17f958d2ee523a2206206994597c13d831ec7",
						Denoms:       []string{"usdt"},
						Enabled:      true,
					},
					{
						ERC20Address: "0xB8c77482e45F1F44dE1745F52C74426C631bDD52",
						Denoms:       []string{"usdt"},
						Enabled:      true,
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid genesis - invalid token pair",
			genState: &GenesisState{
				Params: DefaultParams(),
				TokenPairs: []TokenPair{
					{
						ERC20Address: "0xinvalidaddress",
						Denoms:       []string{"bad"},
						Enabled:      true,
					},
				},
			},
			expPass: false,
		},
		{
			// Voting period cant be zero
			name:     "empty genesis",
			genState: &GenesisState{},
			expPass:  true,
		},
	}

	for _, tc := range testCases {
		if tc.expPass {
			suite.Require().NoError(tc.genState.Validate(), tc.name)
		} else {
			suite.Require().Error(tc.genState.Validate(), tc.name)
		}
	}
}
