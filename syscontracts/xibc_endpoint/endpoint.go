package endpoint

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	"github.com/teleport-network/teleport/syscontracts"
)

var (
	//go:embed Endpoint.json
	EndpointJSON []byte // nolint: golint

	// EndpointContract is the compiled Endpoint contract
	EndpointContract evmtypes.CompiledContract

	// EndpointContractAddress is the deployed Endpoint contract address
	EndpointContractAddress common.Address
)

func init() {
	EndpointContractAddress = common.HexToAddress(syscontracts.EndpointContractAddress)

	if err := json.Unmarshal(EndpointJSON, &EndpointContract); err != nil {
		panic(err)
	}

	if len(EndpointContract.Bin) == 0 {
		panic("load contract failed")
	}
}
