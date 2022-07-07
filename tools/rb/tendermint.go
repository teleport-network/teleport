package rb

import (
	"errors"
	"fmt"

	"github.com/bitdao-io/bitnetwork/tools/common"

	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/privval"
	tmstate "github.com/tendermint/tendermint/proto/tendermint/state"
	tmversion "github.com/tendermint/tendermint/proto/tendermint/version"
	"github.com/tendermint/tendermint/state"
	"github.com/tendermint/tendermint/store"
	"github.com/tendermint/tendermint/version"
	dbm "github.com/tendermint/tm-db"
)

func RollbackBlockAndState(home string, height int64, config *cfg.Config) error {
	// rollback tendermint state
	blockStoreDB, err := common.OpenBlockStoreDB(home)
	if err != nil {
		panic(err)
	}

	stateDB, err := common.OpenStateDB(home)
	if err != nil {
		panic(err)
	}

	_, hash, err := restoreStateFromBlock(blockStoreDB, stateDB, height)
	if err != nil {
		return fmt.Errorf("failed to rollback state: %w", err)
	}
	fmt.Printf("Rolled back state to height %d and hash %X \n", height, hash)
	modifyPrivValidatorsFile(config, height)
	return nil
}

func restoreStateFromBlock(blockStoreDB dbm.DB, stateDB dbm.DB, rollbackHeight int64) (int64, []byte, error) {
	bs := store.NewBlockStore(blockStoreDB)
	ss := state.NewStore(stateDB)
	currentState, err := ss.Load() // current state
	if err != nil {
		return -1, nil, err
	}
	if currentState.IsEmpty() {
		return -1, nil, errors.New("no state found")
	}

	height := bs.Height()

	// NOTE: persistence of state and blocks don't happen atomically. Therefore it is possible that
	// when the user stopped the node the state wasn't updated but the blockstore was. In this situation
	// we don't need to rollback any state and can just return early

	//if height == invalidState.LastBlockHeight+1 {
	//	return invalidState.LastBlockHeight, invalidState.AppHash, nil
	//}

	// If the state store isn't one below nor equal to the blockstore height than this violates the
	// invariant
	//if height != invalidState.LastBlockHeight {
	//	return -1, nil, fmt.Errorf("statestore height (%d) is not one below or equal to blockstore height (%d)",
	//		invalidState.LastBlockHeight, height)
	//}

	// state store height is equal to blockstore height. We're good to proceed with rolling back state
	//rollbackHeight := invalidState.LastBlockHeight - 1

	rollbackBlock := bs.LoadBlockMeta(rollbackHeight)
	if rollbackBlock == nil {
		return -1, nil, fmt.Errorf("block at height %d not found", rollbackHeight)
	}

	currentBlock := bs.LoadBlockMeta(rollbackHeight + 1)
	if currentBlock == nil {
		return -1, nil, fmt.Errorf("block at height %d not found", rollbackHeight+1)
	}

	lastValidatorSet, err := ss.LoadValidators(rollbackHeight)
	if err != nil {
		return -1, nil, err
	}

	validatorSet, err := ss.LoadValidators(rollbackHeight + 1)
	if err != nil {
		return -1, nil, err
	}

	nextValidatorSet, err := ss.LoadValidators(rollbackHeight + 2)

	if err != nil {
		return -1, nil, err
	}

	previousParams, err := ss.LoadConsensusParams(rollbackHeight + 1)
	if err != nil {
		return -1, nil, err
	}

	valChangeHeight := currentState.LastHeightValidatorsChanged
	// this can only happen if the validator set changed since the last block
	if valChangeHeight > rollbackHeight {
		valChangeHeight = rollbackHeight + 1
	}

	paramsChangeHeight := currentState.LastHeightConsensusParamsChanged
	// this can only happen if params changed from the last block
	if paramsChangeHeight > rollbackHeight {
		paramsChangeHeight = rollbackHeight + 1
	}

	// build the new state from the old state and the prior block
	rolledBackState := state.State{
		Version: tmstate.Version{
			Consensus: tmversion.Consensus{
				Block: version.BlockProtocol,
				App:   previousParams.Version.AppVersion,
			},
			Software: version.TMCoreSemVer,
		},
		// immutable fields
		ChainID:       currentState.ChainID,
		InitialHeight: rollbackHeight,

		LastBlockHeight: rollbackBlock.Header.Height,
		LastBlockID:     rollbackBlock.BlockID,
		LastBlockTime:   rollbackBlock.Header.Time,

		NextValidators:              nextValidatorSet,
		Validators:                  validatorSet,
		LastValidators:              lastValidatorSet,
		LastHeightValidatorsChanged: valChangeHeight,

		ConsensusParams:                  previousParams,
		LastHeightConsensusParamsChanged: paramsChangeHeight,

		LastResultsHash: currentBlock.Header.LastResultsHash,
		AppHash:         currentBlock.Header.AppHash,
	}

	// persist the new state. This overrides the invalid one. NOTE: this will also
	// persist the validator set and consensus params over the existing structures,
	// but both should be the same
	if err := ss.Save(rolledBackState); err != nil {
		return -1, nil, fmt.Errorf("failed to save rolled back state: %w", err)
	}
	_, err = PruneRangeBlocks(blockStoreDB, rollbackHeight+1, height)

	//bs.S{Height: int64(config.BaseConfig.RollbackHeight)}.Save(blockStoreDB)

	// remove block meta first as this is used to indicate whether the block exists.
	// For this reason, we also use ony block meta as a measure of the amount of blocks pruned

	if err != nil {
		panic("pruned failed")
	}

	return rolledBackState.LastBlockHeight, rolledBackState.AppHash, nil
}

func modifyPrivValidatorsFile(config *cfg.Config, rollbackHeight int64) {
	pval := privval.LoadOrGenFilePV(config.PrivValidatorKeyFile(), config.PrivValidatorStateFile())
	pval.LastSignState.Signature = nil
	pval.LastSignState.SignBytes = nil
	pval.LastSignState.Step = 0
	pval.LastSignState.Round = 0
	pval.LastSignState.Height = rollbackHeight
	pval.LastSignState.Save()
}
