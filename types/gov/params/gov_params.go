package govparams

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/params"
	stakeTypes "github.com/irisnet/irishub/modules/stake/types"
	"github.com/irisnet/irishub/types"
	sdk "github.com/irisnet/irishub/types"
)

var DepositProcedureParameter DepositProcedureParam

const (
	CRITICAL_DEPOSIT     = 4000
	IMPORTANT_DEPOSIT    = 2000
	NORMAL_DEPOSIT       = 1000
	CRITICAL             = "Critical"
	IMPORTANT            = "Important"
	NORMAL               = "Normal"
	LOWER_BOUND_AMOUNT   = 10
	UPPER_BOUND_AMOUNT   = 10000
	THREE_DAYS           = 3 * 3600 * 24
	SIXTY_HOURS          = 60 * 3600
	TWO_DAYS             = 2 * 3600 * 24 //
	ONE_DAY              = 1 * 3600 * 24
	STABLE_CRITIACAL_NUM = 1
	MIN_IMPORTANT_NUM    = 2
	MIN_NORMAL_NUM       = 1
)

var _ params.GovParameter = (*DepositProcedureParam)(nil)

type Params struct {
	DepositProcedure  DepositProcedure  `json:"Gov/govDepositProcedure"`
	VotingProcedure   VotingProcedure   `json:"Gov/govVotingProcedure"`
	TallyingProcedure TallyingProcedure `json:"Gov/govTallyingProcedure"`
}

// Procedure around Deposits for governance
type DepositProcedure struct {
	CriticalMinDeposit  sdk.Coins     `json:"critical_min_deposit"`  //  Minimum deposit for a critical proposal to enter voting period.
	ImportantMinDeposit sdk.Coins     `json:"important_min_deposit"` //  Minimum deposit for a important proposal to enter voting period.
	NormalMinDeposit    sdk.Coins     `json:"normal_min_deposit"`    //  Minimum deposit for a normal proposal to enter voting period.
	MaxDepositPeriod    time.Duration `json:"max_deposit_period"`    //  Maximum period for Atom holders to deposit on a proposal. Initial value: 2 months
}

type DepositProcedureParam struct {
	Value      DepositProcedure
	paramSpace params.Subspace
}

func NewDepositProcedure() DepositProcedure {
	var ciriticalMinDeposit, _ = types.NewDefaultCoinType(stakeTypes.StakeDenomName).ConvertToMinCoin(fmt.Sprintf("%d%s", CRITICAL_DEPOSIT, stakeTypes.StakeDenomName))
	var importantMinDeposit, _ = types.NewDefaultCoinType(stakeTypes.StakeDenomName).ConvertToMinCoin(fmt.Sprintf("%d%s", IMPORTANT_DEPOSIT, stakeTypes.StakeDenomName))
	var normalMinDeposit, _ = types.NewDefaultCoinType(stakeTypes.StakeDenomName).ConvertToMinCoin(fmt.Sprintf("%d%s", NORMAL_DEPOSIT, stakeTypes.StakeDenomName))

	return DepositProcedure{
		CriticalMinDeposit:  sdk.Coins{ciriticalMinDeposit},
		ImportantMinDeposit: sdk.Coins{importantMinDeposit},
		NormalMinDeposit:    sdk.Coins{normalMinDeposit},
		MaxDepositPeriod:    time.Duration(ONE_DAY) * time.Second}
}

func (param *DepositProcedureParam) GetValueFromRawData(cdc *codec.Codec, res []byte) interface{} {
	cdc.UnmarshalJSON(res, &param.Value)
	return param.Value
}

func (param *DepositProcedureParam) InitGenesis(genesisState interface{}) {
	if value, ok := genesisState.(DepositProcedure); ok {
		param.Value = value
	} else {
		panic("The govDepositProcedure in GenesisState is empty. ")
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
		jsonBytes, _ = json.Marshal(param.Value)
		return string(jsonBytes)
	}

	if err := json.Unmarshal([]byte(jsonStr), &param.Value); err == nil {
		jsonBytes, _ = json.Marshal(param.Value)
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

		if param.Value.CriticalMinDeposit[0].Denom != stakeTypes.StakeDenom ||
			param.Value.ImportantMinDeposit[0].Denom != stakeTypes.StakeDenom ||
			param.Value.NormalMinDeposit[0].Denom != stakeTypes.StakeDenom {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMinDepositDenom, fmt.Sprintf("It should be %s!", stakeTypes.StakeDenom))
		}

		LowerBound, _ := types.NewDefaultCoinType(stakeTypes.StakeDenomName).ConvertToMinCoin(fmt.Sprintf("%d%s", LOWER_BOUND_AMOUNT, stakeTypes.StakeDenomName))
		UpperBound, _ := types.NewDefaultCoinType(stakeTypes.StakeDenomName).ConvertToMinCoin(fmt.Sprintf("%d%s", UPPER_BOUND_AMOUNT, stakeTypes.StakeDenomName))

		if param.Value.CriticalMinDeposit[0].Amount.LT(LowerBound.Amount) || param.Value.CriticalMinDeposit[0].Amount.GT(UpperBound.Amount) {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMinDepositAmount, fmt.Sprintf(CRITICAL+"MinDepositAmount"+param.Value.CriticalMinDeposit[0].String()+" should be larger than 10iris and less than 10000iris"))
		}

		if param.Value.ImportantMinDeposit[0].Amount.LT(LowerBound.Amount) || param.Value.ImportantMinDeposit[0].Amount.GT(UpperBound.Amount) {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMinDepositAmount, fmt.Sprintf(IMPORTANT+"MinDepositAmount"+param.Value.CriticalMinDeposit[0].String()+" should be larger than 10iris and less than 10000iris"))
		}

		if param.Value.NormalMinDeposit[0].Amount.LT(LowerBound.Amount) || param.Value.NormalMinDeposit[0].Amount.GT(UpperBound.Amount) {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMinDepositAmount, fmt.Sprintf(NORMAL+"MinDepositAmount"+param.Value.CriticalMinDeposit[0].String()+" should be larger than 10iris and less than 10000iris"))
		}

		if param.Value.MaxDepositPeriod.Seconds() < 20 || param.Value.MaxDepositPeriod.Seconds() > THREE_DAYS {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidDepositPeriod, fmt.Sprintf("MaxDepositPeriod (%s) should be between 20s and %ds", strconv.Itoa(int(param.Value.MaxDepositPeriod.Seconds())), THREE_DAYS))
		}

		return nil

	}
	return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMinDeposit, fmt.Sprintf("Json is not valid"))
}

var VotingProcedureParameter VotingProcedureParam
var _ params.GovParameter = (*VotingProcedureParam)(nil)

// Procedure around Voting in governance
type VotingProcedure struct {
	CriticalVotingPeriod  time.Duration `json:"critical_voting_period"`  //  Length of the critical voting period.
	ImportantVotingPeriod time.Duration `json:"important_voting_period"` //  Length of the important voting period.
	NormalVotingPeriod    time.Duration `json:"normal_voting_period"`    //  Length of the normal voting period.
	CriticalMaxNum        uint64        `json:"critical_max_num"`
	ImportantMaxNum       uint64        `json:"important_max_num"`
	NormalMaxNum          uint64        `json:"normal_max_num"`
}

type VotingProcedureParam struct {
	Value      VotingProcedure
	paramSpace params.Subspace
}

func NewVotingProcedure() VotingProcedure {
	return VotingProcedure{
		CriticalVotingPeriod:  time.Duration(THREE_DAYS) * time.Second,
		ImportantVotingPeriod: time.Duration(SIXTY_HOURS) * time.Second,
		NormalVotingPeriod:    time.Duration(TWO_DAYS) * time.Second,
		CriticalMaxNum:        STABLE_CRITIACAL_NUM,
		ImportantMaxNum:       MIN_IMPORTANT_NUM,
		NormalMaxNum:          MIN_NORMAL_NUM,
	}
}

func (param *VotingProcedureParam) GetValueFromRawData(cdc *codec.Codec, res []byte) interface{} {
	cdc.UnmarshalJSON(res, &param.Value)
	return param.Value
}

func (param *VotingProcedureParam) InitGenesis(genesisState interface{}) {
	if value, ok := genesisState.(VotingProcedure); ok {
		param.Value = value
	} else {
		panic("The govVotingProcedure in GenesisState is empty. ")
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
		jsonBytes, _ = json.Marshal(param.Value)
		return string(jsonBytes)
	}

	if err := json.Unmarshal([]byte(jsonStr), &param.Value); err == nil {
		jsonBytes, _ = json.Marshal(param.Value)
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

		if param.Value.CriticalVotingPeriod.Seconds() < 20 || param.Value.CriticalVotingPeriod.Seconds() > THREE_DAYS {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidVotingPeriod, fmt.Sprintf(CRITICAL+"VotingPeriod (%s) should be between 20s and %ds", strconv.Itoa(int(param.Value.CriticalVotingPeriod.Seconds())), THREE_DAYS))
		}

		if param.Value.ImportantVotingPeriod.Seconds() < 20 || param.Value.ImportantVotingPeriod.Seconds() > THREE_DAYS {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidVotingPeriod, fmt.Sprintf(IMPORTANT+"VotingPeriod (%s) should be between 20s and %ds", strconv.Itoa(int(param.Value.ImportantVotingPeriod.Seconds())), THREE_DAYS))
		}

		if param.Value.NormalVotingPeriod.Seconds() < 20 || param.Value.NormalVotingPeriod.Seconds() > THREE_DAYS {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidVotingPeriod, fmt.Sprintf(NORMAL+"VotingPeriod (%s) should be between 20s and %ds", strconv.Itoa(int(param.Value.NormalVotingPeriod.Seconds())), THREE_DAYS))
		}

		if param.Value.CriticalMaxNum != STABLE_CRITIACAL_NUM {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMaxProposalNum, fmt.Sprintf("The num of Max"+CRITICAL+"Proposal  only can be %v.", STABLE_CRITIACAL_NUM))
		}

		if param.Value.ImportantMaxNum < MIN_IMPORTANT_NUM {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMaxProposalNum, fmt.Sprintf("The num of Max"+IMPORTANT+"Proposal must be more than %v.", MIN_IMPORTANT_NUM))
		}

		if param.Value.NormalMaxNum < MIN_NORMAL_NUM {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMaxProposalNum, fmt.Sprintf("The num of Max"+NORMAL+"Proposal should be no less than %v.", MIN_NORMAL_NUM))
		}
		return nil

	}
	return sdk.NewError(params.DefaultCodespace, params.CodeInvalidVotingProcedure, fmt.Sprintf("Json is not valid"))
}

var TallyingProcedureParameter TallyingProcedureParam
var _ params.GovParameter = (*TallyingProcedureParam)(nil)

type TallyCondition struct {
	Threshold     sdk.Dec `json:"threshold"`     //  Minimum propotion of Yes votes for proposal to pass. Initial value: 0.5
	Veto          sdk.Dec `json:"veto"`          //  Minimum value of Veto votes to Total votes ratio for proposal to be vetoed. Initial value: 1/3
	Participation sdk.Dec `json:"participation"` //
	Penalty       sdk.Dec `json:"penalty"`       //  Penalty if validator does not vote
}

// Procedure around Tallying votes in governance
type TallyingProcedure struct {
	CriticalCondition  TallyCondition `json:"critical_condition"`
	ImportantCondition TallyCondition `json:"important_condition"`
	NormalCondition    TallyCondition `json:"normal_condition"`
}

type TallyingProcedureParam struct {
	Value      TallyingProcedure
	paramSpace params.Subspace
}

func NewTallyingProcedure() TallyingProcedure {
	return TallyingProcedure{
		CriticalCondition: TallyCondition{
			Threshold:     sdk.NewDecWithPrec(834, 3),
			Veto:          sdk.NewDecWithPrec(334, 3),
			Participation: sdk.NewDecWithPrec(8572, 4),
			Penalty:       sdk.NewDecWithPrec(9, 4),
		},
		ImportantCondition: TallyCondition{
			Threshold:     sdk.NewDecWithPrec(8, 1),
			Veto:          sdk.NewDecWithPrec(334, 3),
			Participation: sdk.NewDecWithPrec(834, 3),
			Penalty:       sdk.NewDecWithPrec(7, 4),
		},
		NormalCondition: TallyCondition{
			Threshold:     sdk.NewDecWithPrec(667, 3),
			Veto:          sdk.NewDecWithPrec(334, 3),
			Participation: sdk.NewDecWithPrec(75, 2),
			Penalty:       sdk.NewDecWithPrec(5, 4),
		},
	}
}

func (param *TallyingProcedureParam) GetValueFromRawData(cdc *codec.Codec, res []byte) interface{} {
	cdc.UnmarshalJSON(res, &param.Value)
	return param.Value
}

func (param *TallyingProcedureParam) InitGenesis(genesisState interface{}) {
	if value, ok := genesisState.(TallyingProcedure); ok {
		param.Value = value
	} else {
		panic("The govTallyingProcedure in GenesisState is empty. ")
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
		jsonBytes, _ = json.Marshal(param.Value)
		return string(jsonBytes)
	}

	if err := json.Unmarshal([]byte(jsonStr), &param.Value); err == nil {
		jsonBytes, _ = json.Marshal(param.Value)
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

		if err := validTallyCondition(param.Value.CriticalCondition, CRITICAL); err != nil {
			return err
		}
		if err := validTallyCondition(param.Value.ImportantCondition, IMPORTANT); err != nil {
			return err
		}
		if err := validTallyCondition(param.Value.NormalCondition, NORMAL); err != nil {
			return err
		}

		return nil

	}
	return sdk.NewError(params.DefaultCodespace, params.CodeInvalidTallyingProcedure, fmt.Sprintf("Json is not valid"))
}

func validTallyCondition(tc TallyCondition, str string) sdk.Error {
	if tc.Threshold.LTE(sdk.ZeroDec()) || tc.Threshold.GTE(sdk.NewDec(1)) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidThreshold, fmt.Sprintf("Invalid "+str+" Threshold ( "+tc.Threshold.String()+" ) should be between 0 and 1"))
	}
	if tc.Participation.LTE(sdk.ZeroDec()) || tc.Participation.GTE(sdk.NewDec(1)) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidParticipation, fmt.Sprintf("Invalid "+str+" participation ( "+tc.Participation.String()+" ) should be between 0 and 1"))
	}
	if tc.Veto.LTE(sdk.ZeroDec()) || tc.Veto.GTE(sdk.NewDec(1)) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidVeto, fmt.Sprintf("Invalid "+str+" Veto ( "+tc.Veto.String()+" ) should be between 0 and 1"))
	}
	if tc.Penalty.LTE(sdk.ZeroDec()) || tc.Penalty.GTE(sdk.NewDec(1)) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidGovernancePenalty, fmt.Sprintf("Invalid "+str+" GovernancePenalty ( "+tc.Penalty.String()+" ) should be between 0 and 1"))
	}
	return nil
}
