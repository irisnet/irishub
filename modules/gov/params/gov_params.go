package govparams

import (
	"encoding/json"
	"fmt"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/params"
	"github.com/irisnet/irishub/types"
	"strconv"
	"time"
	stakeTypes "github.com/irisnet/irishub/modules/stake/types"
)

var DepositProcedureParameter DepositProcedureParam

const LOWER_BOUND_AMOUNT = 10
const UPPER_BOUND_AMOUNT = 10000
const THREE_DAYS = 3*3600*24
var _ params.GovParameter = (*DepositProcedureParam)(nil)

type ParamSet struct {
	DepositProcedure  DepositProcedure  `json:"Gov/govDepositProcedure"`
	VotingProcedure   VotingProcedure   `json:"Gov/govVotingProcedure"`
	TallyingProcedure TallyingProcedure `json:"Gov/govTallyingProcedure"`
}

// Procedure around Deposits for governance
type DepositProcedure struct {
	MinDeposit       sdk.Coins        `json:"min_deposit"`        //  Minimum deposit for a proposal to enter voting period.
	MaxDepositPeriod time.Duration    `json:"max_deposit_period"` //  Maximum period for Atom holders to deposit on a proposal. Initial value: 2 months
}

type DepositProcedureParam struct {
	Value   DepositProcedure
	paramSpace params.Subspace
}

func (param *DepositProcedureParam) GetValueFromRawData(cdc *codec.Codec, res []byte) interface{} {
	cdc.UnmarshalJSON(res, &param.Value)
	return param.Value
}

func (param *DepositProcedureParam) InitGenesis(genesisState interface{}) {
	if value, ok := genesisState.(DepositProcedure); ok {
		param.Value = value
	} else {
		var minDeposit, _ = types.NewDefaultCoinType(stakeTypes.StakeDenomName).ConvertToMinCoin(fmt.Sprintf("%d%s", 10, stakeTypes.StakeDenomName))

		param.Value = DepositProcedure{
			MinDeposit:       sdk.Coins{minDeposit},
			MaxDepositPeriod: time.Duration(172800) * time.Second}
	}
}

func (param *DepositProcedureParam) SetReadWriter(paramSpace params.Subspace) {
	param.paramSpace = paramSpace
}

func (param *DepositProcedureParam) GetStoreKey() []byte {
	return []byte("govDepositProcedure")
}

func (param *DepositProcedureParam) SaveValue(ctx sdk.Context) {
	param.paramSpace.Set(ctx, param.GetStoreKey(), param.Value)
}

func (param *DepositProcedureParam) LoadValue(ctx sdk.Context) bool {
	if param.paramSpace.Has(ctx, param.GetStoreKey()) == false {
		return false
	}
	param.paramSpace.Get(ctx, param.GetStoreKey(), &param.Value)
	return true
}

func (param *DepositProcedureParam) ToJson(jsonStr string) string {

	var jsonBytes []byte

	if len(jsonStr) == 0 {
		jsonBytes, _  = json.Marshal(param.Value)
		return string(jsonBytes)
	}

	if err := json.Unmarshal([]byte(jsonStr), &param.Value); err == nil {
		jsonBytes, _  = json.Marshal(param.Value)
		return string(jsonBytes)
	}
	return string(jsonBytes)
}

func (param *DepositProcedureParam) Update(ctx sdk.Context, jsonStr string) {
	if err := json.Unmarshal([]byte(jsonStr), &param.Value); err == nil {
		param.SaveValue(ctx)
	}
}

func (param *DepositProcedureParam) Valid(jsonStr string) sdk.Error {

	var err error

	if err = json.Unmarshal([]byte(jsonStr), &param.Value); err == nil {

		if param.Value.MinDeposit[0].Denom != stakeTypes.StakeDenom {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMinDepositDenom, fmt.Sprintf("It should be %s!", stakeTypes.StakeDenom))
		}

		LowerBound, _ := types.NewDefaultCoinType(stakeTypes.StakeDenomName).ConvertToMinCoin(fmt.Sprintf("%d%s", LOWER_BOUND_AMOUNT, stakeTypes.StakeDenomName))
		UpperBound, _ := types.NewDefaultCoinType(stakeTypes.StakeDenomName).ConvertToMinCoin(fmt.Sprintf("%d%s", UPPER_BOUND_AMOUNT, stakeTypes.StakeDenomName))

		if param.Value.MinDeposit[0].Amount.LT(LowerBound.Amount) || param.Value.MinDeposit[0].Amount.GT(UpperBound.Amount) {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMinDepositAmount, fmt.Sprintf("MinDepositAmount"+param.Value.MinDeposit[0].String()+" should be larger than 10iris and less than 10000iris"))

		}

		if param.Value.MaxDepositPeriod.Seconds() < 20 || param.Value.MaxDepositPeriod.Seconds() > THREE_DAYS {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidDepositPeriod, fmt.Sprintf("MaxDepositPeriod (%s) should be larger than 20s and less than %ds",strconv.Itoa(int(param.Value.MaxDepositPeriod.Seconds())), THREE_DAYS))
		}

		return nil

	}
	return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMinDeposit, fmt.Sprintf("Json is not valid"))
}

var VotingProcedureParameter VotingProcedureParam
var _ params.GovParameter = (*VotingProcedureParam)(nil)

// Procedure around Voting in governance
type VotingProcedure struct {
	VotingPeriod time.Duration `json:"voting_period"` //  Length of the voting period.
}

type VotingProcedureParam struct {
	Value   VotingProcedure
	paramSpace params.Subspace
}

func (param *VotingProcedureParam) GetValueFromRawData(cdc *codec.Codec, res []byte) interface{} {
	cdc.UnmarshalJSON(res, &param.Value)
	return param.Value
}

func (param *VotingProcedureParam) InitGenesis(genesisState interface{}) {
	if value, ok := genesisState.(VotingProcedure); ok {
		param.Value = value
	} else {
		param.Value = VotingProcedure{VotingPeriod: time.Duration(172800) * time.Second}
	}
}

func (param *VotingProcedureParam) SetReadWriter(paramSpace params.Subspace) {
	param.paramSpace = paramSpace
}

func (param *VotingProcedureParam) GetStoreKey() []byte {
	return []byte("govVotingProcedure")
}

func (param *VotingProcedureParam) SaveValue(ctx sdk.Context) {
	param.paramSpace.Set(ctx, param.GetStoreKey(), param.Value)
}

func (param *VotingProcedureParam) LoadValue(ctx sdk.Context) bool {
	if param.paramSpace.Has(ctx, param.GetStoreKey()) == false {
		return false
	}
	param.paramSpace.Get(ctx, param.GetStoreKey(), &param.Value)
	return true
}

func (param *VotingProcedureParam) ToJson(jsonStr string) string {
	var jsonBytes []byte

	if len(jsonStr) == 0 {
		jsonBytes, _  = json.Marshal(param.Value)
		return string(jsonBytes)
	}

	if err := json.Unmarshal([]byte(jsonStr), &param.Value); err == nil {
		jsonBytes, _  = json.Marshal(param.Value)
		return string(jsonBytes)
	}
	return string(jsonBytes)
}

func (param *VotingProcedureParam) Update(ctx sdk.Context, jsonStr string) {
	if err := json.Unmarshal([]byte(jsonStr), &param.Value); err == nil {
		param.SaveValue(ctx)
	}
}

func (param *VotingProcedureParam) Valid(jsonStr string) sdk.Error {

	var err error

	if err = json.Unmarshal([]byte(jsonStr), &param.Value); err == nil {

		if param.Value.VotingPeriod.Seconds() < 20 || param.Value.VotingPeriod.Seconds() > THREE_DAYS {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidVotingPeriod, fmt.Sprintf("VotingPeriod (%s) should be larger than 20s and less than %ds",strconv.Itoa(int(param.Value.VotingPeriod.Seconds())), THREE_DAYS))
		}

		return nil

	}
	return sdk.NewError(params.DefaultCodespace, params.CodeInvalidVotingProcedure, fmt.Sprintf("Json is not valid"))
}

var TallyingProcedureParameter TallyingProcedureParam
var _ params.GovParameter = (*TallyingProcedureParam)(nil)

// Procedure around Tallying votes in governance
type TallyingProcedure struct {
	Threshold         sdk.Dec `json:"threshold"`          //  Minimum propotion of Yes votes for proposal to pass. Initial value: 0.5
	Veto              sdk.Dec `json:"veto"`               //  Minimum value of Veto votes to Total votes ratio for proposal to be vetoed. Initial value: 1/3
	Participation     sdk.Dec `json:"participation"`      //
}

type TallyingProcedureParam struct {
	Value   TallyingProcedure
	paramSpace params.Subspace
}

func (param *TallyingProcedureParam) GetValueFromRawData(cdc *codec.Codec, res []byte) interface{} {
	cdc.UnmarshalJSON(res, &param.Value)
	return param.Value
}

func (param *TallyingProcedureParam) InitGenesis(genesisState interface{}) {
	if value, ok := genesisState.(TallyingProcedure); ok {
		param.Value = value
	} else {
		param.Value = TallyingProcedure{
			Threshold:         sdk.NewDecWithPrec(5, 1),
			Veto:              sdk.NewDecWithPrec(334, 3),
			Participation:     sdk.NewDecWithPrec(667, 3),
		}
	}
}

func (param *TallyingProcedureParam) SetReadWriter(paramSpace params.Subspace) {
	param.paramSpace = paramSpace
}

func (param *TallyingProcedureParam) GetStoreKey() []byte {
	return []byte("govTallyingProcedure")
}

func (param *TallyingProcedureParam) SaveValue(ctx sdk.Context) {
	param.paramSpace.Set(ctx, param.GetStoreKey(), param.Value)
}

func (param *TallyingProcedureParam) LoadValue(ctx sdk.Context) bool {
	if param.paramSpace.Has(ctx, param.GetStoreKey()) == false {
		return false
	}
	param.paramSpace.Get(ctx, param.GetStoreKey(), &param.Value)
	return true
}

func (param *TallyingProcedureParam) ToJson(jsonStr string) string {
	var jsonBytes []byte

	if len(jsonStr) == 0 {
		jsonBytes, _  = json.Marshal(param.Value)
		return string(jsonBytes)
	}

	if err := json.Unmarshal([]byte(jsonStr), &param.Value); err == nil {
		jsonBytes, _  = json.Marshal(param.Value)
		return string(jsonBytes)
	}
	return string(jsonBytes)
}

func (param *TallyingProcedureParam) Update(ctx sdk.Context, jsonStr string) {
	if err := json.Unmarshal([]byte(jsonStr), &param.Value); err == nil {
		param.SaveValue(ctx)
	}
}

func (param *TallyingProcedureParam) Valid(jsonStr string) sdk.Error {

	var err error

	if err = json.Unmarshal([]byte(jsonStr), &param.Value); err == nil {

		if param.Value.Threshold.LTE(sdk.ZeroDec()) || param.Value.Threshold.GTE(sdk.NewDec(1)) {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidThreshold, fmt.Sprintf("Invalid Threshold ( "+param.Value.Threshold.String()+" ) should be between 0 and 1"))
		}
		if param.Value.Participation.LTE(sdk.ZeroDec()) || param.Value.Participation.GTE(sdk.NewDec(1)) {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidParticipation, fmt.Sprintf("Invalid participation ( "+param.Value.Participation.String()+" ) should be between 0 and 1"))
		}
		if param.Value.Veto.LTE(sdk.ZeroDec()) || param.Value.Veto.GTE(sdk.NewDec(1)) {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidVeto, fmt.Sprintf("Invalid Veto ( "+param.Value.Veto.String()+" ) should be between 0 and 1"))
		}

		return nil

	}
	return sdk.NewError(params.DefaultCodespace, params.CodeInvalidTallyingProcedure, fmt.Sprintf("Json is not valid"))
}
