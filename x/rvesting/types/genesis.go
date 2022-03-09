package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:     DefaultParams(),
		From:       "",
		InitReward: sdk.NewCoins(),
	}
}

func ValidateGenesis(data *GenesisState) error {
	if err := data.Params.validate(); err != nil {
		return err
	}
	if len(data.From) != 0 {
		if _, err := sdk.AccAddressFromBech32(data.From); err != nil {
			return err
		}
		return data.InitReward.Validate()
	}
	return nil
}

func NewGenesisState(params Params) *GenesisState {
	return &GenesisState{
		Params:     params,
		From:       "",
		InitReward: sdk.NewCoins(),
	}
}

// GetGenesisStateFromAppState returns x/rvesting GenesisState given raw application
// genesis state.
func GetGenesisStateFromAppState(cdc codec.JSONCodec, appState map[string]json.RawMessage) *GenesisState {
	var genesisState GenesisState

	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}

	return &genesisState
}
