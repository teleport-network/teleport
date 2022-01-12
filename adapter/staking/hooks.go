package staking

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

func (h *HookAdapter) PostTxProcessing(
	ctx sdk.Context,
	from common.Address,
	to *common.Address,
	receipt *ethtypes.Receipt,
) error {
	if len(receipt.Logs) != 1 { // staking hook only handles one event to prevent invocation abuse
		return nil
	}

	if len(receipt.Logs[0].Topics) == 0 {
		return nil
	}

	if !bytes.Equal(receipt.Logs[0].Address.Bytes(), h.stakingContract.Bytes()) {
		return nil
	}

	handler, ok := h.handlers[receipt.Logs[0].Topics[0]]
	if !ok {
		return nil // return if no related handler found
	}

	return handler(ctx, receipt.Logs[0])
}
