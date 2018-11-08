package iservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/irisnet/irishub/tools/protoidl"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"fmt"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *wire.Codec
	ck       bank.Keeper

	// codespace
	codespace sdk.CodespaceType
}

func NewKeeper(cdc *wire.Codec, key sdk.StoreKey, ck bank.Keeper, codespace sdk.CodespaceType) Keeper {
	keeper := Keeper{
		storeKey:  key,
		cdc:       cdc,
		ck:        ck,
		codespace: codespace,
	}
	return keeper
}

// return the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

func (k Keeper) AddServiceDefinition(ctx sdk.Context, svcDef SvcDef) {
	kvStore := ctx.KVStore(k.storeKey)

	svcDefBytes, err := k.cdc.MarshalBinary(svcDef)
	if err != nil {
		panic(err)
	}

	kvStore.Set(GetServiceDefinitionKey(svcDef.ChainId, svcDef.Name), svcDefBytes)
}

func (k Keeper) AddMethods(ctx sdk.Context, svcDef SvcDef) sdk.Error {
	methods, err := protoidl.GetMethods(svcDef.IDLContent)
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
		kvStore.Set(GetMethodPropertyKey(svcDef.ChainId, svcDef.Name, methodProperty.ID), methodBytes)
	}
	return nil
}

func (k Keeper) GetServiceDefinition(ctx sdk.Context, chainId, name string) (svcDef SvcDef, found bool) {
	kvStore := ctx.KVStore(k.storeKey)

	serviceDefBytes := kvStore.Get(GetServiceDefinitionKey(chainId, name))
	if serviceDefBytes != nil {
		var serviceDef SvcDef
		k.cdc.MustUnmarshalBinary(serviceDefBytes, &serviceDef)
		return serviceDef, true
	}
	return svcDef, false
}

// Gets all the methods in a specific service
func (k Keeper) GetMethods(ctx sdk.Context, chainId, name string) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, GetMethodsSubspaceKey(chainId, name))
}

func (k Keeper) AddServiceBinding(ctx sdk.Context, svcBinding SvcBinding) (sdk.Error, bool) {
	kvStore := ctx.KVStore(k.storeKey)
	_, found := k.GetServiceDefinition(ctx, svcBinding.DefChainID, svcBinding.DefName)
	if !found {
		return ErrSvcDefNotExists(k.Codespace(), svcBinding.DefChainID, svcBinding.DefName), false
	}

	_, found = k.GetServiceBinding(ctx, svcBinding.DefChainID, svcBinding.DefName, svcBinding.BindChainID, svcBinding.Provider)
	if found {
		return ErrSvcBindingExists(k.Codespace()), false
	}

	if !svcBinding.Deposit.IsGTE(iserviceParams.MinProviderDeposit) {
		return ErrLtMinProviderDeposit(k.Codespace(), iserviceParams.MinProviderDeposit), false
	}

	err := k.ValidateMethodPrices(ctx, svcBinding)
	if err != nil {
		return err, false
	}

	// Subtract coins from provider's account
	_, _, err = k.ck.SubtractCoins(ctx, svcBinding.Provider, svcBinding.Deposit)
	if err != nil {
		return err, false
	}

	svcBindingBytes := k.cdc.MustMarshalBinary(svcBinding)
	kvStore.Set(GetServiceBindingKey(svcBinding.DefChainID, svcBinding.DefName, svcBinding.BindChainID, svcBinding.Provider), svcBindingBytes)
	return nil, true
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

func (k Keeper) UpdateServiceBinding(ctx sdk.Context, svcBinding SvcBinding) (sdk.Error, bool) {
	kvStore := ctx.KVStore(k.storeKey)
	oldBinding, found := k.GetServiceBinding(ctx, svcBinding.DefChainID, svcBinding.DefName, svcBinding.BindChainID, svcBinding.Provider)
	if !found {
		return ErrSvcBindingNotExists(k.Codespace()), false
	}

	if len(svcBinding.Prices) > 0 {
		err := k.ValidateMethodPrices(ctx, svcBinding)
		if err != nil {
			return err, false
		}
		oldBinding.Prices = svcBinding.Prices
	}

	// can't update type Global to Local
	if oldBinding.BindingType == Global && svcBinding.BindingType == Local {
		return ErrInvalidUpdate(k.Codespace(), "can't update binding type from Global to Local"), false
	}

	if svcBinding.BindingType != 0x00 {
		oldBinding.BindingType = svcBinding.BindingType
	}

	// Add coins to svcBinding deposit
	if svcBinding.Deposit.IsNotNegative() {
		oldBinding.Deposit = oldBinding.Deposit.Plus(svcBinding.Deposit)
	}

	if !oldBinding.Deposit.IsGTE(iserviceParams.MinProviderDeposit) {
		return ErrLtMinProviderDeposit(k.Codespace(), iserviceParams.MinProviderDeposit.Minus(oldBinding.Deposit)), false
	}

	// Subtract coins from provider's account
	_, _, err := k.ck.SubtractCoins(ctx, svcBinding.Provider, svcBinding.Deposit)
	if err != nil {
		return err, false
	}

	if svcBinding.Expiration != 0 {
		height := ctx.BlockHeader().Height
		if svcBinding.Expiration >= 0 && svcBinding.Expiration < height {
			oldBinding.Expiration = height
		} else {
			oldBinding.Expiration = svcBinding.Expiration
		}
	}
	if svcBinding.Level.UsableTime != 0 {
		oldBinding.Level.UsableTime = svcBinding.Level.UsableTime
	}
	if svcBinding.Level.AvgRspTime != 0 {
		oldBinding.Level.AvgRspTime = svcBinding.Level.AvgRspTime
	}

	svcBindingBytes := k.cdc.MustMarshalBinary(oldBinding)
	kvStore.Set(GetServiceBindingKey(svcBinding.DefChainID, svcBinding.DefName, svcBinding.BindChainID, svcBinding.Provider), svcBindingBytes)
	return nil, true
}

func (k Keeper) RefundDeposit(ctx sdk.Context, defChainID, defName, bindChainID string, provider sdk.AccAddress) (sdk.Error, bool) {
	kvStore := ctx.KVStore(k.storeKey)
	binding, found := k.GetServiceBinding(ctx, defChainID, defName, bindChainID, provider)
	if !found {
		return ErrSvcBindingNotExists(k.Codespace()), false
	}

	if binding.Expiration < 0 {
		return ErrRefundDeposit(k.Codespace(), "service binding don`t set expiration height"), false
	}

	if binding.Deposit.IsZero() {
		return ErrRefundDeposit(k.Codespace(), "service binding deposit is zero"), false
	}

	height := ctx.BlockHeader().Height
	refundHeight := binding.Expiration + int64(iserviceParams.MaxRequestTimeout)
	if refundHeight >= height {
		return ErrRefundDeposit(k.Codespace(), fmt.Sprintf("you can refund deposit util block height greater than %d", refundHeight)), false
	}

	// Add coins to provider's account
	_, _, err := k.ck.AddCoins(ctx, binding.Provider, binding.Deposit)
	if err != nil {
		return err, false
	}

	binding.Deposit = sdk.Coins{}

	svcBindingBytes := k.cdc.MustMarshalBinary(binding)
	kvStore.Set(GetServiceBindingKey(binding.DefChainID, binding.DefName, binding.BindChainID, binding.Provider), svcBindingBytes)
	return nil, true
}

func (k Keeper) ValidateMethodPrices(ctx sdk.Context, svcBinding SvcBinding) sdk.Error {
	methodIterator := k.GetMethods(ctx, svcBinding.DefChainID, svcBinding.DefName)
	var methods []MethodProperty
	for ; methodIterator.Valid(); methodIterator.Next() {
		var method MethodProperty
		k.cdc.MustUnmarshalBinary(methodIterator.Value(), &method)
		methods = append(methods, method)
	}

	if len(methods) != len(svcBinding.Prices) {
		return ErrInvalidPriceCount(k.Codespace(), len(svcBinding.Prices), len(methods))
	}
	return nil
}
