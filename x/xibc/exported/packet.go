package exported

// PacketI defines the standard interface for XIBC packets
type PacketI interface {
	GetSequence() uint64
	GetSrcChain() string
	GetDestChain() string
	GetSender() string
	GetTransferData() []byte
	GetCallData() []byte
	GetCallbackAddress() string
	GetFeeOption() uint64
	ABIPack() ([]byte, error)
	ABIDecode(bz []byte) error
	ValidateBasic() error
}
