package types

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

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
func (data RCCPacketData) GetBytes() ([]byte, error) {
	pack, err := abi.Arguments{{Type: TupleRCCPacketData}}.Pack(data)
	if err != nil {
		return nil, err
	}
	return pack, nil
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
	if err = json.Unmarshal(bzTmp, &data); err != nil {
		return err
	}
	return nil
}
