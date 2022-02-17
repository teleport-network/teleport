package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrWritingEthTxData        = sdkerrors.Register(ModuleName, 2, "writing ethereum tx data error")
	ErrScChainEqualToDestChain = sdkerrors.Register(ModuleName, 3, "source chain equals to destination chain")
)
