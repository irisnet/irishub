package tests

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/distribution/keeper"
	"github.com/irisnet/irishub/modules/distribution/types"
	"github.com/irisnet/irishub/modules/stake"
	sdk "github.com/irisnet/irishub/types"
)

const custom = "custom"
const QuerierRoute = "distr"

func getQueriedWithdrawAddr(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, delegatorAddr sdk.AccAddress) (address sdk.AccAddress) {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QuerierRoute, keeper.QueryWithdrawAddr}, "/"),
		Data: cdc.MustMarshalJSON(keeper.NewQueryDelegatorParams(delegatorAddr)),
	}

	bz, err := querier(ctx, []string{keeper.QueryWithdrawAddr}, query)
	require.Nil(t, err)
	require.Nil(t, cdc.UnmarshalJSON(bz, &address))
	return
}

func getQueriedDelegationDistInfo(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress) (ddi types.DelegationDistInfo) {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QuerierRoute, keeper.QueryDelegationDistInfo}, "/"),
		Data: cdc.MustMarshalJSON(keeper.NewQueryDelegationDistInfoParams(delegatorAddr, validatorAddr)),
	}

	bz, err := querier(ctx, []string{keeper.QueryDelegationDistInfo}, query)
	require.Nil(t, err)
	require.Nil(t, cdc.UnmarshalJSON(bz, &ddi))
	return
}

func getQueriedAllDelegationDistInfo(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, delegatorAddr sdk.AccAddress) (ddis []types.DelegationDistInfo) {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QuerierRoute, keeper.QueryAllDelegationDistInfo}, "/"),
		Data: cdc.MustMarshalJSON(keeper.NewQueryDelegatorParams(delegatorAddr)),
	}

	bz, err := querier(ctx, []string{keeper.QueryAllDelegationDistInfo}, query)
	require.Nil(t, err)
	require.Nil(t, cdc.UnmarshalJSON(bz, &ddis))

	return
}

func getQueriedValidatorDistInfo(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, validatorAddr sdk.ValAddress) (vdi types.ValidatorDistInfo) {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QuerierRoute, keeper.QueryValidatorDistInfo}, "/"),
		Data: cdc.MustMarshalJSON(keeper.NewQueryValidatorDistInfoParams(validatorAddr)),
	}

	bz, err := querier(ctx, []string{keeper.QueryValidatorDistInfo}, query)
	require.Nil(t, err)
	require.Nil(t, cdc.UnmarshalJSON(bz, &vdi))

	return
}

func getQueriedRewards(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, address sdk.AccAddress) (rewards keeper.Rewards) {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QuerierRoute, keeper.QueryRewards}, "/"),
		Data: cdc.MustMarshalJSON(keeper.NewQueryRewardsParams(address)),
	}

	bz, err := querier(ctx, []string{keeper.QueryRewards}, query)
	require.Nil(t, err)
	require.Nil(t, cdc.UnmarshalJSON(bz, &rewards))

	return
}

func TestQueries(t *testing.T) {
	cdc := MakeTestCodec()
	initCoins, _ := sdk.NewIntFromString("100000000000000000000")
	ctx, _, dk, sk, feeKeeper := CreateTestInputDefault(t, false, initCoins)

	querier := keeper.NewQuerier(dk)

	// test param queries
	communityTax := sdk.NewDecWithPrec(3, 1)
	baseProposerReward := sdk.NewDecWithPrec(2, 1)
	bonusProposerReward := sdk.NewDecWithPrec(1, 1)
	dk.SetParams(ctx, types.Params{
		CommunityTax:        communityTax,
		BaseProposerReward:  baseProposerReward,
		BonusProposerReward: bonusProposerReward,
	})
	params := dk.GetParams(ctx)
	require.Equal(t, communityTax, params.CommunityTax)
	require.Equal(t, baseProposerReward, params.BaseProposerReward)
	require.Equal(t, bonusProposerReward, params.BonusProposerReward)

	// test withdraw address query
	dk.SetDelegatorWithdrawAddr(ctx, delAddr1, delAddr2)
	withdrawAddr := getQueriedWithdrawAddr(t, ctx, cdc, querier, delAddr1)
	require.Equal(t, delAddr2, withdrawAddr)

	// test delegation distribution info query
	delegation := types.NewDelegationDistInfo(delAddr1, valOpAddr1, ctx.BlockHeight())
	dk.SetDelegationDistInfo(ctx, delegation)
	delegationQuery := getQueriedDelegationDistInfo(t, ctx, cdc, querier, delAddr1, valOpAddr1)
	require.Equal(t, delegation.DelegatorAddr, delegationQuery.DelegatorAddr)
	require.Equal(t, delegation.ValOperatorAddr, delegationQuery.ValOperatorAddr)
	require.Equal(t, delegation.DelPoolWithdrawalHeight, delegationQuery.DelPoolWithdrawalHeight)

	// test all delegation distribution info query
	delegation1 := types.NewDelegationDistInfo(delAddr1, valOpAddr2, ctx.BlockHeight())
	dk.SetDelegationDistInfo(ctx, delegation1)
	delegations := getQueriedAllDelegationDistInfo(t, ctx, cdc, querier, delAddr1)
	require.Equal(t, 2, len(delegations))

	// test validator
	validatorDistInfo := types.NewValidatorDistInfo(valOpAddr3, ctx.BlockHeight())
	dk.SetValidatorDistInfo(ctx, validatorDistInfo)
	validatorDistInfoQuery := getQueriedValidatorDistInfo(t, ctx, cdc, querier, valOpAddr3)
	require.Equal(t, validatorDistInfo.String(), validatorDistInfoQuery.String())

	// test rewards query
	sh := stake.NewHandler(sk)
	comm := stake.NewCommissionMsg(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	msg := stake.NewMsgCreateValidator(valOpAddr1, valConsPk1,
		sdk.NewCoin(stake.BondDenom, initCoins), stake.Description{}, comm)
	require.True(t, sh(ctx, msg).IsOK())
	stake.EndBlocker(ctx, sk)
	rewards := getQueriedRewards(t, ctx, cdc, querier, sdk.AccAddress(valOpAddr1))
	require.True(t, rewards.Total.IsZero())
	initial := int64(20)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	sk.SetLastTotalPower(ctx, sdk.NewInt(100))
	tokens := sdk.Coins{{stake.BondDenom, sdk.NewInt(initial)}}
	feeKeeper.SetCollectedFees(tokens)
	dk.AllocateTokens(ctx, sdk.NewDec(1), valConsAddr1)
	rewards1 := getQueriedRewards(t, ctx, cdc, querier, sdk.AccAddress(valOpAddr1))
	require.Equal(t, true, sdk.Coins{sdk.NewCoin(stake.BondDenom, sdk.NewInt(14))}.IsEqual(rewards1.Total))
	require.Equal(t, true, sdk.Coins{sdk.NewCoin(stake.BondDenom, sdk.NewInt(7))}.IsEqual(rewards1.Commission))
	require.Equal(t, true, sdk.Coins{sdk.NewCoin(stake.BondDenom, sdk.NewInt(7))}.IsEqual(rewards1.Delegations[0].Reward))
}
