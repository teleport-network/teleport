package types

import (
	"fmt"
)

// NewGenesisState creates a new genesis state
func NewGenesisState(params Params, pairs []TokenPair) GenesisState {
	return GenesisState{
		Params:     params,
		TokenPairs: pairs,
	}
}

// DefaultGenesisState sets default evm genesis state with empty accounts and default params and chain config values.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any failure
func (gs GenesisState) Validate() error {
	seenErc20 := make(map[string]bool)
	seenDenom := make(map[string]bool)

	for _, b := range gs.TokenPairs {
		if seenErc20[b.ERC20Address] {
			return fmt.Errorf("token ERC20 contract duplicated on genesis '%s'", b.ERC20Address)
		}
		if seenDenom[b.Denoms[0]] {
			return fmt.Errorf("coin denomination duplicated on genesis: '%s'", b.Denoms[0])
		}

		if err := b.Validate(); err != nil {
			return err
		}

		seenErc20[b.ERC20Address] = true
		seenDenom[b.Denoms[0]] = true
	}

	return gs.Params.Validate()
}
