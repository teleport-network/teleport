package host_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	host "github.com/teleport-network/teleport/x/xibc/core/host"
)

func TestParseIdentifier(t *testing.T) {
	testCases := []struct {
		name       string
		identifier string
		prefix     string
		expSeq     uint64
		expPass    bool
	}{
		{"valid 0", "chain-0", "chain-", 0, true},
		{"valid 1", "chain-1", "chain-", 1, true},
		// one above uint64 max
		{"invalid uint64", "chain-18446744073709551616", "chain-", 0, false},
		// uint64 == 20 characters
		{"invalid large sequence", "chain-2345682193567182931243", "chain-", 0, false},
		{"capital prefix", "Chain-0", "chain-", 0, false},
		{"double prefix", "chain-chain-0", "chain-", 0, false},
		{"doesn't have prefix", "chain-0", "prefix", 0, false},
		{"missing dash", "chain0", "chain-", 0, false},
		{"blank id", "               ", "chain-", 0, false},
		{"empty id", "", "chain-", 0, false},
		{"negative sequence", "chain--1", "chain-", 0, false},
	}

	for _, tc := range testCases {
		seq, err := host.ParseIdentifier(tc.identifier, tc.prefix)
		require.Equal(t, tc.expSeq, seq)

		if tc.expPass {
			require.NoError(t, err, tc.name)
		} else {
			require.Error(t, err, tc.name)
		}
	}
}
