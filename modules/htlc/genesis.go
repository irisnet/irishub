package htlc

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	"github.com/irisnet/irismod/modules/htlc/keeper"
	"github.com/irisnet/irismod/modules/htlc/types"
)

// InitGenesis stores the genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	if err := types.ValidateGenesis(data); err != nil {
		panic(err.Error())
	}

	k.SetPreviousBlockTime(ctx, data.PreviousBlockTime)
	k.SetParams(ctx, data.Params)
	for _, supply := range data.Supplies {
		k.SetAssetSupply(ctx, supply, supply.CurrentSupply.Denom)
	}

	var incomingSupplies sdk.Coins
	var outgoingSupplies sdk.Coins
	for _, htlc := range data.PendingHtlcs {
		id, err := hex.DecodeString(htlc.Id)
		if err != nil {
			panic(err.Error())
		}

		if htlc.State != types.Open {
			panic(fmt.Sprintf("htlc %s has invalid status %s", htlc.Id, htlc.State.String()))
		}

		if !htlc.Transfer {
			k.SetHTLC(ctx, htlc, id)
			k.AddHTLCToExpiredQueue(ctx, htlc.ExpirationHeight, id)
			continue
		}

		// htlt assets must be both supported and active
		if err := k.ValidateLiveAsset(ctx, htlc.Amount[0]); err != nil {
			panic(err.Error())
		}
		k.SetHTLC(ctx, htlc, id)
		k.AddHTLCToExpiredQueue(ctx, htlc.ExpirationHeight, id)

		switch htlc.Direction {
		case types.Incoming:
			incomingSupplies = incomingSupplies.Add(htlc.Amount...)
		case types.Outgoing:
			outgoingSupplies = outgoingSupplies.Add(htlc.Amount...)
		default:
			panic(fmt.Sprintf("htlt %s has invalid direction %s", htlc.Id, htlc.Direction.String()))
		}
	}

	// Asset's given incoming/outgoing supply much match the amount of coins in incoming/outgoing HTLTs
	supplies := k.GetAllAssetSupplies(ctx)
	for _, supply := range supplies {
		incomingSupply := incomingSupplies.AmountOf(supply.CurrentSupply.Denom)
		if !supply.IncomingSupply.Amount.Equal(incomingSupply) {
			panic(fmt.Sprintf(
				"asset's incoming supply %s does not match amount %s in incoming atomic swaps",
				supply.IncomingSupply, incomingSupply,
			))
		}
		outgoingSupply := outgoingSupplies.AmountOf(supply.CurrentSupply.Denom)
		if !supply.OutgoingSupply.Amount.Equal(outgoingSupply) {
			panic(fmt.Sprintf(
				"asset's outgoing supply %s does not match amount %s in outgoing atomic swaps",
				supply.OutgoingSupply, outgoingSupply,
			))
		}
		limit, err := k.GetSupplyLimit(ctx, supply.CurrentSupply.Denom)
		if err != nil {
			panic(err)
		}
		if supply.CurrentSupply.Amount.GT(limit.Limit) {
			panic(fmt.Sprintf("asset's current supply %s is over the supply limit %s", supply.CurrentSupply, limit.Limit))
		}
		if supply.IncomingSupply.Amount.GT(limit.Limit) {
			panic(fmt.Sprintf("asset's incoming supply %s is over the supply limit %s", supply.IncomingSupply, limit.Limit))
		}
		if supply.IncomingSupply.Amount.Add(supply.CurrentSupply.Amount).GT(limit.Limit) {
			panic(fmt.Sprintf("asset's incoming supply + current supply %s is over the supply limit %s", supply.IncomingSupply.Add(supply.CurrentSupply), limit.Limit))
		}
		if supply.OutgoingSupply.Amount.GT(limit.Limit) {
			panic(fmt.Sprintf("asset's outgoing supply %s is over the supply limit %s", supply.OutgoingSupply, limit.Limit))
		}
	}
}

// ExportGenesis outputs the genesis state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	pendingHTLCs := []types.HTLC{}
	k.IterateHTLCs(
		ctx,
		func(_ tmbytes.HexBytes, h types.HTLC) (stop bool) {
			pendingHTLCs = append(pendingHTLCs, h)
			return false
		},
	)

	supplies := k.GetAllAssetSupplies(ctx)
	previousBlockTime, found := k.GetPreviousBlockTime(ctx)
	if !found {
		previousBlockTime = types.DefaultPreviousBlockTime
	}

	return types.NewGenesisState(
		k.GetParams(ctx),
		pendingHTLCs,
		supplies,
		previousBlockTime,
	)
}

func PrepForZeroHeightGenesis(ctx sdk.Context, k keeper.Keeper) {
	k.IterateHTLCs(
		ctx,
		func(id tmbytes.HexBytes, h types.HTLC) (stop bool) {
			if h.State == types.Open {
				h.ExpirationHeight = h.ExpirationHeight - uint64(ctx.BlockHeight()) + 1
				k.SetHTLC(ctx, h, id)
			}
			return false
		},
	)
	// TODO: update asset supplies and previous block time
}
