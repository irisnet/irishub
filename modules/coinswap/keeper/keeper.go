package keeper

import (
	"fmt"
	"strconv"

	gogotypes "github.com/gogo/protobuf/types"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/irisnet/irismod/modules/coinswap/types"
)

// Keeper of the coinswap store
type Keeper struct {
	cdc          codec.Marshaler
	storeKey     sdk.StoreKey
	bk           types.BankKeeper
	ak           types.AccountKeeper
	paramSpace   paramstypes.Subspace
	blockedAddrs map[string]bool
}

// NewKeeper returns a coinswap keeper. It handles:
// - creating new ModuleAccounts for each trading pair
// - burning and minting liquidity coins
// - sending to and from ModuleAccounts
func NewKeeper(
	cdc codec.Marshaler,
	key sdk.StoreKey,
	paramSpace paramstypes.Subspace,
	bk types.BankKeeper,
	ak types.AccountKeeper,
	blockedAddrs map[string]bool,
) Keeper {
	// ensure coinswap module account is set
	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:     key,
		bk:           bk,
		ak:           ak,
		cdc:          cdc,
		paramSpace:   paramSpace,
		blockedAddrs: blockedAddrs,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", types.ModuleName)
}

// Swap execute swap order in specified pool
func (k Keeper) Swap(ctx sdk.Context, msg *types.MsgSwapOrder) error {
	var amount sdk.Int
	var err error

	standardDenom := k.GetStandardDenom(ctx)
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
			sdk.NewAttribute(types.AttributeValueSender, msg.Input.Address),
			sdk.NewAttribute(types.AttributeValueRecipient, msg.Output.Address),
			sdk.NewAttribute(types.AttributeValueIsBuyOrder, strconv.FormatBool(msg.IsBuyOrder)),
			sdk.NewAttribute(types.AttributeValueTokenPair, types.GetTokenPairByDenom(msg.Input.Coin.Denom, msg.Output.Coin.Denom)),
		),
	)

	return nil
}

// AddLiquidity adds liquidity to the specified pool
func (k Keeper) AddLiquidity(ctx sdk.Context, msg *types.MsgAddLiquidity) (sdk.Coin, error) {
	standardDenom := k.GetStandardDenom(ctx)
	if standardDenom == msg.MaxToken.Denom {
		return sdk.Coin{}, sdkerrors.Wrapf(types.ErrInvalidDenom,
			"MaxToken: %s should not be StandardDenom", msg.MaxToken.String())
	}
	var mintLiquidityAmt sdk.Int
	var depositToken sdk.Coin
	var standardCoin = sdk.NewCoin(standardDenom, msg.ExactStandardAmt)

	poolId := types.GetPoolId(msg.MaxToken.Denom)
	pool, exists := k.GetPool(ctx, poolId)

	// calculate amount of UNI to be minted for sender
	// and coin amount to be deposited
	if !exists {
		mintLiquidityAmt = msg.ExactStandardAmt
		if mintLiquidityAmt.LT(msg.MinLiquidity) {
			return sdk.Coin{}, sdkerrors.Wrap(types.ErrConstraintNotMet, fmt.Sprintf("liquidity amount not met, user expected: no less than %s, actual: %s", msg.MinLiquidity.String(), mintLiquidityAmt.String()))
		}
		depositToken = sdk.NewCoin(msg.MaxToken.Denom, msg.MaxToken.Amount)
		pool = k.CreatePool(ctx, msg.MaxToken.Denom)
	} else {
		balances, err := k.GetPoolBalances(ctx, pool.EscrowAddress)
		if err != nil {
			return sdk.Coin{}, err
		}

		standardReserveAmt := balances.AmountOf(standardDenom)
		tokenReserveAmt := balances.AmountOf(msg.MaxToken.Denom)
		liquidity := k.bk.GetSupply(ctx).GetTotal().AmountOf(pool.LptDenom)

		mintLiquidityAmt = (liquidity.Mul(msg.ExactStandardAmt)).Quo(standardReserveAmt)
		if mintLiquidityAmt.LT(msg.MinLiquidity) {
			return sdk.Coin{}, sdkerrors.Wrap(types.ErrConstraintNotMet, fmt.Sprintf("liquidity amount not met, user expected: no less than %s, actual: %s", msg.MinLiquidity.String(), mintLiquidityAmt.String()))
		}
		depositAmt := (tokenReserveAmt.Mul(msg.ExactStandardAmt)).Quo(standardReserveAmt).AddRaw(1)
		depositToken = sdk.NewCoin(msg.MaxToken.Denom, depositAmt)

		if depositAmt.GT(msg.MaxToken.Amount) {
			return sdk.Coin{}, sdkerrors.Wrap(types.ErrConstraintNotMet, fmt.Sprintf("token amount not met, user expected: no more than %s, actual: %s", msg.MaxToken.String(), depositToken.String()))
		}
	}

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdk.Coin{}, err
	}

	reservePoolAddress, err := sdk.AccAddressFromBech32(pool.EscrowAddress)
	if err != nil {
		return sdk.Coin{}, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAddLiquidity,
			sdk.NewAttribute(types.AttributeValueSender, msg.Sender),
			sdk.NewAttribute(types.AttributeValueTokenPair, types.GetTokenPairByDenom(msg.MaxToken.Denom, standardDenom)),
		),
	)
	return k.addLiquidity(ctx, sender, reservePoolAddress, standardCoin, depositToken, pool.LptDenom, mintLiquidityAmt)
}

func (k Keeper) addLiquidity(ctx sdk.Context,
	sender sdk.AccAddress,
	reservePoolAddress sdk.AccAddress,
	standardCoin, token sdk.Coin,
	lptDenom string,
	mintLiquidityAmt sdk.Int,
) (sdk.Coin, error) {
	depositedTokens := sdk.NewCoins(standardCoin, token)
	// transfer deposited token into coinswaps Account
	if err := k.bk.SendCoins(ctx, sender, reservePoolAddress, depositedTokens); err != nil {
		return sdk.Coin{}, err
	}

	mintToken := sdk.NewCoin(lptDenom, mintLiquidityAmt)
	mintTokens := sdk.NewCoins(mintToken)
	if err := k.bk.MintCoins(ctx, types.ModuleName, mintTokens); err != nil {
		return sdk.Coin{}, err
	}
	if err := k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, mintTokens); err != nil {
		return sdk.Coin{}, err
	}

	return mintToken, nil
}

// RemoveLiquidity removes liquidity from the specified pool
func (k Keeper) RemoveLiquidity(ctx sdk.Context, msg *types.MsgRemoveLiquidity) (sdk.Coins, error) {
	standardDenom := k.GetStandardDenom(ctx)

	pool, exists := k.GetPoolByLptDenom(ctx, msg.WithdrawLiquidity.Denom)
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrReservePoolNotExists, "liquidity pool token: %s", msg.WithdrawLiquidity.Denom)
	}

	balances, err := k.GetPoolBalances(ctx, pool.EscrowAddress)
	if err != nil {
		return nil, err
	}

	lptDenom := msg.WithdrawLiquidity.Denom
	minTokenDenom := pool.CounterpartyDenom

	standardReserveAmt := balances.AmountOf(standardDenom)
	tokenReserveAmt := balances.AmountOf(minTokenDenom)
	liquidityReserve := k.bk.GetSupply(ctx).GetTotal().AmountOf(lptDenom)
	if standardReserveAmt.LT(msg.MinStandardAmt) {
		return nil, sdkerrors.Wrap(types.ErrInsufficientFunds, fmt.Sprintf("insufficient %s funds, user expected: %s, actual: %s", standardDenom, msg.MinStandardAmt.String(), standardReserveAmt.String()))
	}
	if tokenReserveAmt.LT(msg.MinToken) {
		return nil, sdkerrors.Wrap(types.ErrInsufficientFunds, fmt.Sprintf("insufficient %s funds, user expected: %s, actual: %s", minTokenDenom, msg.MinToken.String(), tokenReserveAmt.String()))
	}
	if liquidityReserve.LT(msg.WithdrawLiquidity.Amount) {
		return nil, sdkerrors.Wrap(types.ErrInsufficientFunds, fmt.Sprintf("insufficient %s funds, user expected: %s, actual: %s", lptDenom, msg.WithdrawLiquidity.Amount.String(), liquidityReserve.String()))
	}

	// calculate amount of UNI to be burned for sender
	// and coin amount to be returned
	irisWithdrawnAmt := msg.WithdrawLiquidity.Amount.Mul(standardReserveAmt).Quo(liquidityReserve)
	tokenWithdrawnAmt := msg.WithdrawLiquidity.Amount.Mul(tokenReserveAmt).Quo(liquidityReserve)

	irisWithdrawCoin := sdk.NewCoin(standardDenom, irisWithdrawnAmt)
	tokenWithdrawCoin := sdk.NewCoin(minTokenDenom, tokenWithdrawnAmt)
	deductUniCoin := msg.WithdrawLiquidity

	if irisWithdrawCoin.Amount.LT(msg.MinStandardAmt) {
		return nil, sdkerrors.Wrap(types.ErrConstraintNotMet, fmt.Sprintf("iris amount not met, user expected: no less than %s, actual: %s", sdk.NewCoin(standardDenom, msg.MinStandardAmt).String(), irisWithdrawCoin.String()))
	}
	if tokenWithdrawCoin.Amount.LT(msg.MinToken) {
		return nil, sdkerrors.Wrap(types.ErrConstraintNotMet, fmt.Sprintf("token amount not met, user expected: no less than %s, actual: %s", sdk.NewCoin(minTokenDenom, msg.MinToken).String(), tokenWithdrawCoin.String()))
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRemoveLiquidity,
			sdk.NewAttribute(types.AttributeValueSender, msg.Sender),
			sdk.NewAttribute(types.AttributeValueTokenPair, types.GetTokenPairByDenom(minTokenDenom, standardDenom)),
		),
	)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	poolAddr, err := sdk.AccAddressFromBech32(pool.EscrowAddress)
	if err != nil {
		return nil, err
	}

	return k.removeLiquidity(ctx, poolAddr, sender, deductUniCoin, irisWithdrawCoin, tokenWithdrawCoin)
}

func (k Keeper) removeLiquidity(ctx sdk.Context, poolAddr, sender sdk.AccAddress, deductUniCoin, irisWithdrawCoin, tokenWithdrawCoin sdk.Coin) (sdk.Coins, error) {
	deltaCoins := sdk.NewCoins(deductUniCoin)

	// send liquidity vouchers to be burned from sender account to module account
	if err := k.bk.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, deltaCoins); err != nil {
		return nil, err
	}
	// burn liquidity vouchers of reserve pool from module account
	if err := k.bk.BurnCoins(ctx, types.ModuleName, deltaCoins); err != nil {
		return nil, err
	}

	// transfer withdrawn liquidity from coinswap reserve pool account to sender account
	coins := sdk.NewCoins(irisWithdrawCoin, tokenWithdrawCoin)

	return coins, k.bk.SendCoins(ctx, poolAddr, sender, coins)
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

// SetStandardDenom sets the standard denom for the coinswap module.
func (k Keeper) SetStandardDenom(ctx sdk.Context, denom string) {
	store := ctx.KVStore(k.storeKey)
	denomWrap := gogotypes.StringValue{Value: denom}
	bz := k.cdc.MustMarshalBinaryBare(&denomWrap)
	store.Set(types.KeyStandardDenom, bz)
}

// GetStandardDenom returns the standard denom of the coinswap module.
func (k Keeper) GetStandardDenom(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyStandardDenom)

	var denomWrap = gogotypes.StringValue{}
	k.cdc.MustUnmarshalBinaryBare(bz, &denomWrap)
	return denomWrap.Value
}
