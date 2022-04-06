package types

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type Amount struct {
	Value *big.Int
}

type TransferData struct {
	TokenAddress common.Address
	Receiver     string
	Amount       *big.Int
	DestChain    string
	RelayChain   string
}

type Fee struct {
	TokenAddress common.Address
	Amount       *big.Int
}

var (
	TupleFTPacketData                abi.Type
	TupleTransferEventSendPacketData abi.Type
)

func init() {
	initTupleFTPacketData()
	initTupleTransferEventSendPacketData()
}

func initTupleFTPacketData() {
	tupleFTPacketData, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "src_chain", Type: "string"},
			{Name: "dest_chain", Type: "string"},
			{Name: "sequence", Type: "uint64"},
			{Name: "sender", Type: "string"},
			{Name: "receiver", Type: "string"},
			{Name: "amount", Type: "bytes"},
			{Name: "token", Type: "string"},
			{Name: "ori_token", Type: "string"},
		},
	)
	if err != nil {
		panic(err)
	}
	if tupleFTPacketData.T != abi.TupleTy {
		panic("New TupleERC20TransferData type err")
	}
	TupleFTPacketData = tupleFTPacketData
}

func initTupleTransferEventSendPacketData() {
	tupleTransferEventSendPacketData, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "src_chain", Type: "string"},
			{Name: "dest_chain", Type: "string"},
			{Name: "relay_chain", Type: "string"},
			{Name: "sequence", Type: "uint64"},
			{Name: "sender", Type: "string"},
			{Name: "receiver", Type: "string"},
			{Name: "amount", Type: "bytes"},
			{Name: "token", Type: "string"},
			{Name: "ori_token", Type: "string"},
		},
	)
	if err != nil {
		panic(err)
	}
	if tupleTransferEventSendPacketData.T != abi.TupleTy {
		panic("New TupleERC20TransferData type err")
	}
	TupleTransferEventSendPacketData = tupleTransferEventSendPacketData
}

type TransferEventSendPacketData struct {
	SrcChain   string   `json:"srcChain"`
	DestChain  string   `json:"destChain"`
	RelayChain string   `json:"relayChain"`
	Sequence   uint64   `json:"sequence"`
	Sender     string   `json:"sender"`
	Receiver   string   `json:"receiver"`
	Amount     *big.Int `json:"amount"`
	Token      string   `json:"token"`
	OriToken   string   `json:"oriToken"`
}

func (data *TransferEventSendPacketData) DecodeBytes(bz []byte) error {
	dataBz, err := abi.Arguments{{Type: TupleTransferEventSendPacketData}}.Unpack(bz)
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

func (data *TransferEventSendPacketData) DecodeInterface(bz interface{}) error {
	bzTmp, err := json.Marshal(bz)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(bzTmp, &data); err != nil {
		return err
	}
	return nil
}
