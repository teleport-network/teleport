package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/teleport-network/teleport/x/xibc/exported"
)

// ClientKeeper expected account XIBC client keeper
type ClientKeeper interface {
	GetClientState(ctx sdk.Context, chainName string) (exported.ClientState, bool)
	GetClientConsensusState(ctx sdk.Context, chainName string, height exported.Height) (exported.ConsensusState, bool)
	ClientStore(ctx sdk.Context, chainName string) sdk.KVStore
	GetChainName(ctx sdk.Context) string
}
