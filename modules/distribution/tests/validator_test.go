package tests

import (
	"testing"

	"github.com/irisnet/irishub/modules/stake"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
)

func TestWithdrawValidatorRewardsAllNoDelegator(t *testing.T) {
	ctx, accMapper, keeper, sk, fck := CreateTestInputAdvanced(t, false, sdk.NewIntWithDecimal(100, 18), sdk.ZeroDec())
	stakeHandler := stake.NewHandler(sk)
	denom := sk.BondDenom()

	//first make a validator
	msgCreateValidator := stake.NewTestMsgCreateValidator(valOpAddr1, valConsPk1, sdk.NewIntWithDecimal(10, 18))
	got := stakeHandler(ctx, msgCreateValidator)
	require.True(t, got.IsOK(), "expected msg to be ok, got %v", got)
	_ = sk.ApplyAndReturnValidatorSetUpdates(ctx)

	// allocate 100 denom of fees
	feeInputs := sdk.NewIntWithDecimal(100, 18)
	fck.SetCollectedFees(sdk.Coins{sdk.NewCoin(denom, feeInputs)})
	require.Equal(t, feeInputs, fck.GetCollectedFees(ctx).AmountOf(denom))
	keeper.AllocateTokens(ctx, sdk.OneDec(), valConsAddr1)

	// withdraw self-delegation reward
	ctx = ctx.WithBlockHeight(1)
	keeper.WithdrawValidatorRewardsAll(ctx, valOpAddr1)
	amt := accMapper.GetAccount(ctx, valAccAddr1).GetCoins().AmountOf(denom)
	expRes := sdk.NewDecFromInt(sdk.NewIntWithDecimal(90, 18)).Add(sdk.NewDecFromInt(sdk.NewIntWithDecimal(100, 18))).TruncateInt()
	require.True(sdk.IntEq(t, expRes, amt))
}

func TestWithdrawValidatorRewardsAllDelegatorNoCommission(t *testing.T) {
	ctx, accMapper, keeper, sk, fck := CreateTestInputAdvanced(t, false, sdk.NewIntWithDecimal(100, 18), sdk.ZeroDec())
	stakeHandler := stake.NewHandler(sk)
	denom := sk.BondDenom()

	//first make a validator
	msgCreateValidator := stake.NewTestMsgCreateValidator(valOpAddr1, valConsPk1, sdk.NewIntWithDecimal(10, 18))
	got := stakeHandler(ctx, msgCreateValidator)
	require.True(t, got.IsOK(), "expected msg to be ok, got %v", got)
	_ = sk.ApplyAndReturnValidatorSetUpdates(ctx)

	// delegate
	msgDelegate := stake.NewTestMsgDelegate(delAddr1, valOpAddr1, sdk.NewIntWithDecimal(10, 18))
	got = stakeHandler(ctx, msgDelegate)
	require.True(t, got.IsOK())
	amt := accMapper.GetAccount(ctx, delAddr1).GetCoins().AmountOf(denom)
	require.Equal(t, sdk.NewIntWithDecimal(90, 18), amt)

	// allocate 100 denom of fees
	feeInputs := sdk.NewIntWithDecimal(100, 18)
	fck.SetCollectedFees(sdk.Coins{sdk.NewCoin(denom, feeInputs)})
	require.Equal(t, feeInputs, fck.GetCollectedFees(ctx).AmountOf(denom))
	keeper.AllocateTokens(ctx, sdk.OneDec(), valConsAddr1)

	// withdraw self-delegation reward
	ctx = ctx.WithBlockHeight(1)
	keeper.WithdrawValidatorRewardsAll(ctx, valOpAddr1)
	amt = accMapper.GetAccount(ctx, valAccAddr1).GetCoins().AmountOf(denom)
	expRes := sdk.NewDecFromInt(sdk.NewIntWithDecimal(90, 18)).Add(sdk.NewDecFromInt(sdk.NewIntWithDecimal(100, 18)).Quo(sdk.NewDec(2))).TruncateInt() // 90 + 100 tokens * 10/20
	require.True(sdk.IntEq(t, expRes, amt))
}

func TestWithdrawValidatorRewardsAllDelegatorWithCommission(t *testing.T) {
	ctx, accMapper, keeper, sk, fck := CreateTestInputAdvanced(t, false, sdk.NewIntWithDecimal(100, 18), sdk.ZeroDec())
	stakeHandler := stake.NewHandler(sk)
	denom := sk.BondDenom()

	//first make a validator
	commissionRate := sdk.NewDecWithPrec(1, 1)
	msgCreateValidator := stake.NewTestMsgCreateValidatorWithCommission(
		valOpAddr1, valConsPk1, sdk.NewIntWithDecimal(10, 18), commissionRate)
	got := stakeHandler(ctx, msgCreateValidator)
	require.True(t, got.IsOK(), "expected msg to be ok, got %v", got)
	_ = sk.ApplyAndReturnValidatorSetUpdates(ctx)

	// delegate
	msgDelegate := stake.NewTestMsgDelegate(delAddr1, valOpAddr1, sdk.NewIntWithDecimal(10, 18))
	got = stakeHandler(ctx, msgDelegate)
	require.True(t, got.IsOK())
	amt := accMapper.GetAccount(ctx, delAddr1).GetCoins().AmountOf(denom)
	require.Equal(t, sdk.NewIntWithDecimal(90, 18), amt)

	// allocate 100 denom of fees
	feeInputs := sdk.NewIntWithDecimal(100, 18)
	fck.SetCollectedFees(sdk.Coins{sdk.NewCoin(denom, feeInputs)})
	require.Equal(t, feeInputs, fck.GetCollectedFees(ctx).AmountOf(denom))
	keeper.AllocateTokens(ctx, sdk.OneDec(), valConsAddr1)

	// withdraw validator reward
	ctx = ctx.WithBlockHeight(1)
	keeper.WithdrawValidatorRewardsAll(ctx, valOpAddr1)
	amt = accMapper.GetAccount(ctx, valAccAddr1).GetCoins().AmountOf(denom)
	commissionTaken := sdk.NewDecFromInt(sdk.NewIntWithDecimal(100, 18)).Mul(commissionRate)
	afterCommission := sdk.NewDecFromInt(sdk.NewIntWithDecimal(100, 18)).Sub(commissionTaken)
	selfDelegationReward := afterCommission.Quo(sdk.NewDec(2))
	expRes := sdk.NewDecFromInt(sdk.NewIntWithDecimal(90, 18)).Add(commissionTaken).Add(selfDelegationReward).TruncateInt() // 90 + 100 tokens * 10/20
	require.True(sdk.IntEq(t, expRes, amt))
}

func TestWithdrawValidatorRewardsAllMultipleValidator(t *testing.T) {
	ctx, accMapper, keeper, sk, fck := CreateTestInputAdvanced(t, false, sdk.NewIntWithDecimal(100, 18), sdk.ZeroDec())
	stakeHandler := stake.NewHandler(sk)
	denom := sk.BondDenom()

	// Make some  validators with different commissions.
	// Bond 10 of 100 with 0.1 commission.
	msgCreateValidator := stake.NewTestMsgCreateValidatorWithCommission(
		valOpAddr1, valConsPk1, sdk.NewIntWithDecimal(10, 18), sdk.NewDecWithPrec(1, 1))
	got := stakeHandler(ctx, msgCreateValidator)
	require.True(t, got.IsOK(), "expected msg to be ok, got %v", got)

	// Bond 50 of 100 with 0.2 commission.
	msgCreateValidator = stake.NewTestMsgCreateValidatorWithCommission(
		valOpAddr2, valConsPk2, sdk.NewIntWithDecimal(50, 18), sdk.NewDecWithPrec(2, 1))
	got = stakeHandler(ctx, msgCreateValidator)
	require.True(t, got.IsOK(), "expected msg to be ok, got %v", got)

	// Bond 40 of 100 with 0.3 commission.
	msgCreateValidator = stake.NewTestMsgCreateValidatorWithCommission(
		valOpAddr3, valConsPk3, sdk.NewIntWithDecimal(40, 18), sdk.NewDecWithPrec(3, 1))
	got = stakeHandler(ctx, msgCreateValidator)
	require.True(t, got.IsOK(), "expected msg to be ok, got %v", got)

	_ = sk.ApplyAndReturnValidatorSetUpdates(ctx)

	// allocate 1000 denom of fees
	feeInputs := sdk.NewIntWithDecimal(1000, 18)
	fck.SetCollectedFees(sdk.Coins{sdk.NewCoin(denom, feeInputs)})
	require.Equal(t, feeInputs, fck.GetCollectedFees(ctx).AmountOf(denom))
	keeper.AllocateTokens(ctx, sdk.OneDec(), valConsAddr1)

	// withdraw validator reward
	ctx = ctx.WithBlockHeight(1)
	keeper.WithdrawValidatorRewardsAll(ctx, valOpAddr1)
	amt := accMapper.GetAccount(ctx, valAccAddr1).GetCoins().AmountOf(denom)

	feesInNonProposer := sdk.NewDecFromInt(feeInputs).Mul(sdk.NewDecWithPrec(95, 2))
	feesInProposer := sdk.NewDecFromInt(feeInputs).Mul(sdk.NewDecWithPrec(5, 2))
	expRes := sdk.NewDecFromInt(sdk.NewIntWithDecimal(90, 18)). // orig tokens (100 - 10)
									Add(feesInNonProposer.Quo(sdk.NewDec(10))). // validator 1 has 1/10 total power (non-proposer rewards = 95)
									Add(feesInProposer).                        // (proposer rewards = 50)
									TruncateInt()
	require.True(sdk.IntEq(t, expRes, amt))
}

func TestWithdrawValidatorRewardsAllMultipleDelegator(t *testing.T) {
	ctx, accMapper, keeper, sk, fck := CreateTestInputAdvanced(t, false, sdk.NewIntWithDecimal(100, 18), sdk.ZeroDec())
	stakeHandler := stake.NewHandler(sk)
	denom := sk.BondDenom()

	//first make a validator with 10% commission
	commissionRate := sdk.NewDecWithPrec(1, 1)
	msgCreateValidator := stake.NewTestMsgCreateValidatorWithCommission(
		valOpAddr1, valConsPk1, sdk.NewIntWithDecimal(10, 18), sdk.NewDecWithPrec(1, 1))
	got := stakeHandler(ctx, msgCreateValidator)
	require.True(t, got.IsOK(), "expected msg to be ok, got %v", got)
	_ = sk.ApplyAndReturnValidatorSetUpdates(ctx)

	// delegate
	msgDelegate := stake.NewTestMsgDelegate(delAddr1, valOpAddr1, sdk.NewIntWithDecimal(10, 18))
	got = stakeHandler(ctx, msgDelegate)
	require.True(t, got.IsOK())
	amt := accMapper.GetAccount(ctx, delAddr1).GetCoins().AmountOf(denom)
	require.Equal(t, sdk.NewIntWithDecimal(90, 18), amt)

	msgDelegate = stake.NewTestMsgDelegate(delAddr2, valOpAddr1, sdk.NewIntWithDecimal(20, 18))
	got = stakeHandler(ctx, msgDelegate)
	require.True(t, got.IsOK())
	amt = accMapper.GetAccount(ctx, delAddr2).GetCoins().AmountOf(denom)
	require.Equal(t, sdk.NewIntWithDecimal(80, 18), amt)

	// allocate 100 denom of fees
	feeInputs := sdk.NewIntWithDecimal(100, 18)
	fck.SetCollectedFees(sdk.Coins{sdk.NewCoin(denom, feeInputs)})
	require.Equal(t, feeInputs, fck.GetCollectedFees(ctx).AmountOf(denom))
	keeper.AllocateTokens(ctx, sdk.OneDec(), valConsAddr1)

	// withdraw validator reward
	ctx = ctx.WithBlockHeight(1)
	keeper.WithdrawValidatorRewardsAll(ctx, valOpAddr1)
	amt = accMapper.GetAccount(ctx, valAccAddr1).GetCoins().AmountOf(denom)

	commissionTaken := sdk.NewDecFromInt(sdk.NewIntWithDecimal(100, 18)).Mul(commissionRate)
	afterCommission := sdk.NewDecFromInt(sdk.NewIntWithDecimal(100, 18)).Sub(commissionTaken)
	expRes := sdk.NewDecFromInt(sdk.NewIntWithDecimal(90, 18)).
		Add(afterCommission.Quo(sdk.NewDec(4))).
		Add(commissionTaken).
		TruncateInt() // 90 + 100*90% tokens * 10/40
	require.True(sdk.IntEq(t, expRes, amt))
}
