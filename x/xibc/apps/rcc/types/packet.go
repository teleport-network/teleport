package types

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/abi"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewRCCPacketData contructs a new RCCPacketData instance
func NewRCCPacketData(
	srcChain string,
	destChain string,
	sequence uint64,
	sender string,
	contractAddress string,
	data []byte,
) RCCPacketData {
	return RCCPacketData{
		SrcChain:        srcChain,
		DestChain:       destChain,
		Sequence:        sequence,
		Sender:          sender,
		ContractAddress: contractAddress,
		Data:            data,
	}
}

// ValidateBasic is used for validating the ft transfer.
func (data RCCPacketData) ValidateBasic() error {
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

	if len(data.ContractAddress) == 0 {
		return sdkerrors.Wrap(ErrInvalidAddress, "address is invalid")
	}

	return nil
}

// GetBytes is a helper for serialising
func (data RCCPacketData) GetBytes() ([]byte, error) {
	return abi.Arguments{{Type: TupleRCCPacketData}}.Pack(data)
}

func (data *RCCPacketData) DecodeBytes(bz []byte) error {
	dataBz, err := abi.Arguments{{Type: TupleRCCPacketData}}.Unpack(bz)
	if err != nil {
		return err
	}
	bzTmp, err := json.Marshal(dataBz[0])
	if err != nil {
		return err
	}
	return json.Unmarshal(bzTmp, &data)
}
