package keeper

import (
	"github.com/irisnet/irishub/app/v3/rand/internal/types"
	"github.com/irisnet/irishub/app/v3/service/exported"
	sdk "github.com/irisnet/irishub/types"
)

// RequestService ...
func (k Keeper) RequestService(ctx sdk.Context, reqID []byte, consumer sdk.AccAddress) sdk.Error {

	// TODO: rand provider
	provider := []sdk.AccAddress{}

	timeout := k.sk.GetParamSet(ctx).MaxRequestTimeout

	requestContextID, err := k.sk.CreateRequestContext(
		ctx,               //
		types.ServiceName, //
		provider,          //
		consumer,          //
		"{}",              // TODO  input 				string
		sdk.Coins{},       // TODO  serviceFeeCap 		sdk.Coins
		timeout,           //
		false,             // TODO  superMode 			bool
		false,             // TODO  repeated			bool
		0,                 // TODO  repeatedFrequency 	uint64
		0,                 // TODO  repeatedTotal 		int64
		exported.RUNNING,  // TODO  state 				exported.RequestContextState
		0,                 // TODO  respThreshold 		uint16
		types.ModuleName,  // TODO  respHandler 		string
	)
	if err != nil {
		return err
	}

	return k.sk.StartRequestContext(ctx, requestContextID)
}

// HandlerResponse ...
func (k Keeper) HandlerResponse(ctx sdk.Context, requestContextID []byte, responseOutput []string) {

	// TODO: Generate random

	// TODO: DequeueOracleTimeoutRandRequest
}

// GetRequestContext ...
func (k Keeper) GetRequestContext(ctx sdk.Context, requestContextID []byte) (exported.RequestContext, bool) {
	return k.sk.GetRequestContext(ctx, requestContextID)
}

// GetMaxServiceRequestTimeout ...
func (k Keeper) GetMaxServiceRequestTimeout(ctx sdk.Context) int64 {
	return k.sk.GetParamSet(ctx).MaxRequestTimeout
}
