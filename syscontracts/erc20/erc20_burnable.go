package erc20

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"

	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	"github.com/teleport-network/teleport/syscontracts"
)

var (
	// ERC20BurnableContract is the compiled ERC20Burnable contract
	ERC20BurnableContract evmtypes.CompiledContract
)

func init() {
	if err := json.Unmarshal(syscontracts.ERC20BurnableJSON, &ERC20BurnableContract); err != nil {
		panic(err)
	}
}
