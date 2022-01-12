package multicall

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	"github.com/teleport-network/teleport/syscontracts"
)

var (
	//go:embed MultiCall.json
	MultiCallJSON []byte // nolint: golint

	// MultiCallContract is the compiled multicall contract
	MultiCallContract evmtypes.CompiledContract

	// MultiCallContractAddress is the deployed multicall contract address
	MultiCallContractAddress common.Address
)

func init() {
	MultiCallContractAddress = common.HexToAddress(syscontracts.MultiCallContractAddress)

	if err := json.Unmarshal(MultiCallJSON, &MultiCallContract); err != nil {
		panic(err)
	}

	if len(MultiCallContract.Bin) == 0 {
		panic("load contract failed")
	}
}
