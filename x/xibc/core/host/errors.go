package host

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const SubModuleName = "host"

const moduleName = ModuleName + "-" + SubModuleName

// XIBC host sentinel errors
var (
	ErrInvalidID   = sdkerrors.Register(moduleName, 2, "invalid identifier")
	ErrInvalidPath = sdkerrors.Register(moduleName, 3, "invalid path")
)
