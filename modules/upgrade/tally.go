package upgrade

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var Threshold = sdk.NewRat(95, 100)

func tally(ctx sdk.Context, k Keeper) (passes bool) {

	proposalID := k.GetCurrentProposalID(ctx)

	if proposalID!=-1{

		totalVotingPower := sdk.ZeroRat()
		switchVotingPower:= sdk.ZeroRat()

	    for _,validator :=range k.sk.GetAllValidators(ctx) {
			totalVotingPower.Add(validator.GetPower())
	   	    if _,ok := k.GetSwitch(ctx,proposalID,validator.Owner);ok {
				switchVotingPower.Add(validator.GetPower())
		    }
	    }

		// If more than 95% of validator update , do switch
		if switchVotingPower.Quo(totalVotingPower).GT(Threshold) {
			k.SetCurrentProposalID(ctx,-1)
			return true
		}

	}

	k.SetCurrentProposalID(ctx,-1)
	return  false
}
