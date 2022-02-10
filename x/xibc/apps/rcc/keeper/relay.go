package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	rcc "github.com/teleport-network/teleport/syscontracts/xibc_rcc"
	"github.com/teleport-network/teleport/x/xibc/apps/rcc/types"
	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
)

func (k Keeper) SendRemoteContractCall(
	ctx sdk.Context,
	destChain string,
	relayChain string,
	sender string,
	contractAddress string,
	data []byte,
) error {
	sourceChain := k.clientKeeper.GetChainName(ctx)
	if sourceChain == destChain {
		return sdkerrors.Wrapf(types.ErrScChainEqualToDestChain, "invalid destChain %s equals to scChain %s", destChain, sourceChain)
	}

	// get the next sequence
	sequence := k.packetKeeper.GetNextSequenceSend(ctx, sourceChain, destChain)

	// TODO: validate packetData
	packetData := types.NewRCCPacketData(
		sourceChain,
		destChain,
		sender,
		contractAddress,
		data,
	)
	packetBz, err := packetData.GetBytes()
	if err != nil {
		return err
	}
	packet := packettypes.NewPacket(sequence, sourceChain, destChain, relayChain, []string{types.PortID}, [][]byte{packetBz})

	return k.packetKeeper.SendPacket(ctx, packet)
}

func (k Keeper) OnRecvPacket(ctx sdk.Context, data types.RCCPacketData) (packettypes.Result, error) {
	// validate packet data upon receiving
	if err := data.ValidateBasic(); err != nil {
		return packettypes.Result{}, err
	}

	res, err := k.RecvPacket(ctx, data)
	if err != nil {
		return packettypes.Result{}, err
	}

	var r struct{ Result packettypes.Result }
	if err := rcc.RCCContract.ABI.UnpackIntoInterface(&r, "onRecvPacket", res.Ret); err != nil {
		return packettypes.Result{}, err
	}

	return r.Result, nil
}

func (k Keeper) RecvPacket(ctx sdk.Context, data types.RCCPacketData) (*evmtypes.MsgEthereumTxResponse, error) {
	return k.CallRCC(ctx, "onRecvPacket", data)
}

func (k Keeper) OnAcknowledgementPacket(ctx sdk.Context, dataHash [32]byte, result []byte) error {
	if _, err := k.AcknowledgementPacket(ctx, dataHash, result); err != nil {
		return err
	}
	return nil
}

func (k Keeper) AcknowledgementPacket(
	ctx sdk.Context,
	dataHash [32]byte,
	result []byte,
) (
	*evmtypes.MsgEthereumTxResponse,
	error,
) {
	return k.CallRCC(ctx, "onAcknowledgementPacket", dataHash, result)
}
