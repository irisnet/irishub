package types

import (
	sdk "github.com/irisnet/irishub/types"
)

// the address for where distributions rewards are withdrawn to by default
// this struct is only used at genesis to feed in default withdraw addresses
type DelegatorWithdrawInfo struct {
	DelegatorAddr sdk.AccAddress `json:"delegator_addr"`
	WithdrawAddr  sdk.AccAddress `json:"withdraw_addr"`
}

// GenesisState - all distribution state that must be provided at genesis
type GenesisState struct {
	Params                 Params                  `json:"params"`
	FeePool                FeePool                 `json:"fee_pool"`
	ValidatorDistInfos     []ValidatorDistInfo     `json:"validator_dist_infos"`
	DelegationDistInfos    []DelegationDistInfo    `json:"delegator_dist_infos"`
	DelegatorWithdrawInfos []DelegatorWithdrawInfo `json:"delegator_withdraw_infos"`
	PreviousProposer       sdk.ConsAddress         `json:"previous_proposer"`
}

func NewGenesisState(params Params, feePool FeePool, vdis []ValidatorDistInfo,
	ddis []DelegationDistInfo, dwis []DelegatorWithdrawInfo, pp sdk.ConsAddress) GenesisState {

	return GenesisState{
		Params:                 params,
		FeePool:                feePool,
		ValidatorDistInfos:     vdis,
		DelegationDistInfos:    ddis,
		DelegatorWithdrawInfos: dwis,
		PreviousProposer:       pp,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:  DefaultParams(),
		FeePool: InitialFeePool(),
	}
}

// default genesis utility function, initialize for starting validator set
func DefaultGenesisWithValidators(valAddrs []sdk.ValAddress) GenesisState {

	vdis := make([]ValidatorDistInfo, len(valAddrs))
	ddis := make([]DelegationDistInfo, len(valAddrs))

	for i, valAddr := range valAddrs {
		vdis[i] = NewValidatorDistInfo(valAddr, 0)
		accAddr := sdk.AccAddress(valAddr)
		ddis[i] = NewDelegationDistInfo(accAddr, valAddr, 0)
	}

	return GenesisState{
		Params:              DefaultParams(),
		FeePool:             InitialFeePool(),
		ValidatorDistInfos:  vdis,
		DelegationDistInfos: ddis,
	}
}
