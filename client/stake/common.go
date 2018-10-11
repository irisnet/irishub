package stake

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/client/context"
	"time"
)

// defines a delegation without type Rat for shares
type DelegationOutput struct {
	DelegatorAddr sdk.AccAddress `json:"delegator_addr"`
	ValidatorAddr sdk.AccAddress `json:"validator_addr"`
	Shares        string         `json:"shares"`
	Height        int64          `json:"height"`
}

func (d DelegationOutput) HumanReadableString() (string, error) {
	resp := "Delegation \n"
	resp += fmt.Sprintf("Delegator: %s\n", d.DelegatorAddr)
	resp += fmt.Sprintf("Validator: %s\n", d.ValidatorAddr)
	resp += fmt.Sprintf("Shares: %s", d.Shares)
	resp += fmt.Sprintf("Height: %d", d.Height)

	return resp, nil
}

// UnbondingDelegation reflects a delegation's passive unbonding queue.
type UnbondingDelegationOutput struct {
	DelegatorAddr  sdk.AccAddress `json:"delegator_addr"`  // delegator
	ValidatorAddr  sdk.AccAddress `json:"validator_addr"`  // validator unbonding from owner addr
	CreationHeight int64          `json:"creation_height"` // height which the unbonding took place
	MinTime        time.Time      `json:"min_time"`        // unix time for unbonding completion
	InitialBalance string         `json:"initial_balance"` // atoms initially scheduled to receive at completion
	Balance        string         `json:"balance"`         // atoms to receive at completion
}

func (d UnbondingDelegationOutput) HumanReadableString() (string, error) {
	resp := "Unbonding Delegation \n"
	resp += fmt.Sprintf("Delegator: %s\n", d.DelegatorAddr)
	resp += fmt.Sprintf("Validator: %s\n", d.ValidatorAddr)
	resp += fmt.Sprintf("Creation height: %v\n", d.CreationHeight)
	resp += fmt.Sprintf("Min time to unbond (unix): %v\n", d.MinTime)
	resp += fmt.Sprintf("Expected balance: %s", d.Balance)

	return resp, nil

}

type RedelegationOutput struct {
	DelegatorAddr    sdk.AccAddress `json:"delegator_addr"`     // delegator
	ValidatorSrcAddr sdk.AccAddress `json:"validator_src_addr"` // validator redelegation source owner addr
	ValidatorDstAddr sdk.AccAddress `json:"validator_dst_addr"` // validator redelegation destination owner addr
	CreationHeight   int64          `json:"creation_height"`    // height which the redelegation took place
	MinTime          time.Time      `json:"min_time"`           // unix time for redelegation completion
	InitialBalance   string         `json:"initial_balance"`    // initial balance when redelegation started
	Balance          string         `json:"balance"`            // current balance
	SharesSrc        sdk.Rat        `json:"shares_src"`         // amount of source shares redelegating
	SharesDst        sdk.Rat        `json:"shares_dst"`         // amount of destination shares redelegating
}

func (d RedelegationOutput) HumanReadableString() (string, error) {
	resp := "Redelegation \n"
	resp += fmt.Sprintf("Delegator: %s\n", d.DelegatorAddr)
	resp += fmt.Sprintf("Source Validator: %s\n", d.ValidatorSrcAddr)
	resp += fmt.Sprintf("Destination Validator: %s\n", d.ValidatorDstAddr)
	resp += fmt.Sprintf("Creation height: %v\n", d.CreationHeight)
	resp += fmt.Sprintf("Min time to unbond (unix): %v\n", d.MinTime)
	resp += fmt.Sprintf("Source shares: %s", d.SharesSrc.String())
	resp += fmt.Sprintf("Destination shares: %s", d.SharesDst.String())

	return resp, nil

}

type ValidatorOutput struct {
	Owner   sdk.AccAddress `json:"owner"`   // in bech32
	PubKey  string         `json:"pub_key"` // in bech32
	Revoked bool           `json:"revoked"` // has the validator been revoked from bonded status?

	Status          sdk.BondStatus `json:"status"`           // validator status (bonded/unbonding/unbonded)
	Tokens          string        `json:"tokens"`           // delegated tokens (incl. self-delegation)
	DelegatorShares string        `json:"delegator_shares"` // total shares issued to a validator's delegators

	Description        stake.Description `json:"description"`           // description terms for the validator
	BondHeight         int64       `json:"bond_height"`           // earliest height as a bonded validator
	BondIntraTxCounter int16       `json:"bond_intra_tx_counter"` // block-local tx index of validator change
	ProposerRewardPool []string   `json:"proposer_reward_pool"`  // XXX reward pool collected from being the proposer

	Commission            sdk.Rat `json:"commission"`              // XXX the commission rate of fees charged to any delegators
	CommissionMax         sdk.Rat `json:"commission_max"`          // XXX maximum commission rate which this validator can ever charge
	CommissionChangeRate  sdk.Rat `json:"commission_change_rate"`  // XXX maximum daily increase of the validator commission
	CommissionChangeToday sdk.Rat `json:"commission_change_today"` // XXX commission rate change today, reset each day (UTC time)

	// fee related
	LastBondedTokens sdk.Rat `json:"prev_bonded_tokens"` // last bonded token amount
}

func (v ValidatorOutput) HumanReadableString() (string, error) {
	resp := "Validator \n"
	resp += fmt.Sprintf("Owner: %s\n", v.Owner)
	resp += fmt.Sprintf("Validator: %s\n", v.PubKey)
	resp += fmt.Sprintf("Revoked: %v\n", v.Revoked)
	resp += fmt.Sprintf("Status: %s\n", sdk.BondStatusToString(v.Status))
	resp += fmt.Sprintf("Tokens: %s\n", v.Tokens)
	resp += fmt.Sprintf("Delegator Shares: %s\n", v.DelegatorShares)
	resp += fmt.Sprintf("Description: %s\n", v.Description)
	resp += fmt.Sprintf("Bond Height: %d\n", v.BondHeight)
	resp += fmt.Sprintf("Proposer Reward Pool: %s\n", v.ProposerRewardPool)
	resp += fmt.Sprintf("Commission: %s\n", v.Commission.String())
	resp += fmt.Sprintf("Max Commission Rate: %s\n", v.CommissionMax.String())
	resp += fmt.Sprintf("Commission Change Rate: %s\n", v.CommissionChangeRate.String())
	resp += fmt.Sprintf("Commission Change Today: %s\n", v.CommissionChangeToday.String())
	resp += fmt.Sprintf("Previous Bonded Tokens: %s\n", v.LastBondedTokens.String())

	return resp, nil
}

func ExRateFromStakeTokenToMainUnit(cliCtx context.CLIContext) sdk.Rat {
	stakeTokenDenom, err := cliCtx.GetCoinType(app.Denom)
	if err != nil {
		panic(err)
	}
	decimalDiff := stakeTokenDenom.MinUnit.Decimal - stakeTokenDenom.GetMainUnit().Decimal
	exRate := sdk.NewRat(1).Quo(sdk.NewRatFromInt(sdk.NewIntWithDecimal(1, decimalDiff)))
	return exRate
}

func ConvertValidatorToValidatorOutput(cliCtx context.CLIContext, v stake.Validator) (ValidatorOutput, error) {
	exRate := ExRateFromStakeTokenToMainUnit(cliCtx)
	poolToken, err := cliCtx.ConvertCoinToMainUnit(v.ProposerRewardPool.String())
	if err != nil {
		return ValidatorOutput{}, err
	}
	bechValPubkey, err := sdk.Bech32ifyValPub(v.PubKey)
	if err != nil {
		return ValidatorOutput{}, err
	}
	return ValidatorOutput{
		Owner:   v.Owner,
		PubKey:  bechValPubkey,
		Revoked: v.Revoked,

		Status:          v.Status,
		Tokens:          v.Tokens.Mul(exRate).FloatString(),
		DelegatorShares: v.DelegatorShares.Mul(exRate).FloatString(),

		Description:        v.Description,
		BondHeight:         v.BondHeight,
		BondIntraTxCounter: v.BondIntraTxCounter,
		ProposerRewardPool: poolToken,

		Commission:            v.Commission,
		CommissionMax:         v.CommissionMax,
		CommissionChangeRate:  v.CommissionChangeRate,
		CommissionChangeToday: v.CommissionChangeToday,

		LastBondedTokens: v.LastBondedTokens,
	}, nil
}

func ConvertDelegationToDelegationOutput(cliCtx context.CLIContext, delegation stake.Delegation) DelegationOutput {
	exRate := ExRateFromStakeTokenToMainUnit(cliCtx)
	return DelegationOutput{
		DelegatorAddr: delegation.DelegatorAddr,
		ValidatorAddr: delegation.ValidatorAddr,
		Shares:        delegation.Shares.Mul(exRate).FloatString(),
		Height:        delegation.Height,
	}
}

func ConvertUBDToUBDOutput(cliCtx context.CLIContext, ubd stake.UnbondingDelegation) UnbondingDelegationOutput {
	initialBalance, err := cliCtx.ConvertCoinToMainUnit(sdk.Coins{ubd.InitialBalance}.String())
	if err != nil && len(initialBalance) != 1 {
		panic(err)
	}
	balance, err := cliCtx.ConvertCoinToMainUnit(sdk.Coins{ubd.Balance}.String())
	if err != nil && len(balance) != 1 {
		panic(err)
	}
	return UnbondingDelegationOutput{
		DelegatorAddr:  ubd.DelegatorAddr,
		ValidatorAddr:  ubd.ValidatorAddr,
		CreationHeight: ubd.CreationHeight,
		MinTime:        ubd.MinTime,
		InitialBalance: initialBalance[0],
		Balance:        balance[0],
	}
}

func ConvertREDToREDOutput(cliCtx context.CLIContext, red stake.Redelegation) RedelegationOutput {
	exRate := ExRateFromStakeTokenToMainUnit(cliCtx)
	initialBalance, err := cliCtx.ConvertCoinToMainUnit(sdk.Coins{red.InitialBalance}.String())
	if err != nil && len(initialBalance) != 1 {
		panic(err)
	}
	balance, err := cliCtx.ConvertCoinToMainUnit(sdk.Coins{red.Balance}.String())
	if err != nil && len(balance) != 1 {
		panic(err)
	}
	return RedelegationOutput{
		DelegatorAddr:    red.DelegatorAddr,
		ValidatorSrcAddr: red.ValidatorSrcAddr,
		ValidatorDstAddr: red.ValidatorDstAddr,
		CreationHeight:   red.CreationHeight,
		MinTime:          red.MinTime,
		InitialBalance:   initialBalance[0],
		Balance:          balance[0],
		SharesSrc:        red.SharesSrc.Mul(exRate),
		SharesDst:        red.SharesDst.Mul(exRate),
	}
}
