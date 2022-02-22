package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/teleport-network/teleport/x/xibc/core/host"
)

const (
	SubModuleName = "tss-client"
	moduleName    = host.ModuleName + "-" + SubModuleName
)

// XIBC tss client sentinel errors
var (
	ErrInvalidChainID       = sdkerrors.Register(moduleName, 2, "invalid chain-id")
	ErrDelayPeriodNotPassed = sdkerrors.Register(moduleName, 3, "packet-specified delay period has not been reached")
	ErrInvalidProofSpecs    = sdkerrors.Register(moduleName, 4, "invalid proof specs")
)
