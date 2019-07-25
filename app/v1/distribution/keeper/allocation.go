package keeper

import (
	"fmt"
	"github.com/irisnet/irishub/app/v1/auth"
	"strconv"

	"github.com/irisnet/irishub/app/v1/distribution/types"
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
		k.bankKeeper.AddCoins(ctx, auth.CommunityTaxCoinsAccAddr, feesCollected)
		//		feePool.CommunityPool = feePool.CommunityPool.Plus(feesCollectedDec)
		//		k.SetFeePool(ctx, feePool)
		k.feeKeeper.ClearCollectedFees(ctx)
		ctx.CoinFlowTags().AppendCoinFlowTag(ctx, "", auth.CommunityTaxCoinsAccAddr.String(), feesCollected.String(), sdk.CommunityTaxCollectFlow, "")
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

	//	feePool.CommunityPool = feePool.CommunityPool.Plus(communityFunding)
	fundingCoins, change := communityFunding.TruncateDecimal()
	k.bankKeeper.AddCoins(ctx, auth.CommunityTaxCoinsAccAddr, fundingCoins)
	ctx.CoinFlowTags().AppendCoinFlowTag(ctx, "", auth.CommunityTaxCoinsAccAddr.String(), fundingCoins.String(), sdk.CommunityTaxCollectFlow, "")

	communityTaxCoins := k.bankKeeper.GetCoins(ctx, auth.CommunityTaxCoinsAccAddr)
	communityTaxDec := sdk.NewDecFromInt(communityTaxCoins.AmountOf(sdk.IrisAtto))
	communityTaxFloat, err := strconv.ParseFloat(communityTaxDec.QuoInt(sdk.AttoScaleFactor).String(), 64)
	//communityTaxAmount, err := strconv.ParseFloat(feePool.CommunityPool.AmountOf(sdk.IrisAtto).QuoInt(sdk.AttoScaleFactor).String(), 64)
	if err == nil {
		k.metrics.CommunityTax.Set(communityTaxFloat)
	}

	logger.Info("Allocate reward to community tax fund", "allocate_amount", fundingCoins.String(), "total_community_tax", communityTaxCoins.String())

	// set the global pool within the distribution module
	poolReceived := feesCollectedDec.Minus(proposerReward).Minus(communityFunding).Plus(change)
	feePool.ValPool = feePool.ValPool.Plus(poolReceived)
	k.SetFeePool(ctx, feePool)

	logger.Info("Allocate reward to global validator pool", "allocate_amount", poolReceived.ToString(), "total_global_validator_pool", feePool.ValPool.ToString())

	// clear the now distributed fees
	k.feeKeeper.ClearCollectedFees(ctx)
}

// Allocate fee tax from the community fee pool, burn or send to trustee account
func (k Keeper) AllocateFeeTax(ctx sdk.Context, destAddr sdk.AccAddress, percent sdk.Dec, burn bool) {
	logger := ctx.Logger()
	//feePool := k.GetFeePool(ctx)
	//communityPool := feePool.CommunityPool
	//allocateCoins, _ := communityPool.MulDec(percent).TruncateDecimal()

	//feePool.CommunityPool = communityPool.Minus(types.NewDecCoins(allocateCoins))
	taxCoins := k.bankKeeper.GetCoins(ctx, auth.CommunityTaxCoinsAccAddr)
	taxDecCoins := types.NewDecCoins(taxCoins)
	allocatedDecCoins := taxDecCoins.MulDec(percent)
	allocatedCoins, _ := allocatedDecCoins.TruncateDecimal()
	taxLeftDecCoins := taxDecCoins.Minus(allocatedDecCoins)

	taxLeftDec := taxLeftDecCoins.AmountOf(sdk.IrisAtto)
	taxLeftFloat, err := strconv.ParseFloat(taxLeftDec.QuoInt(sdk.AttoScaleFactor).String(), 64)
	//communityTaxAmount, err := strconv.ParseFloat(feePool.CommunityPool.AmountOf(sdk.IrisAtto).QuoInt(sdk.AttoScaleFactor).String(), 64)
	if err == nil {
		k.metrics.CommunityTax.Set(taxLeftFloat)
	}

	//k.SetFeePool(ctx, feePool)
	logger.Info("Spend community tax fund", "total_community_tax_fund", taxCoins.String(), "left_community_tax_fund", taxLeftDecCoins.String())
	if burn {
		logger.Info("Burn community tax", "burn_amount", allocatedCoins.String())
		_, err := k.bankKeeper.BurnCoins(ctx, auth.CommunityTaxCoinsAccAddr, allocatedCoins)
		if err != nil {
			panic(err)
		}
		if !allocatedCoins.IsZero() {
			ctx.CoinFlowTags().AppendCoinFlowTag(ctx, auth.CommunityTaxCoinsAccAddr.String(), "", allocatedCoins.String(), sdk.BurnFlow, "")
		}
	} else {
		logger.Info("Grant community tax to account", "grant_amount", allocatedCoins.String(), "grant_address", destAddr.String())
		if !allocatedCoins.IsZero() {
			ctx.CoinFlowTags().AppendCoinFlowTag(ctx, "", destAddr.String(), allocatedCoins.String(), sdk.CommunityTaxUseFlow, "")
		}
		_, err := k.bankKeeper.SendCoins(ctx, auth.CommunityTaxCoinsAccAddr, destAddr, allocatedCoins)
		if err != nil {
			panic(err)
		}
	}

}
