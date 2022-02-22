package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/teleport-network/teleport/x/xibc/core/host"
	"github.com/teleport-network/teleport/x/xibc/exported"
)

const (
	ProposalTypeClientCreate    = "CreateClient"
	ProposalTypeClientUpgrade   = "UpgradeClient"
	ProposalTypeClientToggle    = "ToggleClient"
	ProposalTypeRelayerRegister = "RegisterRelayer"
)

var (
	_ govtypes.Content = &CreateClientProposal{}
	_ govtypes.Content = &UpgradeClientProposal{}
	_ govtypes.Content = &ToggleClientProposal{}
	_ govtypes.Content = &RegisterRelayerProposal{}

	_ codectypes.UnpackInterfacesMessage = &CreateClientProposal{}
	_ codectypes.UnpackInterfacesMessage = &UpgradeClientProposal{}
	_ codectypes.UnpackInterfacesMessage = &ToggleClientProposal{}
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeClientCreate)
	govtypes.RegisterProposalType(ProposalTypeClientUpgrade)
	govtypes.RegisterProposalType(ProposalTypeClientToggle)
	govtypes.RegisterProposalType(ProposalTypeRelayerRegister)
}

// NewCreateClientProposal creates a new creating client proposal.
func NewCreateClientProposal(
	title string,
	description string,
	chainName string,
	clientState exported.ClientState,
	consensusState exported.ConsensusState,
) (
	*CreateClientProposal, error,
) {
	clientStateAny, err := PackClientState(clientState)
	if err != nil {
		return nil, err
	}

	consensusStateAny, err := PackConsensusState(consensusState)
	if err != nil {
		return nil, err
	}

	return &CreateClientProposal{
		Title:          title,
		Description:    description,
		ChainName:      chainName,
		ClientState:    clientStateAny,
		ConsensusState: consensusStateAny,
	}, nil
}

// GetTitle returns the title of a client update proposal.
func (cup *CreateClientProposal) GetTitle() string { return cup.Title }

// GetDescription returns the description of a client update proposal.
func (cup *CreateClientProposal) GetDescription() string { return cup.Description }

// ProposalRoute returns the routing key of a client update proposal.
func (cup *CreateClientProposal) ProposalRoute() string { return GovRouterKey }

// ProposalType returns the type of a client update proposal.
func (cup *CreateClientProposal) ProposalType() string { return ProposalTypeClientCreate }

// ValidateBasic runs basic stateless validity checks
func (cup *CreateClientProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(cup); err != nil {
		return err
	}

	if err := host.ClientIdentifierValidator(cup.ChainName); err != nil {
		return err
	}

	clientState, err := UnpackClientState(cup.ClientState)
	if err != nil {
		return err
	}

	return clientState.Validate()
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (cup CreateClientProposal) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if err := unpacker.UnpackAny(cup.ClientState, new(exported.ClientState)); err != nil {
		return err
	}
	return unpacker.UnpackAny(cup.ConsensusState, new(exported.ConsensusState))
}

// NewUpgradeClientProposal create a upgrade client proposal.
func NewUpgradeClientProposal(
	title string,
	description string,
	chainName string,
	clientState exported.ClientState,
	consensusState exported.ConsensusState,
) (
	*UpgradeClientProposal, error,
) {
	clientStateAny, err := PackClientState(clientState)
	if err != nil {
		return nil, err
	}

	consensusStateAny, err := PackConsensusState(consensusState)
	if err != nil {
		return nil, err
	}

	return &UpgradeClientProposal{
		Title:          title,
		Description:    description,
		ChainName:      chainName,
		ClientState:    clientStateAny,
		ConsensusState: consensusStateAny,
	}, nil
}

// GetTitle returns the title of a client upgrade proposal.
func (cup *UpgradeClientProposal) GetTitle() string { return cup.Title }

// GetDescription returns the description of a client upgrade proposal.
func (cup *UpgradeClientProposal) GetDescription() string { return cup.Description }

// ProposalRoute returns the routing key of a client upgrade proposal.
func (cup *UpgradeClientProposal) ProposalRoute() string { return GovRouterKey }

// ProposalType returns the type of a client upgrade proposal.
func (cup *UpgradeClientProposal) ProposalType() string { return ProposalTypeClientUpgrade }

// ValidateBasic runs basic stateless validity checks
func (cup *UpgradeClientProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(cup); err != nil {
		return err
	}

	if err := host.ClientIdentifierValidator(cup.ChainName); err != nil {
		return err
	}

	clientState, err := UnpackClientState(cup.ClientState)
	if err != nil {
		return err
	}

	return clientState.Validate()
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (cup UpgradeClientProposal) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if err := unpacker.UnpackAny(cup.ClientState, new(exported.ClientState)); err != nil {
		return err
	}
	return unpacker.UnpackAny(cup.ConsensusState, new(exported.ConsensusState))
}

// NewToggleClientProposal creates a new creating client proposal.
func NewToggleClientProposal(
	title string,
	description string,
	chainName string,
	clientState exported.ClientState,
	consensusState exported.ConsensusState,
) (
	*ToggleClientProposal, error,
) {
	clientStateAny, err := PackClientState(clientState)
	if err != nil {
		return nil, err
	}

	consensusStateAny, err := PackConsensusState(consensusState)
	if err != nil {
		return nil, err
	}

	return &ToggleClientProposal{
		Title:          title,
		Description:    description,
		ChainName:      chainName,
		ClientState:    clientStateAny,
		ConsensusState: consensusStateAny,
	}, nil
}

// GetTitle returns the title of a client toggle proposal.
func (cup *ToggleClientProposal) GetTitle() string { return cup.Title }

// GetDescription returns the description of a client toggle proposal.
func (cup *ToggleClientProposal) GetDescription() string { return cup.Description }

// ProposalRoute returns the routing key of a client toggle proposal.
func (cup *ToggleClientProposal) ProposalRoute() string { return GovRouterKey }

// ProposalType returns the type of a client toggle proposal.
func (cup *ToggleClientProposal) ProposalType() string { return ProposalTypeClientToggle }

// ValidateBasic runs basic stateless validity checks
func (cup *ToggleClientProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(cup); err != nil {
		return err
	}

	if err := host.ClientIdentifierValidator(cup.ChainName); err != nil {
		return err
	}

	clientState, err := UnpackClientState(cup.ClientState)
	if err != nil {
		return err
	}

	return clientState.Validate()
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (cup ToggleClientProposal) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if err := unpacker.UnpackAny(cup.ClientState, new(exported.ClientState)); err != nil {
		return err
	}
	return unpacker.UnpackAny(cup.ConsensusState, new(exported.ConsensusState))
}

// NewRegisterRelayerProposal creates a new registering relayer proposal.
func NewRegisterRelayerProposal(
	title string,
	description string,
	chainName string,
	relayers []string,
) *RegisterRelayerProposal {
	return &RegisterRelayerProposal{
		Title:       title,
		Description: description,
		ChainName:   chainName,
		Relayers:    relayers,
	}
}

// GetTitle returns the title of a registering relayer proposal.
func (rrp *RegisterRelayerProposal) GetTitle() string { return rrp.Title }

// GetDescription returns the description of a registering relayer proposal.
func (rrp *RegisterRelayerProposal) GetDescription() string { return rrp.Description }

// ProposalRoute returns the routing key of a registering relayer proposal.
func (rrp *RegisterRelayerProposal) ProposalRoute() string { return GovRouterKey }

// ProposalType returns the type of a client registering relayer proposal.
func (rrp *RegisterRelayerProposal) ProposalType() string { return ProposalTypeRelayerRegister }

// ValidateBasic runs basic stateless validity checks
func (rrp *RegisterRelayerProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(rrp); err != nil {
		return err
	}

	if err := host.ClientIdentifierValidator(rrp.ChainName); err != nil {
		return err
	}

	if len(rrp.Relayers) == 0 {
		return govtypes.ErrInvalidLengthGov
	}

	for _, relayer := range rrp.Relayers {
		if _, err := sdk.AccAddressFromBech32(relayer); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "string could not be parsed as address: %v", err)
		}
	}

	return nil
}
