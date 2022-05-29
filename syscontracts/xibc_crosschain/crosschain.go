package crosschain

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	"github.com/teleport-network/teleport/syscontracts"
)

var (
	//go:embed CrossChain.json
	CrossChainJSON []byte // nolint: golint

	// CrossChainContract is the compiled CrossChain contract
	CrossChainContract evmtypes.CompiledContract

	// CrossChainContractAddress is the deployed CrossChain contract address
	CrossChainContractAddress common.Address
)

func init() {
	CrossChainContractAddress = common.HexToAddress(syscontracts.CrossChainContractAddress)

	if err := json.Unmarshal(CrossChainJSON, &CrossChainContract); err != nil {
		panic(err)
	}

	if len(CrossChainContract.Bin) == 0 {
		panic("load contract failed")
	}
}
