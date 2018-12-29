package service

import (
	sdk "github.com/irisnet/irishub/types"
)

// GenesisState - all service state that must be provided at genesis
type GenesisState struct {
	Params Params `json:"params"` // service params
}

func NewGenesisState(params Params) GenesisState {
	return GenesisState{
		Params: params,
	}
}

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	k.SetParamSet(ctx, data.Params)
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return NewGenesisState(k.GetParamSet(ctx))
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

// get raw genesis raw message for testing
func DefaultGenesisStateForTest() GenesisState {
	return GenesisState{
		Params: DefaultParamsForTest(),
	}
}

// refund deposit from all bindings
// refund service fee from all request
// refund all incoming/return fee
// no process for service fee tax account
func PrepForZeroHeightGenesis(ctx sdk.Context, k Keeper) {
	store := ctx.KVStore(k.storeKey)

	// refund deposit from all bindings
	bindingIterator := sdk.KVStorePrefixIterator(store, bindingPropertyKey)
	defer bindingIterator.Close()
	for ; bindingIterator.Valid(); bindingIterator.Next() {
		var binding SvcBinding
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bindingIterator.Value(), &binding)
		k.ck.SendCoins(ctx, DepositedCoinsAccAddr, binding.Provider, binding.Deposit)
	}

	// refund service fee from all active request
	requestIterator := sdk.KVStorePrefixIterator(store, activeRequestKey)
	defer requestIterator.Close()
	for ; requestIterator.Valid(); requestIterator.Next() {
		var request SvcRequest
		k.cdc.MustUnmarshalBinaryLengthPrefixed(requestIterator.Value(), &request)
		k.ck.SendCoins(ctx, RequestCoinsAccAddr, request.Consumer, request.ServiceFee)
	}

	// refund all incoming fee
	incomingFeeIterator := sdk.KVStorePrefixIterator(store, incomingFeeKey)
	defer incomingFeeIterator.Close()
	for ; incomingFeeIterator.Valid(); incomingFeeIterator.Next() {
		var incomingFee IncomingFee
		k.cdc.MustUnmarshalBinaryLengthPrefixed(incomingFeeIterator.Value(), &incomingFee)
		k.ck.SendCoins(ctx, RequestCoinsAccAddr, incomingFee.Address, incomingFee.Coins)
	}

	// refund all return fee
	returnedFeeIterator := sdk.KVStorePrefixIterator(store, returnedFeeKey)
	defer returnedFeeIterator.Close()
	for ; returnedFeeIterator.Valid(); returnedFeeIterator.Next() {
		var returnedFee ReturnedFee
		k.cdc.MustUnmarshalBinaryLengthPrefixed(returnedFeeIterator.Value(), &returnedFee)
		k.ck.SendCoins(ctx, RequestCoinsAccAddr, returnedFee.Address, returnedFee.Coins)
	}
}

// ValidateGenesis validates the provided service genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	err := validateParams(data.Params)
	if err != nil {
		return err
	}
	return nil
}
