package types

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	"github.com/teleport-network/teleport/x/xibc/core/host"
	"github.com/teleport-network/teleport/x/xibc/exported"
)

func (cs ClientState) CheckHeaderAndUpdateState(
	ctx sdk.Context,
	cdc codec.BinaryCodec,
	store sdk.KVStore,
	header exported.Header,
) (
	exported.ClientState,
	exported.ConsensusState,
	error,
) {
	ethHeader, ok := header.(*Header)
	if !ok {
		return nil, nil, sdkerrors.Wrapf(clienttypes.ErrInvalidHeader, "expected type %T, got %T", &Header{}, header)
	}
	height := cs.GetLatestHeight()
	// get consensus state from clientStore
	ethConsState, err := GetConsensusState(store, cdc, height)
	if err != nil {
		return nil, nil, sdkerrors.Wrapf(
			err, "could not get consensus state from clientStore at TrustedHeight: %s,please upgrade", cs.GetLatestHeight(),
		)
	}
	if err := checkValidity(ctx, cdc, store, &cs, ethConsState, *ethHeader); err != nil {
		return nil, nil, err
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
			// this return just for get out of the func
			return true
		}
		blockTime := uint64(ctx.BlockTime().Unix())
		if consState.Timestamp+cs.TrustingPeriod < blockTime {
			pruneHeight = height
		}
		return true
	}
	IterateConsensusStateAscending(store, pruneCb)
	if pruneError != nil {
		return nil, nil, pruneError
	}
	// if pruneHeight is set, delete consensus state and metadata
	if pruneHeight != nil {
		if err = deleteConsensusStateAndIndexHeader(cdc, store, pruneHeight); err != nil {
			return nil, nil, err
		}
	}
	newClientState, consensusState, err := update(ctx, cdc, store, &cs, ethHeader)
	if err != nil {
		return nil, nil, err
	}
	//Check the bifurcation
	if !bytes.Equal(cs.Header.Hash().Bytes(), ethHeader.ParentHash) {
		if err = cs.RestrictChain(cdc, store, *ethHeader); err != nil {
			return nil, nil, err
		}
	}
	cs.Header = *ethHeader
	newClientState.Header = *ethHeader

	return newClientState, consensusState, nil
}

// checkValidity checks if the eth header is valid.
func checkValidity(
	ctx sdk.Context,
	cdc codec.BinaryCodec,
	store sdk.KVStore,
	clientState *ClientState,
	consState *ConsensusState,
	header Header,
) error {
	if err := header.ValidateBasic(); err != nil {
		return err
	}
	if err := verifyHeader(ctx, cdc, store, clientState, header); err != nil {
		return err
	}

	if clientState.ChainId != 4 {
		// Ensure that the header's extra-data section is of a reasonable size
		if uint64(len(header.Extra)) > params.MaximumExtraDataSize {
			return fmt.Errorf("extra-data too long: %d > %d", len(header.Extra), params.MaximumExtraDataSize)
		}
		if err := VerifyCascadingFields(header); err != nil {
			return err
		}
	}
	return nil
}

// update the RecentSingers and the ConsensusState.
func update(
	ctx sdk.Context,
	cdc codec.BinaryCodec,
	store sdk.KVStore,
	clientState *ClientState,
	header *Header,
) (
	*ClientState,
	*ConsensusState,
	error,
) {
	cs := &ConsensusState{
		Timestamp: header.Time,
		Height:    header.Height,
		Root:      header.Root,
	}
	headerInterface, err := cdc.MarshalInterface(header)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(ErrInvalidGenesisBlock, "marshal consensus to interface failed")
	}
	SetEthHeaderIndex(store, *header, headerInterface)
	SetEthConsensusRoot(store, header.GetHeight().GetRevisionHeight(), header.ToEthHeader().Root, header.Hash())
	return clientState, cs, nil
}

func (m ClientState) RestrictChain(cdc codec.BinaryCodec, store sdk.KVStore, new Header) error {
	si, ti := m.Header.Height, new.Height
	var err error
	current := m.Header
	//si > ti
	if si.RevisionHeight > ti.RevisionHeight {
		ConsensusTmp := store.Get(host.ConsensusStateKey(ti))
		if ConsensusTmp == nil {
			return sdkerrors.Wrapf(
				clienttypes.ErrInvalidConsensus, "can not find consensus state for height %s in RestrictChain", ti)
		}
		var tiConsensus exported.ConsensusState
		if err = cdc.UnmarshalInterface(ConsensusTmp, &tiConsensus); err != nil {
			return sdkerrors.Wrapf(ErrUnmarshalInterface, "can not unmarshal ConsensusState interface in RestrictChain ")

		}
		tmpConsensus, ok := tiConsensus.(*ConsensusState)
		if !ok {
			return sdkerrors.Wrapf(
				clienttypes.ErrInvalidConsensus, "can not find consensus state for height %s in RestrictChain", ti)
		}
		root := tmpConsensus.Root
		headerIndexKey := GetHeaderIndexKeyByEthConsensusRoot(store, common.BytesToHash(root), ti.GetRevisionHeight())
		currentBytes := store.Get(headerIndexKey)
		if currentBytes == nil {
			return sdkerrors.Wrapf(
				clienttypes.ErrInvalidConsensus, "can not find Header for height %s in RestrictChain", ti)
		}
		var currentHeaderInterface exported.Header
		if err = cdc.UnmarshalInterface(currentBytes, &currentHeaderInterface); err != nil {
			return sdkerrors.Wrapf(ErrUnmarshalInterface, "can not unmarshal ConsensusState interface in RestrictChain ")

		}
		currentTmp, ok := currentHeaderInterface.(*Header)
		if !ok {
			return sdkerrors.Wrapf(
				clienttypes.ErrInvalidConsensus, "can not find consensus state for height %s in RestrictChain", ti)
		}
		current = *currentTmp
		si = ti
	}
	newHashes := make([]common.Hash, 0)

	for ti.RevisionHeight > si.RevisionHeight {
		newHashes = append(newHashes, new.Hash())
		newTmp := GetParentHeaderFromIndex(store, new)
		if newTmp == nil {
			return sdkerrors.Wrapf(
				clienttypes.ErrInvalidConsensus, "can not find consensus state for hash %s in RestrictChain in RestrictChain", new.ToEthHeader().ParentHash,
			)
		}
		var currently exported.Header
		if err := cdc.UnmarshalInterface(newTmp, &currently); err != nil {
			return sdkerrors.Wrapf(ErrUnmarshalInterface, "can not unmarshal ConsensusState interface in RestrictChain ")
		}
		tmpConsensus, ok := currently.(*Header)
		if !ok {
			return sdkerrors.Wrapf(
				clienttypes.ErrInvalidConsensus, "can not find consensus state for hash %s in RestrictChain", new.ToEthHeader().ParentHash)
		}
		new = *tmpConsensus
		ti.RevisionHeight--
	}
	// si.parent != ti.parent
	for !bytes.Equal(current.ParentHash, new.ParentHash) {
		newHashes = append(newHashes, new.Hash())
		newTmp := GetParentHeaderFromIndex(store, new)
		if newTmp == nil {
			return sdkerrors.Wrapf(
				clienttypes.ErrInvalidConsensus, "can not find consensus state for hash %s in RestrictChain", new.ToEthHeader().ParentHash)
		}
		var currently exported.Header
		if err = cdc.UnmarshalInterface(newTmp, &currently); err != nil {
			return sdkerrors.Wrapf(ErrUnmarshalInterface, "can not unmarshal ConsensusState interface in RestrictChain ")
		}
		tmpConsensus, ok := currently.(*Header)
		if !ok {
			return sdkerrors.Wrapf(
				clienttypes.ErrInvalidConsensus, "can not find consensus state for hash %s in RestrictChain", new.ToEthHeader().ParentHash)
		}
		new = *tmpConsensus
		ti.RevisionHeight--
		si.RevisionHeight--
		currentTmp := GetParentHeaderFromIndex(store, current)
		if currentTmp == nil {
			return sdkerrors.Wrapf(
				clienttypes.ErrInvalidConsensus, "can not find consensus state for height %s in RestrictChain", si)
		}
		if err = cdc.UnmarshalInterface(currentTmp, &currently); err != nil {
			return sdkerrors.Wrapf(ErrUnmarshalInterface, "can not unmarshal ConsensusState interface in RestrictChain ")
		}
		tmpConsensus = currently.(*Header)
		if !ok {
			return sdkerrors.Wrapf(
				clienttypes.ErrInvalidConsensus, "can not  consensus state for height %s in RestrictChain", si)
		}
		current = *tmpConsensus
	}
	for i := len(newHashes) - 1; i >= 0; i-- {
		newTmp := store.Get(EthHeaderIndexKey(newHashes[i], ti.GetRevisionHeight()))
		if newTmp == nil {
			return sdkerrors.Wrapf(
				clienttypes.ErrInvalidConsensus, "can not find consensus state for hash %s in RestrictChain", newHashes[i])

		}
		var currently exported.Header
		if err := cdc.UnmarshalInterface(newTmp, &currently); err != nil {
			return sdkerrors.Wrapf(ErrUnmarshalInterface, "can not unmarshal ConsensusState interface in RestrictChain ")
		}
		tmpHeader, ok := currently.(*Header)
		if !ok {
			return sdkerrors.Wrapf(
				clienttypes.ErrInvalidConsensus, "can not find consensus state for hash %s in RestrictChain", new.ToEthHeader().ParentHash)
		}
		consensusState := &ConsensusState{
			Timestamp: tmpHeader.Time,
			Height:    tmpHeader.Height,
			Root:      tmpHeader.Root[:],
		}
		consensusStateBytes, err := cdc.MarshalInterface(consensusState)
		if err != nil {
			return sdkerrors.Wrap(ErrInvalidGenesisBlock, "marshal consensus to byte failed")
		}
		// set main_chain
		store.Set(host.ConsensusStateKey(ti), consensusStateBytes)
		ti.RevisionHeight++
	}
	return err
}
