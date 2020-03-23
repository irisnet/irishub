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

	for _, definition := range data.Definitions {
		k.SetServiceDefinition(ctx, definition)
	}

	for _, binding := range data.Bindings {
		k.SetServiceBinding(ctx, binding)
	}

	for providerAddressStr, withdrawAddress := range data.WithdrawAddresses {
		providerAddress, _ := sdk.AccAddressFromBech32(providerAddressStr)
		k.SetWithdrawAddress(ctx, providerAddress, withdrawAddress)
	}

	for reqContextIDStr, requestContext := range data.RequestContexts {
		requestContextID, _ := hex.DecodeString(reqContextIDStr)
		k.SetRequestContext(ctx, requestContextID, requestContext)
	}
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	definitions := []ServiceDefinition{}
	bindings := []ServiceBinding{}
	withdrawAddresses := make(map[string]sdk.AccAddress)
	requestContexts := make(map[string]RequestContext)

	k.IterateServiceDefinitions(
		ctx,
		func(definition ServiceDefinition) bool {
			definitions = append(definitions, definition)
			return false
		},
	)

	k.IterateServiceBindings(
		ctx,
		func(binding ServiceBinding) bool {
			bindings = append(bindings, binding)
			return false
		},
	)

	k.IterateWithdrawAddresses(
		ctx,
		func(providerAddress sdk.AccAddress, withdrawAddress sdk.AccAddress) bool {
			withdrawAddresses[providerAddress.String()] = withdrawAddress
			return false
		},
	)

	k.IterateRequestContexts(
		ctx,
		func(requestContextID cmn.HexBytes, requestContext RequestContext) bool {
			if requestContext.State != COMPLETED {
				requestContext.State = PAUSED
				requestContexts[requestContextID.String()] = requestContext
			}
			return false
		},
	)

	return NewGenesisState(
		k.GetParamSet(ctx),
		definitions,
		bindings,
		withdrawAddresses,
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
