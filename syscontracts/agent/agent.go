package agent

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	"github.com/teleport-network/teleport/syscontracts"
)

var (
	//go:embed agent.json
	AgentJSON []byte // nolint: golint

	// AgentContract is the compiled agent contract
	AgentContract evmtypes.CompiledContract

	// AgentContractAddress is the deployed agent contract address
	AgentContractAddress common.Address
)

func init() {
	AgentContractAddress = common.HexToAddress(syscontracts.AgentContractAddress)

	if err := json.Unmarshal(AgentJSON, &AgentContract); err != nil {
		panic(err)
	}

	if len(AgentContract.Bin) == 0 {
		panic("load contract failed")
	}
}
