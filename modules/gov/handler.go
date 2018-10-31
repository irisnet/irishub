package gov

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/gov/tags"
	"github.com/irisnet/irishub/modules/gov/params"
	"strconv"
	"encoding/json"
)

// Handle all "gov" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgDeposit:
			return handleMsgDeposit(ctx, keeper, msg)
		case MsgSubmitProposal:
			return handleMsgSubmitProposal(ctx, keeper, msg)
		case MsgVote:
			return handleMsgVote(ctx, keeper, msg)
		default:
			errMsg := "Unrecognized gov msg type"
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSubmitProposal(ctx sdk.Context, keeper Keeper, msg MsgSubmitProposal) sdk.Result {
	////////////////////  iris begin  ///////////////////////////
	proposal := keeper.NewProposal(ctx, msg.Title, msg.Description, msg.ProposalType,msg.Param)
	////////////////////  iris end  /////////////////////////////


	err, votingStarted := keeper.AddDeposit(ctx, proposal.GetProposalID(), msg.Proposer, msg.InitialDeposit)
	if err != nil {
		return err.Result()
	}
	////////////////////  iris begin  ///////////////////////////
	proposalIDBytes := []byte(strconv.FormatInt(proposal.GetProposalID(), 10))

	var paramBytes []byte
	if msg.ProposalType == ProposalTypeParameterChange {
		paramBytes, _ = json.Marshal(proposal.(*ParameterProposal).Param)
	}
	////////////////////  iris end  /////////////////////////////
	resTags := sdk.NewTags(
		tags.Action, tags.ActionSubmitProposal,
		tags.Proposer, []byte(msg.Proposer.String()),
		tags.ProposalID, proposalIDBytes,
		////////////////////  iris begin  ///////////////////////////
		tags.Param, paramBytes,
		////////////////////  iris end  /////////////////////////////
	)

	if votingStarted {
		resTags.AppendTag(tags.VotingPeriodStart, proposalIDBytes)
	}

	return sdk.Result{
		Data: proposalIDBytes,
		Tags: resTags,
	}
}

func handleMsgDeposit(ctx sdk.Context, keeper Keeper, msg MsgDeposit) sdk.Result {

	err, votingStarted := keeper.AddDeposit(ctx, msg.ProposalID, msg.Depositer, msg.Amount)
	if err != nil {
		return err.Result()
	}

	////////////////////  iris begin  ///////////////////////////
	proposalIDBytes := []byte(strconv.FormatInt(msg.ProposalID, 10))
	////////////////////  iris end  /////////////////////////////


	// TODO: Add tag for if voting period started
	resTags := sdk.NewTags(
		tags.Action, tags.ActionDeposit,
		tags.Depositer, []byte(msg.Depositer.String()),
		tags.ProposalID, proposalIDBytes,
	)

	if votingStarted {
		resTags.AppendTag(tags.VotingPeriodStart, proposalIDBytes)
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

	////////////////////  iris begin  ///////////////////////////
	proposalIDBytes := []byte(strconv.FormatInt(msg.ProposalID, 10))
	////////////////////  iris end  /////////////////////////////

	resTags := sdk.NewTags(
		tags.Action, tags.ActionVote,
		tags.Voter, []byte(msg.Voter.String()),
		tags.ProposalID, proposalIDBytes,
	)
	return sdk.Result{
		Tags: resTags,
	}
}

// Called every block, process inflation, update validator set
func EndBlocker(ctx sdk.Context, keeper Keeper) (resTags sdk.Tags) {

	logger := ctx.Logger().With("module", "x/gov")

	resTags = sdk.NewTags()

	// Delete proposals that haven't met minDeposit
	for shouldPopInactiveProposalQueue(ctx, keeper) {
		inactiveProposal := keeper.InactiveProposalQueuePop(ctx)
		if inactiveProposal.GetStatus() != StatusDepositPeriod {
			continue
		}
		////////////////////  iris begin  ///////////////////////////
		proposalIDBytes := []byte(strconv.FormatInt(inactiveProposal.GetProposalID(), 10))
		////////////////////  iris end  /////////////////////////////
		keeper.DeleteProposal(ctx, inactiveProposal)
		resTags.AppendTag(tags.Action, tags.ActionProposalDropped)
		resTags.AppendTag(tags.ProposalID, proposalIDBytes)

		logger.Info(
			fmt.Sprintf("proposal %d (%s) didn't meet minimum deposit of %v iris-atto (had only %v iris-atto); deleted",
				inactiveProposal.GetProposalID(),
				inactiveProposal.GetTitle(),
				////////////////////  iris begin  ///////////////////////////
				govparams.GetDepositProcedure(ctx).MinDeposit.AmountOf("iris-atto"),
				////////////////////  iris end  /////////////////////////////
				inactiveProposal.GetTotalDeposit().AmountOf("iris-atto"),
			),
		)
	}

	// Check if earliest Active Proposal ended voting period yet
	for shouldPopActiveProposalQueue(ctx, keeper) {
		activeProposal := keeper.ActiveProposalQueuePop(ctx)

		proposalStartTime := activeProposal.GetVotingStartTime()
		////////////////////  iris begin  ///////////////////////////
		votingPeriod := govparams.GetVotingProcedure(ctx).VotingPeriod
		////////////////////  iris end  /////////////////////////////
		if ctx.BlockHeader().Time.Before(proposalStartTime.Add(votingPeriod)) {
			continue
		}

		passes, tallyResults := tally(ctx, keeper, activeProposal)
		////////////////////  iris begin  ///////////////////////////
		proposalIDBytes := []byte(strconv.FormatInt(activeProposal.GetProposalID(), 10))
		////////////////////  iris end  /////////////////////////////
		var action []byte
		if passes {
			keeper.RefundDeposits(ctx, activeProposal.GetProposalID())
			activeProposal.SetStatus(StatusPassed)
			action = tags.ActionProposalPassed
			////////////////////  iris begin  ///////////////////////////
			activeProposal.Execute(ctx, keeper)
			////////////////////  iris end  /////////////////////////////
		} else {
			keeper.DeleteDeposits(ctx, activeProposal.GetProposalID())
			activeProposal.SetStatus(StatusRejected)
			action = tags.ActionProposalRejected
		}
		activeProposal.SetTallyResult(tallyResults)
		keeper.SetProposal(ctx, activeProposal)

		logger.Info(fmt.Sprintf("proposal %d (%s) tallied; passed: %v",
			activeProposal.GetProposalID(), activeProposal.GetTitle(), passes))

		resTags.AppendTag(tags.Action, action)
		resTags.AppendTag(tags.ProposalID, proposalIDBytes)
	}

	return resTags
}
func shouldPopInactiveProposalQueue(ctx sdk.Context, keeper Keeper) bool {
	////////////////////  iris begin  ///////////////////////////
	depositProcedure := govparams.GetDepositProcedure(ctx)
	////////////////////  iris end  /////////////////////////////
	peekProposal := keeper.InactiveProposalQueuePeek(ctx)

	if peekProposal == nil {
		return false
	} else if peekProposal.GetStatus() != StatusDepositPeriod {
		return true
	} else if !ctx.BlockHeader().Time.Before(peekProposal.GetSubmitTime().Add(depositProcedure.MaxDepositPeriod)) {
		return true
	}
	return false
}

func shouldPopActiveProposalQueue(ctx sdk.Context, keeper Keeper) bool {
	////////////////////  iris begin  ///////////////////////////
	votingProcedure := govparams.GetVotingProcedure(ctx)
	////////////////////  iris end  /////////////////////////////
	peekProposal := keeper.ActiveProposalQueuePeek(ctx)

	if peekProposal == nil {
		return false
	} else if !ctx.BlockHeader().Time.Before(peekProposal.GetVotingStartTime().Add(votingProcedure.VotingPeriod)) {
		return true
	}
	return false
}
