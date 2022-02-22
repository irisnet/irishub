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
) types.Denom {

	denom := types.Denom{
		Id:    k.genDenomID(ctx),
		Name:  name,
		Owner: sednder.String(),
		Data:  data,
	}

	// store denom
	k.SetDenom(ctx, denom)

	return denom
}

// IssueMT issues a new MT
func (k Keeper) IssueMT(ctx sdk.Context,
	denomID string,
	amount uint64,
	data []byte,
	recipient sdk.AccAddress,
) types.MT {

	mt := types.NewMT(k.genMTID(ctx), amount, data)

	// store MT
	k.SetMT(ctx, denomID, mt)

	// increase denom supply
	k.IncreaseDenomSupply(ctx, denomID)

	// increase MT supply
	k.IncreaseMTSupply(ctx, denomID, mt.GetID(), amount)

	// mint amounts to the recipient
	k.AddBalance(ctx, denomID, mt.GetID(), amount, recipient)

	return mt
}

// MintMT mints amounts of an existing MT
func (k Keeper) MintMT(ctx sdk.Context,
	denomID, mtID string,
	amount uint64,
	recipient sdk.AccAddress,
) {

	// increase MT supply
	k.IncreaseMTSupply(ctx, denomID, mtID, amount)

	// mint amounts to the recipient
	k.AddBalance(ctx, denomID, mtID, amount, recipient)
}

// EditMT updates an existing MT
func (k Keeper) EditMT(ctx sdk.Context,
	denomID, mtID string,
	metadata []byte,
	sender sdk.AccAddress,
) error {
	mt, err := k.GetMT(ctx, denomID, mtID)
	if err != nil {
		return err
	}

	if types.Modified(string(metadata)) {
		newMT := types.NewMT(mt.GetID(), mt.GetSupply(), metadata)
		k.SetMT(ctx, denomID, newMT)
	}

	return nil
}

// TransferOwner transfers the ownership of the given MT to the new owner
func (k Keeper) TransferOwner(ctx sdk.Context,
	denomID, mtID string,
	amount uint64,
	srcOwner, dstOwner sdk.AccAddress,
) error {

	srcOwnerAmount := k.GetBalance(ctx, denomID, mtID, srcOwner)
	if srcOwnerAmount < amount {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "insufficient balance: %d", srcOwnerAmount)
	}

	k.Transfer(ctx, denomID, mtID, amount, srcOwner, dstOwner)
	return nil
}

// BurnMT burn amounts of MT from a owner
func (k Keeper) BurnMT(ctx sdk.Context,
	denomID, mtID string,
	amount uint64,
	owner sdk.AccAddress) error {

	srcOwnerAmount := k.GetBalance(ctx, denomID, mtID, owner)
	if srcOwnerAmount < amount {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "insufficient balance: %d", srcOwnerAmount)
	}

	// sub balance
	k.SubBalance(ctx, denomID, mtID, amount, owner)

	// sub total supply
	k.decreaseMTSupply(ctx, denomID, mtID, amount)

	return nil
}

// TransferDenomOwner transfers the ownership of the given denom to the new owner
func (k Keeper) TransferDenomOwner(
	ctx sdk.Context, denomID string, srcOwner, dstOwner sdk.AccAddress,
) error {

	// authorize
	if err := k.Authorize(ctx, denomID, srcOwner); err != nil {
		return err
	}

	denom, _ := k.GetDenom(ctx, denomID)
	denom.Owner = dstOwner.String()

	err := k.UpdateDenom(ctx, denom)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) ExportGenesisState(ctx sdk.Context) *types.GenesisState {

	var collections []types.Collection
	denoms := k.GetDenoms(ctx)
	for _, d := range denoms {
		mts := k.GetMTs(ctx, d.Id)
		collections = append(collections, types.NewCollection(d, mts))
	}

	owners := k.getBalances(ctx)

	return types.NewGenesisState(collections, owners)
}
