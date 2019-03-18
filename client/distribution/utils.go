package distribution

import (
	"github.com/irisnet/irishub/client/context"
	tendermint "github.com/irisnet/irishub/client/tendermint/rpc"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/modules/distribution"
	distrKeeper "github.com/irisnet/irishub/modules/distribution/keeper"
	"github.com/irisnet/irishub/modules/distribution/types"
	"github.com/irisnet/irishub/modules/stake"
	stakeKeeper "github.com/irisnet/irishub/modules/stake/keeper"
	stakeTypes "github.com/irisnet/irishub/modules/stake/types"
	sdk "github.com/irisnet/irishub/types"
	"github.com/tendermint/tendermint/libs/log"
)

// distribution info for a particular validator
type ValidatorDistInfoOutput struct {
	OperatorAddr            sdk.ValAddress          `json:"operator_addr"`
	FeePoolWithdrawalHeight int64                   `json:"fee_pool_withdrawal_height"`
	DelAccum                distribution.TotalAccum `json:"del_accum"`
	DelPool                 string                  `json:"del_pool"`
	ValCommission           string                  `json:"val_commission"`
}

type RewardsOutput struct {
	Total       sdk.Coins           `json:"total"`
	Delegations []DelegationsReward `json:"delegations"`
	Commission  sdk.Coins           `json:"commission"`
}

type DelegationsReward struct {
	Validator sdk.ValAddress `json:"validator"`
	Reward    sdk.Coins      `json:"reward"`
}

func ConvertToValidatorDistInfoOutput(cliCtx context.CLIContext, vdi distribution.ValidatorDistInfo) ValidatorDistInfoOutput {
	exRate := utils.ExRateFromStakeTokenToMainUnit(cliCtx)
	delPool := utils.ConvertDecToRat(vdi.DelPool.AmountOf(stakeTypes.StakeDenom)).Mul(exRate).FloatString() + stakeTypes.StakeTokenName
	valCommission := utils.ConvertDecToRat(vdi.ValCommission.AmountOf(stakeTypes.StakeDenom)).Mul(exRate).FloatString() + stakeTypes.StakeTokenName
	return ValidatorDistInfoOutput{
		OperatorAddr:            vdi.OperatorAddr,
		FeePoolWithdrawalHeight: vdi.FeePoolWithdrawalHeight,
		DelAccum:                vdi.DelAccum,
		DelPool:                 delPool,
		ValCommission:           valCommission,
	}
}

func GetRewards(distrStoreName string, stakeStoreName string, cliCtx context.CLIContext, account sdk.AccAddress) RewardsOutput {
	totalWithdraw := types.DecCoins{}
	rewardsOutput := RewardsOutput{}

	// get all delegator rewards
	res, err := cliCtx.QuerySubspace(stakeKeeper.GetDelegationsKey(account), stakeStoreName)
	if err != nil {
		return RewardsOutput{}
	}

	feePool := GetFeePool(distrStoreName, cliCtx)
	chainHeight, err := tendermint.GetChainHeight(cliCtx)
	for _, re := range res {
		del := stakeTypes.MustUnmarshalDelegation(cliCtx.Codec, re.Key, re.Value)
		valAddr := del.GetValidatorAddr()
		validator := GetValidator(stakeStoreName, cliCtx, valAddr)
		if !validator.OperatorAddr.Equals(valAddr) {
			continue
		}
		vdi := GetValidatorDistInfo(distrStoreName, cliCtx, valAddr)
		ddi := GetDelegationDistInfo(distrStoreName, cliCtx, del.DelegatorAddr, del.ValidatorAddr)
		wc := GetWithdrawContext(stakeStoreName, cliCtx, feePool, chainHeight, validator)
		_, _, _, diWithdraw := ddi.WithdrawRewards(log.NewNopLogger(), wc, vdi, validator.GetDelegatorShares(), del.GetShares())
		totalWithdraw = totalWithdraw.Plus(diWithdraw)
		rewardTruncate, _ := diWithdraw.TruncateDecimal()
		rewardsOutput.Delegations = append(rewardsOutput.Delegations, DelegationsReward{valAddr, rewardTruncate})
	}

	// get all validator rewards
	validator := GetValidator(stakeStoreName, cliCtx, sdk.ValAddress(account))
	if validator.OperatorAddr.Equals(sdk.ValAddress(account)) {
		wc := GetWithdrawContext(stakeStoreName, cliCtx, feePool, chainHeight, validator)
		valInfo := GetValidatorDistInfo(distrStoreName, cliCtx, validator.GetOperator())
		valInfo, _, commission := valInfo.WithdrawCommission(log.NewNopLogger(), wc)
		totalWithdraw = totalWithdraw.Plus(commission)
		rewardTruncate, _ := commission.TruncateDecimal()
		rewardsOutput.Commission = rewardTruncate
	}

	rewardTruncate, _ := totalWithdraw.TruncateDecimal()
	rewardsOutput.Total = rewardTruncate
	return rewardsOutput
}

func GetFeePool(storeName string, cliCtx context.CLIContext) (feePool types.FeePool) {
	res, err := cliCtx.QueryStore(distribution.FeePoolKey, storeName)
	if err != nil {
		return
	}
	if res == nil {
		panic("Stored fee pool should not have been nil")
	}
	cliCtx.Codec.MustUnmarshalBinaryLengthPrefixed(res, &feePool)
	return
}

func GetValidatorDistInfo(storeName string, cliCtx context.CLIContext,
	operatorAddr sdk.ValAddress) (vdi types.ValidatorDistInfo) {
	res, err := cliCtx.QueryStore(distrKeeper.GetValidatorDistInfoKey(operatorAddr), storeName)
	if err != nil || len(res) == 0 {
		return
	}
	cliCtx.Codec.MustUnmarshalBinaryLengthPrefixed(res, &vdi)
	return
}

func GetDelegationDistInfo(storeName string, cliCtx context.CLIContext, delAddr sdk.AccAddress,
	valOperatorAddr sdk.ValAddress) (ddi types.DelegationDistInfo) {
	res, err := cliCtx.QueryStore(distrKeeper.GetDelegationDistInfoKey(delAddr, valOperatorAddr), storeName)
	if err != nil || len(res) == 0 {
		return
	}
	cliCtx.Codec.MustUnmarshalBinaryLengthPrefixed(res, &ddi)
	return
}

func GetValidator(storeName string, cliCtx context.CLIContext, addr sdk.ValAddress) (validator stake.Validator) {
	key := stake.GetValidatorKey(addr)
	res, err := cliCtx.QueryStore(key, storeName)
	if err != nil || len(res) == 0 {
		return
	}

	validator = stakeTypes.MustUnmarshalValidator(cliCtx.Codec, addr, res)
	return
}

func GetWithdrawContext(storeName string, cliCtx context.CLIContext, feePool types.FeePool, height int64, validator stakeTypes.Validator) types.WithdrawContext {
	key := stakeKeeper.GetLastValidatorPowerKey(validator.OperatorAddr)
	res, err := cliCtx.QueryStore(key, storeName)
	lastValPower := sdk.ZeroInt()
	if err == nil && len(res) != 0 {
		cliCtx.Codec.MustUnmarshalBinaryLengthPrefixed(res, &lastValPower)
	}

	key1 := stakeKeeper.LastTotalPowerKey
	res, err = cliCtx.QueryStore(key1, storeName)
	lastTotalPower := sdk.ZeroInt()
	if err == nil && len(res) != 0 {
		cliCtx.Codec.MustUnmarshalBinaryLengthPrefixed(res, &lastTotalPower)
	}

	return types.NewWithdrawContext(
		feePool, height, sdk.NewDecFromInt(lastTotalPower), sdk.NewDecFromInt(lastValPower),
		validator.GetCommission())
}
