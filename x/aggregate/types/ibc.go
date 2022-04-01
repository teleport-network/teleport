package types

import (
	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
)

// EventTypeIBCAggregate IBC aggregate events
const (
	EventTypeIBCAggregate = "ibc_aggregate"

	AttributeKeyStatus = "status"
	AttributeKeyMsg    = "message"
)

func IBCDenom(port, channel, denom string) (string, error) {
	// since SendPacket did not prefix the denomination, we must prefix denomination here
	sourcePrefix := transfertypes.GetDenomPrefix(port, channel)
	// NOTE: sourcePrefix contains the trailing "/"
	prefixedDenom := sourcePrefix + denom

	// construct the denomination trace from the full raw denomination
	denomTrace := transfertypes.ParseDenomTrace(prefixedDenom)
	voucherDenom := denomTrace.IBCDenom()
	return voucherDenom, nil
}
