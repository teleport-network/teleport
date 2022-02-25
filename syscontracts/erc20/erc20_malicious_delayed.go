package erc20

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	"github.com/teleport-network/teleport/syscontracts"
	"github.com/teleport-network/teleport/x/aggregate/types"
)

// This is an evil token. Whenever an A -> B transfer is called,
// a predefined C is given a massive allowance on B.
var (
	// ERC20MaliciousDelayedContract is the compiled erc20 contract
	ERC20MaliciousDelayedContract evmtypes.CompiledContract

	// ERC20MaliciousDelayedAddress is the erc20 module address
	ERC20MaliciousDelayedAddress common.Address
)

func init() {
	ERC20MaliciousDelayedAddress = types.ModuleAddress

	if err := json.Unmarshal(syscontracts.ERC20MaliciousDelayedJSON, &ERC20MaliciousDelayedContract); err != nil {
		panic(err)
	}

	if len(ERC20MaliciousDelayedContract.Bin) == 0 {
		panic("load contract failed")
	}
}
