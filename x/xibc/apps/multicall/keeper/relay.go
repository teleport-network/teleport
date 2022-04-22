package keeper

import (
	"encoding/json"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/teleport-network/teleport/x/xibc/apps/multicall/types"
	rcctypes "github.com/teleport-network/teleport/x/xibc/apps/rcc/types"
	transfertypes "github.com/teleport-network/teleport/x/xibc/apps/transfer/types"
	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
)

func (k Keeper) SendMultiCall(ctx sdk.Context, sender common.Address, calldata types.MultiCallData) error {
	sourceChain := k.clientKeeper.GetChainName(ctx)
	if sourceChain == calldata.DestChain {
		return sdkerrors.Wrapf(types.ErrScChainEqualToDestChain, "invalid destChain %s equals to scChain %s", calldata.DestChain, sourceChain)
	}

	// get the next sequence
	sequence := k.packetKeeper.GetNextSequenceSend(ctx, sourceChain, calldata.DestChain)

	TupleTransferData, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "token_address", Type: "address"},
			{Name: "receiver", Type: "string"},
			{Name: "amount", Type: "uint256"},
		},
	)
	if err != nil {
		return err
	}

	TupleRCCData, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "contract_address", Type: "string"},
			{Name: "data", Type: "bytes"},
		},
	)
	if err != nil {
		return err
	}

	var ports []string
	var dataList [][]byte

	for i, fid := range calldata.Functions {
		switch fid {
		case types.Transfer:
			data, err := abi.Arguments{{Type: TupleTransferData}}.Unpack(calldata.Data[i])
			if err != nil {
				return sdkerrors.Wrapf(types.ErrInvalidMultiCallEvent, "unpack failed, function ID %d", fid)
			}
			var transferData types.TransferData
			bz, err := json.Marshal(data[0])
			if err != nil {
				return sdkerrors.Wrapf(types.ErrInvalidMultiCallEvent, "unpack failed, function ID %d", fid)
			}
			if err := json.Unmarshal(bz, &transferData); err != nil {
				return sdkerrors.Wrapf(types.ErrInvalidMultiCallEvent, "unpack failed, function ID %d", fid)
			}
			ports = append(ports, transfertypes.PortID)
			oriToken, _, _, err := k.aggregateKeeper.QueryERC20Trace(ctx, transferData.TokenAddress, calldata.DestChain)
			if err != nil {
				return err
			}
			transferPacketData := transfertypes.NewFungibleTokenPacketData(
				sourceChain,
				calldata.DestChain,
				sequence,
				strings.ToLower(sender.String()),
				strings.ToLower(transferData.Receiver),
				transferData.Amount.Bytes(),
				strings.ToLower(transferData.TokenAddress.String()),
				oriToken,
			)
			if err := transferPacketData.ValidateBasic(); err != nil {
				return sdkerrors.Wrapf(transfertypes.ErrInvalidPacket, "invalid packet data")
			}
			transferPacketDataBz, err := transferPacketData.GetBytes()
			if err != nil {
				return sdkerrors.Wrapf(transfertypes.ErrABIPack, "get packet bytes error")
			}
			dataList = append(dataList, transferPacketDataBz)
		case types.RemoteCall:
			data, err := abi.Arguments{{Type: TupleRCCData}}.Unpack(calldata.Data[i])
			if err != nil {
				return sdkerrors.Wrapf(types.ErrInvalidMultiCallEvent, "unpack failed, function ID %d", fid)
			}
			var rccData types.RCCData
			bz, err := json.Marshal(data[0])
			if err != nil {
				return sdkerrors.Wrapf(types.ErrInvalidMultiCallEvent, "unpack failed, function ID %d", fid)
			}
			if err := json.Unmarshal(bz, &rccData); err != nil {
				return sdkerrors.Wrapf(types.ErrInvalidMultiCallEvent, "unpack failed, function ID %d", fid)
			}
			ports = append(ports, rcctypes.PortID)
			rccPacketData := rcctypes.NewRCCPacketData(
				sourceChain,
				calldata.DestChain,
				sequence,
				strings.ToLower(sender.String()),
				strings.ToLower(rccData.ContractAddress),
				rccData.Data,
			)
			if err := rccPacketData.ValidateBasic(); err != nil {
				return sdkerrors.Wrapf(rcctypes.ErrInvalidPacket, "invalid packet data")
			}
			rccPacketDataBz, err := rccPacketData.GetBytes()
			if err != nil {
				return sdkerrors.Wrapf(types.ErrInvalidMultiCallEvent, "unpack failed, function ID %d", fid)
			}
			dataList = append(
				dataList,
				rccPacketDataBz,
			)
		default:
			return sdkerrors.Wrapf(types.ErrInvalidMultiCallEvent, "invalid function ID %d", fid)
		}
	}

	packet := packettypes.NewPacket(
		sequence,
		sourceChain,
		calldata.DestChain,
		calldata.RelayChain,
		ports,
		dataList,
	)

	return k.packetKeeper.SendPacket(ctx, packet)
}
