package iservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *wire.Codec
}

func NewKeeper(cdc *wire.Codec, key sdk.StoreKey) Keeper {
	keeper := Keeper{
		storeKey: key,
		cdc:      cdc,
	}
	return keeper
}

func (k Keeper) AddServiceDefinition(ctx sdk.Context, serviceDef MsgSvcDef) {
	kvStore := ctx.KVStore(k.storeKey)

	serviceDefBytes, err := k.cdc.MarshalBinary(serviceDef)
	if err != nil {
		panic(err)
	}

	kvStore.Set(GetServiceDefinitionKey(serviceDef.ChainId, serviceDef.Name), serviceDefBytes)
}

func (k Keeper) GetServiceDefinition(ctx sdk.Context, chainId, name string) *MsgSvcDef {
	kvStore := ctx.KVStore(k.storeKey)

	serviceDefBytes := kvStore.Get(GetServiceDefinitionKey(chainId, name))
	if serviceDefBytes != nil {
		var serviceDef MsgSvcDef
		err := k.cdc.UnmarshalBinary(serviceDefBytes, &serviceDef)
		if err != nil {
			panic(err)
		}

		return &serviceDef

	}
	return nil
}
