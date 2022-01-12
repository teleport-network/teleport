package xibctesting

import (
	"time"

	xibctmtypes "github.com/teleport-network/teleport/x/xibc/clients/light-clients/tendermint/types"
	"github.com/teleport-network/teleport/x/xibc/exported"
)

type ClientConfig interface {
	GetClientType() string
}

type TendermintConfig struct {
	TrustLevel                   xibctmtypes.Fraction
	TrustingPeriod               time.Duration
	UnbondingPeriod              time.Duration
	MaxClockDrift                time.Duration
	AllowUpdateAfterExpiry       bool
	AllowUpdateAfterMisbehaviour bool
}

func NewTendermintConfig() *TendermintConfig {
	return &TendermintConfig{
		TrustLevel:                   DefaultTrustLevel,
		TrustingPeriod:               TrustingPeriod,
		UnbondingPeriod:              UnbondingPeriod,
		MaxClockDrift:                MaxClockDrift,
		AllowUpdateAfterExpiry:       false,
		AllowUpdateAfterMisbehaviour: false,
	}
}

func (tmcfg *TendermintConfig) GetClientType() string {
	return exported.Tendermint
}
