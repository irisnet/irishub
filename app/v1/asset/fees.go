//nolint
package asset

import (
	"fmt"
	"math"
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
	feeFactor := calcFeeFactor(name)
	actualFee := int64(math.Round(float64(baseFee.Int64()) / feeFactor))

	return sdk.NewInt(actualFee)
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

// GatewayFeeOutput is for the gateway fee query output
type GatewayFeeOutput struct {
	Exist bool     `exist` // indicate if the gateway has existed
	Fee   sdk.Coin `fee`   // creation fee
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
