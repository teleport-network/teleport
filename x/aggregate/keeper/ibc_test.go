package keeper_test

import (
	"fmt"

	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
)

func (suite *KeeperTestSuite) TestDenom() {
	sourcePrefix := transfertypes.GetDenomPrefix("transfer", "channel-0")
	// NOTE: sourcePrefix contains the trailing "/"
	prefixedDenom := sourcePrefix + "atele"

	// construct the denomination trace from the full raw denomination
	denomTrace := transfertypes.ParseDenomTrace(prefixedDenom)

	traceHash := denomTrace.Hash()
	voucherDenom := denomTrace.IBCDenom()
	fmt.Println(traceHash)
	fmt.Println(voucherDenom)
}
