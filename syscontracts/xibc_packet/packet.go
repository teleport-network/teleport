package packet

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	"github.com/teleport-network/teleport/syscontracts"
)

var (
	//go:embed packet.json
	PacketJSON []byte // nolint: golint

	// PacketContract is the compiled packet contract
	PacketContract evmtypes.CompiledContract

	// PacketContractAddress is the deployed packet contract address
	PacketContractAddress common.Address
)

func init() {
	PacketContractAddress = common.HexToAddress(syscontracts.PacketContractAddress)

	if err := json.Unmarshal(PacketJSON, &PacketContract); err != nil {
		panic(err)
	}

	if len(PacketContract.Bin) == 0 {
		panic("load contract failed")
	}
}
