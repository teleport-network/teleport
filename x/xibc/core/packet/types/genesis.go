package types

import (
	"errors"
	"fmt"

	"github.com/teleport-network/teleport/x/xibc/core/host"
)

// NewPacketState creates a new PacketState instance.
func NewPacketState(srcChain, dstChain string, seq uint64, data []byte) PacketState {
	return PacketState{
		SrcChain: srcChain,
		DstChain: dstChain,
		Sequence: seq,
		Data:     data,
	}
}

// Validate performs basic validation of fields returning an error upon any
// failure.
func (pa PacketState) Validate() error {
	if pa.Data == nil {
		return errors.New("data bytes cannot be nil")
	}
	return validateGenFields(pa.SrcChain, pa.DstChain, pa.Sequence)
}

// NewPacketSequence creates a new PacketSequences instance.
func NewPacketSequence(srcChain string, dstChain string, seq uint64) PacketSequence {
	return PacketSequence{
		SrcChain: srcChain,
		DstChain: dstChain,
		Sequence: seq,
	}
}

// Validate performs basic validation of fields returning an error upon any
// failure.
func (ps PacketSequence) Validate() error {
	return validateGenFields(ps.SrcChain, ps.DstChain, ps.Sequence)
}

// NewGenesisState creates a GenesisState instance.
func NewGenesisState(
	acks, commitments, receipts []PacketState,
	sendSeqs, recvSeqs, ackSeqs []PacketSequence,
) GenesisState {
	return GenesisState{
		Acknowledgements: acks,
		Commitments:      commitments,
		Receipts:         receipts,
		SendSequences:    sendSeqs,
		RecvSequences:    recvSeqs,
		AckSequences:     ackSeqs,
	}
}

// DefaultGenesisState returns the xibc packet submodule's default genesis state.
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Acknowledgements: []PacketState{},
		Receipts:         []PacketState{},
		Commitments:      []PacketState{},
		SendSequences:    []PacketSequence{},
		RecvSequences:    []PacketSequence{},
		AckSequences:     []PacketSequence{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	for i, ack := range gs.Acknowledgements {
		if err := ack.Validate(); err != nil {
			return fmt.Errorf("invalid acknowledgement %v ack index %d: %w", ack, i, err)
		}
		if len(ack.Data) == 0 {
			return fmt.Errorf("invalid acknowledgement %v ack index %d: data bytes cannot be empty", ack, i)
		}
	}

	for i, receipt := range gs.Receipts {
		if err := receipt.Validate(); err != nil {
			return fmt.Errorf("invalid acknowledgement %v ack index %d: %w", receipt, i, err)
		}
	}

	for i, commitment := range gs.Commitments {
		if err := commitment.Validate(); err != nil {
			return fmt.Errorf("invalid commitment %v index %d: %w", commitment, i, err)
		}
		if len(commitment.Data) == 0 {
			return fmt.Errorf("invalid acknowledgement %v ack index %d: data bytes cannot be empty", commitment, i)
		}
	}

	for i, ss := range gs.SendSequences {
		if err := ss.Validate(); err != nil {
			return fmt.Errorf("invalid send sequence %v index %d: %w", ss, i, err)
		}
	}

	for i, rs := range gs.RecvSequences {
		if err := rs.Validate(); err != nil {
			return fmt.Errorf("invalid receive sequence %v index %d: %w", rs, i, err)
		}
	}

	for i, as := range gs.AckSequences {
		if err := as.Validate(); err != nil {
			return fmt.Errorf("invalid acknowledgement sequence %v index %d: %w", as, i, err)
		}
	}

	return nil
}

func validateGenFields(srcChain string, dstChain string, sequence uint64) error {
	if err := host.SrcChainValidator(srcChain); err != nil {
		return fmt.Errorf("invalid src chain ID: %w", err)
	}
	if err := host.DstChainValidator(dstChain); err != nil {
		return fmt.Errorf("invalid dst chain ID: %w", err)
	}
	if sequence == 0 {
		return errors.New("sequence cannot be 0")
	}
	return nil
}
