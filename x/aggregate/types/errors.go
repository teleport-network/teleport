package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// errors
var (
	ErrAggregateDisabled      = sdkerrors.Register(ModuleName, 2, "aggregate module is disabled")
	ErrInternalTokenPair      = sdkerrors.Register(ModuleName, 3, "internal ethereum token mapping error")
	ErrTokenPairNotFound      = sdkerrors.Register(ModuleName, 4, "token pair not found")
	ErrTokenPairAlreadyExists = sdkerrors.Register(ModuleName, 5, "token pair already exists")
	ErrUndefinedOwner         = sdkerrors.Register(ModuleName, 6, "undefined owner of contract pair")
	ErrBalanceInvariance      = sdkerrors.Register(ModuleName, 7, "post transfer balance invariant failed")
	ErrUnexpectedEvent        = sdkerrors.Register(ModuleName, 8, "unexpected event")
	ErrABIPack                = sdkerrors.Register(ModuleName, 9, "contract ABI pack failed")
	ErrABIUnpack              = sdkerrors.Register(ModuleName, 10, "contract ABI unpack failed")
	ErrEVMDenom               = sdkerrors.Register(ModuleName, 11, "EVM denomination registration")
	ErrEVMCall                = sdkerrors.Register(ModuleName, 12, "EVM call unexpected error")
	ErrERC20TokenPairDisabled = sdkerrors.Register(ModuleName, 13, "erc20 token pair is disabled")
	ErrInvalidOriginToken     = sdkerrors.Register(ModuleName, 14, "invalid origin token")
	ErrInvalidOriginChain     = sdkerrors.Register(ModuleName, 15, "invalid origin chain")
	ErrERC20TraceExist        = sdkerrors.Register(ModuleName, 16, "erc20 trace already exist")
	ErrERC20TraceScale        = sdkerrors.Register(ModuleName, 17, "invalid erc20 trace scale")
	ErrInvalidTimePeriod      = sdkerrors.Register(ModuleName, 18, "invalid time period")
	ErrInvalidTimeBasedLimit  = sdkerrors.Register(ModuleName, 19, "invalid time based limit")
	ErrInvalidMaxAmount       = sdkerrors.Register(ModuleName, 20, "invalid max amount")
	ErrInvalidMinAmount       = sdkerrors.Register(ModuleName, 21, "invalid min amount")
)
