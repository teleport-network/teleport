package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPacketAbi(t *testing.T) {
	packet := NewPacket(
		"srcChain",
		"destChain",
		"relayChain",
		1,
		[]byte("mock Transfer Data"),
		[]byte("mock Call Data"),
		"",
		0,
	)
	packData, err := packet.AbiPack()
	require.NoError(t, err)
	require.NotNil(t, packData)
	var p Packet
	err = p.DecodeAbiBytes(packData)
	require.NoError(t, err)
}

func TestAckAbi(t *testing.T) {
	ack := NewResultAcknowledgement(
		0,
		[]byte("nodata"),
		"",
		"address",
	)
	ackData, err := ack.AbiPack()
	require.NoError(t, err)
	require.NotNil(t, ackData)
	var p Acknowledgement
	err = p.DecodeAbiBytes(ackData)
	require.NoError(t, err)
}
