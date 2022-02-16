package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

const (
	TransferERC20 = uint8(0)
	TransferBase  = uint8(1)
	RemoteCall    = uint8(2)
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

type ERC20TransferData struct {
	TokenAddress common.Address `json:"token_address"`
	Receiver     string         `json:"receiver"`
	Amount       *big.Int       `json:"amount"`
}

type BaseTransferData struct {
	Receiver string   `json:"receiver"`
	Amount   *big.Int `json:"amount"`
}

type RCCData struct {
	ContractAddress string `json:"contract_address"`
	Data            []byte `json:"data"`
}

var (
	TupleERC20TransferData abi.Type
	TupleBaseTransferData  abi.Type
	TupleRCCData           abi.Type
)

func init() {
	initERC20TransferData()
	initBaseTransferData()
	initRCCData()
}

func initERC20TransferData() {
	tupleERC20TransferData, err := abi.NewType(
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
	if tupleERC20TransferData.T != abi.TupleTy {
		panic("New TupleERC20TransferData type err")
	}
	TupleERC20TransferData = tupleERC20TransferData
}

func initBaseTransferData() {
	tupleBaseTransferData, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "receiver", Type: "string"},
			{Name: "amount", Type: "uint256"},
		},
	)
	if err != nil {
		panic(err)
	}
	if tupleBaseTransferData.T != abi.TupleTy {
		panic("New TupleBaseTransferData type err")
	}
	TupleBaseTransferData = tupleBaseTransferData
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
