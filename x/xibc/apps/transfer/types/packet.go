package types

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// NewFungibleTokenPacketData contructs a new FungibleTokenPacketData instance
func NewFungibleTokenPacketData(
	srcChain string,
	destChain string,
	sequence uint64,
	sender string,
	receiver string,
	amount []byte,
	token string,
	oriToken string,
) FungibleTokenPacketData {
	return FungibleTokenPacketData{
		SrcChain:  srcChain,
		DestChain: destChain,
		Sequence:  sequence,
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
func (data FungibleTokenPacketData) GetBytes() ([]byte, error) {
	dataBz, err := abi.Arguments{{Type: TupleFTPacketData}}.Pack(
		data,
	)
	if err != nil {
		return nil, err
	}
	return dataBz, nil
}

func (data *FungibleTokenPacketData) DecodeBytes(bz []byte) error {
	dataBz, err := abi.Arguments{{Type: TupleFTPacketData}}.Unpack(bz)
	if err != nil {
		return err
	}
	bzTmp, err := json.Marshal(dataBz[0])
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bzTmp, &data); err != nil {
		return err
	}
	return nil
}
