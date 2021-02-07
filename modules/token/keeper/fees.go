//nolint
package keeper

import (
	"math"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/token/types"
)

// fee factor formula: (ln(len({name}))/ln{base})^{exp}
const (
	FeeFactorBase = 3
	FeeFactorExp  = 4
)

// DeductIssueTokenFee performs fee handling for issuing token
func (k Keeper) DeductIssueTokenFee(ctx sdk.Context, owner sdk.AccAddress, symbol string) error {
	// get the required issuance fee
	fee, err := k.GetTokenIssueFee(ctx, symbol)
	if err != nil {
		return err
	}
	return feeHandler(ctx, k, owner, fee)
}

// DeductMintTokenFee performs fee handling for minting token
func (k Keeper) DeductMintTokenFee(ctx sdk.Context, owner sdk.AccAddress, symbol string) error {
	// get the required minting fee
	fee, err := k.GetTokenMintFee(ctx, symbol)
	if err != nil {
		return err
	}
	return feeHandler(ctx, k, owner, fee)
}

// GetTokenIssueFee returns the token issurance fee
func (k Keeper) GetTokenIssueFee(ctx sdk.Context, symbol string) (sdk.Coin, error) {
	fee, _ := k.calcTokenIssueFee(ctx, symbol)
	token, err := k.GetToken(ctx, fee.Denom)
	if err != nil {
		return sdk.Coin{}, err
	}
	return token.ToMinCoin(sdk.NewDecCoinFromCoin(fee))
}

// GetTokenMintFee returns the token minting fee
func (k Keeper) GetTokenMintFee(ctx sdk.Context, symbol string) (sdk.Coin, error) {
	fee, params := k.calcTokenIssueFee(ctx, symbol)
	token, err := k.GetToken(ctx, fee.Denom)
	if err != nil {
		return sdk.Coin{}, err
	}

	mintFee := sdk.NewDecFromInt(fee.Amount).Mul(params.MintTokenFeeRatio).TruncateInt()
	return token.ToMinCoin(sdk.NewDecCoinFromDec(params.IssueTokenBaseFee.Denom, sdk.NewDecFromInt(mintFee)))
}

func (k Keeper) calcTokenIssueFee(ctx sdk.Context, symbol string) (sdk.Coin, types.Params) {
	// get params
	params := k.GetParamSet(ctx)
	issueTokenBaseFee := params.IssueTokenBaseFee

	// compute the fee
	feeAmt := calcFeeByBase(symbol, issueTokenBaseFee.Amount)
	if feeAmt.GT(sdk.NewDec(1)) {
		return sdk.NewCoin(issueTokenBaseFee.Denom, feeAmt.TruncateInt()), params
	}
	return sdk.NewCoin(issueTokenBaseFee.Denom, sdk.OneInt()), params
}

// feeHandler handles the fee of token
func feeHandler(ctx sdk.Context, k Keeper, feeAcc sdk.AccAddress, fee sdk.Coin) error {
	params := k.GetParamSet(ctx)
	tokenTaxRate := params.TokenTaxRate

	// compute community tax and burned coin
	communityTaxCoin := sdk.NewCoin(fee.Denom,
		sdk.NewDecFromInt(fee.Amount).Mul(tokenTaxRate).TruncateInt())
	burnedCoins := sdk.NewCoins(fee.Sub(communityTaxCoin))

	// send all fees to module account
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, feeAcc, types.ModuleName, sdk.NewCoins(fee),
	); err != nil {
		return err
	}

	// send community tax to collectedFees
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, sdk.NewCoins(communityTaxCoin)); err != nil {
		return err
	}

	// burn burnedCoin
	return k.bankKeeper.BurnCoins(ctx, types.ModuleName, burnedCoins)
}

// calcFeeByBase computes the actual fee according to the given base fee
func calcFeeByBase(name string, baseFee sdk.Int) sdk.Dec {
	feeFactor := calcFeeFactor(name)
	actualFee := sdk.NewDecFromInt(baseFee).Quo(feeFactor)

	return actualFee
}

// calcFeeFactor computes the fee factor of the given name
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
