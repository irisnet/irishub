package types

import (
	sdk "github.com/irisnet/irishub/types"
)

// GenesisState - all staking state that must be provided at genesis
type GenesisState struct {
	BondedPool           BondedPool            `json:"pool"`
	Params               Params                `json:"params"`
	IntraTxCounter       int16                 `json:"intra_tx_counter"`
	LastTotalPower       sdk.Int               `json:"last_total_power"`
	LastValidatorPowers  []LastValidatorPower  `json:"last_validator_powers"`
	Validators           []Validator           `json:"validators"`
	Bonds                []Delegation          `json:"bonds"`
	UnbondingDelegations []UnbondingDelegation `json:"unbonding_delegations"`
	Redelegations        []Redelegation        `json:"redelegations"`
	Exported             bool                  `json:"exported"`
}

// Last validator power, needed for validator set update logic
type LastValidatorPower struct {
	Address sdk.ValAddress
	Power   sdk.Int
}

func NewGenesisState(bondedPool BondedPool, params Params, validators []Validator, bonds []Delegation) GenesisState {
	return GenesisState{
		BondedPool: bondedPool,
		Params:     params,
		Validators: validators,
		Bonds:      bonds,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		BondedPool: InitialBondedPool(),
		Params:     DefaultParams(),
	}
}
