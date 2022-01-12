package utils

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	ibcclient "github.com/teleport-network/teleport/x/xibc/client"
	host "github.com/teleport-network/teleport/x/xibc/core/host"
	"github.com/teleport-network/teleport/x/xibc/core/packet/types"
)

// QueryPacketCommitment returns a packet commitment.
// If prove is true, it performs an ABCI store query in order to retrieve the merkle proof. Otherwise,
// it uses the gRPC query client.
func QueryPacketCommitment(
	clientCtx client.Context, sourceChain, destChain string, sequence uint64, prove bool,
) (
	*types.QueryPacketCommitmentResponse, error,
) {
	if prove {
		return queryPacketCommitmentABCI(clientCtx, sourceChain, destChain, sequence)
	}

	queryClient := types.NewQueryClient(clientCtx)
	req := &types.QueryPacketCommitmentRequest{
		SourceChain: sourceChain,
		DestChain:   destChain,
		Sequence:    sequence,
	}

	return queryClient.PacketCommitment(context.Background(), req)
}

func queryPacketCommitmentABCI(
	clientCtx client.Context, sourceChain, destChain string, sequence uint64,
) (
	*types.QueryPacketCommitmentResponse, error,
) {
	key := host.PacketCommitmentKey(sourceChain, destChain, sequence)

	value, proofBz, proofHeight, err := ibcclient.QueryTendermintProof(clientCtx, key)
	if err != nil {
		return nil, err
	}

	// check if packet commitment exists
	if len(value) == 0 {
		return nil, sdkerrors.Wrapf(
			types.ErrPacketCommitmentNotFound,
			"source chain name  (%s), dest chain name (%s), sequence (%d)",
			sourceChain, destChain, sequence,
		)
	}

	return types.NewQueryPacketCommitmentResponse(value, proofBz, proofHeight), nil
}

// QueryPacketReceipt returns data about a packet receipt.
// If prove is true, it performs an ABCI store query in order to retrieve the merkle proof. Otherwise,
// it uses the gRPC query client.
func QueryPacketReceipt(
	clientCtx client.Context, sourceChain, destChain string, sequence uint64, prove bool,
) (
	*types.QueryPacketReceiptResponse, error,
) {
	if prove {
		return queryPacketReceiptABCI(clientCtx, sourceChain, destChain, sequence)
	}

	queryClient := types.NewQueryClient(clientCtx)
	req := &types.QueryPacketReceiptRequest{
		SourceChain: sourceChain,
		DestChain:   destChain,
		Sequence:    sequence,
	}

	return queryClient.PacketReceipt(context.Background(), req)
}

func queryPacketReceiptABCI(
	clientCtx client.Context, sourceChain, destChain string, sequence uint64,
) (
	*types.QueryPacketReceiptResponse, error,
) {
	key := host.PacketReceiptKey(sourceChain, destChain, sequence)

	value, proofBz, proofHeight, err := ibcclient.QueryTendermintProof(clientCtx, key)
	if err != nil {
		return nil, err
	}

	return types.NewQueryPacketReceiptResponse(value != nil, proofBz, proofHeight), nil
}

// QueryPacketAcknowledgement returns the data about a packet acknowledgement.
// If prove is true, it performs an ABCI store query in order to retrieve the merkle proof. Otherwise,
// it uses the gRPC query client
func QueryPacketAcknowledgement(
	clientCtx client.Context, sourceChain, destChain string, sequence uint64, prove bool,
) (
	*types.QueryPacketAcknowledgementResponse, error,
) {
	if prove {
		return queryPacketAcknowledgementABCI(clientCtx, sourceChain, destChain, sequence)
	}

	queryClient := types.NewQueryClient(clientCtx)
	req := &types.QueryPacketAcknowledgementRequest{
		SourceChain: sourceChain,
		DestChain:   destChain,
		Sequence:    sequence,
	}

	return queryClient.PacketAcknowledgement(context.Background(), req)
}

func queryPacketAcknowledgementABCI(
	clientCtx client.Context, sourceChain, destChain string, sequence uint64,
) (
	*types.QueryPacketAcknowledgementResponse, error,
) {
	key := host.PacketAcknowledgementKey(sourceChain, destChain, sequence)

	value, proofBz, proofHeight, err := ibcclient.QueryTendermintProof(clientCtx, key)
	if err != nil {
		return nil, err
	}

	if len(value) == 0 {
		return nil, sdkerrors.Wrapf(
			types.ErrInvalidAcknowledgement,
			"source chain name  (%s), dest chain name (%s), sequence (%d)",
			sourceChain, destChain, sequence,
		)
	}

	return types.NewQueryPacketAcknowledgementResponse(value, proofBz, proofHeight), nil
}
