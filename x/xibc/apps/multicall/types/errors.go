package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrScChainEqualToDestChain = sdkerrors.Register(ModuleName, 2, "source chain equals to destination chain")
	ErrInvalidMultiCallEvent   = sdkerrors.Register(ModuleName, 3, "invalid multicall event")
)
