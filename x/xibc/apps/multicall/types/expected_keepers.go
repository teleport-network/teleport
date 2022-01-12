package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"

	"github.com/teleport-network/teleport/x/xibc/exported"
)

// PacketKeeper defines the expected packet keeper
type PacketKeeper interface {
	GetNextSequenceSend(ctx sdk.Context, sourceChain, destChain string) uint64
	SendPacket(ctx sdk.Context, packet exported.PacketI) error
}

// ClientKeeper defines the expected client keeper
type ClientKeeper interface {
	GetChainName(ctx sdk.Context) string
}

// AggregateKeeper defines the expected Aggregate keeper
type AggregateKeeper interface {
	QueryERC20Trace(ctx sdk.Context, contract common.Address, originChain string) (string, *big.Int, bool, error)
}
