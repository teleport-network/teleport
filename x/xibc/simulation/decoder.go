package simulation

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/kv"

	clientsim "github.com/teleport-network/teleport/x/xibc/core/client/simulation"
	host "github.com/teleport-network/teleport/x/xibc/core/host"
	packetsim "github.com/teleport-network/teleport/x/xibc/core/packet/simulation"
	"github.com/teleport-network/teleport/x/xibc/keeper"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding xibc type.
func NewDecodeStore(k keeper.Keeper) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		if res, found := clientsim.NewDecodeStore(k.ClientKeeper, kvA, kvB); found {
			return res
		}

		if res, found := packetsim.NewDecodeStore(k.Codec(), kvA, kvB); found {
			return res
		}

		panic(fmt.Sprintf("invalid %s key prefix: %s", host.ModuleName, string(kvA.Key)))
	}
}
