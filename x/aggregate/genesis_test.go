package aggregate_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmversion "github.com/tendermint/tendermint/proto/tendermint/version"
	"github.com/tendermint/tendermint/version"

	"github.com/tharsis/ethermint/tests"
	feemarkettypes "github.com/tharsis/ethermint/x/feemarket/types"

	"github.com/teleport-network/teleport/app"
	"github.com/teleport-network/teleport/x/aggregate"
	"github.com/teleport-network/teleport/x/aggregate/types"
)

type GenesisTestSuite struct {
	suite.Suite
	ctx     sdk.Context
	app     *app.Teleport
	genesis types.GenesisState
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

func (suite *GenesisTestSuite) SetupTest() {
	// consensus key
	consAddress := sdk.ConsAddress(tests.GenerateAddress().Bytes())

	suite.app = app.Setup(false, feemarkettypes.DefaultGenesisState())
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{
		Height:          1,
		ChainID:         "teleport_9000-1",
		Time:            time.Now().UTC(),
		ProposerAddress: consAddress.Bytes(),

		Version: tmversion.Consensus{
			Block: version.BlockProtocol,
		},
		LastBlockId: tmproto.BlockID{
			Hash: tmhash.Sum([]byte("block_id")),
			PartSetHeader: tmproto.PartSetHeader{
				Total: 11,
				Hash:  tmhash.Sum([]byte("partset_header")),
			},
		},
		AppHash:            tmhash.Sum([]byte("app")),
		DataHash:           tmhash.Sum([]byte("data")),
		EvidenceHash:       tmhash.Sum([]byte("evidence")),
		ValidatorsHash:     tmhash.Sum([]byte("validators")),
		NextValidatorsHash: tmhash.Sum([]byte("next_validators")),
		ConsensusHash:      tmhash.Sum([]byte("consensus")),
		LastResultsHash:    tmhash.Sum([]byte("last_result")),
	})

	suite.genesis = *types.DefaultGenesisState()
}

func (suite *GenesisTestSuite) TestERC20InitGenesis() {
	testCases := []struct {
		name         string
		genesisState types.GenesisState
	}{
		{
			"empty genesis",
			types.GenesisState{},
		},
		{
			"default genesis",
			*types.DefaultGenesisState(),
		},
		{
			"custom genesis",
			types.NewGenesisState(
				types.DefaultParams(),
				[]types.TokenPair{
					{
						ERC20Address:  "0x5dCA2483280D9727c80b5518faC4556617fb19ZZ",
						Denoms:        []string{"coin"},
						Enabled:       true,
						ContractOwner: types.OWNER_MODULE,
					},
				}),
		},
	}

	for _, tc := range testCases {

		suite.Require().NotPanics(func() {
			aggregate.InitGenesis(suite.ctx, suite.app.AggregateKeeper, suite.app.AccountKeeper, tc.genesisState)
		})
		params := suite.app.AggregateKeeper.GetParams(suite.ctx)

		tokenPairs := suite.app.AggregateKeeper.GetTokenPairs(suite.ctx)
		suite.Require().Equal(tc.genesisState.Params, params)
		if len(tokenPairs) > 0 {
			suite.Require().Equal(tc.genesisState.TokenPairs, tokenPairs)
		} else {
			suite.Require().Len(tc.genesisState.TokenPairs, 0)
		}
	}
}

func (suite *GenesisTestSuite) TestErc20ExportGenesis() {
	testGenCases := []struct {
		name         string
		genesisState types.GenesisState
	}{
		{
			"empty genesis",
			types.GenesisState{},
		},
		{
			"default genesis",
			*types.DefaultGenesisState(),
		},
		{
			"custom genesis",
			types.NewGenesisState(
				types.DefaultParams(),
				[]types.TokenPair{
					{
						ERC20Address:  "0x5dCA2483280D9727c80b5518faC4556617fb19ZZ",
						Denoms:        []string{"coin"},
						Enabled:       true,
						ContractOwner: types.OWNER_MODULE,
					},
				}),
		},
	}

	for _, tc := range testGenCases {
		aggregate.InitGenesis(suite.ctx, suite.app.AggregateKeeper, suite.app.AccountKeeper, tc.genesisState)
		suite.Require().NotPanics(func() {
			genesisExported := aggregate.ExportGenesis(suite.ctx, suite.app.AggregateKeeper)
			params := suite.app.AggregateKeeper.GetParams(suite.ctx)
			suite.Require().Equal(genesisExported.Params, params)

			tokenPairs := suite.app.AggregateKeeper.GetTokenPairs(suite.ctx)
			if len(tokenPairs) > 0 {
				suite.Require().Equal(genesisExported.TokenPairs, tokenPairs)
			} else {
				suite.Require().Len(genesisExported.TokenPairs, 0)
			}
		})
		// }
	}
}
