package upgrade

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/upgrade/params"
)

var Threshold = sdk.NewDecWithPrec(95, 2)

func tally(ctx sdk.Context, k Keeper) (passes bool) {

	proposalID := upgradeparams.GetCurrentUpgradeProposalId(ctx)

	if proposalID != 0 {

		totalVotingPower := sdk.ZeroDec()
		switchVotingPower := sdk.ZeroDec()
		for _, validator := range k.sk.GetAllValidators(ctx) {
			totalVotingPower = totalVotingPower.Add(validator.GetPower())

			valAcc := sdk.AccAddress(validator.OperatorAddr)
			if _, ok := k.GetSwitch(ctx, proposalID, valAcc); ok {
				switchVotingPower = switchVotingPower.Add(validator.GetPower())
			}
		}
		// If more than 95% of validator update , do switch
		if switchVotingPower.Quo(totalVotingPower).GT(Threshold) {
			return true
		}
	}
	return false
}
