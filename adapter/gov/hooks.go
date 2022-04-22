package gov

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

func (h *HookAdapter) PostTxProcessing(
	ctx sdk.Context,
	_ core.Message,
	receipt *ethtypes.Receipt,
) error {
	for _, log := range receipt.Logs {
		if bytes.Equal(log.Address.Bytes(), h.govContract.Bytes()) { // only care the logs from gov contract
			handler, ok := h.handlers[log.Topics[0]]
			if !ok {
				continue
			}
			if err := handler(ctx, log); err != nil {
				return err
			}
		}
	}
	return nil
}
