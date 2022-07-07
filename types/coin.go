package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// AttoBit defines the default coin denomination used in BitNetwork in:
	//
	// - Staking parameters: denomination used as stake in the dPoS chain
	// - Mint parameters: denomination minted due to fee distribution rewards
	// - Governance parameters: denomination used for spam prevention in proposal deposits
	// - Crisis parameters: constant fee denomination used for spam prevention to check broken invariant
	// - EVM parameters: denomination used for running EVM state transitions in BitNetwork.
	AttoBit = "abit"

	// DisplayDenom defines the denomination displayed to users in client applications.
	DisplayDenom = "bit"
	// BaseDenom defines to the default denomination used in BitNetwork (staking, EVN, governance, etc)
	BaseDenom = AttoBit

	// BaseDenomUnit defines the base denomination unit for Bits.
	// 1 bit = 1x10^{BaseDenomUnit} abit
	BaseDenomUnit = 18
)

// PowerReduction defines the default power reduction value for staking
// (BaseDenomUnit - 2) indicates that staking tokens:power = 1:100, which means we consider 0.01 BaseDenom as 1 power
var PowerReduction = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(BaseDenomUnit-2), nil))

// NewBitCoin is a utility function that returns an "abit" coin with the given sdk.Int amount.
// The function will panic if the provided amount is negative.
func NewBitCoin(amount sdk.Int) sdk.Coin {
	return sdk.NewCoin(AttoBit, amount)
}

// NewBitDecCoin is a utility function that returns an "abit" decimal coin with the given sdk.Int amount.
// The function will panic if the provided amount is negative.
func NewBitDecCoin(amount sdk.Int) sdk.DecCoin {
	return sdk.NewDecCoin(AttoBit, amount)
}

// NewBitCoinInt64 is a utility function that returns an "abit" coin with the given int64 amount.
// The function will panic if the provided amount is negative.
func NewBitCoinInt64(amount int64) sdk.Coin {
	return sdk.NewInt64Coin(AttoBit, amount)
}
