package keeper

import (
	"fmt"
	"strconv"

	"github.com/irisnet/irishub/modules/distribution/types"
	sdk "github.com/irisnet/irishub/types"
)

// Allocate fees handles distribution of the collected fees
func (k Keeper) AllocateTokens(ctx sdk.Context, percentVotes sdk.Dec, proposer sdk.ConsAddress) {
	logger := ctx.Logger()
	// get the proposer of this block
	proposerValidator := k.stakeKeeper.ValidatorByConsAddr(ctx, proposer)

	if proposerValidator == nil {
		panic(fmt.Sprintf("Can't find proposer %s in validator set", proposerValidator.GetConsAddr()))
	}

	proposerDist := k.GetValidatorDistInfo(ctx, proposerValidator.GetOperator())

	// get the fees which have been getting collected through all the
	// transactions in the block
	feesCollected := k.feeKeeper.GetCollectedFees(ctx)
	feesCollectedDec := types.NewDecCoins(feesCollected)

	logger.Info("Get collected transaction fee token and minted token", "collected_token", feesCollected)

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
		logger.Info("Allocate reward to proposer", "proposer_address", proposerValidator.GetOperator().String())
		baseProposerReward := k.GetBaseProposerReward(ctx)
		bonusProposerReward := k.GetBonusProposerReward(ctx)
		proposerMultiplier := baseProposerReward.Add(bonusProposerReward.Mul(percentVotes))
		proposerReward = feesCollectedDec.MulDec(proposerMultiplier)

		// apply commission
		commission := proposerReward.MulDec(proposerValidator.GetCommission())
		remaining := proposerReward.Minus(commission)
		proposerDist.ValCommission = proposerDist.ValCommission.Plus(commission)
		proposerDist.DelPool = proposerDist.DelPool.Plus(remaining)
		logger.Info("Allocate commission to proposer commission pool", "commission", commission.ToString())
		logger.Info("Allocate reward to proposer delegation reward pool", "delegation_reward", remaining.ToString())

		// save validator distribution info
		k.SetValidatorDistInfo(ctx, proposerDist)
	} else {
		logger.Info("The block proposer is jailed, distribute no reward to it", "proposer_address", proposerValidator.GetOperator().String())
	}

	// allocate community funding
	communityTax := k.GetCommunityTax(ctx)
	communityFunding := feesCollectedDec.MulDec(communityTax)
	feePool.CommunityPool = feePool.CommunityPool.Plus(communityFunding)

	communityTaxAmount, err := strconv.ParseFloat(feePool.CommunityPool.AmountOf(sdk.IrisAtto).QuoInt(sdk.AttoScaleFactor).String(), 64)
	if err == nil {
		k.metrics.CommunityTax.Set(communityTaxAmount)
	}

	logger.Info("Allocate reward to community tax fund", "allocate_amount", communityFunding.ToString(), "total_community_tax", feePool.CommunityPool.ToString())

	// set the global pool within the distribution module
	poolReceived := feesCollectedDec.Minus(proposerReward).Minus(communityFunding)
	feePool.ValPool = feePool.ValPool.Plus(poolReceived)
	k.SetFeePool(ctx, feePool)

	logger.Info("Allocate reward to global validator pool", "allocate_amount", poolReceived.ToString(), "total_global_validator_pool", feePool.ValPool.ToString())

	// clear the now distributed fees
	k.feeKeeper.ClearCollectedFees(ctx)
}

// Allocate fee tax from the community fee pool, burn or send to trustee account
func (k Keeper) AllocateFeeTax(ctx sdk.Context, destAddr sdk.AccAddress, percent sdk.Dec, burn bool) {
	logger := ctx.Logger()
	feePool := k.GetFeePool(ctx)
	communityPool := feePool.CommunityPool
	allocateCoins, _ := communityPool.MulDec(percent).TruncateDecimal()
	feePool.CommunityPool = communityPool.Minus(types.NewDecCoins(allocateCoins))

	communityTaxAmount, err := strconv.ParseFloat(feePool.CommunityPool.AmountOf(sdk.IrisAtto).QuoInt(sdk.AttoScaleFactor).String(), 64)
	if err == nil {
		k.metrics.CommunityTax.Set(communityTaxAmount)
	}

	k.SetFeePool(ctx, feePool)
	logger.Info("Spend community tax fund", "total_community_tax_fund", communityPool.ToString(), "left_community_tax_fund", feePool.CommunityPool.ToString())
	if burn {
		logger.Info("Burn community tax", "burn_amount", allocateCoins.String())
		_, err := k.bankKeeper.BurnCoinsFromPool(ctx, "communityTax", allocateCoins)
		if err != nil {
			panic(err)
		}
	} else {
		logger.Info("Grant community tax to account", "grant_amount", allocateCoins.String(), "grant_address", destAddr.String())
		if !allocateCoins.IsZero() {
			ctx.CoinFlowTags().AppendCoinFlowTag(ctx, "", destAddr.String(), allocateCoins.String(), sdk.CommunityTaxCollectFlow, "")
		}
		_, _, err := k.bankKeeper.AddCoins(ctx, destAddr, allocateCoins)
		if err != nil {
			panic(err)
		}
	}

}
