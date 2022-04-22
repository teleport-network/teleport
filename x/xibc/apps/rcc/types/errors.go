package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidPacket           = sdkerrors.Register(ModuleName, 2, "invalid packet")
	ErrABIPack                 = sdkerrors.Register(ModuleName, 3, "contract ABI pack failed")
	ErrScChainEqualToDestChain = sdkerrors.Register(ModuleName, 4, "source chain equals to destination chain")
	ErrInvalidSequence         = sdkerrors.Register(ModuleName, 5, "invalid sequence")
	ErrInvalidSrcChain         = sdkerrors.Register(ModuleName, 6, "invalid source chain")
	ErrInvalidDestChain        = sdkerrors.Register(ModuleName, 7, "invalid destination chain")
	ErrInvalidSender           = sdkerrors.Register(ModuleName, 8, "invalid sender")
	ErrInvalidAddress          = sdkerrors.Register(ModuleName, 9, "invalid address")
)
