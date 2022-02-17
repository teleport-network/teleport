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
	// ERC20DirectBalanceManipulationContract is the compiled erc20 contract
	ERC20DirectBalanceManipulationContract evmtypes.CompiledContract

	// ERC20DirectBalanceManipulationAddress is the erc20 module address
	ERC20DirectBalanceManipulationAddress common.Address
)

func init() {
	ERC20DirectBalanceManipulationAddress = types.ModuleAddress

	if err := json.Unmarshal(syscontracts.ERC20DirectBalanceManipulationJSON, &ERC20DirectBalanceManipulationContract); err != nil {
		panic(err)
	}

	if len(ERC20DirectBalanceManipulationContract.Bin) == 0 {
		panic("load contract failed")
	}
}
