package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	host "github.com/teleport-network/teleport/x/xibc/core/host"
)

const (
	SubModuleName = "bsc-client"
	moduleName    = host.ModuleName + "-" + SubModuleName
)

var (
	ErrInvalidGenesisBlock   = sdkerrors.Register(moduleName, 2, "invalid genesis block")
	ErrInvalidValidatorBytes = sdkerrors.Register(moduleName, 3, "invalid validators bytes length")
	ErrMissingVanity         = sdkerrors.Register(moduleName, 4, "extra-data 32 byte vanity prefix missing")    // ErrMissingVanity is returned if a block's extra-data section is shorter than 32 bytes, which is required to store the signer vanity
	ErrMissingSignature      = sdkerrors.Register(moduleName, 5, "extra-data 65 byte signature suffix missing") // ErrMissingSignature is returned if a block's extra-data section doesn't seem to contain a 65 byte secp256k1 signature
	ErrInvalidMixDigest      = sdkerrors.Register(moduleName, 6, "non-zero mix digest")                         // ErrInvalidMixDigest is returned if a block's mix digest is non-zero
	ErrInvalidUncleHash      = sdkerrors.Register(moduleName, 7, "non empty uncle hash")                        // ErrInvalidUncleHash is returned if a block contains an non-empty uncle list
	ErrInvalidDifficulty     = sdkerrors.Register(moduleName, 8, "invalid difficulty")                          // ErrInvalidDifficulty is returned if the difficulty of a block is missing
	ErrUnknownAncestor       = sdkerrors.Register(moduleName, 9, "unknown ancestor")
	ErrCoinBaseMisMatch      = sdkerrors.Register(moduleName, 10, "coinbase do not match with signature")               // ErrCoinBaseMisMatch is returned if a header's coinbase do not match with signature
	ErrUnauthorizedValidator = sdkerrors.Register(moduleName, 11, "unauthorized validator")                             // ErrUnauthorizedValidator is returned if a header is signed by a non-authorized entity
	ErrRecentlySigned        = sdkerrors.Register(moduleName, 12, "recently signed")                                    // ErrRecentlySigned is returned if a header is signed by an authorized entity that already signed a header recently, thus is temporarily not allowed to
	ErrWrongDifficulty       = sdkerrors.Register(moduleName, 13, "wrong difficulty")                                   // ErrWrongDifficulty is returned if the difficulty of a block doesn't match the turn of the signer
	ErrExtraValidators       = sdkerrors.Register(moduleName, 14, "non-sprint-end block contains extra validator list") // ErrExtraValidators is returned if non-sprint-end block contain validator data in their extra-data fields
	ErrInvalidSpanValidators = sdkerrors.Register(moduleName, 15, "invalid validator list on sprint end block")         // ErrInvalidSpanValidators is returned if a block contains an invalid list of validators (i.e. non divisible by 20 bytes)
	ErrInvalidProof          = sdkerrors.Register(moduleName, 16, "invalid proof")
)
