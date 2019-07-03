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

// Returns the current Deposit Procedure from the global param store
func (p ProposalLevel) GetDepositProcedure(ctx sdk.Context, k Keeper) DepositProcedure {
	params := k.GetParamSet(ctx)
	switch p {
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
		panic("There is no level for this proposal which type is " + p.string())
	}
}

// Returns the current Voting Procedure from the global param store
func (p ProposalLevel) GetVotingProcedure(ctx sdk.Context, k Keeper) VotingProcedure {
	params := k.GetParamSet(ctx)
	switch p {
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
		panic("There is no level for this proposal which type is " + p.string())
	}
}

func (p ProposalLevel) GetMaxNumByProposalLevel(ctx sdk.Context, k Keeper) uint64 {
	params := k.GetParamSet(ctx)
	switch p {
	case ProposalLevelCritical:
		return params.CriticalMaxNum

	case ProposalLevelImportant:
		return params.ImportantMaxNum

	case ProposalLevelNormal:
		return params.NormalMaxNum
	default:
		panic("There is no level for this proposal which type is " + p.string())
	}
}

// Returns the current Tallying Procedure from the global param store
func (p ProposalLevel) GetTallyingProcedure(ctx sdk.Context, k Keeper) TallyingProcedure {
	params := k.GetParamSet(ctx)
	switch p {
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
		panic("There is no level for this proposal which type is " + p.string())
	}
}

func (p ProposalLevel) AddProposalNum(ctx sdk.Context, k Keeper, args ...interface{}) {
	switch p {
	case ProposalLevelCritical:
		proposalID := args[0].(uint64)
		k.AddCriticalProposalNum(ctx, proposalID)
	case ProposalLevelImportant:
		k.AddImportantProposalNum(ctx)
	case ProposalLevelNormal:
		k.AddNormalProposalNum(ctx)
	default:
		panic("There is no level for this proposal which type is " + p.string())
	}
}

func (p ProposalLevel) SubProposalNum(ctx sdk.Context, k Keeper) {
	switch p {
	case ProposalLevelCritical:
		k.SubCriticalProposalNum(ctx)
	case ProposalLevelImportant:
		k.SubImportantProposalNum(ctx)
	case ProposalLevelNormal:
		k.SubNormalProposalNum(ctx)
	default:
		panic("There is no level for this proposal which type is " + p.string())
	}
}

func (p ProposalLevel) HasReachedTheMaxProposalNum(ctx sdk.Context, k Keeper) (uint64, bool) {
	ctx.Logger().Debug("Proposals Distribution",
		"CriticalProposalNum", k.GetCriticalProposalNum(ctx),
		"ImportantProposalNum", k.GetImportantProposalNum(ctx),
		"NormalProposalNum", k.GetNormalProposalNum(ctx))

	maxNum := p.GetMaxNumByProposalLevel(ctx, k)
	switch p {
	case ProposalLevelCritical:
		return k.GetCriticalProposalNum(ctx), k.GetCriticalProposalNum(ctx) == maxNum
	case ProposalLevelImportant:
		return k.GetImportantProposalNum(ctx), k.GetImportantProposalNum(ctx) == maxNum
	case ProposalLevelNormal:
		return k.GetNormalProposalNum(ctx), k.GetNormalProposalNum(ctx) == maxNum
	default:
		panic("There is no level for this proposal")
	}
}
