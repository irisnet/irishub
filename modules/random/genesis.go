package random

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"mods.irisnet.org/random/keeper"
	"mods.irisnet.org/random/types"
)

// InitGenesis stores the genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	if err := types.ValidateGenesis(data); err != nil {
		panic(fmt.Errorf("failed to initialize random genesis state: %s", err.Error()))
	}
	for height, requests := range data.PendingRandomRequests {
		for _, request := range requests.Requests {
			h, _ := strconv.ParseInt(height, 10, 64)
			k.EnqueueRandomRequest(ctx, h, types.GenerateRequestID(request), request)
		}
	}
}

// ExportGenesis outputs the genesis state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	pendingRequests := make(map[string]types.Requests)

	k.IterateRandomRequestQueue(ctx, func(height int64, reqID []byte, request types.Request) bool {
		leftHeight := fmt.Sprintf("%d", height)
		heightRequests, ok := pendingRequests[leftHeight]
		if ok {
			heightRequests.Requests = append(heightRequests.Requests, request)
		} else {
			heightRequests = types.Requests{
				Requests: []types.Request{request},
			}
		}
		pendingRequests[leftHeight] = heightRequests
		return false
	})

	return &types.GenesisState{PendingRandomRequests: pendingRequests}
}

func PrepForZeroHeightGenesis(ctx sdk.Context, k keeper.Keeper) {
	k.IterateRandomRequestQueue(ctx, func(height int64, reqID []byte, request types.Request) bool {
		leftHeight := height - ctx.BlockHeight() + 1
		k.DequeueRandomRequest(ctx, height, reqID)
		k.EnqueueRandomRequest(ctx, leftHeight, reqID, request)
		return false
	})
}
