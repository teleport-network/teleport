package simulation_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/teleport-network/teleport/app"
	xibctmtypes "github.com/teleport-network/teleport/x/xibc/clients/light-clients/tendermint/types"
	"github.com/teleport-network/teleport/x/xibc/core/host"
	"github.com/teleport-network/teleport/x/xibc/simulation"
)

func TestDecodeStore(t *testing.T) {
	teleport := app.Setup(false, nil)
	dec := simulation.NewDecodeStore(*teleport.XIBCKeeper)

	chainName := "clientidone"

	clientState := &xibctmtypes.ClientState{}

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{{
			Key:   host.FullClientStateKey(chainName),
			Value: teleport.XIBCKeeper.ClientKeeper.MustMarshalClientState(clientState),
		}, {
			Key:   []byte{0x99},
			Value: []byte{0x99},
		}},
	}
	tests := []struct {
		name        string
		expectedLog string
	}{
		{"ClientState", fmt.Sprintf("ClientState A: %v\nClientState B: %v", clientState, clientState)},
		{"other", ""},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if i == len(tests)-1 {
				require.Panics(t, func() { dec(kvPairs.Pairs[i], kvPairs.Pairs[i]) }, tt.name)
			} else {
				require.Equal(t, tt.expectedLog, dec(kvPairs.Pairs[i], kvPairs.Pairs[i]), tt.name)
			}
		})
	}
}
