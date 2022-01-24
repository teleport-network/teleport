package rcc

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	"github.com/teleport-network/teleport/syscontracts"
)

var (
	//go:embed rcc.json
	RCCJSON []byte // nolint: golint

	// RCCContract is the compiled rcc contract
	RCCContract evmtypes.CompiledContract

	// RCCContract is the deployed rcc contract address
	RCCContractAddress common.Address
)

func init() {
	RCCContractAddress = common.HexToAddress(syscontracts.RCCContractAddress)

	if err := json.Unmarshal(RCCJSON, &RCCContract); err != nil {
		panic(err)
	}

	if len(RCCContract.Bin) == 0 {
		panic("load contract failed")
	}
}
