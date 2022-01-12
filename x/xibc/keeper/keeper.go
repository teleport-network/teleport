package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	clientkeeper "github.com/teleport-network/teleport/x/xibc/core/client/keeper"
	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	packetkeeper "github.com/teleport-network/teleport/x/xibc/core/packet/keeper"
	routingkeeper "github.com/teleport-network/teleport/x/xibc/core/routing/keeper"
	routingtypes "github.com/teleport-network/teleport/x/xibc/core/routing/types"
	"github.com/teleport-network/teleport/x/xibc/types"
)

var _ types.QueryServer = (*Keeper)(nil)

// Keeper defines each keeper for XIBC
type Keeper struct {
	// implements gRPC QueryServer interface
	types.QueryServer
	cdc           codec.BinaryCodec
	ClientKeeper  clientkeeper.Keeper
	PacketKeeper  packetkeeper.Keeper
	RoutingKeeper routingkeeper.Keeper
}

// NewKeeper creates a new xibc Keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	key sdk.StoreKey,
	paramSpace paramtypes.Subspace,
	stakingKeeper clienttypes.StakingKeeper,
) *Keeper {
	clientKeeper := clientkeeper.NewKeeper(cdc, key, paramSpace, stakingKeeper)
	routingKeeper := routingkeeper.NewKeeper(key)
	packetkeeper := packetkeeper.NewKeeper(cdc, key, clientKeeper)

	return &Keeper{
		cdc:           cdc,
		ClientKeeper:  clientKeeper,
		PacketKeeper:  packetkeeper,
		RoutingKeeper: routingKeeper,
	}
}

// Codec returns the XIBC module codec.
func (k Keeper) Codec() codec.BinaryCodec {
	return k.cdc
}

// SetRouter sets the Router in XIBC Keeper and seals it. The method panics if
// there is an existing router that's already sealed.
func (k *Keeper) SetRouter(rtr *routingtypes.Router) {
	k.RoutingKeeper.SetRouter(rtr)
}
