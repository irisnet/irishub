package gov

import (
	"fmt"

	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/gov/tags"
	"strconv"
	"encoding/json"
	"github.com/irisnet/irishub/modules/gov/params"
	tmstate "github.com/tendermint/tendermint/state"
	protocolKeeper "github.com/irisnet/irishub/app/protocol/keeper"
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
		default:
			errMsg := "Unrecognized gov msg type"
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSubmitProposal(ctx sdk.Context, keeper Keeper, msg MsgSubmitProposal) sdk.Result {
	////////////////////  iris begin  ///////////////////////////
	if msg.ProposalType == ProposalTypeSoftwareHalt {
		_, found := keeper.gk.GetProfiler(ctx, msg.Proposer)
		if !found {
			return ErrNotProfiler(keeper.codespace, msg.Proposer).Result()
		}
	}
	proposal := keeper.NewProposal(ctx, msg.Title, msg.Description, msg.ProposalType, msg.Param)

	////////////////////  iris end  /////////////////////////////

	err, votingStarted := keeper.AddDeposit(ctx, proposal.GetProposalID(), msg.Proposer, msg.InitialDeposit)
	if err != nil {
		return err.Result()
	}
	////////////////////  iris begin  ///////////////////////////
	proposalIDBytes := []byte(strconv.FormatUint(proposal.GetProposalID(), 10))

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
		resTags = resTags.AppendTag(tags.VotingPeriodStart, proposalIDBytes)
	}

	return sdk.Result{
		Data: proposalIDBytes,
		Tags: resTags,
	}
}

func handleMsgSubmitTxTaxUsageProposal(ctx sdk.Context, keeper Keeper, msg MsgSubmitTxTaxUsageProposal) sdk.Result {
	_, found := keeper.gk.GetTrustee(ctx, msg.DestAddress)
	if !found {
		return ErrNotTrustee(keeper.codespace, msg.DestAddress).Result()
	}

	proposal := keeper.NewUsageProposal(ctx, msg)

	err, votingStarted := keeper.AddDeposit(ctx, proposal.GetProposalID(), msg.Proposer, msg.InitialDeposit)
	if err != nil {
		return err.Result()
	}
	proposalIDBytes := []byte(strconv.FormatUint(proposal.GetProposalID(), 10))

	resTags := sdk.NewTags(
		tags.Action, tags.ActionSubmitProposal,
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

	_, found := keeper.gk.GetProfiler(ctx, msg.Proposer)
	if !found {
		return ErrNotProfiler(keeper.codespace, msg.Proposer).Result()
	}

	if msg.ProposalType == ProposalTypeSoftwareUpgrade {
		emptyUpgradeConfig := protocolKeeper.UpgradeConfig{}
		if keeper.pk.GetUpgradeConfig(ctx) != emptyUpgradeConfig {
			return ErrSwitchPeriodInProcess(keeper.codespace).Result()
		}
	}


	proposal := keeper.NewSoftwareUpgradeProposal(ctx, msg)

	err, votingStarted := keeper.AddDeposit(ctx, proposal.GetProposalID(), msg.Proposer, msg.InitialDeposit)
	if err != nil {
		return err.Result()
	}
	proposalIDBytes := []byte(strconv.FormatUint(proposal.GetProposalID(), 10))

	resTags := sdk.NewTags(
		tags.Action, tags.ActionSubmitProposal,
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

	////////////////////  iris begin  ///////////////////////////
	proposalIDBytes := []byte(strconv.FormatUint(msg.ProposalID, 10))
	////////////////////  iris end  /////////////////////////////

	// TODO: Add tag for if voting period started
	resTags := sdk.NewTags(
		tags.Action, tags.ActionDeposit,
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

	////////////////////  iris begin  ///////////////////////////
	proposalIDBytes := []byte(strconv.FormatUint(msg.ProposalID, 10))
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

	logger := ctx.Logger().With("module", "gov")

	resTags = sdk.NewTags()

	if ctx.BlockHeight() == keeper.GetTerminatorHeight(ctx) {
		resTags = resTags.AppendTag(tmstate.HaltTagKey, []byte(tmstate.HaltTagValue))
		logger.Info(fmt.Sprintf("Terminator Start!!!"))
	}

	inactiveIterator := keeper.InactiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	for ; inactiveIterator.Valid(); inactiveIterator.Next() {
		var proposalID uint64
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(inactiveIterator.Value(), &proposalID)
		inactiveProposal := keeper.GetProposal(ctx, proposalID)
		keeper.DeleteProposal(ctx, proposalID)
		keeper.DeleteDeposits(ctx, proposalID) // delete any associated deposits (burned)

		resTags = resTags.AppendTag(tags.Action, tags.ActionProposalDropped)
		resTags = resTags.AppendTag(tags.ProposalID, []byte(string(proposalID)))

		logger.Info(
			fmt.Sprintf("proposal %d (%s) didn't meet minimum deposit of %s (had only %s); deleted",
				inactiveProposal.GetProposalID(),
				inactiveProposal.GetTitle(),
				govparams.GetDepositProcedure(ctx).MinDeposit,
				inactiveProposal.GetTotalDeposit(),
			),
		)
	}
	inactiveIterator.Close()

	activeIterator := keeper.ActiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	for ; activeIterator.Valid(); activeIterator.Next() {
		var proposalID uint64
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(activeIterator.Value(), &proposalID)
		activeProposal := keeper.GetProposal(ctx, proposalID)
		passes, tallyResults := tally(ctx, keeper, activeProposal)

		var action []byte
		if passes {
			keeper.RefundDeposits(ctx, activeProposal.GetProposalID())
			activeProposal.SetStatus(StatusPassed)
			action = tags.ActionProposalPassed
			activeProposal.Execute(ctx, keeper)
		} else {
			keeper.DeleteDeposits(ctx, activeProposal.GetProposalID())
			activeProposal.SetStatus(StatusRejected)
			action = tags.ActionProposalRejected
		}
		activeProposal.SetTallyResult(tallyResults)
		keeper.SetProposal(ctx, activeProposal)

		keeper.RemoveFromActiveProposalQueue(ctx, activeProposal.GetVotingEndTime(), activeProposal.GetProposalID())

		logger.Info(fmt.Sprintf("proposal %d (%s) tallied; passed: %v",
			activeProposal.GetProposalID(), activeProposal.GetTitle(), passes))

		resTags = resTags.AppendTag(tags.Action, action)
		resTags = resTags.AppendTag(tags.ProposalID, []byte(string(proposalID)))
	}
	activeIterator.Close()

	return resTags
}
