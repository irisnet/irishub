//nolint
package keeper

import (
	"fmt"
	"math"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/asset/01-token/internal/types"
)

// fee factor formula: (ln(len({name}))/ln{base})^{exp}
const (
	FeeFactorBase = 3
	FeeFactorExp  = 4
)

// IssueTokenFeeHandler performs fee handling for issuing token
func IssueTokenFeeHandler(ctx sdk.Context, k Keeper, owner sdk.AccAddress, symbol string) error {
	// get the required issuance fee
	fee := GetTokenIssueFee(ctx, k, symbol)
	return feeHandler(ctx, k, owner, fee)
}

// MintTokenFeeHandler performs fee handling for minting token
func MintTokenFeeHandler(ctx sdk.Context, k Keeper, owner sdk.AccAddress, symbol string) error {
	// get the required minting fee
	fee := GetTokenMintFee(ctx, k, symbol)
	return feeHandler(ctx, k, owner, fee)
}

// feeHandler handles the fee of asset
func feeHandler(ctx sdk.Context, k Keeper, feeAcc sdk.AccAddress, fee sdk.Coin) error {
	assetTaxRate := k.AssetTaxRate(ctx)

	// compute community tax and burned coin
	communityTaxCoin := sdk.NewCoin(fee.Denom, sdk.NewDecFromInt(fee.Amount).Mul(assetTaxRate).TruncateInt())
	burnedCoins := sdk.NewCoins(fee.Sub(communityTaxCoin))

	// send all fees to module account
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(
		ctx, feeAcc, types.SubModuleName, sdk.NewCoins(fee),
	); err != nil {
		return err
	}

	// send community tax to collectedFees
	if err := k.addCollectedFees(ctx, sdk.NewCoins(communityTaxCoin)); err != nil {
		return err
	}

	// burn burnedCoin
	return k.supplyKeeper.BurnCoins(ctx, types.SubModuleName, burnedCoins)
}

// getTokenIssueFee returns the token issurance fee
func GetTokenIssueFee(ctx sdk.Context, k Keeper, symbol string) sdk.Coin {
	// compute the fee
	issueTokenBaseFee := k.IssueTokenBaseFee(ctx)
	token, found := k.GetTokenByMintUint(ctx, issueTokenBaseFee.Denom)
	if !found {
		panic(fmt.Sprintf("token [%s] not found", issueTokenBaseFee.Denom))
	}
	fee := calcFeeByBase(symbol, issueTokenBaseFee.Amount)
	return sdk.NewCoin(issueTokenBaseFee.Denom, convertFeeToInt(fee, token.Scale))
}

// getTokenMintFee returns the token mint fee
func GetTokenMintFee(ctx sdk.Context, k Keeper, symbol string) sdk.Coin {
	// compute the issurance fee and mint fee
	issueFee := GetTokenIssueFee(ctx, k, symbol)
	token, found := k.GetTokenByMintUint(ctx, issueFee.Denom)
	if !found {
		panic(fmt.Sprintf("token [%s] not found", issueFee.Denom))
	}
	mintFee := sdk.NewDecFromInt(issueFee.Amount).Mul(k.MintTokenFeeRatio(ctx))

	return sdk.NewCoin(issueFee.Denom, convertFeeToInt(mintFee, token.Scale))
}

// calcFeeByBase computes the actual fee according to the given base fee
func calcFeeByBase(name string, baseFee sdk.Int) sdk.Dec {
	feeFactor := calcFeeFactor(name)
	return sdk.NewDecFromInt(baseFee).Quo(feeFactor)
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
func convertFeeToInt(fee sdk.Dec, scale uint8) sdk.Int {
	scaleFactor := sdk.NewIntWithDecimal(1, int(scale))
	feeNativeToken := fee.Quo(scaleFactor.ToDec()).TruncateInt()
	if feeNativeToken.GT(sdk.NewInt(1)) {
		return feeNativeToken.Mul(scaleFactor)
	} else {
		return scaleFactor
	}
}
