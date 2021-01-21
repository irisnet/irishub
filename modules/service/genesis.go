package service

import (
	"encoding/hex"
	"fmt"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/service/keeper"
	"github.com/irisnet/irismod/modules/service/types"
)

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	if err := types.ValidateGenesis(data); err != nil {
		panic(err.Error())
	}

	k.SetParams(ctx, data.Params)

	for _, definition := range data.Definitions {
		k.SetServiceDefinition(ctx, definition)
	}

	for _, binding := range data.Bindings {
		if err := k.SetServiceBindingForGenesis(ctx, binding); err != nil {
			panic(err.Error())
		}
	}

	for ownerAddressStr, withdrawAddress := range data.WithdrawAddresses {
		ownerAddress, _ := sdk.AccAddressFromBech32(ownerAddressStr)
		withdrawAddress, _ := sdk.AccAddressFromBech32(withdrawAddress)
		k.SetWithdrawAddress(ctx, ownerAddress, withdrawAddress)
	}

	for reqContextIDStr, requestContext := range data.RequestContexts {
		requestContextID, _ := hex.DecodeString(reqContextIDStr)
		k.SetRequestContext(ctx, requestContextID, *requestContext)
	}
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	var definitions []types.ServiceDefinition
	var bindings []types.ServiceBinding
	withdrawAddresses := make(map[string]string)
	requestContexts := make(map[string]*types.RequestContext)

	k.IterateServiceDefinitions(
		ctx,
		func(definition types.ServiceDefinition) bool {
			definitions = append(definitions, definition)
			return false
		},
	)

	k.IterateServiceBindings(
		ctx,
		func(binding types.ServiceBinding) bool {
			bindings = append(bindings, binding)
			return false
		},
	)

	k.IterateWithdrawAddresses(
		ctx,
		func(ownerAddress sdk.AccAddress, withdrawAddress sdk.AccAddress) bool {
			withdrawAddresses[ownerAddress.String()] = withdrawAddress.String()
			return false
		},
	)

	k.IterateRequestContexts(
		ctx,
		func(requestContextID tmbytes.HexBytes, requestContext types.RequestContext) bool {
			requestContexts[requestContextID.String()] = &requestContext
			return false
		},
	)

	return types.NewGenesisState(
		k.GetParams(ctx),
		definitions,
		bindings,
		withdrawAddresses,
		requestContexts,
	)
}

// PrepForZeroHeightGenesis refunds the deposits, service fees and earned fees
func PrepForZeroHeightGenesis(ctx sdk.Context, k keeper.Keeper) {
	// refund service fees from all active requests
	if err := k.RefundServiceFees(ctx); err != nil {
		panic(fmt.Sprintf("failed to refund the service fees: %s", err))
	}

	// refund all the earned fees
	if err := k.RefundEarnedFees(ctx); err != nil {
		panic(fmt.Sprintf("failed to refund the earned fees: %s", err))
	}

	// reset request contexts state and batch
	if err := k.ResetRequestContextsStateAndBatch(ctx); err != nil {
		panic(fmt.Sprintf("failed to reset the request context state: %s", err))
	}
}
