package gov

import (
	sdk "github.com/irisnet/irishub/types"
	govtypes "github.com/irisnet/irishub/types/gov"
)

type ProposalResult string

const (
	PASS       ProposalResult = "pass"
	REJECT     ProposalResult = "reject"
	REJECTVETO ProposalResult = "reject-veto"
)

// validatorGovInfo used for tallying
type validatorGovInfo struct {
	Address         sdk.ValAddress      // address of the validator operator
	Power           sdk.Dec             // Power of a Validator
	Vote            govtypes.VoteOption // Vote of the validator
}

func tally(ctx sdk.Context, keeper Keeper, proposal govtypes.Proposal) (result ProposalResult, tallyResults govtypes.TallyResult) {
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
		valAddrStr := sdk.ValAddress(vote.Voter).String()
		if val, ok := currValidators[valAddrStr]; ok {
			val.Vote = vote.Option
			currValidators[valAddrStr] = val
			results[val.Vote] = results[val.Vote].Add(val.Power)
			totalVotingPower = totalVotingPower.Add(val.Power)
		}
	}

	////////////////////  iris begin  ///////////////////////////
	tallyingProcedure := GetTallyingProcedure(ctx)
	////////////////////  iris end  /////////////////////////////

	tallyResults = govtypes.TallyResult{
		Yes:        results[govtypes.OptionYes],
		Abstain:    results[govtypes.OptionAbstain],
		No:         results[govtypes.OptionNo],
		NoWithVeto: results[govtypes.OptionNoWithVeto],
	}

	// If no one votes, proposal fails
	if totalVotingPower.Sub(results[govtypes.OptionAbstain]).Equal(sdk.ZeroDec()) {
		return REJECT, tallyResults
	}
	////////////////////  iris begin  ///////////////////////////
	//if more than 1/3 of voters abstain, proposal fails
	if tallyingProcedure.Participation.GT(totalVotingPower.Quo(systemVotingPower)) {
		return REJECT, tallyResults
	}
	////////////////////  iris end  ///////////////////////////

	// If more than 1/3 of voters veto, proposal fails
	if results[govtypes.OptionNoWithVeto].Quo(totalVotingPower).GT(tallyingProcedure.Veto) {
		return REJECTVETO, tallyResults
	}

	// If more than 1/2 of non-abstaining voters vote Yes, proposal passes
	if results[govtypes.OptionYes].Quo(totalVotingPower).GT(tallyingProcedure.Threshold) {
		return PASS, tallyResults
	}
	// If more than 1/2 of non-abstaining voters vote No, proposal fails

	return REJECT, tallyResults
}
