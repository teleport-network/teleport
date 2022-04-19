package keeper

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/teleport-network/teleport/x/xibc/core/client/types"
)

// RegisterRelayers saves the relayers under the specified chainname
func (k Keeper) RegisterRelayers(
	ctx sdk.Context,
	address string,
	chains []string,
	addresses []string,
) {
	store := k.RelayerStore(ctx)
	ir := &types.IdentifiedRelayer{
		Address:   address,
		Chains:    chains,
		Addresses: addresses,
	}
	irBz := k.cdc.MustMarshal(ir)
	store.Set([]byte(address), irBz)
}

// AuthRelayer asserts whether a relayer is already registered
func (k Keeper) AuthRelayer(ctx sdk.Context, chainName string, relayer string) bool {
	if ir, found := k.GetRelayer(ctx, relayer); found {
		for _, chain := range ir.Chains {
			if chain == chainName {
				return true
			}
		}
	}
	return false
}

// GetAllRelayers returns all registered relayer addresses
func (k Keeper) GetAllRelayers(ctx sdk.Context) (relayers []types.IdentifiedRelayer) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.KeyRelayers))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var ir = &types.IdentifiedRelayer{}
		k.cdc.MustUnmarshal(iterator.Value(), ir)
		relayers = append(relayers, *ir)
	}
	return
}

func (k Keeper) GetRelayer(ctx sdk.Context, address string) (types.IdentifiedRelayer, bool) {
	store := k.RelayerStore(ctx)
	bz := store.Get([]byte(address))
	if bz == nil {
		return types.IdentifiedRelayer{}, false
	}

	var ir types.IdentifiedRelayer
	k.cdc.MustUnmarshal(bz, &ir)
	return ir, true
}

func (k Keeper) GetRelayerAddressOnOtherChain(ctx sdk.Context, chainName string, address string) (string, bool) {
	if ir, found := k.GetRelayer(ctx, address); found {
		for i, chain := range ir.Chains {
			if chain == chainName {
				return ir.Addresses[i], true
			}
		}
	}

	return "", false
}

func (k Keeper) GetRelayerAddressOnTeleport(ctx sdk.Context, chainName string, address string) (string, bool) {
	for _, ir := range k.GetAllRelayers(ctx) {
		for i, chain := range ir.Chains {
			if chain == chainName && strings.EqualFold(ir.Addresses[i], address) {
				return ir.Address, true
			}
		}
	}
	return "", false
}
