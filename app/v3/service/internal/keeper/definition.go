package keeper

import (
	"github.com/irisnet/irishub/app/v3/service/internal/types"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/tools/protoidl"
)

// AddServiceDefinition
func (k Keeper) AddServiceDefinition(
	ctx sdk.Context,
	name,
	chainID,
	description string,
	tags []string,
	author sdk.AccAddress,
	authorDescription,
	idlContent string,
) sdk.Error {
	if _, found := k.GetServiceDefinition(ctx, chainID, name); found {
		return types.ErrSvcDefExists(k.codespace, chainID, name)
	}

	svcDef := types.NewSvcDef(name, chainID, description, tags, author, authorDescription, idlContent)
	k.SetServiceDefinition(ctx, svcDef)

	return k.AddMethods(ctx, svcDef)
}

// SetServiceDefinition
func (k Keeper) SetServiceDefinition(ctx sdk.Context, svcDef types.SvcDef) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(svcDef)
	store.Set(GetServiceDefinitionKey(svcDef.ChainID, svcDef.Name), bz)
}

// AddMethods
func (k Keeper) AddMethods(ctx sdk.Context, svcDef types.SvcDef) sdk.Error {
	methods, err := protoidl.GetMethods(svcDef.IDLContent)
	if err != nil {
		panic(err)
	}

	for index, method := range methods {
		methodProperty, err := types.MethodToMethodProperty(index+1, method)
		if err != nil {
			return err
		}

		k.SetMethod(ctx, svcDef.ChainID, svcDef.Name, methodProperty)
	}

	return nil
}

// SetMethod
func (k Keeper) SetMethod(ctx sdk.Context, chainID, svcName string, method types.MethodProperty) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(method)
	store.Set(GetMethodPropertyKey(chainID, svcName, method.ID), bz)
}

// GetServiceDefinition
func (k Keeper) GetServiceDefinition(ctx sdk.Context, chainID, name string) (svcDef types.SvcDef, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(GetServiceDefinitionKey(chainID, name))
	if bz == nil {
		return svcDef, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &svcDef)
	return svcDef, true
}

// GetMethod gets the method in a specific service and methodID
func (k Keeper) GetMethod(ctx sdk.Context, chainID, svcName string, id int16) (method types.MethodProperty, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(GetMethodPropertyKey(chainID, svcName, id))
	if bz == nil {
		return method, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &method)
	return method, true
}

// GetMethods gets all the methods in a specific service
func (k Keeper) GetMethods(ctx sdk.Context, chainID, svcName string) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, GetMethodsSubspaceKey(chainID, svcName))
}
