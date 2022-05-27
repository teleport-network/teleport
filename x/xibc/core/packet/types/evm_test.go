package types_test

import (
	"testing"

	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"

	"github.com/stretchr/testify/require"
)

func TestPacketAbi(t *testing.T) {
	packet := packettypes.NewPacket(
		"srcChain",
		"destChain",
		1,
		[]byte("mock Transfer Data"),
		[]byte("mock Call Data"),
		"",
		0,
	)
	packData, err := packet.AbiPack()
	require.NoError(t, err)
	require.NotNil(t, packData)
	var p packettypes.Packet
	err = p.DecodeAbiBytes(packData)
	require.NoError(t, err)
}

func TestAckAbi(t *testing.T) {
	ack := packettypes.NewResultAcknowledgement(
		0,
		[]byte("nodata"),
		"",
		"address",
	)
	ackData, err := ack.AbiPack()
	require.NoError(t, err)
	require.NotNil(t, ackData)
	var p packettypes.Acknowledgement
	err = p.DecodeAbiBytes(ackData)
	require.NoError(t, err)
}
