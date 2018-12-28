package govparams

import (
	sdk "github.com/irisnet/irishub/types"
)

// Returns the current Deposit Procedure from the global param store
func GetDepositProcedure(ctx sdk.Context) DepositProcedure {
	DepositProcedureParameter.LoadValue(ctx)
	return DepositProcedureParameter.Value
}

// Returns the current Voting Procedure from the global param store
func GetVotingProcedure(ctx sdk.Context) VotingProcedure {
	VotingProcedureParameter.LoadValue(ctx)
	return VotingProcedureParameter.Value
}

// Returns the current Tallying Procedure from the global param store
func GetTallyingProcedure(ctx sdk.Context) TallyingProcedure {
	TallyingProcedureParameter.LoadValue(ctx)
	return TallyingProcedureParameter.Value
}
