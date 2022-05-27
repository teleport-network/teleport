package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

const PacketSendEvent = "PacketSent"

type InToken struct {
	OriToken string
	Amount   *big.Int
	Scale    uint8
	Bound    bool
}

type CrossChainData struct {
	// path data
	DestChain string
	// transfer token data
	TokenAddress common.Address // zero address if base token
	Receiver     string
	Amount       *big.Int
	// contract call data
	ContractAddress string
	CallData        []byte
	// callback data
	CallbackAddress common.Address
	// fee option
	FeeOption uint64
}

type Fee struct {
	TokenAddress common.Address // zero address if base token
	Amount       *big.Int
}

type Ack struct {
	Code      uint64
	Result    []byte
	Message   string
	Relayer   string
	FeeOption uint64
}

var (
	TupleRecvPacketResultData abi.Type
	TuplePacketData           abi.Type
	TupleAckData              abi.Type
	TuplePacketSendData       abi.Type
	TupleTransferData         abi.Type
	TupleCallData             abi.Type
)

func init() {
	initTupleAckData()
	initRecvPacketResultData()
	initPacketData()
	initTuplePacketSendData()
	initTupleTransferData()
	initTupleCallData()
}

func initRecvPacketResultData() {
	tupleRecvPacketResultData, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "code", Type: "uint64"},
			{Name: "result", Type: "bytes"},
			{Name: "message", Type: "string"},
		},
	)
	if err != nil {
		panic(err)
	}
	if tupleRecvPacketResultData.T != abi.TupleTy {
		panic("New TupleAckData type err")
	}
	TupleRecvPacketResultData = tupleRecvPacketResultData
}

func initPacketData() {
	tuplePacketData, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "source_chain", Type: "string"},
			{Name: "destination_chain", Type: "string"},
			{Name: "sequence", Type: "uint64"},
			{Name: "sender", Type: "string"},
			{Name: "transfer_data", Type: "bytes"},
			{Name: "call_data", Type: "bytes"},
			{Name: "callback_address", Type: "string"},
			{Name: "fee_option", Type: "uint64"},
		},
	)
	if err != nil {
		panic(err)
	}
	if tuplePacketData.T != abi.TupleTy {
		panic("New TupleAckData type err")
	}
	TuplePacketData = tuplePacketData
}

func initTupleAckData() {
	tupleAckData, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "code", Type: "uint64"},
			{Name: "result", Type: "bytes"},
			{Name: "message", Type: "string"},
			{Name: "relayer", Type: "string"},
			{Name: "fee_option", Type: "uint64"},
		},
	)
	if err != nil {
		panic(err)
	}
	if tupleAckData.T != abi.TupleTy {
		panic("New TupleAckData type err")
	}
	TupleAckData = tupleAckData
}

func initTuplePacketSendData() {
	tuplePacketSendData, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "packet", Type: "bytes"},
		},
	)
	if err != nil {
		panic(err)
	}
	if tuplePacketSendData.T != abi.TupleTy {
		panic("New TupleAckData type err")
	}
	TuplePacketSendData = tuplePacketSendData
}

func initTupleTransferData() {
	tupleTransferData, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "receiver", Type: "string"},
			{Name: "amount", Type: "bytes"},
			{Name: "token", Type: "string"},
			{Name: "ori_token", Type: "string"},
		},
	)
	if err != nil {
		panic(err)
	}
	if tupleTransferData.T != abi.TupleTy {
		panic("New TupleAckData type err")
	}
	TupleTransferData = tupleTransferData
}

func initTupleCallData() {
	tupleCallData, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "contract_address", Type: "string"},
			{Name: "call_data", Type: "bytes"},
		},
	)
	if err != nil {
		panic(err)
	}
	if tupleCallData.T != abi.TupleTy {
		panic("New TupleAckData type err")
	}
	TupleCallData = tupleCallData
}
