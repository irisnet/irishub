package service

import (
	"encoding/hex"
	"fmt"

	cmn "github.com/tendermint/tendermint/libs/common"

	sdk "github.com/irisnet/irishub/types"
)

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	if err := ValidateGenesis(data); err != nil {
		panic(err.Error())
	}

	k.SetParamSet(ctx, data.Params)

	// TODO: save service definitions

	// TODO: save service definition bindings

	// TODO: save withdraw addresses

	for reqContextIDStr, requestContext := range data.RequestContexts {
		requestContextID, _ := hex.DecodeString(reqContextIDStr)
		k.SetRequestContext(ctx, requestContextID, requestContext)
	}
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	definitions := []ServiceDefinition{}
	bindings := make(map[string][]ServiceBinding)
	withdrawAddress := []sdk.AccAddress{}
	requestContexts := make(map[string]RequestContext)

	k.IterateServiceDefinitions(
		ctx,
		func(name string, definition ServiceDefinition) bool {
			// TODO
			return false
		},
	)

	k.IterateServiceBindings(
		ctx,
		func(name string, binding ServiceBinding) bool {
			// TODO
			return false
		},
	)

	k.IterateWithdrawAddresses(
		ctx,
		func(providerAddress sdk.AccAddress, withdrawAddress sdk.AccAddress) bool {
			// TODO
			return false
		},
	)

	k.IterateRequestContexts(
		ctx,
		func(requestContextID cmn.HexBytes, requestContext RequestContext) bool {
			requestContexts[requestContextID.String()] = requestContext
			return false
		},
	)

	k.IterateNewRequestBatch(
		ctx,
		func(height cmn.HexBytes, requestContextID cmn.HexBytes) bool {
			// TODO
			return false
		},
	)

	k.IterateExpiredRequestBatch(
		ctx,
		func(height cmn.HexBytes, requestContextID cmn.HexBytes) bool {
			// TODO
			return false
		},
	)

	return NewGenesisState(
		k.GetParamSet(ctx),
		definitions,
		bindings,
		withdrawAddress,
		requestContexts,
	)
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
