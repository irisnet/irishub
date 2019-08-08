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

// CreateReservePool initializes a new reserve pool by creating a special account.
func (k Keeper) CreateReservePool(ctx sdk.Context, reservePoolName string) {
	reservePoolAccAddr := getReservePoolAddr(reservePoolName)
	moduleAcc := k.ak.GetAccount(ctx, reservePoolAccAddr)
	if moduleAcc != nil {
		panic(fmt.Sprintf("reserve pool for %s already exists", reservePoolName))
	}
	k.bk.AddCoins(ctx, reservePoolAccAddr, sdk.Coins{})
}

func (k Keeper) Swap(ctx sdk.Context, msg types.MsgSwapOrder) (sdk.Tags, sdk.Error) {
	tags := sdk.EmptyTags()
	handler := k.GetHandler(msg)
	amount, err := handler(ctx, msg.Input, msg.Output, msg.Sender, msg.Recipient)
	if err != nil {
		return nil, err
	}
	tags.AppendTag(types.TagAmount, []byte(amount.String()))
	return tags, nil
}

func (k Keeper) AddLiquidity(ctx sdk.Context, msg types.MsgAddLiquidity) sdk.Error {
	reservePoolName, err := k.GetReservePoolName(sdk.IrisAtto, msg.Deposit.Denom)

	if err != nil {
		return err
	}

	// create reserve pool if it does not exist
	reservePool, found := k.GetReservePool(ctx, reservePoolName)
	if !found {
		k.CreateReservePool(ctx, reservePoolName)
	}

	nativeReserveAmt := reservePool.AmountOf(sdk.IrisAtto)
	depositedReserveAmt := reservePool.AmountOf(msg.Deposit.Denom)
	liquidity := reservePool.AmountOf(reservePoolName)

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
			return types.ErrConstraintNotMet(fmt.Sprintf("liquidity[%s] is less than user 's min reward[%s]", mintLiquidityAmt.String(), msg.MinReward.String()))
		}

		mod := depositedReserveAmt.Mul(msg.DepositAmount).Mod(nativeReserveAmt)
		depositAmt := (depositedReserveAmt.Mul(msg.DepositAmount)).Div(nativeReserveAmt)
		if !mod.IsZero() {
			depositAmt = depositAmt.AddRaw(1)
		}

		depositedCoin = sdk.NewCoin(msg.Deposit.Denom, depositAmt)

		if depositAmt.GT(msg.Deposit.Amount) {
			return types.ErrConstraintNotMet(fmt.Sprintf("amount[%s] of token depositd is greater than user 's max deposited amount[%s]", depositedCoin.String(), msg.Deposit.String()))
		}
	}
	if !k.bk.HasCoins(ctx, msg.Sender, sdk.NewCoins(nativeCoin, depositedCoin)) {
		return sdk.ErrInsufficientCoins("sender does not have sufficient funds to add liquidity")
	}

	// transfer deposited liquidity into coinswaps ModuleAccount
	k.SendCoins(ctx, msg.Sender, reservePoolName, sdk.NewCoins(nativeCoin, depositedCoin))

	// mint liquidity vouchers for sender
	coins := k.MintCoins(ctx, reservePoolName, mintLiquidityAmt)
	_, _, err = k.bk.AddCoins(ctx, msg.Sender, coins)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) RemoveLiquidity(ctx sdk.Context, msg types.MsgRemoveLiquidity) sdk.Error {
	reservePoolName, err := k.GetReservePoolName(sdk.IrisAtto, msg.Withdraw.Denom)

	if err != nil {
		return err
	}

	// check if reserve pool exists
	reservePool, found := k.GetReservePool(ctx, reservePoolName)
	if !found {
		return types.ErrReservePoolNotExists("")
	}

	nativeReserveAmt := reservePool.AmountOf(sdk.IrisAtto)
	depositedReserveAmt := reservePool.AmountOf(msg.Withdraw.Denom)
	liquidityReserve := reservePool.AmountOf(reservePoolName)

	// calculate amount of UNI to be burned for sender
	// and coin amount to be returned
	nativeWithdrawnAmt := msg.WithdrawAmount.Mul(nativeReserveAmt).Div(liquidityReserve)
	depositedWithdrawnAmt := msg.WithdrawAmount.Mul(depositedReserveAmt).Div(liquidityReserve)

	nativeWithdrawCoin := sdk.NewCoin(sdk.IrisAtto, nativeWithdrawnAmt)
	depositedWithdrawCoin := sdk.NewCoin(msg.Withdraw.Denom, depositedWithdrawnAmt)
	deltaLiquidityCoin := sdk.NewCoin(reservePoolName, msg.WithdrawAmount)

	if nativeWithdrawCoin.Amount.LT(msg.MinNative) {
		return types.ErrConstraintNotMet(fmt.Sprintf("The amount of cash available [%s] is less than the minimum amount specified [%s] by the user.", nativeWithdrawCoin.String(), sdk.NewCoin(sdk.IrisAtto, msg.MinNative).String()))
	}
	if depositedWithdrawCoin.Amount.LT(msg.Withdraw.Amount) {
		return types.ErrConstraintNotMet(fmt.Sprintf("The amount of cash available [%s] is less than the minimum amount specified [%s] by the user.", depositedWithdrawCoin.String(), msg.Withdraw.String()))
	}

	// burn liquidity from reserve Pool
	err = k.BurnLiquidity(ctx, reservePoolName, deltaLiquidityCoin)
	if err != nil {
		return err
	}
	// burn liquidity from account
	_, err = k.bk.BurnCoins(ctx, msg.Sender, sdk.Coins{deltaLiquidityCoin})

	if err != nil {
		return err
	}

	// transfer withdrawn liquidity from coinswaps ModuleAccount to sender's account
	coins := sdk.NewCoins(nativeWithdrawCoin, depositedWithdrawCoin)
	k.ReceiveCoins(ctx, msg.Sender, reservePoolName, coins)
	return nil
}

// BurnCoins burns liquidity coins from the ModuleAccount at reservePoolName. The
// reservePoolName and denomination of the liquidity coins are the same.
func (k Keeper) BurnLiquidity(ctx sdk.Context, reservePoolName string, deltaCoin sdk.Coin) sdk.Error {
	swapPoolAccAddr := getReservePoolAddr(reservePoolName)
	if !k.bk.HasCoins(ctx, swapPoolAccAddr, sdk.NewCoins(deltaCoin)) {
		return sdk.ErrInsufficientCoins("sender does not have sufficient funds to remove liquidity")
	}
	coins := sdk.NewCoins(deltaCoin)
	_, err := k.bk.BurnCoins(ctx, swapPoolAccAddr, coins)
	if err != nil {
		return err
	}
	return nil
}

// MintCoins mints liquidity coins to the ModuleAccount at reservePoolName. The
// reservePoolName and denomination of the liquidity coins are the same.
func (k Keeper) MintCoins(ctx sdk.Context, reservePoolName string, amt sdk.Int) sdk.Coins {
	reservePoolAccAddr := getReservePoolAddr(reservePoolName)
	coins := sdk.NewCoins(sdk.NewCoin(reservePoolName, amt))
	_, _, err := k.bk.AddCoins(ctx, reservePoolAccAddr, coins)
	if err != nil {
		panic(err)
	}
	return coins
}

// SendCoin sends coins from the address to the ModuleAccount at reservePoolName.
func (k Keeper) SendCoins(ctx sdk.Context, addr sdk.AccAddress, reservePoolName string, coins sdk.Coins) sdk.Error {
	swapPoolAccAddr := getReservePoolAddr(reservePoolName)
	_, err := k.bk.SendCoins(ctx, addr, swapPoolAccAddr, coins)
	return err
}

// RecieveCoin sends coins from the ModuleAccount at reservePoolName to the
// address provided.
func (k Keeper) ReceiveCoins(ctx sdk.Context, addr sdk.AccAddress, reservePoolName string, coins sdk.Coins) sdk.Error {
	swapPoolAccAddr := getReservePoolAddr(reservePoolName)
	_, err := k.bk.SendCoins(ctx, swapPoolAccAddr, addr, coins)
	return err
}

// GetReservePool returns the total balance of an reserve pool at the
// provided denomination.
func (k Keeper) GetReservePool(ctx sdk.Context, reservePoolName string) (coins sdk.Coins, found bool) {
	swapPoolAccAddr := getReservePoolAddr(reservePoolName)
	acc := k.ak.GetAccount(ctx, swapPoolAccAddr)
	if acc == nil {
		return nil, false
	}
	return acc.GetCoins(), true
}

// GetFeeParam returns the current FeeParam from the global param store
func (k Keeper) GetFeeParam(ctx sdk.Context) (feeParam types.Params) {
	return k.GetParams(ctx)
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
