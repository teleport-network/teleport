package keeper

import (
	"encoding/hex"

	"github.com/armon/go-metrics"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/teleport-network/teleport/x/xibc/core/client/types"
	"github.com/teleport-network/teleport/x/xibc/exported"
)

// CreateClient creates a new client state and populates it with a given
// client state and consensus state
func (k Keeper) CreateClient(
	ctx sdk.Context,
	chainName string,
	clientState exported.ClientState,
	consensusState exported.ConsensusState,
) error {
	k.SetClientState(ctx, chainName, clientState)
	// verifies initial consensus state against client state and initializes client store with any client-specific metadata
	// e.g. set ProcessedTime in Tendermint clients
	if err := clientState.Initialize(ctx, k.cdc, k.ClientStore(ctx, chainName), consensusState); err != nil {
		return err
	}

	// check if consensus state is nil in case the created client is Localhost
	k.SetClientConsensusState(ctx, chainName, clientState.GetLatestHeight(), consensusState)
	k.Logger(ctx).Info(
		"client created at height",
		"chain-name", chainName,
		"height", clientState.GetLatestHeight().String(),
	)

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{"xibc", "client", "create"},
			1,
			[]metrics.Label{
				telemetry.NewLabel(types.LabelClientType, clientState.ClientType()),
				telemetry.NewLabel(types.LabelChainName, chainName),
			},
		)
	}()

	return nil
}

// UpgradeClient upgrades the client to a new client state.
func (k Keeper) UpgradeClient(
	ctx sdk.Context,
	chainName string,
	newClientState exported.ClientState,
	newConsensusState exported.ConsensusState,
) error {
	clientState, found := k.GetClientState(ctx, chainName)
	if !found {
		return sdkerrors.Wrapf(types.ErrClientNotFound, "cannot upgrade client %s", chainName)
	}

	if clientState.ClientType() != newClientState.ClientType() {
		return sdkerrors.Wrapf(types.ErrInvalidClientType, "cannot upgrade client %s, client-type not match", chainName)
	}

	if err := newClientState.UpgradeState(ctx, k.cdc, k.ClientStore(ctx, chainName), newConsensusState); err != nil {
		return sdkerrors.Wrapf(types.ErrUpgradeClient, "cannot upgrade client %s", chainName)
	}

	k.SetClientState(ctx, chainName, newClientState)
	k.SetClientConsensusState(ctx, chainName, newClientState.GetLatestHeight(), newConsensusState)

	k.Logger(ctx).Info(
		"client state upgraded",
		"chain-name", chainName,
		"height", newClientState.GetLatestHeight().String(),
	)

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{"xibc", "client", "upgrade"},
			1,
			[]metrics.Label{
				telemetry.NewLabel(types.LabelClientType, newClientState.ClientType()),
				telemetry.NewLabel(types.LabelChainName, chainName),
			},
		)
	}()

	return nil
}

// ToggleClient creates a new client with different type
func (k Keeper) ToggleClient(
	ctx sdk.Context,
	chainName string,
	newClientState exported.ClientState,
	newConsensusState exported.ConsensusState,
) error {
	clientState, found := k.GetClientState(ctx, chainName)
	if !found {
		return sdkerrors.Wrapf(types.ErrClientNotFound, "cannot toggle client %s", chainName)
	}

	if clientState.ClientType() != newClientState.ClientType() {
		return sdkerrors.Wrapf(types.ErrInvalidClientType, "cannot toggle client %s, client-type can't be the same", chainName)
	}

	k.SetClientState(ctx, chainName, newClientState)
	if err := clientState.Initialize(ctx, k.cdc, k.ClientStore(ctx, chainName), newConsensusState); err != nil {
		return err
	}
	k.SetClientConsensusState(ctx, chainName, newClientState.GetLatestHeight(), newConsensusState)

	k.Logger(ctx).Info(
		"client state toggled",
		"chain-name", chainName,
		"height", newClientState.GetLatestHeight().String(),
	)

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{"xibc", "client", "toggle"},
			1,
			[]metrics.Label{
				telemetry.NewLabel(types.LabelClientType, newClientState.ClientType()),
				telemetry.NewLabel(types.LabelChainName, chainName),
			},
		)
	}()

	return nil
}

// UpdateClient updates the consensus state and the state root from a provided header.
func (k Keeper) UpdateClient(
	ctx sdk.Context,
	chainName string,
	header exported.Header,
) error {
	clientState, found := k.GetClientState(ctx, chainName)
	if !found {
		return sdkerrors.Wrapf(types.ErrClientNotFound, "cannot update client %s", chainName)
	}

	clientStore := k.ClientStore(ctx, chainName)
	if status := clientState.Status(ctx, clientStore, k.cdc); status != exported.Active {
		return sdkerrors.Wrapf(types.ErrClientNotActive, "cannot update client (%s) with status %s", chainName, status)
	}

	// Any writes made in CheckHeaderAndUpdateState are persisted on both valid updates
	// Light client implementations are responsible for writing the correct metadata (if any) in either case.
	newClientState, newConsensusState, err := clientState.CheckHeaderAndUpdateState(ctx, k.cdc, k.ClientStore(ctx, chainName), header)
	if err != nil {
		return sdkerrors.Wrapf(err, "cannot update client %s", chainName)
	}

	// set new client state regardless of if update is valid update
	k.SetClientState(ctx, chainName, newClientState)

	// set new consensus state regardless of if update is valid update
	var consensusHeight = header.GetHeight()
	k.SetClientConsensusState(ctx, chainName, header.GetHeight(), newConsensusState)
	k.Logger(ctx).Info(
		"client state updated",
		"chain-name", chainName,
		"height", consensusHeight.String(),
	)

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{"xibc", "client", "update"},
			1,
			[]metrics.Label{
				telemetry.NewLabel(types.LabelClientType, clientState.ClientType()),
				telemetry.NewLabel(types.LabelChainName, chainName),
				telemetry.NewLabel(types.LabelUpdateType, "msg"),
			},
		)
	}()

	_ = ctx.EventManager().EmitTypedEvent(&types.EventUpdateClient{
		ChainName:       chainName,
		ClientType:      clientState.ClientType(),
		ConsensusHeight: consensusHeight.String(),
		Header:          hex.EncodeToString(types.MustMarshalHeader(k.cdc, header)),
	})

	return nil
}
