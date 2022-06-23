package types

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	ethermint "github.com/tharsis/ethermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
)

// constants
const (
	ProposalTypeRegisterCoin                string = "RegisterCoin"
	ProposalTypeAddCoin                     string = "AddCoin"
	ProposalTypeRegisterERC20               string = "RegisterERC20"
	ProposalTypeToggleTokenConversion       string = "ToggleTokenConversion" // #nosec
	ProposalTypeRegisterERC20Trace          string = "RegisterERC20Trace"
	ProposalTypeEnableTimeBasedSupplyLimit  string = "EnableTimeBasedSupplyLimit"
	ProposalTypeDisableTimeBasedSupplyLimit string = "DisableTimeBasedSupplyLimit"
)

// Implements Proposal Interface
var (
	_ govtypes.Content = &RegisterCoinProposal{}
	_ govtypes.Content = &AddCoinProposal{}
	_ govtypes.Content = &RegisterERC20Proposal{}
	_ govtypes.Content = &ToggleTokenConversionProposal{}
	_ govtypes.Content = &RegisterERC20TraceProposal{}
	_ govtypes.Content = &EnableTimeBasedSupplyLimitProposal{}
	_ govtypes.Content = &DisableTimeBasedSupplyLimitProposal{}
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeRegisterCoin)
	govtypes.RegisterProposalType(ProposalTypeAddCoin)
	govtypes.RegisterProposalType(ProposalTypeRegisterERC20)
	govtypes.RegisterProposalType(ProposalTypeToggleTokenConversion)
	govtypes.RegisterProposalType(ProposalTypeRegisterERC20Trace)
	govtypes.RegisterProposalType(ProposalTypeEnableTimeBasedSupplyLimit)
	govtypes.RegisterProposalType(ProposalTypeDisableTimeBasedSupplyLimit)

	govtypes.RegisterProposalTypeCodec(&RegisterCoinProposal{}, "aggregate/RegisterCoinProposal")
	govtypes.RegisterProposalTypeCodec(&AddCoinProposal{}, "aggregate/AddCoinProposal")
	govtypes.RegisterProposalTypeCodec(&RegisterERC20Proposal{}, "aggregate/RegisterERC20Proposal")
	govtypes.RegisterProposalTypeCodec(&ToggleTokenConversionProposal{}, "aggregate/ToggleTokenConversionProposal")
	govtypes.RegisterProposalTypeCodec(&RegisterERC20TraceProposal{}, "aggregate/RegisterERC20TraceProposal")
	govtypes.RegisterProposalTypeCodec(&EnableTimeBasedSupplyLimitProposal{}, "aggregate/EnableTimeBasedSupplyLimitProposal")
	govtypes.RegisterProposalTypeCodec(&DisableTimeBasedSupplyLimitProposal{}, "aggregate/DisableTimeBasedSupplyLimitProposal")
}

// CreateDenomDescription generates a string with the coin description
func CreateDenomDescription(address string) string {
	return fmt.Sprintf("Cosmos coin token representation of %s", address)
}

// CreateDenom generates a string the module name plus the address to avoid conflicts with names staring with a number
func CreateDenom(address string) string {
	return fmt.Sprintf("%s/%s", ModuleName, address)
}

// ================================================================================================================

// NewRegisterCoinProposal returns new instance of RegisterCoinProposal
func NewRegisterCoinProposal(title, description string, coinMetadata banktypes.Metadata) govtypes.Content {
	return &RegisterCoinProposal{
		Title:       title,
		Description: description,
		Metadata:    coinMetadata,
	}
}

// ProposalRoute returns router key for this proposal
func (*RegisterCoinProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns proposal type for this proposal
func (*RegisterCoinProposal) ProposalType() string {
	return ProposalTypeRegisterCoin
}

// ValidateBasic performs a stateless check of the proposal fields
func (p *RegisterCoinProposal) ValidateBasic() error {
	if err := p.Metadata.Validate(); err != nil {
		return err
	}

	if err := ibctransfertypes.ValidateIBCDenom(p.Metadata.Base); err != nil {
		return err
	}

	if err := validateIBCVoucherMetadata(p.Metadata); err != nil {
		return err
	}

	return govtypes.ValidateAbstract(p)
}

// validateIBCVoucherMetadata checks that the coin metadata fields are consistent
// with an IBC voucher denomination.
func validateIBCVoucherMetadata(metadata banktypes.Metadata) error {
	// Check ibc/ denom
	denomSplit := strings.SplitN(metadata.Base, "/", 2)
	if denomSplit[0] == metadata.Base && strings.TrimSpace(metadata.Base) != "" {
		// Not IBC
		return nil
	}
	if len(denomSplit) != 2 || denomSplit[0] != ibctransfertypes.DenomPrefix {
		// NOTE: should be unaccessible (covered on ValidateIBCDenom)
		return fmt.Errorf("invalid metadata. %s denomination should be prefixed with the format 'ibc/", metadata.Base)
	}
	return nil
}

// ValidateAggregateDenom checks if a denom is a valid 'aggregate/' denomination
func ValidateAggregateDenom(denom string) error {
	denomSplit := strings.SplitN(denom, "/", 2)

	if len(denomSplit) != 2 || denomSplit[0] != ModuleName {
		return fmt.Errorf("invalid denom. %s denomination should be prefixed with the format 'aggregate/", denom)
	}

	return ethermint.ValidateAddress(denomSplit[1])
}

// ================================================================================================================

// NewAddCoinProposal returns new instance of AddCoinProposal
func NewAddCoinProposal(title, description string, coinMetadata banktypes.Metadata, contractAddr string) govtypes.Content {
	return &AddCoinProposal{
		Title:           title,
		Description:     description,
		Metadata:        coinMetadata,
		ContractAddress: contractAddr,
	}
}

// ProposalRoute returns router key for this proposal
func (*AddCoinProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns proposal type for this proposal
func (*AddCoinProposal) ProposalType() string {
	return ProposalTypeAddCoin
}

// ValidateBasic performs a stateless check of the proposal fields
func (rtbp *AddCoinProposal) ValidateBasic() error {
	if err := rtbp.Metadata.Validate(); err != nil {
		return err
	}
	if err := ibctransfertypes.ValidateIBCDenom(rtbp.Metadata.Base); err != nil {
		return err
	}
	if err := validateIBCVoucherMetadata(rtbp.Metadata); err != nil {
		return err
	}
	if check := common.IsHexAddress(rtbp.ContractAddress); !check {
		return sdkerrors.Wrap(ErrAggregateDisabled, "ERC20 address invalid")
	}
	return govtypes.ValidateAbstract(rtbp)
}

// ================================================================================================================

// NewRegisterERC20Proposal returns new instance of RegisterERC20Proposal
func NewRegisterERC20Proposal(title, description, erc20Addr string) govtypes.Content {
	return &RegisterERC20Proposal{
		Title:        title,
		Description:  description,
		ERC20Address: erc20Addr,
	}
}

// ProposalRoute returns router key for this proposal
func (*RegisterERC20Proposal) ProposalRoute() string { return RouterKey }

// ProposalType returns proposal type for this proposal
func (*RegisterERC20Proposal) ProposalType() string {
	return ProposalTypeRegisterERC20
}

// ValidateBasic performs a stateless check of the proposal fields
func (p *RegisterERC20Proposal) ValidateBasic() error {
	if err := ethermint.ValidateAddress(p.ERC20Address); err != nil {
		return sdkerrors.Wrap(err, "ERC20 address")
	}
	return govtypes.ValidateAbstract(p)
}

// ================================================================================================================

// NewToggleTokenConversionProposal returns new instance of ToggleTokenConversionProposal
func NewToggleTokenConversionProposal(title, description string, token string) govtypes.Content {
	return &ToggleTokenConversionProposal{
		Title:       title,
		Description: description,
		Token:       token,
	}
}

// ProposalRoute returns router key for this proposal
func (*ToggleTokenConversionProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns proposal type for this proposal
func (*ToggleTokenConversionProposal) ProposalType() string {
	return ProposalTypeToggleTokenConversion
}

// ValidateBasic performs a stateless check of the proposal fields
func (ttcp *ToggleTokenConversionProposal) ValidateBasic() error {
	// check if the token is a hex address, if not, check if it is a valid SDK denom
	if err := ethermint.ValidateAddress(ttcp.Token); err != nil {
		if err := sdk.ValidateDenom(ttcp.Token); err != nil {
			return err
		}
	}
	return govtypes.ValidateAbstract(ttcp)
}

// ================================================================================================================

// NewRegisterERC20TraceProposal returns new instance of RegisterERC20Proposal
func NewRegisterERC20TraceProposal(
	title string,
	description string,
	erc20Addr string,
	originToken string,
	originChain string,
	scale uint64,
) govtypes.Content {
	return &RegisterERC20TraceProposal{
		Title:        title,
		Description:  description,
		ERC20Address: erc20Addr,
		OriginToken:  originToken,
		OriginChain:  originChain,
		Scale:        scale,
	}
}

// ProposalRoute returns router key for this proposal
func (*RegisterERC20TraceProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns proposal type for this proposal
func (*RegisterERC20TraceProposal) ProposalType() string {
	return ProposalTypeRegisterERC20Trace
}

// ValidateBasic performs a stateless check of the proposal fields
func (p *RegisterERC20TraceProposal) ValidateBasic() error {
	if err := ethermint.ValidateAddress(p.ERC20Address); err != nil {
		return sdkerrors.Wrap(err, "ERC20 address")
	}

	// TODO: validate originToken
	if len(strings.TrimSpace(p.OriginToken)) == 0 {
		return sdkerrors.Wrap(ErrInvalidOriginToken, "originToken cannot be blank")
	}

	// TODO: validate originChain
	if len(strings.TrimSpace(p.OriginChain)) == 0 {
		return sdkerrors.Wrap(ErrInvalidOriginChain, "originChain cannot be blank")
	}

	if p.Scale > 18 {
		return sdkerrors.Wrap(ErrERC20TraceScale, "ERC20 trace scale should be smaller than 18")
	}

	return govtypes.ValidateAbstract(p)
}

// ================================================================================================================

// NewRegisterERC20TraceProposal returns new instance of RegisterERC20Proposal
func NewEnableTimeBasedSupplyLimitProposal(
	title string,
	description string,
	erc20Address string,
	timePeriod string,
	timeBasedLimit string,
	maxAmount string,
	minAmount string,
) govtypes.Content {
	return &EnableTimeBasedSupplyLimitProposal{
		Title:          title,
		Description:    description,
		ERC20Address:   erc20Address,
		TimePeriod:     timePeriod,
		TimeBasedLimit: timeBasedLimit,
		MaxAmount:      maxAmount,
		MinAmount:      minAmount,
	}
}

// ProposalRoute returns router key for this proposal
func (*EnableTimeBasedSupplyLimitProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns proposal type for this proposal
func (*EnableTimeBasedSupplyLimitProposal) ProposalType() string {
	return ProposalTypeEnableTimeBasedSupplyLimit
}

// ValidateBasic performs a stateless check of the proposal fields
func (p *EnableTimeBasedSupplyLimitProposal) ValidateBasic() error {
	if err := ethermint.ValidateAddress(p.ERC20Address); err != nil {
		return sdkerrors.Wrap(err, "ERC20 address")
	}

	timePeriod, valid := new(big.Int).SetString(p.TimePeriod, 10)
	if !valid || timePeriod.Cmp(big.NewInt(0)) <= 0 {
		return sdkerrors.Wrapf(ErrInvalidTimePeriod, "timePeriod: %s", p.TimePeriod)
	}

	minAmount, valid := new(big.Int).SetString(p.MinAmount, 10)
	if !valid || minAmount.Cmp(big.NewInt(0)) <= 0 {
		return sdkerrors.Wrapf(ErrInvalidMinAmount, "minAmount: %s", p.TimePeriod)
	}

	maxAmount, valid := new(big.Int).SetString(p.MaxAmount, 10)
	if !valid || maxAmount.Cmp(minAmount) <= 0 {
		return sdkerrors.Wrapf(ErrInvalidMaxAmount, "maxAmount: %s", p.TimePeriod)
	}

	timeBasedLimit, valid := new(big.Int).SetString(p.TimeBasedLimit, 10)
	if !valid || timeBasedLimit.Cmp(maxAmount) <= 0 {
		return sdkerrors.Wrapf(ErrInvalidTimeBasedLimit, "timeBasedLimit: %s", p.TimePeriod)
	}

	return govtypes.ValidateAbstract(p)
}

// ================================================================================================================

// NewRegisterERC20TraceProposal returns new instance of RegisterERC20Proposal
func NewDisableTimeBasedSupplyLimitProposal(
	title string,
	description string,
	erc20Address string,
) govtypes.Content {
	return &DisableTimeBasedSupplyLimitProposal{
		Title:        title,
		Description:  description,
		ERC20Address: erc20Address,
	}
}

// ProposalRoute returns router key for this proposal
func (*DisableTimeBasedSupplyLimitProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns proposal type for this proposal
func (*DisableTimeBasedSupplyLimitProposal) ProposalType() string {
	return ProposalTypeDisableTimeBasedSupplyLimit
}

// ValidateBasic performs a stateless check of the proposal fields
func (p *DisableTimeBasedSupplyLimitProposal) ValidateBasic() error {
	if err := ethermint.ValidateAddress(p.ERC20Address); err != nil {
		return sdkerrors.Wrap(err, "ERC20 address")
	}
	return govtypes.ValidateAbstract(p)
}
