package gov

import (
	"strconv"

	"github.com/irisnet/irishub/app/v1/gov/tags"
	sdk "github.com/irisnet/irishub/types"
)

// Handle all "gov" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSubmitProposal,
			MsgSubmitSoftwareUpgradeProposal,
			MsgSubmitAddTokenProposal,
			MsgSubmitTxTaxUsageProposal:
			return handleMsgSubmitProposal(ctx, keeper, msg)
		case MsgDeposit:
			return handleMsgDeposit(ctx, keeper, msg)
		case MsgVote:
			return handleMsgVote(ctx, keeper, msg)
		default:
			errMsg := "Unrecognized gov msg type"
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSubmitProposal(ctx sdk.Context, keeper Keeper, msg sdk.Msg) sdk.Result {
	resTags, err := keeper.SubmitProposal(ctx, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
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
