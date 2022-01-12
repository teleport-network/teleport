package types

import (
	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
)

// NewQueryPacketCommitmentResponse creates a new QueryPacketCommitmentResponse instance
func NewQueryPacketCommitmentResponse(
	commitment []byte, proof []byte, height clienttypes.Height,
) *QueryPacketCommitmentResponse {
	return &QueryPacketCommitmentResponse{
		Commitment:  commitment,
		Proof:       proof,
		ProofHeight: height,
	}
}

// NewQueryPacketReceiptResponse creates a new QueryPacketReceiptResponse instance
func NewQueryPacketReceiptResponse(
	recvd bool, proof []byte, height clienttypes.Height,
) *QueryPacketReceiptResponse {
	return &QueryPacketReceiptResponse{
		Received:    recvd,
		Proof:       proof,
		ProofHeight: height,
	}
}

// NewQueryPacketAcknowledgementResponse creates a new QueryPacketAcknowledgementResponse instance
func NewQueryPacketAcknowledgementResponse(
	acknowledgement []byte, proof []byte, height clienttypes.Height,
) *QueryPacketAcknowledgementResponse {
	return &QueryPacketAcknowledgementResponse{
		Acknowledgement: acknowledgement,
		Proof:           proof,
		ProofHeight:     height,
	}
}
