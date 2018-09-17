package parameter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace                    sdk.CodespaceType = 5
	CodeInvalidMinDeposit               sdk.CodeType      = 100
	CodeInvalidMinDepositDenom          sdk.CodeType      = 101
	CodeInvalidMinDepositAmount         sdk.CodeType      = 102
	CodeInvalidDepositPeriod            sdk.CodeType      = 103
	CodeInvalidCurrentUpgradeProposalID sdk.CodeType      = 103
)
