package types

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

type Amount struct {
	Value *big.Int
}

type CallRCCData struct {
	ContractAddress string
	Data            []byte
	DestChain       string
	RelayChain      string
}

var (
	TupleRCCPacketData          abi.Type
	TupleRCCEventSendPacketData abi.Type
)

func init() {
	initRCCPacketData()
	initRCCEventSendPacketData()
}

func initRCCPacketData() {
	tupleRCCPacketData, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "src_chain", Type: "string"},
			{Name: "dest_chain", Type: "string"},
			{Name: "sequence", Type: "uint64"},
			{Name: "sender", Type: "string"},
			{Name: "contract_address", Type: "string"},
			{Name: "data", Type: "bytes"},
		},
	)
	if err != nil {
		panic(err)
	}
	if tupleRCCPacketData.T != abi.TupleTy {
		panic("New TupleRCCPacketData type err")
	}
	TupleRCCPacketData = tupleRCCPacketData
}

func initRCCEventSendPacketData() {
	tupleRCCEventSendPacketData, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "src_chain", Type: "string"},
			{Name: "dest_chain", Type: "string"},
			{Name: "relay_chain", Type: "string"},
			{Name: "sequence", Type: "uint64"},
			{Name: "sender", Type: "string"},
			{Name: "contract_address", Type: "string"},
			{Name: "data", Type: "bytes"},
		},
	)
	if err != nil {
		panic(err)
	}
	if tupleRCCEventSendPacketData.T != abi.TupleTy {
		panic("New TupleRCCPacketData type err")
	}
	TupleRCCEventSendPacketData = tupleRCCEventSendPacketData
}

type RCCEventSendPacketData struct {
	SrcChain        string `json:"srcChain"`
	DestChain       string `json:"destChain"`
	RelayChain      string `json:"relayChain"`
	Sequence        uint64 `json:"sequence"`
	Sender          string `json:"sender"`
	ContractAddress string `json:"contractAddress"`
	Data            []byte `json:"data"`
}

func (data *RCCEventSendPacketData) DecodeBytes(bz []byte) error {
	dataBz, err := abi.Arguments{{Type: TupleRCCEventSendPacketData}}.Unpack(bz)
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

func (data *RCCEventSendPacketData) DecodeInterface(bz interface{}) error {
	bzTmp, err := json.Marshal(bz)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(bzTmp, &data); err != nil {
		return err
	}
	return nil
}
