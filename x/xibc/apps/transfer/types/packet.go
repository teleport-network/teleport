package types

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/abi"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
	if len(data.SrcChain) == 0 {
		return sdkerrors.Wrap(ErrInvalidSrcChain, "srcChain is empty")
	}

	if len(data.DestChain) == 0 {
		return sdkerrors.Wrap(ErrInvalidDestChain, "destChain is empty")
	}

	if data.SrcChain == data.DestChain {
		return sdkerrors.Wrap(ErrScChainEqualToDestChain, "srcChain equals to destChain")
	}

	if data.Sequence == 0 {
		return sdkerrors.Wrap(ErrInvalidSequence, "packet sequence cannot be 0")
	}

	if len(data.Sender) == 0 {
		return sdkerrors.Wrap(ErrInvalidSender, "sender is empty")
	}

	if len(data.Token) == 0 {
		return sdkerrors.Wrap(ErrInvalidAddress, "address is invalid")
	}

	return nil
}

// GetBytes is a helper for serialising
func (data FungibleTokenPacketData) GetBytes() ([]byte, error) {
	return abi.Arguments{{Type: TupleFTPacketData}}.Pack(data)
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
	return json.Unmarshal(bzTmp, &data)
}
