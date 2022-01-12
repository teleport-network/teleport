package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	ethtypes "github.com/teleport-network/teleport/x/xibc/clients/light-clients/eth/types"
	tmtypes "github.com/teleport-network/teleport/x/xibc/clients/light-clients/tendermint/types"
	tsstypes "github.com/teleport-network/teleport/x/xibc/clients/tss-client/types"
	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	commitmenttypes "github.com/teleport-network/teleport/x/xibc/core/commitment/types"
	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
)

// RegisterInterfaces registers x/ibc interfaces into protobuf Any.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	clienttypes.RegisterInterfaces(registry)
	packettypes.RegisterInterfaces(registry)
	commitmenttypes.RegisterInterfaces(registry)
	tmtypes.RegisterInterfaces(registry)
	ethtypes.RegisterInterfaces(registry)
	tsstypes.RegisterInterfaces(registry)
}
