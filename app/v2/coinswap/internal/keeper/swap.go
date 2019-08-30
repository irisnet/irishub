package keeper

import (
	"fmt"
	"github.com/irisnet/irishub/app/v2/coinswap/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

func (k Keeper) swapCoins(ctx sdk.Context, sender, recipient sdk.AccAddress, coinSold, coinBought sdk.Coin) sdk.Error {
	uniId, err := types.GetUniId(coinSold.Denom, coinBought.Denom)
	if err != nil {
		return err
	}

	poolAddr := getReservePoolAddr(uniId)
	_, err = k.bk.SendCoins(ctx, sender, poolAddr, sdk.NewCoins(coinSold))
	if err != nil {
		return err
	}

	ctx.CoinFlowTags().AppendCoinFlowTag(ctx, sender.String(), poolAddr.String(), coinSold.String(), sdk.CoinSwapInputFlow, "")

	if recipient.Empty() {
		recipient = sender
	}
	_, err = k.bk.SendCoins(ctx, poolAddr, recipient, sdk.NewCoins(coinBought))

	ctx.CoinFlowTags().AppendCoinFlowTag(ctx, poolAddr.String(), recipient.String(), coinBought.String(), sdk.CoinSwapOutputFlow, "")

	return err
}

/**
Calculate the amount of another token to be received based on the exact amount of tokens sold
@param exactSoldCoin : sold coin
@param soldTokenDenom : received token's denom
@return : token amount that will to be received
*/
func (k Keeper) calculateWithExactInput(ctx sdk.Context, exactSoldCoin sdk.Coin, boughtTokenDenom string) (sdk.Int, sdk.Error) {
	uniId, err := types.GetUniId(exactSoldCoin.Denom, boughtTokenDenom)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	reservePool := k.GetReservePool(ctx, uniId)
	if reservePool == nil {
		return sdk.ZeroInt(), types.ErrReservePoolNotExists(fmt.Sprintf("reserve pool for %s not found", uniId))
	}
	inputReserve := reservePool.AmountOf(exactSoldCoin.Denom)
	outputReserve := reservePool.AmountOf(boughtTokenDenom)

	if !inputReserve.IsPositive() {
		return sdk.ZeroInt(), types.ErrInsufficientFunds(fmt.Sprintf("reserve pool insufficient funds, actual [%s%s]", inputReserve.String(), exactSoldCoin.Denom))
	}
	if !outputReserve.IsPositive() {
		return sdk.ZeroInt(), types.ErrInsufficientFunds(fmt.Sprintf("reserve pool insufficient funds, actual [%s%s]", outputReserve.String(), boughtTokenDenom))
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
	// assert that the calculated amount is more than the
	// minimum amount the buyer is willing to buy.
	if boughtTokenAmt.LT(output.Coin.Amount) {
		return sdk.ZeroInt(), types.ErrConstraintNotMet(fmt.Sprintf("the actual amount (%s) of tokens (%s) to be bought was less than the minimum amount (%s) of buyer bought", boughtTokenAmt, output.Coin.Denom, output.Coin.Amount))
	}
	boughtToken := sdk.NewCoin(output.Coin.Denom, boughtTokenAmt)
	err = k.swapCoins(ctx, input.Address, output.Address, input.Coin, boughtToken)
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
	err = k.swapCoins(ctx, input.Address, output.Address, input.Coin, nativeCoin)
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
		return sdk.ZeroInt(), types.ErrConstraintNotMet(fmt.Sprintf("the actual amount (%s) of tokens (%s) to be bought was less than the minimum amount (%s) of buyer bought", boughtAmt, output.Coin.Denom, output.Coin.Amount))
	}

	err = k.swapCoins(ctx, input.Address, output.Address, nativeCoin, boughtToken)
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
	uniId, err := types.GetUniId(exactBoughtCoin.Denom, soldTokenDenom)
	if err != nil {
		return sdk.ZeroInt(), types.ErrReservePoolNotExists(fmt.Sprintf("reserve pool not found: %s", err.Error()))
	}
	reservePool := k.GetReservePool(ctx, uniId)
	if reservePool == nil {
		return sdk.ZeroInt(), types.ErrReservePoolNotExists(fmt.Sprintf("reserve pool for %s not found", uniId))
	}
	outputReserve := reservePool.AmountOf(exactBoughtCoin.Denom)
	inputReserve := reservePool.AmountOf(soldTokenDenom)

	if !inputReserve.IsPositive() {
		return sdk.ZeroInt(), types.ErrInsufficientFunds(fmt.Sprintf("reserve pool insufficient funds, actual [%s%s]", inputReserve.String(), soldTokenDenom))
	}
	if !outputReserve.IsPositive() {
		return sdk.ZeroInt(), types.ErrInsufficientFunds(fmt.Sprintf("reserve pool insufficient funds, actual [%s%s]", outputReserve.String(), exactBoughtCoin.Denom))
	}
	if exactBoughtCoin.Amount.GTE(outputReserve) {
		return sdk.ZeroInt(), types.ErrInsufficientFunds(fmt.Sprintf("reserve pool insufficient funds, tokens[%s] to be brought could be not greater than or equal to %s", exactBoughtCoin.Denom, outputReserve.String()))
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
	// max amount the buyer is willing to sold.
	if soldTokenAmt.GT(input.Coin.Amount) {
		return sdk.ZeroInt(), types.ErrConstraintNotMet(fmt.Sprintf("the actual amount (%s) of tokens (%s) to be paid was greater than the maximum amount (%s) of sellers sold", soldTokenAmt, output.Coin.Denom, input.Coin.Amount))
	}
	soldToken := sdk.NewCoin(input.Coin.Denom, soldTokenAmt)
	err = k.swapCoins(ctx, input.Address, output.Address, soldToken, output.Coin)
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
	// max amount the buyer is willing to sold.
	if soldTokenAmt.GT(input.Coin.Amount) {
		return sdk.ZeroInt(), types.ErrConstraintNotMet(fmt.Sprintf("the actual amount (%s) of tokens (%s) to be paid was greater than the max amount (%s) of sellers sold", soldTokenAmt, input.Coin.Denom, input.Coin.Amount))
	}

	err = k.swapCoins(ctx, input.Address, output.Address, soldTokenCoin, soldIrisCoin)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	err = k.swapCoins(ctx, input.Address, output.Address, soldIrisCoin, output.Coin)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	return soldTokenAmt, nil
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
