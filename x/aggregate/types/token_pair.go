package types

import (
	"github.com/tendermint/tendermint/crypto/tmhash"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"
	ethermint "github.com/tharsis/ethermint/types"
)

// NewTokenPair returns an instance of TokenPair
func NewTokenPair(erc20Address common.Address, denom string, enabled bool, contractOwner Owner) TokenPair {
	return TokenPair{
		ERC20Address:  erc20Address.String(),
		Denom:         denom,
		Enabled:       true,
		ContractOwner: contractOwner,
	}
}

// GetID returns the SHA256 hash of the ERC20 address and denomination
func (tp TokenPair) GetID() []byte {
	id := tp.ERC20Address + "|" + tp.Denom
	return tmhash.Sum([]byte(id))
}

// GetErc20Contract casts the hex string address of the ERC20 to common.Address
func (tp TokenPair) GetERC20Contract() common.Address {
	return common.HexToAddress(tp.ERC20Address)
}

// Validate performs a stateless validation of a TokenPair
func (tp TokenPair) Validate() error {
	if err := sdk.ValidateDenom(tp.Denom); err != nil {
		return err
	}

	if err := ethermint.ValidateAddress(tp.ERC20Address); err != nil {
		return err
	}

	return nil
}

// IsNativeCoin returns true if the owner of the ERC20 contract is the
// aggregate module account
func (tp TokenPair) IsNativeCoin() bool {
	return tp.ContractOwner == OWNER_MODULE
}

// IsNativeERC20 returns true if the owner of the ERC20 contract not the
// aggregate module account
func (tp TokenPair) IsNativeERC20() bool {
	return tp.ContractOwner == OWNER_EXTERNAL
}
