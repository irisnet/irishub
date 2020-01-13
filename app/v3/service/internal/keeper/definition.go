package keeper

import (
	"github.com/irisnet/irishub/app/v3/service/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

// AddServiceDefinition creates a new service definition
func (k Keeper) AddServiceDefinition(
	ctx sdk.Context,
	name,
	description string,
	tags []string,
	author sdk.AccAddress,
	authorDescription,
	schemas string,
) sdk.Error {
	if _, found := k.GetServiceDefinition(ctx, name); found {
		return types.ErrServiceDefinitionExists(k.codespace, name)
	}

	svcDef := types.NewServiceDefinition(name, description, tags, author, authorDescription, schemas)
	k.SetServiceDefinition(ctx, svcDef)

	return nil
}

// SetServiceDefinition sets the service definition
func (k Keeper) SetServiceDefinition(ctx sdk.Context, svcDef types.ServiceDefinition) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(svcDef)
	store.Set(GetServiceDefinitionKey(svcDef.Name), bz)
}

// GetServiceDefinition retrieves a service definition of the specified name
func (k Keeper) GetServiceDefinition(ctx sdk.Context, name string) (svcDef types.ServiceDefinition, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(GetServiceDefinitionKey(name))
	if bz == nil {
		return svcDef, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &svcDef)
	return svcDef, true
}
