package gov

import (
	"fmt"

	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"strconv"
	"time"
)

const (
	CRITICAL_DEPOSIT      = 4000
	IMPORTANT_DEPOSIT     = 2000
	NORMAL_DEPOSIT        = 1000
	CRITICAL              = "Critical"
	IMPORTANT             = "Important"
	NORMAL                = "Normal"
	LOWER_BOUND_AMOUNT    = 10
	UPPER_BOUND_AMOUNT    = 10000
	STABLE_CRITIACAL_NUM  = 1
	DEFAULT_IMPORTANT_NUM = 5
	DEFAULT_NORMAL_NUM    = 7
	MIN_IMPORTANT_NUM     = 1
	MIN_NORMAL_NUM        = 1
)

var _ params.ParamSet = (*GovParams)(nil)

// default paramspace for params keeper
const (
	DefaultParamSpace = "gov"
)

//Parameter store key
var (
	KeyCriticalDepositPeriod = []byte(CRITICAL + "DepositPeriod")
	KeyCriticalMinDeposit    = []byte(CRITICAL + "MinDeposit")
	KeyCriticalVotingPeriod  = []byte(CRITICAL + "VotingPeriod")
	KeyCriticalMaxNum        = []byte(CRITICAL + "MaxNum")
	KeyCriticalThreshold     = []byte(CRITICAL + "Threshold")
	KeyCriticalVeto          = []byte(CRITICAL + "Veto")
	KeyCriticalParticipation = []byte(CRITICAL + "Participation")
	KeyCriticalPenalty       = []byte(CRITICAL + "Penalty")

	KeyImportantDepositPeriod = []byte(IMPORTANT + "DepositPeriod")
	KeyImportantMinDeposit    = []byte(IMPORTANT + "MinDeposit")
	KeyImportantVotingPeriod  = []byte(IMPORTANT + "VotingPeriod")
	KeyImportantMaxNum        = []byte(IMPORTANT + "MaxNum")
	KeyImportantThreshold     = []byte(IMPORTANT + "Threshold")
	KeyImportantVeto          = []byte(IMPORTANT + "Veto")
	KeyImportantParticipation = []byte(IMPORTANT + "Participation")
	KeyImportantPenalty       = []byte(IMPORTANT + "Penalty")

	KeyNormalDepositPeriod = []byte(NORMAL + "DepositPeriod")
	KeyNormalMinDeposit    = []byte(NORMAL + "MinDeposit")
	KeyNormalVotingPeriod  = []byte(NORMAL + "VotingPeriod")
	KeyNormalMaxNum        = []byte(NORMAL + "MaxNum")
	KeyNormalThreshold     = []byte(NORMAL + "Threshold")
	KeyNormalVeto          = []byte(NORMAL + "Veto")
	KeyNormalParticipation = []byte(NORMAL + "Participation")
	KeyNormalPenalty       = []byte(NORMAL + "Penalty")

	KeySystemHaltPeriod = []byte("SystemHaltPeriod")
)

// ParamTable for mint module
func ParamTypeTable() params.TypeTable {
	return params.NewTypeTable().RegisterParamSet(&GovParams{})
}

// mint parameters
type GovParams struct {
	CriticalDepositPeriod time.Duration `json:"critical_deposit_period"` //  Maximum period for Atom holders to deposit on a proposal. Initial value: 2 months
	CriticalMinDeposit    sdk.Coins     `json:"critical_min_deposit"`    //  Minimum deposit for a critical proposal to enter voting period.
	CriticalVotingPeriod  time.Duration `json:"critical_voting_period"`  //  Length of the critical voting period.
	CriticalMaxNum        uint64        `json:"critical_max_num"`
	CriticalThreshold     sdk.Dec       `json:"critical_threshold"`     //  Minimum propotion of Yes votes for proposal to pass. Initial value: 0.5
	CriticalVeto          sdk.Dec       `json:"critical_veto"`          //  Minimum value of Veto votes to Total votes ratio for proposal to be vetoed. Initial value: 1/3
	CriticalParticipation sdk.Dec       `json:"critical_participation"` //
	CriticalPenalty       sdk.Dec       `json:"critical_penalty"`       //  Penalty if validator does not vote

	ImportantDepositPeriod time.Duration `json:"important_deposit_period"` //  Maximum period for Atom holders to deposit on a proposal. Initial value: 2 months
	ImportantMinDeposit    sdk.Coins     `json:"important_min_deposit"`    //  Minimum deposit for a important proposal to enter voting period.
	ImportantVotingPeriod  time.Duration `json:"important_voting_period"`  //  Length of the important voting period.
	ImportantMaxNum        uint64        `json:"important_max_num"`
	ImportantThreshold     sdk.Dec       `json:"important_threshold"`     //  Minimum propotion of Yes votes for proposal to pass. Initial value: 0.5
	ImportantVeto          sdk.Dec       `json:"important_veto"`          //  Minimum value of Veto votes to Total votes ratio for proposal to be vetoed. Initial value: 1/3
	ImportantParticipation sdk.Dec       `json:"important_participation"` //
	ImportantPenalty       sdk.Dec       `json:"important_penalty"`       //  Penalty if validator does not vote

	NormalDepositPeriod time.Duration `json:"normal_deposit_period"` //  Maximum period for Atom holders to deposit on a proposal. Initial value: 2 months
	NormalMinDeposit    sdk.Coins     `json:"normal_min_deposit"`    //  Minimum deposit for a normal proposal to enter voting period.
	NormalVotingPeriod  time.Duration `json:"normal_voting_period"`  //  Length of the normal voting period.
	NormalMaxNum        uint64        `json:"normal_max_num"`
	NormalThreshold     sdk.Dec       `json:"normal_threshold"`     //  Minimum propotion of Yes votes for proposal to pass. Initial value: 0.5
	NormalVeto          sdk.Dec       `json:"normal_veto"`          //  Minimum value of Veto votes to Total votes ratio for proposal to be vetoed. Initial value: 1/3
	NormalParticipation sdk.Dec       `json:"normal_participation"` //
	NormalPenalty       sdk.Dec       `json:"normal_penalty"`       //  Penalty if validator does not vote

	SystemHaltPeriod int64 `json:"system_halt_period"`
}

func (p GovParams) String() string {
	return fmt.Sprintf(`Gov Params:
System Halt Period:     %v
Proposal Parameter:    [Critical]         [Important]        [Normal]
  DepositPeriod:        %v         %v        %v
  MinDeposit:           %s         %s        %s
  Voting Period:        %v         %v        %v
  Max Num:              %v         %v        %v
  Threshold:            %s         %s        %s
  Veto:                 %s         %s        %s
  Participation:        %s         %s        %s
  Penalty:              %s         %s        %s
`, p.SystemHaltPeriod,
		p.CriticalDepositPeriod, p.ImportantDepositPeriod, p.NormalDepositPeriod,
		p.CriticalMinDeposit.String(), p.ImportantMinDeposit.String(), p.NormalMinDeposit.String(),
		p.CriticalVotingPeriod, p.ImportantVotingPeriod, p.NormalVotingPeriod,
		p.CriticalMaxNum, p.ImportantMaxNum, p.NormalMaxNum,
		p.CriticalThreshold.String(), p.ImportantThreshold.String(), p.NormalThreshold.String(),
		p.CriticalVeto.String(), p.ImportantVeto.String(), p.NormalVeto.String(),
		p.CriticalParticipation.String(), p.ImportantParticipation.String(), p.NormalParticipation.String(),
		p.CriticalPenalty.String(), p.ImportantPenalty.String(), p.NormalPenalty.String())
}

// Implements params.ParamStruct
func (p *GovParams) GetParamSpace() string {
	return DefaultParamSpace
}

func (p *GovParams) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{
		{KeyCriticalDepositPeriod, &p.CriticalDepositPeriod},
		{KeyCriticalMinDeposit, &p.CriticalMinDeposit},
		{KeyCriticalVotingPeriod, &p.CriticalVotingPeriod},
		{KeyCriticalMaxNum, &p.CriticalMaxNum},
		{KeyCriticalThreshold, &p.CriticalThreshold},
		{KeyCriticalVeto, &p.CriticalVeto},
		{KeyCriticalParticipation, &p.CriticalParticipation},
		{KeyCriticalPenalty, &p.CriticalPenalty},

		{KeyImportantDepositPeriod, &p.ImportantDepositPeriod},
		{KeyImportantMinDeposit, &p.ImportantMinDeposit},
		{KeyImportantVotingPeriod, &p.ImportantVotingPeriod},
		{KeyImportantMaxNum, &p.ImportantMaxNum},
		{KeyImportantThreshold, &p.ImportantThreshold},
		{KeyImportantVeto, &p.ImportantVeto},
		{KeyImportantParticipation, &p.ImportantParticipation},
		{KeyImportantPenalty, &p.ImportantPenalty},

		{KeyNormalDepositPeriod, &p.NormalDepositPeriod},
		{KeyNormalMinDeposit, &p.NormalMinDeposit},
		{KeyNormalVotingPeriod, &p.NormalVotingPeriod},
		{KeyNormalMaxNum, &p.NormalMaxNum},
		{KeyNormalThreshold, &p.NormalThreshold},
		{KeyNormalVeto, &p.NormalVeto},
		{KeyNormalParticipation, &p.NormalParticipation},
		{KeyNormalPenalty, &p.NormalPenalty},

		{KeySystemHaltPeriod, &p.SystemHaltPeriod},
	}
}

func (p *GovParams) Validate(key string, value string) (interface{}, sdk.Error) {
	return nil, nil
}

func (p *GovParams) StringFromBytes(cdc *codec.Codec, key string, bytes []byte) (string, error) {
	switch key {
	case string(KeyCriticalDepositPeriod):
		err := cdc.UnmarshalJSON(bytes, &p.CriticalDepositPeriod)
		return p.CriticalDepositPeriod.String(), err
	case string(KeyCriticalMinDeposit):
		err := cdc.UnmarshalJSON(bytes, &p.CriticalMinDeposit)
		return p.CriticalMinDeposit.String(), err
	case string(KeyCriticalVotingPeriod):
		err := cdc.UnmarshalJSON(bytes, &p.CriticalVotingPeriod)
		return p.CriticalDepositPeriod.String(), err
	case string(KeyCriticalMaxNum):
		err := cdc.UnmarshalJSON(bytes, &p.CriticalMaxNum)
		return strconv.FormatUint(p.CriticalMaxNum, 10), err
	case string(KeyCriticalThreshold):
		err := cdc.UnmarshalJSON(bytes, &p.CriticalThreshold)
		return p.CriticalThreshold.String(), err
	case string(KeyCriticalVeto):
		err := cdc.UnmarshalJSON(bytes, &p.CriticalVeto)
		return p.CriticalThreshold.String(), err
	case string(KeyCriticalParticipation):
		err := cdc.UnmarshalJSON(bytes, &p.CriticalParticipation)
		return p.CriticalParticipation.String(), err
	case string(KeyCriticalPenalty):
		err := cdc.UnmarshalJSON(bytes, &p.CriticalPenalty)
		return p.CriticalPenalty.String(), err

	case string(KeyImportantDepositPeriod):
		err := cdc.UnmarshalJSON(bytes, &p.ImportantDepositPeriod)
		return p.ImportantDepositPeriod.String(), err
	case string(KeyImportantMinDeposit):
		err := cdc.UnmarshalJSON(bytes, &p.ImportantMinDeposit)
		return p.ImportantMinDeposit.String(), err
	case string(KeyImportantVotingPeriod):
		err := cdc.UnmarshalJSON(bytes, &p.ImportantVotingPeriod)
		return p.ImportantDepositPeriod.String(), err
	case string(KeyImportantMaxNum):
		err := cdc.UnmarshalJSON(bytes, &p.ImportantMaxNum)
		return strconv.FormatUint(p.ImportantMaxNum, 10), err
	case string(KeyImportantThreshold):
		err := cdc.UnmarshalJSON(bytes, &p.ImportantThreshold)
		return p.ImportantThreshold.String(), err
	case string(KeyImportantVeto):
		err := cdc.UnmarshalJSON(bytes, &p.ImportantVeto)
		return p.ImportantThreshold.String(), err
	case string(KeyImportantParticipation):
		err := cdc.UnmarshalJSON(bytes, &p.ImportantParticipation)
		return p.ImportantParticipation.String(), err
	case string(KeyImportantPenalty):
		err := cdc.UnmarshalJSON(bytes, &p.ImportantPenalty)
		return p.ImportantPenalty.String(), err

	case string(KeyNormalDepositPeriod):
		err := cdc.UnmarshalJSON(bytes, &p.NormalDepositPeriod)
		return p.NormalDepositPeriod.String(), err
	case string(KeyNormalMinDeposit):
		err := cdc.UnmarshalJSON(bytes, &p.NormalMinDeposit)
		return p.NormalMinDeposit.String(), err
	case string(KeyNormalVotingPeriod):
		err := cdc.UnmarshalJSON(bytes, &p.NormalVotingPeriod)
		return p.NormalDepositPeriod.String(), err
	case string(KeyNormalMaxNum):
		err := cdc.UnmarshalJSON(bytes, &p.NormalMaxNum)
		return strconv.FormatUint(p.NormalMaxNum, 10), err
	case string(KeyNormalThreshold):
		err := cdc.UnmarshalJSON(bytes, &p.NormalThreshold)
		return p.NormalThreshold.String(), err
	case string(KeyNormalVeto):
		err := cdc.UnmarshalJSON(bytes, &p.NormalVeto)
		return p.NormalThreshold.String(), err
	case string(KeyNormalParticipation):
		err := cdc.UnmarshalJSON(bytes, &p.NormalParticipation)
		return p.NormalParticipation.String(), err
	case string(KeyNormalPenalty):
		err := cdc.UnmarshalJSON(bytes, &p.NormalPenalty)
		return p.NormalPenalty.String(), err

	case string(KeySystemHaltPeriod):
		err := cdc.UnmarshalJSON(bytes, &p.SystemHaltPeriod)
		return strconv.FormatInt(p.SystemHaltPeriod, 10), err
	default:
		return "", fmt.Errorf("%s is not existed", key)
	}
}

func (p *GovParams) ReadOnly() bool {
	return true
}

// default minting module parameters
func DefaultParams() GovParams {
	var criticalMinDeposit, _ = sdk.IrisCoinType.ConvertToMinDenomCoin(fmt.Sprintf("%d%s", CRITICAL_DEPOSIT, sdk.Iris))
	var importantMinDeposit, _ = sdk.IrisCoinType.ConvertToMinDenomCoin(fmt.Sprintf("%d%s", IMPORTANT_DEPOSIT, sdk.Iris))
	var normalMinDeposit, _ = sdk.IrisCoinType.ConvertToMinDenomCoin(fmt.Sprintf("%d%s", NORMAL_DEPOSIT, sdk.Iris))

	if sdk.NetworkType == sdk.Mainnet {
		return GovParams{
			CriticalDepositPeriod: time.Duration(sdk.Day),
			CriticalMinDeposit:    sdk.Coins{criticalMinDeposit},
			CriticalVotingPeriod:  time.Duration(sdk.FiveDays),
			CriticalMaxNum:        STABLE_CRITIACAL_NUM,
			CriticalThreshold:     sdk.NewDecWithPrec(75, 2),
			CriticalVeto:          sdk.NewDecWithPrec(33, 2),
			CriticalParticipation: sdk.NewDecWithPrec(50, 2),
			CriticalPenalty:       sdk.ZeroDec(),

			ImportantDepositPeriod: time.Duration(sdk.Day),
			ImportantMinDeposit:    sdk.Coins{importantMinDeposit},
			ImportantVotingPeriod:  time.Duration(sdk.FiveDays),
			ImportantMaxNum:        DEFAULT_IMPORTANT_NUM,
			ImportantThreshold:     sdk.NewDecWithPrec(67, 2),
			ImportantVeto:          sdk.NewDecWithPrec(33, 2),
			ImportantParticipation: sdk.NewDecWithPrec(50, 2),
			ImportantPenalty:       sdk.ZeroDec(),

			NormalDepositPeriod: time.Duration(sdk.Day),
			NormalMinDeposit:    sdk.Coins{normalMinDeposit},
			NormalVotingPeriod:  time.Duration(sdk.FiveDays),
			NormalMaxNum:        DEFAULT_NORMAL_NUM,
			NormalThreshold:     sdk.NewDecWithPrec(50, 2),
			NormalVeto:          sdk.NewDecWithPrec(33, 2),
			NormalParticipation: sdk.NewDecWithPrec(50, 2),
			NormalPenalty:       sdk.ZeroDec(),
			SystemHaltPeriod:    20000,
		}
	} else {
		return GovParams{
			CriticalDepositPeriod: time.Duration(sdk.Day),
			CriticalMinDeposit:    sdk.Coins{criticalMinDeposit},
			CriticalVotingPeriod:  time.Duration(2 * time.Minute),
			CriticalMaxNum:        STABLE_CRITIACAL_NUM,
			CriticalThreshold:     sdk.NewDecWithPrec(75, 2),
			CriticalVeto:          sdk.NewDecWithPrec(33, 2),
			CriticalParticipation: sdk.NewDecWithPrec(50, 2),
			CriticalPenalty:       sdk.ZeroDec(),

			ImportantDepositPeriod: time.Duration(sdk.Day),
			ImportantMinDeposit:    sdk.Coins{importantMinDeposit},
			ImportantVotingPeriod:  time.Duration(2 * time.Minute),
			ImportantMaxNum:        DEFAULT_IMPORTANT_NUM,
			ImportantThreshold:     sdk.NewDecWithPrec(67, 2),
			ImportantVeto:          sdk.NewDecWithPrec(33, 2),
			ImportantParticipation: sdk.NewDecWithPrec(50, 2),
			ImportantPenalty:       sdk.ZeroDec(),

			NormalDepositPeriod: time.Duration(sdk.Day),
			NormalMinDeposit:    sdk.Coins{normalMinDeposit},
			NormalVotingPeriod:  time.Duration(2 * time.Minute),
			NormalMaxNum:        DEFAULT_NORMAL_NUM,
			NormalThreshold:     sdk.NewDecWithPrec(50, 2),
			NormalVeto:          sdk.NewDecWithPrec(33, 2),
			NormalParticipation: sdk.NewDecWithPrec(50, 2),
			NormalPenalty:       sdk.ZeroDec(),
			SystemHaltPeriod:    60,
		}
	}
}

func DefaultParamsForTest() GovParams {
	var criticalMinDeposit, _ = sdk.IrisCoinType.ConvertToMinDenomCoin(fmt.Sprintf("%d%s", 10, sdk.Iris))
	var importantMinDeposit, _ = sdk.IrisCoinType.ConvertToMinDenomCoin(fmt.Sprintf("%d%s", 10, sdk.Iris))
	var normalMinDeposit, _ = sdk.IrisCoinType.ConvertToMinDenomCoin(fmt.Sprintf("%d%s", 10, sdk.Iris))

	return GovParams{
		CriticalDepositPeriod: time.Duration(30 * time.Second),
		CriticalMinDeposit:    sdk.Coins{criticalMinDeposit},
		CriticalVotingPeriod:  time.Duration(30 * time.Second),
		CriticalMaxNum:        STABLE_CRITIACAL_NUM,
		CriticalThreshold:     sdk.NewDecWithPrec(857, 3),
		CriticalVeto:          sdk.NewDecWithPrec(334, 3),
		CriticalParticipation: sdk.NewDecWithPrec(875, 3),
		CriticalPenalty:       sdk.ZeroDec(),

		ImportantDepositPeriod: time.Duration(30 * time.Second),
		ImportantMinDeposit:    sdk.Coins{importantMinDeposit},
		ImportantVotingPeriod:  time.Duration(30 * time.Second),
		ImportantMaxNum:        DEFAULT_IMPORTANT_NUM,
		ImportantThreshold:     sdk.NewDecWithPrec(8, 1),
		ImportantVeto:          sdk.NewDecWithPrec(334, 3),
		ImportantParticipation: sdk.NewDecWithPrec(834, 3),
		ImportantPenalty:       sdk.ZeroDec(),

		NormalDepositPeriod: time.Duration(30 * time.Second),
		NormalMinDeposit:    sdk.Coins{normalMinDeposit},
		NormalVotingPeriod:  time.Duration(30 * time.Second),
		NormalMaxNum:        DEFAULT_NORMAL_NUM,
		NormalThreshold:     sdk.NewDecWithPrec(667, 3),
		NormalVeto:          sdk.NewDecWithPrec(334, 3),
		NormalParticipation: sdk.NewDecWithPrec(75, 2),
		NormalPenalty:       sdk.ZeroDec(),
		SystemHaltPeriod:    60,
	}
}

func validateParams(p GovParams) sdk.Error {
	if err := validateDepositProcedure(DepositProcedure{
		MaxDepositPeriod: p.CriticalDepositPeriod,
		MinDeposit:       p.CriticalMinDeposit,
	}, CRITICAL); err != nil {
		return err
	}

	if err := validatorVotingProcedure(VotingProcedure{
		VotingPeriod: p.CriticalVotingPeriod,
	}, CRITICAL); err != nil {
		return err
	}

	if err := validateTallyingProcedure(TallyingProcedure{
		Threshold:     p.CriticalThreshold,
		Veto:          p.CriticalVeto,
		Participation: p.CriticalParticipation,
		Penalty:       p.CriticalPenalty,
	}, CRITICAL); err != nil {
		return err
	}

	if err := validateDepositProcedure(DepositProcedure{
		MaxDepositPeriod: p.ImportantDepositPeriod,
		MinDeposit:       p.ImportantMinDeposit,
	}, IMPORTANT); err != nil {
		return err
	}

	if err := validatorVotingProcedure(VotingProcedure{
		VotingPeriod: p.ImportantVotingPeriod,
	}, IMPORTANT); err != nil {
		return err
	}

	if err := validateTallyingProcedure(TallyingProcedure{
		Threshold:     p.ImportantThreshold,
		Veto:          p.ImportantVeto,
		Participation: p.ImportantParticipation,
		Penalty:       p.ImportantPenalty,
	}, IMPORTANT); err != nil {
		return err
	}

	if err := validateDepositProcedure(DepositProcedure{
		MaxDepositPeriod: p.NormalDepositPeriod,
		MinDeposit:       p.NormalMinDeposit,
	}, NORMAL); err != nil {
		return err
	}

	if err := validatorVotingProcedure(VotingProcedure{
		VotingPeriod: p.NormalVotingPeriod,
	}, NORMAL); err != nil {
		return err
	}

	if err := validateTallyingProcedure(TallyingProcedure{
		Threshold:     p.NormalThreshold,
		Veto:          p.NormalVeto,
		Participation: p.NormalParticipation,
		Penalty:       p.NormalPenalty,
	}, NORMAL); err != nil {
		return err
	}

	if err := validateMaxNum(p); err != nil {
		return err
	}

	if p.SystemHaltPeriod < 0 || p.SystemHaltPeriod > 50000 {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidSystemHaltPeriod, fmt.Sprintf("SystemHaltPeriod should be between [0, 50000]"))
	}

	return nil
}

//______________________________________________________________________

type DepositProcedure struct {
	MinDeposit       sdk.Coins
	MaxDepositPeriod time.Duration
}

type VotingProcedure struct {
	VotingPeriod time.Duration `json:"critical_voting_period"` //  Length of the critical voting period.
}

type TallyingProcedure struct {
	Threshold     sdk.Dec `json:"threshold"`     //  Minimum propotion of Yes votes for proposal to pass. Initial value: 0.5
	Veto          sdk.Dec `json:"veto"`          //  Minimum value of Veto votes to Total votes ratio for proposal to be vetoed. Initial value: 1/3
	Participation sdk.Dec `json:"participation"` //
	Penalty       sdk.Dec `json:"penalty"`       //  Penalty if validator does not vote
}

func validateDepositProcedure(dp DepositProcedure, level string) sdk.Error {
	if dp.MinDeposit[0].Denom != sdk.IrisAtto {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMinDepositDenom, fmt.Sprintf(level+"MinDeposit denom should be %s!", sdk.IrisAtto))
	}

	LowerBound, _ := sdk.IrisCoinType.ConvertToMinDenomCoin(fmt.Sprintf("%d%s", LOWER_BOUND_AMOUNT, sdk.Iris))
	UpperBound, _ := sdk.IrisCoinType.ConvertToMinDenomCoin(fmt.Sprintf("%d%s", UPPER_BOUND_AMOUNT, sdk.Iris))

	if dp.MinDeposit[0].Amount.LT(LowerBound.Amount) || dp.MinDeposit[0].Amount.GT(UpperBound.Amount) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMinDepositAmount, fmt.Sprintf(level+"MinDepositAmount"+dp.MinDeposit[0].String()+" should be larger than 10iris and less than 10000iris"))
	}

	if dp.MaxDepositPeriod < sdk.TwentySeconds || dp.MaxDepositPeriod > sdk.ThreeDays {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidDepositPeriod, fmt.Sprintf(level+"MaxDepositPeriod (%s) should be between 20s and %s", dp.MaxDepositPeriod.String(), sdk.ThreeDays.String()))
	}
	return nil
}

func validatorVotingProcedure(vp VotingProcedure, level string) sdk.Error {
	if vp.VotingPeriod < sdk.TwentySeconds || vp.VotingPeriod > sdk.Week {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidVotingPeriod, fmt.Sprintf(level+"VotingPeriod (%s) should be between 20s and 1 week", vp.VotingPeriod.String()))
	}
	return nil
}

func validateTallyingProcedure(tp TallyingProcedure, level string) sdk.Error {
	if tp.Threshold.LTE(sdk.ZeroDec()) || tp.Threshold.GTE(sdk.NewDec(1)) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidThreshold, fmt.Sprintf("Invalid "+level+" Threshold ( "+tp.Threshold.String()+" ) should be (0,1)"))
	}
	if tp.Participation.LTE(sdk.ZeroDec()) || tp.Participation.GTE(sdk.NewDec(1)) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidParticipation, fmt.Sprintf("Invalid "+level+" participation ( "+tp.Participation.String()+" ) should be (0,1)"))
	}
	if tp.Veto.LTE(sdk.ZeroDec()) || tp.Veto.GTE(sdk.NewDec(1)) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidVeto, fmt.Sprintf("Invalid "+level+" Veto ( "+tp.Veto.String()+" ) should be (0,1)"))
	}
	if tp.Penalty.LT(sdk.ZeroDec()) || tp.Penalty.GTE(sdk.NewDec(1)) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidGovernancePenalty, fmt.Sprintf("Invalid "+level+" GovernancePenalty ( "+tp.Penalty.String()+" ) should be [0,1)"))
	}
	return nil
}

func validateMaxNum(gp GovParams) sdk.Error {
	if gp.CriticalMaxNum != STABLE_CRITIACAL_NUM {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMaxProposalNum, fmt.Sprintf("The num of Max"+CRITICAL+"Proposal [%v] can only be %v.", gp.CriticalMaxNum, STABLE_CRITIACAL_NUM))
	}
	if gp.ImportantMaxNum < MIN_IMPORTANT_NUM {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMaxProposalNum, fmt.Sprintf("The num of Max"+IMPORTANT+"Proposal [%v] should be no less than %v.", gp.CriticalMaxNum, MIN_IMPORTANT_NUM))
	}
	if gp.NormalMaxNum < MIN_NORMAL_NUM {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMaxProposalNum, fmt.Sprintf("The num of Max"+NORMAL+"Proposal [%v] should be no less than %v.", gp.NormalMaxNum, MIN_NORMAL_NUM))
	}
	return nil
}
