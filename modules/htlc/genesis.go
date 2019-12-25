package htlc

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis stores genesis data
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	if err := ValidateGenesis(data); err != nil {
		panic(fmt.Errorf("failed to initialize HTLC genesis state: %s", err.Error()))
	}
	for hashLockHex, htlc := range data.PendingHTLCs {
		hashLock, _ := hex.DecodeString(hashLockHex)
		k.SetHTLC(ctx, htlc, hashLock)
		k.AddHTLCToExpireQueue(ctx, htlc.ExpireHeight, hashLock)
	}
}

// ExportGenesis outputs genesis data
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	pendingHTLCs := make(map[string]HTLC)

	k.IterateHTLCs(ctx, func(hlock HTLCHashLock, h HTLC) (stop bool) {
		if h.State == OPEN {
			h.ExpireHeight = h.ExpireHeight - uint64(ctx.BlockHeight()) + 1
			pendingHTLCs[hex.EncodeToString(hlock)] = h
		}
		if h.State == EXPIRED {
			if _, err := k.RefundHTLC(ctx, hlock); err != nil {
				panic(fmt.Errorf("failed to export HTLC genesis state: %s", hex.EncodeToString(hlock)))
			}
		}
		return false
	})

	return GenesisState{PendingHTLCs: pendingHTLCs}
}
