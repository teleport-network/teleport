package common

import (
	"crypto/sha256"
	"path/filepath"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/iavl"

	dbm "github.com/tendermint/tm-db"

	gogotypes "github.com/gogo/protobuf/types"
)

const (
	Int64Size = 8
	HashSize  = sha256.Size

	LatestVersionKey = "s/latest"
	CommitInfoKeyFmt = "s/%d" // s/<version>
)

// Root nodes are indexed separately by their version
var (
	RootKeyFormat   = iavl.NewKeyFormat('r', Int64Size)                      // r<version>
	NodeKeyFormat   = iavl.NewKeyFormat('n', HashSize)                       // n<hash>
	OrphanKeyFormat = iavl.NewKeyFormat('o', Int64Size, Int64Size, HashSize) // o<last-version><first-version><hash>

)

func OpenApplicationDB(rootDir string) (dbm.DB, error) {
	dataDir := filepath.Join(rootDir, "data")
	return sdk.NewLevelDB("application", dataDir)
}

func PrefixDB(key *sdk.KVStoreKey, originDB dbm.DB) *dbm.PrefixDB {
	prefix := "s/k:" + key.Name() + "/"
	return dbm.NewPrefixDB(originDB, []byte(prefix))
}

func GetLatestVersion(db dbm.DB) int64 {
	bz, err := db.Get([]byte(LatestVersionKey))
	if err != nil {
		panic(err)
	} else if bz == nil {
		return 0
	}

	var latestVersion int64

	if err := gogotypes.StdInt64Unmarshal(&latestVersion, bz); err != nil {
		panic(err)
	}

	return latestVersion
}

func Commit(batches []dbm.Batch) error {
	for _, batch := range batches {
		if err := batch.WriteSync(); err != nil {
			return err
		}
		if err := batch.Close(); err != nil {
			return err
		}
	}
	return nil
}
