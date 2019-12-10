package service

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	if err := ValidateGenesis(data); err != nil {
		panic(err.Error())
	}

	k.SetParams(ctx, data.Params)
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return NewGenesisState(k.GetParams(ctx))
}

// PrepForZeroHeightGenesis refunds all deposits, service fees, returned fees and incoming fees
func PrepForZeroHeightGenesis(ctx sdk.Context, k Keeper) {
	// refund deposits from all binding services
	if err := k.RefundDeposits(ctx); err != nil {
		panic(fmt.Sprintf("failed to refund deposits: %s", err))
	}

	// refund service fees from all active requests
	if err := k.RefundServiceFees(ctx); err != nil {
		panic(fmt.Sprintf("failed to refund service fees: %s", err))
	}

	// refund all incoming fees
	if err := k.RefundIncomingFees(ctx); err != nil {
		panic(fmt.Sprintf("failed to refund incoming fees: %s", err))
	}

	// refund all returned fees
	if err := k.RefundReturnedFees(ctx); err != nil {
		panic(fmt.Sprintf("failed to refund rerurned fees: %s", err))
	}
}
