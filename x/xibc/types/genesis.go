package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
)

var _ codectypes.UnpackInterfacesMessage = GenesisState{}

// DefaultGenesisState returns the xibc module's default genesis state.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		ClientGenesis: clienttypes.DefaultGenesisState(),
		PacketGenesis: packettypes.DefaultGenesisState(),
	}
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (gs GenesisState) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return gs.ClientGenesis.UnpackInterfaces(unpacker)
}

// Validate performs basic genesis state validation returning an error upon any failure
func (gs *GenesisState) Validate() error {
	if err := gs.ClientGenesis.Validate(); err != nil {
		return err
	}

	return gs.PacketGenesis.Validate()
}
