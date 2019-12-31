package keeper

import (
	"fmt"
	"strconv"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/irisnet/irishub/modules/coinswap/internal/types"
)

// Keeper of the coinswap store
type Keeper struct {
	cdc        *codec.Codec
	storeKey   sdk.StoreKey
	bk         types.BankKeeper
	ak         types.AccountKeeper
	sk         types.SupplyKeeper
	paramSpace params.Subspace
}

// NewKeeper returns a coinswap keeper. It handles:
// - creating new ModuleAccounts for each trading pair
// - burning minting liquidity coins
// - sending to and from ModuleAccounts
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, bk types.BankKeeper, ak types.AccountKeeper, sk types.SupplyKeeper, paramSpace params.Subspace) Keeper {
	// ensure coinswap module account is set
	if addr := sk.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		storeKey:   key,
		bk:         bk,
		ak:         ak,
		sk:         sk,
		cdc:        cdc,
		paramSpace: paramSpace.WithKeyTable(types.ParamKeyTable()),
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("%s", types.ModuleName))
}

// Swap execute swap order in specified pool
func (k Keeper) Swap(ctx sdk.Context, msg types.MsgSwapOrder) error {
	var amount sdk.Int
	var err error

	standardDenom := k.GetParams(ctx).StandardDenom
	isDoubleSwap := (msg.Input.Coin.Denom != standardDenom) && (msg.Output.Coin.Denom != standardDenom)

	if msg.IsBuyOrder && isDoubleSwap {
		amount, err = k.doubleTradeInputForExactOutput(ctx, msg.Input, msg.Output)
	} else if msg.IsBuyOrder && !isDoubleSwap {
		amount, err = k.TradeInputForExactOutput(ctx, msg.Input, msg.Output)
	} else if !msg.IsBuyOrder && isDoubleSwap {
		amount, err = k.doubleTradeExactInputForOutput(ctx, msg.Input, msg.Output)
	} else if !msg.IsBuyOrder && !isDoubleSwap {
		amount, err = k.TradeExactInputForOutput(ctx, msg.Input, msg.Output)
	}
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSwap,
			sdk.NewAttribute(types.AttributeValueAmount, amount.String()),
			sdk.NewAttribute(types.AttributeValueSender, msg.Input.Address.String()),
			sdk.NewAttribute(types.AttributeValueRecipient, msg.Output.Address.String()),
			sdk.NewAttribute(types.AttributeValueIsBuyOrder, strconv.FormatBool(msg.IsBuyOrder)),
			sdk.NewAttribute(types.AttributeValueTokenPair, types.GetTokenPairByDenom(msg.Input.Coin.Denom, msg.Output.Coin.Denom)),
		),
	)

	return nil
}

// AddLiquidity add liquidity to specified pool
func (k Keeper) AddLiquidity(ctx sdk.Context, msg types.MsgAddLiquidity) error {
	standardDenom := k.GetParams(ctx).StandardDenom
	uniDenom, err := types.GetUniDenomFromDenom(msg.MaxToken.Denom)
	if err != nil {
		return err
	}

	reservePool := k.GetReservePool(ctx, uniDenom)
	standardReserveAmt := reservePool.AmountOf(standardDenom)
	tokenReserveAmt := reservePool.AmountOf(msg.MaxToken.Denom)
	liquidity := k.sk.GetSupply(ctx).GetTotal().AmountOf(uniDenom)

	var mintLiquidityAmt sdk.Int
	var depositToken sdk.Coin
	var standardCoin = sdk.NewCoin(standardDenom, msg.ExactStandardAmt)

	// calculate amount of UNI to be minted for sender
	// and coin amount to be deposited
	if liquidity.IsZero() {
		mintLiquidityAmt = msg.ExactStandardAmt
		depositToken = sdk.NewCoin(msg.MaxToken.Denom, msg.MaxToken.Amount)
	} else {
		mintLiquidityAmt = (liquidity.Mul(msg.ExactStandardAmt)).Quo(standardReserveAmt)
		if mintLiquidityAmt.LT(msg.MinLiquidity) {
			return sdkerrors.Wrap(types.ErrConstraintNotMet, fmt.Sprintf("liquidity amount not met, user expected: no less than %s, actual: %s", msg.MinLiquidity.String(), mintLiquidityAmt.String()))
		}
		depositAmt := (tokenReserveAmt.Mul(msg.ExactStandardAmt)).Quo(standardReserveAmt).AddRaw(1)
		depositToken = sdk.NewCoin(msg.MaxToken.Denom, depositAmt)

		if depositAmt.GT(msg.MaxToken.Amount) {
			return sdkerrors.Wrap(types.ErrConstraintNotMet, fmt.Sprintf("token amount not met, user expected: no more than %s, actual: %s", msg.MaxToken.String(), depositToken.String()))
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAddLiquidity,
			sdk.NewAttribute(types.AttributeValueSender, msg.Sender.String()),
			sdk.NewAttribute(types.AttributeValueTokenPair, types.GetTokenPairByDenom(msg.MaxToken.Denom, standardDenom)),
		),
	)

	return k.addLiquidity(ctx, msg.Sender, standardCoin, depositToken, uniDenom, mintLiquidityAmt)
}

func (k Keeper) addLiquidity(ctx sdk.Context, sender sdk.AccAddress, standardCoin, token sdk.Coin, uniDenom string, mintLiquidityAmt sdk.Int) error {
	depositedTokens := sdk.NewCoins(standardCoin, token)
	poolAddr := types.GetReservePoolAddr(uniDenom)
	// transfer deposited token into coinswaps Account
	if err := k.bk.SendCoins(ctx, sender, poolAddr, depositedTokens); err != nil {
		return err
	}

	mintToken := sdk.NewCoins(sdk.NewCoin(uniDenom, mintLiquidityAmt))
	if err := k.sk.MintCoins(ctx, types.ModuleName, mintToken); err != nil {
		return err
	}
	// send half of the liquidity vouchers from module account to sender
	if err := k.sk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, mintToken); err != nil {
		return err
	}

	return nil
}

// RemoveLiquidity remove liquidity from specified pool
func (k Keeper) RemoveLiquidity(ctx sdk.Context, msg types.MsgRemoveLiquidity) error {
	standardDenom := k.GetParams(ctx).StandardDenom
	uniDenom := msg.WithdrawLiquidity.Denom

	minTokenDenom, err := types.GetCoinDenomFromUniDenom(uniDenom)
	if err != nil {
		return err
	}

	// check if reserve pool exists
	reservePool := k.GetReservePool(ctx, uniDenom)
	if reservePool == nil {
		return sdkerrors.Wrap(types.ErrReservePoolNotExists, uniDenom)
	}

	standardReserveAmt := reservePool.AmountOf(standardDenom)
	tokenReserveAmt := reservePool.AmountOf(minTokenDenom)
	liquidityReserve := k.sk.GetSupply(ctx).GetTotal().AmountOf(uniDenom)
	if standardReserveAmt.LT(msg.MinStandardAmt) {
		return sdkerrors.Wrap(types.ErrInsufficientFunds, fmt.Sprintf("insufficient %s funds, user expected: %s, actual: %s", standardDenom, msg.MinStandardAmt.String(), standardReserveAmt.String()))
	}
	if tokenReserveAmt.LT(msg.MinToken) {
		return sdkerrors.Wrap(types.ErrInsufficientFunds, fmt.Sprintf("insufficient %s funds, user expected: %s, actual: %s", minTokenDenom, msg.MinToken.String(), tokenReserveAmt.String()))
	}
	if liquidityReserve.LT(msg.WithdrawLiquidity.Amount) {
		return sdkerrors.Wrap(types.ErrInsufficientFunds, fmt.Sprintf("insufficient %s funds, user expected: %s, actual: %s", uniDenom, msg.WithdrawLiquidity.Amount.String(), liquidityReserve.String()))
	}

	// calculate amount of UNI to be burned for sender
	// and coin amount to be returned
	irisWithdrawnAmt := msg.WithdrawLiquidity.Amount.Mul(standardReserveAmt).Quo(liquidityReserve)
	tokenWithdrawnAmt := msg.WithdrawLiquidity.Amount.Mul(tokenReserveAmt).Quo(liquidityReserve)

	irisWithdrawCoin := sdk.NewCoin(standardDenom, irisWithdrawnAmt)
	tokenWithdrawCoin := sdk.NewCoin(minTokenDenom, tokenWithdrawnAmt)
	deductUniCoin := msg.WithdrawLiquidity

	if irisWithdrawCoin.Amount.LT(msg.MinStandardAmt) {
		return sdkerrors.Wrap(types.ErrConstraintNotMet, fmt.Sprintf("iris amount not met, user expected: no less than %s, actual: %s", sdk.NewCoin(standardDenom, msg.MinStandardAmt).String(), irisWithdrawCoin.String()))
	}
	if tokenWithdrawCoin.Amount.LT(msg.MinToken) {
		return sdkerrors.Wrap(types.ErrConstraintNotMet, fmt.Sprintf("token amount not met, user expected: no less than %s, actual: %s", sdk.NewCoin(minTokenDenom, msg.MinToken).String(), tokenWithdrawCoin.String()))
	}
	poolAddr := types.GetReservePoolAddr(uniDenom)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRemoveLiquidity,
			sdk.NewAttribute(types.AttributeValueSender, msg.Sender.String()),
			sdk.NewAttribute(types.AttributeValueTokenPair, types.GetTokenPairByDenom(minTokenDenom, standardDenom)),
		),
	)

	return k.removeLiquidity(ctx, poolAddr, msg.Sender, deductUniCoin, irisWithdrawCoin, tokenWithdrawCoin)
}

func (k Keeper) removeLiquidity(ctx sdk.Context, poolAddr, sender sdk.AccAddress, deductUniCoin, irisWithdrawCoin, tokenWithdrawCoin sdk.Coin) error {

	deltaCoins := sdk.NewCoins(deductUniCoin)

	// send liquidity vouchers to be burned from sender account to module account
	if err := k.sk.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, deltaCoins); err != nil {
		return err
	}
	// burn liquidity vouchers of reserve pool form module account
	if err := k.sk.BurnCoins(ctx, types.ModuleName, deltaCoins); err != nil {
		return err
	}

	// transfer withdrawn liquidity from coinswaps reserve pool account to sender account
	coins := sdk.NewCoins(irisWithdrawCoin, tokenWithdrawCoin)

	return k.bk.SendCoins(ctx, poolAddr, sender, coins)
}

// GetReservePool returns the total balance of an reserve pool at the
// provided denomination.
func (k Keeper) GetReservePool(ctx sdk.Context, uniDenom string) (coins sdk.Coins) {
	swapPoolAccAddr := types.GetReservePoolAddr(uniDenom)
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
