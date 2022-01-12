package types

import (
	"math/big"

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
