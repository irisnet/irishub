package upgrade

import (
	sdk "github.com/irisnet/irishub/types"
)

func tally(ctx sdk.Context,versionProtocol uint64, k Keeper) (passes bool) {

	totalVotingPower := sdk.ZeroDec()
	switchVotingPower := sdk.ZeroDec()

	k.sk.IterateBondedValidatorsByPower(ctx, func(index int64, validator sdk.Validator) (stop bool) {
		totalVotingPower = totalVotingPower.Add(validator.GetPower())
		valAcc := validator.GetConsAddr().String()
		if ok := k.GetSignal(ctx, versionProtocol, valAcc); ok {
			switchVotingPower = switchVotingPower.Add(validator.GetPower())
		}

		return false
	})

	// If more than 95% of validator update , do switch
	if switchVotingPower.Quo(totalVotingPower).GT(GetUpgradeThreshlod(ctx)) {
		return true
	}
	return false
}
