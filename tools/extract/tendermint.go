package extract

import (
	"github.com/teleport-network/teleport/tools/common"
)

func CopyBlockStore(home string, targetHome string) {
	originBlockStoreDB, err := common.OpenBlockStoreDB(home)
	if err != nil {
		panic(err)
	}
	targetBlockStoreDB, err := common.OpenBlockStoreDB(targetHome)
	if err != nil {
		panic(err)
	}
	migrateBlockStore(originBlockStoreDB, targetBlockStoreDB)
}
