package stake

import (
	"fmt"
	"time"

	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/stake"
	"github.com/irisnet/irishub/modules/stake/types"
	sdk "github.com/irisnet/irishub/types"
	"github.com/pkg/errors"
)

// defines a delegation without type Rat for shares
type DelegationOutput struct {
	DelegatorAddr sdk.AccAddress `json:"delegator_addr"`
	ValidatorAddr sdk.ValAddress `json:"validator_addr"`
	Shares        string         `json:"shares"`
	Height        int64          `json:"height"`
}

func (d DelegationOutput) HumanReadableString() (string, error) {
	resp := "Delegation \n"
	resp += fmt.Sprintf("Delegator: %s\n", d.DelegatorAddr)
	resp += fmt.Sprintf("Validator: %s\n", d.ValidatorAddr)
	resp += fmt.Sprintf("Shares: %s\n", d.Shares)
	resp += fmt.Sprintf("Height: %d", d.Height)

	return resp, nil
}

// UnbondingDelegation reflects a delegation's passive unbonding queue.
type UnbondingDelegationOutput struct {
	DelegatorAddr  sdk.AccAddress `json:"delegator_addr"`  // delegator
	ValidatorAddr  sdk.ValAddress `json:"validator_addr"`  // validator unbonding from owner addr
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
	ValidatorSrcAddr sdk.ValAddress `json:"validator_src_addr"` // validator redelegation source owner addr
	ValidatorDstAddr sdk.ValAddress `json:"validator_dst_addr"` // validator redelegation destination owner addr
	CreationHeight   int64          `json:"creation_height"`    // height which the redelegation took place
	MinTime          time.Time      `json:"min_time"`           // unix time for redelegation completion
	InitialBalance   string         `json:"initial_balance"`    // initial balance when redelegation started
	Balance          string         `json:"balance"`            // current balance
	SharesSrc        string         `json:"shares_src"`         // amount of source shares redelegating
	SharesDst        string         `json:"shares_dst"`         // amount of destination shares redelegating
}

func (d RedelegationOutput) HumanReadableString() (string, error) {
	resp := "Redelegation \n"
	resp += fmt.Sprintf("Delegator: %s\n", d.DelegatorAddr)
	resp += fmt.Sprintf("Source Validator: %s\n", d.ValidatorSrcAddr)
	resp += fmt.Sprintf("Destination Validator: %s\n", d.ValidatorDstAddr)
	resp += fmt.Sprintf("Creation height: %v\n", d.CreationHeight)
	resp += fmt.Sprintf("Min time to unbond (unix): %v\n", d.MinTime)
	resp += fmt.Sprintf("Source shares: %s\n", d.SharesSrc)
	resp += fmt.Sprintf("Destination shares: %s", d.SharesDst)

	return resp, nil

}

type Commission struct {
	Rate          string    `json:"rate"`
	MaxRate       string    `json:"max_rate"`
	MaxChangeRate string    `json:"max_change_rate"`
	UpdateTime    time.Time `json:"update_time"`
}

type ValidatorOutput struct {
	OperatorAddr     sdk.ValAddress    `json:"operator_address"`
	ConsPubKey       string            `json:"consensus_pubkey"`
	Jailed           bool              `json:"jailed"`
	Status           sdk.BondStatus    `json:"status"`
	Tokens           string            `json:"tokens"`
	DelegatorShares  string            `json:"delegator_shares"`
	Description      stake.Description `json:"description"`
	BondHeight       int64             `json:"bond_height"`
	UnbondingHeight  int64             `json:"unbonding_height"`
	UnbondingMinTime time.Time         `json:"unbonding_time"`
	Commission       Commission        `json:"commission"`
}

func (v ValidatorOutput) HumanReadableString() (string, error) {
	resp := "Validator \n"
	resp += fmt.Sprintf("Operator Address: %s\n", v.OperatorAddr)
	resp += fmt.Sprintf("Validator Consensus Pubkey: %s\n", v.ConsPubKey)
	resp += fmt.Sprintf("Jailed: %v\n", v.Jailed)
	resp += fmt.Sprintf("Status: %s\n", sdk.BondStatusToString(v.Status))
	resp += fmt.Sprintf("Tokens: %s\n", v.Tokens)
	resp += fmt.Sprintf("Delegator Shares: %s\n", v.DelegatorShares)
	resp += fmt.Sprintf("Description: %s\n", v.Description)
	resp += fmt.Sprintf("Bond Height: %d\n", v.BondHeight)
	resp += fmt.Sprintf("Unbonding Height: %d\n", v.UnbondingHeight)
	resp += fmt.Sprintf("Minimum Unbonding Time: %v\n", v.UnbondingMinTime)
	resp += fmt.Sprintf("Commission: {%s}", v.Commission)

	return resp, nil
}

type PoolOutput struct {
	LooseTokens  string `json:"loose_tokens"`
	BondedTokens string `json:"bonded_tokens"`
	TokenSupply  string `json:"total_supply"`
	BondedRatio  string `json:"bonded_ratio"`
}

func (p PoolOutput) HumanReadableString() string {

	resp := "Pool \n"
	resp += fmt.Sprintf("Loose Tokens: %s\n", p.LooseTokens)
	resp += fmt.Sprintf("Bonded Tokens: %s\n", p.BondedTokens)
	resp += fmt.Sprintf("Token Supply: %s\n", p.TokenSupply)
	resp += fmt.Sprintf("Bonded Ratio: %v", p.BondedRatio)
	return resp
}

func ConvertValidatorToValidatorOutput(cliCtx context.CLIContext, v stake.Validator) ValidatorOutput {
	exRate := utils.ExRateFromStakeTokenToMainUnit(cliCtx)

	bechValPubkey, err := sdk.Bech32ifyValPub(v.ConsPubKey)
	if err != nil {
		panic(err)
	}

	commission := Commission{
		Rate:          utils.ConvertDecToRat(v.Commission.Rate).FloatString(),
		MaxRate:       utils.ConvertDecToRat(v.Commission.MaxRate).FloatString(),
		MaxChangeRate: utils.ConvertDecToRat(v.Commission.MaxChangeRate).FloatString(),
		UpdateTime:    v.Commission.UpdateTime,
	}
	return ValidatorOutput{
		OperatorAddr:     v.OperatorAddr,
		ConsPubKey:       bechValPubkey,
		Jailed:           v.Jailed,
		Status:           v.Status,
		Tokens:           utils.ConvertDecToRat(v.Tokens).Mul(exRate).FloatString(),
		DelegatorShares:  utils.ConvertDecToRat(v.DelegatorShares).Mul(exRate).FloatString(),
		Description:      v.Description,
		BondHeight:       v.UnbondingHeight,
		UnbondingHeight:  v.UnbondingHeight,
		UnbondingMinTime: v.UnbondingMinTime,
		Commission:       commission,
	}
}

func ConvertDelegationToDelegationOutput(cliCtx context.CLIContext, delegation stake.Delegation) DelegationOutput {
	exRate := utils.ExRateFromStakeTokenToMainUnit(cliCtx)
	return DelegationOutput{
		DelegatorAddr: delegation.DelegatorAddr,
		ValidatorAddr: delegation.ValidatorAddr,
		Shares:        utils.ConvertDecToRat(delegation.Shares).Mul(exRate).FloatString(),
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
	exRate := utils.ExRateFromStakeTokenToMainUnit(cliCtx)
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
		SharesSrc:        utils.ConvertDecToRat(red.SharesSrc).Mul(exRate).FloatString(),
		SharesDst:        utils.ConvertDecToRat(red.SharesDst).Mul(exRate).FloatString(),
	}
}

func ConvertPoolToPoolOutput(cliCtx context.CLIContext, pool stake.PoolStatus) PoolOutput {
	exRate := utils.ExRateFromStakeTokenToMainUnit(cliCtx)
	return PoolOutput{
		LooseTokens:  utils.ConvertDecToRat(pool.LooseTokens).Mul(exRate).FloatString(),
		BondedTokens: utils.ConvertDecToRat(pool.BondedTokens).Mul(exRate).FloatString(),
		TokenSupply:  utils.ConvertDecToRat(pool.BondedTokens.Add(pool.LooseTokens)).Mul(exRate).FloatString(),
		BondedRatio:  utils.ConvertDecToRat(pool.BondedTokens.Quo(pool.BondedTokens.Add(pool.LooseTokens))).FloatString(),
	}
}

func BuildCommissionMsg(rateStr, maxRateStr, maxChangeRateStr string) (commission types.CommissionMsg, err error) {
	if rateStr == "" || maxRateStr == "" || maxChangeRateStr == "" {
		return commission, fmt.Errorf("must specify all validator commission parameters")
	}

	rate, err := sdk.NewDecFromStr(rateStr)
	if err != nil {
		return commission, err
	}

	maxRate, err := sdk.NewDecFromStr(maxRateStr)
	if err != nil {
		return commission, err
	}

	maxChangeRate, err := sdk.NewDecFromStr(maxChangeRateStr)
	if err != nil {
		return commission, err
	}

	commission = types.NewCommissionMsg(rate, maxRate, maxChangeRate)
	return commission, nil
}

// nolint: gocyclo
// TODO: Make this pass gocyclo linting
func GetShares(
	storeName string, cliCtx context.CLIContext, cdc *codec.Codec, sharesAmountStr,
	sharesPercentStr string, delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress,
) (sharesAmount sdk.Dec, err error) {
	switch {
	case sharesAmountStr != "" && sharesPercentStr != "":
		return sharesAmount, errors.Errorf("can either specify the amount OR the percent of the shares, not both")

	case sharesAmountStr == "" && sharesPercentStr == "":
		return sharesAmount, errors.Errorf("can either specify the amount OR the percent of the shares, not both")

	case sharesAmountStr != "":
		sharesAmount, err = sdk.NewDecFromStr(sharesAmountStr)
		if err != nil {
			return sharesAmount, err
		}
		if !sharesAmount.GT(sdk.ZeroDec()) {
			return sharesAmount, errors.Errorf("shares amount must be positive number (ex. 123, 1.23456789)")
		}

		stakeToken, err := cliCtx.GetCoinType(types.StakeTokenName)
		if err != nil {
			panic(err)
		}
		decimalDiff := stakeToken.MinUnit.Decimal - stakeToken.GetMainUnit().Decimal
		exRate := sdk.NewDecFromInt(sdk.NewIntWithDecimal(1, decimalDiff))
		sharesAmount = sharesAmount.Mul(exRate)
	case sharesPercentStr != "":
		var sharesPercent sdk.Dec
		sharesPercent, err = sdk.NewDecFromStr(sharesPercentStr)
		if err != nil {
			return sharesAmount, err
		}
		if !sharesPercent.GT(sdk.ZeroDec()) || !sharesPercent.LTE(sdk.OneDec()) {
			return sharesAmount, errors.Errorf("shares percent must be >0 and <=1 (ex. 0.01, 0.75, 1)")
		}

		// make a query to get the existing delegation shares
		key := stake.GetDelegationKey(delegatorAddr, validatorAddr)

		resQuery, err := cliCtx.QueryStore(key, storeName)
		if err != nil {
			return sharesAmount, errors.Errorf("cannot find delegation to determine percent Error: %v", err)
		} else if len(resQuery) == 0 {
			return sharesAmount, errors.Errorf("delegation (from delegator %s to validator %s) doesn't exist", delegatorAddr.String(), validatorAddr.String())
		}

		delegation, err := types.UnmarshalDelegation(cdc, key, resQuery)
		if err != nil {
			return sdk.ZeroDec(), err
		}

		sharesAmount = sharesPercent.Mul(delegation.Shares)
	}
	return
}
