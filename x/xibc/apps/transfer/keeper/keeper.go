package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/teleport-network/teleport/x/xibc/apps/transfer/types"
)

type Keeper struct {
	cdc           codec.BinaryCodec
	accountKeeper types.AccountKeeper
	packetKeeper  types.PacketKeeper
	clientKeeper  types.ClientKeeper
	evmKeeper     types.EVMKeeper
}

// NewKeeper creates a new XIBC transfer Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	accountKeeper types.AccountKeeper,
	packetKeeper types.PacketKeeper,
	clientKeeper types.ClientKeeper,
	evmKeeper types.EVMKeeper,
) Keeper {
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the XIBC transfer module account has not been set")
	}

	return Keeper{
		cdc:           cdc,
		accountKeeper: accountKeeper,
		packetKeeper:  packetKeeper,
		clientKeeper:  clientKeeper,
		evmKeeper:     evmKeeper,
	}
}

// GetTransferModuleAddress returns the transfers ModuleAddress
func (k Keeper) GetTransferModuleAddress() sdk.AccAddress {
	return k.accountKeeper.GetModuleAddress(types.ModuleName)
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
