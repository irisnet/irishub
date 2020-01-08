package gov

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	Params Params `json:"params"` // inflation params
}

type Params struct {
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
