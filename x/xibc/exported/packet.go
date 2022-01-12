package exported

// PacketI defines the standard interface for XIBC packets
type PacketI interface {
	GetSequence() uint64
	GetSourceChain() string
	GetDestChain() string
	GetRelayChain() string
	GetDataList() [][]byte
	GetPorts() []string
	ValidateBasic() error
}
