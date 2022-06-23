package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/cosmos/cosmos-sdk/x/gov/client/rest"

	"github.com/teleport-network/teleport/x/aggregate/client/cli"
)

func EmptyProposalRESTHandler(client.Context) rest.ProposalRESTHandler {
	return rest.ProposalRESTHandler{}
}

var (
	RegisterCoinProposalHandler                = govclient.NewProposalHandler(cli.NewRegisterCoinProposalCmd, EmptyProposalRESTHandler)
	AddCoinProposalHandler                     = govclient.NewProposalHandler(cli.NewAddCoinProposalCmd, EmptyProposalRESTHandler)
	RegisterERC20PairProposalHandler           = govclient.NewProposalHandler(cli.NewRegisterERC20ProposalCmd, EmptyProposalRESTHandler)
	ToggleTokenConversionProposalHandler       = govclient.NewProposalHandler(cli.NewToggleTokenConversionProposalCmd, EmptyProposalRESTHandler)
	RegisterERC20TraceProposalHandler          = govclient.NewProposalHandler(cli.NewRegisterERC20TraceProposalCmd, EmptyProposalRESTHandler)
	EnableTimeBasedSupplyLimitProposalHandler  = govclient.NewProposalHandler(cli.NewEnableTimeBasedSupplyLimitProposalCmd, EmptyProposalRESTHandler)
	DisableTimeBasedSupplyLimitProposalHandler = govclient.NewProposalHandler(cli.NewDisableTimeBasedSupplyLimitProposalCmd, EmptyProposalRESTHandler)
)
