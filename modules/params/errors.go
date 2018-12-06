package params

import (
	sdk "github.com/irisnet/irishub/types"
)

const (
	DefaultCodespace                    sdk.CodespaceType = "params"
	CodeInvalidMinDeposit               sdk.CodeType      = 100
	CodeInvalidMinDepositDenom          sdk.CodeType      = 101
	CodeInvalidMinDepositAmount         sdk.CodeType      = 102
	CodeInvalidDepositPeriod            sdk.CodeType      = 103
	CodeInvalidCurrentUpgradeProposalID sdk.CodeType      = 104
	CodeInvalidVotingPeriod             sdk.CodeType      = 105
	CodeInvalidVotingProcedure          sdk.CodeType      = 106
	CodeInvalidThreshold                sdk.CodeType      = 107
	CodeInvalidParticipation            sdk.CodeType      = 108
	CodeInvalidVeto                     sdk.CodeType      = 109
	CodeInvalidTallyingProcedure        sdk.CodeType      = 110
	CodeInvalidKey                      sdk.CodeType      = 111
	CodeInvalidParamString              sdk.CodeType      = 112
	CodeInvalidModule                   sdk.CodeType      = 113
	CodeInvalidQueryParams              sdk.CodeType      = 114
	CodeInvalidMaxRequestTimeout        sdk.CodeType      = 115
	CodeInvalidMinDepositMultiple       sdk.CodeType      = 116
)
