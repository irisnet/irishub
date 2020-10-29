package keeper

import (
	"fmt"

	"github.com/tidwall/gjson"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/service/types"
)

// GetPrice gets the current price for the specified consumer and binding
// Note: ensure that the binding is valid
func (k Keeper) GetExchangedPrice(
	ctx sdk.Context, consumer sdk.AccAddress, binding types.ServiceBinding,
) (
	sdk.Coins, string, error,
) {
	provider, _ := sdk.AccAddressFromBech32(binding.Provider)
	pricing := k.GetPricing(ctx, binding.ServiceName, provider)

	// get discounts
	discountByTime := types.GetDiscountByTime(pricing, ctx.BlockTime())
	discountByVolume := types.GetDiscountByVolume(
		pricing, k.GetRequestVolume(ctx, consumer, binding.ServiceName, provider),
	)

	baseDenom := k.BaseDenom(ctx)
	rawDenom := pricing.Price.GetDenomByIndex(0)

	rawPrice := pricing.Price.AmountOf(rawDenom)
	price := sdk.NewDecFromInt(rawPrice).Mul(discountByTime).Mul(discountByVolume)

	realPrice := price
	if baseDenom != rawDenom {
		exchangeRateSvc, exist := k.GetModuleServiceByModuleName(types.RegisterModuleName)
		if !exist {
			return nil, rawDenom, sdkerrors.Wrapf(types.ErrInvalidModuleService, "module service not exist: %s", types.RegisterModuleName)
		}
		inputBody := fmt.Sprintf(`{"pair":"%s-%s"}`, rawDenom, baseDenom)
		input := fmt.Sprintf(`{"header":{},"body":%s`, inputBody)
		if err := types.ValidateRequestInputBody(types.OraclePriceSchemas, inputBody); err != nil {
			return nil, rawDenom, err
		}
		result, output := exchangeRateSvc.ReuquestService(ctx, input)
		if code, msg := CheckResult(result); code != "200" {
			return nil, rawDenom, sdkerrors.Wrapf(types.ErrInvalidModuleService, msg)
		}
		outputBody := gjson.Get(output, types.PATH_BODY).String()
		if err := types.ValidateResponseOutputBody(types.OraclePriceSchemas, outputBody); err != nil {
			return nil, rawDenom, err
		}
		rate, err := GetExchangeRate(outputBody)
		if err != nil {
			return nil, rawDenom, err
		}
		if rate.IsZero() {
			return nil, rawDenom, sdkerrors.Wrapf(types.ErrInvalidResponseOutputBody, "rate can not be zero")
		}
		realPrice = price.Mul(rate)
	}

	return sdk.NewCoins(sdk.NewCoin(baseDenom, realPrice.TruncateInt())), rawDenom, nil
}

func CheckResult(jsonStr string) (string, string) {
	code := gjson.Get(jsonStr, "code").String()
	msg := gjson.Get(jsonStr, "message").String()
	return code, msg
}

func GetExchangeRate(jsonStr string) (sdk.Dec, error) {
	result := gjson.Get(jsonStr, types.OraclePriceValueJSONPath)
	return sdk.NewDecFromStr(result.String())
}
