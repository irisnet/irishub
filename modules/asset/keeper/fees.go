//nolint
package keeper

import (
	"math"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// fee factor formula: (ln(len({name}))/ln{base})^{exp}
const (
	FeeFactorBase = 3
	FeeFactorExp  = 4
)

// TokenIssueFeeHandler performs fee handling for issuing token
func TokenIssueFeeHandler(ctx sdk.Context, k Keeper, owner sdk.AccAddress, symbol string) sdk.Error {
	// get the required issuance fee
	fee := GetTokenIssueFee(ctx, k, symbol)

	return feeHandler(ctx, k, owner, fee)
}

// TokenMintFeeHandler performs fee handling for minting token
func TokenMintFeeHandler(ctx sdk.Context, k Keeper, owner sdk.AccAddress, symbol string) sdk.Error {
	// get the required minting fee
	fee := GetTokenMintFee(ctx, k, symbol)

	return feeHandler(ctx, k, owner, fee)
}

// feeHandler handles the fee of asset
func feeHandler(ctx sdk.Context, k Keeper, feeAcc sdk.AccAddress, fee sdk.Coin) sdk.Error {
	//assetTaxRate := k.AssetTaxRate(ctx)
	//
	// compute community tax and burned coin
	//communityTaxCoin := sdk.NewCoin(fee.Denom, sdk.NewDecFromInt(fee.Amount).Mul(assetTaxRate).TruncateInt())
	//burnedCoin := fee.Sub(communityTaxCoin)

	// send community tax
	feePool := k.distributionKeeper.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(sdk.NewDecCoins(sdk.NewCoins(fee)))
	k.distributionKeeper.SetFeePool(ctx, feePool)
	//if _, err := k.distributionKeeper.SendCoins(ctx, feeAcc, auth.CommunityTaxCoinsAccAddr, sdk.Coins{communityTaxCoin}); err != nil {
	//	return err
	//}
	//ctx.CoinFlowTags().AppendCoinFlowTag(ctx, feeAcc.String(), auth.CommunityTaxCoinsAccAddr.String(), communityTaxCoin.String(), sdk.CommunityTaxCollectFlow, "")
	//
	//// burn burnedCoin
	//k.supplyKeeper.b
	//if _, err := k.bk.BurnCoins(ctx, feeAcc, sdk.Coins{burnedCoin}); err != nil {
	//	return err
	//}
	//ctx.CoinFlowTags().AppendCoinFlowTag(ctx, feeAcc.String(), auth.BurnedCoinsAccAddr.String(), burnedCoin.String(), sdk.BurnFlow, "")

	return nil
}

// getTokenIssueFee returns the token issurance fee
func GetTokenIssueFee(ctx sdk.Context, k Keeper, symbol string) sdk.Coin {
	// compute the fee
	fee := calcFeeByBase(symbol, k.IssueTokenBaseFee(ctx))

	return sdk.NewCoin(k.AssetFeeDenom(ctx), convertFeeToInt(fee))
}

// getTokenMintFee returns the token mint fee
func GetTokenMintFee(ctx sdk.Context, k Keeper, symbol string) sdk.Coin {
	// compute the issurance fee and mint fee
	issueFee := GetTokenIssueFee(ctx, k, symbol)
	mintFee := sdk.NewDecFromInt(issueFee.Amount).Mul(k.MintTokenFeeRatio(ctx))

	return sdk.NewCoin(k.AssetFeeDenom(ctx), convertFeeToInt(mintFee))
}

// calcFeeByBase computes the actual fee according to the given base fee
func calcFeeByBase(name string, baseFee sdk.Int) sdk.Dec {
	feeFactor := calcFeeFactor(name)
	actualFee := sdk.NewDecFromInt(baseFee).Quo(feeFactor)

	return actualFee
}

// calcFeeFactor computes the fee factor of the given name(common for gateway and asset)
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
