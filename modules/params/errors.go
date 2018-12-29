package params

import (
	sdk "github.com/irisnet/irishub/types"
	"fmt"
)

const (
	DefaultCodespace sdk.CodespaceType = "params"
	//
	CodeInvalidString sdk.CodeType = 0
	//gov
	CodeInvalidMinDeposit        sdk.CodeType = 100
	CodeInvalidMinDepositDenom   sdk.CodeType = 101
	CodeInvalidMinDepositAmount  sdk.CodeType = 102
	CodeInvalidDepositPeriod     sdk.CodeType = 103
	CodeInvalidVotingPeriod      sdk.CodeType = 104
	CodeInvalidVotingProcedure   sdk.CodeType = 105
	CodeInvalidThreshold         sdk.CodeType = 106
	CodeInvalidParticipation     sdk.CodeType = 107
	CodeInvalidVeto              sdk.CodeType = 108
	CodeInvalidGovernancePenalty sdk.CodeType = 109
	CodeInvalidTallyingProcedure sdk.CodeType = 110
	CodeInvalidKey               sdk.CodeType = 111
	CodeInvalidModule            sdk.CodeType = 112
	CodeInvalidQueryParams       sdk.CodeType = 113
	CodeInvalidMaxProposalNum    sdk.CodeType = 114
	CodeInvalidSystemHaltPeriod  sdk.CodeType = 115
	//service
	CodeInvalidMaxRequestTimeout    sdk.CodeType = 200
	CodeInvalidMinDepositMultiple   sdk.CodeType = 201
	CodeInvalidServiceFeeTax        sdk.CodeType = 202
	CodeInvalidSlashFraction        sdk.CodeType = 203
	CodeInvalidServiceParams        sdk.CodeType = 204
	CodeInvalidArbitrationTimeLimit sdk.CodeType = 205
	CodeComplaintRetrospect         sdk.CodeType = 206
	//upgrade
	CodeInvalidUpgradeParams sdk.CodeType = 300
	//mint
	CodeInvalidMintInflation sdk.CodeType = 400
	//slash
	CodeInvalidSlashParams   sdk.CodeType = 500
)

func ErrInvalidString(valuestr string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidString, fmt.Sprintf("%s can't convert to a specific type", valuestr))
}
