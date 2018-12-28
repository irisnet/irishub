package service

import (
	"github.com/irisnet/irishub/modules/params"
	"github.com/irisnet/irishub/modules/service/params"
	sdk "github.com/irisnet/irishub/types"
)

// GenesisState - all service state that must be provided at genesis
type GenesisState struct {
	ServiceParams serviceparams.Params `json:service_govparams`
}

func NewGenesisState(maxRequestTimeout int64, minDepositMultiple int64, serviceFeeTax, slashFraction sdk.Dec) GenesisState {
	return GenesisState{
		ServiceParams: serviceparams.Params{
			MaxRequestTimeout:  maxRequestTimeout,
			MinDepositMultiple: minDepositMultiple,
			ServiceFeeTax:      serviceFeeTax,
			SlashFraction:      slashFraction,
		},
	}
}

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	params.InitGenesisParameter(&serviceparams.ServiceParameter, ctx, data.ServiceParams)
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return GenesisState{
		ServiceParams: serviceparams.GetSericeParams(ctx),
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		ServiceParams: serviceparams.NewSericeParams(),
	}
}

// get raw genesis raw message for testing
func DefaultGenesisStateForTest() GenesisState {
	serviceParams := serviceparams.NewSericeParams()
	serviceParams.MaxRequestTimeout = 10
	serviceParams.MinDepositMultiple = 10
	return GenesisState{
		ServiceParams: serviceParams,
	}
}

// refund deposit from all bindings
// refund service fee from all request
// no process for service fee tax account
func PrepForZeroHeightGenesis(ctx sdk.Context, k Keeper) {
	store := ctx.KVStore(k.storeKey)

	bindingIterator := sdk.KVStorePrefixIterator(store, bindingPropertyKey)
	defer bindingIterator.Close()

	// refund deposit from all bindings
	for ; bindingIterator.Valid(); bindingIterator.Next() {
		var binding SvcBinding
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bindingIterator.Value(), &binding)
		k.ck.SendCoins(ctx, DepositedCoinsAccAddr, binding.Provider, binding.Deposit)
	}

	// refund service fee from all active request
	requestIterator := sdk.KVStorePrefixIterator(store, activeRequestKey)
	for ; requestIterator.Valid(); requestIterator.Next() {
		var request SvcRequest
		k.cdc.MustUnmarshalBinaryLengthPrefixed(requestIterator.Value(), &request)
		k.ck.SendCoins(ctx, RequestCoinsAccAddr, request.Consumer, request.ServiceFee)
	}

	// refund all incoming fee
	incomingFeeIterator := sdk.KVStorePrefixIterator(store, incomingFeeKey)
	for ; incomingFeeIterator.Valid(); incomingFeeIterator.Next() {
		var incomingFee IncomingFee
		k.cdc.MustUnmarshalBinaryLengthPrefixed(incomingFeeIterator.Value(), &incomingFee)
		k.ck.SendCoins(ctx, RequestCoinsAccAddr, incomingFee.Address, incomingFee.Coins)
	}

	// refund all return fee
	returnedFeeIterator := sdk.KVStorePrefixIterator(store, returnedFeeKey)
	for ; returnedFeeIterator.Valid(); returnedFeeIterator.Next() {
		var returnedFee ReturnedFee
		k.cdc.MustUnmarshalBinaryLengthPrefixed(returnedFeeIterator.Value(), &returnedFee)
		k.ck.SendCoins(ctx, RequestCoinsAccAddr, returnedFee.Address, returnedFee.Coins)
	}
}
