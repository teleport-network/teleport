package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	"github.com/teleport-network/teleport/x/xibc/core/packet/types"
)

func TestCommitPacket(t *testing.T) {
	packet := types.NewPacket(srcChain, dstChain, 1, "sender", mockTransferData, mockCallData, "", 0)

	registry := codectypes.NewInterfaceRegistry()
	clienttypes.RegisterInterfaces(registry)
	types.RegisterInterfaces(registry)

	commitment, err := types.CommitPacket(packet)
	require.NoError(t, err)
	require.NotNil(t, commitment)
}

func TestPacketValidateBasic(t *testing.T) {
	testCases := []struct {
		packet  *types.Packet
		expPass bool
		errMsg  string
	}{
		{types.NewPacket(srcChain, dstChain, 1, "sender", mockTransferData, mockCallData, "", 0), true, ""},
		{types.NewPacket(srcChain, dstChain, 0, "sender", mockTransferData, mockCallData, "", 0), false, "invalid sequence"},
		{types.NewPacket(srcChain, dstChain, 1, "sender", []byte(""), mockCallData, "", 0), true, ""},
	}

	for i, tc := range testCases {
		if tc.expPass {
			require.NoError(t, tc.packet.ValidateBasic(), "Case %d failed: %s", i, tc.errMsg)
		} else {
			require.Error(t, tc.packet.ValidateBasic(), "Invalid Case %d passed: %s", i, tc.errMsg)
		}
	}
}
