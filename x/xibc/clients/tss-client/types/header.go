package types

import (
	"github.com/teleport-network/teleport/x/xibc/exported"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ exported.Header = (*Header)(nil)

func (h Header) ClientType() string {
	return exported.TSS
}

func (h Header) GetHeight() exported.Height {
	return nil
}

func (h Header) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(h.TssAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "string could not be parsed as address: %v", err)
	}
	return nil
}
