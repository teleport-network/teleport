package types

import (
	"math/big"

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
