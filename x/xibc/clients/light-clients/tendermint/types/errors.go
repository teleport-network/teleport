package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	host "github.com/teleport-network/teleport/x/xibc/core/host"
)

const (
	SubModuleName = "tendermint-client"
	moduleName    = host.ModuleName + "-" + SubModuleName
)

// XIBC tendermint client sentinel errors
var (
	ErrInvalidChainID         = sdkerrors.Register(moduleName, 2, "invalid chain-id")
	ErrInvalidTrustingPeriod  = sdkerrors.Register(moduleName, 3, "invalid trusting period")
	ErrInvalidUnbondingPeriod = sdkerrors.Register(moduleName, 4, "invalid unbonding period")
	ErrInvalidHeaderHeight    = sdkerrors.Register(moduleName, 5, "invalid header height")
	ErrInvalidMaxClockDrift   = sdkerrors.Register(moduleName, 6, "invalid max clock drift")
	ErrProcessedTimeNotFound  = sdkerrors.Register(moduleName, 7, "processed time not found")
	ErrDelayPeriodNotPassed   = sdkerrors.Register(moduleName, 8, "packet-specified delay period has not been reached")
	ErrInvalidProofSpecs      = sdkerrors.Register(moduleName, 9, "invalid proof specs")
	ErrInvalidValidatorSet    = sdkerrors.Register(moduleName, 10, "invalid validator set")
)
