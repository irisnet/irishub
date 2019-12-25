package rand

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis stores genesis data
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	if err := ValidateGenesis(data); err != nil {
		panic(fmt.Errorf("failed to initialize rand genesis state: %s", err.Error()))
	}
	for height, requests := range data.PendingRandRequests {
		for _, request := range requests {
			h, _ := strconv.ParseInt(height, 10, 64)
			k.EnqueueRandRequest(ctx, h, GenerateRequestID(request), request)
		}
	}
}

// ExportGenesis outputs genesis data
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	pendingRequests := make(map[string][]Request)

	k.IterateRandRequestQueue(ctx, func(height int64, request Request) bool {
		leftHeight := fmt.Sprintf("%d", height-ctx.BlockHeight()+1)
		pendingRequests[leftHeight] = append(pendingRequests[leftHeight], request)
		return false
	})

	return GenesisState{PendingRandRequests: pendingRequests}
}
