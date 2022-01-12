package types

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/ethereum/go-ethereum/common"
)

const (
	// ModuleName defines the XIBC remote contract call name
	ModuleName = "CONTRACT"

	// PortID is the default port id that remote contract call module binds to
	PortID = ModuleName
)

var (
	// ModuleAddress is the native module address for EVM
	// 0xfef812Ed2Bf63E7eE056931d54A6292fcbbaDFaA
	ModuleAddress common.Address
)

func init() {
	ModuleAddress = common.BytesToAddress(authtypes.NewModuleAddress(ModuleName).Bytes())
}
