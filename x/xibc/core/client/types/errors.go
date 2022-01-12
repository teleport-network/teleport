package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	host "github.com/teleport-network/teleport/x/xibc/core/host"
)

const moduleName = host.ModuleName + "-" + SubModuleName

// XIBC client sentinel errors
var (
	ErrClientExists           = sdkerrors.Register(moduleName, 2, "light client already exists")
	ErrClientNotFound         = sdkerrors.Register(moduleName, 3, "light client not found")
	ErrInvalidClientMetadata  = sdkerrors.Register(moduleName, 4, "invalid client metadata")
	ErrConsensusStateNotFound = sdkerrors.Register(moduleName, 5, "consensus state not found")
	ErrInvalidConsensus       = sdkerrors.Register(moduleName, 6, "invalid consensus state")
	ErrInvalidClientType      = sdkerrors.Register(moduleName, 7, "invalid client type")
	ErrInvalidHeader          = sdkerrors.Register(moduleName, 8, "invalid client header")
	ErrClientNotActive        = sdkerrors.Register(moduleName, 9, "client is not active")
	ErrUpgradeClient          = sdkerrors.Register(moduleName, 10, "Upgrade client failed")
)
