package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	transfer "github.com/teleport-network/teleport/syscontracts/xibc_transfer"
	"github.com/teleport-network/teleport/x/aggregate/types"
)

func (k Keeper) AddERC20TraceToTransferContract(
	ctx sdk.Context,
	contract common.Address,
	originToken string,
	originChain string,
	scale uint8,
) (
	*evmtypes.MsgEthereumTxResponse,
	error,
) {
	payload, err := transfer.TransferContract.ABI.Pack("bindToken", contract, originToken, originChain, scale)
	if err != nil {
		return nil, err
	}

	return k.CallEVMWithData(ctx, types.ModuleAddress, &transfer.TransferContractAddress, payload)
}

func (k Keeper) EnableTimeBasedSupplyLimitInTransferContract(
	ctx sdk.Context,
	erc20Address common.Address,
	timePeriod *big.Int,
	timeBasedLimit *big.Int,
	maxAmount *big.Int,
	minAmount *big.Int,
) (
	*evmtypes.MsgEthereumTxResponse,
	error,
) {
	payload, err := transfer.TransferContract.ABI.Pack(
		"enableTimeBasedSupplyLimit",
		erc20Address,
		timePeriod,
		timeBasedLimit,
		maxAmount,
		minAmount,
	)
	if err != nil {
		return nil, err
	}

	return k.CallEVMWithData(ctx, types.ModuleAddress, &transfer.TransferContractAddress, payload)
}

func (k Keeper) DisableTimeBasedSupplyLimitInTransferContract(
	ctx sdk.Context,
	erc20Address common.Address,
) (
	*evmtypes.MsgEthereumTxResponse,
	error,
) {
	payload, err := transfer.TransferContract.ABI.Pack("disableTimeBasedSupplyLimit", erc20Address)
	if err != nil {
		return nil, err
	}

	return k.CallEVMWithData(ctx, types.ModuleAddress, &transfer.TransferContractAddress, payload)
}
