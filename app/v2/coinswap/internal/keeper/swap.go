package keeper

import (
	"fmt"
	"strings"

	"github.com/irisnet/irishub/app/v2/coinswap/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

func (k Keeper) SwapCoins(ctx sdk.Context, sender sdk.AccAddress, coinSold, coinBought sdk.Coin) error {
	if !k.HasCoins(ctx, sender, coinSold) {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("sender account does not have sufficient amount of %s to fulfill the swap order", coinSold.Denom))
	}

	moduleName, err := k.GetModuleName(coinSold.Denom, coinBought.Denom)
	if err != nil {
		return err
	}

	k.SendCoins(ctx, sender, moduleName, coinSold)
	k.RecieveCoins(ctx, sender, moduleName, coinBought)
	return nil
}

// GetInputPrice returns the amount of coins bought (calculated) given the input amount being sold (exact)
// The fee is included in the input coins being bought
// https://github.com/runtimeverification/verified-smart-contracts/blob/uniswap/uniswap/x-y-k.pdf
// TODO: continue using numerator/denominator -> open issue for eventually changing to sdk.Dec
func (k Keeper) GetInputPrice(ctx sdk.Context, soldCoin sdk.Coin, boughtDenom string) sdk.Int {
	moduleName, err := k.GetModuleName(soldCoin.Denom, boughtDenom)
	if err != nil {
		panic(err)
	}
	reservePool, found := k.GetReservePool(ctx, moduleName)
	if !found {
		panic(fmt.Sprintf("reserve pool for %s not found", moduleName))
	}
	inputBalance := reservePool.AmountOf(soldCoin.Denom)
	outputBalance := reservePool.AmountOf(boughtDenom)
	fee := k.GetFeeParam(ctx)

	inputAmtWithFee := soldCoin.Amount.Mul(fee.Numerator)
	numerator := inputAmtWithFee.Mul(outputBalance)
	denominator := inputBalance.Mul(fee.Denominator).Add(inputAmtWithFee)
	return numerator.Div(denominator)
}

// GetOutputPrice returns the amount of coins sold (calculated) given the output amount being bought (exact)
// The fee is included in the output coins being bought
// https://github.com/runtimeverification/verified-smart-contracts/blob/uniswap/uniswap/x-y-k.pdf
// TODO: continue using numerator/denominator -> open issue for eventually changing to sdk.Dec
func (k Keeper) GetOutputPrice(ctx sdk.Context, boughtCoin sdk.Coin, soldDenom string) sdk.Int {
	moduleName, err := k.GetModuleName(boughtCoin.Denom, soldDenom)
	if err != nil {
		panic(err)
	}
	reservePool, found := k.GetReservePool(ctx, moduleName)
	if !found {
		panic(fmt.Sprintf("reserve pool for %s not found", moduleName))
	}
	inputBalance := reservePool.AmountOf(boughtCoin.Denom)
	outputBalance := reservePool.AmountOf(soldDenom)
	fee := k.GetFeeParam(ctx)

	numerator := inputBalance.Mul(boughtCoin.Amount).Mul(fee.Denominator)
	denominator := (outputBalance.Sub(boughtCoin.Amount)).Mul(fee.Numerator)
	return numerator.Div(denominator).Add(sdk.OneInt())
}

// IsDoubleSwap returns true if the trade requires a double swap.
func (k Keeper) IsDoubleSwap(ctx sdk.Context, denom1, denom2 string) bool {
	nativeDenom := k.GetNativeDenom(ctx)
	return denom1 != nativeDenom && denom2 != nativeDenom
}

// GetModuleName returns the ModuleAccount name for the provided denominations.
// The module name is in the format of 'swap:denom:denom' where the denominations
// are sorted alphabetically.
func (k Keeper) GetModuleName(denom1, denom2 string) (string, error) {
	switch strings.Compare(denom1, denom2) {
	case -1:
		return "swap:" + denom1 + ":" + denom2, nil
	case 1:
		return "swap:" + denom2 + ":" + denom1, nil
	default:
		return "", types.ErrEqualDenom(types.DefaultCodespace, "denomnations for forming module name are equal")
	}
}

func (k Keeper) GetUNIDenom(moduleName string) (string, error) {
	s := strings.Split(moduleName, ":")
	if len(s) != 3 {
		return "", types.ErrEqualDenom(types.DefaultCodespace, "denomnations for forming module name are invalid")
	}

	denom1 := strings.Split(s[1], "-")[0]
	denom2 := strings.Split(s[2], "-")[0]
	uniDenom := strings.Join([]string{"uni", denom1, denom2}, "-")

	//TODO should validate by asset module
	return uniDenom, nil
}
