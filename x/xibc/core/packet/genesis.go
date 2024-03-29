package packet

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/teleport-network/teleport/x/xibc/core/packet/keeper"
	"github.com/teleport-network/teleport/x/xibc/core/packet/types"
)

// InitGenesis initializes the xibc packet submodule's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, gs types.GenesisState) {
	// ensure xibc transfer module account is set on genesis
	if acc := k.GetModuleAccount(ctx); acc == nil {
		panic("the xibc packet module account has not been set")
	}

	for _, ack := range gs.Acknowledgements {
		k.SetPacketAcknowledgement(ctx, ack.SrcChain, ack.DstChain, ack.Sequence, ack.Data)
	}
	for _, commitment := range gs.Commitments {
		k.SetPacketCommitment(ctx, commitment.SrcChain, commitment.DstChain, commitment.Sequence, commitment.Data)
	}
	for _, receipt := range gs.Receipts {
		k.SetPacketReceipt(ctx, receipt.SrcChain, receipt.DstChain, receipt.Sequence)
	}
	for _, ss := range gs.SendSequences {
		k.SetNextSequenceSend(ctx, ss.SrcChain, ss.DstChain, ss.Sequence)
	}
}

// ExportGenesis returns the xibc packet submodule's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {
	return types.GenesisState{
		Acknowledgements: k.GetAllPacketAcks(ctx),
		Commitments:      k.GetAllPacketCommitments(ctx),
		Receipts:         k.GetAllPacketReceipts(ctx),
		SendSequences:    k.GetAllPacketSendSeqs(ctx),
	}
}
