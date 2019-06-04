package asset

import (
	"fmt"

	"github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/guardian"
	sdk "github.com/irisnet/irishub/types"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
	ck       bank.Keeper
	gk       guardian.Keeper

	// codespace
	codespace sdk.CodespaceType
	// params subspace
	paramSpace params.Subspace
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, ck bank.Keeper, gk guardian.Keeper, codespace sdk.CodespaceType, paramSpace params.Subspace) Keeper {
	return Keeper{
		storeKey:   key,
		cdc:        cdc,
		ck:         ck,
		gk:         gk,
		codespace:  codespace,
		paramSpace: paramSpace.WithTypeTable(ParamTypeTable()),
	}
}

// return the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// CreateGateway creates a gateway
func (k Keeper) CreateGateway(ctx sdk.Context, msg MsgCreateGateway) (sdk.Tags, sdk.Error) {
	// check if the moniker already exists
	if k.HasGateway(ctx, msg.Moniker) {
		return nil, ErrGatewayAlreadyExists(k.codespace, fmt.Sprintf("moniker already exists:%s", msg.Moniker))
	}

	// get the next gateway id
	gatewayID, err := k.getNewGatewayID(ctx)
	if err != nil {
		return nil, err
	}

	var gateway = Gateway{
		ID:         gatewayID,
		Owner:      msg.Owner,
		Identity:   msg.Identity,
		Moniker:    msg.Moniker,
		Details:    msg.Details,
		Website:    msg.Website,
		RedeemAddr: msg.RedeemAddr,
		Operators:  msg.Operators,
	}

	// set the gateway and related keys
	k.SetGateway(ctx, gateway)
	k.SetMoniker(ctx, msg.Moniker, gatewayID)
	k.SetOwnerGatewayID(ctx, msg.Owner, gatewayID)

	// TODO
	createTags := sdk.NewTags(
		"id", []byte{gatewayID},
		"moniker", []byte(msg.Moniker),
	)

	return createTags, nil
}

// EditGateway edits the specified gateway
func (k Keeper) EditGateway(ctx sdk.Context, msg MsgEditGateway) (sdk.Tags, sdk.Error) {
	// get the destination gateway
	gateway, err := k.GetGatewayByMoniker(ctx, msg.Moniker)
	if err != nil {
		return nil, err
	}

	// check if the given owner matches with the owner of the destination gateway
	if !msg.Owner.Equals(gateway.Owner) {
		return nil, ErrInvalidOwner(k.codespace, fmt.Sprintf("the address %d is not the owner of the gateway %s", msg.Owner, msg.Moniker))
	}

	// update the gateway
	if msg.Identity != nil {
		gateway.Identity = *msg.Identity
	}
	if msg.Details != nil {
		gateway.Details = *msg.Details
	}
	if msg.Website != nil {
		gateway.Website = *msg.Website
	}

	// set the new gateway
	k.SetGateway(ctx, gateway)

	// TODO
	editTags := sdk.NewTags(
		"id", []byte{gateway.ID},
		"moniker", []byte(msg.Moniker),
	)

	return editTags, nil
}

// GetGateway retrieves the gateway of the given id
func (k Keeper) GetGateway(ctx sdk.Context, gatewayID uint8) (Gateway, sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyGateway(gatewayID))
	if bz == nil {
		return Gateway{}, ErrUnkwownGateway(k.codespace, fmt.Sprintf("Unknown gateway id:%d", gatewayID))
	}

	var gateway Gateway
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &gateway)

	return gateway, nil
}

// GetGatewayByMoniker retrieves the gateway of the given moniker
func (k Keeper) GetGatewayByMoniker(ctx sdk.Context, moniker string) (Gateway, sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyMoniker(moniker))
	if bz == nil {
		return Gateway{}, ErrUnkwownGateway(k.codespace, fmt.Sprintf("Unknown gateway moniker:%s", moniker))
	}

	var gatewayID uint8
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &gatewayID)

	return k.GetGateway(ctx, gatewayID)
}

// HasGateway checks if the given gateway exists. Return true if exists, false otherwise
func (k Keeper) HasGateway(ctx sdk.Context, moniker string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(KeyMoniker(moniker))
}

// SetGateway stores the given gateway into underlying storage
func (k Keeper) SetGateway(ctx sdk.Context, gateway Gateway) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(gateway)

	// set KeyGateway
	store.Set(KeyGateway(gateway.ID), bz)
}

// SetMoniker stores the gateway ID into storage by the key KeyMoniker
func (k Keeper) SetMoniker(ctx sdk.Context, moniker string, gatewayID uint8) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(gatewayID)

	// set KeyMoniker
	store.Set(KeyMoniker(moniker), bz)
}

// SetOwnerGatewayID stores the gateway ID into storage by the key KeyOwnerGateway. Intended for iteration on ids of an owner
func (k Keeper) SetOwnerGatewayID(ctx sdk.Context, owner sdk.AccAddress, gatewayID uint8) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(gatewayID)

	// set KeyOwnerGatewayID
	store.Set(KeyOwnerGatewayID(owner, gatewayID), bz)
}

// getNewGatewayID gets the next available gateway ID and increments it
func (k Keeper) getNewGatewayID(ctx sdk.Context) (gatewayID uint8, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyNextGatewayID)
	if bz == nil {
		return 0, ErrInvalidGenesis(k.codespace, "Initial gateway ID never set")
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &gatewayID)
	bz = k.cdc.MustMarshalBinaryLengthPrefixed(gatewayID + 1)
	store.Set(KeyNextGatewayID, bz)

	return gatewayID, nil
}

// setInitialGatewayID sets the initial gateway id in genesis
func (k Keeper) setInitialGatewayID(ctx sdk.Context, gatewayID uint8) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyNextGatewayID)
	if bz != nil {
		return ErrInvalidGenesis(k.codespace, "Initial gateway ID already set")
	}

	bz = k.cdc.MustMarshalBinaryLengthPrefixed(gatewayID)
	store.Set(KeyNextGatewayID, bz)

	return nil
}
