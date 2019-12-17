package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/service/internal/types"
)

func (k Keeper) AddServiceDefinition(
	ctx sdk.Context,
	name,
	chainId,
	description string,
	tags []string,
	author sdk.AccAddress,
	authorDescription,
	idlContent string,
) sdk.Error {
	_, found := k.GetServiceDefinition(ctx, chainId, name)
	if found {
		return types.ErrSvcDefExists(k.codespace, chainId, name)
	}

	svcDef := types.NewSvcDef(name, chainId, description, tags, author, authorDescription, idlContent)
	k.SetServiceDefinition(ctx, svcDef)

	return k.AddMethods(ctx, svcDef)
}

func (k Keeper) SetServiceDefinition(ctx sdk.Context, svcDef types.SvcDef) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(svcDef)
	store.Set(types.GetServiceDefinitionKey(svcDef.ChainId, svcDef.Name), bz)
}

// TODO
func (k Keeper) AddMethods(ctx sdk.Context, svcDef types.SvcDef) sdk.Error {
	methods, err := types.ParseMethods(svcDef.IDLContent)
	if err != nil {
		panic(err)
	}

	for index, method := range methods {
		methodProperty, err := types.MethodToMethodProperty(index+1, method)
		if err != nil {
			return err
		}

		k.SetMethod(ctx, svcDef.ChainId, svcDef.Name, methodProperty)
	}

	return nil
}

func (k Keeper) SetMethod(ctx sdk.Context, chainId, svcName string, method types.MethodProperty) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(method)
	store.Set(types.GetMethodPropertyKey(chainId, svcName, method.ID), bz)
}

func (k Keeper) GetServiceDefinition(ctx sdk.Context, chainId, name string) (svcDef types.SvcDef, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetServiceDefinitionKey(chainId, name))
	if bz == nil {
		return svcDef, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &svcDef)
	return svcDef, true
}

// Gets the method in a specific service and methodID
func (k Keeper) GetMethod(ctx sdk.Context, chainId, svcName string, id int16) (method types.MethodProperty, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetMethodPropertyKey(chainId, svcName, id))
	if bz == nil {
		return method, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &method)
	return method, true
}

// Gets all the methods in a specific service
func (k Keeper) GetMethods(ctx sdk.Context, chainId, svcName string) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.GetMethodsSubspaceKey(chainId, svcName))
}
