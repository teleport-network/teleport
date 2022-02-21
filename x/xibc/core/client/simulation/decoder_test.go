package simulation_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/teleport-network/teleport/app"
	xibctmtypes "github.com/teleport-network/teleport/x/xibc/clients/light-clients/tendermint/types"
	"github.com/teleport-network/teleport/x/xibc/core/client/simulation"
	"github.com/teleport-network/teleport/x/xibc/core/client/types"
	"github.com/teleport-network/teleport/x/xibc/core/host"
)

func TestDecodeStore(t *testing.T) {
	teleport := app.Setup(false, nil)
	chainName := "clientidone"

	height := types.NewHeight(0, 10)

	clientState := &xibctmtypes.ClientState{}
	consState := &xibctmtypes.ConsensusState{
		Timestamp: time.Now().UTC(),
	}

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{{
			Key:   host.FullClientStateKey(chainName),
			Value: teleport.XIBCKeeper.ClientKeeper.MustMarshalClientState(clientState),
		}, {
			Key:   host.FullConsensusStateKey(chainName, height),
			Value: teleport.XIBCKeeper.ClientKeeper.MustMarshalConsensusState(consState),
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
		{"ConsensusState", fmt.Sprintf("ConsensusState A: %v\nConsensusState B: %v", consState, consState)},
		{"other", ""},
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			res, found := simulation.NewDecodeStore(teleport.XIBCKeeper.ClientKeeper, kvPairs.Pairs[i], kvPairs.Pairs[i])
			if i == len(tests)-1 {
				require.False(t, found, string(kvPairs.Pairs[i].Key))
				require.Empty(t, res, string(kvPairs.Pairs[i].Key))
			} else {
				require.True(t, found, string(kvPairs.Pairs[i].Key))
				require.Equal(t, tt.expectedLog, res, string(kvPairs.Pairs[i].Key))
			}
		})
	}
}
