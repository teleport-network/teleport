package extract

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitdao-io/bitnetwork/tools/common"

	dbm "github.com/tendermint/tm-db"

	"github.com/pkg/errors"
)

// CopyApplication extract the latest version data of application db, and store it into target db
func CopyApplication(home string, targetHome string, storeKeys map[string]*sdk.KVStoreKey) error {
	originDB, err := common.OpenApplicationDB(home)
	if err != nil {
		return err
	}
	targetDB, err := common.OpenApplicationDB(targetHome)
	if err != nil {
		return err
	}
	targetBatch := targetDB.NewBatch()
	latestVersion := common.GetLatestVersion(originDB)

	// migrate latestVersion to target db
	if err = migrateLatestVersion(targetBatch, latestVersion); err != nil {
		return err
	}

	// migrate commit info to target db
	if err = migrateCommitInfo(originDB, targetBatch, latestVersion); err != nil {
		return err
	}

	targetBatchSlice := []dbm.Batch{targetBatch}
	for _, key := range storeKeys {
		prefixOriginDB := common.PrefixDB(key, originDB)
		rootHash, err := prefixOriginDB.Get(common.RootKeyFormat.Key(latestVersion))
		if err != nil {
			return errors.New("fail to get root from original DB")
		}

		prefixTargetBatch := common.PrefixDB(key, targetDB).NewBatch()

		// migrate each root of prefixed stores(module store) to target db
		if err = migrateStoreRoot(prefixTargetBatch, latestVersion, rootHash); err != nil {
			return err
		}

		// migrate data of prefixed stores(module store) to target db
		if err = migrateStoreData(prefixOriginDB, prefixTargetBatch, rootHash); err != nil {
			return err
		}
		targetBatchSlice = append(targetBatchSlice, prefixTargetBatch)
	}
	return common.Commit(targetBatchSlice)
}
