package rand

import (
	"fmt"
	"strconv"

	sdk "github.com/irisnet/irishub/types"
)

// InitGenesis stores genesis data
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	if err := ValidateGenesis(data); err != nil {
		panic(fmt.Errorf("failed to initialize rand genesis state: %s", err.Error()))
	}

	for height, requests := range data.PendingRandRequests {
		for _, request := range requests {

			// check request context exists
			if request.Oracle {
				if _, success := k.GetRequestContext(ctx, request.ServiceContextID); !success {
					panic(fmt.Errorf("no service request context"))
				}
			}

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

	return GenesisState{
		PendingRandRequests: pendingRequests,
	}
}

// DefaultGenesisState gets the default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		PendingRandRequests: map[string][]Request{},
	}
}

// DefaultGenesisStateForTest gets the default genesis state for test
func DefaultGenesisStateForTest() GenesisState {
	return GenesisState{
		PendingRandRequests: map[string][]Request{},
	}
}
