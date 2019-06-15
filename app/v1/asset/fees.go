//nolint
package asset

import (
	"fmt"
	"math"
	"math/big"

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
	// get params
	params := k.GetParamSet(ctx)
	gatewayBaseFee := params.CreateGatewayBaseFee

	// check if the denom of fee is same as that of gatewayBaseFee
	if fee.Denom != gatewayBaseFee.Denom {
		return ErrIncorrectFeeDenom(k.Codespace(), fmt.Sprintf("incorrect fee denom: expected %s, got %s", gatewayBaseFee.Denom, fee.Denom))
	}

	// compute the actual fee
	actualFee := calcFee(moniker, gatewayBaseFee.Amount)

	// check if the provided fee is enough
	if fee.Amount.LT(actualFee) {
		return ErrInsufficientFee(k.Codespace(), fmt.Sprintf("insufficient gateway create fee: expected %d, got %d", actualFee, fee.Amount))
	}

	return feeHandler(ctx, k, owner, sdk.NewCoin(fee.Denom, actualFee))
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

// calcFee computes the actual fee according to the given base fee
func calcFee(name string, baseFee sdk.Int) sdk.Int {
	feeFactor, _ := big.NewFloat(calcFeeFactor(name)).Int64()
	actualFee := baseFee.Div(sdk.NewInt(feeFactor))

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
