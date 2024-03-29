package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/teleport-network/teleport/x/xibc/core/host"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's Value
func NewDecodeStore(cdc codec.BinaryCodec, kvA, kvB kv.Pair) (string, bool) {
	switch {
	case bytes.HasPrefix(kvA.Key, []byte(host.KeyNextSeqSendPrefix)):
		seqA := sdk.BigEndianToUint64(kvA.Value)
		seqB := sdk.BigEndianToUint64(kvB.Value)
		return fmt.Sprintf("NextSeqSend A: %d\nNextSeqSend B: %d", seqA, seqB), true

	case bytes.HasPrefix(kvA.Key, []byte(host.KeyPacketCommitmentPrefix)):
		return fmt.Sprintf("CommitmentHash A: %X\nCommitmentHash B: %X", kvA.Value, kvB.Value), true

	case bytes.HasPrefix(kvA.Key, []byte(host.KeyPacketAckPrefix)):
		return fmt.Sprintf("AckHash A: %X\nAckHash B: %X", kvA.Value, kvB.Value), true

	default:
		return "", false
	}
}
