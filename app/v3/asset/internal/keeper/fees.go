//nolint
package keeper

import (
	"math"
	"strconv"

	"github.com/irisnet/irishub/app/v1/auth"
	sdk "github.com/irisnet/irishub/types"
)

// fee factor formula: (ln(len({name}))/ln{base})^{exp}
const (
	FeeFactorBase = 3
	FeeFactorExp  = 4
)

// DeductIssueTokenFee performs fee handling for issuing token
func (k Keeper) DeductIssueTokenFee(ctx sdk.Context, owner sdk.AccAddress, symbol string) sdk.Error {
	// get the required issuance fee
	fee := k.getTokenIssueFee(ctx, symbol)
	return feeHandler(ctx, k, owner, fee)
}

// DeductMintTokenFee performs fee handling for minting token
func (k Keeper) DeductMintTokenFee(ctx sdk.Context, owner sdk.AccAddress, symbol string) sdk.Error {
	// get the required minting fee
	fee := k.getTokenMintFee(ctx, symbol)
	return feeHandler(ctx, k, owner, fee)
}

// GetTokenIssueFee returns the token issurance fee
func (k Keeper) getTokenIssueFee(ctx sdk.Context, symbol string) sdk.Coin {
	// get params
	params := k.GetParamSet(ctx)
	issueTokenBaseFee := params.IssueTokenBaseFee

	// compute the fee
	fee := calcFeeByBase(symbol, issueTokenBaseFee.Amount)

	return sdk.NewCoin(sdk.IrisAtto, convertFeeToInt(fee))
}

// getTokenMintFee returns the token minting fee
func (k Keeper) getTokenMintFee(ctx sdk.Context, symbol string) sdk.Coin {
	// get params
	params := k.GetParamSet(ctx)
	mintTokenFeeRatio := params.MintTokenFeeRatio

	// compute the issurance and minting fees
	issueFee := k.getTokenIssueFee(ctx, symbol)
	mintFee := sdk.NewDecFromInt(issueFee.Amount).Mul(mintTokenFeeRatio)

	return sdk.NewCoin(sdk.IrisAtto, convertFeeToInt(mintFee))
}

// feeHandler handles the fee
func feeHandler(ctx sdk.Context, k Keeper, feeAcc sdk.AccAddress, fee sdk.Coin) sdk.Error {
	params := k.GetParamSet(ctx)
	assetTaxRate := params.AssetTaxRate

	// compute community tax and burned coin
	communityTaxCoin := sdk.NewCoin(fee.Denom, sdk.NewDecFromInt(fee.Amount).Mul(assetTaxRate).TruncateInt())
	burnedCoin := fee.Sub(communityTaxCoin)

	// send community tax
	if _, err := k.bk.SendCoins(ctx, feeAcc, auth.CommunityTaxCoinsAccAddr, sdk.Coins{communityTaxCoin}); err != nil {
		return err
	}
	ctx.CoinFlowTags().AppendCoinFlowTag(ctx, feeAcc.String(), auth.CommunityTaxCoinsAccAddr.String(), communityTaxCoin.String(), sdk.CommunityTaxCollectFlow, "")

	// burn burnedCoin
	if _, err := k.bk.BurnCoins(ctx, feeAcc, sdk.Coins{burnedCoin}); err != nil {
		return err
	}
	ctx.CoinFlowTags().AppendCoinFlowTag(ctx, feeAcc.String(), auth.BurnedCoinsAccAddr.String(), burnedCoin.String(), sdk.BurnFlow, "")

	return nil
}

// calcFeeByBase computes the actual fee according to the given base fee
func calcFeeByBase(name string, baseFee sdk.Int) sdk.Dec {
	feeFactor := calcFeeFactor(name)
	actualFee := sdk.NewDecFromInt(baseFee).Quo(feeFactor)

	return actualFee
}

// calcFeeFactor computes the fee factor of the given name
// Note: make sure that the name size is examined before invoking the function
func calcFeeFactor(name string) sdk.Dec {
	nameLen := len(name)
	if nameLen == 0 {
		panic("the length of name must be greater than 0")
	}

	denominator := math.Log(FeeFactorBase)
	numerator := math.Log(float64(nameLen))

	feeFactor := math.Pow(numerator/denominator, FeeFactorExp)
	feeFactorDec, err := sdk.NewDecFromStr(strconv.FormatFloat(feeFactor, 'f', 2, 64))
	if err != nil {
		panic("invalid string")
	}

	return feeFactorDec
}

// convertFeeToInt converts the given fee to Int.
// if greater than 1, rounds it; returns 1 otherwise
func convertFeeToInt(fee sdk.Dec) sdk.Int {
	feeNativeToken := fee.Quo(sdk.NewDecFromInt(sdk.NewIntWithDecimal(1, 18)))

	if feeNativeToken.GT(sdk.NewDec(1)) {
		return feeNativeToken.TruncateInt().Mul(sdk.NewIntWithDecimal(1, 18))
	} else {
		return sdk.NewInt(1).Mul(sdk.NewIntWithDecimal(1, 18))
	}
}
