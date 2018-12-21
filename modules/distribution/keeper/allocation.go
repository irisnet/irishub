package keeper

import (
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/distribution/types"
)

// Allocate fees handles distribution of the collected fees
func (k Keeper) AllocateTokens(ctx sdk.Context, percentVotes sdk.Dec, proposer sdk.ConsAddress) {

	// get the proposer of this block
	proposerValidator := k.stakeKeeper.ValidatorByConsAddr(ctx, proposer)

	proposerDist := k.GetValidatorDistInfo(ctx, proposerValidator.GetOperator())

	// get the fees which have been getting collected through all the
	// transactions in the block
	feesCollected := k.feeCollectionKeeper.GetCollectedFees(ctx)
	feesCollectedDec := types.NewDecCoins(feesCollected)

	feePool := k.GetFeePool(ctx)
	if k.stakeKeeper.GetLastTotalPower(ctx).IsZero() {
		feePool.CommunityPool = feePool.CommunityPool.Plus(feesCollectedDec)
		k.SetFeePool(ctx, feePool)
		k.feeCollectionKeeper.ClearCollectedFees(ctx)
		return
	}

	// allocated rewards to proposer
	baseProposerReward := k.GetBaseProposerReward(ctx)
	bonusProposerReward := k.GetBonusProposerReward(ctx)
	proposerMultiplier := baseProposerReward.Add(bonusProposerReward.Mul(percentVotes))
	proposerReward := feesCollectedDec.MulDec(proposerMultiplier)

	// apply commission
	commission := proposerReward.MulDec(proposerValidator.GetCommission())
	remaining := proposerReward.Minus(commission)
	proposerDist.ValCommission = proposerDist.ValCommission.Plus(commission)
	proposerDist.DelPool = proposerDist.DelPool.Plus(remaining)

	// allocate community funding
	communityTax := k.GetCommunityTax(ctx)
	communityFunding := feesCollectedDec.MulDec(communityTax)
	feePool.CommunityPool = feePool.CommunityPool.Plus(communityFunding)

	// set the global pool within the distribution module
	poolReceived := feesCollectedDec.Minus(proposerReward).Minus(communityFunding)
	feePool.ValPool = feePool.ValPool.Plus(poolReceived)

	k.SetValidatorDistInfo(ctx, proposerDist)
	k.SetFeePool(ctx, feePool)

	// clear the now distributed fees
	k.feeCollectionKeeper.ClearCollectedFees(ctx)
}

// Allocate fee tax from the community fee pool, burn or send to trustee account
func (k Keeper) AllocateFeeTax(ctx sdk.Context, destAddr sdk.AccAddress, percent sdk.Dec, burn bool) {
	feePool := k.GetFeePool(ctx)
	communityPool := feePool.CommunityPool
	allocateCoins, _ := communityPool.MulDec(percent).TruncateDecimal()
	feePool.CommunityPool = communityPool.Minus(types.NewDecCoins(allocateCoins))
	k.SetFeePool(ctx, feePool)

	if burn {
		stakeDenom := k.stakeKeeper.GetStakeDenom(ctx)
		for _, coin := range allocateCoins {
			if coin.Denom == stakeDenom {
				k.stakeKeeper.BurnAmount(ctx, sdk.NewDecFromInt(coin.Amount))
			}
		}
	} else {
		k.bankKeeper.AddCoins(ctx, destAddr, allocateCoins)
	}

}
