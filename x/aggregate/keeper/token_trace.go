package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

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
		return nil, sdkerrors.Wrap(
			types.ErrWritingEthTxData,
			sdkerrors.Wrap(err, "failed to create transaction payload").Error(),
		)
	}

	return k.CallEVMWithData(ctx, types.ModuleAddress, &transfer.TransferContractAddress, payload)
}
