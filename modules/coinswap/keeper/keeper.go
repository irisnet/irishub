package keeper

import (
	"fmt"
	"math/big"
	"strconv"

	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/libs/log"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/coinswap/types"
)

// Keeper of the coinswap store
type Keeper struct {
	cdc              codec.BinaryCodec
	storeKey         storetypes.StoreKey
	bk               types.BankKeeper
	ak               types.AccountKeeper
	feeCollectorName string
	authority        string
	blockedAddrs     map[string]bool
}

// NewKeeper returns a coinswap keeper. It handles:
// - creating new ModuleAccounts for each trading pair
// - burning and minting liquidity coins
// - sending to and from ModuleAccounts
func NewKeeper(
	cdc codec.BinaryCodec,
	key storetypes.StoreKey,
	bk types.BankKeeper,
	ak types.AccountKeeper,
	feeCollectorName string,
	authority string,
) Keeper {
	// ensure coinswap module account is set
	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		storeKey:         key,
		bk:               bk,
		ak:               ak,
		cdc:              cdc,
		blockedAddrs:     bk.GetBlockedAddresses(),
		feeCollectorName: feeCollectorName,
		authority:        authority,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", types.ModuleName)
}

// Swap execute swap order in specified pool
func (k Keeper) Swap(ctx sdk.Context, msg *types.MsgSwapOrder) error {
	var amount sdkmath.Int
	var err error

	standardDenom := k.GetStandardDenom(ctx)
	isDoubleSwap := (msg.Input.Coin.Denom != standardDenom) &&
		(msg.Output.Coin.Denom != standardDenom)

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
			sdk.NewAttribute(
				types.AttributeValueTokenPair,
				types.GetTokenPairByDenom(msg.Input.Coin.Denom, msg.Output.Coin.Denom),
			),
		),
	)

	return nil
}

// AddLiquidity adds liquidity to the specified pool
func (k Keeper) AddLiquidity(ctx sdk.Context, msg *types.MsgAddLiquidity) (sdk.Coin, error) {
	standardDenom := k.GetStandardDenom(ctx)
	if standardDenom == msg.MaxToken.Denom {
		return sdk.Coin{}, errorsmod.Wrapf(types.ErrInvalidDenom,
			"MaxToken: %s should not be StandardDenom", msg.MaxToken.String())
	}

	var (
		mintLiquidityAmt sdkmath.Int
		depositToken     sdk.Coin
		standardCoin     = sdk.NewCoin(standardDenom, msg.ExactStandardAmt)
	)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdk.Coin{}, err
	}

	poolID := types.GetPoolId(msg.MaxToken.Denom)
	pool, exists := k.GetPool(ctx, poolID)

	// calculate amount of UNI to be minted for sender
	// and coin amount to be deposited
	if !exists {
		// deduct the user's fee for creating a Liquidity pool
		if err := k.DeductPoolCreationFee(ctx, sender); err != nil {
			return sdk.Coin{}, err
		}

		mintLiquidityAmt = msg.ExactStandardAmt
		if mintLiquidityAmt.LT(msg.MinLiquidity) {
			return sdk.Coin{}, errorsmod.Wrapf(
				types.ErrConstraintNotMet,
				"liquidity amount not met, user expected: no less than %s, actual: %s",
				msg.MinLiquidity.String(),
				mintLiquidityAmt.String(),
			)
		}
		depositToken = sdk.NewCoin(msg.MaxToken.Denom, msg.MaxToken.Amount)
		pool = k.CreatePool(ctx, msg.MaxToken.Denom)
		return k.addLiquidity(
			ctx,
			sender,
			pool.EscrowAddress,
			standardCoin,
			depositToken,
			pool.LptDenom,
			mintLiquidityAmt,
		)
	}

	balances, err := k.GetPoolBalances(ctx, pool.EscrowAddress)
	if err != nil {
		return sdk.Coin{}, err
	}

	//pool exist but has no balances,so do same operations as firist addLiquidity(but without creating pool)
	if balances == nil || balances.IsZero() {
		mintLiquidityAmt = msg.ExactStandardAmt
		if mintLiquidityAmt.LT(msg.MinLiquidity) {
			return sdk.Coin{}, errorsmod.Wrapf(
				types.ErrConstraintNotMet,
				"liquidity amount not met, user expected: no less than %s, actual: %s",
				msg.MinLiquidity.String(),
				mintLiquidityAmt.String(),
			)
		}
		depositToken = sdk.NewCoin(msg.MaxToken.Denom, msg.MaxToken.Amount)
		return k.addLiquidity(
			ctx,
			sender,
			pool.EscrowAddress,
			standardCoin,
			depositToken,
			pool.LptDenom,
			mintLiquidityAmt,
		)
	}

	// add liquidity
	standardReserveAmt := balances.AmountOf(standardDenom)
	tokenReserveAmt := balances.AmountOf(msg.MaxToken.Denom)
	liquidity := k.bk.GetSupply(ctx, pool.LptDenom).Amount
	if standardReserveAmt.IsZero() || tokenReserveAmt.IsZero() || liquidity.IsZero() {
		return sdk.Coin{},
			errorsmod.Wrapf(types.ErrConstraintNotMet, "liquidity pool invalid")
	}

	mintLiquidityAmt = (liquidity.Mul(msg.ExactStandardAmt)).Quo(standardReserveAmt)
	if mintLiquidityAmt.LT(msg.MinLiquidity) {
		return sdk.Coin{}, errorsmod.Wrapf(
			types.ErrConstraintNotMet,
			"liquidity amount not met, user expected: no less than %s, actual: %s",
			msg.MinLiquidity.String(),
			mintLiquidityAmt.String(),
		)
	}
	depositAmt := (tokenReserveAmt.Mul(msg.ExactStandardAmt)).Quo(standardReserveAmt).AddRaw(1)
	depositToken = sdk.NewCoin(msg.MaxToken.Denom, depositAmt)

	if depositAmt.GT(msg.MaxToken.Amount) {
		return sdk.Coin{},
			errorsmod.Wrapf(
				types.ErrConstraintNotMet,
				"token amount not met, user expected: no more than %s, actual: %s",
				msg.MaxToken.String(),
				depositToken.String(),
			)
	}
	return k.addLiquidity(
		ctx,
		sender,
		pool.EscrowAddress,
		standardCoin,
		depositToken,
		pool.LptDenom,
		mintLiquidityAmt,
	)
}

func (k Keeper) addLiquidity(ctx sdk.Context,
	sender sdk.AccAddress,
	poolAddress string,
	standardCoin, token sdk.Coin,
	lptDenom string,
	mintLiquidityAmt sdkmath.Int,
) (sdk.Coin, error) {
	reservePoolAddress, err := sdk.AccAddressFromBech32(poolAddress)
	if err != nil {
		return sdk.Coin{}, err
	}

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

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAddLiquidity,
			sdk.NewAttribute(types.AttributeValueSender, sender.String()),
			sdk.NewAttribute(
				types.AttributeValueTokenPair,
				types.GetTokenPairByDenom(token.Denom, standardCoin.Denom),
			),
		),
	)
	return mintToken, nil
}

// AddUnilateralLiquidity adds liquidity unilaterally to the specified pool
func (k Keeper) AddUnilateralLiquidity(
	ctx sdk.Context,
	msg *types.MsgAddUnilateralLiquidity,
) (sdk.Coin, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdk.Coin{}, err
	}

	poolID := types.GetPoolId(msg.CounterpartyDenom)
	pool, exist := k.GetPool(ctx, poolID)
	if !exist {
		return sdk.Coin{}, errorsmod.Wrapf(
			types.ErrReservePoolNotExists,
			"liquidity pool: %s ",
			poolID,
		)
	}

	poolAddr, err := sdk.AccAddressFromBech32(pool.EscrowAddress)
	if err != nil {
		return sdk.Coin{}, err
	}

	balances, err := k.GetPoolBalances(ctx, pool.EscrowAddress)
	if err != nil {
		return sdk.Coin{}, err
	}

	if msg.ExactToken.Denom != msg.CounterpartyDenom &&
		msg.ExactToken.Denom != k.GetStandardDenom(ctx) {
		return sdk.Coin{}, errorsmod.Wrapf(
			types.ErrInvalidDenom,
			"liquidity pool %s has no %s",
			poolID,
			msg.ExactToken.Denom,
		)
	}

	if balances == nil || balances.IsZero() {
		return sdk.Coin{}, errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"When the liquidity is empty, can not add unilateral liquidity",
		)
	}

	// square = ( token_balance + ( 1- fee_unilateral ) * exact_token ) / token_balance * lpt_balance^2
	// 1 - fee_unilateral = numerator / denominator
	tokenBalanceAmt := balances.AmountOf(msg.ExactToken.Denom)
	lptBalanceAmt := k.bk.GetSupply(ctx, pool.LptDenom).Amount
	exactTokenAmt := msg.ExactToken.Amount

	deltaFeeUnilateral := sdk.OneDec().Sub(k.GetParams(ctx).UnilateralLiquidityFee)
	numerator := sdkmath.NewIntFromBigInt(deltaFeeUnilateral.BigInt())
	denominator := sdkmath.NewIntWithDecimal(1, sdk.Precision)

	square := denominator.Mul(tokenBalanceAmt).
		Add(numerator.Mul(exactTokenAmt)).
		Mul(lptBalanceAmt).
		Mul(lptBalanceAmt).
		Quo(denominator.Mul(tokenBalanceAmt))

	// lpt = square^0.5 - lpt_balance
	var squareBigInt = &big.Int{}
	squareBigInt.Sqrt(square.BigInt())
	mintLptAmt := sdkmath.NewIntFromBigInt(squareBigInt).Sub(lptBalanceAmt)

	if mintLptAmt.LT(msg.MinLiquidity) {
		return sdk.Coin{}, errorsmod.Wrapf(
			types.ErrConstraintNotMet,
			"liquidity amount not met, user expected: no less than %s, actual: %s",
			msg.MinLiquidity.String(),
			mintLptAmt.String(),
		)
	}

	// event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAddUnilateralLiquidity,
			sdk.NewAttribute(types.AttributeValueSender, msg.Sender),
			sdk.NewAttribute(types.AttributeValueTokenUnilateral, msg.ExactToken.Denom),
			sdk.NewAttribute(types.AttributeValueLptDenom, pool.LptDenom),
		),
	)

	return k.addUnilateralLiquidity(
		ctx,
		sender,
		poolAddr,
		msg.ExactToken,
		pool.LptDenom,
		mintLptAmt,
	)
}

func (k Keeper) addUnilateralLiquidity(ctx sdk.Context,
	sender sdk.AccAddress,
	poolAddr sdk.AccAddress,
	exactToken sdk.Coin,
	lptDenom string,
	mintLptAmt sdkmath.Int,
) (sdk.Coin, error) {
	// add liquidity
	exactCoins := sdk.NewCoins(exactToken)
	if err := k.bk.SendCoins(ctx, sender, poolAddr, exactCoins); err != nil {
		return sdk.Coin{}, err
	}

	// mint and send lpt
	mintLpt := sdk.NewCoin(lptDenom, mintLptAmt)
	mintLpts := sdk.NewCoins(mintLpt)
	if err := k.bk.MintCoins(ctx, types.ModuleName, mintLpts); err != nil {
		return sdk.Coin{}, err
	}
	if err := k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, mintLpts); err != nil {
		return sdk.Coin{}, err
	}

	return mintLpt, nil
}

// RemoveLiquidity removes liquidity from the specified pool
func (k Keeper) RemoveLiquidity(ctx sdk.Context, msg *types.MsgRemoveLiquidity) (sdk.Coins, error) {
	standardDenom := k.GetStandardDenom(ctx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	pool, exists := k.GetPoolByLptDenom(ctx, msg.WithdrawLiquidity.Denom)
	if !exists {
		return nil, errorsmod.Wrapf(
			types.ErrReservePoolNotExists,
			"liquidity pool token: %s",
			msg.WithdrawLiquidity.Denom,
		)
	}

	balances, err := k.GetPoolBalances(ctx, pool.EscrowAddress)
	if err != nil {
		return nil, err
	}

	lptDenom := msg.WithdrawLiquidity.Denom
	minTokenDenom := pool.CounterpartyDenom

	standardReserveAmt := balances.AmountOf(standardDenom)
	tokenReserveAmt := balances.AmountOf(minTokenDenom)
	liquidityReserve := k.bk.GetSupply(ctx, lptDenom).Amount
	if standardReserveAmt.LT(msg.MinStandardAmt) {
		return nil, errorsmod.Wrapf(
			types.ErrInsufficientFunds,
			"insufficient %s funds, user expected: %s, actual: %s",
			standardDenom,
			msg.MinStandardAmt.String(),
			standardReserveAmt.String(),
		)
	}
	if tokenReserveAmt.LT(msg.MinToken) {
		return nil, errorsmod.Wrapf(
			types.ErrInsufficientFunds,
			"insufficient %s funds, user expected: %s, actual: %s",
			minTokenDenom,
			msg.MinToken.String(),
			tokenReserveAmt.String(),
		)
	}
	if liquidityReserve.LT(msg.WithdrawLiquidity.Amount) {
		return nil, errorsmod.Wrapf(
			types.ErrInsufficientFunds,
			"insufficient %s funds, user expected: %s, actual: %s",
			lptDenom,
			msg.WithdrawLiquidity.Amount.String(),
			liquidityReserve.String(),
		)
	}

	// calculate amount of UNI to be burned for sender
	// and coin amount to be returned
	irisWithdrawnAmt := msg.WithdrawLiquidity.Amount.Mul(standardReserveAmt).Quo(liquidityReserve)
	tokenWithdrawnAmt := msg.WithdrawLiquidity.Amount.Mul(tokenReserveAmt).Quo(liquidityReserve)

	irisWithdrawCoin := sdk.NewCoin(standardDenom, irisWithdrawnAmt)
	tokenWithdrawCoin := sdk.NewCoin(minTokenDenom, tokenWithdrawnAmt)
	deductUniCoin := msg.WithdrawLiquidity

	if irisWithdrawCoin.Amount.LT(msg.MinStandardAmt) {
		return nil, errorsmod.Wrapf(
			types.ErrConstraintNotMet,
			"iris amount not met, user expected: no less than %s, actual: %s",
			sdk.NewCoin(standardDenom, msg.MinStandardAmt).String(),
			irisWithdrawCoin.String(),
		)
	}
	if tokenWithdrawCoin.Amount.LT(msg.MinToken) {
		return nil, errorsmod.Wrapf(
			types.ErrConstraintNotMet,
			"token amount not met, user expected: no less than %s, actual: %s",
			sdk.NewCoin(minTokenDenom, msg.MinToken).String(),
			tokenWithdrawCoin.String(),
		)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRemoveLiquidity,
			sdk.NewAttribute(types.AttributeValueSender, msg.Sender),
			sdk.NewAttribute(
				types.AttributeValueTokenPair,
				types.GetTokenPairByDenom(minTokenDenom, standardDenom),
			),
		),
	)

	poolAddr, err := sdk.AccAddressFromBech32(pool.EscrowAddress)
	if err != nil {
		return nil, err
	}

	return k.removeLiquidity(
		ctx,
		poolAddr,
		sender,
		deductUniCoin,
		irisWithdrawCoin,
		tokenWithdrawCoin,
	)
}

func (k Keeper) removeLiquidity(
	ctx sdk.Context,
	poolAddr, sender sdk.AccAddress,
	deductUniCoin, irisWithdrawCoin, tokenWithdrawCoin sdk.Coin,
) (sdk.Coins, error) {
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

// RemoveUnilateralLiquidity removes liquidity unilaterally from the specified pool
func (k Keeper) RemoveUnilateralLiquidity(
	ctx sdk.Context,
	msg *types.MsgRemoveUnilateralLiquidity,
) (sdk.Coins, error) {
	var targetTokenDenom string

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdk.Coins{}, err
	}

	poolID := types.GetPoolId(msg.CounterpartyDenom)
	pool, exist := k.GetPool(ctx, poolID)
	if !exist {
		return sdk.Coins{}, errorsmod.Wrapf(
			types.ErrReservePoolNotExists,
			"liquidity pool: %s ",
			poolID,
		)
	}

	poolAddr, err := sdk.AccAddressFromBech32(pool.EscrowAddress)
	if err != nil {
		return sdk.Coins{}, err
	}

	balances, err := k.GetPoolBalances(ctx, pool.EscrowAddress)
	if err != nil {
		return sdk.Coins{}, err
	}

	if msg.MinToken.Denom != msg.CounterpartyDenom &&
		msg.MinToken.Denom != k.GetStandardDenom(ctx) {
		return sdk.Coins{}, errorsmod.Wrapf(types.ErrInvalidDenom,
			"liquidity pool %s has no %s", poolID, msg.MinToken.Denom)
	}

	lptDenom := pool.LptDenom
	targetTokenDenom = msg.MinToken.Denom

	targetBalanceAmt := balances.AmountOf(targetTokenDenom)
	lptBalanceAmt := k.bk.GetSupply(ctx, lptDenom).Amount

	if lptBalanceAmt.LT(msg.ExactLiquidity) {
		return sdk.Coins{}, errorsmod.Wrapf(types.ErrInsufficientFunds,
			"insufficient %s funds, user expected: %s, actual: %s",
			lptDenom, msg.ExactLiquidity.String(), lptBalanceAmt.String())
	}

	if lptBalanceAmt.Equal(msg.ExactLiquidity) {
		return sdk.Coins{}, errorsmod.Wrapf(
			types.ErrConstraintNotMet,
			"forbid to withdraw all liquidity unilaterally, should be less than: %s",
			lptBalanceAmt.String(),
		)
	}

	if targetBalanceAmt.LT(msg.MinToken.Amount) {
		return sdk.Coins{}, errorsmod.Wrapf(types.ErrInsufficientFunds,
			"insufficient %s funds, user expected: %s, actual: %s",
			targetTokenDenom, msg.MinToken.Amount.String(), targetBalanceAmt.String())
	}

	// Calculate Withdrawn Amount
	// t_withdrawn = t_balance * delta_lpt / lpt_balance
	// c_withdrawn = c_balance * delta_lpt / lpt_balance
	//
	// Calculate Swap Amount
	// As `(t_balance - t_withdraw)(c_balance - c_withdraw) = (t_balance - t_withdraw - t_swap) * c_balance`,
	// we get `t_swap = (t_balance - t_withdraw) * c_withdraw / c_balance`
	//
	// Simplify the formula:
	// target_amt = t_balance * (2 * lpt_balance - delta_lpt) * delta_lpt / (lpt_balance^2)
	//
	// Deduce with fee
	// target_amt' = target_amt * ( 1 - fee_unilateral)
	// fee_unilateral = numerator / denominator
	deltaFeeUnilateral := sdk.OneDec().Sub(k.GetParams(ctx).UnilateralLiquidityFee)
	feeNumerator := sdkmath.NewIntFromBigInt(deltaFeeUnilateral.BigInt())
	feeDenominator := sdkmath.NewIntWithDecimal(1, sdk.Precision)

	targetTokenNumerator := lptBalanceAmt.Add(lptBalanceAmt).Sub(msg.ExactLiquidity).
		Mul(msg.ExactLiquidity).Mul(targetBalanceAmt).Mul(feeNumerator)
	targetTokenDenominator := lptBalanceAmt.Mul(lptBalanceAmt).Mul(feeDenominator)

	targetTokenAmtAfterFee := targetTokenNumerator.Quo(targetTokenDenominator)

	if targetTokenAmtAfterFee.LT(msg.MinToken.Amount) {
		return nil, errorsmod.Wrapf(types.ErrConstraintNotMet,
			"token withdrawn amount not met, user expected: no less than %s, actual: %s",
			msg.MinToken.String(), sdk.NewCoin(targetTokenDenom, targetTokenAmtAfterFee).String())
	}

	// event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRemoveUnilateralLiquidity,
			sdk.NewAttribute(types.AttributeValueSender, msg.Sender),
			sdk.NewAttribute(types.AttributeValueTokenUnilateral, targetTokenDenom),
			sdk.NewAttribute(types.AttributeValueLptDenom, pool.LptDenom),
		),
	)

	return k.removeUnilateralLiquidity(
		ctx,
		sender,
		poolAddr,
		lptDenom,
		targetTokenDenom,
		msg.ExactLiquidity,
		targetTokenAmtAfterFee,
	)
}

func (k Keeper) removeUnilateralLiquidity(ctx sdk.Context,
	sender sdk.AccAddress,
	poolAddr sdk.AccAddress,
	lptDenom string,
	targetTokenDenom string,
	exactLiquidity sdkmath.Int,
	targetTokenAmtAfterFee sdkmath.Int,
) (sdk.Coins, error) {
	// send lpt and burn lpt
	lptCoins := sdk.NewCoins(sdk.NewCoin(lptDenom, exactLiquidity))
	if err := k.bk.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, lptCoins); err != nil {
		return nil, err
	}
	if err := k.bk.BurnCoins(ctx, types.ModuleName, lptCoins); err != nil {
		return nil, err
	}

	// send withdraw coins
	coins := sdk.NewCoins(sdk.NewCoin(targetTokenDenom, targetTokenAmtAfterFee))

	return coins, k.bk.SendCoins(ctx, poolAddr, sender, coins)
}
