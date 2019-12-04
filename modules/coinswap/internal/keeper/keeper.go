package keeper

import (
	"fmt"
	"strconv"

	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/irisnet/irishub/config"
	"github.com/irisnet/irishub/modules/coinswap/internal/types"
)

// Keeper of the coinswap store
type Keeper struct {
	cdc        *codec.Codec
	storeKey   sdk.StoreKey
	bk         types.BankKeeper
	ak         types.AccountKeeper
	paramSpace params.Subspace
}

// NewKeeper returns a coinswap keeper. It handles:
// - creating new ModuleAccounts for each trading pair
// - burning minting liquidity coins
// - sending to and from ModuleAccounts
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, bk types.BankKeeper, ak types.AccountKeeper, paramSpace params.Subspace) Keeper {
	return Keeper{
		storeKey:   key,
		bk:         bk,
		ak:         ak,
		cdc:        cdc,
		paramSpace: paramSpace.WithKeyTable(types.ParamKeyTable()),
	}
}

func (k Keeper) Swap(ctx sdk.Context, msg types.MsgSwapOrder) sdk.Error {
	var amount sdk.Int
	var err sdk.Error
	var isDoubleSwap = msg.Input.Coin.Denom != config.IrisAtto && msg.Output.Coin.Denom != config.IrisAtto

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
		return err
	}

	swapEvent := sdk.NewEvent(
		types.EventSwap,
		sdk.NewAttribute(types.AttributeValueAmount, amount.String()),
		sdk.NewAttribute(types.AttributeValueSender, msg.Input.Address.String()),
		sdk.NewAttribute(types.AttributeValueRecipient, msg.Output.Address.String()),
		sdk.NewAttribute(types.AttributeValueIsBuyOrder, strconv.FormatBool(msg.IsBuyOrder)),
		sdk.NewAttribute(types.AttributeValueTokenPair, getTokenPairByDenom(msg.Input.Coin.Denom, msg.Output.Coin.Denom)),
	)
	ctx.EventManager().EmitEvents(sdk.Events{swapEvent})

	return nil
}

func (k Keeper) AddLiquidity(ctx sdk.Context, msg types.MsgAddLiquidity) sdk.Error {
	uniId, err := types.GetUniId(config.IrisAtto, msg.MaxToken.Denom)
	if err != nil {
		return err
	}
	uniDenom, err := types.GetUniDenom(uniId)
	if err != nil {
		return err
	}

	reservePool := k.GetReservePool(ctx, uniId)
	irisReserveAmt := reservePool.AmountOf(config.IrisAtto)
	tokenReserveAmt := reservePool.AmountOf(msg.MaxToken.Denom)
	liquidity := reservePool.AmountOf(uniDenom)

	var mintLiquidityAmt sdk.Int
	var depositToken sdk.Coin
	var irisCoin = sdk.NewCoin(config.IrisAtto, msg.ExactIrisAmt)

	// calculate amount of UNI to be minted for sender
	// and coin amount to be deposited
	if liquidity.IsZero() {
		mintLiquidityAmt = msg.ExactIrisAmt
		depositToken = sdk.NewCoin(msg.MaxToken.Denom, msg.MaxToken.Amount)
	} else {
		mintLiquidityAmt = (liquidity.Mul(msg.ExactIrisAmt)).Quo(irisReserveAmt)
		if mintLiquidityAmt.LT(msg.MinLiquidity) {
			return types.ErrConstraintNotMet(fmt.Sprintf("liquidity amount not met, user expected: no less than %s, actual: %s", msg.MinLiquidity.String(), mintLiquidityAmt.String()))
		}
		depositAmt := (tokenReserveAmt.Mul(msg.ExactIrisAmt)).Quo(irisReserveAmt).AddRaw(1)
		depositToken = sdk.NewCoin(msg.MaxToken.Denom, depositAmt)

		if depositAmt.GT(msg.MaxToken.Amount) {
			return types.ErrConstraintNotMet(fmt.Sprintf("token amount not met, user expected: no more than %s, actual: %s", msg.MaxToken.String(), depositToken.String()))
		}
	}

	addLiquidityEvent := sdk.NewEvent(
		types.EventAddLiquidity,
		sdk.NewAttribute(types.AttributeValueSender, msg.Sender.String()),
		sdk.NewAttribute(types.AttributeValueTokenPair, getTokenPairByDenom(msg.MaxToken.Denom, config.IrisAtto)),
	)
	ctx.EventManager().EmitEvents(sdk.Events{addLiquidityEvent})

	return k.addLiquidity(ctx, msg.Sender, irisCoin, depositToken, uniId, mintLiquidityAmt)
}

func (k Keeper) addLiquidity(ctx sdk.Context, sender sdk.AccAddress, irisCoin, token sdk.Coin, uniId string, mintLiquidityAmt sdk.Int) sdk.Error {
	depositedTokens := sdk.NewCoins(irisCoin, token)
	poolAddr := getReservePoolAddr(uniId)
	// transfer deposited token into coinswaps Account
	err := k.bk.SendCoins(ctx, sender, poolAddr, depositedTokens)
	if err != nil {
		return err
	}

	uniDenom, err := types.GetUniDenom(uniId)
	if err != nil {
		return err
	}
	// mint liquidity vouchers for reserve Pool
	mintToken := sdk.NewCoins(sdk.NewCoin(uniDenom, mintLiquidityAmt))
	_, _ = k.bk.AddCoins(ctx, poolAddr, mintToken)

	// mint liquidity vouchers for sender
	_, _ = k.bk.AddCoins(ctx, sender, mintToken)

	return nil
}

func (k Keeper) RemoveLiquidity(ctx sdk.Context, msg types.MsgRemoveLiquidity) sdk.Error {
	uniDenom := msg.WithdrawLiquidity.Denom
	uniId, err1 := sdk.GetCoinNameByDenom(uniDenom)
	if err1 != nil {
		return types.ErrIllegalDenom(err1.Error())
	}
	minTokenDenom, err := types.GetCoinMinDenomFromUniDenom(uniDenom)
	if err != nil {
		return err
	}

	// check if reserve pool exists
	reservePool := k.GetReservePool(ctx, uniId)
	if reservePool == nil {
		return types.ErrReservePoolNotExists("")
	}

	irisReserveAmt := reservePool.AmountOf(config.IrisAtto)
	tokenReserveAmt := reservePool.AmountOf(minTokenDenom)
	liquidityReserve := reservePool.AmountOf(uniDenom)
	if irisReserveAmt.LT(msg.MinIrisAmt) {
		return types.ErrInsufficientFunds(fmt.Sprintf("insufficient %s funds, user expected: %s, actual: %s", config.IrisAtto, msg.MinIrisAmt.String(), irisReserveAmt.String()))
	}
	if tokenReserveAmt.LT(msg.MinToken) {
		return types.ErrInsufficientFunds(fmt.Sprintf("insufficient %s funds, user expected: %s, actual: %s", minTokenDenom, msg.MinToken.String(), tokenReserveAmt.String()))
	}
	if liquidityReserve.LT(msg.WithdrawLiquidity.Amount) {
		return types.ErrInsufficientFunds(fmt.Sprintf("insufficient %s funds, user expected: %s, actual: %s", uniDenom, msg.WithdrawLiquidity.Amount.String(), liquidityReserve.String()))
	}

	// calculate amount of UNI to be burned for sender
	// and coin amount to be returned
	irisWithdrawnAmt := msg.WithdrawLiquidity.Amount.Mul(irisReserveAmt).Quo(liquidityReserve)
	tokenWithdrawnAmt := msg.WithdrawLiquidity.Amount.Mul(tokenReserveAmt).Quo(liquidityReserve)

	irisWithdrawCoin := sdk.NewCoin(config.IrisAtto, irisWithdrawnAmt)
	tokenWithdrawCoin := sdk.NewCoin(minTokenDenom, tokenWithdrawnAmt)
	deductUniCoin := msg.WithdrawLiquidity

	if irisWithdrawCoin.Amount.LT(msg.MinIrisAmt) {
		return types.ErrConstraintNotMet(fmt.Sprintf("iris amount not met, user expected: no less than %s, actual: %s", sdk.NewCoin(config.IrisAtto, msg.MinIrisAmt).String(), irisWithdrawCoin.String()))
	}
	if tokenWithdrawCoin.Amount.LT(msg.MinToken) {
		return types.ErrConstraintNotMet(fmt.Sprintf("token amount not met, user expected: no less than %s, actual: %s", sdk.NewCoin(minTokenDenom, msg.MinToken).String(), tokenWithdrawCoin.String()))
	}
	poolAddr := getReservePoolAddr(uniId)

	removeLiquidityEvent := sdk.NewEvent(
		types.EventRemoveLiquidity,
		sdk.NewAttribute(types.AttributeValueSender, msg.Sender.String()),
		sdk.NewAttribute(types.AttributeValueTokenPair, getTokenPairByDenom(minTokenDenom, config.IrisAtto)),
	)
	ctx.EventManager().EmitEvents(sdk.Events{removeLiquidityEvent})

	return k.removeLiquidity(ctx, poolAddr, msg.Sender, deductUniCoin, irisWithdrawCoin, tokenWithdrawCoin)
}

func (k Keeper) removeLiquidity(ctx sdk.Context, poolAddr, sender sdk.AccAddress, deductUniCoin, irisWithdrawCoin, tokenWithdrawCoin sdk.Coin) sdk.Error {
	// burn liquidity from reserve Pool
	deltaCoins := sdk.NewCoins(deductUniCoin)
	_, err := k.bk.SubtractCoins(ctx, poolAddr, deltaCoins)
	if err != nil {
		return err
	}

	// burn liquidity from account
	_, err = k.bk.SubtractCoins(ctx, sender, deltaCoins)
	if err != nil {
		return err
	}

	// transfer withdrawn liquidity from coinswaps special account to sender's account
	coins := sdk.NewCoins(irisWithdrawCoin, tokenWithdrawCoin)
	err = k.bk.SendCoins(ctx, poolAddr, sender, coins)

	return err
}

// GetReservePool returns the total balance of an reserve pool at the
// provided denomination.
func (k Keeper) GetReservePool(ctx sdk.Context, uniId string) (coins sdk.Coins) {
	swapPoolAccAddr := getReservePoolAddr(uniId)
	acc := k.ak.GetAccount(ctx, swapPoolAccAddr)
	if acc == nil {
		return nil
	}
	return acc.GetCoins()
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

func (k Keeper) Init(ctx sdk.Context) {
	paramSet := types.DefaultParams()
	k.paramSpace.SetParamSet(ctx, &paramSet)
}

func getReservePoolAddr(uniDenom string) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(uniDenom)))
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
