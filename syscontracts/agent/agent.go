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
	AGENTJSON []byte // nolint: golint

	// AGENTContract is the compiled agent contract
	AGENTContract evmtypes.CompiledContract

	// AGENTContractAddress is the deployed agent contract address
	AGENTContractAddress common.Address
)

func init() {
	AGENTContractAddress = common.HexToAddress(syscontracts.CUSTODIANContractAddress)

	if err := json.Unmarshal(AGENTJSON, &AGENTContract); err != nil {
		panic(err)
	}

	if len(AGENTContract.Bin) == 0 {
		panic("load contract failed")
	}
}
