package params

import (
	sdk "github.com/irisnet/irishub/types"
)

const (
	DefaultCodespace                    sdk.CodespaceType = "params"
	//gov
	CodeInvalidMinDeposit               sdk.CodeType      = 100
	CodeInvalidMinDepositDenom          sdk.CodeType      = 101
	CodeInvalidMinDepositAmount         sdk.CodeType      = 102
	CodeInvalidDepositPeriod            sdk.CodeType      = 103
	CodeInvalidVotingPeriod             sdk.CodeType      = 104
	CodeInvalidVotingProcedure          sdk.CodeType      = 105
	CodeInvalidThreshold                sdk.CodeType      = 106
	CodeInvalidParticipation            sdk.CodeType      = 107
	CodeInvalidVeto                     sdk.CodeType      = 108
	CodeInvalidGovernancePenalty        sdk.CodeType      = 109
	CodeInvalidTallyingProcedure        sdk.CodeType      = 110
	CodeInvalidKey                      sdk.CodeType      = 111
	CodeInvalidModule                   sdk.CodeType      = 112
	CodeInvalidQueryParams              sdk.CodeType      = 113
	CodeInvalidMaxProposalNum           sdk.CodeType      = 114
    //service
	CodeInvalidMaxRequestTimeout        sdk.CodeType      = 200
	CodeInvalidMinDepositMultiple       sdk.CodeType      = 201
	CodeInvalidServiceFeeTax            sdk.CodeType      = 202
	CodeInvalidSlashFraction            sdk.CodeType      = 203
	CodeInvalidServiceParams            sdk.CodeType      = 204
	//upgrade
	CodeInvalidUpgradeParams        sdk.CodeType      = 300
)
