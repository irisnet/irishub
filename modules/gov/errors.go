//nolint
package gov

import (
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

const (
	DefaultCodespace sdk.CodespaceType = "gov"

	CodeUnknownProposal         sdk.CodeType = 1
	CodeInactiveProposal        sdk.CodeType = 2
	CodeAlreadyActiveProposal   sdk.CodeType = 3
	CodeAlreadyFinishedProposal sdk.CodeType = 4
	CodeAddressNotStaked        sdk.CodeType = 5
	CodeInvalidTitle            sdk.CodeType = 6
	CodeInvalidDescription      sdk.CodeType = 7
	CodeInvalidProposalType     sdk.CodeType = 8
	CodeInvalidVote             sdk.CodeType = 9
	CodeInvalidGenesis          sdk.CodeType = 10
	CodeInvalidProposalStatus   sdk.CodeType = 11
	////////////////////  iris begin  ///////////////////////////
	CodeInvalidParam          sdk.CodeType = 12
	CodeInvalidParamOp        sdk.CodeType = 13
	CodeSwitchPeriodInProcess sdk.CodeType = 14
	CodeInvalidPercent        sdk.CodeType = 15
	CodeInvalidUsageType      sdk.CodeType = 16
	CodeInvalidInput          sdk.CodeType = 16
	////////////////////  iris end  /////////////////////////////
)

//----------------------------------------
// Error constructors

func ErrUnknownProposal(codespace sdk.CodespaceType, proposalID uint64) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownProposal, fmt.Sprintf("Unknown proposal with id %d", proposalID))
}

func ErrInactiveProposal(codespace sdk.CodespaceType, proposalID uint64) sdk.Error {
	return sdk.NewError(codespace, CodeInactiveProposal, fmt.Sprintf("Inactive proposal with id %d", proposalID))
}

func ErrAlreadyActiveProposal(codespace sdk.CodespaceType, proposalID uint64) sdk.Error {
	return sdk.NewError(codespace, CodeAlreadyActiveProposal, fmt.Sprintf("Proposal %d has been already active", proposalID))
}

func ErrAlreadyFinishedProposal(codespace sdk.CodespaceType, proposalID uint64) sdk.Error {
	return sdk.NewError(codespace, CodeAlreadyFinishedProposal, fmt.Sprintf("Proposal %d has already passed its voting period", proposalID))
}

func ErrAddressNotStaked(codespace sdk.CodespaceType, address sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeAddressNotStaked, fmt.Sprintf("Address %s is not staked and is thus ineligible to vote", address))
}

func ErrInvalidTitle(codespace sdk.CodespaceType, title string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTitle, fmt.Sprintf("Proposal Title '%s' is not valid", title))
}

func ErrInvalidDescription(codespace sdk.CodespaceType, description string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDescription, fmt.Sprintf("Proposal Desciption '%s' is not valid", description))
}

func ErrInvalidProposalType(codespace sdk.CodespaceType, proposalType ProposalKind) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidProposalType, fmt.Sprintf("Proposal Type '%s' is not valid", proposalType))
}

func ErrInvalidVote(codespace sdk.CodespaceType, voteOption VoteOption) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidVote, fmt.Sprintf("'%v' is not a valid voting option", voteOption))
}

func ErrInvalidGenesis(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidVote, msg)
}

////////////////////  iris begin  ///////////////////////////
func ErrInvalidParam(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidParam, fmt.Sprintf("Param is not valid"))
}

func ErrInvalidParamOp(codespace sdk.CodespaceType, opStr string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidParamOp, fmt.Sprintf("Op '%s' is not valid", opStr))
}
func ErrSwitchPeriodInProcess(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSwitchPeriodInProcess, fmt.Sprintf("Software Upgrade Switch Period is in process."))
}

func ErrInvalidPercent(codespace sdk.CodespaceType, percent sdk.Dec) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidPercent, fmt.Sprintf("invalid percent [%s], must >0 & <=1", percent))
}

func ErrInvalidUsageType(codespace sdk.CodespaceType, usageType UsageType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidUsageType, fmt.Sprintf("Usage Type '%s' is not valid", usageType))
}

func ErrNotTrustee(codespace sdk.CodespaceType, trustee sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, fmt.Sprintf("[%s] is not a trustee address", trustee))
}

func ErrNotProfiler(codespace sdk.CodespaceType, profiler sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, fmt.Sprintf("[%s] is not a profiler address", profiler))
}

////////////////////  iris end  /////////////////////////////
