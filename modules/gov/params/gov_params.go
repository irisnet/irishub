package govparams

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/parameter"

	"strconv"
)

var DepositProcedureParameter DepositProcedureParam


var (
	minDeposit, _ = types.NewDefaultCoinType("iris").ConvertToMinCoin(fmt.Sprintf("%d%s", 10, "iris"))
)

const LOWER_BOUND_AMOUNT = 1
const UPPER_BOUND_AMOUNT = 200

var _ parameter.GovParameter = (*DepositProcedureParam)(nil)

type ParamSet struct {
	DepositProcedure   DepositProcedure  `json:"Gov/gov/DepositProcedure"`
	VotingProcedure    VotingProcedure   `json:"Gov/gov/VotingProcedure"`
	TallyingProcedure  TallyingProcedure `json:"Gov/gov/TallyingProcedure"`
}

// Procedure around Deposits for governance
type DepositProcedure struct {
	MinDeposit       sdk.Coins `json:"min_deposit"`        //  Minimum deposit for a proposal to enter voting period.
	MaxDepositPeriod int64     `json:"max_deposit_period"` //  Maximum period for Atom holders to deposit on a proposal. Initial value: 2 months
}

type DepositProcedureParam struct {
	Value   DepositProcedure
	psetter params.Setter
	pgetter params.Getter
}

func (param *DepositProcedureParam) InitGenesis(genesisState interface{}) {
	if value, ok := genesisState.(DepositProcedure); ok {
		param.Value = value
	} else {
		param.Value = DepositProcedure{
			MinDeposit:       sdk.Coins{minDeposit},
			MaxDepositPeriod: 1440}
	}
}

func (param *DepositProcedureParam) SetReadWriter(setter params.Setter) {
	param.psetter = setter
	param.pgetter = setter.Getter
}

func (param *DepositProcedureParam) GetStoreKey() string {
	return "Gov/gov/DepositProcedure"
}

func (param *DepositProcedureParam) SaveValue(ctx sdk.Context) {
	param.psetter.Set(ctx, param.GetStoreKey(), param.Value)
}

func (param *DepositProcedureParam) LoadValue(ctx sdk.Context) bool {
	err := param.pgetter.Get(ctx, param.GetStoreKey(), &param.Value)
	if err != nil {
		return false
	}
	return true
}

func (param *DepositProcedureParam) ToJson() string {
	jsonBytes, _ := json.Marshal(param.Value)
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


		if param.Value.MinDeposit[0].Denom != "iris-atto" {
			return sdk.NewError(parameter.DefaultCodespace, parameter.CodeInvalidMinDepositDenom, fmt.Sprintf("It should be iris-atto! git"))
		}

		LowerBound, _ := types.NewDefaultCoinType("iris").ConvertToMinCoin(fmt.Sprintf("%d%s", LOWER_BOUND_AMOUNT, "iris"))
		UpperBound, _ := types.NewDefaultCoinType("iris").ConvertToMinCoin(fmt.Sprintf("%d%s", UPPER_BOUND_AMOUNT, "iris"))

		if param.Value.MinDeposit[0].Amount.LT(LowerBound.Amount) || param.Value.MinDeposit[0].Amount.GT(UpperBound.Amount) {
			return sdk.NewError(parameter.DefaultCodespace, parameter.CodeInvalidMinDepositAmount, fmt.Sprintf("MinDepositAmount"+param.Value.MinDeposit[0].String()+" should be larger than 10 and less than 20000"))

		}

		if param.Value.MaxDepositPeriod < 20 || param.Value.MaxDepositPeriod > 20000 {
			return sdk.NewError(parameter.DefaultCodespace, parameter.CodeInvalidDepositPeriod, fmt.Sprintf("MaxDepositPeriod ("+strconv.Itoa(int(param.Value.MaxDepositPeriod))+") should be larger than 20 and less than 20000"))
		}

		return nil

	}
	return sdk.NewError(parameter.DefaultCodespace, parameter.CodeInvalidMinDeposit, fmt.Sprintf("Json is not valid"))
}


var VotingProcedureParameter VotingProcedureParam
var _ parameter.GovParameter = (*VotingProcedureParam)(nil)

// Procedure around Voting in governance
type VotingProcedure struct {
	VotingPeriod int64 `json:"voting_period"` //  Length of the voting period.
}

type VotingProcedureParam struct {
	Value   VotingProcedure
	psetter params.Setter
	pgetter params.Getter
}

func (param *VotingProcedureParam) InitGenesis(genesisState interface{}) {
	if value, ok := genesisState.(VotingProcedure); ok {
		param.Value = value
	} else {
		param.Value = VotingProcedure{VotingPeriod: 1000,
		}
	}
}

func (param *VotingProcedureParam) SetReadWriter(setter params.Setter) {
	param.psetter = setter
	param.pgetter = setter.Getter
}

func (param *VotingProcedureParam) GetStoreKey() string {
	return "Gov/gov/votingProcedure"
}

func (param *VotingProcedureParam) SaveValue(ctx sdk.Context) {
	param.psetter.Set(ctx, param.GetStoreKey(), param.Value)
}

func (param *VotingProcedureParam) LoadValue(ctx sdk.Context) bool {
	err := param.pgetter.Get(ctx, param.GetStoreKey(), &param.Value)
	if err != nil {
		return false
	}
	return true
}

func (param *VotingProcedureParam) ToJson() string {
	jsonBytes, _ := json.Marshal(param.Value)
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

		if param.Value.VotingPeriod < 20 || param.Value.VotingPeriod > 20000 {
			return sdk.NewError(parameter.DefaultCodespace, parameter.CodeInvalidVotingPeriod, fmt.Sprintf("VotingPeriod ("+strconv.Itoa(int(param.Value.VotingPeriod))+") should be larger than 20 and less than 20000"))
		}

		return nil

	}
	return sdk.NewError(parameter.DefaultCodespace, parameter.CodeInvalidVotingProcedure, fmt.Sprintf("Json is not valid"))
}

var TallyingProcedureParameter TallyingProcedureParam
var _ parameter.GovParameter = (*TallyingProcedureParam)(nil)

// Procedure around Tallying votes in governance
type TallyingProcedure struct {
	Threshold         sdk.Rat `json:"threshold"`          //  Minimum propotion of Yes votes for proposal to pass. Initial value: 0.5
	Veto              sdk.Rat `json:"veto"`               //  Minimum value of Veto votes to Total votes ratio for proposal to be vetoed. Initial value: 1/3
	GovernancePenalty sdk.Rat `json:"governance_penalty"` //  Penalty if validator does not vote
}


type TallyingProcedureParam struct {
	Value   TallyingProcedure
	psetter params.Setter
	pgetter params.Getter
}

func (param *TallyingProcedureParam) InitGenesis(genesisState interface{}) {
	if value, ok := genesisState.(TallyingProcedure); ok {
		param.Value = value
	} else {
		param.Value = TallyingProcedure{
			Threshold:         sdk.NewRat(1, 2),
			Veto:              sdk.NewRat(1, 3),
			GovernancePenalty: sdk.NewRat(1, 100),
		}
	}
}

func (param *TallyingProcedureParam) SetReadWriter(setter params.Setter) {
	param.psetter = setter
	param.pgetter = setter.Getter
}

func (param *TallyingProcedureParam) GetStoreKey() string {
	return "Gov/gov/TallyingProcedure"
}

func (param *TallyingProcedureParam) SaveValue(ctx sdk.Context) {
	param.psetter.Set(ctx, param.GetStoreKey(), param.Value)
}

func (param *TallyingProcedureParam) LoadValue(ctx sdk.Context) bool {
	err := param.pgetter.Get(ctx, param.GetStoreKey(), &param.Value)
	if err != nil {
		return false
	}
	return true
}

func (param *TallyingProcedureParam) ToJson() string {
	jsonBytes, _ := json.Marshal(param.Value)
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

		if param.Value.Threshold.LT(sdk.NewRat(0)) || param.Value.Threshold.GT(sdk.NewRat(1)) {
			return sdk.NewError(parameter.DefaultCodespace, parameter.CodeInvalidThreshold, fmt.Sprintf("VotingPeriod ( "+param.Value.Threshold.String()+" ) should be between 0 and 1"))
		}
		if param.Value.GovernancePenalty.LT(sdk.NewRat(0)) || param.Value.GovernancePenalty.GT(sdk.NewRat(1)) {
			return sdk.NewError(parameter.DefaultCodespace, parameter.CodeInvalidGovernancePenalty, fmt.Sprintf("VotingPeriod ( "+param.Value.GovernancePenalty.String()+" ) should be between 0 and 1"))
		}
		if param.Value.Veto.LT(sdk.NewRat(0)) || param.Value.Veto.GT(sdk.NewRat(1)) {
			return sdk.NewError(parameter.DefaultCodespace, parameter.CodeInvalidVeto, fmt.Sprintf("VotingPeriod ( "+param.Value.Veto.String()+" ) should be between 0 and 1"))
		}

		return nil

	}
	return sdk.NewError(parameter.DefaultCodespace, parameter.CodeInvalidTallyingProcedure, fmt.Sprintf("Json is not valid"))
}



