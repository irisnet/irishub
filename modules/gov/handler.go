package gov

import (
	"encoding/json"
	"strconv"

	"github.com/irisnet/irishub/modules/gov/tags"
	sdk "github.com/irisnet/irishub/types"
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
	if num, ok := keeper.HasReachedTheMaxProposalNum(ctx, proposalLevel); ok {
		return ErrMoreThanMaxProposal(keeper.codespace, num, proposalLevel.string()).Result()
	}

	if msg.ProposalType == ProposalTypeSystemHalt {
		_, found := keeper.guardianKeeper.GetProfiler(ctx, msg.Proposer)
		if !found {
			return ErrNotProfiler(keeper.codespace, msg.Proposer).Result()
		}
	}

	if msg.ProposalType == ProposalTypeParameterChange {
		for _, param := range msg.Params {
			if p, ok := keeper.paramsKeeper.GetParamSet(param.Subspace); ok {
				if _, err := p.Validate(param.Key, param.Value); err != nil {
					return err.Result()
				}
			} else {
				return ErrInvalidParam(DefaultCodespace, param.Subspace).Result()
			}
		}
	}
	proposal := keeper.NewProposal(ctx, msg.Title, msg.Description, msg.ProposalType, msg.Params)

	err, votingStarted := keeper.AddInitialDeposit(ctx, proposal, msg.Proposer, msg.InitialDeposit)
	if err != nil {
		return err.Result()
	}

	proposalIDBytes := []byte(strconv.FormatUint(proposal.GetProposalID(), 10))

	var paramBytes []byte
	if msg.ProposalType == ProposalTypeParameterChange {
		paramBytes, _ = json.Marshal(proposal.(*ParameterProposal).Params)
	}

	resTags := sdk.NewTags(
		tags.Proposer, []byte(msg.Proposer.String()),
		tags.ProposalID, proposalIDBytes,

		tags.Param, paramBytes,
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
	if num, ok := keeper.HasReachedTheMaxProposalNum(ctx, proposalLevel); ok {
		return ErrMoreThanMaxProposal(keeper.codespace, num, proposalLevel.string()).Result()
	}

	if msg.Usage != UsageTypeBurn {
		_, found := keeper.guardianKeeper.GetTrustee(ctx, msg.DestAddress)
		if !found {
			return ErrNotTrustee(keeper.codespace, msg.DestAddress).Result()
		}
	}

	proposal := keeper.NewUsageProposal(ctx, msg)

	err, votingStarted := keeper.AddInitialDeposit(ctx, proposal, msg.Proposer, msg.InitialDeposit)
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

	keeper.AddProposalNum(ctx, proposal)
	return sdk.Result{
		Data: proposalIDBytes,
		Tags: resTags,
	}
}

func handleMsgSubmitSoftwareUpgradeProposal(ctx sdk.Context, keeper Keeper, msg MsgSubmitSoftwareUpgradeProposal) sdk.Result {
	proposalLevel := GetProposalLevelByProposalKind(msg.ProposalType)
	if num, ok := keeper.HasReachedTheMaxProposalNum(ctx, proposalLevel); ok {
		return ErrMoreThanMaxProposal(keeper.codespace, num, proposalLevel.string()).Result()
	}

	if !keeper.protocolKeeper.IsValidVersion(ctx, msg.Version) {
		return ErrCodeInvalidVersion(keeper.codespace, msg.Version).Result()
	}

	if uint64(ctx.BlockHeight()) > msg.SwitchHeight {
		return ErrCodeInvalidSwitchHeight(keeper.codespace, uint64(ctx.BlockHeight()), msg.SwitchHeight).Result()
	}
	_, found := keeper.guardianKeeper.GetProfiler(ctx, msg.Proposer)
	if !found {
		return ErrNotProfiler(keeper.codespace, msg.Proposer).Result()
	}

	if _, ok := keeper.protocolKeeper.GetUpgradeConfig(ctx); ok {
		return ErrSwitchPeriodInProcess(keeper.codespace).Result()
	}

	proposal := keeper.NewSoftwareUpgradeProposal(ctx, msg)

	err, votingStarted := keeper.AddInitialDeposit(ctx, proposal, msg.Proposer, msg.InitialDeposit)
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

// Called every block, process inflation, update validator set
func EndBlocker(ctx sdk.Context, keeper Keeper) (resTags sdk.Tags) {
	ctx = ctx.WithCoinFlowTrigger(sdk.GovEndBlocker)
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "endBlock").With("module", "iris/gov"))
	resTags = sdk.NewTags()

	if ctx.BlockHeight() == keeper.GetSystemHaltHeight(ctx) {
		resTags = resTags.AppendTag(tmstate.HaltTagKey, []byte(tmstate.HaltTagValue))
		ctx.Logger().Info("SystemHalt Start!!!")
	}

	inactiveIterator := keeper.InactiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	defer inactiveIterator.Close()
	for ; inactiveIterator.Valid(); inactiveIterator.Next() {
		var proposalID uint64
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(inactiveIterator.Value(), &proposalID)
		inactiveProposal := keeper.GetProposal(ctx, proposalID)
		keeper.SubProposalNum(ctx, inactiveProposal)
		keeper.DeleteDeposits(ctx, proposalID)
		keeper.DeleteProposal(ctx, proposalID)

		resTags = resTags.AppendTag(tags.Action, tags.ActionProposalDropped)
		resTags = resTags.AppendTag(tags.ProposalID, []byte(string(proposalID)))

		keeper.RemoveFromInactiveProposalQueue(ctx, inactiveProposal.GetDepositEndTime(), inactiveProposal.GetProposalID())
		ctx.Logger().Info("Proposal didn't meet minimum deposit; deleted", "ProposalID",
			inactiveProposal.GetProposalID(), "MinDeposit", keeper.GetDepositProcedure(ctx, inactiveProposal).MinDeposit,
			"ActualDeposit", inactiveProposal.GetTotalDeposit(),
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
			keeper.metrics.ProposalStatus.With(ProposalIDLabel, strconv.FormatUint(proposalID, 10)).Set(2)
			keeper.RefundDeposits(ctx, activeProposal.GetProposalID())
			activeProposal.SetStatus(StatusPassed)
			action = tags.ActionProposalPassed
			Execute(ctx, keeper, activeProposal)
		} else if result == REJECT {
			keeper.metrics.ProposalStatus.With(ProposalIDLabel, strconv.FormatUint(proposalID, 10)).Set(3)
			keeper.RefundDeposits(ctx, activeProposal.GetProposalID())
			activeProposal.SetStatus(StatusRejected)
			action = tags.ActionProposalRejected
		} else if result == REJECTVETO {
			keeper.metrics.ProposalStatus.With(ProposalIDLabel, strconv.FormatUint(proposalID, 10)).Set(3)
			keeper.DeleteDeposits(ctx, activeProposal.GetProposalID())
			activeProposal.SetStatus(StatusRejected)
			action = tags.ActionProposalRejected
		}
		keeper.RemoveFromActiveProposalQueue(ctx, activeProposal.GetVotingEndTime(), activeProposal.GetProposalID())
		activeProposal.SetTallyResult(tallyResults)
		keeper.SetProposal(ctx, activeProposal)
		ctx.Logger().Info("Proposal tallied", "ProposalID", activeProposal.GetProposalID(), "result", result)
		resTags = resTags.AppendTag(tags.Action, action)
		resTags = resTags.AppendTag(tags.ProposalID, []byte(strconv.FormatUint(proposalID, 10)))

		for _, valAddr := range keeper.GetValidatorSet(ctx, proposalID) {
			if _, ok := votingVals[valAddr.String()]; !ok {
				val := keeper.ds.GetValidatorSet().Validator(ctx, valAddr)
				if val != nil && val.GetStatus() == sdk.Bonded {
					keeper.ds.GetValidatorSet().Slash(ctx,
						val.GetConsAddr(),
						ctx.BlockHeight(),
						val.GetPower().RoundInt64(),
						keeper.GetTallyingProcedure(ctx, activeProposal).Penalty)
				}
			}
		}

		keeper.SubProposalNum(ctx, activeProposal)
		keeper.DeleteValidatorSet(ctx, activeProposal.GetProposalID())
	}
	return resTags
}
