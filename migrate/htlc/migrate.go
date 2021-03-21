package htlc

import (
	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	htlckeeper "github.com/irisnet/irismod/modules/htlc/keeper"
	htlctypes "github.com/irisnet/irismod/modules/htlc/types"
)

func Migrate(ctx sdk.Context, cdc codec.Marshaler, k htlckeeper.Keeper, bk bankkeeper.Keeper) error {
	store := ctx.KVStore(sdk.NewKVStoreKey(htlctypes.StoreKey))

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
		id := htlctypes.GetID(sender, receiver, htlc.Amount, hashLock)
		expirationHeight := htlc.ExpirationHeight
		closedBlock := uint64(0)

		var state htlctypes.HTLCState
		switch htlc.State {
		case Open:
			state = htlctypes.Open
			// Add to expired queue
			k.AddHTLCToExpiredQueue(ctx, expirationHeight, id)
		case Completed:
			state = htlctypes.Completed
		case Expired:
			// Refund expired htlc
			state = htlctypes.Refunded
			if err := bk.SendCoinsFromModuleToAccount(ctx, htlctypes.ModuleName, sender, htlc.Amount); err != nil {
				return err
			}
			closedBlock = uint64(ctx.BlockHeight())
		case Refunded:
			state = htlctypes.Refunded
		}

		// Delete origin htlc
		store.Delete(GetHTLCKey(hashLock))

		newHTLC := htlctypes.HTLC{
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
			Direction:            htlctypes.Invalid,
		}
		// Set new htlc
		k.SetHTLC(ctx, newHTLC, id)
	}

	return nil
}
