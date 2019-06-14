//nolint
package asset

import (
	"fmt"

	bank "github.com/irisnet/irishub/app/v1/bank"
	sdk "github.com/irisnet/irishub/types"
)

var (
	FeeFactorSetBySize = []string{
		"1.00",  //(ln3/ln3)^4
		"2.53",  //(ln4/ln3)^4
		"4.60",  //(ln5/ln3)^4
		"7.07",  //(ln6/ln3)^4
		"9.84",  //(ln7/ln3)^4
		"12.83", //(ln8/ln3)^4
	}

	startingSize = 3
)

// GatewayFeeHandler performs fee handling for creating a gateway
func GatewayFeeHandler(ctx sdk.Context, k Keeper, owner sdk.AccAddress, moniker string, fee sdk.Coins) sdk.Error {
	// get params
	params := k.GetParamSet(ctx)
	gatewayBaseFee := params.CreateGatewayBaseFee
	assetTaxRate := params.AssetTaxRate

	// compute fee
	totalFee := sdk.NewDec(int64(gatewayBaseFee)).Quo(calcGatewayFeeFactor(moniker))

	// check if the provided fee is enough
	if fee.AmountOf(sdk.NativeTokenName).LT(totalFee.TruncateInt()) {
		return ErrInsufficientFee(k.Codespace(), fmt.Sprintf("insufficient gateway create fee; expected %s, got %s", totalFee.TruncateInt(), fee.AmountOf(sdk.NativeTokenMinDenom)))
	}

	// compute tax and burned coin
	communityTax := sdk.NewCoin(sdk.NativeTokenName, assetTaxRate.Mul(totalFee).TruncateInt())
	burnedCoin := sdk.NewCoin(sdk.NativeTokenName, sdk.NewDec(1).Sub(assetTaxRate).Mul(totalFee).TruncateInt())

	// send community tax
	if _, err := k.bk.SendCoins(ctx, owner, bank.CommunityTaxCoinsAccAddr, sdk.Coins{communityTax}); err != nil {
		return err
	}

	// burn burnedCoin
	if _, err := k.bk.BurnCoins(ctx, owner, sdk.Coins{burnedCoin}); err != nil {
		return err
	}

	// decrease loosen tokens
	k.bk.DecreaseLoosenToken(ctx, sdk.Coins{burnedCoin})

	return nil
}

// calcGatewayFeeFactor computes the fee factor of the given moniker
func calcGatewayFeeFactor(moniker string) sdk.Dec {
	len := len(moniker)

	if len < MinimumGatewayMonikerSize || len > MaximumGatewayMonikerSize {
		return sdk.ZeroDec()
	}

	// error ignored
	feeFactor, _ := sdk.NewDecFromStr(FeeFactorSetBySize[len-startingSize])
	return feeFactor
}
