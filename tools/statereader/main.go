package main

import (
	"encoding/json"
	"fmt"
	sm "github.com/tendermint/tendermint/state"
	dbm "github.com/tendermint/tm-db"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	args := os.Args[1:]
	argNum := len(args)
	if argNum != 2 {
		panic(fmt.Sprintf("expected 2 args, found %d", argNum))
	}
	heightStr := args[0]
	h, err := strconv.Atoi(heightStr)
	if err != nil {
		panic(err)
	}
	var height = int64(h)
	homeDir := args[1]

	dbDir := filepath.Join(homeDir, "data")
	stateDB, err := dbm.NewDB("state", "goleveldb", dbDir)
	if err != nil {
		panic(err)
	}
	stateStore := sm.NewStore(stateDB)
	res, err := stateStore.LoadABCIResponses(height)
	if err != nil {
		panic(err)
	}

	jsonStr, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonStr))
	//for _, e := range res.BeginBlock.GetEvents() {
	//	fmt.Println("type: ", e.GetType())
	//	for _, ea := range e.Attributes {
	//		fmt.Println(fmt.Sprintf("key: %s, value: %s", string(ea.Key), string(ea.GetValue())))
	//	}
	//	fmt.Println()
	//}

}
