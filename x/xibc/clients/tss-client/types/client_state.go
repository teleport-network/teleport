package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/teleport-network/teleport/x/xibc/core/client/types"
	"github.com/teleport-network/teleport/x/xibc/exported"
)

var _ exported.ClientState = (*ClientState)(nil)

func (cs ClientState) ClientType() string {
	return exported.TSS
}

func (cs ClientState) GetLatestHeight() exported.Height {
	return types.Height{}
}

func (cs ClientState) CheckMsg(msg sdk.Msg) error {
	signer := msg.GetSigners()[0].String()
	if cs.TssAddress != signer {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidAddress,
			"invalid TSS address, expected %s, got %s",
			cs.TssAddress, signer)
	}
	return nil
}

func (cs ClientState) Validate() error {
	if _, err := sdk.AccAddressFromBech32(cs.TssAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "string could not be parsed as address: %v", err)
	}
	return nil
}

func (cs ClientState) GetDelayTime() uint64 {
	return 0
}

func (cs ClientState) GetDelayBlock() uint64 {
	return 0
}

func (cs ClientState) GetPrefix() exported.Prefix {
	return nil
}

func (cs ClientState) Initialize(
	ctx sdk.Context,
	cdc codec.BinaryCodec,
	store sdk.KVStore,
	state exported.ConsensusState,
) error {
	return nil
}

func (cs ClientState) Status(
	ctx sdk.Context,
	store sdk.KVStore,
	cdc codec.BinaryCodec,
) exported.Status {
	// always return active
	return exported.Active
}

func (cs ClientState) UpgradeState(
	ctx sdk.Context,
	cdc codec.BinaryCodec,
	store sdk.KVStore,
	state exported.ConsensusState,
) error {
	return nil
}
func (cs ClientState) ExportMetadata(store sdk.KVStore) []exported.GenesisMetadata {
	return nil
}

func (cs ClientState) VerifyPacketCommitment(
	ctx sdk.Context,
	store sdk.KVStore,
	cdc codec.BinaryCodec,
	height exported.Height,
	proof []byte,
	sourceChain, destChain string,
	sequence uint64,
	commitment []byte,
) error {
	if string(proof) != cs.TssAddress {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidAddress,
			"invalid TSS address, expected %s, got %s",
			cs.TssAddress, string(proof),
		)
	}
	return nil
}

func (m ClientState) VerifyPacketAcknowledgement(
	ctx sdk.Context,
	store sdk.KVStore,
	cdc codec.BinaryCodec,
	height exported.Height,
	proof []byte,
	sourceChain, destChain string,
	sequence uint64,
	ackBytes []byte,
) error {
	if string(proof) != m.TssAddress {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidAddress,
			"invalid TSS address, expected %s, got %s",
			m.TssAddress, string(proof),
		)
	}
	return nil
}

// ClientType returns tss
func (ConsensusState) ClientType() string {
	return exported.TSS
}

// GetRoot returns the commitment Root for the specific
func (cs ConsensusState) GetRoot() []byte {
	return nil
}

func (cs ConsensusState) GetTimestamp() uint64 {
	return 0
}

func (cs ConsensusState) ValidateBasic() error {
	return nil
}
