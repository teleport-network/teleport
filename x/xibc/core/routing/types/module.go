package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
)

// XIBCModule defines an interface that implements all the callbacks
type XIBCModule interface {
	// OnRecvPacket must return the acknowledgement bytes
	// In the case of an asynchronous acknowledgement, nil should be returned.
	OnRecvPacket(ctx sdk.Context, packetData []byte) (res *sdk.Result, result packettypes.Result, err error)

	OnAcknowledgementPacket(ctx sdk.Context, packetData []byte, result []byte) (*sdk.Result, error)
}
