package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/teleport-network/teleport/x/xibc/core/host"
)

const (
	SubModuleName = "eth-client"
	moduleName    = host.ModuleName + "-" + SubModuleName
)

// XIBC evm client sentinel errors
var (
	ErrInvalidGenesisBlock = sdkerrors.Register(moduleName, 2, "invalid genesis block")
	ErrFutureBlock         = sdkerrors.Register(moduleName, 3, "block in the future")
	ErrInvalidMixDigest    = sdkerrors.Register(moduleName, 4, "non-zero mix digest")
	ErrInvalidDifficulty   = sdkerrors.Register(moduleName, 5, "invalid difficulty")
	ErrWrongDifficulty     = sdkerrors.Register(moduleName, 6, "wrong difficulty")
	ErrInvalidProof        = sdkerrors.Register(moduleName, 7, "invalid proof")
	ErrUnmarshalInterface  = sdkerrors.Register(moduleName, 8, "unmarshal field")
	ErrHeader              = sdkerrors.Register(moduleName, 9, "header invalid")
)
