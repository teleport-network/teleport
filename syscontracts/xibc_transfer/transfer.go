package transfer

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	"github.com/teleport-network/teleport/syscontracts"
)

var (
	//go:embed Transfer.json
	TransferJSON []byte // nolint: golint

	// TransferContract is the compiled transfer contract
	TransferContract evmtypes.CompiledContract

	// TransferContract is the deployed transfer contract address
	TransferContractAddress common.Address
)

func init() {
	TransferContractAddress = common.HexToAddress(syscontracts.TransferContractAddress)

	if err := json.Unmarshal(TransferJSON, &TransferContract); err != nil {
		panic(err)
	}

	if len(TransferContract.Bin) == 0 {
		panic("load contract failed")
	}
}
