package rand

import (
	sdk "github.com/irisnet/irishub/types"
)

// InitGenesis stores genesis data
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	for height, requests := range data.PendingRandRequests {
		for _, request := range requests {
			reqID := GenerateRequestID(request)
			k.EnqueueRandRequest(ctx, height, reqID, request)
		}
	}
}

// ExportGenesis outputs genesis data
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	pendingRequests := make(map[int64][]Request)

	k.IterateRandRequestQueue(ctx, func(height int64, request Request) bool {
		leftHeight := height - ctx.BlockHeight() + 1
		pendingRequests[leftHeight] = append(pendingRequests[leftHeight], request)

		return false
	})

	return GenesisState{
		PendingRandRequests: pendingRequests,
	}
}

// DefaultGenesisState gets the default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		PendingRandRequests: map[int64][]Request{},
	}
}

// DefaultGenesisStateForTest gets the default genesis state for test
func DefaultGenesisStateForTest() GenesisState {
	return GenesisState{
		PendingRandRequests: map[int64][]Request{},
	}
}
