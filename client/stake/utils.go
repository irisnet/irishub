package stake

import (
	"fmt"
	"strings"
	"time"

	"github.com/irisnet/irishub/app/v1/stake"
	"github.com/irisnet/irishub/app/v1/stake/types"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
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

func (d DelegationOutput) String() string {
	return fmt.Sprintf(`Delegation:
  Delegator:  %s
  Validator:  %s
  Shares:     %s
  Height:     %v`, d.DelegatorAddr,
		d.ValidatorAddr, d.Shares, d.Height)
}

// Delegations is a collection of delegations
type DelegationsOutput []DelegationOutput

func (ds DelegationsOutput) String() (out string) {
	if len(ds) == 0 {
		return "[]"
	}
	for _, del := range ds {
		out += del.String() + "\n"
	}
	return strings.TrimSpace(out)
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

func (d UnbondingDelegationOutput) String() string {
	return fmt.Sprintf(`Unbonding Delegation:
  Delegator Address:          %s
  Validator Address:          %s
  Creation Height:            %v
  Min time to unbond (unix):  %s
  Initial Balance:            %s
  Balance:                    %s`,
		d.DelegatorAddr, d.ValidatorAddr, d.CreationHeight, d.MinTime, d.InitialBalance, d.Balance)
}

// Validators is a collection of Validator
type UnbondingDelegationsOutput []UnbondingDelegationOutput

func (ubds UnbondingDelegationsOutput) String() (out string) {
	if len(ubds) == 0 {
		return "[]"
	}
	for _, val := range ubds {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
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

func (d RedelegationOutput) String() string {
	return fmt.Sprintf(`Redelegation:
  Delegator:                  %s
  Source Validator:           %s
  Destination Validator:      %s
  Creation height:            %v
  Min time to unbond (unix):  %v
  Source shares:              %s
  Destination shares:         %s`, d.DelegatorAddr, d.ValidatorSrcAddr, d.ValidatorDstAddr,
		d.CreationHeight, d.MinTime, d.SharesSrc, d.SharesDst)
}

// Redelegations are a collection of Redelegation
type RedelegationsOutput []RedelegationOutput

func (reds RedelegationsOutput) String() (out string) {
	if len(reds) == 0 {
		return "[]"
	}
	for _, red := range reds {
		out += red.String() + "\n"
	}
	return strings.TrimSpace(out)
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
	Commission       stake.Commission  `json:"commission"`
}

func (v ValidatorOutput) String() string {
	return fmt.Sprintf(`Validator
  Operator Address:            %s
  Validator Consensus Pubkey:  %s
  Jailed:                      %v
  Status:                      %s
  Tokens:                      %s
  Delegator Shares:            %s
  Description:                 %s
  Unbonding Height:            %d
  Minimum Unbonding Time:      %v
  Commission:                  %s`, v.OperatorAddr, v.ConsPubKey,
		v.Jailed, sdk.BondStatusToString(v.Status), v.Tokens,
		v.DelegatorShares, v.Description,
		v.UnbondingHeight, v.UnbondingMinTime, v.Commission)
}

// Validators is a collection of Validator
type ValidatorsOutput []ValidatorOutput

func (v ValidatorsOutput) String() (out string) {
	if len(v) == 0 {
		return "[]"
	}
	for _, val := range v {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
}

type PoolOutput struct {
	LooseTokens  string `json:"loose_tokens"`
	BondedTokens string `json:"bonded_tokens"`
	TokenSupply  string `json:"total_supply"`
	BondedRatio  string `json:"bonded_ratio"`
}

func (p PoolOutput) String() string {
	return fmt.Sprintf(`Pool:
  Loose Tokens:   %s
  Bonded Tokens:  %s
  Token Supply:   %s
  Bonded Ratio:   %v`, p.LooseTokens,
		p.BondedTokens, p.TokenSupply, p.BondedRatio)
}

func ConvertValidatorToValidatorOutput(cliCtx context.CLIContext, v stake.Validator) ValidatorOutput {
	exRate := utils.ExRateFromStakeTokenToMainUnit(cliCtx)

	bechConsPubkey, err := sdk.Bech32ifyConsPub(v.ConsPubKey)
	if err != nil {
		panic(err)
	}

	return ValidatorOutput{
		OperatorAddr:     v.OperatorAddr,
		ConsPubKey:       bechConsPubkey,
		Jailed:           v.Jailed,
		Status:           v.Status,
		Tokens:           utils.ConvertDecToRat(v.Tokens).Mul(exRate).FloatString(),
		DelegatorShares:  utils.ConvertDecToRat(v.DelegatorShares).Mul(exRate).FloatString(),
		Description:      v.Description,
		BondHeight:       v.BondHeight,
		UnbondingHeight:  v.UnbondingHeight,
		UnbondingMinTime: v.UnbondingMinTime,
		Commission:       v.Commission,
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
	initialBalance, err := cliCtx.ConvertToMainUnit(sdk.Coins{ubd.InitialBalance}.String())
	if err != nil && len(initialBalance) != 1 {
		panic(err)
	}
	balance, err := cliCtx.ConvertToMainUnit(sdk.Coins{ubd.Balance}.String())
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
	initialBalance, err := cliCtx.ConvertToMainUnit(sdk.Coins{red.InitialBalance}.String())
	if err != nil && len(initialBalance) != 1 {
		panic(err)
	}
	balance, err := cliCtx.ConvertToMainUnit(sdk.Coins{red.Balance}.String())
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
		exRate := sdk.NewDecFromInt(sdk.NewIntWithDecimal(1, int(decimalDiff)))
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
