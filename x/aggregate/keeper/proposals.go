package keeper

import (
	"bytes"
	"fmt"
	"math/big"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/ethereum/go-ethereum/common"

	"github.com/teleport-network/teleport/x/aggregate/types"
)

// RegisterCoin deploys an erc20 contract and creates the token pair for the existing cosmos coin
func (k Keeper) RegisterCoin(ctx sdk.Context, coinMetadata banktypes.Metadata) (*types.TokenPair, error) {
	// Check if the conversion is globally enabled
	params := k.GetParams(ctx)
	if !params.EnableAggregate {
		return nil, sdkerrors.Wrap(
			types.ErrAggregateDisabled, "registration is currently disabled by governance",
		)
	}

	// Prohibit denominations that contain the evm denom
	if strings.Contains(coinMetadata.Base, "tele") {
		return nil, sdkerrors.Wrapf(
			types.ErrEVMDenom, "cannot register the TELE denomination %s", coinMetadata.Base,
		)
	}

	// Check if denomination is already registered
	if k.IsDenomRegistered(ctx, coinMetadata.Name) {
		return nil, sdkerrors.Wrapf(
			types.ErrTokenPairAlreadyExists, "coin denomination already registered: %s", coinMetadata.Name,
		)
	}

	// Check if the coin exists by ensuring the supply is set
	if !k.bankKeeper.HasSupply(ctx, coinMetadata.Base) {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrInvalidCoins, "base denomination '%s' cannot have a supply of 0", coinMetadata.Base,
		)
	}

	if err := k.verifyMetadata(ctx, coinMetadata); err != nil {
		return nil, sdkerrors.Wrapf(
			types.ErrInternalTokenPair, "coin metadata is invalid %s", coinMetadata.Name,
		)
	}

	addr, err := k.DeployERC20Contract(ctx, coinMetadata)
	if err != nil {
		return nil, sdkerrors.Wrap(
			err, "failed to create wrapped coin denom metadata for ERC20",
		)
	}

	pair := types.NewTokenPair(addr, []string{coinMetadata.Base}, true, types.OWNER_MODULE)
	k.SetTokenPair(ctx, pair)
	k.SetDenomsMap(ctx, pair.Denoms, pair.GetID())
	k.SetERC20Map(ctx, common.HexToAddress(pair.ERC20Address), pair.GetID())

	return &pair, nil
}

// AddCoin add coin to token pair
func (k Keeper) AddCoin(ctx sdk.Context, coinMetadata banktypes.Metadata, contractAddr string) (*types.TokenPair, error) {
	// check if the contract is validate
	if check := common.IsHexAddress(contractAddr); !check {
		return nil, sdkerrors.Wrapf(
			types.ErrAggregateDisabled, "erc20 address %s not valid", contractAddr,
		)
	}
	// check if the conversion is globally enabled
	params := k.GetParams(ctx)
	if !params.EnableAggregate {
		return nil, sdkerrors.Wrap(types.ErrAggregateDisabled, "registration is currently disabled by governance")
	}

	// Prohibit denominations that contain the evm denom
	if strings.Contains(coinMetadata.Base, "tele") {
		return nil, sdkerrors.Wrapf(
			types.ErrEVMDenom, "cannot register the TELE denomination %s", coinMetadata.Base,
		)
	}

	// check if the denomination already registered
	if k.IsDenomRegistered(ctx, coinMetadata.Name) {
		return nil, sdkerrors.Wrapf(
			types.ErrTokenPairAlreadyExists, "coin denomination already registered: %s", coinMetadata.Name,
		)
	}

	// check if the coin exists by ensuring the supply is set
	if !k.bankKeeper.HasSupply(ctx, coinMetadata.Base) {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrInvalidCoins, "base denomination '%s' cannot have a supply of 0", coinMetadata.Base,
		)
	}

	if err := k.verifyMetadata(ctx, coinMetadata); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInternalTokenPair, "coin metadata is invalid %s", coinMetadata.Name)
	}

	id := k.GetERC20Map(ctx, common.HexToAddress(contractAddr))
	pair, found := k.GetTokenPair(ctx, id)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrTokenPairNotFound, "token '%s' not registered", contractAddr)
	}
	pair.Denoms = append(pair.Denoms, coinMetadata.Base)
	// id not change, just overwrite
	if !bytes.Equal(id, pair.GetID()) {
		return nil, sdkerrors.Wrapf(types.ErrTokenPairNotFound, "pair id not equal ,cannot change it")
	}
	k.SetTokenPair(ctx, pair)
	k.SetDenomMap(ctx, coinMetadata.Base, pair.GetID())

	return &pair, nil
}

// verifyMetadata verifies if the metadata matches the existing one, if not it sets it to the store
func (k Keeper) verifyMetadata(ctx sdk.Context, coinMetadata banktypes.Metadata) error {
	meta, found := k.bankKeeper.GetDenomMetaData(ctx, coinMetadata.Base)
	if !found {
		k.bankKeeper.SetDenomMetaData(ctx, coinMetadata)
		return nil
	}
	// If it already existed, check that is equal to what is stored
	return types.EqualMetadata(meta, coinMetadata)
}

// RegisterERC20 creates a cosmos coin and registers the token pair between the coin and the ERC20
func (k Keeper) RegisterERC20(ctx sdk.Context, contract common.Address) (*types.TokenPair, error) {
	params := k.GetParams(ctx)
	if !params.EnableAggregate {
		return nil, sdkerrors.Wrap(types.ErrAggregateDisabled, "registration is currently disabled by governance")
	}

	if k.IsERC20Registered(ctx, contract) {
		return nil, sdkerrors.Wrapf(types.ErrTokenPairAlreadyExists, "token ERC20 contract already registered: %s", contract.String())
	}

	metadata, err := k.CreateCoinMetadata(ctx, contract)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to create wrapped coin denom metadata for ERC20")
	}

	pair := types.NewTokenPair(contract, []string{metadata.Name}, true, types.OWNER_EXTERNAL)
	k.SetTokenPair(ctx, pair)
	k.SetDenomsMap(ctx, pair.Denoms, pair.GetID())
	k.SetERC20Map(ctx, common.HexToAddress(pair.ERC20Address), pair.GetID())
	return &pair, nil
}

// CreateCoinMetadata generates the metadata to represent the ERC20 token on teleport.
func (k Keeper) CreateCoinMetadata(ctx sdk.Context, contract common.Address) (*banktypes.Metadata, error) {
	strContract := contract.String()

	erc20Data, err := k.QueryERC20(ctx, contract)
	if err != nil {
		return nil, err
	}

	// Check if metadata already exists
	if _, found := k.bankKeeper.GetDenomMetaData(ctx, types.CreateDenom(strContract)); found {
		return nil, sdkerrors.Wrap(types.ErrInternalTokenPair, "denom metadata already registered")
	}

	if k.IsDenomRegistered(ctx, types.CreateDenom(strContract)) {
		return nil, sdkerrors.Wrapf(types.ErrInternalTokenPair, "coin denomination already registered: %s", erc20Data.Name)
	}

	// base denomination
	base := types.CreateDenom(strContract)

	// create a bank denom metadata based on the ERC20 token ABI details
	// metadata name is should always be the contract since it's the key
	// to the bank store
	metadata := banktypes.Metadata{
		Description: types.CreateDenomDescription(strContract),
		Base:        base,
		// NOTE: Denom units MUST be increasing
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    base,
				Exponent: 0,
			},
		},
		Name:    types.CreateDenom(strContract),
		Symbol:  erc20Data.Symbol,
		Display: base,
	}

	// only append metadata if decimals > 0, otherwise validation fails
	if erc20Data.Decimals > 0 {
		nameSanitized := types.SanitizeERC20Name(erc20Data.Name)
		metadata.DenomUnits = append(
			metadata.DenomUnits,
			&banktypes.DenomUnit{
				Denom:    nameSanitized,
				Exponent: uint32(erc20Data.Decimals),
			},
		)
		metadata.Display = nameSanitized
	}

	if err := metadata.Validate(); err != nil {
		return nil, sdkerrors.Wrapf(err, "ERC20 token data is invalid for contract %s", strContract)
	}

	k.bankKeeper.SetDenomMetaData(ctx, metadata)

	return &metadata, nil
}

// ToggleConversion toggles conversion for a given token pair
func (k Keeper) ToggleConversion(ctx sdk.Context, token string) (types.TokenPair, error) {
	id := k.GetTokenPairID(ctx, token)
	if len(id) == 0 {
		return types.TokenPair{}, sdkerrors.Wrapf(types.ErrTokenPairNotFound, "token '%s' not registered by id", token)
	}

	pair, found := k.GetTokenPair(ctx, id)
	if !found {
		return types.TokenPair{}, sdkerrors.Wrapf(types.ErrTokenPairNotFound, "token '%s' not registered", token)
	}

	pair.Enabled = !pair.Enabled

	k.SetTokenPair(ctx, pair)
	return pair, nil
}

// RegisterERC20Trace creates a ERC20 trace for XIBC
func (k Keeper) RegisterERC20Trace(
	ctx sdk.Context,
	contract common.Address,
	originToken string,
	originChain string,
	scale uint8,
) error {
	if _, err := k.AddERC20TraceToTransferContract(ctx, contract, originToken, originChain, scale); err != nil {
		return fmt.Errorf("call bindToken failed: %s", err)
	}
	return nil
}

// EnableTimeBasedSupplyLimit enables time based supply limit
func (k Keeper) EnableTimeBasedSupplyLimit(
	ctx sdk.Context,
	erc20Address common.Address,
	timePeriod *big.Int,
	timeBasedLimit *big.Int,
	maxAmount *big.Int,
	minAmount *big.Int,
) error {
	if _, err := k.EnableTimeBasedSupplyLimitInTransferContract(
		ctx,
		erc20Address,
		timePeriod,
		timeBasedLimit,
		maxAmount,
		minAmount,
	); err != nil {
		return fmt.Errorf("call bindToken failed: %s", err)
	}
	return nil
}

// DisableTimeBasedSupplyLimit disables time based supply limit
func (k Keeper) DisableTimeBasedSupplyLimit(
	ctx sdk.Context,
	erc20Address common.Address,
) error {
	if _, err := k.DisableTimeBasedSupplyLimitInTransferContract(ctx, erc20Address); err != nil {
		return fmt.Errorf("call bindToken failed: %s", err)
	}
	return nil
}
