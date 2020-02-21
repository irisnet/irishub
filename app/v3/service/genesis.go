package service

import (
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	if err := ValidateGenesis(data); err != nil {
		panic(err.Error())
	}

	k.SetParamSet(ctx, data.Params)
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return NewGenesisState(k.GetParamSet(ctx))
}

// PrepForZeroHeightGenesis refunds the deposits, service fees and earned fees
func PrepForZeroHeightGenesis(ctx sdk.Context, k Keeper) {
	// refund deposits from all binding services
	if err := k.RefundDeposits(ctx); err != nil {
		panic(fmt.Sprintf("failed to refund the deposits: %s", err))
	}

	// refund service fees from all active requests
	if err := k.RefundServiceFees(ctx); err != nil {
		panic(fmt.Sprintf("failed to refund the service fees: %s", err))
	}

	// refund all the earned fees
	if err := k.RefundEarnedFees(ctx); err != nil {
		panic(fmt.Sprintf("failed to refund the earned fees: %s", err))
	}
}
