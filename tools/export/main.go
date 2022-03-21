package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	icacontrollertypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/controller/types"
	icahosttypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/host/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibchost "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	aggregatetypes "github.com/teleport-network/teleport/x/aggregate/types"
	xibchost "github.com/teleport-network/teleport/x/xibc/core/host"
	dbm "github.com/tendermint/tm-db"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
	feemarkettypes "github.com/tharsis/ethermint/x/feemarket/types"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

func main() {
	args := os.Args[1:]
	argNum := len(args)
	if argNum != 3 {
		panic(fmt.Sprintf("expected 3 args, found %d", argNum))
	}
	heightStr := args[0]
	h, err := strconv.Atoi(heightStr)
	if err != nil {
		panic(err)
	}
	var height = int64(h)
	homeDir := args[1]
	outputDir := args[2]
	dataDir := filepath.Join(homeDir, "data")
	db, err := sdk.NewLevelDB("application", dataDir)
	if err != nil {
		panic(err)
	}
	outputStore(db, height, outputDir)

}

func outputStore(db dbm.DB, height int64, outputDir string) {
	cms := store.NewCommitMultiStore(db)
	keys := sdk.NewKVStoreKeys(
		// SDK keys
		authtypes.StoreKey,
		banktypes.StoreKey,
		stakingtypes.StoreKey,
		distrtypes.StoreKey,
		slashingtypes.StoreKey,
		govtypes.StoreKey,
		paramstypes.StoreKey,
		upgradetypes.StoreKey,
		evidencetypes.StoreKey,
		capabilitytypes.StoreKey,
		feegrant.StoreKey,
		authzkeeper.StoreKey,
		// ibc keys
		ibchost.StoreKey,
		ibctransfertypes.StoreKey,
		icacontrollertypes.StoreKey,
		icahosttypes.StoreKey,
		// xibc keys
		xibchost.StoreKey,
		// ethermint keys
		evmtypes.StoreKey,
		feemarkettypes.StoreKey,
		// teleport keys
		aggregatetypes.StoreKey,
	)

	for _, v := range keys {
		cms.MountStoreWithDB(v, sdk.StoreTypeIAVL, nil)
	}

	err := cms.LoadVersion(height)
	if err != nil {
		panic(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(len(keys))
	for _, storeKey := range keys {
		ckvstore := cms.GetCommitKVStore(storeKey)
		commitID := ckvstore.LastCommitID()
		fmt.Printf("store key name: %s, commitId: %s \n", storeKey.Name(), hex.EncodeToString(commitID.Hash))

		go func(sk types.StoreKey) {
			defer wg.Done()
			filePath := filepath.Join(outputDir, sk.Name()+"_"+hex.EncodeToString(commitID.Hash))
			file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				panic(err)
			}
			defer file.Close()
			write := bufio.NewWriter(file)

			itr := ckvstore.Iterator(nil, nil)
			defer itr.Close()

			for ; itr.Valid(); itr.Next() {
				k, v := itr.Key(), itr.Value()
				key := hex.EncodeToString(k)
				value := hex.EncodeToString(v)
				write.WriteString(fmt.Sprintf("%s,%s\n", key, value))
			}
			if err = write.Flush(); err != nil {
				panic(err)
			}
			if err = itr.Error(); err != nil {
				panic(err)
			}
		}(storeKey)
	}
	wg.Wait()
	fmt.Println("job finished")
}
