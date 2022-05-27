package crossChain

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	"github.com/teleport-network/teleport/syscontracts"
)

var (
	//go:embed cross_chain.json
	CrossChainJSON []byte // nolint: golint

	// CrossChainContract is the compiled agent contract
	CrossChainContract evmtypes.CompiledContract

	// CrossChainAddress is the deployed agent contract address
	CrossChainAddress common.Address
)

func init() {
	CrossChainAddress = common.HexToAddress(syscontracts.CrossChainAddress)

	if err := json.Unmarshal(CrossChainJSON, &CrossChainContract); err != nil {
		panic(err)
	}

	if len(CrossChainContract.Bin) == 0 {
		panic("load contract failed")
	}
}
