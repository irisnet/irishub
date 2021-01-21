package htlc

import (
	"encoding/hex"
	"fmt"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/htlc/keeper"
	"github.com/irisnet/irismod/modules/htlc/types"
)

// InitGenesis stores the genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	if err := types.ValidateGenesis(data); err != nil {
		panic(err.Error())
	}

	for hashLockStr, htlc := range data.PendingHtlcs {
		hashLock, _ := hex.DecodeString(hashLockStr)

		k.SetHTLC(ctx, htlc, hashLock)
		k.AddHTLCToExpiredQueue(ctx, htlc.ExpirationHeight, hashLock)
	}
}

// ExportGenesis outputs the genesis state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	pendingHtlcs := make(map[string]types.HTLC)

	k.IterateHTLCs(ctx, func(hlock tmbytes.HexBytes, h types.HTLC) (stop bool) {
		if h.State == types.Open {
			pendingHtlcs[hlock.String()] = h
		}
		return false
	})

	return &types.GenesisState{
		PendingHtlcs: pendingHtlcs,
	}
}

func PrepForZeroHeightGenesis(ctx sdk.Context, k keeper.Keeper) {
	k.IterateHTLCs(
		ctx,
		func(hlock tmbytes.HexBytes, h types.HTLC) (stop bool) {
			if h.State == types.Open {
				h.ExpirationHeight = h.ExpirationHeight - uint64(ctx.BlockHeight()) + 1
				k.SetHTLC(ctx, h, hlock)
			} else if h.State == types.Expired {
				if err := k.RefundHTLC(ctx, hlock); err != nil {
					panic(fmt.Errorf("failed to export the HTLC genesis state: %s", hlock.String()))
				}
			}
			return false
		})
}
