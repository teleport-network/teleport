package types

import "github.com/ethereum/go-ethereum/accounts/abi"

var (
	TupleRecvPacketResultData abi.Type
	TuplePacketData           abi.Type
	TupleAckData              abi.Type
)

func init() {
	initTupleAckData()
	initRecvPacketResultData()
	initPacketData()
}

func initRecvPacketResultData() {
	tupleRecvPacketResultData, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "code", Type: "uint64"},
			{Name: "results", Type: "bytes"},
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
			{Name: "destination_port", Type: "string"},
			{Name: "relay_chain", Type: "string"},
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
