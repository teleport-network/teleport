package types

import (
	"crypto/sha256"
	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/teleport-network/teleport/x/xibc/exported"
)

// CommitPacket returns the packet commitment bytes. The commitment consists of:
// sha256_hash(timeout_timestamp + timeout_height.RevisionNumber + timeout_height.RevisionHeight + sha256_hash(port + data))
// from a given packet. This results in a fixed length preimage.
// NOTE: sdk.Uint64ToBigEndian sets the uint64 to a slice of length 8.
func CommitPacket(packet exported.PacketI) []byte {
	var dataSum []byte
	for i, data := range packet.GetDataList() {
		dataHash := sha256.Sum256(
			append([]byte(packet.GetPorts()[i]), data...),
		)
		dataSum = append(dataSum, dataHash[:]...)
	}
	hash := sha256.Sum256(dataSum)
	return hash[:]
}

// CommitAcknowledgement returns the hash of commitment bytes
func CommitAcknowledgement(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

var _ exported.PacketI = (*Packet)(nil)

var ModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())

// NewPacket creates a new Packet instance. It panics if the provided packet data interface is not registered.
func NewPacket(
	sequence uint64,
	sourceChain string,
	destinationChain string,
	relayChain string,
	ports []string,
	dataList [][]byte,
) Packet {
	return Packet{
		Sequence:         sequence,
		SourceChain:      sourceChain,
		DestinationChain: destinationChain,
		RelayChain:       relayChain,
		Ports:            ports,
		DataList:         dataList,
	}
}

// GetSequence implements PacketI interface
func (p Packet) GetSequence() uint64 { return p.Sequence }

// GetSourceChain implements PacketI interface
func (p Packet) GetSourceChain() string { return p.SourceChain }

// GetDestinationChain implements PacketI interface
func (p Packet) GetDestChain() string { return p.DestinationChain }

// GetRelayChain implements PacketI interface
func (p Packet) GetRelayChain() string { return p.RelayChain }

// GetPorts implements PacketI interface
func (p Packet) GetPorts() []string { return p.Ports }

// GetDataList implements PacketI interface
func (p Packet) GetDataList() [][]byte { return p.DataList }

// ValidateBasic implements PacketI interface
func (p Packet) ValidateBasic() error {
	if p.Sequence == 0 {
		return sdkerrors.Wrap(ErrInvalidPacket, "packet sequence cannot be 0")
	}
	if len(p.Ports) != len(p.DataList) {
		return sdkerrors.Wrap(ErrInvalidPacket, "the number of ports must be as many as the data")
	}
	if len(p.DataList) == 0 {
		return sdkerrors.Wrap(ErrInvalidPacket, "empty packet data list")
	}
	for _, data := range p.DataList {
		if len(data) == 0 {
			return sdkerrors.Wrap(ErrInvalidPacket, "packet data bytes cannot be empty")
		}
	}
	return nil
}

// NewResultAcknowledgement returns a new instance of Acknowledgement using an Acknowledgement_Result type in the Response field.
func NewResultAcknowledgement(results [][]byte) Acknowledgement {
	return Acknowledgement{Results: results}
}

// NewErrorAcknowledgement returns a new instance of Acknowledgement using an Acknowledgement_Error type in the Response field.
func NewErrorAcknowledgement(message string) Acknowledgement {
	return Acknowledgement{
		Message: message,
	}
}

// GetBytes is a helper for serialising acknowledgements
func (ack Acknowledgement) GetBytes() ([]byte, error) {
	pack, err := abi.Arguments{{Type: TupleAckData}}.Pack(ack)
	if err != nil {
		return nil, err
	}
	return pack, nil
}

func (ack *Acknowledgement) DecodeBytes(bz []byte) error {
	dataBz, err := abi.Arguments{{Type: TupleAckData}}.Unpack(bz)
	if err != nil {
		return err
	}
	bzTmp, err := json.Marshal(dataBz[0])
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bzTmp, &ack); err != nil {
		return err
	}
	return nil
}

// Result is the execution result of packet data
type Result struct {
	// the execution result
	Result []byte
	// error message
	Message string
}
