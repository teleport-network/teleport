package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrABIPack                 = sdkerrors.Register(ModuleName, 2, "contract ABI pack failed")
	ErrScChainEqualToDestChain = sdkerrors.Register(ModuleName, 3, "source chain equals to destination chain")
)
