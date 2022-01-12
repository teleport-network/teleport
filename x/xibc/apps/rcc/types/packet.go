package types

// NewRCCPacketData contructs a new RCCPacketData instance
func NewRCCPacketData(
	srcChain string,
	destChain string,
	sender string,
	contractAddress string,
	data []byte,
) RCCPacketData {
	return RCCPacketData{
		SrcChain:        srcChain,
		DestChain:       destChain,
		Sender:          sender,
		ContractAddress: contractAddress,
		Data:            data,
	}
}

// ValidateBasic is used for validating the ft transfer.
func (data RCCPacketData) ValidateBasic() error {
	// TODO
	return nil
}

// GetBytes is a helper for serialising
func (data RCCPacketData) GetBytes() []byte {
	return ModuleCdc.MustMarshal(&data)
}
