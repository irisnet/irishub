package keeper

import (
	"github.com/irisnet/irishub/app/v3/service/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

func (k Keeper) Init(ctx sdk.Context, svcDefinitions []types.ServiceDefinition) {
	for _, definition := range svcDefinitions {
		k.SetServiceDefinition(ctx, definition)
	}
}
