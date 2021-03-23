package keeper

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"time"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/htlc/types"
)

// CreateHTLC creates an HTLC
func (k Keeper) CreateHTLC(
	ctx sdk.Context,
	sender sdk.AccAddress,
	to sdk.AccAddress,
	receiverOnOtherChain string,
	senderOnOtherChain string,
	amount sdk.Coins,
	hashLock tmbytes.HexBytes,
	timestamp uint64,
	timeLock uint64,
	transfer bool,
) (
	id tmbytes.HexBytes,
	err error,
) {
	id = types.GetID(sender, to, amount, hashLock)

	// check if the HTLC already exists
	if k.HasHTLC(ctx, id) {
		return id, sdkerrors.Wrap(types.ErrHTLCExists, id.String())
	}

	expirationHeight := uint64(ctx.BlockHeight()) + timeLock

	var direction types.SwapDirection
	if transfer {
		// create HTLT
		if direction, err = k.createHTLT(
			ctx, sender, to, receiverOnOtherChain, senderOnOtherChain,
			amount, hashLock, timestamp, timeLock,
		); err != nil {
			return id, err
		}
	} else {
		// create HTLT
		if err = k.createHTLC(ctx, sender, amount); err != nil {
			return id, err
		}
	}

	htlc := types.NewHTLC(
		id, sender, to, receiverOnOtherChain,
		senderOnOtherChain, amount, hashLock,
		nil, timestamp, expirationHeight,
		types.Open, 0, transfer, direction,
	)

	// set the HTLC
	k.SetHTLC(ctx, htlc, id)

	// add to the expiration queue
	k.AddHTLCToExpiredQueue(ctx, htlc.ExpirationHeight, id)

	return id, nil
}

func (k Keeper) createHTLC(
	ctx sdk.Context,
	sender sdk.AccAddress,
	amount sdk.Coins,
) error {
	// transfer the specified tokens to the HTLC module account
	return k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, amount)
}

func (k Keeper) createHTLT(
	ctx sdk.Context,
	sender sdk.AccAddress,
	to sdk.AccAddress,
	receiverOnOtherChain string,
	senderOnOtherChain string,
	amount sdk.Coins,
	hashLock tmbytes.HexBytes,
	timestamp uint64,
	timeLock uint64,
) (
	types.SwapDirection,
	error,
) {
	var direction types.SwapDirection

	if len(amount) != 1 {
		return direction, sdkerrors.Wrapf(types.ErrInvalidAmount, amount.String())
	}

	asset, err := k.GetAsset(ctx, amount[0].Denom)
	if err != nil {
		return direction, err
	}

	if err = k.ValidateLiveAsset(ctx, amount[0]); err != nil {
		return direction, err
	}

	// Swap amount must be within the specified swap amount limits
	if amount[0].Amount.LT(asset.MinSwapAmount) || amount[0].Amount.GT(asset.MaxSwapAmount) {
		return direction, sdkerrors.Wrapf(types.ErrInvalidAmount, "amount %d outside range [%s, %s]", amount[0].Amount, asset.MinSwapAmount, asset.MaxSwapAmount)
	}

	// Unix timestamp must be in range [-15 mins, 30 mins) of the current time
	pastTimestampLimit := ctx.BlockTime().Add(-15 * time.Minute).Unix()
	futureTimestampLimit := ctx.BlockTime().Add(30 * time.Minute).Unix()
	if timestamp < uint64(pastTimestampLimit) || timestamp >= uint64(futureTimestampLimit) {
		return direction, sdkerrors.Wrap(
			types.ErrInvalidTimestamp,
			fmt.Sprintf(
				"timestamp can neither be 15 minutes ahead of the current time, nor 30 minutes later. block time: %s, timestamp: %s",
				ctx.BlockTime().String(), time.Unix(int64(timestamp), 0).UTC().String(),
			),
		)
	}

	deputyAddress, _ := sdk.AccAddressFromBech32(asset.DeputyAddress)

	if sender.Equals(deputyAddress) {
		if to.Equals(deputyAddress) {
			return direction, sdkerrors.Wrapf(types.ErrInvalidAccount, "deputy cannot be both sender and receiver: %s", asset.DeputyAddress)
		}
		direction = types.Incoming
	} else {
		if !to.Equals(deputyAddress) {
			return direction, sdkerrors.Wrapf(types.ErrInvalidAccount, "deputy must be recipient for outgoing account: %s", to)
		}
		direction = types.Outgoing
	}

	switch direction {
	case types.Incoming:
		// If recipient's account doesn't exist, register it in state so that the address can send
		// a claim swap tx without needing to be registered in state by receiving a coin transfer.
		recipientAcc := k.accountKeeper.GetAccount(ctx, deputyAddress)
		if recipientAcc == nil {
			acc := k.accountKeeper.NewAccountWithAddress(ctx, deputyAddress)
			k.accountKeeper.SetAccount(ctx, acc)
		}
		// Incoming swaps have already had their fees collected by the deputy during the relay process.
		if err := k.IncrementIncomingAssetSupply(ctx, amount[0]); err != nil {
			return direction, err
		}
	case types.Outgoing:
		// Outgoing swaps must have a time lock within the accepted range
		if timeLock < asset.MinBlockLock || timeLock > asset.MaxBlockLock {
			return direction, sdkerrors.Wrapf(types.ErrInvalidTimeLock, "time lock %d outside range [%d, %d]", timeLock, asset.MinBlockLock, asset.MaxBlockLock)
		}
		// Amount in outgoing swaps must be able to pay the deputy's fixed fee.
		if amount[0].Amount.LT(asset.FixedFee.Add(asset.MinSwapAmount)) {
			return direction, sdkerrors.Wrapf(
				types.ErrInsufficientAmount,
				"amount %s < fixed fee %s + min swap amount %s",
				amount[0].String(), asset.FixedFee.String(), asset.MinSwapAmount,
			)
		}
		if err := k.IncrementOutgoingAssetSupply(ctx, amount[0]); err != nil {
			return direction, err
		}
		// Transfer coins to module - only needed for outgoing swaps
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, amount); err != nil {
			return direction, err
		}
	default:
		return direction, sdkerrors.Wrapf(types.ErrInvalidDirection, direction.String())
	}

	return direction, nil
}

// ClaimHTLC claims the specified HTLC with the given secret
func (k Keeper) ClaimHTLC(
	ctx sdk.Context,
	id tmbytes.HexBytes,
	secret tmbytes.HexBytes,
) (
	string,
	bool,
	types.SwapDirection,
	error,
) {
	// query the HTLC
	htlc, found := k.GetHTLC(ctx, id)
	if !found {
		return "", false, types.None, sdkerrors.Wrap(types.ErrUnknownHTLC, id.String())
	}

	// check if the HTLC is open
	if htlc.State != types.Open {
		return "", false, types.None, sdkerrors.Wrap(types.ErrHTLCNotOpen, id.String())
	}

	hashLock, _ := hex.DecodeString(htlc.HashLock)

	// check if the secret matches with the hash lock
	if !bytes.Equal(types.GetHashLock(secret, htlc.Timestamp), hashLock) {
		return "", false, types.None, sdkerrors.Wrap(types.ErrInvalidSecret, secret.String())
	}

	to, err := sdk.AccAddressFromBech32(htlc.To)
	if err != nil {
		return "", false, types.None, err
	}

	if htlc.Transfer {
		if err := k.claimHTLT(ctx, htlc); err != nil {
			return "", false, types.None, err
		}
	} else {
		if err := k.claimHTLC(ctx, htlc.Amount, to); err != nil {
			return "", false, types.None, err
		}
	}

	// update the secret and state of the HTLC
	htlc.Secret = secret.String()
	htlc.State = types.Completed
	htlc.ClosedBlock = uint64(ctx.BlockHeight())
	k.SetHTLC(ctx, htlc, id)

	// delete from the expiration queue
	k.DeleteHTLCFromExpiredQueue(ctx, htlc.ExpirationHeight, id)

	return htlc.HashLock, htlc.Transfer, htlc.Direction, nil
}

func (k Keeper) claimHTLC(ctx sdk.Context, amount sdk.Coins, to sdk.AccAddress) error {
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, to, amount)
}

func (k Keeper) claimHTLT(ctx sdk.Context, htlc types.HTLC) error {
	switch htlc.Direction {
	case types.Incoming:
		if err := k.DecrementIncomingAssetSupply(ctx, htlc.Amount[0]); err != nil {
			return err
		}
		if err := k.IncrementCurrentAssetSupply(ctx, htlc.Amount[0]); err != nil {
			return err
		}
		// incoming case - coins should be MINTED, then sent to user
		if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, htlc.Amount); err != nil {
			return err
		}
		// Send intended recipient coins
		toAddr, _ := sdk.AccAddressFromBech32(htlc.To)
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, toAddr, htlc.Amount); err != nil {
			return err
		}
	case types.Outgoing:
		if err := k.DecrementOutgoingAssetSupply(ctx, htlc.Amount[0]); err != nil {
			return err
		}
		if err := k.DecrementCurrentAssetSupply(ctx, htlc.Amount[0]); err != nil {
			return err
		}
		// outgoing case  - coins should be burned
		if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, htlc.Amount); err != nil {
			return err
		}
	default:
		return sdkerrors.Wrapf(types.ErrInvalidDirection, htlc.Direction.String())
	}

	return nil
}

// RefundHTLC refunds the specified HTLC
func (k Keeper) RefundHTLC(ctx sdk.Context, h types.HTLC, id tmbytes.HexBytes) error {
	sender, err := sdk.AccAddressFromBech32(h.Sender)
	if err != nil {
		return err
	}

	if h.Transfer {
		if err := k.refundHTLT(ctx, h.Direction, sender, h.Amount); err != nil {
			return err
		}
	} else {
		if err := k.refundHTLC(ctx, sender, h.Amount); err != nil {
			return err
		}
	}

	// update the state of the HTLC
	h.State = types.Refunded
	h.ClosedBlock = uint64(ctx.BlockHeight())
	k.SetHTLC(ctx, h, id)

	return nil
}

func (k Keeper) refundHTLC(ctx sdk.Context, sender sdk.AccAddress, amount sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, amount)
}

func (k Keeper) refundHTLT(ctx sdk.Context, direction types.SwapDirection, sender sdk.AccAddress, amount sdk.Coins) error {
	switch direction {
	case types.Incoming:
		if err := k.DecrementIncomingAssetSupply(ctx, amount[0]); err != nil {
			return err
		}
	case types.Outgoing:
		if err := k.DecrementOutgoingAssetSupply(ctx, amount[0]); err != nil {
			return err
		}
		// Refund coins to original swap sender for outgoing swaps
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, amount); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid direction")
	}
	return nil
}

// HasHTLC checks if the given HTLC exists
func (k Keeper) HasHTLC(ctx sdk.Context, id tmbytes.HexBytes) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetHTLCKey(id))
}

// SetHTLC sets the given HTLC
func (k Keeper) SetHTLC(ctx sdk.Context, htlc types.HTLC, id tmbytes.HexBytes) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(&htlc)
	store.Set(types.GetHTLCKey(id), bz)
}

// GetHTLC retrieves the specified HTLC
func (k Keeper) GetHTLC(ctx sdk.Context, id tmbytes.HexBytes) (htlc types.HTLC, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetHTLCKey(id))
	if bz == nil {
		return htlc, false
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &htlc)
	return htlc, true
}

// AddHTLCToExpiredQueue adds the specified HTLC to the expiration queue
func (k Keeper) AddHTLCToExpiredQueue(ctx sdk.Context, expirationHeight uint64, id tmbytes.HexBytes) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetHTLCExpiredQueueKey(expirationHeight, id), []byte{})
}

// DeleteHTLCFromExpiredQueue removes the specified HTLC from the expiration queue
func (k Keeper) DeleteHTLCFromExpiredQueue(ctx sdk.Context, expirationHeight uint64, id tmbytes.HexBytes) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetHTLCExpiredQueueKey(expirationHeight, id))
}

// IterateHTLCs iterates through the HTLCs
func (k Keeper) IterateHTLCs(
	ctx sdk.Context,
	op func(id tmbytes.HexBytes, h types.HTLC) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.HTLCKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		id := tmbytes.HexBytes(iterator.Key()[1:])

		var htlc types.HTLC
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &htlc)

		if stop := op(id, htlc); stop {
			break
		}
	}
}

// IterateHTLCExpiredQueueByHeight iterates through the HTLC expiration queue by the specified height
func (k Keeper) IterateHTLCExpiredQueueByHeight(
	ctx sdk.Context, height uint64,
	op func(id tmbytes.HexBytes, h types.HTLC) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.GetHTLCExpiredQueueSubspace(height))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		id := tmbytes.HexBytes(iterator.Key()[9:])
		htlc, _ := k.GetHTLC(ctx, id)

		if stop := op(id, htlc); stop {
			break
		}
	}
}
