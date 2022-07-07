package helpers

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	bitapp "github.com/bitdao-io/bitnetwork/app"
)

// SimAppChainID hardcoded chainID for simulation
const (
	SimAppChainID = "bitnetwork-app"
)

// DefaultConsensusParams defines the default Tendermint consensus params used
// in bitapp testing.
var DefaultConsensusParams = &abci.ConsensusParams{
	Block: &abci.BlockParams{
		MaxBytes: 200000,
		MaxGas:   2000000,
	},
	Evidence: &tmproto.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
		MaxBytes:        10000,
	},
	Validator: &tmproto.ValidatorParams{
		PubKeyTypes: []string{
			tmtypes.ABCIPubKeyTypeEd25519,
		},
	},
}

type EmptyAppOptions struct{}

func (EmptyAppOptions) Get(o string) interface{} { return nil }

func Setup(t *testing.T, isCheckTx bool, invCheckPeriod uint) *bitapp.BitNetwork {
	t.Helper()

	app, genesisState := setup(!isCheckTx, invCheckPeriod)
	if !isCheckTx {
		// InitChain must be called to stop deliverState from being nil
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		require.NoError(t, err)

		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return app
}

func setup(withGenesis bool, invCheckPeriod uint) (*bitapp.BitNetwork, bitapp.GenesisState) {
	db := dbm.NewMemDB()
	encCdc := params.MakeTestEncodingConfig()
	app := bitapp.NewBitNetwork(
		log.NewNopLogger(),
		db,
		nil,
		true,
		map[int64]bool{},
		bitapp.DefaultNodeHome,
		invCheckPeriod,
		encCdc,
		EmptyAppOptions{},
	)
	if withGenesis {
		return app, bitapp.NewDefaultGenesisState()
	}

	return app, bitapp.GenesisState{}
}
