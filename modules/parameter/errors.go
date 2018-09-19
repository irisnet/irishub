package parameter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace                    sdk.CodespaceType = 6
	CodeInvalidMinDeposit               sdk.CodeType      = 100
	CodeInvalidMinDepositDenom          sdk.CodeType      = 101
	CodeInvalidMinDepositAmount         sdk.CodeType      = 102
	CodeInvalidDepositPeriod            sdk.CodeType      = 103
	CodeInvalidCurrentUpgradeProposalID sdk.CodeType      = 104
	CodeInvalidVotingPeriod             sdk.CodeType      = 105
	CodeInvalidVotingProcedure          sdk.CodeType      = 106
)
