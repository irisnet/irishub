package iservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/irisnet/irishub/tools/protoidl"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec

	// codespace
	codespace sdk.CodespaceType
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, codespace sdk.CodespaceType) Keeper {
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
	for index, method := range methods {
		methodProperty, err := methodToMethodProperty(index+1, method)
		if err != nil {
			return err
		}
		methodBytes := k.cdc.MustMarshalBinary(methodProperty)
		kvStore.Set(GetMethodPropertyKey(serviceDef.ChainId, serviceDef.Name, methodProperty.ID), methodBytes)
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

func (k Keeper) AddServiceBinding(ctx sdk.Context, svcBinding SvcBinding) {
	kvStore := ctx.KVStore(k.storeKey)

	svcBindingBytes, err := k.cdc.MarshalBinary(svcBinding)
	if err != nil {
		panic(err)
	}

	kvStore.Set(GetServiceBindingKey(svcBinding.DefChainID, svcBinding.DefName, svcBinding.BindChainID, svcBinding.Provider), svcBindingBytes)
}

func (k Keeper) GetServiceBinding(ctx sdk.Context, defChainID, defName, bindChainID string, provider sdk.AccAddress) (svcBinding SvcBinding, found bool) {
	kvStore := ctx.KVStore(k.storeKey)

	svcBindingBytes := kvStore.Get(GetServiceBindingKey(defChainID, defName, bindChainID, provider))
	if svcBindingBytes != nil {
		var svcBinding SvcBinding
		k.cdc.MustUnmarshalBinary(svcBindingBytes, &svcBinding)
		return svcBinding, true
	}
	return svcBinding, false
}
