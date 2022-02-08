package types

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/ethereum/go-ethereum/common"
)

const (
	// SubModuleName defines the XIBC packets name
	SubModuleName = "packet"
)

var (
	// ModuleAddress is the native module address for EVM
	// 0x7426aFC489D0eeF99a0B438DEF226aD139F75235
	ModuleAddress common.Address
)

func init() {
	ModuleAddress = common.BytesToAddress(authtypes.NewModuleAddress(SubModuleName).Bytes())
}
