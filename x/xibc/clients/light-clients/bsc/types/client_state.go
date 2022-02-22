package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

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

// NewClientState creates a new ClientState instance
func NewClientState(
	header Header,
	chainID uint64,
	epoch uint64,
	blockInteval uint64,
	validators [][]byte,
	contractAddress []byte,
	trustingPeriod uint64,
) *ClientState {
	return &ClientState{
		Header:          header,
		ChainId:         chainID,
		Epoch:           epoch,
		BlockInteval:    blockInteval,
		Validators:      validators,
		ContractAddress: contractAddress,
		TrustingPeriod:  trustingPeriod,
	}
}

func (m ClientState) ClientType() string {
	return exported.BSC
}

func (m ClientState) GetLatestHeight() exported.Height {
	return m.Header.Height
}

func (m ClientState) Validate() error {
	return m.Header.ValidateBasic()
}

func (m ClientState) CheckMsg(msg sdk.Msg) error {
	return nil
}

func (m ClientState) GetDelayTime() uint64 {
	return uint64(2*len(m.Validators)/3+1) * m.BlockInteval
}

func (m ClientState) GetDelayBlock() uint64 {
	return uint64(2*len(m.Validators)/3 + 1)
}

func (m ClientState) GetPrefix() exported.Prefix {
	return commitmenttypes.MerklePrefix{}
}

func (m ClientState) Initialize(
	ctx sdk.Context,
	cdc codec.BinaryCodec,
	store sdk.KVStore,
	state exported.ConsensusState,
) error {
	if m.Header.Height.RevisionHeight%m.Epoch != 0 {
		return sdkerrors.Wrap(ErrInvalidGenesisBlock, "header")
	}
	// Resolve the authorization key and check against validators
	signer, err := ecrecover(m.Header, big.NewInt(int64(m.ChainId)))
	if err != nil {
		return err
	}
	if signer != common.BytesToAddress(m.Header.Coinbase) {
		return sdkerrors.Wrap(ErrCoinBaseMisMatch, "header.Coinbase")
	}
	SetSigner(store, Signer{
		Height:    m.Header.Height,
		Validator: signer.Bytes(),
	})
	validators, err := ParseValidators(m.Header.Extra)
	if err != nil {
		return err
	}
	SetPendingValidators(store, cdc, validators)
	return nil
}

func (m ClientState) UpgradeState(
	ctx sdk.Context,
	cdc codec.BinaryCodec,
	store sdk.KVStore,
	state exported.ConsensusState,
) error {
	if m.Header.Height.RevisionHeight%m.Epoch != 0 {
		return sdkerrors.Wrap(ErrInvalidGenesisBlock, "header")
	}
	// Check the earliest consensus state to see if it is expired, if so then set the prune height
	// so that we can delete consensus state and all associated metadata.
	var (
		pruneHeight exported.Height
		pruneError  error
	)
	pruneCb := func(height exported.Height) bool {
		consState, err := GetConsensusState(store, cdc, height)
		// this error should never occur
		if err != nil {
			pruneError = err
			return true
		}
		blockTime := uint64(ctx.BlockTime().Unix())
		if consState.Timestamp+m.TrustingPeriod < blockTime {
			pruneHeight = height
		}
		return true
	}
	IterateConsensusStateAscending(store, pruneCb)
	if pruneError != nil {
		return pruneError
	}
	// if pruneHeight is set, delete consensus state and metadata
	if pruneHeight != nil {
		deleteConsensusState(store, pruneHeight)
	}
	// Delete all signer
	if err := DeleteAllSigner(store); err != nil {
		return err
	}

	// Resolve the authorization key and check against validators
	signer, err := ecrecover(m.Header, big.NewInt(int64(m.ChainId)))
	if err != nil {
		return err
	}
	if signer != common.BytesToAddress(m.Header.Coinbase) {
		return sdkerrors.Wrap(ErrCoinBaseMisMatch, "header.Coinbase")
	}
	SetSigner(store, Signer{
		Height:    m.Header.Height,
		Validator: signer.Bytes(),
	})
	validators, err := ParseValidators(m.Header.Extra)
	if err != nil {
		return err
	}
	SetPendingValidators(store, cdc, validators)
	return nil
}

func (m ClientState) Status(
	ctx sdk.Context,
	store sdk.KVStore,
	cdc codec.BinaryCodec,
) exported.Status {
	onsState, err := GetConsensusState(store, cdc, m.GetLatestHeight())
	if err != nil {
		return exported.Unknown
	}
	if onsState.Timestamp+m.TrustingPeriod < uint64(ctx.BlockTime().Unix()) {
		return exported.Expired
	}
	return exported.Active
}

// ExportMetadata exports RecentSingers and PendingValidators
func (m ClientState) ExportMetadata(store sdk.KVStore) []exported.GenesisMetadata {
	gm := make([]exported.GenesisMetadata, 0)
	callback := func(key, val []byte) bool {
		gm = append(gm, clienttypes.NewGenesisMetadata(key, val))
		return false
	}

	IteratorTraversal(store, PrefixKeyRecentSingers, callback)
	IteratorTraversal(store, PrefixPendingValidators, callback)

	if len(gm) == 0 {
		return nil
	}
	return gm
}

func (m ClientState) VerifyPacketCommitment(
	ctx sdk.Context,
	store sdk.KVStore,
	cdc codec.BinaryCodec,
	height exported.Height,
	proof []byte,
	sourceChain, destChain string,
	sequence uint64,
	commitment []byte,
) error {
	bscProof, consensusState, err := produceVerificationArgs(store, cdc, m, height, proof)
	if err != nil {
		return err
	}

	// check delay period has passed
	delayBlock := m.Header.Height.RevisionHeight - height.GetRevisionHeight()
	if delayBlock < m.GetDelayBlock() {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidHeight,
			"delay block (%d) < client state delay block (%d)",
			delayBlock, m.GetDelayBlock(),
		)
	}

	constructor := NewProofKeyConstructor(sourceChain, destChain, sequence)

	// verify that the provided commitment has been stored
	return verifyMerkleProof(bscProof, consensusState, m.ContractAddress, commitment, constructor.GetPacketCommitmentProofKey())
}

func (m ClientState) VerifyPacketAcknowledgement(
	ctx sdk.Context,
	store sdk.KVStore,
	cdc codec.BinaryCodec,
	height exported.Height,
	proof []byte,
	sourceChain, destChain string,
	sequence uint64,
	ackBytes []byte,
) error {
	bscProof, consensusState, err := produceVerificationArgs(store, cdc, m, height, proof)
	if err != nil {
		return err
	}

	delayBlock := m.Header.Height.RevisionHeight - height.GetRevisionHeight()
	if delayBlock < m.GetDelayBlock() {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidHeight,
			"delay block (%d) < client state delay block (%d)",
			delayBlock, m.GetDelayBlock(),
		)
	}
	constructor := NewProofKeyConstructor(sourceChain, destChain, sequence)
	return verifyMerkleProof(bscProof, consensusState, m.ContractAddress, ackBytes, constructor.GetAckProofKey())
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
	bscProof Proof,
	consensusState *ConsensusState,
	contractAddr []byte,
	commitment []byte,
	proofKey []byte,
) error {
	//1. prepare verify account
	nodeList := new(light.NodeList)

	for _, s := range bscProof.AccountProof {
		_ = nodeList.Put(nil, common.FromHex(s))
	}
	ns := nodeList.NodeSet()

	addr := common.FromHex(bscProof.Address)
	if !bytes.Equal(addr, contractAddr) {
		return fmt.Errorf(
			"verifyMerkleProof, contract address is error, proof address: %s, side chain address: %s",
			bscProof.Address, hex.EncodeToString(contractAddr),
		)
	}
	acctKey := crypto.Keccak256(addr)

	//2. verify account proof
	root := common.BytesToHash(consensusState.Root)
	acctVal, err := trie.VerifyProof(root, acctKey, ns)
	if err != nil {
		return fmt.Errorf("verifyMerkleProof, verify account proof error:%s", err)
	}

	storageHash := common.HexToHash(bscProof.StorageHash)
	codeHash := common.HexToHash(bscProof.CodeHash)
	nonce := common.HexToHash(bscProof.Nonce).Big()
	balance := common.HexToHash(bscProof.Balance).Big()

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

	//3.verify storage proof
	nodeList = new(light.NodeList)
	if len(bscProof.StorageProof) != 1 {
		return fmt.Errorf("verifyMerkleProof, invalid storage proof format")
	}

	sp := bscProof.StorageProof[0]

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
	// TODO
	// hash := crypto.Keccak256(value)
	return bytes.Equal(s, value)
}
