package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/teleport-network/teleport/x/xibc/core/host"
)

const moduleName = host.ModuleName + "-" + SubModuleName

// XIBC packet sentinel errors
var (
	ErrInvalidPacket            = sdkerrors.Register(moduleName, 2, "invalid packet")
	ErrInvalidAcknowledgement   = sdkerrors.Register(moduleName, 3, "invalid acknowledgement")
	ErrPacketCommitmentNotFound = sdkerrors.Register(moduleName, 4, "packet commitment not found")
	ErrAcknowledgementExists    = sdkerrors.Register(moduleName, 5, "packet acknowledgement already exists")
	ErrWritingEthTxPayload      = sdkerrors.Register(moduleName, 6, "writing ethereum tx payload error")
	ErrScChainEqualToDestChain  = sdkerrors.Register(moduleName, 7, "source chain equals to destination chain")
)
