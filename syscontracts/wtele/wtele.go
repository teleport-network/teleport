package wtele

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	"github.com/teleport-network/teleport/syscontracts"
)

var (
	// WTELEContract is the compiled wtele contract
	WTELEContract evmtypes.CompiledContract

	// WTELEContract is the deployed wtele contract address
	WTELEContractAddress common.Address
)

func init() {
	WTELEContractAddress = common.HexToAddress(syscontracts.WTELEContractAddress)

	var contractBinRuntime syscontracts.CompiledContract
	if err := json.Unmarshal(syscontracts.WTELEJSON, &contractBinRuntime); err != nil {
		panic(err)
	}

	WTELEContract.ABI = contractBinRuntime.ABI
	WTELEContract.Bin = contractBinRuntime.Bin

	if len(WTELEContract.Bin) == 0 {
		panic("load contract failed")
	}
}
