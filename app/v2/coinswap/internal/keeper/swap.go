package keeper

import (
	"fmt"
	"github.com/irisnet/irishub/app/v2/coinswap/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

const PrefixReservePool = "u-%s"

func (k Keeper) SwapCoins(ctx sdk.Context, sender, recipient sdk.AccAddress, coinSold, coinBought sdk.Coin) sdk.Error {
	reservePoolName, err := k.GetReservePoolName(coinSold.Denom, coinBought.Denom)
	if err != nil {
		return err
	}

	poolAddr := getReservePoolAddr(reservePoolName)
	_, err = k.bk.SendCoins(ctx, sender, poolAddr, sdk.NewCoins(coinSold))
	if err != nil {
		return err
	}

	if recipient.Empty() {
		recipient = sender
	}
	_, err = k.bk.SendCoins(ctx, poolAddr, recipient, sdk.NewCoins(coinBought))
	return err
}

/**
Calculate the amount of another token to be received based on the exact amount of tokens sold
@param exactSoldCoin : sold coin
@param soldTokenDenom : received token's denom
@return : token amount that will to be received
*/
func (k Keeper) calculateWithExactInput(ctx sdk.Context, exactSoldCoin sdk.Coin, boughtTokenDenom string) (sdk.Int, sdk.Error) {
	reservePoolName, err := k.GetReservePoolName(exactSoldCoin.Denom, boughtTokenDenom)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	reservePool := k.GetReservePool(ctx, reservePoolName)
	if reservePool == nil {
		return sdk.ZeroInt(), types.ErrReservePoolNotExists(fmt.Sprintf("reserve pool for %s not found", reservePoolName))
	}
	inputReserve := reservePool.AmountOf(exactSoldCoin.Denom)
	outputReserve := reservePool.AmountOf(boughtTokenDenom)

	if !inputReserve.IsPositive() {
		return sdk.ZeroInt(), types.ErrInsufficientFunds(fmt.Sprintf("insufficient funds,actual:%s", inputReserve.String()))
	}
	if !outputReserve.IsPositive() {
		return sdk.ZeroInt(), types.ErrInsufficientFunds(fmt.Sprintf("insufficient funds,actual:%s", outputReserve.String()))
	}
	param := k.GetParams(ctx)

	boughtTokenAmt := getInputPrice(exactSoldCoin.Amount, inputReserve, outputReserve, param.Fee)
	return boughtTokenAmt, nil
}

/**
sell a exact amount of another token with a token,one of token denom is iris
@param input : sold MaxToken
@param output : another token received,user specified minimum amount
@param sender : address of transaction sender
@param receipt : address of  receiver bought MaxToken
@return : token amount received
*/
func (k Keeper) tradeExactInputForOutput(ctx sdk.Context, input types.Input, output types.Output) (sdk.Int, sdk.Error) {
	boughtTokenAmt, err := k.calculateWithExactInput(ctx, input.Coin, output.Coin.Denom)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	// assert that the calculated amount is less than the
	// minimum amount the buyer is willing to buy.
	if boughtTokenAmt.LT(output.Coin.Amount) {
		return sdk.ZeroInt(), types.ErrConstraintNotMet(fmt.Sprintf("token amount (%s) to be bought was less than the minimum amount (%s)", boughtTokenAmt, output.Coin.Amount))
	}
	boughtToken := sdk.NewCoin(output.Coin.Denom, boughtTokenAmt)
	err = k.SwapCoins(ctx, input.Address, output.Address, input.Coin, boughtToken)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	return boughtTokenAmt, nil
}

/**
sell a exact amount of another non-iris token with a non-iris token
@param input : sold MaxToken
@param output : another token received,user specified minimum amount
@param sender : address of transaction sender
@param receipt : address of  receiver bought MaxToken
@return : token amount received
*/
func (k Keeper) doubleTradeExactInputForOutput(ctx sdk.Context, input types.Input, output types.Output) (sdk.Int, sdk.Error) {
	nativeAmount, err := k.calculateWithExactInput(ctx, input.Coin, sdk.IrisAtto)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	nativeCoin := sdk.NewCoin(sdk.IrisAtto, nativeAmount)
	err = k.SwapCoins(ctx, input.Address, output.Address, input.Coin, nativeCoin)
	if err != nil {
		return sdk.ZeroInt(), err
	}

	boughtAmt, err := k.calculateWithExactInput(ctx, nativeCoin, output.Coin.Denom)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	boughtToken := sdk.NewCoin(output.Coin.Denom, boughtAmt)
	// assert that the calculated amount is less than the
	// minimum amount the buyer is willing to buy.
	if boughtAmt.LT(output.Coin.Amount) {
		return sdk.ZeroInt(), types.ErrConstraintNotMet(fmt.Sprintf("token amount (%s) to be bought was less than the minimum amount (%s)", boughtAmt, output.Coin.Amount))
	}

	err = k.SwapCoins(ctx, input.Address, output.Address, nativeCoin, boughtToken)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	return boughtAmt, nil
}

/**
Calculate the amount of another token to be spent based on the exact amount of tokens bought
@param exactBoughtCoin : bought coin
@param soldTokenDenom : sold token
@return : token amount that needs to be sold
*/
func (k Keeper) calculateWithExactOutput(ctx sdk.Context, exactBoughtCoin sdk.Coin, soldTokenDenom string) (sdk.Int, sdk.Error) {
	reservePoolName, err := k.GetReservePoolName(exactBoughtCoin.Denom, soldTokenDenom)
	if err != nil {
		return sdk.ZeroInt(), types.ErrReservePoolNotExists(fmt.Sprintf("reserve pool for %s not found", reservePoolName))
	}
	reservePool := k.GetReservePool(ctx, reservePoolName)
	if reservePool == nil {
		return sdk.ZeroInt(), types.ErrReservePoolNotExists(fmt.Sprintf("reserve pool for %s not found", reservePoolName))
	}
	inputReserve := reservePool.AmountOf(exactBoughtCoin.Denom)
	outputReserve := reservePool.AmountOf(soldTokenDenom)

	if !inputReserve.IsPositive() {
		return sdk.ZeroInt(), types.ErrInsufficientFunds(fmt.Sprintf("insufficient funds,actual:%s", inputReserve.String()))
	}
	if !outputReserve.IsPositive() {
		return sdk.ZeroInt(), types.ErrInsufficientFunds(fmt.Sprintf("insufficient funds,actual:%s", outputReserve.String()))
	}
	if exactBoughtCoin.Amount.GTE(outputReserve) {
		return sdk.ZeroInt(), types.ErrInsufficientFunds(fmt.Sprintf("insufficient funds,want:%s,actual:%s", exactBoughtCoin.String(), outputReserve.String()))
	}
	param := k.GetParams(ctx)

	soldTokenAmt := getOutputPrice(exactBoughtCoin.Amount, inputReserve, outputReserve, param.Fee)
	return soldTokenAmt, nil
}

/**
Purchase a exact amount of another token with a token,one of token denom is iris
@param input : bought MaxToken
@param output : another token that needs to be spent,user specified maximum
@param sender : address of transaction sender
@param receipt : address of  receiver bought MaxToken
@return : token amount that needs to be spent
*/
func (k Keeper) tradeInputForExactOutput(ctx sdk.Context, input types.Input, output types.Output) (sdk.Int, sdk.Error) {
	soldTokenAmt, err := k.calculateWithExactOutput(ctx, output.Coin, input.Coin.Denom)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	// assert that the calculated amount is less than the
	// minimum amount the buyer is willing to buy.
	if soldTokenAmt.GT(input.Coin.Amount) {
		return sdk.ZeroInt(), types.ErrConstraintNotMet(fmt.Sprintf("token amount (%s) to be bought was less than the minimum amount (%s)", soldTokenAmt, input.Coin.Amount))
	}
	soldToken := sdk.NewCoin(input.Coin.Denom, soldTokenAmt)
	err = k.SwapCoins(ctx, input.Address, output.Address, soldToken, output.Coin)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	return soldTokenAmt, nil
}

/**
Purchase a exact amount of another non-iris token with a non-iris token
@param input : bought MaxToken
@param output : another token that needs to be spent,user specified maximum
@param sender : address of transaction sender
@param receipt : address of  receiver bought MaxToken
@return : token amount that needs to be spent
*/
func (k Keeper) doubleTradeInputForExactOutput(ctx sdk.Context, input types.Input, output types.Output) (sdk.Int, sdk.Error) {
	soldIrisAmount, err := k.calculateWithExactOutput(ctx, output.Coin, sdk.IrisAtto)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	soldIrisCoin := sdk.NewCoin(sdk.IrisAtto, soldIrisAmount)

	soldTokenAmt, err := k.calculateWithExactOutput(ctx, soldIrisCoin, input.Coin.Denom)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	soldTokenCoin := sdk.NewCoin(input.Coin.Denom, soldTokenAmt)

	// assert that the calculated amount is less than the
	// minimum amount the buyer is willing to buy.
	if soldTokenAmt.GT(input.Coin.Amount) {
		return sdk.ZeroInt(), types.ErrConstraintNotMet(fmt.Sprintf("token amount (%s) to be bought was less than the minimum amount (%s)", soldTokenAmt, input.Coin.Amount))
	}

	err = k.SwapCoins(ctx, input.Address, output.Address, soldTokenCoin, soldIrisCoin)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	err = k.SwapCoins(ctx, input.Address, output.Address, soldIrisCoin, output.Coin)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	return soldTokenAmt, nil
}

// GetReservePoolName returns the reserve pool name for the provided denominations.
// The reserve pool name is in the format of 'u-denom' which the denomination
// is not iris-atto.
func (k Keeper) GetReservePoolName(denom1, denom2 string) (string, sdk.Error) {
	if denom1 == denom2 {
		return "", types.ErrEqualDenom("denomnations for forming reserve pool name are equal")
	}

	if denom1 != sdk.IrisAtto && denom2 != sdk.IrisAtto {
		return "", types.ErrIllegalDenom(fmt.Sprintf("illegal denomnations for forming reserve pool name, must have one native denom: %s", sdk.IrisAtto))
	}

	var denom = denom2
	if denom1 != sdk.IrisAtto {
		denom = denom1
	}
	return fmt.Sprintf(PrefixReservePool, denom), nil
}

// getInputPrice returns the amount of coins bought (calculated) given the input amount being sold (exact)
// The fee is included in the input coins being bought
// https://github.com/runtimeverification/verified-smart-contracts/blob/uniswap/uniswap/x-y-k.pdf
func getInputPrice(inputAmt, inputReserve, outputReserve sdk.Int, fee sdk.Rat) sdk.Int {
	deltaFee := sdk.OneRat().Sub(fee)
	inputAmtWithFee := inputAmt.Mul(deltaFee.Num())
	numerator := inputAmtWithFee.Mul(outputReserve)
	denominator := inputReserve.Mul(deltaFee.Denom()).Add(inputAmtWithFee)
	return numerator.Div(denominator)
}

// getOutputPrice returns the amount of coins sold (calculated) given the output amount being bought (exact)
// The fee is included in the output coins being bought
func getOutputPrice(outputAmt, inputReserve, outputReserve sdk.Int, fee sdk.Rat) sdk.Int {
	deltaFee := sdk.OneRat().Sub(fee)
	numerator := inputReserve.Mul(outputAmt).Mul(deltaFee.Denom())
	denominator := (outputReserve.Sub(outputAmt)).Mul(deltaFee.Num())
	return numerator.Div(denominator).Add(sdk.OneInt())
}
