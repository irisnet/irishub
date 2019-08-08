package keeper

import (
	"fmt"
	"github.com/irisnet/irishub/app/v2/coinswap/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

type HandlerFun func(ctx sdk.Context, input, output sdk.Coin, sender, receipt sdk.AccAddress) (sdk.Int, sdk.Error)

func (k Keeper) SwapCoins(ctx sdk.Context, sender, recipient sdk.AccAddress, coinSold, coinBought sdk.Coin) sdk.Error {
	if !k.bk.HasCoins(ctx, sender, sdk.NewCoins(coinSold)) {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("sender account does not have sufficient amount of %s to fulfill the swap order", coinSold.Denom))
	}

	reservePoolName, err := k.GetReservePoolName(coinSold.Denom, coinBought.Denom)
	if err != nil {
		return err
	}

	err = k.SendCoins(ctx, sender, reservePoolName, sdk.NewCoins(coinSold))
	if err != nil {
		return err
	}

	if recipient.Empty() {
		recipient = sender
	}
	err = k.ReceiveCoins(ctx, recipient, reservePoolName, sdk.NewCoins(coinBought))
	return err
}

func (k Keeper) GetHandler(msg types.MsgSwapOrder) HandlerFun {
	var handlerMap = map[bool]map[bool]HandlerFun{
		// BuyOrder
		true: map[bool]HandlerFun{
			// Double swap
			true: k.SwapDoubleByOutput,
			// Single swap
			false: k.SwapByOutput,
		},
		// SellOrder
		false: map[bool]HandlerFun{
			// Double swap
			true: k.SwapDoubleByInput,
			// Single swap
			false: k.SwapByInput,
		},
	}
	hMap := handlerMap[msg.IsBuyOrder]
	return hMap[k.IsDoubleSwap(msg.Input.Denom, msg.Output.Denom)]
}

/**
Calculate the amount of another token to be received based on the exact amount of tokens sold
@param exactSoldCoin : sold coin
@param soldTokenDenom : received token's denom
@return : token amount that will to be received
*/
func (k Keeper) GetPriceByInput(ctx sdk.Context, exactSoldCoin sdk.Coin, boughtTokenDenom string) (sdk.Int, sdk.Error) {
	reservePoolName, err := k.GetReservePoolName(exactSoldCoin.Denom, boughtTokenDenom)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	reservePool, found := k.GetReservePool(ctx, reservePoolName)
	if !found {
		return sdk.ZeroInt(), types.ErrReservePoolNotExists(fmt.Sprintf("reserve pool for %s not found", reservePoolName))
	}
	inputReserve := reservePool.AmountOf(exactSoldCoin.Denom)
	outputReserve := reservePool.AmountOf(boughtTokenDenom)

	if !inputReserve.IsPositive() {
		return sdk.ZeroInt(), types.ErrInsufficientFunds(fmt.Sprintf("the bought token is insufficient in the reserve Pool"))
	}
	if !outputReserve.IsPositive() {
		return sdk.ZeroInt(), types.ErrInsufficientFunds(fmt.Sprintf("the bought token is insufficient in the reserve Pool"))
	}
	param := k.GetFeeParam(ctx)

	boughtTokenAmt := GetInputPrice(exactSoldCoin.Amount, inputReserve, outputReserve, param.Fee)
	return boughtTokenAmt, nil
}

/**
sell a exact amount of another token with a token,one of token denom is iris
@param exactSoldCoin : sold Token
@param minExpect : another token received,user specified minimum amount
@param sender : address of transaction sender
@param receipt : address of  receiver bought Token
@return : token amount received
*/
func (k Keeper) SwapByInput(ctx sdk.Context, exactSoldCoin, minExpect sdk.Coin, sender, receipt sdk.AccAddress) (sdk.Int, sdk.Error) {
	boughtTokenAmt, err := k.GetPriceByInput(ctx, exactSoldCoin, minExpect.Denom)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	// assert that the calculated amount is less than the
	// minimum amount the buyer is willing to buy.
	if boughtTokenAmt.LT(minExpect.Amount) {
		return sdk.ZeroInt(), types.ErrConstraintNotMet(fmt.Sprintf("token amount (%s) to be bought was less than the minimum amount (%s)", boughtTokenAmt, minExpect.Amount))
	}
	boughtToken := sdk.NewCoin(minExpect.Denom, boughtTokenAmt)
	err = k.SwapCoins(ctx, sender, receipt, exactSoldCoin, boughtToken)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	return boughtTokenAmt, nil
}

/**
sell a exact amount of another non-iris token with a non-iris token
@param exactSoldCoin : sold Token
@param minExpect : another token received,user specified minimum amount
@param sender : address of transaction sender
@param receipt : address of  receiver bought Token
@return : token amount received
*/
func (k Keeper) SwapDoubleByInput(ctx sdk.Context, exactSoldCoin, minExpect sdk.Coin, sender, receipt sdk.AccAddress) (sdk.Int, sdk.Error) {
	nativeAmount, err := k.GetPriceByInput(ctx, exactSoldCoin, sdk.IrisAtto)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	nativeCoin := sdk.NewCoin(sdk.IrisAtto, nativeAmount)
	err = k.SwapCoins(ctx, sender, receipt, exactSoldCoin, nativeCoin)
	if err != nil {
		return sdk.ZeroInt(), err
	}

	boughtAmt, err := k.GetPriceByInput(ctx, nativeCoin, minExpect.Denom)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	boughtToken := sdk.NewCoin(minExpect.Denom, boughtAmt)
	// assert that the calculated amount is less than the
	// minimum amount the buyer is willing to buy.
	if boughtAmt.LT(minExpect.Amount) {
		return sdk.ZeroInt(), types.ErrConstraintNotMet(fmt.Sprintf("token amount (%s) to be bought was less than the minimum amount (%s)", boughtAmt, minExpect.Amount))
	}

	err = k.SwapCoins(ctx, sender, receipt, nativeCoin, boughtToken)
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
func (k Keeper) GetPriceByOutput(ctx sdk.Context, exactBoughtCoin sdk.Coin, soldTokenDenom string) (sdk.Int, sdk.Error) {
	reservePoolName, err := k.GetReservePoolName(exactBoughtCoin.Denom, soldTokenDenom)
	if err != nil {
		panic(err)
	}
	reservePool, found := k.GetReservePool(ctx, reservePoolName)
	if !found {
		panic(fmt.Sprintf("reserve pool for %s not found", reservePoolName))
	}
	inputReserve := reservePool.AmountOf(exactBoughtCoin.Denom)
	outputReserve := reservePool.AmountOf(soldTokenDenom)

	if !inputReserve.IsPositive() {
		return sdk.ZeroInt(), types.ErrInsufficientFunds(fmt.Sprintf("the bought token is insufficient in the reserve Pool"))
	}
	if !outputReserve.IsPositive() {
		return sdk.ZeroInt(), types.ErrInsufficientFunds(fmt.Sprintf("the bought token is insufficient in the reserve Pool"))
	}
	param := k.GetFeeParam(ctx)

	soldTokenAmt := GetOutputPrice(exactBoughtCoin.Amount, inputReserve, outputReserve, param.Fee)
	return soldTokenAmt, nil
}

/**
Purchase a exact amount of another token with a token,one of token denom is iris
@param exactBoughtCoin : bought Token
@param maxExpect : another token that needs to be spent,user specified maximum
@param sender : address of transaction sender
@param receipt : address of  receiver bought Token
@return : token amount that needs to be spent
*/
func (k Keeper) SwapByOutput(ctx sdk.Context, maxExpect, exactBoughtCoin sdk.Coin, sender, receipt sdk.AccAddress) (sdk.Int, sdk.Error) {
	soldTokenAmt, err := k.GetPriceByOutput(ctx, exactBoughtCoin, maxExpect.Denom)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	// assert that the calculated amount is less than the
	// minimum amount the buyer is willing to buy.
	if soldTokenAmt.GT(maxExpect.Amount) {
		return sdk.ZeroInt(), types.ErrConstraintNotMet(fmt.Sprintf("token amount (%s) to be bought was less than the minimum amount (%s)", soldTokenAmt, maxExpect.Amount))
	}
	soldToken := sdk.NewCoin(maxExpect.Denom, soldTokenAmt)
	err = k.SwapCoins(ctx, sender, receipt, soldToken, exactBoughtCoin)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	return soldTokenAmt, nil
}

/**
Purchase a exact amount of another non-iris token with a non-iris token
@param exactBoughtCoin : bought Token
@param maxExpect : another token that needs to be spent,user specified maximum
@param sender : address of transaction sender
@param receipt : address of  receiver bought Token
@return : token amount that needs to be spent
*/
func (k Keeper) SwapDoubleByOutput(ctx sdk.Context, maxExpect, exactBoughtCoin sdk.Coin, sender, receipt sdk.AccAddress) (sdk.Int, sdk.Error) {
	soldIrisAmount, err := k.GetPriceByOutput(ctx, exactBoughtCoin, sdk.IrisAtto)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	soldIrisCoin := sdk.NewCoin(sdk.IrisAtto, soldIrisAmount)

	soldTokenAmt, err := k.GetPriceByOutput(ctx, soldIrisCoin, maxExpect.Denom)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	soldTokenCoin := sdk.NewCoin(maxExpect.Denom, soldTokenAmt)

	// assert that the calculated amount is less than the
	// minimum amount the buyer is willing to buy.
	if soldTokenAmt.GT(maxExpect.Amount) {
		return sdk.ZeroInt(), types.ErrConstraintNotMet(fmt.Sprintf("token amount (%s) to be bought was less than the minimum amount (%s)", soldTokenAmt, maxExpect.Amount))
	}

	err = k.SwapCoins(ctx, sender, receipt, soldTokenCoin, soldIrisCoin)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	err = k.SwapCoins(ctx, sender, receipt, soldIrisCoin, exactBoughtCoin)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	return soldTokenAmt, nil
}

// IsDoubleSwap returns true if the trade requires a double swap.
func (k Keeper) IsDoubleSwap(denom1, denom2 string) bool {
	return denom1 != sdk.IrisAtto && denom2 != sdk.IrisAtto
}

// GetReservePoolName returns the reserve pool name for the provided denominations.
// The reserve pool name is in the format of 's-denom' which the denomination
// is not iris-atto.
func (k Keeper) GetReservePoolName(denom1, denom2 string) (string, sdk.Error) {
	if denom1 == denom2 {
		return "", types.ErrEqualDenom("denomnations for forming reserve pool name are equal")
	}

	if denom1 != sdk.IrisAtto && denom2 != sdk.IrisAtto {
		return "", types.ErrIllegalDenom(fmt.Sprintf("illegal denomnations for forming reserve pool name, must have one native denom: %s", sdk.IrisAtto))
	}

	if denom1 != sdk.IrisAtto {
		return k.GetUniDenom(denom1)
	} else {
		return k.GetUniDenom(denom2)
	}
}

// GetUniDenom returns the liquidity token denom, which is the same as the reserve pool name
func (k Keeper) GetUniDenom(denom string) (string, sdk.Error) {
	if denom == sdk.IrisAtto {
		return "", types.ErrIllegalDenom("illegal denomnation for forming liquidity token denom")
	}
	return fmt.Sprintf("s-%s", denom), nil
}

// GetInputPrice returns the amount of coins bought (calculated) given the input amount being sold (exact)
// The fee is included in the input coins being bought
// https://github.com/runtimeverification/verified-smart-contracts/blob/uniswap/uniswap/x-y-k.pdf
func GetInputPrice(inputAmt, inputReserve, outputReserve sdk.Int, fee sdk.Rat) sdk.Int {
	deltaFee := sdk.OneRat().Sub(fee)
	inputAmtWithFee := inputAmt.Mul(deltaFee.Num())
	numerator := inputAmtWithFee.Mul(outputReserve)
	denominator := inputReserve.Mul(deltaFee.Denom()).Add(inputAmtWithFee)
	return numerator.Div(denominator)
}

// GetOutputPrice returns the amount of coins sold (calculated) given the output amount being bought (exact)
// The fee is included in the output coins being bought
func GetOutputPrice(outputAmt, inputReserve, outputReserve sdk.Int, fee sdk.Rat) sdk.Int {
	deltaFee := sdk.OneRat().Sub(fee)
	numerator := inputReserve.Mul(outputAmt).Mul(deltaFee.Denom())
	denominator := (outputReserve.Sub(outputAmt)).Mul(deltaFee.Num())
	return numerator.Div(denominator).Add(sdk.OneInt())
}
