//nolint
package asset

import (
	"fmt"

	bank "github.com/irisnet/irishub/app/v1/bank"
	sdk "github.com/irisnet/irishub/types"
)

var (
	FeeFactorSetBySize = []string{
		"1.45",  //(ln3)^4
		"3.69",  //(ln4)^4
		"6.70",  //(ln5)^4
		"10.30", //(ln6)^4
		"14.33", //(ln7)^4
		"18.69", //(ln8)^4
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
	feeCoin := sdk.NewCoin(sdk.NativeTokenMinDenom, totalFee.TruncateInt())

	// check if the provided fee is enough
	if fee.AmountOf(sdk.NativeTokenMinDenom).LT(totalFee.TruncateInt()) {
		return ErrInsufficientFee(k.Codespace(), fmt.Sprintf("insufficient gateway create fee; expected %s, got %s", totalFee.TruncateInt(), fee.AmountOf(sdk.NativeTokenMinDenom)))
	}

	// compute tax and burned coin
	communityTax := sdk.NewCoin(sdk.NativeTokenMinDenom, assetTaxRate.Mul(totalFee).TruncateInt())
	burnedCoin := sdk.NewCoin(sdk.NativeTokenMinDenom, sdk.NewDec(1).Sub(assetTaxRate).Mul(totalFee).TruncateInt())

	// substract coin from owner
	if _, _, err := k.bk.SubtractCoins(ctx, owner, sdk.Coins{feeCoin}); err != nil {
		return err
	}

	// add community tax
	if _, _, err := k.bk.AddCoins(ctx, bank.CommunityTaxCoinsAccAddr, sdk.Coins{communityTax}); err != nil {
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
