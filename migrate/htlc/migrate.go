package htlc

import (
	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	htlckeeper "github.com/irisnet/irismod/modules/htlc/keeper"
	"github.com/irisnet/irismod/modules/htlc/types"
)

func Migrate(ctx sdk.Context, cdc codec.Marshaler, k htlckeeper.Keeper, bk bankkeeper.Keeper) error {
	store := ctx.KVStore(sdk.NewKVStoreKey(types.StoreKey))

	// Delete expired queue
	store.Delete(HTLCExpiredQueueKey)

	iterator := sdk.KVStorePrefixIterator(store, HTLCKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		hashLock := tmbytes.HexBytes(iterator.Key()[1:])

		var htlc HTLC
		cdc.MustUnmarshalBinaryBare(iterator.Value(), &htlc)

		sender, err := sdk.AccAddressFromBech32(htlc.Sender)
		if err != nil {
			return err
		}
		receiver, err := sdk.AccAddressFromBech32(htlc.To)
		if err != nil {
			return err
		}
		id := types.GetID(sender, receiver, htlc.Amount, hashLock)
		expirationHeight := htlc.ExpirationHeight
		closedBlock := uint64(0)

		var state types.HTLCState
		switch htlc.State {
		case Open:
			state = types.Open
			// Add to expired queue
			k.AddHTLCToExpiredQueue(ctx, expirationHeight, id)
		case Completed:
			state = types.Completed
		case Expired:
			// Refund expired htlc
			state = types.Refunded
			if err = bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, htlc.Amount); err != nil {
				return err
			}
			closedBlock = uint64(ctx.BlockHeight())
		case Refunded:
			state = types.Refunded
		}

		// Delete origin htlc
		store.Delete(GetHTLCKey(hashLock))

		newHTLC := types.HTLC{
			Id:                   id.String(),
			Sender:               htlc.Sender,
			To:                   htlc.To,
			ReceiverOnOtherChain: htlc.ReceiverOnOtherChain,
			SenderOnOtherChain:   "",
			Amount:               htlc.Amount,
			HashLock:             hashLock.String(),
			Secret:               htlc.Secret,
			Timestamp:            htlc.Timestamp,
			ExpirationHeight:     expirationHeight,
			State:                state,
			ClosedBlock:          closedBlock,
			Transfer:             false,
			Direction:            types.Invalid,
		}
		// Set new htlc
		k.SetHTLC(ctx, newHTLC, id)
	}

	return nil
}
