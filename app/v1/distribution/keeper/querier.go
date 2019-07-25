package keeper

import (
	"fmt"

	"github.com/irisnet/irishub/app/v1/distribution/types"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)

// nolint
const (
	QueryWithdrawAddr          = "withdraw_addr"
	QueryDelegationDistInfo    = "delegation_dist_info"
	QueryAllDelegationDistInfo = "all_delegation_dist_info"
	QueryValidatorDistInfo     = "validator_dist_info"
	QueryRewards               = "rewards"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryWithdrawAddr:
			return queryDelegatorWithdrawAddress(ctx, path[1:], req, k)

		case QueryDelegationDistInfo:
			return queryDelegationDistInfo(ctx, path[1:], req, k)

		case QueryAllDelegationDistInfo:
			return queryAllDelegationDistInfo(ctx, path[1:], req, k)

		case QueryValidatorDistInfo:
			return queryValidatorDistInfo(ctx, path[1:], req, k)

		case QueryRewards:
			return queryRewards(ctx, path[1:], req, k)

		default:
			return nil, sdk.ErrUnknownRequest("unknown distr query endpoint")
		}
	}
}

// params for query 'custom/distr/delegation_dist_info', 'custom/distr/all_delegation_dist_info' and 'withdraw_addr'
type QueryDelegatorParams struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address"`
}

func NewQueryDelegatorParams(delegatorAddr sdk.AccAddress) QueryDelegatorParams {
	return QueryDelegatorParams{DelegatorAddress: delegatorAddr}
}

func queryDelegatorWithdrawAddress(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params QueryDelegatorParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	// cache-wrap context as to not persist state changes during querying
	ctx, _ = ctx.CacheContext()
	withdrawAddr := k.GetDelegatorWithdrawAddr(ctx, params.DelegatorAddress)

	bz, err := codec.MarshalJSONIndent(k.cdc, withdrawAddr)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}

// params for query 'custom/distr/delegation_dist_info'
type QueryDelegationDistInfoParams struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address"`
	ValidatorAddress sdk.ValAddress `json:"validator_address"`
}

// creates a new instance of QueryDelegationDistInfoParams
func NewQueryDelegationDistInfoParams(delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress) QueryDelegationDistInfoParams {
	return QueryDelegationDistInfoParams{
		DelegatorAddress: delegatorAddr,
		ValidatorAddress: validatorAddr,
	}
}

func queryDelegationDistInfo(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params QueryDelegationDistInfoParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	// cache-wrap context as to not persist state changes during querying
	ctx, _ = ctx.CacheContext()
	if !k.HasDelegationDistInfo(ctx, params.DelegatorAddress, params.ValidatorAddress) {
		return []byte{}, types.ErrNoDelegationDistInfo(types.DefaultCodespace)
	}
	ddi := k.GetDelegationDistInfo(ctx, params.DelegatorAddress, params.ValidatorAddress)
	res, errRes := codec.MarshalJSONIndent(k.cdc, ddi)
	if errRes != nil {
		return nil, sdk.MarshalResultErr(err)
	}
	return res, nil
}

func queryAllDelegationDistInfo(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params QueryDelegatorParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	// cache-wrap context as to not persist state changes during querying
	ctx, _ = ctx.CacheContext()
	var distInfos []types.DelegationDistInfo
	ddiIter := func(_ int64, distInfo types.DelegationDistInfo) (stop bool) {
		distInfos = append(distInfos, distInfo)
		if err != nil {
			panic(err)
		}
		return false
	}
	k.IterateDelegatorDistInfos(ctx, params.DelegatorAddress, ddiIter)
	res, errRes := codec.MarshalJSONIndent(k.cdc, distInfos)
	if errRes != nil {
		return nil, sdk.MarshalResultErr(err)
	}
	return res, nil
}

// params for query 'custom/distr/validator_dist_info'
type QueryValidatorDistInfoParams struct {
	ValidatorAddress sdk.ValAddress `json:"validator_address"`
}

// creates a new instance of QueryValidatorDistInfoParams
func NewQueryValidatorDistInfoParams(validatorAddr sdk.ValAddress) QueryValidatorDistInfoParams {
	return QueryValidatorDistInfoParams{
		ValidatorAddress: validatorAddr,
	}
}

func queryValidatorDistInfo(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params QueryValidatorDistInfoParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	// cache-wrap context as to not persist state changes during querying
	ctx, _ = ctx.CacheContext()
	if !k.HasValidatorDistInfo(ctx, params.ValidatorAddress) {
		return []byte{}, types.ErrNoValidatorDistInfo(types.DefaultCodespace)
	}
	vdi := k.GetValidatorDistInfo(ctx, params.ValidatorAddress)
	res, errRes := codec.MarshalJSONIndent(k.cdc, vdi)
	if errRes != nil {
		return nil, sdk.MarshalResultErr(err)
	}
	return res, nil
}

// params for query 'custom/distr/rewards'
type QueryRewardsParams struct {
	Address sdk.AccAddress `json:"address"`
}

// creates a new instance of QueryRewardsParams
func NewQueryRewardsParams(address sdk.AccAddress) QueryRewardsParams {
	return QueryRewardsParams{
		Address: address,
	}
}

type Rewards struct {
	Total       sdk.Coins           `json:"total"`
	Delegations []DelegationsReward `json:"delegations"`
	Commission  sdk.Coins           `json:"commission"`
}

func (r Rewards) String() string {
	var delegations string
	for _, val := range r.Delegations {
		delegations += "\n  " + val.String()
	}
	return fmt.Sprintf(`Total:        %s
Delegations:  %s
Commission:   %s`, r.Total.MainUnitString(), delegations, r.Commission.MainUnitString())
}

func (r Rewards) HumanString(converter sdk.CoinsConverter) string {
	var delegations string
	for _, val := range r.Delegations {
		delegations += "\n  " + val.HumanString(converter)
	}
	return fmt.Sprintf(`Total:        %s
Delegations:  %s
Commission:   %s`, converter.ToMainUnit(r.Total), delegations, converter.ToMainUnit(r.Commission))
}

type DelegationsReward struct {
	Validator sdk.ValAddress `json:"validator"`
	Reward    sdk.Coins      `json:"reward"`
}

func (dr DelegationsReward) String() string {
	return fmt.Sprintf(`validator: %s, reward: %s`,
		dr.Validator, dr.Reward.MainUnitString())
}

func (dr DelegationsReward) HumanString(converter sdk.CoinsConverter) string {
	return fmt.Sprintf(`validator: %s, reward: %s`,
		dr.Validator, converter.ToMainUnit(dr.Reward))
}

func queryRewards(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	totalWithdraw := types.DecCoins{}
	rewards := Rewards{}

	var params QueryRewardsParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	// cache-wrap context as to not persist state changes during querying
	ctx, _ = ctx.CacheContext()
	var selfVdi types.ValidatorDistInfo
	selfValidator := k.stakeKeeper.Validator(ctx, sdk.ValAddress(params.Address))
	if selfValidator != nil && selfValidator.GetOperator().Equals(sdk.ValAddress(params.Address)) {
		selfVdi = k.GetValidatorDistInfo(ctx, selfValidator.GetOperator())
	}

	feePool := k.GetFeePool(ctx)

	// get all delegator rewards
	operationAtDelegation := func(_ int64, del sdk.Delegation) (stop bool) {
		validator := k.stakeKeeper.Validator(ctx, del.GetValidatorAddr())
		vdi := k.GetValidatorDistInfo(ctx, del.GetValidatorAddr())
		wc := k.GetWithdrawContext(ctx, del.GetValidatorAddr())
		wc.FeePool = feePool
		distInfo := k.GetDelegationDistInfo(ctx, del.GetDelegatorAddr(), del.GetValidatorAddr())
		_, vdi, newFeePool, diWithdraw := distInfo.WithdrawRewards(log.NewNopLogger(), wc, vdi, validator.GetDelegatorShares(), del.GetShares())
		totalWithdraw = totalWithdraw.Plus(diWithdraw)
		rewardTruncate, _ := diWithdraw.TruncateDecimal()
		rewards.Delegations = append(rewards.Delegations, DelegationsReward{del.GetValidatorAddr(), rewardTruncate})
		if selfValidator != nil && selfValidator.GetOperator().Equals(vdi.OperatorAddr) {
			selfVdi = vdi
		}
		feePool = newFeePool
		return false
	}
	k.stakeKeeper.IterateDelegations(ctx, params.Address, operationAtDelegation)

	// get all validator rewards
	if !selfVdi.OperatorAddr.Empty() && selfVdi.OperatorAddr.Equals(sdk.ValAddress(params.Address)) {
		wc := k.GetWithdrawContext(ctx, selfValidator.GetOperator())
		wc.FeePool = feePool
		_, _, commission := selfVdi.WithdrawCommission(log.NewNopLogger(), wc)
		totalWithdraw = totalWithdraw.Plus(commission)
		rewardTruncate, _ := commission.TruncateDecimal()
		rewards.Commission = rewardTruncate
	}

	rewardTruncate, _ := totalWithdraw.TruncateDecimal()
	rewards.Total = rewardTruncate
	res, errRes := codec.MarshalJSONIndent(k.cdc, rewards)
	if errRes != nil {
		return nil, sdk.MarshalResultErr(err)
	}
	return res, nil
}

type CommunityTax struct {
	Amount types.DecCoins `json:"amount"`
}

func (ct CommunityTax) String() string {
	return fmt.Sprintf(`Amount:  %s`,
		ct.Amount.MainUnitString())
}
