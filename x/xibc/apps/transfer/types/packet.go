package types

// NewFungibleTokenPacketData contructs a new FungibleTokenPacketData instance
func NewFungibleTokenPacketData(
	srcChain string,
	destChain string,
	sender string,
	receiver string,
	amount []byte,
	token string,
	oriToken string,
) FungibleTokenPacketData {
	return FungibleTokenPacketData{
		SrcChain:  srcChain,
		DestChain: destChain,
		Sender:    sender,
		Receiver:  receiver,
		Amount:    amount,
		Token:     token,
		OriToken:  oriToken,
	}
}

// ValidateBasic is used for validating the ft transfer.
func (data FungibleTokenPacketData) ValidateBasic() error {
	// TODO
	return nil
}

// GetBytes is a helper for serialising
func (data FungibleTokenPacketData) GetBytes() []byte {
	return ModuleCdc.MustMarshal(&data)
}
