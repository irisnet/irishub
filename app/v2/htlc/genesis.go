package htlc

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

// InitGenesis stores genesis data
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	for hashLockHex, htlc := range data.PendingHTLCs {
		hashLock, err := hex.DecodeString(hashLockHex)
		if err != nil {
			panic(fmt.Errorf("failed to initialize HTLC genesis state: %s", err.Error()))
		}

		k.SetHTLC(ctx, htlc, hashLock)
		k.AddHTLCToExpireQueue(ctx, htlc.ExpireHeight, hashLock)
	}
}

// ExportGenesis outputs genesis data
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	pendingHTLCs := make(map[string]HTLC)

	k.IterateHTLCs(ctx, func(hlock []byte, h HTLC) (stop bool) {
		if h.State == OPEN || h.State == EXPIRED {
			if h.State == OPEN {
				h.ExpireHeight = h.ExpireHeight - uint64(ctx.BlockHeight()) + 1
				pendingHTLCs[hex.EncodeToString(hlock)] = h
			} else {
				_, err := k.RefundHTLC(ctx, hlock)
				if err != nil {
					panic(fmt.Errorf("failed to export HTLC genesis state: %s", hex.EncodeToString(hlock)))
				}
			}
		}

		return false
	})

	return GenesisState{
		PendingHTLCs: pendingHTLCs,
	}
}

// DefaultGenesisState gets the default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		PendingHTLCs: map[string]HTLC{},
	}
}

// DefaultGenesisStateForTest gets the default genesis state for test
func DefaultGenesisStateForTest() GenesisState {
	return GenesisState{
		PendingHTLCs: map[string]HTLC{},
	}
}
