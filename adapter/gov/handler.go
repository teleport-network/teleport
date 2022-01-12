package gov

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/x/gov/types"

	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/teleport-network/teleport/adapter/common"
	"github.com/teleport-network/teleport/syscontracts"
	"github.com/teleport-network/teleport/syscontracts/gov"
)

func (h *HookAdapter) HandleVoted(ctx sdk.Context, log *ethtypes.Log) error {
	event := new(gov.GovVoted)
	if err := syscontracts.ParseLog(event, h.abi, log, "Voted"); err != nil {
		return err
	}
	voter, err := bech32.ConvertAndEncode(sdk.GetConfig().GetBech32AccountAddrPrefix(), event.Voter.Bytes())
	if err != nil {
		return err
	}
	msg := &types.MsgVote{
		ProposalId: event.ProposalId,
		Voter:      voter,
		Option:     types.VoteOption(event.VoteOption),
	}
	return common.ExecuteMsg(ctx, h.router, msg)
}

func (h *HookAdapter) HandleVotedWeighted(ctx sdk.Context, log *ethtypes.Log) error {
	event := new(gov.GovVotedWeighted)
	if err := syscontracts.ParseLog(event, h.abi, log, "VotedWeighted"); err != nil {
		return err
	}

	if len(event.Options) == 0 {
		return errors.New("Must have options ")
	}
	sdkWeightedVoteOption := make([]types.WeightedVoteOption, 0, len(event.Options))
	for _, option := range event.Options {
		sdkWeightedVoteOption = append(sdkWeightedVoteOption, types.WeightedVoteOption{
			Option: types.VoteOption(option.Option),
			Weight: sdk.NewDecWithPrec(int64(option.Weight), 2),
		})
	}
	voter, err := bech32.ConvertAndEncode(sdk.GetConfig().GetBech32AccountAddrPrefix(), event.Voter.Bytes())
	if err != nil {
		return err
	}

	msg := &types.MsgVoteWeighted{
		ProposalId: event.ProposalId,
		Voter:      voter,
		Options:    sdkWeightedVoteOption,
	}
	return common.ExecuteMsg(ctx, h.router, msg)
}
