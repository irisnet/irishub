package stake

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

type GenesisState struct {
	BondedPool           BondedPool            `json:"pool"`
	Params               Params                `json:"params"`
	LastTotalPower       sdk.Int               `json:"last_total_power"`
	LastValidatorPowers  []LastValidatorPower  `json:"last_validator_powers"`
	Validators           []Validator           `json:"validators"`
	Bonds                []Delegation          `json:"bonds"`
	UnbondingDelegations []UnbondingDelegation `json:"unbonding_delegations"`
	Redelegations        []Redelegation        `json:"redelegations"`
	Exported             bool                  `json:"exported"`
}

type BondedPool struct {
	BondedTokens sdk.Dec `json:"bonded_tokens"` // reserve of bonded tokens
}

type Params struct {
	UnbondingTime time.Duration `json:"unbonding_time"`
	MaxValidators uint16        `json:"max_validators"` // maximum number of validators
}

type LastValidatorPower struct {
	Address sdk.ValAddress `json:"address"`
	Power   sdk.Int        `json:"power"`
}

type Validator struct {
	OperatorAddr     sdk.ValAddress `json:"operator_address"` // address of the validator's operator; bech encoded in JSON
	ConsPubKey       crypto.PubKey  `json:"consensus_pubkey"` // the consensus public key of the validator; bech encoded in JSON
	Jailed           bool           `json:"jailed"`           // has the validator been jailed from bonded status?
	Status           sdk.BondStatus `json:"status"`           // validator status (bonded/unbonding/unbonded)
	Tokens           sdk.Dec        `json:"tokens"`           // delegated tokens (incl. self-delegation)
	DelegatorShares  sdk.Dec        `json:"delegator_shares"` // total shares issued to a validator's delegators
	Description      Description    `json:"description"`      // description terms for the validator
	BondHeight       int64          `json:"bond_height"`      // earliest height as a bonded validator
	UnbondingHeight  int64          `json:"unbonding_height"` // if unbonding, height at which this validator has begun unbonding
	UnbondingMinTime time.Time      `json:"unbonding_time"`   // if unbonding, min time for the validator to complete unbonding
	Commission       Commission     `json:"commission"`       // commission parameters
}

type Description struct {
	Moniker  string `json:"moniker"`  // name
	Identity string `json:"identity"` // optional identity signature (ex. UPort or Keybase)
	Website  string `json:"website"`  // optional website link
	Details  string `json:"details"`  // optional details
}

type Commission struct {
	Rate          sdk.Dec   `json:"rate"`            // the commission rate charged to delegators
	MaxRate       sdk.Dec   `json:"max_rate"`        // maximum commission rate which validator can ever charge
	MaxChangeRate sdk.Dec   `json:"max_change_rate"` // maximum daily increase of the validator commission
	UpdateTime    time.Time `json:"update_time"`     // the last time the commission rate was changed
}

type Delegation struct {
	DelegatorAddr sdk.AccAddress `json:"delegator_addr"`
	ValidatorAddr sdk.ValAddress `json:"validator_addr"`
	Shares        sdk.Dec        `json:"shares"`
	Height        int64          `json:"height"` // Last height bond updated
}

type UnbondingDelegation struct {
	TxHash         string         `json:"tx_hash"`
	DelegatorAddr  sdk.AccAddress `json:"delegator_addr"`  // delegator
	ValidatorAddr  sdk.ValAddress `json:"validator_addr"`  // validator unbonding from operator addr
	CreationHeight int64          `json:"creation_height"` // height which the unbonding took place
	MinTime        time.Time      `json:"min_time"`        // unix time for unbonding completion
	InitialBalance sdk.Coin       `json:"initial_balance"` // atoms initially scheduled to receive at completion
	Balance        sdk.Coin       `json:"balance"`         // atoms to receive at completion
}

type Redelegation struct {
	DelegatorAddr    sdk.AccAddress `json:"delegator_addr"`     // delegator
	ValidatorSrcAddr sdk.ValAddress `json:"validator_src_addr"` // validator redelegation source operator addr
	ValidatorDstAddr sdk.ValAddress `json:"validator_dst_addr"` // validator redelegation destination operator addr
	CreationHeight   int64          `json:"creation_height"`    // height which the redelegation took place
	MinTime          time.Time      `json:"min_time"`           // unix time for redelegation completion
	InitialBalance   sdk.Coin       `json:"initial_balance"`    // initial balance when redelegation started
	Balance          sdk.Coin       `json:"balance"`            // current balance
	SharesSrc        sdk.Dec        `json:"shares_src"`         // amount of source shares redelegating
	SharesDst        sdk.Dec        `json:"shares_dst"`         // amount of destination shares redelegating
}
