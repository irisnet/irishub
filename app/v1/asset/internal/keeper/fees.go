//nolint
package keeper

import (
	"github.com/irisnet/irishub/app/v1/auth"
	sdk "github.com/irisnet/irishub/types"
	"math"
	"strconv"
)

// fee factor formula: (ln(len({name}))/ln{base})^{exp}
const (
	FeeFactorBase = 3
	FeeFactorExp  = 4
)

// GatewayCreateFeeHandler performs fee handling for creating a gateway
func GatewayCreateFeeHandler(ctx sdk.Context, k Keeper, owner sdk.AccAddress, moniker string) sdk.Error {
	// get the required creation fee
	fee := GetGatewayCreateFee(ctx, k, moniker)

	return feeHandler(ctx, k, owner, fee)
}

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

// GatewayTokenIssueFeeHandler performs fee handling for issuing gateway token
func GatewayTokenIssueFeeHandler(ctx sdk.Context, k Keeper, owner sdk.AccAddress, symbol string) sdk.Error {
	// get the required issuance fee
	fee := GetGatewayTokenIssueFee(ctx, k, symbol)

	return feeHandler(ctx, k, owner, fee)
}

// GatewayTokenMintFeeHandler performs fee handling for minting gateway token
func GatewayTokenMintFeeHandler(ctx sdk.Context, k Keeper, owner sdk.AccAddress, symbol string) sdk.Error {
	// get the required minting fee
	fee := GetGatewayTokenMintFee(ctx, k, symbol)

	return feeHandler(ctx, k, owner, fee)
}

// feeHandler handles the fee of gateway or asset
func feeHandler(ctx sdk.Context, k Keeper, feeAcc sdk.AccAddress, fee sdk.Coin) sdk.Error {
	params := k.GetParamSet(ctx)
	assetTaxRate := params.AssetTaxRate

	// compute community tax and burned coin
	communityTaxCoin := sdk.NewCoin(fee.Denom, sdk.NewDecFromInt(fee.Amount).Mul(assetTaxRate).TruncateInt())
	burnedCoin := fee.Minus(communityTaxCoin)

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

// GetGatewayCreateFee returns the gateway creation fee
func GetGatewayCreateFee(ctx sdk.Context, k Keeper, moniker string) sdk.Coin {
	// get params
	params := k.GetParamSet(ctx)
	gatewayBaseFee := params.CreateGatewayBaseFee

	// compute the fee
	fee := calcFeeByBase(moniker, gatewayBaseFee.Amount)

	return sdk.NewCoin(sdk.IrisAtto, convertFeeToInt(fee))
}

// getTokenIssueFee returns the token issurance fee
func GetTokenIssueFee(ctx sdk.Context, k Keeper, symbol string) sdk.Coin {
	// get params
	params := k.GetParamSet(ctx)
	issueTokenBaseFee := params.IssueTokenBaseFee

	// compute the fee
	fee := calcFeeByBase(symbol, issueTokenBaseFee.Amount)

	return sdk.NewCoin(sdk.IrisAtto, convertFeeToInt(fee))
}

// getTokenMintFee returns the token mint fee
func GetTokenMintFee(ctx sdk.Context, k Keeper, symbol string) sdk.Coin {
	// get params
	params := k.GetParamSet(ctx)
	mintTokenFeeRatio := params.MintTokenFeeRatio

	// compute the issurance fee and mint fee
	issueFee := GetTokenIssueFee(ctx, k, symbol)
	mintFee := sdk.NewDecFromInt(issueFee.Amount).Mul(mintTokenFeeRatio)

	return sdk.NewCoin(sdk.IrisAtto, convertFeeToInt(mintFee))
}

// getGatewayTokenIssueFee returns the gateway token issurance fee
func GetGatewayTokenIssueFee(ctx sdk.Context, k Keeper, symbol string) sdk.Coin {
	// get params
	params := k.GetParamSet(ctx)
	gatewayAssetFeeRatio := params.GatewayAssetFeeRatio

	// compute the native token issurance fee and gateway token issurance fee
	nativeTokenIssueFee := GetTokenIssueFee(ctx, k, symbol)
	gatewayTokenIssueFee := sdk.NewDecFromInt(nativeTokenIssueFee.Amount).Mul(gatewayAssetFeeRatio)

	return sdk.NewCoin(sdk.IrisAtto, convertFeeToInt(gatewayTokenIssueFee))
}

// getGatewayTokenMintFee returns the gateway token mint fee
func GetGatewayTokenMintFee(ctx sdk.Context, k Keeper, symbol string) sdk.Coin {
	// get params
	params := k.GetParamSet(ctx)
	gatewayAssetFeeRatio := params.GatewayAssetFeeRatio

	// compute the native token mint fee and gateway token mint fee
	nativeTokenMintFee := GetTokenMintFee(ctx, k, symbol)
	gatewayTokenMintFee := sdk.NewDecFromInt(nativeTokenMintFee.Amount).Mul(gatewayAssetFeeRatio)

	return sdk.NewCoin(sdk.IrisAtto, convertFeeToInt(gatewayTokenMintFee))
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
