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
	id, name string, sednder sdk.AccAddress, data []byte,
) error {
	return k.SetDenom(ctx, types.Denom{
		Id:               id,
		Name:             name,
		Owner:            sednder.String(),
		Data:             data,
	})
}

// MintMT mints an MT and manages the MT's existence within Collections and Owners
func (k Keeper) MintMT(
	ctx sdk.Context, denomID, tokenID string, amout uint64, data []byte, owner sdk.AccAddress,
) error {

	// TODO create if not exists, mint if exists
	if k.HasMT(ctx, denomID, tokenID) {
		return sdkerrors.Wrapf(types.ErrMTAlreadyExists, "MT %s already exists in collection %s", tokenID, denomID)
	}

	k.setMT(
		ctx, denomID,
		types.NewMT(
			tokenID,
			amout,
			owner,
			data,
		),
	)
	k.setOwner(ctx, denomID, tokenID, owner)
	k.increaseSupply(ctx, denomID)

	return nil
}

// EditMT updates an existing MT
func (k Keeper) EditMT(
	ctx sdk.Context, denomID, tokenID, tokenNm,
	tokenURI, tokenURIHash, tokenData string, owner sdk.AccAddress,
) error {
	_, found := k.GetDenom(ctx, denomID)
	if !found {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "Denom not found: %s", denomID)
	}

	// just the owner of MT can edit
	mt, err := k.Authorize(ctx, denomID, tokenID, owner)
	if err != nil {
		return err
	}

	if types.Modified(tokenData) {
		mt.Data = []byte(tokenData)
	}

	k.setMT(ctx, denomID, mt)

	return nil
}

// TODO add amount
// TransferOwner transfers the ownership of the given MT to the new owner
func (k Keeper) TransferOwner(
	ctx sdk.Context, denomID, tokenID string, srcOwner, dstOwner sdk.AccAddress,
) error {
	_, found := k.GetDenom(ctx, denomID)
	if !found {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "Denom not found: %s", denomID)
	}

	mt, err := k.Authorize(ctx, denomID, tokenID, srcOwner)
	if err != nil {
		return err
	}

	mt.Owner = dstOwner.String()

	k.setMT(ctx, denomID, mt)
	k.swapOwner(ctx, denomID, tokenID, srcOwner, dstOwner)
	return nil
}

// TODO add amount
// BurnMT deletes a specified MT
func (k Keeper) BurnMT(ctx sdk.Context, denomID, tokenID string, owner sdk.AccAddress) error {
	if !k.HasDenomID(ctx, denomID) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "Denom not found: %s", denomID)
	}

	mt, err := k.Authorize(ctx, denomID, tokenID, owner)
	if err != nil {
		return err
	}

	k.deleteMT(ctx, denomID, mt)
	k.deleteOwner(ctx, denomID, tokenID, owner)
	k.decreaseSupply(ctx, denomID)

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
