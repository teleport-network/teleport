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
	ErrWritingEthTxData         = sdkerrors.Register(moduleName, 6, "writing ethereum tx data error")
	ErrScChainEqualToDestChain  = sdkerrors.Register(moduleName, 7, "source chain equals to destination chain")
	ErrRelayerNotFound          = sdkerrors.Register(moduleName, 8, "relayer not found")
	ErrInvalidRelayer           = sdkerrors.Register(moduleName, 9, "invalid relayer")
	ErrInvalidSequence          = sdkerrors.Register(moduleName, 10, "invalid sequence")
	ErrInvalidSrcChain          = sdkerrors.Register(moduleName, 11, "invalid source chain")
	ErrInvalidDestChain         = sdkerrors.Register(moduleName, 12, "invalid destination chain")
	ErrInvalidRelayChain        = sdkerrors.Register(moduleName, 13, "invalid relay chain")
	ErrAbiPack                  = sdkerrors.Register(moduleName, 14, "err pack to bytes")
	ErrDecodeAbi                = sdkerrors.Register(moduleName, 15, "err decode abi ")
)
