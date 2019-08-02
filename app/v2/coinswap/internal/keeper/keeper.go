package keeper

import (
	"fmt"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/app/v2/coinswap/internal/types"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/tendermint/tendermint/crypto"
)

// Keeper of the coinswap store
type Keeper struct {
	cdc        *codec.Codec
	storeKey   sdk.StoreKey
	bk         bank.Keeper
	ak         auth.AccountKeeper
	paramSpace params.Subspace
}

// NewKeeper returns a coinswap keeper. It handles:
// - creating new ModuleAccounts for each trading pair
// - burning minting liquidity coins
// - sending to and from ModuleAccounts
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, bk bank.Keeper, ak auth.AccountKeeper, paramSpace params.Subspace) Keeper {
	return Keeper{
		storeKey:   key,
		bk:         bk,
		ak:         ak,
		cdc:        cdc,
		paramSpace: paramSpace.WithTypeTable(types.ParamTypeTable()),
	}
}

// CreateReservePool initializes a new reserve pool by creating a
// ModuleAccount with minting and burning permissions.
func (k Keeper) CreateReservePool(ctx sdk.Context, moduleName string) {
	swapPoolAccAddr := getPoolAccAddr(moduleName)
	moduleAcc := k.ak.GetAccount(ctx, swapPoolAccAddr)
	if moduleAcc != nil {
		panic(fmt.Sprintf("reserve pool for %s already exists", moduleName))
	}
	k.bk.AddCoins(ctx, swapPoolAccAddr, sdk.Coins{})
}

func (k Keeper) SwapOrder(ctx sdk.Context, msg types.MsgSwapOrder) sdk.Error {
	var targetAmt sdk.Int
	doubleSwap := k.IsDoubleSwap(ctx, msg.Input.Denom, msg.Output.Denom)

	if msg.IsBuyOrder {
		if doubleSwap {
			nativeAmount := k.GetInputPrice(ctx, msg.Input, sdk.IrisAtto)
			nativeCoin := sdk.NewCoin(sdk.IrisAtto, nativeAmount)

			targetAmt = k.GetInputPrice(ctx, nativeCoin, msg.Output.Denom)
			targetCoin := sdk.NewCoin(msg.Output.Denom, targetAmt)

			k.SwapCoins(ctx, msg.Sender, msg.Input, nativeCoin)
			k.SwapCoins(ctx, msg.Sender, nativeCoin, targetCoin)

		} else {
			targetAmt = k.GetInputPrice(ctx, msg.Input, msg.Output.Denom)
			targetCoin := sdk.NewCoin(msg.Output.Denom, targetAmt)
			k.SwapCoins(ctx, msg.Sender, msg.Input, targetCoin)
		}

		// assert that the calculated amount is greater than or equal to the
		// minimum amount the buyer is willing to buy.
		if targetAmt.LT(msg.Output.Amount) {
			return types.ErrConstraintNotMet(types.DefaultCodespace, fmt.Sprintf("minimum amount (%s) to be bought was not met (%s)", targetAmt, msg.Output.Amount))
		}
	} else {
		if doubleSwap {
			nativeAmount := k.GetOutputPrice(ctx, msg.Output, sdk.IrisAtto)
			nativeCoin := sdk.NewCoin(sdk.IrisAtto, nativeAmount)

			targetAmt = k.GetOutputPrice(ctx, nativeCoin, msg.Input.Denom)
			targetCoin := sdk.NewCoin(msg.Input.Denom, targetAmt)

			k.SwapCoins(ctx, msg.Sender, targetCoin, nativeCoin)
			k.SwapCoins(ctx, msg.Sender, nativeCoin, msg.Output)
		} else {
			targetAmt = k.GetOutputPrice(ctx, msg.Output, msg.Input.Denom)
			targetCoin := sdk.NewCoin(msg.Input.Denom, targetAmt)
			k.SwapCoins(ctx, msg.Sender, targetCoin, msg.Output)
		}

		// assert that the calculated amount is greater than the
		// maximum amount the sender is willing to sell.
		if targetAmt.GT(msg.Input.Amount) {
			return types.ErrConstraintNotMet(types.DefaultCodespace, fmt.Sprintf("maximum amount (%s) to be sold was exceeded (%s)", targetAmt, msg.Input.Amount))
		}
	}
	return nil
}

func (k Keeper) AddLiquidity(ctx sdk.Context, msg types.MsgAddLiquidity) sdk.Error {
	exchangePair, err := k.GetModuleName(sdk.IrisAtto, msg.Deposit.Denom)

	if err != nil {
		return sdk.NewError(types.DefaultCodespace, types.CodeEqualDenom, err.Error())
	}

	// create reserve pool if it does not exist
	reservePool, found := k.GetReservePool(ctx, exchangePair)
	if !found {
		k.CreateReservePool(ctx, exchangePair)
	}

	nativeReserveAmt := reservePool.AmountOf(sdk.IrisAtto)
	depositedReserveAmt := reservePool.AmountOf(msg.Deposit.Denom)
	liquidity := reservePool.AmountOf(exchangePair)

	var mintLiquidityAmt sdk.Int
	var depositedCoin sdk.Coin
	var nativeCoin = sdk.NewCoin(sdk.IrisAtto, msg.DepositAmount)

	// calculate amount of UNI to be minted for sender
	// and coin amount to be deposited
	if liquidity.IsZero() {
		mintLiquidityAmt = msg.DepositAmount
		depositedCoin = sdk.NewCoin(msg.Deposit.Denom, msg.Deposit.Amount)
	} else {
		mintLiquidityAmt = (liquidity.Mul(msg.DepositAmount)).Div(nativeReserveAmt)
		if mintLiquidityAmt.LT(msg.MinReward) {
			return types.ErrLessThanMinReward(types.DefaultCodespace, fmt.Sprintf("liquidity[%s] is less than user 's min reward[%s]", mintLiquidityAmt.String(), msg.MinReward.String()))
		}

		mod := depositedReserveAmt.Mul(msg.DepositAmount).Mod(nativeReserveAmt)
		depositAmt := (depositedReserveAmt.Mul(msg.DepositAmount)).Div(nativeReserveAmt)
		if !mod.IsZero() {
			depositAmt = depositAmt.AddRaw(1)
		}

		depositedCoin = sdk.NewCoin(msg.Deposit.Denom, depositAmt)

		if depositAmt.GT(msg.Deposit.Amount) {
			return types.ErrGreaterThanMaxDeposit(types.DefaultCodespace, fmt.Sprintf("amount[%s] of token depositd is greater than user 's max deposited amount[%s]", depositedCoin.String(), msg.Deposit.String()))
		}
	}

	if !k.HasCoins(ctx, msg.Sender, nativeCoin, depositedCoin) {
		return sdk.ErrInsufficientCoins("sender does not have sufficient funds to add liquidity")
	}

	// transfer deposited liquidity into coinswaps ModuleAccount
	k.SendCoins(ctx, msg.Sender, exchangePair, nativeCoin, depositedCoin)

	// mint liquidity vouchers for sender
	coins := k.MintCoins(ctx, exchangePair, mintLiquidityAmt)
	_, _, err = k.bk.AddCoins(ctx, msg.Sender, coins.Sort())
	if err != nil {
		return sdk.NewError(types.DefaultCodespace, types.CodeEqualDenom, err.Error())
	}
	return nil
}

func (k Keeper) RemoveLiquidity(ctx sdk.Context, msg types.MsgRemoveLiquidity) sdk.Error {
	exchangePair, err := k.GetModuleName(sdk.IrisAtto, msg.Withdraw.Denom)

	if err != nil {
		return sdk.NewError(types.DefaultCodespace, types.CodeEqualDenom, err.Error())
	}

	// check if reserve pool exists
	reservePool, found := k.GetReservePool(ctx, exchangePair)
	if !found {
		return sdk.NewError(types.DefaultCodespace, types.CodeEqualDenom, fmt.Sprintf("error retrieving reserve pool for ModuleAccoint name: %s", exchangePair))
	}

	nativeReserveAmt := reservePool.AmountOf(sdk.IrisAtto)
	depositedReserveAmt := reservePool.AmountOf(msg.Withdraw.Denom)
	liquidityReserve := reservePool.AmountOf(exchangePair)

	// calculate amount of UNI to be burned for sender
	// and coin amount to be returned
	nativeWithdrawnAmt := msg.WithdrawAmount.Mul(nativeReserveAmt).Div(liquidityReserve)
	depositedWithdrawnAmt := msg.WithdrawAmount.Mul(depositedReserveAmt).Div(liquidityReserve)

	nativeWithdrawCoin := sdk.NewCoin(sdk.IrisAtto, nativeWithdrawnAmt)
	depositedWithdrawCoin := sdk.NewCoin(msg.Withdraw.Denom, depositedWithdrawnAmt)
	deltaLiquidityCoin := sdk.NewCoin(exchangePair, msg.WithdrawAmount)

	if nativeWithdrawCoin.Amount.LT(msg.MinNative) {
		return types.ErrLessThanMinWithdrawAmount(types.DefaultCodespace, fmt.Sprintf("The amount of cash available [%s] is less than the minimum amount specified [%s] by the user.", nativeWithdrawCoin.String(), sdk.NewCoin(sdk.IrisAtto, msg.MinNative).String()))
	}
	if depositedWithdrawCoin.Amount.LT(msg.Withdraw.Amount) {
		return types.ErrLessThanMinWithdrawAmount(types.DefaultCodespace, fmt.Sprintf("The amount of cash available [%s] is less than the minimum amount specified [%s] by the user.", depositedWithdrawCoin.String(), msg.Withdraw.String()))
	}

	// burn liquidity from reserve Pool
	err = k.BurnLiquidity(ctx, exchangePair, deltaLiquidityCoin)
	if err != nil {
		return sdk.NewError(types.DefaultCodespace, types.CodeEqualDenom, err.Error())
	}
	// burn liquidity from account
	_, err = k.bk.BurnCoins(ctx, msg.Sender, sdk.Coins{deltaLiquidityCoin})

	if err != nil {
		return sdk.NewError(types.DefaultCodespace, types.CodeEqualDenom, err.Error())
	}

	// transfer withdrawn liquidity from coinswaps ModuleAccount to sender's account
	k.RecieveCoins(ctx, msg.Sender, exchangePair, nativeWithdrawCoin, depositedWithdrawCoin)
	return nil
}

// HasCoins returns whether or not an account has at least coins.
func (k Keeper) HasCoins(ctx sdk.Context, addr sdk.AccAddress, coins ...sdk.Coin) bool {
	coins = sdk.Coins(coins).Sort()
	return k.bk.HasCoins(ctx, addr, coins)
}

// BurnCoins burns liquidity coins from the ModuleAccount at moduleName. The
// moduleName and denomination of the liquidity coins are the same.
func (k Keeper) BurnLiquidity(ctx sdk.Context, moduleName string, deltaCoin sdk.Coin) sdk.Error {
	swapPoolAccAddr := getPoolAccAddr(moduleName)
	if !k.HasCoins(ctx, swapPoolAccAddr, deltaCoin) {
		return sdk.ErrInsufficientCoins("sender does not have sufficient funds to remove liquidity")
	}
	coins := sdk.Coins{deltaCoin}
	_, err := k.bk.BurnCoins(ctx, swapPoolAccAddr, coins)
	if err != nil {
		return err
	}
	return nil
}

// MintCoins mints liquidity coins to the ModuleAccount at moduleName. The
// moduleName and denomination of the liquidity coins are the same.
func (k Keeper) MintCoins(ctx sdk.Context, moduleName string, amt sdk.Int) sdk.Coins {
	swapPoolAccAddr := getPoolAccAddr(moduleName)
	uniDenom, err := k.GetUNIDenom(moduleName)
	if err != nil {
		panic(err)
	}
	coins := sdk.Coins{sdk.NewCoin(uniDenom, amt)}
	_, _, err = k.bk.AddCoins(ctx, swapPoolAccAddr, coins)
	if err != nil {
		panic(err)
	}
	return coins
}

// SendCoin sends coins from the address to the ModuleAccount at moduleName.
func (k Keeper) SendCoins(ctx sdk.Context, addr sdk.AccAddress, moduleName string, coins ...sdk.Coin) {
	swapPoolAccAddr := getPoolAccAddr(moduleName)
	_, err := k.bk.SendCoins(ctx, addr, swapPoolAccAddr, sdk.Coins(coins).Sort())
	if err != nil {
		panic(err)
	}
}

// RecieveCoin sends coins from the ModuleAccount at moduleName to the
// address provided.
func (k Keeper) RecieveCoins(ctx sdk.Context, addr sdk.AccAddress, moduleName string, coins ...sdk.Coin) {
	swapPoolAccAddr := getPoolAccAddr(moduleName)
	_, err := k.bk.SendCoins(ctx, swapPoolAccAddr, addr, sdk.Coins(coins).Sort())
	if err != nil {
		panic(err)
	}
}

// GetReservePool returns the total balance of an reserve pool at the
// provided denomination.
func (k Keeper) GetReservePool(ctx sdk.Context, liquidityName string) (coins sdk.Coins, found bool) {
	swapPoolAccAddr := getPoolAccAddr(liquidityName)
	acc := k.ak.GetAccount(ctx, swapPoolAccAddr)
	if acc == nil {
		return nil, false
	}
	return acc.GetCoins(), true
}

// GetFeeParam returns the current FeeParam from the global param store
func (k Keeper) GetFeeParam(ctx sdk.Context) (feeParam types.FeeParam) {
	return k.GetParams(ctx).Fee
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

func getPoolAccAddr(liquidityName string) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte("swapPool:" + liquidityName)))
}
