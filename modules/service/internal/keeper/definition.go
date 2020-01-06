package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irishub/modules/service/internal/types"
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
		return sdkerrors.Wrapf(types.ErrUnknownSvcDef, "name: %s", name)
	}

	svcDef := types.NewServiceDefinition(name, description, tags, author, authorDescription, schemas)
	k.SetServiceDefinition(ctx, svcDef)

	return nil
}

// SetServiceDefinition sets the service definition
func (k Keeper) SetServiceDefinition(ctx sdk.Context, svcDef types.ServiceDefinition) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(svcDef)
	store.Set(types.GetServiceDefinitionKey(svcDef.Name), bz)
}

// GetServiceDefinition retrieves a service definition of the specified name
func (k Keeper) GetServiceDefinition(ctx sdk.Context, name string) (svcDef types.ServiceDefinition, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetServiceDefinitionKey(name))
	if bz == nil {
		return svcDef, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &svcDef)
	return svcDef, true
}
