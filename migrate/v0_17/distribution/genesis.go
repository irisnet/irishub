package distr

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	Params                 Params                  `json:"params"`
	FeePool                FeePool                 `json:"fee_pool"`
	ValidatorDistInfos     []ValidatorDistInfo     `json:"validator_dist_infos"`
	DelegationDistInfos    []DelegationDistInfo    `json:"delegator_dist_infos"`
	DelegatorWithdrawInfos []DelegatorWithdrawInfo `json:"delegator_withdraw_infos"`
	PreviousProposer       sdk.ConsAddress         `json:"previous_proposer"`
}

type Params struct {
	CommunityTax        sdk.Dec `json:"community_tax"`
	BaseProposerReward  sdk.Dec `json:"base_proposer_reward"`
	BonusProposerReward sdk.Dec `json:"bonus_proposer_reward"`
}

type FeePool struct {
	TotalValAccum TotalAccum   `json:"val_accum"`      // total valdator accum held by validators
	ValPool       sdk.DecCoins `json:"val_pool"`       // funds for all validators which have yet to be withdrawn
	CommunityPool sdk.DecCoins `json:"community_pool"` // pool for community funds yet to be spent
}

type ValidatorDistInfo struct {
	OperatorAddr            sdk.ValAddress `json:"operator_addr"`
	FeePoolWithdrawalHeight int64          `json:"fee_pool_withdrawal_height"` // last height this validator withdrew from the global pool
	DelAccum                TotalAccum     `json:"del_accum"`                  // total accumulation factor held by delegators
	DelPool                 sdk.DecCoins   `json:"del_pool"`                   // rewards owed to delegators, commission has already been charged (includes proposer reward)
	ValCommission           sdk.DecCoins   `json:"val_commission"`             // commission collected by this validator (pending withdrawal)
}

type TotalAccum struct {
	UpdateHeight int64   `json:"update_height"`
	Accum        sdk.Dec `json:"accum"`
}

type DelegationDistInfo struct {
	DelegatorAddr           sdk.AccAddress `json:"delegator_addr"`
	ValOperatorAddr         sdk.ValAddress `json:"val_operator_addr"`
	DelPoolWithdrawalHeight int64          `json:"del_pool_withdrawal_height"` // last time this delegation withdrew rewards
}

type DelegatorWithdrawInfo struct {
	DelegatorAddr sdk.AccAddress `json:"delegator_addr"`
	WithdrawAddr  sdk.AccAddress `json:"withdraw_addr"`
}
