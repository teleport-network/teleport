package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	"github.com/teleport-network/teleport/x/xibc/exported"
)

func (cs ClientState) CheckHeaderAndUpdateState(
	ctx sdk.Context,
	cdc codec.BinaryCodec,
	store sdk.KVStore,
	header exported.Header,
) (
	exported.ClientState,
	exported.ConsensusState,
	error,
) {
	tssHeader, ok := header.(*Header)
	if !ok {
		return nil, nil, sdkerrors.Wrapf(clienttypes.ErrInvalidHeader, "expected type %T, got %T", &Header{}, header)
	}
	cs.TssAddress = tssHeader.TssAddress
	cs.Pubkey = tssHeader.Pubkey
	cs.PartPubkeys = tssHeader.PartPubkeys
	cs.Threshold = tssHeader.Threshold
	return &cs, nil, nil
}
