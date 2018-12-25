package gov

import (
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/types/gov/params"
	govtypes "github.com/irisnet/irishub/types/gov"
	"time"
)

//-----------------------------------------------------------
// ProposalLevel

// Type that represents Proposal Level as a byte
type ProposalLevel byte

//nolint
const (
	ProposalLevelNil        ProposalLevel = 0x00
	ProposalLevelCritical   ProposalLevel = 0x01
	ProposalLevelImportant  ProposalLevel = 0x02
	////////////////////  iris begin  /////////////////////////////
	ProposalLevelNormal     ProposalLevel = 0x03
	////////////////////  iris end  /////////////////////////////
)

func (p ProposalLevel) string() string {
	switch p {
	case ProposalLevelCritical:
		return "ciritical"
	case ProposalLevelImportant:
		return  "important"
	case ProposalLevelNormal:
		return "normal"
	default:
		return  " "
	}
}

func GetProposalLevel(p govtypes.Proposal) ProposalLevel {
	return GetProposalLevelByProposalKind(p.GetProposalType())
}

func GetProposalLevelByProposalKind(p govtypes.ProposalKind) ProposalLevel {
	switch p {
	case govtypes.ProposalTypeTxTaxUsage:
		return ProposalLevelNormal
	case govtypes.ProposalTypeParameterChange:
		return ProposalLevelImportant
	case govtypes.ProposalTypeSoftwareHalt:
		return ProposalLevelCritical
	case govtypes.ProposalTypeSoftwareUpgrade:
		return ProposalLevelCritical
	default:
		return  ProposalLevelNil
	}
}

// Returns the current Deposit Procedure from the global param store
func GetDepositProcedure(ctx sdk.Context) govparams.DepositProcedure {
	govparams.DepositProcedureParameter.LoadValue(ctx)
	return govparams.DepositProcedureParameter.Value
}

func GetMinDeposit(ctx sdk.Context, p govtypes.Proposal) sdk.Coins {
	govparams.DepositProcedureParameter.LoadValue(ctx)
	switch GetProposalLevel(p) {
	case ProposalLevelCritical:
		return govparams.DepositProcedureParameter.Value.CriticalMinDeposit
	case ProposalLevelImportant:
		return govparams.DepositProcedureParameter.Value.ImportantMinDeposit
	case ProposalLevelNormal:
		return govparams.DepositProcedureParameter.Value.NormalMinDeposit
	default:
		panic("There is no level for this proposal which type is "+ p.GetProposalType().String())
	}
}

func GetDepositPeriod(ctx sdk.Context) time.Duration {
	govparams.DepositProcedureParameter.LoadValue(ctx)
	return govparams.DepositProcedureParameter.Value.MaxDepositPeriod
}


// Returns the current Voting Procedure from the global param store
func GetVotingProcedure(ctx sdk.Context) govparams.VotingProcedure {
	govparams.VotingProcedureParameter.LoadValue(ctx)
	return govparams.VotingProcedureParameter.Value
}

func GetVotingPeriod(ctx sdk.Context, p govtypes.Proposal) time.Duration {
	govparams.VotingProcedureParameter.LoadValue(ctx)
	switch GetProposalLevel(p) {
	case ProposalLevelCritical:
		return govparams.VotingProcedureParameter.Value.CriticalVotingPeriod
	case ProposalLevelImportant:
		return govparams.VotingProcedureParameter.Value.ImportantVotingPeriod
	case ProposalLevelNormal:
		return govparams.VotingProcedureParameter.Value.NormalVotingPeriod
	default:
		panic("There is no level for this proposal which type is "+ p.GetProposalType().String())
	}
}

// Returns the current Tallying Procedure from the global param store
func GetTallyingProcedure(ctx sdk.Context) govparams.TallyingProcedure {
	govparams.TallyingProcedureParameter.LoadValue(ctx)
	return govparams.TallyingProcedureParameter.Value
}
