//nolint
package asset

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	bank "github.com/irisnet/irishub/app/v1/bank"
	sdk "github.com/irisnet/irishub/types"
)

// fee factor formula: (ln(len({name}))/ln{base})^{exp}
const (
	FeeFactorBase = 3
	FeeFactorExp  = 4
)

// GatewayFeeHandler performs fee handling for creating a gateway
func GatewayFeeHandler(ctx sdk.Context, k Keeper, owner sdk.AccAddress, moniker string, fee sdk.Coin) sdk.Error {
	// get the actual fee
	actualFee := getGatewayCreateFee(ctx, k, moniker)

	// convert to native token min denom
	actualFeeMin := convertFeeToNativeTokenMin(actualFee)

	// check if the provided fee is enough
	if fee.Amount.LT(actualFeeMin) {
		return ErrInsufficientFee(k.Codespace(), fmt.Sprintf("insufficient gateway create fee: expected %d, got %d", actualFeeMin, fee.Amount))
	}

	return feeHandler(ctx, k, owner, sdk.NewCoin(fee.Denom, actualFeeMin))
}

// TokenIssueFeeHandler performs fee handling for issuing token
func TokenIssueFeeHandler(ctx sdk.Context, k Keeper, owner sdk.AccAddress, symbol string, fee sdk.Coin) sdk.Error {
	// get the actual fee
	actualFee := getTokenIssueFee(ctx, k, symbol)

	// convert to native token min denom
	actualFeeMin := convertFeeToNativeTokenMin(actualFee)

	// check if the provided fee is enough
	if fee.Amount.LT(actualFeeMin) {
		return ErrInsufficientFee(k.Codespace(), fmt.Sprintf("insufficient token issurance fee: expected %d, got %d", actualFeeMin, fee.Amount))
	}

	return feeHandler(ctx, k, owner, sdk.NewCoin(fee.Denom, actualFeeMin))
}

// TokenMintFeeHandler performs fee handling for minting token
func TokenMintFeeHandler(ctx sdk.Context, k Keeper, owner sdk.AccAddress, symbol string, fee sdk.Coin) sdk.Error {
	// get the actual fee
	actualFee := getTokenMintFee(ctx, k, symbol)

	// convert to native token min denom
	actualFeeMin := convertFeeToNativeTokenMin(actualFee)

	// check if the provided fee is enough
	if fee.Amount.LT(actualFeeMin) {
		return ErrInsufficientFee(k.Codespace(), fmt.Sprintf("insufficient token mint fee: expected %d, got %d", actualFeeMin, fee.Amount))
	}

	return feeHandler(ctx, k, owner, sdk.NewCoin(fee.Denom, actualFeeMin))
}

// GatewayTokenIssueFeeHandler performs fee handling for issuing gateway token
func GatewayTokenIssueFeeHandler(ctx sdk.Context, k Keeper, owner sdk.AccAddress, symbol string, fee sdk.Coin) sdk.Error {
	// get the actual fee
	actualFee := getGatewayTokenIssueFee(ctx, k, symbol)

	// convert to native token min denom
	actualFeeMin := convertFeeToNativeTokenMin(actualFee)

	// check if the provided fee is enough
	if fee.Amount.LT(actualFeeMin) {
		return ErrInsufficientFee(k.Codespace(), fmt.Sprintf("insufficient gateway token issurance fee: expected %d, got %d", actualFeeMin, fee.Amount))
	}

	return feeHandler(ctx, k, owner, sdk.NewCoin(fee.Denom, actualFeeMin))
}

// GatewayTokenMintFeeHandler performs fee handling for minting gateway token
func GatewayTokenMintFeeHandler(ctx sdk.Context, k Keeper, owner sdk.AccAddress, symbol string, fee sdk.Coin) sdk.Error {
	// get the actual fee
	actualFee := getGatewayTokenMintFee(ctx, k, symbol)

	// convert to native token min denom
	actualFeeMin := convertFeeToNativeTokenMin(actualFee)

	// check if the provided fee is enough
	if fee.Amount.LT(actualFeeMin) {
		return ErrInsufficientFee(k.Codespace(), fmt.Sprintf("insufficient gateway token mint fee: expected %d, got %d", actualFeeMin, fee.Amount))
	}

	return feeHandler(ctx, k, owner, sdk.NewCoin(fee.Denom, actualFeeMin))
}

// feeHandler handles the fee of gateway or asset
func feeHandler(ctx sdk.Context, k Keeper, feeAcc sdk.AccAddress, fee sdk.Coin) sdk.Error {
	params := k.GetParamSet(ctx)
	assetTaxRate := params.AssetTaxRate

	// compute community tax and burned coin
	communityTaxCoin := sdk.NewCoin(fee.Denom, sdk.NewDecFromInt(fee.Amount).Mul(assetTaxRate).TruncateInt())
	burnedCoin := fee.Minus(communityTaxCoin)

	// send community tax
	if _, err := k.bk.SendCoins(ctx, feeAcc, bank.CommunityTaxCoinsAccAddr, sdk.Coins{communityTaxCoin}); err != nil {
		return err
	}

	// burn burnedCoin
	if _, err := k.bk.BurnCoins(ctx, feeAcc, sdk.Coins{burnedCoin}); err != nil {
		return err
	}

	return nil
}

// getGatewayCreateFee returns the gateway creation fee
func getGatewayCreateFee(ctx sdk.Context, k Keeper, moniker string) sdk.Int {
	// get params
	params := k.GetParamSet(ctx)
	gatewayBaseFee := params.CreateGatewayBaseFee

	// compute the fee
	fee := calcFeeByBase(moniker, gatewayBaseFee.Amount)

	// convert to native token
	return convertFeeToNativeToken(fee)
}

// getTokenIssueFee returns the token issurance fee
func getTokenIssueFee(ctx sdk.Context, k Keeper, symbol string) sdk.Int {
	// get params
	params := k.GetParamSet(ctx)
	issueFTBaseFee := params.IssueFTBaseFee

	// compute the fee
	fee := calcFeeByBase(symbol, issueFTBaseFee.Amount)

	// convert to native token
	return convertFeeToNativeToken(fee)
}

// getTokenMintFee returns the token mint fee
func getTokenMintFee(ctx sdk.Context, k Keeper, symbol string) sdk.Int {
	// get params
	params := k.GetParamSet(ctx)
	mintFTFeeRate := params.MintFTFeeRatio

	// compute the issurance fee and mint fee
	issueFee := getTokenIssueFee(ctx, k, symbol)
	mintFee := sdk.NewDecFromInt(issueFee).Mul(mintFTFeeRate)

	// error ignored
	mintFeeFloat64, _ := strconv.ParseFloat(mintFee.String(), 64)

	// round fee
	return convertFeeToInt(mintFeeFloat64)
}

// getGatewayTokenIssueFee returns the gateway token issurance fee
func getGatewayTokenIssueFee(ctx sdk.Context, k Keeper, symbol string) sdk.Int {
	// get params
	params := k.GetParamSet(ctx)
	gatewayAssetFeeRatio := params.GatewayAssetFeeRatio

	// compute the native token issurance fee and gateway token issurance fee
	nativeTokenIssueFee := getTokenIssueFee(ctx, k, symbol)
	gatewayTokenIssueFee := sdk.NewDecFromInt(nativeTokenIssueFee).Mul(gatewayAssetFeeRatio)

	// error ignored
	gwTokenIssueFeeF64, _ := strconv.ParseFloat(gatewayTokenIssueFee.String(), 64)

	// round fee
	return convertFeeToInt(gwTokenIssueFeeF64)
}

// getGatewayTokenMintFee returns the gateway token mint fee
func getGatewayTokenMintFee(ctx sdk.Context, k Keeper, symbol string) sdk.Int {
	// get params
	params := k.GetParamSet(ctx)
	gatewayAssetFeeRatio := params.GatewayAssetFeeRatio

	// compute the native token mint fee and gateway token mint fee
	nativeTokenMintFee := getTokenMintFee(ctx, k, symbol)
	gatewayTokenMintFee := sdk.NewDecFromInt(nativeTokenMintFee).Mul(gatewayAssetFeeRatio)

	// error ignored
	gwTokenMintFeeF64, _ := strconv.ParseFloat(gatewayTokenMintFee.String(), 64)

	// round fee
	return convertFeeToInt(gwTokenMintFeeF64)
}

// calcFeeByBase computes the actual fee according to the given base fee
func calcFeeByBase(name string, baseFee sdk.Int) float64 {
	feeFactor := calcFeeFactor(name)
	baseFeeFloat64, _ := strconv.ParseFloat(sdk.NewDecFromInt(baseFee).String(), 64)

	actualFee := baseFeeFloat64 / feeFactor

	return actualFee
}

// calcFeeFactor computes the fee factor of the given name(common for gateway and asset)
// Note: make sure that the name size is examined before invoking the function
func calcFeeFactor(name string) float64 {
	nameLen := len(name)
	if nameLen == 0 {
		panic("the length of name must be greater than 0")
	}

	denominator := math.Log(FeeFactorBase)
	numerator := math.Log(float64(nameLen))

	// error ignored
	feeFactor := math.Pow(numerator/denominator, FeeFactorExp)
	return feeFactor
}

// convertFeeToNativeToken converts fee to native token
func convertFeeToNativeToken(fee float64) sdk.Int {
	nativeTokenAmount := fee / math.Pow10(18)
	return convertFeeToInt(nativeTokenAmount)
}

// convertFeeToNativeTokenMin converts fee to native token min denom
func convertFeeToNativeTokenMin(fee sdk.Int) sdk.Int {
	return sdk.NewIntWithDecimal(fee.Int64(), 18)
}

// convertFeeToInt converts the given fee to Int
// if less than 1, rounds to 1; returns 1 otherwise
func convertFeeToInt(fee float64) sdk.Int {
	var feeInt64 int64

	if fee > 1 {
		feeInt64 = int64(math.Round(fee))
	} else {
		feeInt64 = 1
	}

	return sdk.NewInt(feeInt64)
}

// GatewayFeeOutput is for the gateway fee query output
type GatewayFeeOutput struct {
	Exist bool     `json:"exist"` // indicate if the gateway has existed
	Fee   sdk.Coin `json:"fee"`   // creation fee
}

// String implements stringer
func (gfo GatewayFeeOutput) String() string {
	var out strings.Builder
	if gfo.Exist {
		out.WriteString("The gateway moniker has existed\n")
	}

	out.WriteString(fmt.Sprintf("Fee: %s", gfo.Fee.String()))

	return out.String()
}

// TokenFeesOutput is for the token fees query output
type TokenFeesOutput struct {
	Exist    bool     `exist`            // indicate if the token has existed
	IssueFee sdk.Coin `json:"issue_fee"` // issue fee
	MintFee  sdk.Coin `json:"mint_fee"`  // mint fee
}

// String implements stringer
func (tfo TokenFeesOutput) String() string {
	var out strings.Builder
	if tfo.Exist {
		out.WriteString("The token id has existed\n")
	}

	out.WriteString(fmt.Sprintf(`Fees:
  IssueFee: %s
  MintFee:  %s`,
		tfo.IssueFee.String(), tfo.MintFee.String()))

	return out.String()
}
