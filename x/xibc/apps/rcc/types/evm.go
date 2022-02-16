package types

import (
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
	TupleRCCPacketData abi.Type
)

func init() {
	initRCCPacketData()
}

func initRCCPacketData() {
	tupleRCCPacketData, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "src_chain", Type: "string"},
			{Name: "dest_chain", Type: "string"},
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
