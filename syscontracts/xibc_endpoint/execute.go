package endpoint

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	"github.com/teleport-network/teleport/syscontracts"
)

var (
	//go:embed Execute.json
	ExeceuteJSON []byte // nolint: golint

	// ExecuteContract is the compiled Execute contract
	ExecuteContract evmtypes.CompiledContract

	// ExecuteContractAddress is the deployed Execute contract address
	ExecuteContractAddress common.Address
)

func init() {
	ExecuteContractAddress = common.HexToAddress(syscontracts.ExecuteContractAddress)

	if err := json.Unmarshal(ExeceuteJSON, &ExecuteContract); err != nil {
		panic(err)
	}

	if len(ExecuteContract.Bin) == 0 {
		panic("load contract failed")
	}
}
