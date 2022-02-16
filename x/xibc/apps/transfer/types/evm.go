package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type Amount struct {
	Value *big.Int
}

type ERC20TransferData struct {
	TokenAddress common.Address
	Receiver     string
	Amount       *big.Int
	DestChain    string
	RelayChain   string
}

type BaseTransferData struct {
	Receiver   string
	DestChain  string
	RelayChain string
}

var (
	TupleFTPacketData abi.Type
)

func init() {
	initTupleFTPacketData()
}

func initTupleFTPacketData() {
	tupleFTPacketData, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "src_chain", Type: "string"},
			{Name: "dest_chain", Type: "string"},
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
