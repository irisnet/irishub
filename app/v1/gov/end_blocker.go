package gov

import (
	"strconv"

	"github.com/irisnet/irishub/app/v1/gov/tags"
	sdk "github.com/irisnet/irishub/types"
	tmstate "github.com/tendermint/tendermint/state"
)

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
		inactiveProposal.GetProposalLevel().SubProposalNum(ctx, keeper)
		keeper.DeleteDeposits(ctx, proposalID)
		keeper.DeleteProposal(ctx, proposalID)

		resTags = resTags.AppendTag(tags.Action, tags.ActionProposalDropped)
		resTags = resTags.AppendTag(tags.ProposalID, []byte(string(proposalID)))

		keeper.RemoveFromInactiveProposalQueue(ctx, inactiveProposal.GetDepositEndTime(), inactiveProposal.GetProposalID())
		ctx.Logger().Info("Proposal didn't meet minimum deposit; deleted", "ProposalID",
			inactiveProposal.GetProposalID(), "MinDeposit", inactiveProposal.GetProposalLevel().GetDepositProcedure(ctx, keeper).MinDeposit,
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
			activeProposal.Execute(ctx, keeper)
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
		ctx.Logger().Info("Proposal tallied", "ProposalID", activeProposal.GetProposalID(), "result", result, "tallyResults", tallyResults)
		resTags = resTags.AppendTag(tags.Action, action)
		resTags = resTags.AppendTag(tags.ProposalID, []byte(string(proposalID)))

		for _, valAddr := range keeper.GetValidatorSet(ctx, proposalID) {
			if _, ok := votingVals[valAddr.String()]; !ok {
				val := keeper.ds.GetValidatorSet().Validator(ctx, valAddr)
				if val != nil && val.GetStatus() == sdk.Bonded {
					keeper.ds.GetValidatorSet().Slash(ctx,
						val.GetConsAddr(),
						ctx.BlockHeight(),
						val.GetPower().RoundInt64(),
						activeProposal.GetProposalLevel().GetTallyingProcedure(ctx, keeper).Penalty)
				}
			}
		}

		activeProposal.GetProposalLevel().SubProposalNum(ctx, keeper)
		keeper.DeleteValidatorSet(ctx, activeProposal.GetProposalID())
	}
	return resTags
}
