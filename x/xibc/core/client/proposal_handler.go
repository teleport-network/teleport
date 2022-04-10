package client

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/teleport-network/teleport/x/xibc/core/client/keeper"
	"github.com/teleport-network/teleport/x/xibc/core/client/types"
)

// NewClientProposalHandler defines the client manager proposal handler
func NewClientProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.CreateClientProposal:
			return handleCreateClientProposal(ctx, k, c)
		case *types.UpgradeClientProposal:
			return handleUpgradeClientProposal(ctx, k, c)
		case *types.ToggleClientProposal:
			return handleToggleClientProposal(ctx, k, c)
		case *types.RegisterRelayerProposal:
			return handleRegisterRelayerProposal(ctx, k, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized xibc proposal content type: %T", c)
		}
	}
}

// handleCreateClientProposal will try to create the client with the ClientState and ConsensusState
func handleCreateClientProposal(ctx sdk.Context, k keeper.Keeper, p *types.CreateClientProposal) error {
	clientState, err := k.HandleCreateClient(ctx, p)
	if err != nil {
		return err
	}

	_ = ctx.EventManager().EmitTypedEvent(&types.EventCreateClientProposal{
		ChainName:       p.ChainName,
		ClientType:      clientState.ClientType(),
		ConsensusHeight: clientState.GetLatestHeight().String(),
	})

	return nil
}

// handleUpgradeClientProposal will try to update the client with the new ClientState and ConsensusState if and only if the proposal passes
func handleUpgradeClientProposal(ctx sdk.Context, k keeper.Keeper, p *types.UpgradeClientProposal) error {
	upgradedClientState, err := k.HandleUpgradeClient(ctx, p)
	if err != nil {
		return err
	}

	_ = ctx.EventManager().EmitTypedEvent(&types.EventUpgradeClientProposal{
		ChainName:       p.ChainName,
		ClientType:      upgradedClientState.ClientType(),
		ConsensusHeight: upgradedClientState.GetLatestHeight().String(),
	})

	return nil
}

// handleToggleClientProposal will try to toggle the client type between Light and TSS
func handleToggleClientProposal(ctx sdk.Context, k keeper.Keeper, p *types.ToggleClientProposal) error {
	clientState, err := k.HandleToggleClient(ctx, p)
	if err != nil {
		return err
	}

	_ = ctx.EventManager().EmitTypedEvent(&types.EventToggleClientProposal{
		ChainName:       p.ChainName,
		ClientType:      clientState.ClientType(),
		ConsensusHeight: clientState.GetLatestHeight().String(),
	})

	return nil
}

// handleRegisterRelayerProposal will try to save the registered relayer address under the specified client
func handleRegisterRelayerProposal(ctx sdk.Context, k keeper.Keeper, p *types.RegisterRelayerProposal) error {
	if err := k.HandleRegisterRelayer(ctx, p); err != nil {
		return err
	}

	_ = ctx.EventManager().EmitTypedEvent(&types.EventRegisterRelayerProposal{
		Address:   p.Address,
		Chains:    p.Chains,
		Addresses: p.Addresses,
	})

	return nil
}
