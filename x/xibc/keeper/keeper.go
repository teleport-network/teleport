package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	clientkeeper "github.com/teleport-network/teleport/x/xibc/core/client/keeper"
	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	packetkeeper "github.com/teleport-network/teleport/x/xibc/core/packet/keeper"
	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
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
}

// NewKeeper creates a new xibc Keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	key sdk.StoreKey,
	paramSpace paramtypes.Subspace,
	stakingKeeper clienttypes.StakingKeeper,
	accountKeeper packettypes.AccountKeeper,
	evmKeeper packettypes.EVMKeeper,
) *Keeper {
	clientKeeper := clientkeeper.NewKeeper(cdc, key, paramSpace, stakingKeeper)
	packetkeeper := packetkeeper.NewKeeper(cdc, key, clientKeeper, accountKeeper, evmKeeper)

	return &Keeper{
		cdc:           cdc,
		ClientKeeper:  clientKeeper,
		PacketKeeper:  packetkeeper,
	}
}

// Codec returns the XIBC module codec.
func (k Keeper) Codec() codec.BinaryCodec {
	return k.cdc
}
