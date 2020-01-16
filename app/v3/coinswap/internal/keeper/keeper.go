package keeper

import (
	"fmt"
	"strconv"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/app/v3/coinswap/internal/types"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

// Keeper of the coinswap store
type Keeper struct {
	cdc        *codec.Codec
	storeKey   sdk.StoreKey
	bk         types.BankKeeper
	paramSpace params.Subspace
}

// NewKeeper returns a coinswap keeper. It handles:
// - creating new ModuleAccounts for each trading pair
// - burning minting liquidity coins
// - sending to and from ModuleAccounts
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, bk types.BankKeeper, paramSpace params.Subspace) Keeper {
	return Keeper{
		storeKey:   key,
		bk:         bk,
		cdc:        cdc,
		paramSpace: paramSpace.WithTypeTable(types.ParamTypeTable()),
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("%s/%s", "iris", types.RouterKey))
}

func (k Keeper) HandleSwap(ctx sdk.Context, msg types.MsgSwapOrder) (sdk.Tags, sdk.Error) {
	tags := sdk.EmptyTags()
	var amount sdk.Int
	var err sdk.Error
	var isDoubleSwap = msg.Input.Coin.Denom != sdk.IrisAtto && msg.Output.Coin.Denom != sdk.IrisAtto

	if msg.IsBuyOrder && isDoubleSwap {
		amount, err = k.doubleTradeInputForExactOutput(ctx, msg.Input, msg.Output)
	} else if msg.IsBuyOrder && !isDoubleSwap {
		amount, err = k.tradeInputForExactOutput(ctx, msg.Input, msg.Output)
	} else if !msg.IsBuyOrder && isDoubleSwap {
		amount, err = k.doubleTradeExactInputForOutput(ctx, msg.Input, msg.Output)
	} else if !msg.IsBuyOrder && !isDoubleSwap {
		amount, err = k.tradeExactInputForOutput(ctx, msg.Input, msg.Output)
	}
	if err != nil {
		return tags, err
	}

	tags = sdk.NewTags(
		types.TagAmount, []byte(amount.String()),
		types.TagSender, []byte(msg.Input.Address.String()),
		types.TagRecipient, []byte(msg.Output.Address.String()),
		types.TagIsBuyOrder, []byte(strconv.FormatBool(msg.IsBuyOrder)),
		types.TagTokenPair, []byte(getTokenPairByDenom(msg.Input.Coin.Denom, msg.Output.Coin.Denom)),
	)

	return tags, nil
}

func (k Keeper) HandleAddLiquidity(ctx sdk.Context, msg types.MsgAddLiquidity) (sdk.Tags, sdk.Error) {
	tags := sdk.EmptyTags()
	uniID, err := types.GetUniID(sdk.IrisAtto, msg.MaxToken.Denom)
	if err != nil {
		return tags, err
	}
	uniDenom, err := types.GetUniDenom(uniID)
	if err != nil {
		return tags, err
	}

	pool, existed := k.GetPool(ctx, uniID)
	if !existed {
		_ = k.SetPool(ctx, types.NewPool(uniID, nil))
	}
	irisReserveAmt := pool.BalanceOf(sdk.IrisAtto)
	tokenReserveAmt := pool.BalanceOf(msg.MaxToken.Denom)
	liquidity := pool.BalanceOf(uniDenom)

	var mintLiquidityAmt sdk.Int
	var depositToken sdk.Coin
	var irisCoin = sdk.NewCoin(sdk.IrisAtto, msg.ExactIrisAmt)

	// calculate amount of UNI to be minted for sender
	// and coin amount to be deposited
	if liquidity.IsZero() {
		mintLiquidityAmt = msg.ExactIrisAmt
		depositToken = sdk.NewCoin(msg.MaxToken.Denom, msg.MaxToken.Amount)
	} else {
		mintLiquidityAmt = (liquidity.Mul(msg.ExactIrisAmt)).Div(irisReserveAmt)
		if mintLiquidityAmt.LT(msg.MinLiquidity) {
			return tags, types.ErrConstraintNotMet(fmt.Sprintf("liquidity amount not met, user expected: no less than %s, actual: %s", msg.MinLiquidity.String(), mintLiquidityAmt.String()))
		}
		depositAmt := (tokenReserveAmt.Mul(msg.ExactIrisAmt)).Div(irisReserveAmt).AddRaw(1)
		depositToken = sdk.NewCoin(msg.MaxToken.Denom, depositAmt)

		if depositAmt.GT(msg.MaxToken.Amount) {
			return tags, types.ErrConstraintNotMet(fmt.Sprintf("token amount not met, user expected: no more than %s, actual: %s", msg.MaxToken.String(), depositToken.String()))
		}
	}

	tags = sdk.NewTags(
		types.TagSender, []byte(msg.Sender.String()),
		types.TagTokenPair, []byte(getTokenPairByDenom(msg.MaxToken.Denom, sdk.IrisAtto)),
	)
	return tags, k.addLiquidity(ctx, msg.Sender, irisCoin, depositToken, uniID, mintLiquidityAmt)
}

func (k Keeper) HandleRemoveLiquidity(ctx sdk.Context, msg types.MsgRemoveLiquidity) (sdk.Tags, sdk.Error) {
	tags := sdk.EmptyTags()
	uniDenom := msg.WithdrawLiquidity.Denom
	uniID, err1 := sdk.GetCoinNameByDenom(uniDenom)
	if err1 != nil {
		return tags, types.ErrIllegalDenom(err1.Error())
	}
	minTokenDenom, err := types.GetCoinMinDenomFromUniDenom(uniDenom)
	if err != nil {
		return tags, err
	}

	// check if reserve pool exists
	pool, existed := k.GetPool(ctx, uniID)
	if !existed {
		return tags, types.ErrReservePoolNotExists(fmt.Sprintf("liquidity pool for %s not found", uniID))
	}

	irisReserveAmt := pool.BalanceOf(sdk.IrisAtto)
	tokenReserveAmt := pool.BalanceOf(minTokenDenom)
	liquidityReserve := pool.BalanceOf(uniDenom)
	if irisReserveAmt.LT(msg.MinIrisAmt) {
		return tags, types.ErrInsufficientFunds(fmt.Sprintf("insufficient %s funds, user expected: %s, actual: %s", sdk.IrisAtto, msg.MinIrisAmt.String(), irisReserveAmt.String()))
	}
	if tokenReserveAmt.LT(msg.MinToken) {
		return tags, types.ErrInsufficientFunds(fmt.Sprintf("insufficient %s funds, user expected: %s, actual: %s", minTokenDenom, msg.MinToken.String(), tokenReserveAmt.String()))
	}
	if liquidityReserve.LT(msg.WithdrawLiquidity.Amount) {
		return tags, types.ErrInsufficientFunds(fmt.Sprintf("insufficient %s funds, user expected: %s, actual: %s", uniDenom, msg.WithdrawLiquidity.Amount.String(), liquidityReserve.String()))
	}

	// calculate amount of UNI to be burned for sender
	// and coin amount to be returned
	irisWithdrawnAmt := msg.WithdrawLiquidity.Amount.Mul(irisReserveAmt).Div(liquidityReserve)
	tokenWithdrawnAmt := msg.WithdrawLiquidity.Amount.Mul(tokenReserveAmt).Div(liquidityReserve)

	irisWithdrawCoin := sdk.NewCoin(sdk.IrisAtto, irisWithdrawnAmt)
	tokenWithdrawCoin := sdk.NewCoin(minTokenDenom, tokenWithdrawnAmt)
	deductUniCoin := msg.WithdrawLiquidity

	if irisWithdrawCoin.Amount.LT(msg.MinIrisAmt) {
		return tags, types.ErrConstraintNotMet(fmt.Sprintf("iris amount not met, user expected: no less than %s, actual: %s", sdk.NewCoin(sdk.IrisAtto, msg.MinIrisAmt).String(), irisWithdrawCoin.String()))
	}
	if tokenWithdrawCoin.Amount.LT(msg.MinToken) {
		return tags, types.ErrConstraintNotMet(fmt.Sprintf("token amount not met, user expected: no less than %s, actual: %s", sdk.NewCoin(minTokenDenom, msg.MinToken).String(), tokenWithdrawCoin.String()))
	}
	tags = sdk.NewTags(
		types.TagSender, []byte(msg.Sender.String()),
		types.TagTokenPair, []byte(getTokenPairByDenom(minTokenDenom, sdk.IrisAtto)),
	)
	return tags, k.removeLiquidity(ctx, uniID, msg.Sender, deductUniCoin, irisWithdrawCoin, tokenWithdrawCoin)
}

// GetParams gets the parameters for the coinswap module.
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	var swapParams types.Params
	k.paramSpace.GetParamSet(ctx, &swapParams)
	return swapParams
}

// SetParams sets the parameters for the coinswap module.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

func (k Keeper) addLiquidity(ctx sdk.Context, sender sdk.AccAddress, irisCoin, token sdk.Coin, uniID string, mintLiquidityAmt sdk.Int) sdk.Error {
	// mint liquidity to pool and sender
	if err := k.MintLiquidity(ctx, sender, uniID, mintLiquidityAmt); err != nil {
		return err
	}
	coins := sdk.NewCoins(irisCoin, token)
	// transfer deposited token into liquidity pool
	return k.SendCoinsFromAccountToPool(ctx, sender, uniID, coins)
}

func (k Keeper) removeLiquidity(ctx sdk.Context, uniID string, sender sdk.AccAddress, burnUniCoin, irisWithdrawCoin, tokenWithdrawCoin sdk.Coin) sdk.Error {
	// burn liquidity from pool and sender
	if err := k.BurnLiquidity(ctx, sender, uniID, burnUniCoin.Amount); err != nil {
		return err
	}
	// transfer withdrawn liquidity from pool to sender's account
	coins := sdk.NewCoins(irisWithdrawCoin, tokenWithdrawCoin)
	return k.SendCoinsFromPoolToAccount(ctx, sender, uniID, coins)
}

func getTokenPairByDenom(inputDenom, outputDenom string) string {
	inputToken, err := sdk.GetCoinNameByDenom(inputDenom)
	if err != nil {
		panic(err)
	}
	outputToken, err := sdk.GetCoinNameByDenom(outputDenom)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s-%s", outputToken, inputToken)
}
