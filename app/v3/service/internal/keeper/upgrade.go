package keeper

import (
	"github.com/irisnet/irishub/app/v3/service/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

// Init initializes the service module
func (k Keeper) Init(ctx sdk.Context, svcDefinitions []types.ServiceDefinition) {
	// reset params
	k.SetParamSet(ctx, types.DefaultParams())

	for _, definition := range svcDefinitions {
		k.SetServiceDefinition(ctx, definition)
	}
}
