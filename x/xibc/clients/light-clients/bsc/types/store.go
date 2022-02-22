package types

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	"github.com/teleport-network/teleport/x/xibc/core/host"
	"github.com/teleport-network/teleport/x/xibc/exported"
)

var (
	PrefixKeyRecentSingers  = "recentSingers"
	PrefixPendingValidators = "pendingValidators"
)

func GetIterator(store sdk.KVStore, keyType string) types.Iterator {
	iterator := sdk.KVStorePrefixIterator(store, []byte(keyType))
	return iterator
}

func IteratorTraversal(store sdk.KVStore, keyType string, cb func(key, val []byte) bool) {
	iterator := GetIterator(store, keyType)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		if cb(iterator.Key(), iterator.Value()) {
			break
		}
	}
}

// GetConsensusState retrieves the consensus state from the client prefixed
// store. An error is returned if the consensus state does not exist.
func GetConsensusState(store sdk.KVStore, cdc codec.BinaryCodec, height exported.Height) (*ConsensusState, error) {
	bz := store.Get(host.ConsensusStateKey(height))
	if bz == nil {
		return nil, sdkerrors.Wrapf(
			clienttypes.ErrConsensusStateNotFound,
			"consensus state does not exist for height %s",
			height,
		)
	}

	consensusStateI, err := clienttypes.UnmarshalConsensusState(cdc, bz)
	if err != nil {
		return nil, sdkerrors.Wrapf(clienttypes.ErrInvalidConsensus, "unmarshal error: %v", err)
	}

	consensusState, ok := consensusStateI.(*ConsensusState)
	if !ok {
		return nil, sdkerrors.Wrapf(
			clienttypes.ErrInvalidConsensus,
			"invalid consensus type %T, expected %T",
			consensusState, &ConsensusState{},
		)
	}

	return consensusState, nil
}

func SetSigner(store sdk.KVStore, signer Signer) {
	store.Set(keyRecentSinger(signer), signer.Validator)
}

func DeleteSigner(store sdk.KVStore, height clienttypes.Height) {
	keyBz := []byte(fmt.Sprintf("%s/%s", PrefixKeyRecentSingers, height))
	store.Delete(keyBz)
}

func DeleteAllSigner(store sdk.KVStore) error {
	iterator := sdk.KVStorePrefixIterator(store, []byte(PrefixKeyRecentSingers))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		iterKey := iterator.Key()
		keys := strings.Split(string(iterKey), "/")
		height, err := clienttypes.ParseHeight(keys[1])
		if err != nil {
			return err
		}
		DeleteSigner(store, height)
	}
	return nil
}

// GetRecentSigners retrieves the recent singer list from the client prefixed
func GetRecentSigners(store sdk.KVStore) (recentSingers []Signer, err error) {
	iterator := sdk.KVStorePrefixIterator(store, []byte(PrefixKeyRecentSingers))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		iterKey := iterator.Key()
		keys := strings.Split(string(iterKey), "/")
		height, err := clienttypes.ParseHeight(keys[1])
		if err != nil {
			return nil, err
		}
		recentSingers = append(recentSingers, Signer{
			Height:    height,
			Validator: iterator.Value(),
		})
	}
	return
}

// SetPendingValidators sets the validators to be updated in the client prefixed store
func SetPendingValidators(store sdk.KVStore, cdc codec.BinaryCodec, validators [][]byte,
) {
	validatorSet := ValidatorSet{
		Validators: validators,
	}
	bz := cdc.MustMarshal(&validatorSet)
	store.Set([]byte(PrefixPendingValidators), bz)
}

// GetPendingValidators retrieves the validators to be updated from the client prefixed store
func GetPendingValidators(cdc codec.BinaryCodec, store sdk.KVStore) ValidatorSet {
	bz := store.Get([]byte(PrefixPendingValidators))

	var validatorSet ValidatorSet
	cdc.MustUnmarshal(bz, &validatorSet)
	return validatorSet
}

func keyRecentSinger(singer Signer) []byte {
	return []byte(fmt.Sprintf("%s/%s", PrefixKeyRecentSingers, singer.Height))
}

// IterateConsensusStateAscending iterates through the consensus states in ascending order. It calls the provided
// callback on each height, until stop=true is returned.
func IterateConsensusStateAscending(clientStore sdk.KVStore,
	cb func(height exported.Height) (stop bool)) {
	iterator := sdk.KVStorePrefixIterator(clientStore, []byte(host.KeyConsensusStatePrefix))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		keySplit := strings.Split(string(key), "/")
		// processed time key in prefix store has format: "consensusStates/<height>"
		if len(keySplit) != 2 {
			// ignore all not consensus state keys
			continue
		}
		height := GetHeightFromIterationKey(key)
		if cb(height) {
			return
		}
	}
}

// GetHeightFromIterationKey takes an iteration key and returns the height that it references
func GetHeightFromIterationKey(iterKey []byte) exported.Height {
	bigEndianBytes := iterKey[len([]byte(host.KeyConsensusStatePrefix+"/")):]
	revisionBytes := bigEndianBytes[0:8]
	heightBytes := bigEndianBytes[8:]
	revision := sdk.BigEndianToUint64(revisionBytes)
	height := sdk.BigEndianToUint64(heightBytes)
	return clienttypes.NewHeight(revision, height)
}

// deleteConsensusState deletes the consensus state at the given height
func deleteConsensusState(clientStore sdk.KVStore, height exported.Height) {
	key := host.ConsensusStateKey(height)
	clientStore.Delete(key)
}
