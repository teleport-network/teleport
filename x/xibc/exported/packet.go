package exported

// PacketI defines the standard interface for XIBC packets
type PacketI interface {
	GetSequence() uint64
	GetSourceChain() string
	GetDestChain() string
	GetRelayChain() string
	GetSender() string
	GetTransferData() []byte
	GetCallData() []byte
	GetCallbackAddress() string
	GetFeeOption() uint64
	AbiPack() ([]byte, error)
	DecodeAbiBytes(bz []byte) error
	ValidateBasic() error
}
