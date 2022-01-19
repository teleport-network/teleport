package wtele

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	"github.com/teleport-network/teleport/syscontracts"
)

var (
	//go:embed WTELE.json
	WTELEJSON []byte // nolint: golint

	// WTELEContract is the compiled wtele contract
	WTELEContract evmtypes.CompiledContract

	// WTELEContract is the deployed wtele contract address
	WTELEContractAddress common.Address
)

func init() {
	WTELEContractAddress = common.HexToAddress(syscontracts.WTELEContractAddress)

	if err := json.Unmarshal(WTELEJSON, &WTELEContract); err != nil {
		panic(err)
	}

	if len(WTELEContract.Bin) == 0 {
		panic("load contract failed")
	}
}
