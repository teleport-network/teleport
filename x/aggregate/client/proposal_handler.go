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
	AddCoinProposalHandler            = govclient.NewProposalHandler(cli.NewAddCoinProposalCmd, EmptyProposalRESTHandler)
	RegisterCoinProposalHandler       = govclient.NewProposalHandler(cli.NewRegisterCoinProposalCmd, EmptyProposalRESTHandler)
	RegisterERC20PairProposalHandler  = govclient.NewProposalHandler(cli.NewRegisterERC20ProposalCmd, EmptyProposalRESTHandler)
	ToggleTokenRelayProposalHandler   = govclient.NewProposalHandler(cli.NewToggleTokenRelayProposalCmd, EmptyProposalRESTHandler)
	UpdateTokenPairERC20Proposal      = govclient.NewProposalHandler(cli.NewUpdateTokenPairERC20ProposalCmd, EmptyProposalRESTHandler)
	RegisterERC20TraceProposalHandler = govclient.NewProposalHandler(cli.NewRegisterERC20TraceProposalCmd, EmptyProposalRESTHandler)
)
