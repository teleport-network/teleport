package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/light"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	commitmenttypes "github.com/teleport-network/teleport/x/xibc/core/commitment/types"
	"github.com/teleport-network/teleport/x/xibc/exported"
)

var _ exported.ClientState = (*ClientState)(nil)

func (cs ClientState) ClientType() string {
	return exported.ETH
}

func (cs ClientState) GetLatestHeight() exported.Height {
	return cs.Header.Height
}

func (cs ClientState) Validate() error {
	return cs.Header.ValidateBasic()
}

func (m ClientState) CheckMsg(msg sdk.Msg) error {
	return nil
}

func (cs ClientState) GetDelayTime() uint64 {
	return cs.TimeDelay
}

func (cs ClientState) GetDelayBlock() uint64 {
	return cs.BlockDelay
}

func (cs ClientState) GetPrefix() exported.Prefix {
	return commitmenttypes.MerklePrefix{KeyPrefix: cs.ContractAddress}
}

func (cs ClientState) Initialize(
	ctx sdk.Context,
	cdc codec.BinaryCodec,
	store sdk.KVStore,
	state exported.ConsensusState,
) error {
	header := cs.Header
	headerBytes, err := cdc.MarshalInterface(&header)
	if err != nil {
		return sdkerrors.Wrap(ErrInvalidGenesisBlock, "marshal consensus to interface failed")
	}
	SetEthHeaderIndex(store, header, headerBytes)
	SetEthConsensusRoot(store, header.Height.RevisionHeight, header.ToEthHeader().Root, header.Hash())
	return nil
}

func (cs ClientState) UpgradeState(
	ctx sdk.Context,
	cdc codec.BinaryCodec,
	store sdk.KVStore,
	state exported.ConsensusState,
) error {
	header := cs.Header
	headerBytes, err := cdc.MarshalInterface(&header)
	if err != nil {
		return sdkerrors.Wrap(clienttypes.ErrUpgradeClient, "marshal consensus to interface failed")
	}
	SetEthHeaderIndex(store, header, headerBytes)
	SetEthConsensusRoot(store, header.Height.RevisionHeight, header.ToEthHeader().Root, header.Hash())
	return nil
}

func (cs ClientState) Status(
	ctx sdk.Context,
	store sdk.KVStore,
	cdc codec.BinaryCodec,
) exported.Status {
	onsState, err := GetConsensusState(store, cdc, cs.GetLatestHeight())
	if err != nil {
		return exported.Unknown
	}
	if onsState.Timestamp+cs.TrustingPeriod < uint64(ctx.BlockTime().Unix()) {
		return exported.Expired
	}
	return exported.Active
}

func (cs ClientState) ExportMetadata(store sdk.KVStore) []exported.GenesisMetadata {
	gm := make([]exported.GenesisMetadata, 0)
	callback := func(key, val []byte) bool {
		gm = append(gm, clienttypes.NewGenesisMetadata(key, val))
		return false
	}

	IteratorEthMetaDataByPrefix(store, KeyIndexEthHeaderPrefix, callback)
	IteratorEthMetaDataByPrefix(store, KeyMainRootPrefix, callback)

	if len(gm) == 0 {
		return nil
	}
	return gm
}

func (cs ClientState) VerifyPacketCommitment(
	ctx sdk.Context,
	store sdk.KVStore,
	cdc codec.BinaryCodec,
	height exported.Height,
	proof []byte,
	srcChain string,
	dstChain string,
	sequence uint64,
	commitment []byte,
) error {
	ethProof, consensusState, err := produceVerificationArgs(store, cdc, cs, height, proof)
	if err != nil {
		return err
	}

	// check delay period has passed
	delayBlock := cs.Header.Height.RevisionHeight - height.GetRevisionHeight()
	if delayBlock < cs.GetDelayBlock() {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidHeight,
			"delay block (%d) < client state delay block (%d)",
			delayBlock, cs.GetDelayBlock(),
		)
	}
	constructor := NewProofKeyConstructor(srcChain, dstChain, sequence)
	// verify that the provided commitment has been stored
	return verifyMerkleProof(ethProof, consensusState, cs.ContractAddress, commitment, constructor.GetPacketCommitmentProofKey())
}

func (cs ClientState) VerifyPacketAcknowledgement(
	ctx sdk.Context,
	store sdk.KVStore,
	cdc codec.BinaryCodec,
	height exported.Height,
	proof []byte,
	srcChain string,
	dstChain string,
	sequence uint64,
	ackBytes []byte,
) error {
	ethProof, consensusState, err := produceVerificationArgs(store, cdc, cs, height, proof)
	if err != nil {
		return err
	}

	delayBlock := cs.Header.Height.RevisionHeight - height.GetRevisionHeight()
	if delayBlock < cs.GetDelayBlock() {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidHeight,
			"delay block (%d) < client state delay block (%d)",
			delayBlock, cs.GetDelayBlock(),
		)
	}
	constructor := NewProofKeyConstructor(srcChain, dstChain, sequence)
	return verifyMerkleProof(ethProof, consensusState, cs.ContractAddress, ackBytes, constructor.GetAckProofKey())
}

// produceVerificationArgs performs the basic checks on the arguments that are
// shared between the verification functions and returns the unmarshal
// merkle proof, the consensus state and an error if one occurred.
func produceVerificationArgs(
	store sdk.KVStore,
	cdc codec.BinaryCodec,
	cs ClientState,
	height exported.Height,
	proof []byte,
) (
	merkleProof Proof,
	consensusState *ConsensusState,
	err error,
) {
	if cs.GetLatestHeight().LT(height) {
		return Proof{}, nil, sdkerrors.Wrapf(
			sdkerrors.ErrInvalidHeight,
			"client state height < proof height (%d < %d)",
			cs.GetLatestHeight(), height,
		)
	}

	if proof == nil {
		return Proof{}, nil, sdkerrors.Wrap(ErrInvalidProof, "proof cannot be empty")
	}

	if err = json.Unmarshal(proof, &merkleProof); err != nil {
		return Proof{}, nil, sdkerrors.Wrap(ErrInvalidProof, "failed to unmarshal proof into proof")
	}

	consensusState, err = GetConsensusState(store, cdc, height)
	if err != nil {
		return Proof{}, nil, err
	}

	return merkleProof, consensusState, nil
}

func verifyMerkleProof(
	ethProof Proof,
	consensusState *ConsensusState,
	contractAddr []byte,
	commitment []byte,
	proofKey []byte,
) error {
	// 1. prepare verify account
	nodeList := new(light.NodeList)

	for _, s := range ethProof.AccountProof {
		_ = nodeList.Put(nil, common.FromHex(s))
	}
	ns := nodeList.NodeSet()

	addr := common.FromHex(ethProof.Address)
	if !bytes.Equal(addr, contractAddr) {
		return fmt.Errorf(
			"verifyMerkleProof, contract address is error, proof address: %s, side chain address: %s",
			ethProof.Address, hex.EncodeToString(contractAddr),
		)
	}
	acctKey := crypto.Keccak256(addr)

	// 2. verify account proof
	root := common.BytesToHash(consensusState.Root)
	acctVal, err := trie.VerifyProof(root, acctKey, ns)
	if err != nil {
		return fmt.Errorf("verifyMerkleProof, verify account proof error:%s", err)
	}

	storageHash := common.HexToHash(ethProof.StorageHash)
	codeHash := common.HexToHash(ethProof.CodeHash)
	nonce := common.HexToHash(ethProof.Nonce).Big()
	balance := common.HexToHash(ethProof.Balance).Big()

	acct := &ProofAccount{
		Nonce:    nonce,
		Balance:  balance,
		Storage:  storageHash,
		Codehash: codeHash,
	}

	accRlp, err := rlp.EncodeToBytes(acct)
	if err != nil {
		return err
	}

	if !bytes.Equal(accRlp, acctVal) {
		return fmt.Errorf("verifyMerkleProof, verify account proof failed, wanted:%v, get:%v", accRlp, acctVal)
	}

	// 3.verify storage proof
	nodeList = new(light.NodeList)
	if len(ethProof.StorageProof) != 1 {
		return fmt.Errorf("verifyMerkleProof, invalid storage proof format")
	}

	sp := ethProof.StorageProof[0]

	if !bytes.Equal(common.HexToHash(sp.Key).Bytes(), proofKey) {
		return fmt.Errorf("verifyMerkleProof, storageKey is error, storage key: %s, Key path: %s", common.HexToHash(sp.Key), proofKey)
	}

	storageKey := crypto.Keccak256(common.HexToHash(sp.Key).Bytes())

	for _, prf := range sp.Proof {
		_ = nodeList.Put(nil, common.FromHex(prf))
	}

	ns = nodeList.NodeSet()
	val, err := trie.VerifyProof(storageHash, storageKey, ns)
	if err != nil {
		return fmt.Errorf("verifyMerkleProof, verify storage proof error:%s", err)
	}

	if !checkProofResult(val, commitment) {
		return fmt.Errorf("verifyMerkleProof, verify storage result failed")
	}
	return nil
}

func checkProofResult(result, value []byte) bool {
	var tempBytes []byte
	if err := rlp.DecodeBytes(result, &tempBytes); err != nil {
		return false
	}

	var s []byte
	for i := len(tempBytes); i < 32; i++ {
		s = append(s, 0)
	}
	s = append(s, tempBytes...)

	return bytes.Equal(s, value)
}
