package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

const (
	Transfer   = uint8(0)
	RemoteCall = uint8(1)
)

type Amount struct {
	Value *big.Int
}

type MultiCallData struct {
	DestChain  string   `json:"destChain"`
	RelayChain string   `json:"relayChain"`
	Functions  []uint8  `json:"functions"`
	Data       [][]byte `json:"data"`
}

type TransferData struct {
	TokenAddress common.Address `json:"token_address"`
	Receiver     string         `json:"receiver"`
	Amount       *big.Int       `json:"amount"`
}

type RCCData struct {
	ContractAddress string `json:"contract_address"`
	Data            []byte `json:"data"`
}

type Fee struct {
	TokenAddress common.Address
	Amount       *big.Int
}

var (
	TupleTransferData abi.Type
	TupleRCCData      abi.Type
)

func init() {
	initTransferData()
	initRCCData()
}

func initTransferData() {
	tupleTransferData, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "token_address", Type: "address"},
			{Name: "receiver", Type: "string"},
			{Name: "amount", Type: "uint256"},
		},
	)
	if err != nil {
		panic(err)
	}
	if tupleTransferData.T != abi.TupleTy {
		panic("New TupleERC20TransferData type err")
	}
	TupleTransferData = tupleTransferData
}

func initRCCData() {
	tupleRCCData, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "contract_address", Type: "string"},
			{Name: "data", Type: "bytes"},
		},
	)
	if err != nil {
		panic(err)
	}
	if tupleRCCData.T != abi.TupleTy {
		panic("New TupleRCCData type err")
	}
	TupleRCCData = tupleRCCData
}
