package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"

	"github.com/teleport-network/teleport/x/xibc/core/client/client/cli"
)

func EmptyRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "xibc",
		Handler:  nil,
	}
}

var (
	CreateClientProposalHandler    = govclient.NewProposalHandler(cli.NewCreateClientProposalCmd, EmptyRESTHandler)
	UpgradeClientProposalHandler   = govclient.NewProposalHandler(cli.NewUpgradeClientProposalCmd, EmptyRESTHandler)
	ToggleClientProposalHandler    = govclient.NewProposalHandler(cli.NewToggleClientProposalCmd, EmptyRESTHandler)
	RegisterRelayerProposalHandler = govclient.NewProposalHandler(cli.NewRegisterRelayerProposalCmd, EmptyRESTHandler)
)
