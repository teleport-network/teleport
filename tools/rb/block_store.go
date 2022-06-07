package rb

import (
	"fmt"

	"github.com/gogo/protobuf/proto"

	tmstore "github.com/tendermint/tendermint/proto/tendermint/store"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"github.com/tendermint/tendermint/store"
	"github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

// PruneRangeBlocks removes block up to (but not including) a height. It returns the number of blocks pruned.
func PruneRangeBlocks(db dbm.DB, start int64, height int64) (uint64, error) {
	if start > height {
		return 0, fmt.Errorf("height must be greater than start")
	}

	base := start
	pruned := uint64(0)
	batch := db.NewBatch()
	defer batch.Close()
	flush := func(batch dbm.Batch, base int64, height int64) error {
		// We can't trust batches to be atomic, so update base first to make sure noone
		// tries to access missing blocks.
		bss := tmstore.BlockStoreState{
			Base:   base,
			Height: height,
		}
		store.SaveBlockStoreState(&bss, db)

		err := batch.WriteSync()
		if err != nil {
			return fmt.Errorf("failed to prune up to height %v: %w", base, err)
		}
		batch.Close()
		return nil
	}

	for h := base; h < height; h++ {
		meta := LoadBlockMeta(db, h)
		if meta == nil { // assume already deleted
			continue
		}
		if err := batch.Delete(calcBlockMetaKey(h)); err != nil {
			return 0, err
		}
		if err := batch.Delete(calcBlockHashKey(meta.BlockID.Hash)); err != nil {
			return 0, err
		}
		if err := batch.Delete(calcBlockCommitKey(h)); err != nil {
			return 0, err
		}
		if err := batch.Delete(calcSeenCommitKey(h)); err != nil {
			return 0, err
		}
		for p := 0; p < int(meta.BlockID.PartSetHeader.Total); p++ {
			if err := batch.Delete(calcBlockPartKey(h, p)); err != nil {
				return 0, err
			}
		}
		pruned++

		// flush every 1000 blocks to avoid batches becoming too large
		if pruned%1000 == 0 && pruned > 0 {
			err := flush(batch, h, h)
			if err != nil {
				return 0, err
			}
			batch = db.NewBatch()
			defer batch.Close()
		}
	}

	err := flush(batch, height, height)
	if err != nil {
		return 0, err
	}

	bss := tmstore.BlockStoreState{
		Base:   start - 1,
		Height: start - 1,
	}
	store.SaveBlockStoreState(&bss, db)

	return pruned, nil
}

func LoadBlockMeta(db dbm.DB, height int64) *types.BlockMeta {
	var pbbm = new(tmproto.BlockMeta)
	bz, err := db.Get(calcBlockMetaKey(height))

	if err != nil {
		panic(err)
	}

	if len(bz) == 0 {
		return nil
	}

	err = proto.Unmarshal(bz, pbbm)
	if err != nil {
		panic(fmt.Errorf("unmarshal to tmproto.BlockMeta: %w", err))
	}

	blockMeta, err := types.BlockMetaFromProto(pbbm)
	if err != nil {
		panic(fmt.Errorf("error from proto blockMeta: %w", err))
	}

	return blockMeta
}

func calcBlockMetaKey(height int64) []byte {
	return []byte(fmt.Sprintf("H:%v", height))
}

func calcBlockPartKey(height int64, partIndex int) []byte {
	return []byte(fmt.Sprintf("P:%v:%v", height, partIndex))
}

func calcBlockCommitKey(height int64) []byte {
	return []byte(fmt.Sprintf("C:%v", height))
}

func calcSeenCommitKey(height int64) []byte {
	return []byte(fmt.Sprintf("SC:%v", height))
}

func calcBlockHashKey(hash []byte) []byte {
	return []byte(fmt.Sprintf("BH:%x", hash))
}
