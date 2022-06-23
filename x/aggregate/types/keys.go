package types

import (
	"github.com/ethereum/go-ethereum/common"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// constants
const (
	// module name
	ModuleName = "aggregate"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for gov routing
	RouterKey = ModuleName

	// QuerierRoute is the querier route for aggregate
	QuerierRoute = ModuleName
)

// ModuleAddress is the native module address for EVM
// 0xEE3c65B5c7F4DD0ebeD8bF046725e273e3eeeD3c
var ModuleAddress common.Address

func init() {
	ModuleAddress = common.BytesToAddress(authtypes.NewModuleAddress(ModuleName).Bytes())
}

// prefix bytes for the EVM persistent store
const (
	prefixTokenPair = iota + 1
	prefixTokenPairByERC20
	prefixTokenPairByDenom
)

// KVStore key prefixes
var (
	KeyPrefixTokenPair        = []byte{prefixTokenPair}
	KeyPrefixTokenPairByERC20 = []byte{prefixTokenPairByERC20}
	KeyPrefixTokenPairByDenom = []byte{prefixTokenPairByDenom}
)
