package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/teleport-network/teleport/x/xibc/core/host"
)

const (
	paramsIndex  = 205
	paramsLenght = 32
)

type ProofKeyConstructor struct {
	srcChain string
	dstChain string
	sequence uint64
}

func NewProofKeyConstructor(srcChain string, dstChain string, sequence uint64) ProofKeyConstructor {
	return ProofKeyConstructor{
		srcChain: srcChain,
		dstChain: dstChain,
		sequence: sequence,
	}
}

func (k ProofKeyConstructor) GetPacketCommitmentProofKey() []byte {
	hash := crypto.Keccak256Hash(
		host.PacketCommitmentKey(k.srcChain, k.dstChain, k.sequence),
		common.LeftPadBytes(big.NewInt(paramsIndex).Bytes(), paramsLenght),
	)
	return hash.Bytes()
}

func (k ProofKeyConstructor) GetAckProofKey() []byte {
	hash := crypto.Keccak256Hash(
		host.PacketAcknowledgementKey(k.srcChain, k.dstChain, k.sequence),
		common.LeftPadBytes(big.NewInt(paramsIndex).Bytes(), paramsLenght),
	)
	return hash.Bytes()
}
