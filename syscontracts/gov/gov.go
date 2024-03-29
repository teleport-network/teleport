package gov

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"

	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	"github.com/teleport-network/teleport/syscontracts"
)

var (
	// GovContract is the compiled gov contract
	GovContract evmtypes.CompiledContract

	// GovAddress is the gov contract address
	GovAddress common.Address
)

func init() {
	GovAddress = common.HexToAddress(syscontracts.GovContractAddress)

	var contractBinRuntime syscontracts.CompiledContract
	if err := json.Unmarshal(syscontracts.GovJSON, &contractBinRuntime); err != nil {
		panic(err)
	}

	GovContract.ABI = contractBinRuntime.ABI
	GovContract.Bin = contractBinRuntime.Bin

	if len(GovContract.Bin) == 0 {
		panic("load contract failed")
	}
}
