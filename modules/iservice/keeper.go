package iservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/irisnet/irishub/tools/protoidl"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *wire.Codec

	// codespace
	codespace sdk.CodespaceType
}

func NewKeeper(cdc *wire.Codec, key sdk.StoreKey, codespace sdk.CodespaceType) Keeper {
	keeper := Keeper{
		storeKey:  key,
		cdc:       cdc,
		codespace: codespace,
	}
	return keeper
}

// return the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

func (k Keeper) AddServiceDefinition(ctx sdk.Context, serviceDef MsgSvcDef) {
	kvStore := ctx.KVStore(k.storeKey)

	serviceDefBytes, err := k.cdc.MarshalBinary(serviceDef)
	if err != nil {
		panic(err)
	}

	kvStore.Set(GetServiceDefinitionKey(serviceDef.ChainId, serviceDef.Name), serviceDefBytes)
}

func (k Keeper) AddMethods(ctx sdk.Context, serviceDef MsgSvcDef) sdk.Error {
	methods, err := protoidl.GetMethods(serviceDef.IDLContent)
	if err != nil {
		panic(err)
	}
	kvStore := ctx.KVStore(k.storeKey)
	for _, method := range methods {
		methodProperty, err := methodToMethodProperty(method)
		if err != nil {
			return err
		}
		methodBytes := k.cdc.MustMarshalBinary(methodProperty)
		kvStore.Set(GetMethodPropertyKey(serviceDef.ChainId, serviceDef.Name, method.Name), methodBytes)
	}
	return nil
}

func (k Keeper) GetServiceDefinition(ctx sdk.Context, chainId, name string) (msgSvcDef MsgSvcDef, found bool) {
	kvStore := ctx.KVStore(k.storeKey)

	serviceDefBytes := kvStore.Get(GetServiceDefinitionKey(chainId, name))
	if serviceDefBytes != nil {
		var serviceDef MsgSvcDef
		k.cdc.MustUnmarshalBinary(serviceDefBytes, &serviceDef)
		return serviceDef, true
	}
	return msgSvcDef, false
}

// Gets all the methods in a specific service
func (k Keeper) GetMethods(ctx sdk.Context, chainId, name string) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, GetMethodsSubspaceKey(chainId, name))
}
