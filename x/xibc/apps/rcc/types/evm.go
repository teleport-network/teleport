package types

import (
	"math/big"
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
