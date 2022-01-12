package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// AttoTele defines the default coin denomination used in Teleport in:
	//
	// - Staking parameters: denomination used as stake in the dPoS chain
	// - Mint parameters: denomination minted due to fee distribution rewards
	// - Governance parameters: denomination used for spam prevention in proposal deposits
	// - Crisis parameters: constant fee denomination used for spam prevention to check broken invariant
	// - EVM parameters: denomination used for running EVM state transitions in Teleport.
	AttoTele string = "atele"

	// BaseDenomUnit defines the base denomination unit for Teles.
	// 1 tele = 1x10^{BaseDenomUnit} atele
	BaseDenomUnit = 18

	// DefaultGasPrice is default gas price for evm transactions
	DefaultGasPrice = 20
)

// PowerReduction defines the default power reduction value for staking
var PowerReduction = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(BaseDenomUnit), nil))

// NewTeleCoin is a utility function that returns an "atele" coin with the given sdk.Int amount.
// The function will panic if the provided amount is negative.
func NewTeleCoin(amount sdk.Int) sdk.Coin {
	return sdk.NewCoin(AttoTele, amount)
}

// NewTeleDecCoin is a utility function that returns an "atele" decimal coin with the given sdk.Int amount.
// The function will panic if the provided amount is negative.
func NewTeleDecCoin(amount sdk.Int) sdk.DecCoin {
	return sdk.NewDecCoin(AttoTele, amount)
}

// NewTeleCoinInt64 is a utility function that returns an "atele" coin with the given int64 amount.
// The function will panic if the provided amount is negative.
func NewTeleCoinInt64(amount int64) sdk.Coin {
	return sdk.NewInt64Coin(AttoTele, amount)
}
