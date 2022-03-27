package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/teleport-network/teleport/x/xibc/core/client/types"
	"github.com/teleport-network/teleport/x/xibc/exported"
)

func (k Keeper) HandleCreateClient(ctx sdk.Context, p *types.CreateClientProposal) (exported.ClientState, error) {
	if _, has := k.GetClientState(ctx, p.ChainName); has {
		return nil, sdkerrors.Wrapf(types.ErrClientExists, "chain-name: %s", p.ChainName)
	}

	clientState, err := types.UnpackClientState(p.ClientState)
	if err != nil {
		return nil, err
	}

	consensusState, err := types.UnpackConsensusState(p.ConsensusState)
	if err != nil {
		return nil, err
	}

	if err := k.CreateClient(ctx, p.ChainName, clientState, consensusState); err != nil {
		return nil, err
	}

	return clientState, nil
}

func (k Keeper) HandleUpgradeClient(ctx sdk.Context, p *types.UpgradeClientProposal) (exported.ClientState, error) {
	clientState, err := types.UnpackClientState(p.ClientState)
	if err != nil {
		return nil, err
	}

	consensusState, err := types.UnpackConsensusState(p.ConsensusState)
	if err != nil {
		return nil, err
	}

	if err := k.UpgradeClient(ctx, p.ChainName, clientState, consensusState); err != nil {
		return nil, err
	}

	return clientState, nil
}

func (k Keeper) HandleToggleClient(ctx sdk.Context, p *types.ToggleClientProposal) (exported.ClientState, error) {
	if _, has := k.GetClientState(ctx, p.ChainName); !has {
		return nil, sdkerrors.Wrapf(types.ErrClientNotFound, "chain-name: %s", p.ChainName)
	}

	clientState, err := types.UnpackClientState(p.ClientState)
	if err != nil {
		return nil, err
	}

	consensusState, err := types.UnpackConsensusState(p.ConsensusState)
	if err != nil {
		return nil, err
	}

	if err := k.ToggleClient(ctx, p.ChainName, clientState, consensusState); err != nil {
		return nil, err
	}

	return clientState, nil
}

func (k Keeper) HandleRegisterRelayer(ctx sdk.Context, p *types.RegisterRelayerProposal) error {
	// for _, chain := range p.Chains {
	// 	if _, has := k.GetClientState(ctx, chain); !has {
	// 		return sdkerrors.Wrapf(types.ErrClientNotFound, "chain-name: %s", chain)
	// 	}
	// }

	k.RegisterRelayers(ctx, p.Address, p.Chains, p.Addresses)

	return nil
}
