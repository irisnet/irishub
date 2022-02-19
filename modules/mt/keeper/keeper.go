package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/mt/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey storetypes.StoreKey // Unexposed key to access store from sdk.Context
	cdc      codec.Codec
}

// NewKeeper creates a new instance of the MT Keeper
func NewKeeper(cdc codec.Codec, storeKey storetypes.StoreKey) Keeper {
	return Keeper{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("irismod/%s", types.ModuleName))
}

// IssueDenom issues a denom according to the given params
func (k Keeper) IssueDenom(ctx sdk.Context,
	name string, sednder sdk.AccAddress, data []byte,
) (types.Denom, error) {

	denom := types.Denom{
		Id:    k.genDenomID(ctx),
		Name:  name,
		Owner: sednder.String(),
		Data:  data,
	}

	return denom, k.SetDenom(ctx, denom)
}

// MintMT mints an MT and manages the MT's existence within Collections and Owners
func (k Keeper) MintMT(ctx sdk.Context,
	denomID, tokenID string,
	amout uint64,
	data []byte,
	owner sdk.AccAddress,
) error {

	if k.HasMT(ctx, denomID, tokenID) {
		mt, err := k.GetMT(ctx, denomID, tokenID)
		if err != nil {
			return err
		}

		k.setMT(
			ctx, denomID,
			types.NewMT(
				tokenID,
				amout+mt.GetSupply(),
				owner,
				data,
			),
		)
	} else {
		k.setMT(
			ctx, denomID,
			types.NewMT(
				tokenID,
				amout,
				owner,
				data,
			),
		)

		k.setOwner(ctx, denomID, tokenID, amout, owner)
		// todo 确定是否还需要 collection
		k.increaseSupply(ctx, denomID)
	}

	return nil
}

// EditMT updates an existing MT
func (k Keeper) EditMT(ctx sdk.Context,
	denomID, tokenID string,
	tokenData []byte,
	owner sdk.AccAddress,
) error {
	denom, found := k.GetDenom(ctx, denomID)
	if !found {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "Denom not found: %s", denomID)
	}

	if denom.Owner != owner.String() {
		return sdkerrors.Wrapf(types.ErrUnauthorized, "Denom is owned by %s", denom.Owner)
	}

	mt, err := k.Authorize(ctx, denomID, tokenID, owner)
	if err != nil {
		return err
	}

	if types.Modified(string(tokenData)) {
		mt.Data = tokenData
	}

	k.setMT(ctx, denomID, mt)

	return nil
}

// TransferOwner transfers the ownership of the given MT to the new owner
func (k Keeper) TransferOwner(ctx sdk.Context,
	denomID, tokenID string,
	amount uint64,
	srcOwner, dstOwner sdk.AccAddress,
) error {
	_, found := k.GetDenom(ctx, denomID)
	if !found {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "Denom not found: %s", denomID)
	}

	srcOwnerAmount := k.getOwner(ctx, denomID, tokenID, srcOwner)
	if srcOwnerAmount < amount {
		return sdkerrors.Wrapf(types.ErrInvalidCollection, "Lack of mt: %", srcOwnerAmount)
	}

	k.swapOwner(ctx, denomID, tokenID, amount, srcOwner, dstOwner)
	return nil
}

// BurnMT deletes a specified MT
func (k Keeper) BurnMT(ctx sdk.Context,
	denomID, tokenID string,
	amount uint64,
	owner sdk.AccAddress) error {
	if !k.HasDenomID(ctx, denomID) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "Denom not found: %s", denomID)
	}

	mt, err := k.Authorize(ctx, denomID, tokenID, owner)
	if err != nil {
		return err
	}

	srcOwnerAmount := k.getOwner(ctx, denomID, tokenID, owner)
	if srcOwnerAmount < amount {
		return sdkerrors.Wrapf(types.ErrInvalidCollection, "Lack of mt: %", srcOwnerAmount)
	}

	k.deleteOwner(ctx, denomID, tokenID, amount, owner)
	k.setMT(ctx, denomID, types.MT{
		Id:     mt.Id,
		Supply: mt.Supply - amount,
		Data:   mt.Data,
		Owner:  mt.Owner,
	})

	return nil
}

// TransferDenomOwner transfers the ownership of the given denom to the new owner
func (k Keeper) TransferDenomOwner(
	ctx sdk.Context, denomID string, srcOwner, dstOwner sdk.AccAddress,
) error {
	denom, found := k.GetDenom(ctx, denomID)
	if !found {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "Denom not found: %s", denomID)
	}

	// authorize
	if srcOwner.String() != denom.Owner {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to transfer denom %s", srcOwner.String(), denomID)
	}

	denom.Owner = dstOwner.String()

	err := k.UpdateDenom(ctx, denom)
	if err != nil {
		return err
	}

	return nil
}
