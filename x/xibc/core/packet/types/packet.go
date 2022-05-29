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
func CommitPacket(packet exported.PacketI) ([]byte, error) {
	packeBytes, err := packet.ABIPack()
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(packeBytes)
	return hash[:], nil
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
	sourceChain string,
	destinationChain string,
	sequence uint64,
	sender string,
	transferData []byte,
	callData []byte,
	callbackAddress string,
	feeOption uint64,
) *Packet {
	return &Packet{
		SourceChain:      sourceChain,
		DestinationChain: destinationChain,
		Sequence:         sequence,
		Sender:           sender,
		TransferData:     transferData,
		CallData:         callData,
		CallbackAddress:  callbackAddress,
		FeeOption:        feeOption,
	}
}

// GetSequence implements PacketI interface
func (p Packet) GetSequence() uint64 { return p.Sequence }

// GetSrcChain implements PacketI interface
func (p Packet) GetSrcChain() string { return p.SourceChain }

// GetDestChain implements PacketI interface
func (p Packet) GetDestChain() string { return p.DestinationChain }

// GetRelayChain implements PacketI interface
func (p Packet) GetSender() string { return p.Sender }

// GetTransferData implements PacketI interface
func (p Packet) GetTransferData() []byte { return p.TransferData }

// GetCallData implements PacketI interface
func (p Packet) GetCallData() []byte { return p.CallData }

// GetCallbackAddress implements PacketI interface
func (p Packet) GetCallbackAddress() string { return p.CallbackAddress }

// GetFeeOption implements PacketI interface
func (p Packet) GetFeeOption() uint64 { return p.FeeOption }

// ABIPack implements PacketI interface
func (p Packet) ABIPack() ([]byte, error) {
	pack, err := abi.Arguments{{Type: TuplePacketData}}.Pack(p)
	if err != nil {
		return nil, err
	}
	return pack, nil
}

// ABIDecode implements PacketI interface
func (p *Packet) ABIDecode(bz []byte) error {
	dataBz, err := abi.Arguments{{Type: TuplePacketData}}.Unpack(bz)
	if err != nil {
		return err
	}
	bzTmp, err := json.Marshal(dataBz[0])
	if err != nil {
		return err
	}
	return json.Unmarshal(bzTmp, &p)
}

// ValidateBasic implements PacketI interface
func (p Packet) ValidateBasic() error {
	if len(p.SourceChain) == 0 {
		return sdkerrors.Wrap(ErrInvalidSrcChain, "srcChain is empty")
	}

	if len(p.DestinationChain) == 0 {
		return sdkerrors.Wrap(ErrInvalidDestChain, "destChain is empty")
	}

	if p.SourceChain == p.DestinationChain {
		return sdkerrors.Wrap(ErrScChainEqualToDestChain, "srcChain equals to destChain")
	}

	if p.Sequence == 0 {
		return sdkerrors.Wrap(ErrInvalidPacket, "packet sequence cannot be 0")
	}
	// todo validate packet data
	if len(p.CallData) == 0 && len(p.TransferData) == 0 {
		return sdkerrors.Wrap(ErrInvalidPacket, "packet has no data")
	}
	return nil
}

type WPacket struct {
	// packet base data
	SrcChain  string
	DestChain string
	Sequence  uint64
	Sender    string
	// transfer data. keep empty if not used.
	TransferData []byte
	// call data. keep empty if not used
	CallData []byte
	// callback data
	CallbackAddress string
	// fee option
	FeeOption uint64
}

// ABIDecode implements PacketI interface
func (p *Packet) ToWPacket() WPacket {
	return WPacket{
		p.SourceChain,
		p.DestinationChain,
		p.Sequence,
		p.Sender,
		p.TransferData,
		p.CallData,
		p.CallbackAddress,
		p.FeeOption,
	}
}

// NewResultAcknowledgement returns a new instance of Acknowledgement using an Acknowledgement_Result type in the Response field.
func NewAcknowledgement(code uint64, results []byte, message, relayer string, feeOption uint64) Acknowledgement {
	return Acknowledgement{
		Code:      code,
		Result:    results,
		Message:   message,
		Relayer:   relayer,
		FeeOption: feeOption,
	}
}

// ABIPack is a helper for serialising acknowledgements
func (ack Acknowledgement) ABIPack() ([]byte, error) {
	pack, err := abi.Arguments{{Type: TupleAckData}}.Pack(ack)
	if err != nil {
		return nil, err
	}
	return pack, nil
}

func (ack *Acknowledgement) ABIDecode(bz []byte) error {
	dataBz, err := abi.Arguments{{Type: TupleAckData}}.Unpack(bz)
	if err != nil {
		return err
	}
	bzTmp, err := json.Marshal(dataBz[0])
	if err != nil {
		return err
	}
	return json.Unmarshal(bzTmp, &ack)
}

// Result is the execution result of packet data
type Result struct {
	Code uint64
	// the execution result
	Result []byte
	// error message
	Message string
}

// ABIPack is a helper for serialising Result
func (result Result) ABIPack() ([]byte, error) {
	pack, err := abi.Arguments{{Type: TupleRecvPacketResultData}}.Pack(result)
	if err != nil {
		return nil, err
	}
	return pack, nil
}

func (result *Result) DecodeInterface(bz interface{}) error {
	bzTmp, err := json.Marshal(bz)
	if err != nil {
		return err
	}
	return json.Unmarshal(bzTmp, &result)
}

func (result *Result) ABIDecode(bz []byte) error {
	dataBz, err := abi.Arguments{{Type: TupleRecvPacketResultData}}.Unpack(bz)
	if err != nil {
		return err
	}
	bzTmp, err := json.Marshal(dataBz[0])
	if err != nil {
		return err
	}
	return json.Unmarshal(bzTmp, &result)
}

func (e *EventSendPacket) ABIDecode(bz []byte) error {
	dataBz, err := abi.Arguments{{Type: TuplePacketSendData}}.Unpack(bz)
	if err != nil {
		return err
	}
	bzTmp, err := json.Marshal(dataBz[0])
	if err != nil {
		return err
	}
	return json.Unmarshal(bzTmp, &e)
}

func (e *EventSendPacket) DecodeInterface(bz interface{}) error {
	bzTmp, err := json.Marshal(bz)
	if err != nil {
		return err
	}
	return json.Unmarshal(bzTmp, &e)
}

// ABIPack is a helper for serialising Result
func (e *EventSendPacket) ABIPack() ([]byte, error) {
	pack, err := abi.Arguments{{Type: TuplePacketSendData}}.Pack(e)
	if err != nil {
		return nil, err
	}
	return pack, nil
}

// ABIPack is a helper for serialising Result
func (e *TransferData) ABIPack() ([]byte, error) {
	pack, err := abi.Arguments{{Type: TupleTransferData}}.Pack(e)
	if err != nil {
		return nil, err
	}
	return pack, nil
}

func (e *TransferData) ABIDecode(bz []byte) error {
	dataBz, err := abi.Arguments{{Type: TupleTransferData}}.Unpack(bz)
	if err != nil {
		return err
	}
	bzTmp, err := json.Marshal(dataBz[0])
	if err != nil {
		return err
	}
	return json.Unmarshal(bzTmp, &e)
}

// ABIPack is a helper for serialising Result
func (e *CallData) ABIPack() ([]byte, error) {
	pack, err := abi.Arguments{{Type: TupleCallData}}.Pack(e)
	if err != nil {
		return nil, err
	}
	return pack, nil
}

func (e *CallData) ABIDecode(bz []byte) error {
	dataBz, err := abi.Arguments{{Type: TupleCallData}}.Unpack(bz)
	if err != nil {
		return err
	}
	bzTmp, err := json.Marshal(dataBz[0])
	if err != nil {
		return err
	}
	return json.Unmarshal(bzTmp, &e)
}
