package host

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/teleport-network/teleport/x/xibc/exported"
)

const (
	// ModuleName is the name of the XIBC module
	ModuleName = "xibc"

	// StoreKey is the string store representation
	StoreKey string = ModuleName

	// QuerierRoute is the querier route for the XIBC module
	QuerierRoute string = ModuleName
)

// KVStore key prefixes for XIBC
var (
	KeyClientStorePrefix = []byte("clients")
)

// KVStore key prefixes for XIBC
const (
	KeyClientState            = "clientState"
	KeyConsensusStatePrefix   = "consensusStates"
	KeySequencePrefix         = "sequences"
	KeyNextSeqSendPrefix      = "nextSequenceSend"
	KeyNextSeqRecvPrefix      = "nextSequenceRecv"
	KeyNextSeqAckPrefix       = "nextSequenceAck"
	KeyPacketCommitmentPrefix = "commitments"
	KeyPacketRelayerPrefix    = "relayer"
	KeyPacketAckPrefix        = "acks"
	KeyPacketReceiptPrefix    = "receipts"
)

// FullClientPath returns the full path of a specific client path in the format:
// "clients/{chainName}/{path}" as a string
func FullClientPath(chainName string, path string) string {
	return fmt.Sprintf("%s/%s/%s", KeyClientStorePrefix, chainName, path)
}

// FullClientKey returns the full path of specific client path in the format:
// "clients/{chainName}/{path}" as a byte array
func FullClientKey(chainName string, path []byte) []byte {
	return []byte(FullClientPath(chainName, string(path)))
}

// FullClientStateKey takes a client identifier and returns a Key under which to store a particular client state
func FullClientStateKey(chainName string) []byte {
	return FullClientKey(chainName, []byte(KeyClientState))
}

// ClientStateKey returns a store key under which a particular client state is stored in a client prefixed store
func ClientStateKey() []byte {
	return []byte(KeyClientState)
}

// FullConsensusStateKey returns the store key for the consensus state of a particular client
func FullConsensusStateKey(chainName string, height exported.Height) (key []byte) {
	key = append(key, KeyClientStorePrefix...)
	key = append(key, []byte("/"+chainName+"/")...)
	key = append(key, ConsensusStateKey(height)...)
	return key
}

// ConsensusStatePath returns the suffix store key for the consensus state at a
// particular height stored in a client prefixed store.
func ConsensusStatePath(height exported.Height) string {
	return fmt.Sprintf("%s/%s", KeyConsensusStatePrefix, height)
}

// ConsensusStateKey returns the store key for a the consensus state of a particular
// client stored in a client prefixed store.
func ConsensusStateKey(height exported.Height) (key []byte) {
	reversionNumber := sdk.Uint64ToBigEndian(height.GetRevisionNumber())
	reversionHeight := sdk.Uint64ToBigEndian(height.GetRevisionHeight())

	key = append(key, []byte(KeyConsensusStatePrefix+"/")...)
	key = append(key, reversionNumber...)
	key = append(key, reversionHeight...)
	return key
}

// NextSequenceSendPath defines the next send sequence counter store path
func NextSequenceSendPath(sourceChain, destChain string) string {
	return fmt.Sprintf("%s/%s", KeyNextSeqSendPrefix, packetPath(sourceChain, destChain))
}

// NextSequenceSendKey returns the store key for the send sequence
func NextSequenceSendKey(sourceChain, destChain string) []byte {
	return []byte(NextSequenceSendPath(sourceChain, destChain))
}

// PacketCommitmentPath defines the commitments to packet data fields store path
func PacketCommitmentPath(sourceChain, destinationChain string, sequence uint64) string {
	return fmt.Sprintf("%s/%d", PacketCommitmentPrefixPath(sourceChain, destinationChain), sequence)
}

// PacketCommitmentKey returns the store key of under which a packet commitment is stored
func PacketCommitmentKey(sourceChain, destinationChain string, sequence uint64) []byte {
	return []byte(PacketCommitmentPath(sourceChain, destinationChain, sequence))
}

// PacketCommitmentPrefixPath defines the prefix for commitments to packet data fields store path
func PacketCommitmentPrefixPath(sourceChain, destinationChain string) string {
	return fmt.Sprintf("%s/%s/%s", KeyPacketCommitmentPrefix, packetPath(sourceChain, destinationChain), KeySequencePrefix)
}

// PacketRelayerPath defines the packet data fields store path
func PacketRelayerPath(sourceChain, destinationChain string, sequence uint64) string {
	return fmt.Sprintf("%s/%d", PacketRelayerPrefixPath(sourceChain, destinationChain), sequence)
}

// PacketRelayerKey returns the store key of under which a packet relayer is stored
func PacketRelayerKey(sourceChain, destinationChain string, sequence uint64) []byte {
	return []byte(PacketRelayerPath(sourceChain, destinationChain, sequence))
}

// PacketRelayerPrefixPath defines the prefix for relayer to packet data fields store path
func PacketRelayerPrefixPath(sourceChain, destinationChain string) string {
	return fmt.Sprintf("%s/%s/%s", KeyPacketRelayerPrefix, packetPath(sourceChain, destinationChain), KeySequencePrefix)
}

// PacketAcknowledgementPath defines the packet acknowledgement store path
func PacketAcknowledgementPath(sourceChain, destinationChain string, sequence uint64) string {
	return fmt.Sprintf("%s/%d", PacketAcknowledgementPrefixPath(sourceChain, destinationChain), sequence)
}

// PacketAcknowledgementKey returns the store key of under which a packet
// acknowledgement is stored
func PacketAcknowledgementKey(sourceChain, destinationChain string, sequence uint64) []byte {
	return []byte(PacketAcknowledgementPath(sourceChain, destinationChain, sequence))
}

// PacketAcknowledgementPrefixPath defines the prefix for commitments to packet data fields store path
func PacketAcknowledgementPrefixPath(sourceChain, destinationChain string) string {
	return fmt.Sprintf("%s/%s/%s", KeyPacketAckPrefix, packetPath(sourceChain, destinationChain), KeySequencePrefix)
}

// PacketReceiptPath defines the packet receipt store path
func PacketReceiptPath(sourceChain, destinationChain string, sequence uint64) string {
	return fmt.Sprintf("%s/%d", PacketReceiptPrefixPath(sourceChain, destinationChain), sequence)
}

// PacketReceiptKey returns the store key of under which a packet
// receipt is stored
func PacketReceiptKey(sourceChain, destinationChain string, sequence uint64) []byte {
	return []byte(PacketReceiptPath(sourceChain, destinationChain, sequence))
}

// PacketReceiptKey returns the store key of under which a packet
// receipt is stored
func PacketReceiptPrefixPath(sourceChain, destinationChain string) string {
	return fmt.Sprintf("%s/%s/%s", KeyPacketReceiptPrefix, packetPath(sourceChain, destinationChain), KeySequencePrefix)
}

func packetPath(sourceChain, destinationChain string) string {
	return fmt.Sprintf("%s/%s", sourceChain, destinationChain)
}
