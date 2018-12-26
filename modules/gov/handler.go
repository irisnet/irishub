package gov

import (
	"fmt"

	"encoding/json"
	"strconv"

	"github.com/irisnet/irishub/modules/gov/tags"
	sdk "github.com/irisnet/irishub/types"
	govtypes "github.com/irisnet/irishub/types/gov"
	tmstate "github.com/tendermint/tendermint/state"
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

	proposalLevel := GetProposalLevelByProposalKind(msg.ProposalType)
	if num, ok := keeper.IsMoreThanMaxProposal(ctx, proposalLevel); ok {
		return govtypes.ErrMoreThanMaxProposal(keeper.codespace, num, proposalLevel.string()).Result()
	}
	////////////////////  iris begin  ///////////////////////////
	if msg.ProposalType == govtypes.ProposalTypeSoftwareHalt {
		_, found := keeper.gk.GetProfiler(ctx, msg.Proposer)
		if !found {
			return govtypes.ErrNotProfiler(keeper.codespace, msg.Proposer).Result()
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
	if msg.ProposalType == govtypes.ProposalTypeParameterChange {
		paramBytes, _ = json.Marshal(proposal.(*govtypes.ParameterProposal).Param)
	}
	////////////////////  iris end  /////////////////////////////
	resTags := sdk.NewTags(
		tags.Proposer, []byte(msg.Proposer.String()),
		tags.ProposalID, proposalIDBytes,
		////////////////////  iris begin  ///////////////////////////
		tags.Param, paramBytes,
		////////////////////  iris end  /////////////////////////////
	)

	if votingStarted {
		resTags = resTags.AppendTag(tags.VotingPeriodStart, proposalIDBytes)
	}

	keeper.AddProposalNum(ctx, proposal)
	return sdk.Result{
		Data: proposalIDBytes,
		Tags: resTags,
	}
}

func handleMsgSubmitTxTaxUsageProposal(ctx sdk.Context, keeper Keeper, msg MsgSubmitTxTaxUsageProposal) sdk.Result {
	proposalLevel := GetProposalLevelByProposalKind(msg.ProposalType)
	if num, ok := keeper.IsMoreThanMaxProposal(ctx, proposalLevel); ok {
		return govtypes.ErrMoreThanMaxProposal(keeper.codespace, num, proposalLevel.string()).Result()
	}

	if msg.Usage != govtypes.UsageTypeBurn {
		_, found := keeper.gk.GetTrustee(ctx, msg.DestAddress)
		if !found {
			return govtypes.ErrNotTrustee(keeper.codespace, msg.DestAddress).Result()
		}
	}

	proposal := keeper.NewUsageProposal(ctx, msg)

	err, votingStarted := keeper.AddDeposit(ctx, proposal.GetProposalID(), msg.Proposer, msg.InitialDeposit)
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

	if msg.Usage != govtypes.UsageTypeBurn {
		resTags = resTags.AppendTag(tags.DestAddress, []byte(msg.DestAddress.String()))
	}

	if votingStarted {
		resTags = resTags.AppendTag(tags.VotingPeriodStart, proposalIDBytes)
	}

	keeper.AddProposalNum(ctx, proposal)
	return sdk.Result{
		Data: proposalIDBytes,
		Tags: resTags,
	}
}

func handleMsgSubmitSoftwareUpgradeProposal(ctx sdk.Context, keeper Keeper, msg MsgSubmitSoftwareUpgradeProposal) sdk.Result {
	proposalLevel := GetProposalLevelByProposalKind(msg.ProposalType)
	if num, ok := keeper.IsMoreThanMaxProposal(ctx, proposalLevel); ok {
		return govtypes.ErrMoreThanMaxProposal(keeper.codespace, num, proposalLevel.string()).Result()
	}

	if !keeper.pk.IsValidProtocolVersion(ctx, msg.Version) {
		return govtypes.ErrCodeInvalidVersion(keeper.codespace, msg.Version).Result()
	}

	if uint64(ctx.BlockHeight()) > msg.SwitchHeight {
		return govtypes.ErrCodeInvalidSwitchHeight(keeper.codespace, uint64(ctx.BlockHeight()), msg.SwitchHeight).Result()
	}
	_, found := keeper.gk.GetProfiler(ctx, msg.Proposer)
	if !found {
		return govtypes.ErrNotProfiler(keeper.codespace, msg.Proposer).Result()
	}

	if _, ok := keeper.pk.GetUpgradeConfig(ctx); ok {
		return govtypes.ErrSwitchPeriodInProcess(keeper.codespace).Result()
	}

	proposal := keeper.NewSoftwareUpgradeProposal(ctx, msg)

	err, votingStarted := keeper.AddDeposit(ctx, proposal.GetProposalID(), msg.Proposer, msg.InitialDeposit)
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

	keeper.AddProposalNum(ctx, proposal)
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
	defer inactiveIterator.Close()
	for ; inactiveIterator.Valid(); inactiveIterator.Next() {
		var proposalID uint64
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(inactiveIterator.Value(), &proposalID)
		inactiveProposal := keeper.GetProposal(ctx, proposalID)
		keeper.SubProposalNum(ctx, inactiveProposal)
		keeper.RefundDeposits(ctx, proposalID)
		keeper.DeleteProposal(ctx, proposalID)

		resTags = resTags.AppendTag(tags.Action, tags.ActionProposalDropped)
		resTags = resTags.AppendTag(tags.ProposalID, []byte(string(proposalID)))

		keeper.RemoveFromInactiveProposalQueue(ctx, inactiveProposal.GetDepositEndTime(), inactiveProposal.GetProposalID())
		logger.Info(
			fmt.Sprintf("proposal %d (%s) didn't meet minimum deposit of %s (had only %s); deleted",
				inactiveProposal.GetProposalID(),
				inactiveProposal.GetTitle(),
				GetMinDeposit(ctx, inactiveProposal),
				inactiveProposal.GetTotalDeposit(),
			),
		)
	}

	activeIterator := keeper.ActiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	defer activeIterator.Close()
	for ; activeIterator.Valid(); activeIterator.Next() {
		var proposalID uint64
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(activeIterator.Value(), &proposalID)
		activeProposal := keeper.GetProposal(ctx, proposalID)
		result, tallyResults, votingVals := tally(ctx, keeper, activeProposal)

		var action []byte
		if result == PASS {
			keeper.RefundDeposits(ctx, activeProposal.GetProposalID())
			activeProposal.SetStatus(govtypes.StatusPassed)
			action = tags.ActionProposalPassed
			Execute(ctx, keeper, activeProposal)
		} else if result == REJECT {
			keeper.RefundDeposits(ctx, activeProposal.GetProposalID())
			activeProposal.SetStatus(govtypes.StatusRejected)
			action = tags.ActionProposalRejected
		} else if result == REJECTVETO {
			keeper.DeleteDeposits(ctx, activeProposal.GetProposalID())
			activeProposal.SetStatus(govtypes.StatusRejected)
			action = tags.ActionProposalRejected
		}
		keeper.RemoveFromActiveProposalQueue(ctx, activeProposal.GetVotingEndTime(), activeProposal.GetProposalID())
		activeProposal.SetTallyResult(tallyResults)
		keeper.SetProposal(ctx, activeProposal)
		logger.Info(fmt.Sprintf("proposal %d (%s) tallied; result: %v",
			activeProposal.GetProposalID(), activeProposal.GetTitle(), result))

		resTags = resTags.AppendTag(tags.Action, action)
		resTags = resTags.AppendTag(tags.ProposalID, []byte(string(proposalID)))

		for _, valAddr := range keeper.GetValidatorSet(ctx, proposalID) {
			if _, ok := votingVals[valAddr.String()]; !ok {
				val := keeper.ds.GetValidatorSet().Validator(ctx, valAddr)
				if val != nil {
					keeper.ds.GetValidatorSet().Slash(ctx,
						val.GetConsAddr(),
						ctx.BlockHeight(),
						val.GetPower().RoundInt64(),
						GetTallyingCondition(ctx, activeProposal).Penalty)
				}
			}
		}

		keeper.SubProposalNum(ctx, activeProposal)
		keeper.DeleteValidatorSet(ctx, activeProposal.GetProposalID())
	}
	return resTags
}
