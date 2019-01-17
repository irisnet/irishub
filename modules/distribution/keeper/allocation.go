package keeper

import (
	"fmt"

	"github.com/irisnet/irishub/modules/distribution/types"
	sdk "github.com/irisnet/irishub/types"
)

// Allocate fees handles distribution of the collected fees
func (k Keeper) AllocateTokens(ctx sdk.Context, percentVotes sdk.Dec, proposer sdk.ConsAddress) {
	logger := ctx.Logger().With("module", "x/distr")
	height := ctx.BlockHeight()
	// get the proposer of this block
	proposerValidator := k.stakeKeeper.ValidatorByConsAddr(ctx, proposer)

	proposerDist := k.GetValidatorDistInfo(ctx, proposerValidator.GetOperator())

	// get the fees which have been getting collected through all the
	// transactions in the block
	feesCollected := k.feeKeeper.GetCollectedFees(ctx)
	feesCollectedDec := types.NewDecCoins(feesCollected)

	logger.Info(fmt.Sprintf("Collected fees %s at height %d", feesCollected, height))

	feePool := k.GetFeePool(ctx)
	if k.stakeKeeper.GetLastTotalPower(ctx).IsZero() {
		feePool.CommunityPool = feePool.CommunityPool.Plus(feesCollectedDec)
		k.SetFeePool(ctx, feePool)
		k.feeKeeper.ClearCollectedFees(ctx)
		return
	}

	var proposerReward types.DecCoins
	// If a validator is jailed, distribute no reward to it
	// The jailed validator happen to be a proposer which is a very corner case
	validator := k.stakeKeeper.Validator(ctx, proposerValidator.GetOperator())
	if !validator.GetJailed() {
		// allocated rewards to proposer
		baseProposerReward := k.GetBaseProposerReward(ctx)
		bonusProposerReward := k.GetBonusProposerReward(ctx)
		proposerMultiplier := baseProposerReward.Add(bonusProposerReward.Mul(percentVotes))
		proposerReward = feesCollectedDec.MulDec(proposerMultiplier)

		// apply commission
		commission := proposerReward.MulDec(proposerValidator.GetCommission())
		remaining := proposerReward.Minus(commission)
		proposerDist.ValCommission = proposerDist.ValCommission.Plus(commission)
		proposerDist.DelPool = proposerDist.DelPool.Plus(remaining)

		// save validator distribution info
		k.SetValidatorDistInfo(ctx, proposerDist)
	}

	// allocate community funding
	communityTax := k.GetCommunityTax(ctx)
	communityFunding := feesCollectedDec.MulDec(communityTax)
	feePool.CommunityPool = feePool.CommunityPool.Plus(communityFunding)

	logger.Info(fmt.Sprintf("Collected community tax funding %s at height %d", communityFunding.ToString(), height))
	logger.Info(fmt.Sprintf("Community pool increase to %s at height %d", feePool.CommunityPool.ToString(), height))

	// set the global pool within the distribution module
	poolReceived := feesCollectedDec.Minus(proposerReward).Minus(communityFunding)
	feePool.ValPool = feePool.ValPool.Plus(poolReceived)
	k.SetFeePool(ctx, feePool)

	logger.Info(fmt.Sprintf("Validators pool increase to %s at height %d", feePool.ValPool.ToString(), height))

	// clear the now distributed fees
	k.feeKeeper.ClearCollectedFees(ctx)
}

// Allocate fee tax from the community fee pool, burn or send to trustee account
func (k Keeper) AllocateFeeTax(ctx sdk.Context, destAddr sdk.AccAddress, percent sdk.Dec, burn bool) {
	feePool := k.GetFeePool(ctx)
	communityPool := feePool.CommunityPool
	allocateCoins, _ := communityPool.MulDec(percent).TruncateDecimal()
	feePool.CommunityPool = communityPool.Minus(types.NewDecCoins(allocateCoins))
	k.SetFeePool(ctx, feePool)

	if burn {
		k.bankKeeper.BurnCoinsFromPool(ctx, "communityTax", allocateCoins)
	} else {
		k.bankKeeper.AddCoins(ctx, destAddr, allocateCoins)
	}

}
