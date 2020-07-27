package random

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/random/keeper"
	"github.com/irisnet/irishub/modules/random/types"
)

// InitGenesis stores genesis data
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	if err := types.ValidateGenesis(data); err != nil {
		panic(fmt.Errorf("failed to initialize rand genesis state: %s", err.Error()))
	}
	for height, requests := range data.PendingRandomRequests {
		for _, request := range requests {
			h, _ := strconv.ParseInt(height, 10, 64)
			k.EnqueueRandomRequest(ctx, h, types.GenerateRequestID(request), request)
		}
	}
}

// ExportGenesis outputs genesis data
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {
	pendingRequests := make(map[string][]types.Request)

	k.IterateRandomRequestQueue(ctx, func(height int64, request types.Request) bool {
		leftHeight := fmt.Sprintf("%d", height-ctx.BlockHeight()+1)
		pendingRequests[leftHeight] = append(pendingRequests[leftHeight], request)
		return false
	})

	return types.GenesisState{PendingRandomRequests: pendingRequests}
}
