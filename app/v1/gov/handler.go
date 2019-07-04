package gov

import (
	"encoding/json"
	"strconv"

	"github.com/irisnet/irishub/app/v1/gov/tags"
	sdk "github.com/irisnet/irishub/types"
)

// Handle all "gov" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgDeposit:
			return handleMsgDeposit(ctx, keeper, msg)
		case MsgSubmitProposal:
			return handleMsgSubmitProposal(ctx, keeper, msg)
		case MsgSubmitTxTaxUsageProposal:
			return handleMsgSubmitTxTaxUsageProposal(ctx, keeper, msg)
		case MsgSubmitSoftwareUpgradeProposal:
			return handleMsgSubmitSoftwareUpgradeProposal(ctx, keeper, msg)
		case MsgVote:
			return handleMsgVote(ctx, keeper, msg)
		case MsgSubmitAddTokenProposal:
			return handleMsgSubmitAddTokenProposal(ctx, keeper, msg)
		default:
			errMsg := "Unrecognized gov msg type"
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSubmitProposal(ctx sdk.Context, keeper Keeper, msg MsgSubmitProposal) sdk.Result {
	proposal, votingStarted, err := keeper.SubmitProposal(ctx, msg)
	if err != nil {
		return err.Result()
	}

	proposalIDBytes := []byte(strconv.FormatUint(proposal.GetProposalID(), 10))

	resTags := sdk.NewTags(
		tags.Proposer, []byte(msg.Proposer.String()),
		tags.ProposalID, proposalIDBytes,
	)

	var paramBytes []byte
	if msg.ProposalType == ProposalTypeParameterChange {
		paramBytes, _ = json.Marshal(msg.Params)
		resTags = resTags.AppendTag(tags.Param, paramBytes)
	}

	if votingStarted {
		resTags = resTags.AppendTag(tags.VotingPeriodStart, proposalIDBytes)
	}
	return sdk.Result{
		Data: proposalIDBytes,
		Tags: resTags,
	}
}

func handleMsgSubmitTxTaxUsageProposal(ctx sdk.Context, keeper Keeper, msg MsgSubmitTxTaxUsageProposal) sdk.Result {

	proposal, votingStarted, err := keeper.SubmitProposal(ctx, msg)
	if err != nil {
		return err.Result()
	}

	proposalIDBytes := []byte(strconv.FormatUint(proposal.GetProposalID(), 10))

	resTags := sdk.NewTags(
		tags.Proposer, []byte(msg.Proposer.String()),
		tags.ProposalID, proposalIDBytes,
		tags.Usage, []byte(msg.Usage.String()),
		tags.Percent, []byte(msg.Percent.String()),
	)

	if msg.Usage != UsageTypeBurn {
		resTags = resTags.AppendTag(tags.DestAddress, []byte(msg.DestAddress.String()))
	}

	if votingStarted {
		resTags = resTags.AppendTag(tags.VotingPeriodStart, proposalIDBytes)
	}
	return sdk.Result{
		Data: proposalIDBytes,
		Tags: resTags,
	}
}

func handleMsgSubmitSoftwareUpgradeProposal(ctx sdk.Context, keeper Keeper, msg MsgSubmitSoftwareUpgradeProposal) sdk.Result {
	proposal, votingStarted, err := keeper.SubmitProposal(ctx, msg)
	if err != nil {
		return err.Result()
	}

	proposalIDBytes := []byte(strconv.FormatUint(proposal.GetProposalID(), 10))

	resTags := sdk.NewTags(
		tags.Proposer, []byte(msg.Proposer.String()),
		tags.ProposalID, proposalIDBytes,
	)

	if votingStarted {
		resTags = resTags.AppendTag(tags.VotingPeriodStart, proposalIDBytes)
	}

	return sdk.Result{
		Data: proposalIDBytes,
		Tags: resTags,
	}
}

func handleMsgDeposit(ctx sdk.Context, keeper Keeper, msg MsgDeposit) sdk.Result {

	err, votingStarted := keeper.AddDeposit(ctx, msg.ProposalID, msg.Depositor, msg.Amount)
	if err != nil {
		return err.Result()
	}

	proposalIDBytes := []byte(strconv.FormatUint(msg.ProposalID, 10))

	// TODO: Add tag for if voting period started
	resTags := sdk.NewTags(
		tags.Depositor, []byte(msg.Depositor.String()),
		tags.ProposalID, proposalIDBytes,
	)

	if votingStarted {
		resTags = resTags.AppendTag(tags.VotingPeriodStart, proposalIDBytes)
	}

	return sdk.Result{
		Tags: resTags,
	}
}

func handleMsgVote(ctx sdk.Context, keeper Keeper, msg MsgVote) sdk.Result {

	err := keeper.AddVote(ctx, msg.ProposalID, msg.Voter, msg.Option)
	if err != nil {
		return err.Result()
	}

	proposalIDBytes := []byte(strconv.FormatUint(msg.ProposalID, 10))

	resTags := sdk.NewTags(
		tags.Voter, []byte(msg.Voter.String()),
		tags.ProposalID, proposalIDBytes,
	)
	return sdk.Result{
		Tags: resTags,
	}
}

func handleMsgSubmitAddTokenProposal(ctx sdk.Context, keeper Keeper, msg MsgSubmitAddTokenProposal) sdk.Result {
	proposal, votingStarted, err := keeper.SubmitProposal(ctx, msg)
	if err != nil {
		return err.Result()
	}

	tokenId := proposal.(*AddTokenProposal).FToken.GetUniqueID()
	proposalIDBytes := []byte(strconv.FormatUint(proposal.GetProposalID(), 10))
	resTags := sdk.NewTags(
		tags.Proposer, []byte(msg.Proposer.String()),
		tags.ProposalID, proposalIDBytes,
		tags.TokenId, []byte(tokenId),
	)

	if votingStarted {
		resTags = resTags.AppendTag(tags.VotingPeriodStart, proposalIDBytes)
	}
	return sdk.Result{
		Data: proposalIDBytes,
		Tags: resTags,
	}
}
