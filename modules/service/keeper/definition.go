package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/service/types"
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
) error {
	if _, found := k.GetServiceDefinition(ctx, name); found {
		return sdkerrors.Wrap(types.ErrServiceDefinitionExists, name)
	}

	svcDef := types.NewServiceDefinition(name, description, tags, author, authorDescription, schemas)
	k.SetServiceDefinition(ctx, svcDef)

	return nil
}

// SetServiceDefinition sets the service definition
func (k Keeper) SetServiceDefinition(ctx sdk.Context, svcDef types.ServiceDefinition) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(&svcDef)
	store.Set(types.GetServiceDefinitionKey(svcDef.Name), bz)
}

// GetServiceDefinition retrieves a service definition of the specified service name
func (k Keeper) GetServiceDefinition(ctx sdk.Context, serviceName string) (svcDef types.ServiceDefinition, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetServiceDefinitionKey(serviceName))
	if bz == nil {
		return svcDef, false
	}

	k.cdc.MustUnmarshalBinaryBare(bz, &svcDef)
	return svcDef, true
}

// IterateServiceDefinitions iterates through all service definitions
func (k Keeper) IterateServiceDefinitions(
	ctx sdk.Context,
	op func(definition types.ServiceDefinition) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.ServiceDefinitionKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var definition types.ServiceDefinition
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &definition)

		if stop := op(definition); stop {
			break
		}
	}
}
