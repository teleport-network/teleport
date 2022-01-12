package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	host "github.com/teleport-network/teleport/x/xibc/core/host"
)

const moduleName = host.ModuleName + "-" + SubModuleName

// XIBC routing sentinel errors
var (
	ErrInvalidRoute = sdkerrors.Register(moduleName, 2, "route not found")
)
