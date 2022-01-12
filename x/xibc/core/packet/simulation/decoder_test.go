package simulation_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/teleport-network/teleport/app"
	host "github.com/teleport-network/teleport/x/xibc/core/host"
	"github.com/teleport-network/teleport/x/xibc/core/packet/simulation"
)

func TestDecodeStore(t *testing.T) {
	teleport := app.Setup(false, nil)
	cdc := teleport.AppCodec()

	chainID := "chain"
	portID := "port"

	bz := []byte{0x1, 0x2, 0x3}

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{{
			Key:   host.NextSequenceSendKey(portID, chainID),
			Value: sdk.Uint64ToBigEndian(1),
		}, {
			Key:   host.PacketCommitmentKey(portID, chainID, 1),
			Value: bz,
		}, {
			Key:   host.PacketAcknowledgementKey(portID, chainID, 1),
			Value: bz,
		}, {
			Key:   []byte{0x99},
			Value: []byte{0x99},
		}},
	}
	tests := []struct {
		name        string
		expectedLog string
	}{
		{"NextSeqSend", "NextSeqSend A: 1\nNextSeqSend B: 1"},
		{"CommitmentHash", fmt.Sprintf("CommitmentHash A: %X\nCommitmentHash B: %X", bz, bz)},
		{"AckHash", fmt.Sprintf("AckHash A: %X\nAckHash B: %X", bz, bz)},
		{"other", ""},
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			res, found := simulation.NewDecodeStore(cdc, kvPairs.Pairs[i], kvPairs.Pairs[i])
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
