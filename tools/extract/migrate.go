package extract

import (
	"fmt"

	"github.com/bitdao-io/bitnetwork/tools/common"

	"github.com/tendermint/tendermint/store"
	"github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	gogotypes "github.com/gogo/protobuf/types"
)

// migrate application db
func migrateStoreRoot(prefixDstBatch dbm.Batch, version int64, rootHash []byte) error {
	return prefixDstBatch.Set(common.RootKeyFormat.Key(version), rootHash)
}

func migrateCommitInfo(srcDB dbm.DB, dstDB dbm.Batch, version int64) error {
	cInfoKey := fmt.Sprintf(common.CommitInfoKeyFmt, version)

	bz, err := srcDB.Get([]byte(cInfoKey))
	if err != nil {
		return err
	}
	return dstDB.Set([]byte(cInfoKey), bz)
}

func migrateStoreData(srcDB dbm.DB, dstDB dbm.Batch, hash []byte) error {
	if len(hash) == 0 {
		return nil
	}
	buf, err := srcDB.Get(common.NodeKeyFormat.Key(hash))
	if err != nil {
		return err
	}
	node, err := common.MakeNode(buf)
	if err != nil {
		return err
	}
	if node.LeftHash != nil {
		if err := migrateStoreData(srcDB, dstDB, node.LeftHash); err != nil {
			return err
		}
	}
	if node.RightHash != nil {
		if err := migrateStoreData(srcDB, dstDB, node.RightHash); err != nil {
			return err
		}
	}
	if err = dstDB.Set(common.NodeKeyFormat.Key(hash), buf); err != nil {
		return err
	}
	return nil
}

func migrateLatestVersion(batch dbm.Batch, version int64) error {
	bz, err := gogotypes.StdInt64Marshal(version)
	if err != nil {
		panic(err)
	}
	return batch.Set([]byte(common.LatestVersionKey), bz)
}

// -------------------------------------------------------------------------------------------------------

// migrate tendermint block store
func migrateBlockStore(originBlockStoreDB dbm.DB, targetBlockStoreDB dbm.DB) {
	originBlockStore := store.NewBlockStore(originBlockStoreDB)
	targetBlockStore := store.NewBlockStore(targetBlockStoreDB)

	latestBlockHeight := originBlockStore.Height()
	block := originBlockStore.LoadBlock(latestBlockHeight)
	targetBlockStore.SaveBlock(block, block.MakePartSet(types.BlockPartSizeBytes), originBlockStore.LoadSeenCommit(latestBlockHeight))
}
