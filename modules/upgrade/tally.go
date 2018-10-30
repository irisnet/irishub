package upgrade

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/upgrade/params"
)

var Threshold = sdk.NewDecWithPrec(95, 2)

func tally(ctx sdk.Context, k Keeper) (passes bool) {

	proposalID := upgradeparams.GetCurrentUpgradeProposalId(ctx)

	if proposalID != -1 {

		totalVotingPower := sdk.ZeroDec()
		switchVotingPower := sdk.ZeroDec()
		for _, validator := range k.sk.GetAllValidators(ctx) {
			totalVotingPower = totalVotingPower.Add(validator.GetPower())
			acc, err := sdk.AccAddressFromBech32(validator.OperatorAddr.String())
			if err == nil {
				if _, ok := k.GetSwitch(ctx, proposalID, acc); ok {
					switchVotingPower = switchVotingPower.Add(validator.GetPower())
				}
			}
		}
		// If more than 95% of validator update , do switch
		if switchVotingPower.Quo(totalVotingPower).GT(Threshold) {
			return true
		}
	}
	return false
}
