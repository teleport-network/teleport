package rb

import (
	"fmt"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/gogo/protobuf/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/bitdao-io/bitnetwork/tools/common"
)

func RollbackApp(home string, storeKeys map[string]*sdk.KVStoreKey, height int64) error {
	db, err := common.OpenApplicationDB(home)
	if err != nil {
		return err
	}
	latestVersion := common.GetLatestVersion(db)
	fmt.Println(latestVersion)
	for _, key := range storeKeys {
		prefixDB := common.PrefixDB(key, db)
		fmt.Println("keyName: ", key.Name())
		if err := deleteVersionsFrom(height+1, latestVersion, prefixDB); err != nil {
			return err
		}
	}

	// update latest height
	bz, err := gogotypes.StdInt64Marshal(height)
	if err != nil {
		return err
	}
	if err := db.Set([]byte(common.LatestVersionKey), bz); err != nil {
		return err
	}

	return nil
}

func deleteVersionsFrom(version int64, latestVersion int64, db *dbm.PrefixDB) error {
	latestRoot, err := db.Get(common.RootKeyFormat.Key(latestVersion))
	if err != nil {
		return err
	}

	batch := db.NewBatch()
	if len(latestRoot) > 0 {
		// First, delete all active nodes in the current (latest) version whose node version is after
		// the given version.
		if err := deleteNodesFrom(db, batch, version, common.NodeKeyFormat.Key(latestRoot)); err != nil {
			return err
		}
	}

	// Next, delete orphans:
	// - Delete orphan entries *and referred nodes* with fromVersion >= version
	// - Delete orphan entries with toVersion >= version-1 (since orphans at latest are not orphans)
	traversePrefix(db, common.OrphanKeyFormat.Key(), func(key, hash []byte) {
		var fromVersion, toVersion int64
		common.OrphanKeyFormat.Scan(key, &toVersion, &fromVersion)

		if fromVersion >= version {
			if err = batch.Delete(key); err != nil {
				panic(err)
			}
			if err = batch.Delete(common.NodeKeyFormat.Key(hash)); err != nil {
				panic(err)
			}
		} else if toVersion >= version-1 {
			if err := batch.Delete(key); err != nil {
				panic(err)
			}
		}
	})

	// Finally, delete the version root entries
	traverseRange(db, common.RootKeyFormat.Key(version), common.RootKeyFormat.Key(int64(math.MaxInt64)), func(k, v []byte) {
		if err := batch.Delete(k); err != nil {
			panic(err)
		}
	})

	return batch.WriteSync()
}

// Traverse all keys with a certain prefix.
func traversePrefix(db dbm.DB, prefix []byte, fn func(k, v []byte)) {
	itr, err := dbm.IteratePrefix(db, prefix)
	if err != nil {
		panic(err)
	}
	defer itr.Close()

	for ; itr.Valid(); itr.Next() {
		fn(itr.Key(), itr.Value())
	}
}

// Traverse all keys between a given range (excluding end).
func traverseRange(db dbm.DB, start []byte, end []byte, fn func(k, v []byte)) {
	itr, err := db.Iterator(start, end)
	if err != nil {
		panic(err)
	}
	defer itr.Close()

	for ; itr.Valid(); itr.Next() {
		fn(itr.Key(), itr.Value())
	}

	if err := itr.Error(); err != nil {
		panic(err)
	}
}

func deleteNodesFrom(db *dbm.PrefixDB, batch dbm.Batch, version int64, hash []byte) error {
	buf, err := db.Get(hash)
	if err != nil {
		return err
	}
	node, err := common.MakeNode(buf)
	if err != nil {
		return err
	}
	node.Hash = hash

	if node.LeftHash != nil {
		if err := deleteNodesFrom(db, batch, version, common.NodeKeyFormat.Key(node.LeftHash)); err != nil {
			return err
		}
	}
	if node.RightHash != nil {
		if err := deleteNodesFrom(db, batch, version, common.NodeKeyFormat.Key(node.RightHash)); err != nil {
			return err
		}
	}

	if node.Version >= version {
		if err := batch.Delete(hash); err != nil {
			return err
		}
	}

	return nil
}
