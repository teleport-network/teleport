package types

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/ethereum/go-ethereum/common"
)

const (
	// ModuleName defines the XIBC transfer name
	ModuleName = "FT"

	// PortID is the default port id that transfer module binds to
	PortID = ModuleName
)

var (
	// ModuleAddress is the native module address for EVM
	// 0xDE152Fc3Bc10A8878677FD17c44aE633D9EBF737
	ModuleAddress common.Address
)

func init() {
	ModuleAddress = common.BytesToAddress(authtypes.NewModuleAddress(ModuleName).Bytes())
}
