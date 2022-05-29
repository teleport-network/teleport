package types_test

import (
	"testing"

	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"

	"github.com/stretchr/testify/require"
)

func TestPacketAbi(t *testing.T) {
	crossChainPacket := packettypes.NewPacket(
		"srcChain",
		"destChain",
		1,
		"sender",
		[]byte("mock Transfer Data"),
		[]byte("mock Call Data"),
		"callback_address",
		0,
	)
	packetAbi, err := crossChainPacket.ABIPack()
	require.NoError(t, err)
	require.NotNil(t, packetAbi)
	var p packettypes.Packet
	err = p.ABIDecode(packetAbi)
	require.NoError(t, err)
	require.Equal(t, p.SourceChain, crossChainPacket.SourceChain)
	require.Equal(t, p.DestinationChain, crossChainPacket.DestinationChain)
	require.Equal(t, p.Sequence, crossChainPacket.Sequence)
	require.Equal(t, p.Sender, crossChainPacket.Sender)
	require.Equal(t, p.TransferData, crossChainPacket.TransferData)
	require.Equal(t, p.CallData, crossChainPacket.CallData)
	require.Equal(t, p.CallbackAddress, crossChainPacket.CallbackAddress)
	require.Equal(t, p.FeeOption, crossChainPacket.FeeOption)
}

func TestAckAbi(t *testing.T) {
	ack := packettypes.NewAcknowledgement(
		0,
		[]byte("nodata"),
		"",
		"address",
		0,
	)
	ackData, err := ack.ABIPack()
	require.NoError(t, err)
	require.NotNil(t, ackData)
	var p packettypes.Acknowledgement
	err = p.ABIDecode(ackData)
	require.NotNil(t, p)
	require.NoError(t, err)
}

func TestTransferDataAbi(t *testing.T) {
	transferData := packettypes.TransferData{
		Receiver: "receiver",
		Amount:   []byte("Amount"),
		Token:    "token",
		OriToken: "ori_token",
	}
	data, err := transferData.ABIPack()
	require.NoError(t, err)
	require.NotNil(t, data)
	var p packettypes.TransferData
	err = p.ABIDecode(data)
	require.NoError(t, err)
	require.Equal(t, p.Receiver, transferData.Receiver)
	require.Equal(t, p.Amount, transferData.Amount)
	require.Equal(t, p.Token, transferData.Token)
	require.Equal(t, p.OriToken, transferData.OriToken)
}

func TestCallDataAbi(t *testing.T) {
	callData := packettypes.CallData{
		ContractAddress: "contract_address",
		CallData:        []byte(""),
	}
	data, err := callData.ABIPack()
	require.NoError(t, err)
	require.NotNil(t, data)
	var p packettypes.CallData
	err = p.ABIDecode(data)
	require.NoError(t, err)
	require.Equal(t, p.ContractAddress, callData.ContractAddress)
	require.Equal(t, p.CallData, callData.CallData)
}
