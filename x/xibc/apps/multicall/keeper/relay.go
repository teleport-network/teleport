package keeper

import (
	"encoding/json"
	"math/big"
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

	TupleERC20TransferData, err := abi.NewType(
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

	TupleBaseTransferData, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
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
		case types.TransferERC20:
			data, err := abi.Arguments{{Type: TupleERC20TransferData}}.Unpack(calldata.Data[i])
			if err != nil {
				return sdkerrors.Wrapf(types.ErrInvalidMultiCallEvent, "unpack failed, function ID %d", fid)
			}
			var erc20TransferData types.ERC20TransferData
			bz, err := json.Marshal(data[0])
			if err != nil {
				return sdkerrors.Wrapf(types.ErrInvalidMultiCallEvent, "unpack failed, function ID %d", fid)
			}
			if err := json.Unmarshal(bz, &erc20TransferData); err != nil {
				return sdkerrors.Wrapf(types.ErrInvalidMultiCallEvent, "unpack failed, function ID %d", fid)
			}
			ports = append(ports, transfertypes.PortID)
			oriToken, _, _, err := k.aggregateKeeper.QueryERC20Trace(ctx, erc20TransferData.TokenAddress, calldata.DestChain)
			if err != nil {
				return err
			}
			transferBz, err := transfertypes.NewFungibleTokenPacketData(
				sourceChain,
				calldata.DestChain,
				sequence,
				strings.ToLower(sender.String()),
				strings.ToLower(erc20TransferData.Receiver),
				erc20TransferData.Amount.Bytes(),
				strings.ToLower(erc20TransferData.TokenAddress.String()),
				oriToken,
			).GetBytes()
			if err != nil {
				return err
			}
			dataList = append(
				dataList,
				transferBz,
			)
		case types.TransferBase:
			data, err := abi.Arguments{{Type: TupleBaseTransferData}}.Unpack(calldata.Data[i])
			if err != nil {
				return sdkerrors.Wrapf(types.ErrInvalidMultiCallEvent, "unpack failed, function ID %d", fid)
			}
			var baseTransferData types.BaseTransferData
			bz, err := json.Marshal(data[0])
			if err != nil {
				return sdkerrors.Wrapf(types.ErrInvalidMultiCallEvent, "unpack failed, function ID %d", fid)
			}
			if err := json.Unmarshal(bz, &baseTransferData); err != nil {
				return sdkerrors.Wrapf(types.ErrInvalidMultiCallEvent, "unpack failed, function ID %d", fid)
			}
			ports = append(ports, transfertypes.PortID)
			transferBz, err := transfertypes.NewFungibleTokenPacketData(
				sourceChain,
				calldata.DestChain,
				sequence,
				strings.ToLower(sender.String()),
				strings.ToLower(baseTransferData.Receiver),
				baseTransferData.Amount.Bytes(),
				common.BigToAddress(big.NewInt(0)).String(),
				"",
			).GetBytes()
			if err != nil {
				return err
			}
			dataList = append(
				dataList,
				transferBz,
			)
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
			rccPacketDataBz, err := rcctypes.NewRCCPacketData(
				sourceChain,
				calldata.DestChain,
				sequence,
				strings.ToLower(sender.String()),
				strings.ToLower(rccData.ContractAddress),
				rccData.Data,
			).GetBytes()

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
