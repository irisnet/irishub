package keeper

import (
	"fmt"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/app/v1/coinswap/internal/types"
	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/tendermint/tendermint/libs/log"
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
	moduleAcc := k.ak.GetAccount(ctx, auth.SwapPoolAccAddr)
	if moduleAcc != nil {
		panic(fmt.Sprintf("reserve pool for %s already exists", moduleName))
	}
	k.bk.AddCoins(ctx, auth.SwapPoolAccAddr, sdk.Coins{})
}

func (k Keeper) SwapOrder(ctx sdk.Context, msg types.MsgSwapOrder) sdk.Error {
	var calculatedAmount sdk.Int
	doubleSwap := k.IsDoubleSwap(ctx, msg.Input.Denom, msg.Output.Denom)
	nativeDenom := k.GetNativeDenom(ctx)

	if msg.IsBuyOrder {
		if doubleSwap {
			nativeAmount := k.GetInputAmount(ctx, msg.Output.Amount, nativeDenom, msg.Output.Denom)
			calculatedAmount = k.GetInputAmount(ctx, nativeAmount, msg.Input.Denom, nativeDenom)
			nativeCoin := sdk.NewCoin(nativeDenom, nativeAmount)
			k.SwapCoins(ctx, msg.Sender, sdk.NewCoin(msg.Input.Denom, calculatedAmount), nativeCoin)
			k.SwapCoins(ctx, msg.Sender, nativeCoin, msg.Output)
		} else {
			calculatedAmount = k.GetInputAmount(ctx, msg.Output.Amount, msg.Input.Denom, msg.Output.Denom)
			k.SwapCoins(ctx, msg.Sender, sdk.NewCoin(msg.Input.Denom, calculatedAmount), msg.Output)
		}

		// assert that the calculated amount is less than or equal to the
		// maximum amount the buyer is willing to pay.
		if !calculatedAmount.LTE(msg.Input.Amount) {
			return types.ErrConstraintNotMet(types.DefaultCodespace, fmt.Sprintf("maximum amount (%d) to be sold was exceeded (%d)", msg.Input.Amount, calculatedAmount))
		}
	} else {
		if doubleSwap {
			nativeAmount := k.GetOutputAmount(ctx, msg.Input.Amount, msg.Input.Denom, nativeDenom)
			calculatedAmount = k.GetOutputAmount(ctx, nativeAmount, nativeDenom, msg.Output.Denom)
			nativeCoin := sdk.NewCoin(nativeDenom, nativeAmount)
			k.SwapCoins(ctx, msg.Sender, msg.Input, nativeCoin)
			k.SwapCoins(ctx, msg.Sender, nativeCoin, sdk.NewCoin(msg.Output.Denom, calculatedAmount))
		} else {
			calculatedAmount = k.GetOutputAmount(ctx, msg.Input.Amount, msg.Input.Denom, msg.Output.Denom)
			k.SwapCoins(ctx, msg.Sender, msg.Input, sdk.NewCoin(msg.Output.Denom, calculatedAmount))
		}

		// assert that the calculated amount is greater than or equal to the
		// minimum amount the sender is willing to buy.
		if !calculatedAmount.GTE(msg.Output.Amount) {
			return types.ErrConstraintNotMet(types.DefaultCodespace, fmt.Sprintf("minimum amount (%s) to be sold was not met (%s)", msg.Output.Amount, calculatedAmount))
		}
	}
	return nil
}

func (k Keeper) AddLiquidity(ctx sdk.Context, msg types.MsgAddLiquidity) sdk.Error {
	nativeDenom := k.GetNativeDenom(ctx)
	moduleName, err := k.GetModuleName(nativeDenom, msg.Deposit.Denom)

	if err != nil {
		return sdk.NewError(types.DefaultCodespace, types.CodeEqualDenom, err.Error())
	}

	// create reserve pool if it does not exist
	reservePool, found := k.GetReservePool(ctx, msg.Deposit.Denom)
	if !found {
		k.CreateReservePool(ctx, msg.Deposit.Denom)
	}

	nativeBalance := reservePool.AmountOf(nativeDenom)
	//coinBalance := reservePool.AmountOf(msg.Deposit.Denom)
	liquidityCoinBalance := reservePool.AmountOf(moduleName)

	// calculate amount of UNI to be minted for sender
	// and coin amount to be deposited
	// TODO: verify
	amtToMint := (liquidityCoinBalance.Mul(msg.DepositAmount)).Div(nativeBalance)
	coinAmountDeposited := (liquidityCoinBalance.Mul(msg.DepositAmount)).Div(nativeBalance)
	nativeCoinDeposited := sdk.NewCoin(nativeDenom, msg.DepositAmount)
	coinDeposited := sdk.NewCoin(msg.Deposit.Denom, coinAmountDeposited)

	if !k.HasCoins(ctx, msg.Sender, nativeCoinDeposited, coinDeposited) {
		return sdk.ErrInsufficientCoins("sender does not have sufficient funds to add liquidity")
	}

	// transfer deposited liquidity into coinswaps ModuleAccount
	k.SendCoins(ctx, msg.Sender, moduleName, nativeCoinDeposited, coinDeposited)

	// mint liquidity vouchers for sender
	k.MintCoins(ctx, moduleName, amtToMint)
	k.RecieveCoins(ctx, msg.Sender, moduleName, sdk.NewCoin(moduleName, amtToMint))
	return nil
}

func (k Keeper) RemoveLiquidity(ctx sdk.Context, msg types.MsgRemoveLiquidity) sdk.Error {
	nativeDenom := k.GetNativeDenom(ctx)
	moduleName, err := k.GetModuleName(nativeDenom, msg.Withdraw.Denom)
	if err != nil {
		return sdk.NewError(types.DefaultCodespace, types.CodeEqualDenom, err.Error())
	}

	// check if reserve pool exists
	reservePool, found := k.GetReservePool(ctx, msg.Withdraw.Denom)
	if !found {
		panic(fmt.Sprintf("error retrieving reserve pool for ModuleAccoint name: %s", moduleName))
	}

	nativeBalance := reservePool.AmountOf(nativeDenom)
	coinBalance := reservePool.AmountOf(msg.Withdraw.Denom)
	liquidityCoinBalance := reservePool.AmountOf(moduleName)

	// calculate amount of UNI to be burned for sender
	// and coin amount to be returned
	// TODO: verify, add amt burned
	nativeWithdrawn := msg.WithdrawAmount.Mul(nativeBalance).Div(liquidityCoinBalance)
	coinWithdrawn := msg.WithdrawAmount.Mul(coinBalance).Div(liquidityCoinBalance)
	nativeCoin := sdk.NewCoin(nativeDenom, nativeWithdrawn)
	exchangeCoin := sdk.NewCoin(msg.Withdraw.Denom, coinWithdrawn)
	amtBurned := exchangeCoin.Amount

	if !k.HasCoins(ctx, msg.Sender, sdk.NewCoin(moduleName, amtBurned)) {
		return sdk.ErrInsufficientCoins("sender does not have sufficient funds to remove liquidity")
	}

	// burn liquidity vouchers
	k.SendCoins(ctx, msg.Sender, moduleName, sdk.NewCoin(moduleName, amtBurned))
	k.BurnCoins(ctx, moduleName, amtBurned)

	// transfer withdrawn liquidity from coinswaps ModuleAccount to sender's account
	k.RecieveCoins(ctx, msg.Sender, moduleName, nativeCoin, msg.Withdraw)
	return nil
}

// HasCoins returns whether or not an account has at least coins.
func (k Keeper) HasCoins(ctx sdk.Context, addr sdk.AccAddress, coins ...sdk.Coin) bool {
	return k.bk.HasCoins(ctx, addr, coins)
}

// BurnCoins burns liquidity coins from the ModuleAccount at moduleName. The
// moduleName and denomination of the liquidity coins are the same.
func (k Keeper) BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Int) {
	_, err := k.bk.BurnCoins(ctx, auth.SwapPoolAccAddr, sdk.Coins{sdk.NewCoin(moduleName, amt)})
	if err != nil {
		panic(err)
	}
}

// MintCoins mints liquidity coins to the ModuleAccount at moduleName. The
// moduleName and denomination of the liquidity coins are the same.
func (k Keeper) MintCoins(ctx sdk.Context, moduleName string, amt sdk.Int) {
	//err := k.bk.MintCoins(ctx, moduleName, sdk.NewCoins(sdk.NewCoin(moduleName, amt)))
	//if err != nil {
	//	panic(err)
	//}
}

// SendCoin sends coins from the address to the ModuleAccount at moduleName.
func (k Keeper) SendCoins(ctx sdk.Context, addr sdk.AccAddress, moduleName string, coins ...sdk.Coin) {
	_, err := k.bk.SendCoins(ctx, addr, auth.SwapPoolAccAddr, coins)
	if err != nil {
		panic(err)
	}
}

// RecieveCoin sends coins from the ModuleAccount at moduleName to the
// address provided.
func (k Keeper) RecieveCoins(ctx sdk.Context, addr sdk.AccAddress, moduleName string, coins ...sdk.Coin) {
	_, err := k.bk.SendCoins(ctx, auth.SwapPoolAccAddr, addr, coins)
	if err != nil {
		panic(err)
	}
}

// GetReservePool returns the total balance of an reserve pool at the
// provided denomination.
func (k Keeper) GetReservePool(ctx sdk.Context, moduleName string) (coins sdk.Coins, found bool) {
	acc := k.ak.GetAccount(ctx, auth.SwapPoolAccAddr)
	if acc == nil {
		return nil, false
	}
	return acc.GetCoins(), true
}

// GetNativeDenom returns the native denomination for this module from the
// global param store.
func (k Keeper) GetNativeDenom(ctx sdk.Context) (nativeDenom string) {
	return k.GetParams(ctx).NativeDenom
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

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
