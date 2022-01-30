package types

import (
	context "context"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/tharsis/ethermint/x/evm/types"

	"github.com/teleport-network/teleport/x/xibc/exported"
)

// ClientKeeper expected account XIBC client keeper
type ClientKeeper interface {
	GetClientState(ctx sdk.Context, chainName string) (exported.ClientState, bool)
	GetClientConsensusState(ctx sdk.Context, chainName string, height exported.Height) (exported.ConsensusState, bool)
	ClientStore(ctx sdk.Context, chainName string) sdk.KVStore
	GetChainName(ctx sdk.Context) string
}

// AccountKeeper defines the expected account keeper
type AccountKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
	GetSequence(sdk.Context, sdk.AccAddress) (uint64, error)
}

// EVMKeeper defines the expected EVM keeper interface used on xibc-transfer
type EVMKeeper interface {
	ChainID() *big.Int
	GetNonce(ctx sdk.Context, addr common.Address) uint64
	ApplyMessage(ctx sdk.Context, msg core.Message, tracer vm.Tracer, commit bool) (*types.MsgEthereumTxResponse, error)
	EthereumTx(goCtx context.Context, msg *types.MsgEthereumTx) (*types.MsgEthereumTxResponse, error)
}
