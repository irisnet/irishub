package gov

import (
	sdk "github.com/irisnet/irishub/types"
)

//-----------------------------------------------------------
// ProposalLevel

// Type that represents Proposal Level as a byte
type ProposalLevel byte

//nolint
const (
	ProposalLevelNil       ProposalLevel = 0x00
	ProposalLevelCritical  ProposalLevel = 0x01
	ProposalLevelImportant ProposalLevel = 0x02
	ProposalLevelNormal    ProposalLevel = 0x03
)

func (p ProposalLevel) string() string {
	switch p {
	case ProposalLevelCritical:
		return "critical"
	case ProposalLevelImportant:
		return "important"
	case ProposalLevelNormal:
		return "normal"
	default:
		return " "
	}
}

func GetProposalLevel(p Proposal) ProposalLevel {
	return GetProposalLevelByProposalKind(p.GetProposalType())
}

func GetProposalLevelByProposalKind(p ProposalKind) ProposalLevel {
	switch p {
	case ProposalTypeTxTaxUsage:
		return ProposalLevelNormal
	case ProposalTypeParameterChange:
		return ProposalLevelImportant
	case ProposalTypeSystemHalt:
		return ProposalLevelCritical
	case ProposalTypeSoftwareUpgrade:
		return ProposalLevelCritical
	default:
		return ProposalLevelNil
	}
}

// Returns the current Deposit Procedure from the global param store
func (Keeper Keeper) GetDepositProcedure(ctx sdk.Context, p Proposal) DepositProcedure {
	params := Keeper.GetParamSet(ctx)
	switch GetProposalLevel(p) {
	case ProposalLevelCritical:
		return DepositProcedure{
			MinDeposit:       params.CriticalMinDeposit,
			MaxDepositPeriod: params.CriticalDepositPeriod,
		}
	case ProposalLevelImportant:
		return DepositProcedure{
			MinDeposit:       params.ImportantMinDeposit,
			MaxDepositPeriod: params.ImportantDepositPeriod,
		}
	case ProposalLevelNormal:
		return DepositProcedure{
			MinDeposit:       params.NormalMinDeposit,
			MaxDepositPeriod: params.NormalDepositPeriod,
		}
	default:
		panic("There is no level for this proposal which type is " + p.GetProposalType().String())
	}
}

// Returns the current Voting Procedure from the global param store
func (Keeper Keeper) GetVotingProcedure(ctx sdk.Context, p Proposal) VotingProcedure {
	params := Keeper.GetParamSet(ctx)
	switch GetProposalLevel(p) {
	case ProposalLevelCritical:
		return VotingProcedure{
			VotingPeriod: params.CriticalVotingPeriod,
		}
	case ProposalLevelImportant:
		return VotingProcedure{
			VotingPeriod: params.ImportantVotingPeriod,
		}
	case ProposalLevelNormal:
		return VotingProcedure{
			VotingPeriod: params.NormalVotingPeriod,
		}
	default:
		panic("There is no level for this proposal which type is " + p.GetProposalType().String())
	}
}

func (Keeper Keeper) GetMaxNumByProposalLevel(ctx sdk.Context, pl ProposalLevel) uint64 {
	params := Keeper.GetParamSet(ctx)
	switch pl {
	case ProposalLevelCritical:
		return params.CriticalMaxNum

	case ProposalLevelImportant:
		return params.ImportantMaxNum

	case ProposalLevelNormal:
		return params.NormalMaxNum
	default:
		panic("There is no level for this proposal which type is " + pl.string())
	}
}

// Returns the current Tallying Procedure from the global param store
func (Keeper Keeper) GetTallyingProcedure(ctx sdk.Context, p Proposal) TallyingProcedure {
	params := Keeper.GetParamSet(ctx)
	switch GetProposalLevel(p) {
	case ProposalLevelCritical:
		return TallyingProcedure{
			Threshold:     params.CriticalThreshold,
			Veto:          params.CriticalVeto,
			Participation: params.CriticalParticipation,
			Penalty:       params.CriticalPenalty,
		}
	case ProposalLevelImportant:
		return TallyingProcedure{
			Threshold:     params.ImportantThreshold,
			Veto:          params.ImportantVeto,
			Participation: params.ImportantParticipation,
			Penalty:       params.ImportantPenalty,
		}
	case ProposalLevelNormal:
		return TallyingProcedure{
			Threshold:     params.NormalThreshold,
			Veto:          params.NormalVeto,
			Participation: params.NormalParticipation,
			Penalty:       params.NormalPenalty,
		}
	default:
		panic("There is no level for this proposal which type is " + p.GetProposalType().String())
	}
}

func (keeper Keeper) GetSystemHaltPeriod(ctx sdk.Context) (SystemHaltPeriod int64) {
	keeper.paramSpace.Get(ctx, KeySystemHaltPeriod, &SystemHaltPeriod)
	return
}
