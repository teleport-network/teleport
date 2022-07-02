package keeper

import (
	"context"
	"github.com/teleport-network/teleport/x/aggregate/types"
)

var _ types.MsgServer = &Keeper{}

// ConvertCoin converts ERC20 tokens into Cosmos-native Coins for both
// Cosmos-native and ERC20 TokenPair Owners
func (k Keeper) ConvertCoin(
	goCtx context.Context,
	msg *types.MsgConvertCoin,
) (*types.MsgConvertCoinResponse, error) {
	return nil, types.ErrUndefinedOwner
}

// ConvertERC20 converts ERC20 tokens into Cosmos-native Coins for both
// Cosmos-native and ERC20 TokenPair Owners
func (k Keeper) ConvertERC20(
	goCtx context.Context,
	msg *types.MsgConvertERC20,
) (*types.MsgConvertERC20Response, error) {
	return nil, types.ErrUndefinedOwner
}
