package gov

import (
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/types/gov/params"
	govtypes "github.com/irisnet/irishub/types/gov"
)

// validatorGovInfo used for tallying
type validatorGovInfo struct {
	Address         sdk.ValAddress // address of the validator operator
	Power           sdk.Dec        // Power of a Validator
	DelegatorShares sdk.Dec        // Total outstanding delegator shares
	Minus           sdk.Dec        // Minus of validator, used to compute validator's voting power
	Vote            govtypes.VoteOption     // Vote of the validator
}

func tally(ctx sdk.Context, keeper Keeper, proposal govtypes.Proposal) (passes bool, tallyResults govtypes.TallyResult) {
	results := make(map[govtypes.VoteOption]sdk.Dec)
	results[govtypes.OptionYes] = sdk.ZeroDec()
	results[govtypes.OptionAbstain] = sdk.ZeroDec()
	results[govtypes.OptionNo] = sdk.ZeroDec()
	results[govtypes.OptionNoWithVeto] = sdk.ZeroDec()

	totalVotingPower := sdk.ZeroDec()
	systemVotingPower := sdk.ZeroDec()
	currValidators := make(map[string]validatorGovInfo)

	keeper.vs.IterateBondedValidatorsByPower(ctx, func(index int64, validator sdk.Validator) (stop bool) {
		currValidators[validator.GetOperator().String()] = validatorGovInfo{
			Address:         validator.GetOperator(),
			Power:           validator.GetPower(),
			DelegatorShares: validator.GetDelegatorShares(),
			Minus:           sdk.ZeroDec(),
			Vote:            govtypes.OptionEmpty,
		}
		systemVotingPower = systemVotingPower.Add(validator.GetPower())
		return false
	})

	// iterate over all the votes
	votesIterator := keeper.GetVotes(ctx, proposal.GetProposalID())
	defer votesIterator.Close()
	for ; votesIterator.Valid(); votesIterator.Next() {
		vote := &govtypes.Vote{}
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(votesIterator.Value(), vote)

		// if validator, just record it in the map
		// if delegator tally voting power
		valAddrStr := sdk.ValAddress(vote.Voter).String()
		if val, ok := currValidators[valAddrStr]; ok {
			val.Vote = vote.Option
			currValidators[valAddrStr] = val
		} else {

			keeper.ds.IterateDelegations(ctx, vote.Voter, func(index int64, delegation sdk.Delegation) (stop bool) {
				valAddrStr := delegation.GetValidatorAddr().String()

				if val, ok := currValidators[valAddrStr]; ok {
					val.Minus = val.Minus.Add(delegation.GetShares())
					currValidators[valAddrStr] = val

					delegatorShare := delegation.GetShares().Quo(val.DelegatorShares)
					votingPower := val.Power.Mul(delegatorShare)

					results[vote.Option] = results[vote.Option].Add(votingPower)
					totalVotingPower = totalVotingPower.Add(votingPower)
				}

				return false
			})
		}

		keeper.deleteVote(ctx, vote.ProposalID, vote.Voter)
	}

	// iterate over the validators again to tally their voting power
	for _, val := range currValidators {
		if val.Vote == govtypes.OptionEmpty {
			continue
		}

		sharesAfterMinus := val.DelegatorShares.Sub(val.Minus)
		percentAfterMinus := sharesAfterMinus.Quo(val.DelegatorShares)
		votingPower := val.Power.Mul(percentAfterMinus)

		results[val.Vote] = results[val.Vote].Add(votingPower)
		totalVotingPower = totalVotingPower.Add(votingPower)
	}

	////////////////////  iris begin  ///////////////////////////
	tallyingProcedure := govparams.GetTallyingProcedure(ctx)
	////////////////////  iris end  /////////////////////////////

	tallyResults = govtypes.TallyResult{
		Yes:        results[govtypes.OptionYes],
		Abstain:    results[govtypes.OptionAbstain],
		No:         results[govtypes.OptionNo],
		NoWithVeto: results[govtypes.OptionNoWithVeto],
	}

	// If no one votes, proposal fails
	if totalVotingPower.Sub(results[govtypes.OptionAbstain]).Equal(sdk.ZeroDec()) {
		return false, tallyResults
	}
	////////////////////  iris begin  ///////////////////////////
	//if more than 1/3 of voters abstain, proposal fails
	if tallyingProcedure.Participation.GT(totalVotingPower.Quo(systemVotingPower)) {
		return false, tallyResults
	}
	////////////////////  iris end  ///////////////////////////

	// If more than 1/3 of voters veto, proposal fails
	if results[govtypes.OptionNoWithVeto].Quo(totalVotingPower).GT(tallyingProcedure.Veto) {
		return false, tallyResults
	}

	// If more than 1/2 of non-abstaining voters vote Yes, proposal passes
	if results[govtypes.OptionYes].Quo(totalVotingPower.Sub(results[govtypes.OptionAbstain])).GT(tallyingProcedure.Threshold) {
		return true, tallyResults
	}
	// If more than 1/2 of non-abstaining voters vote No, proposal fails

	return false, tallyResults
}
