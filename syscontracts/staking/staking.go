package staking

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"

	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	"github.com/teleport-network/teleport/syscontracts"
)

var (
	// StakingContract is the compiled staking contract
	StakingContract evmtypes.CompiledContract

	// StakingAddress is the staking address
	StakingAddress common.Address
)

func init() {
	StakingAddress = common.HexToAddress(syscontracts.StakingContractAddress)

	if err := json.Unmarshal(syscontracts.StakingJSON, &StakingContract); err != nil {
		panic(err)
	}

	if len(StakingContract.Bin) == 0 {
		panic("load contract failed")
	}
}
