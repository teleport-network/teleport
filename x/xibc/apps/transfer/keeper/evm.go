package keeper

import (
	"encoding/json"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/tharsis/ethermint/server/config"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	transfer "github.com/teleport-network/teleport/syscontracts/xibc_transfer"
	"github.com/teleport-network/teleport/x/xibc/apps/transfer/types"
)

// CallTransfer call a method of Transfer contract
func (k Keeper) CallTransfer(
	ctx sdk.Context,
	method string,
	args ...interface{},
) (
	*evmtypes.MsgEthereumTxResponse,
	error,
) {
	payload, err := transfer.TransferContract.ABI.Pack(method, args...)
	if err != nil {
		return nil, sdkerrors.Wrap(
			types.ErrWritingEthTxData,
			sdkerrors.Wrap(err, "failed to create transaction payload").Error(),
		)
	}

	res, err := k.CallEVMWithData(ctx, types.ModuleAddress, &transfer.TransferContractAddress, payload)
	if err != nil {
		return nil, fmt.Errorf("contract call failed: method '%s' %s, %s", method, transfer.TransferContractAddress, err)
	}

	return res, nil
}

// CallEVM performs a smart contract method call using given args
func (k Keeper) CallEVM(
	ctx sdk.Context,
	abi abi.ABI,
	from common.Address,
	contract common.Address,
	method string,
	args ...interface{},
) (
	*evmtypes.MsgEthereumTxResponse, error,
) {
	data, err := abi.Pack(method, args...)
	if err != nil {
		return nil, sdkerrors.Wrap(
			types.ErrWritingEthTxData,
			sdkerrors.Wrap(err, "failed to create transaction data").Error(),
		)
	}

	resp, err := k.CallEVMWithData(ctx, from, &contract, data)
	if err != nil {
		return nil, fmt.Errorf("contract call failed: method '%s' %s, %s", method, contract, err)
	}
	return resp, nil
}

// CallEVMWithData performs a smart contract method call using contract data
func (k Keeper) CallEVMWithData(
	ctx sdk.Context,
	from common.Address,
	contract *common.Address,
	data []byte,
) (
	*evmtypes.MsgEthereumTxResponse,
	error,
) {
	nonce, err := k.accountKeeper.GetSequence(ctx, from.Bytes())
	if err != nil {
		return nil, err
	}

	args, err := json.Marshal(evmtypes.TransactionArgs{
		From: &from,
		To:   contract,
		Data: (*hexutil.Bytes)(&data),
	})
	if err != nil {
		return nil, err
	}

	gasRes, err := k.evmKeeper.EstimateGas(
		sdk.WrapSDKContext(ctx),
		&evmtypes.EthCallRequest{
			Args:   args,
			GasCap: config.DefaultGasCap,
		},
	)
	if err != nil {
		return nil, err
	}

	msg := ethtypes.NewMessage(
		from,
		contract,
		nonce,
		big.NewInt(0),         // amount
		gasRes.Gas,            // gasLimit
		big.NewInt(0),         // gasFeeCap
		big.NewInt(0),         // gasTipCap
		big.NewInt(0),         // gasPrice
		data,                  // tx data
		ethtypes.AccessList{}, // accessList
		true,                  // checkNonce
	)

	res, err := k.evmKeeper.ApplyMessage(ctx, msg, evmtypes.NewNoOpTracer(), true)
	if err != nil {
		return nil, err
	}

	if res.Failed() {
		return nil, sdkerrors.Wrap(evmtypes.ErrVMExecution, res.VmError)
	}

	txLogAttrs := make([]sdk.Attribute, len(res.Logs))
	for i, log := range res.Logs {
		value, err := json.Marshal(log)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "failed to encode log")
		}
		txLogAttrs[i] = sdk.NewAttribute(evmtypes.AttributeKeyTxLog, string(value))
	}

	// emit events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			evmtypes.EventTypeTxLog,
			txLogAttrs...,
		),
	})

	return res, nil
}
