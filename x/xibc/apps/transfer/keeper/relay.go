package keeper

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	transfer "github.com/teleport-network/teleport/syscontracts/xibc_transfer"
	"github.com/teleport-network/teleport/x/xibc/apps/transfer/types"
	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
)

func (k Keeper) SendTransfer(
	ctx sdk.Context,
	destChain string,
	relayChain string,
	sequence uint64,
	sender string,
	receiver string,
	amount []byte,
	token string,
	oriToken string,
) error {
	sourceChain := k.clientKeeper.GetChainName(ctx)
	if sourceChain == destChain {
		return sdkerrors.Wrapf(types.ErrScChainEqualToDestChain, "invalid destChain %s equals to scChain %s", destChain, sourceChain)
	}

	// get the next sequence
	sequenceTmp := k.packetKeeper.GetNextSequenceSend(ctx, sourceChain, destChain)
	if sequence != sequenceTmp {
		return sdkerrors.Wrapf(types.ErrScChainEqualToDestChain, "invalid sequence %d, %d", sequence, sequenceTmp)
	}
	// TODO: validate packetData
	packetData := types.NewFungibleTokenPacketData(
		sourceChain,
		destChain,
		sequence,
		strings.ToLower(sender),
		strings.ToLower(receiver),
		amount,
		strings.ToLower(token),
		strings.ToLower(oriToken),
	)
	packetDataBz, err := packetData.GetBytes()
	if err != nil {
		return sdkerrors.Wrapf(types.ErrScChainEqualToDestChain, "Get packet bytes err ")
	}
	packet := packettypes.NewPacket(sequence, sourceChain, destChain, relayChain, []string{types.PortID}, [][]byte{packetDataBz})

	return k.packetKeeper.SendPacket(ctx, packet)
}

func (k Keeper) OnRecvPacket(ctx sdk.Context, data types.FungibleTokenPacketData) (packettypes.Result, error) {
	// validate packet data upon receiving
	if err := data.ValidateBasic(); err != nil {
		return packettypes.Result{}, err
	}

	res, err := k.RecvPacket(ctx, data)
	if err != nil {
		return packettypes.Result{}, err
	}

	var r struct{ Result packettypes.Result }
	if err := transfer.TransferContract.ABI.UnpackIntoInterface(&r, "onRecvPacket", res.Ret); err != nil {
		return packettypes.Result{}, err
	}

	return r.Result, nil
}

func (k Keeper) RecvPacket(ctx sdk.Context, data types.FungibleTokenPacketData) (*evmtypes.MsgEthereumTxResponse, error) {
	return k.CallTransfer(ctx, "onRecvPacket", data)
}

func (k Keeper) OnAcknowledgementPacket(ctx sdk.Context, data types.FungibleTokenPacketData, result []byte) error {
	if _, err := k.AcknowledgementPacket(ctx, data, result); err != nil {
		return err
	}
	return nil
}

func (k Keeper) AcknowledgementPacket(
	ctx sdk.Context,
	data types.FungibleTokenPacketData,
	result []byte,
) (
	*evmtypes.MsgEthereumTxResponse,
	error,
) {
	return k.CallTransfer(ctx, "onAcknowledgementPacket", data, result)
}
