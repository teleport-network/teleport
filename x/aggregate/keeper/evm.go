package keeper

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/tharsis/ethermint/server/config"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	transfercontract "github.com/teleport-network/teleport/syscontracts/xibc_transfer"
	"github.com/teleport-network/teleport/x/aggregate/types"
	"github.com/teleport-network/teleport/x/aggregate/types/contracts"
)

// QueryERC20 returns the data of a deployed ERC20 contract
func (k Keeper) QueryERC20(ctx sdk.Context, contract common.Address) (types.ERC20Data, error) {
	var (
		nameRes    types.ERC20StringResponse
		symbolRes  types.ERC20StringResponse
		decimalRes types.ERC20Uint8Response
	)

	erc20 := contracts.ERC20BurnableContract.ABI

	// Name
	res, err := k.CallEVM(ctx, erc20, types.ModuleAddress, contract, "name")
	if err != nil {
		return types.ERC20Data{}, err
	}

	if err := erc20.UnpackIntoInterface(&nameRes, "name", res.Ret); err != nil {
		return types.ERC20Data{}, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "failed to unpack name: %s", err.Error())
	}

	// Symbol
	res, err = k.CallEVM(ctx, erc20, types.ModuleAddress, contract, "symbol")
	if err != nil {
		return types.ERC20Data{}, err
	}

	if err := erc20.UnpackIntoInterface(&symbolRes, "symbol", res.Ret); err != nil {
		return types.ERC20Data{}, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "failed to unpack symbol: %s", err.Error())
	}

	// Decimals
	res, err = k.CallEVM(ctx, erc20, types.ModuleAddress, contract, "decimals")
	if err != nil {
		return types.ERC20Data{}, err
	}

	if err := erc20.UnpackIntoInterface(&decimalRes, "decimals", res.Ret); err != nil {
		return types.ERC20Data{}, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "failed to unpack decimals: %s", err.Error())
	}

	return types.NewERC20Data(nameRes.Value, symbolRes.Value, decimalRes.Value), nil
}

// QueryERC20Trace returns if the ERC20 trace exist
func (k Keeper) QueryERC20Trace(
	ctx sdk.Context,
	erc20Address common.Address,
	originChain string,
) (
	string, *big.Int, bool, error,
) {
	res, err := k.CallEVM(
		ctx,
		transfercontract.TransferContract.ABI,
		types.ModuleAddress,
		transfercontract.TransferContractAddress,
		"bindings",
		fmt.Sprintf("%s/%s", strings.ToLower(erc20Address.String()), originChain),
	)
	if err != nil {
		return "", nil, false, fmt.Errorf("call findbinding failed: %s", err)
	}

	var binding types.BindingsResponse
	if err := transfercontract.TransferContract.ABI.UnpackIntoInterface(&binding, "bindings", res.Ret); err != nil {
		return "", nil, false, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "failed to unpack findbinding: %s", err.Error())
	}

	return binding.OriToken, binding.Amount, binding.Bound, nil
}

// CallEVM performs a smart contract method call using  given args
func (k Keeper) CallEVM(
	ctx sdk.Context,
	abi abi.ABI,
	from common.Address,
	contract common.Address,
	method string,
	args ...interface{},
) (
	*evmtypes.MsgEthereumTxResponse,
	error,
) {
	payload, err := abi.Pack(method, args...)
	if err != nil {
		return nil, sdkerrors.Wrap(
			types.ErrWritingEthTxPayload,
			sdkerrors.Wrap(err, "failed to create transaction payload").Error(),
		)
	}

	resp, err := k.CallEVMWithPayload(ctx, from, &contract, payload)
	if err != nil {
		return nil, fmt.Errorf("contract call failed: method '%s' %s, %s", method, contract, err)
	}
	return resp, nil
}

// CallEVMWithPayload performs a smart contract method call using contract data
func (k Keeper) CallEVMWithPayload(
	ctx sdk.Context,
	from common.Address,
	contract *common.Address,
	transferData []byte,
) (
	*evmtypes.MsgEthereumTxResponse,
	error,
) {
	nonce, err := k.accountKeeper.GetSequence(ctx, from.Bytes())
	if err != nil {
		return nil, err
	}

	msg := ethtypes.NewMessage(
		from,
		contract,
		nonce,
		big.NewInt(0),        // amount
		config.DefaultGasCap, // gasLimit
		big.NewInt(0),        // gasFeeCap
		big.NewInt(0),        // gasTipCap
		big.NewInt(0),        // gasPrice
		transferData,
		ethtypes.AccessList{}, // AccessList
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
