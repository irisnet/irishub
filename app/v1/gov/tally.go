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
	Address             sdk.ValAddress // address of the validator operator
	Vote                VoteOption     // Vote of the validator
	TokenPerShare       sdk.Dec
	DelegatorShares     sdk.Dec // Total outstanding delegator shares
	DelegatorDeductions sdk.Dec // Delegator deductions from validator's delegators voting independently
}

func tally(ctx sdk.Context, keeper Keeper, proposal Proposal) (result ProposalResult, tallyResults TallyResult, votingVals map[string]bool) {
	results := make(map[VoteOption]sdk.Dec)
	results[OptionYes] = sdk.ZeroDec()
	results[OptionAbstain] = sdk.ZeroDec()
	results[OptionNo] = sdk.ZeroDec()
	results[OptionNoWithVeto] = sdk.ZeroDec()

	//voted votingPower
	totalVotingPower := sdk.ZeroDec()
	//all votingPower
	systemVotingPower := sdk.ZeroDec()
	currValidators := make(map[string]validatorGovInfo)
	votingVals = make(map[string]bool)

	keeper.vs.IterateBondedValidatorsByPower(ctx, func(index int64, validator sdk.Validator) (stop bool) {
		currValidators[validator.GetOperator().String()] = validatorGovInfo{
			Address:             validator.GetOperator(),
			TokenPerShare:       validator.GetTokens().Quo(validator.GetDelegatorShares()),
			Vote:                OptionEmpty,
			DelegatorShares:     validator.GetDelegatorShares(),
			DelegatorDeductions: sdk.ZeroDec(),
		}
		systemVotingPower = systemVotingPower.Add(validator.GetTokens())
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
			votingVals[valAddrStr] = true
			currValidators[valAddrStr] = val
		}
		// if validator is also delegator
		keeper.ds.IterateDelegations(ctx, vote.Voter, func(index int64, delegation sdk.Delegation) (stop bool) {
			valAddr := delegation.GetValidatorAddr().String()
			if valAddr == valAddrStr {
				return false
			}
			//only tally the delegator voting power under the validator
			if val, ok := currValidators[valAddr]; ok {
				val.DelegatorDeductions = val.DelegatorDeductions.Add(delegation.GetShares())
				currValidators[valAddr] = val

				votingPower := delegation.GetShares().Mul(val.TokenPerShare)
				results[vote.Option] = results[vote.Option].Add(votingPower)
				totalVotingPower = totalVotingPower.Add(votingPower)
			}
			return false
		})
	}

	// iterate over the validators again to tally their voting power
	for _, val := range currValidators {
		if val.Vote == OptionEmpty {
			continue
		}

		sharesAfterDeductions := val.DelegatorShares.Sub(val.DelegatorDeductions)
		votingPower := sharesAfterDeductions.Mul(val.TokenPerShare)

		results[val.Vote] = results[val.Vote].Add(votingPower)
		totalVotingPower = totalVotingPower.Add(votingPower)
	}

	tallyingProcedure := keeper.GetTallyingProcedure(ctx, proposal.GetProposalLevel())

	tallyResults = TallyResult{
		Yes:               results[OptionYes].QuoInt(sdk.AttoScaleFactor),
		Abstain:           results[OptionAbstain].QuoInt(sdk.AttoScaleFactor),
		No:                results[OptionNo].QuoInt(sdk.AttoScaleFactor),
		NoWithVeto:        results[OptionNoWithVeto].QuoInt(sdk.AttoScaleFactor),
		SystemVotingPower: systemVotingPower.QuoInt(sdk.AttoScaleFactor),
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
