package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	"github.com/teleport-network/teleport/x/xibc/core/packet/types"
)

func TestCommitPacket(t *testing.T) {
	packet := types.NewPacket(1, sourceChain, destChain, relayChain, []string{port}, [][]byte{validPacketData})

	registry := codectypes.NewInterfaceRegistry()
	clienttypes.RegisterInterfaces(registry)
	types.RegisterInterfaces(registry)

	commitment := types.CommitPacket(&packet)
	require.NotNil(t, commitment)
}

func TestPacketValidateBasic(t *testing.T) {
	testCases := []struct {
		packet  types.Packet
		expPass bool
		errMsg  string
	}{
		{types.NewPacket(1, sourceChain, destChain, relayChain, []string{port}, [][]byte{validPacketData}), true, ""},
		{types.NewPacket(0, sourceChain, destChain, relayChain, []string{port}, [][]byte{validPacketData}), false, "invalid sequence"},
		{types.NewPacket(1, sourceChain, destChain, relayChain, []string{port}, [][]byte{unknownPacketData}), true, ""},
	}

	for i, tc := range testCases {
		if tc.expPass {
			require.NoError(t, tc.packet.ValidateBasic(), "Case %d failed: %s", i, tc.errMsg)
		} else {
			require.Error(t, tc.packet.ValidateBasic(), "Invalid Case %d passed: %s", i, tc.errMsg)
		}
	}
}
