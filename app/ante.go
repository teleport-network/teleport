package app

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethAnte "github.com/tharsis/ethermint/app/ante"
)

func validate(options ethAnte.HandlerOptions) error {
	if options.AccountKeeper == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "account keeper is required for AnteHandler")
	}
	if options.BankKeeper == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "bank keeper is required for AnteHandler")
	}
	if options.SignModeHandler == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for ante builder")
	}
	return nil
}
