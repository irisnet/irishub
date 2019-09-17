package gov

import (
	sdk "github.com/irisnet/irishub/types"
)

type ProposalResult string

const (
	PASS       ProposalResult = "pass"
	REJECT     ProposalResult = "reject"
	REJECTVETO ProposalResult = "reject-veto"
)

// validatorGovInfo used for tallying
type validatorGovInfo struct {
	Address sdk.ValAddress // address of the validator operator
	Power   sdk.Dec        // Power of a Validator
	Vote    VoteOption     // Vote of the validator
}

func tally(ctx sdk.Context, keeper Keeper, proposal Proposal) (result ProposalResult, tallyResults TallyResult, votingVals map[string]bool) {
	results := make(map[VoteOption]sdk.Dec)
	results[OptionYes] = sdk.ZeroDec()
	results[OptionAbstain] = sdk.ZeroDec()
	results[OptionNo] = sdk.ZeroDec()
	results[OptionNoWithVeto] = sdk.ZeroDec()

	totalVotingPower := sdk.ZeroDec()
	systemVotingPower := sdk.ZeroDec()
	currValidators := make(map[string]validatorGovInfo)
	votingVals = make(map[string]bool)
	keeper.vs.IterateBondedValidatorsByPower(ctx, func(index int64, validator sdk.Validator) (stop bool) {
		currValidators[validator.GetOperator().String()] = validatorGovInfo{
			Address: validator.GetOperator(),
			Power:   validator.GetPower(),
			Vote:    OptionEmpty,
		}
		systemVotingPower = systemVotingPower.Add(validator.GetPower())
		return false
	})
	// iterate over all the votes
	votesIterator := keeper.GetVotes(ctx, proposal.GetProposalID())
	defer votesIterator.Close()
	for ; votesIterator.Valid(); votesIterator.Next() {
		vote := &Vote{}
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(votesIterator.Value(), vote)

		// if validator, just record it in the map
		valAddrStr := sdk.ValAddress(vote.Voter).String()
		if val, ok := currValidators[valAddrStr]; ok {
			val.Vote = vote.Option
			results[val.Vote] = results[val.Vote].Add(val.Power)
			totalVotingPower = totalVotingPower.Add(val.Power)
			votingVals[valAddrStr] = true
		}
	}

	tallyingProcedure := keeper.GetTallyingProcedure(ctx, proposal)

	tallyResults = TallyResult{
		Yes:        results[OptionYes],
		Abstain:    results[OptionAbstain],
		No:         results[OptionNo],
		NoWithVeto: results[OptionNoWithVeto],
	}

	// If no one votes, proposal fails
	if totalVotingPower.Sub(results[OptionAbstain]).Equal(sdk.ZeroDec()) {
		return REJECT, tallyResults, votingVals
	}

	//if more than 1/3 of voters abstain, proposal fails
	if tallyingProcedure.Participation.GT(totalVotingPower.Quo(systemVotingPower)) {
		return REJECT, tallyResults, votingVals
	}

	// If more than 1/3 of voters veto, proposal fails
	if results[OptionNoWithVeto].Quo(totalVotingPower).GT(tallyingProcedure.Veto) {
		return REJECTVETO, tallyResults, votingVals
	}

	// If more than 1/2 of non-abstaining voters vote Yes, proposal passes
	if results[OptionYes].Quo(totalVotingPower).GT(tallyingProcedure.Threshold) {
		return PASS, tallyResults, votingVals
	}
	// If more than 1/2 of non-abstaining voters vote No, proposal fails

	return REJECT, tallyResults, votingVals
}
