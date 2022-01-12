package common

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

type EvmEventHandler func(sdk.Context, *ethtypes.Log) error

func ExecuteMsg(ctx sdk.Context, router *baseapp.MsgServiceRouter, msg sdk.Msg) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	handler := router.Handler(msg)
	if handler == nil {
		return errors.New("no handler found to handle Msg")
	}
	_, err := handler(ctx, msg)
	return err
}
