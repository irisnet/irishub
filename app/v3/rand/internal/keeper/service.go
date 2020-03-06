package keeper

import (
	"github.com/irisnet/irishub/app/v3/service/exported"
	sdk "github.com/irisnet/irishub/types"
)

// RequestService ...
func (k Keeper) RequestService(ctx sdk.Context, reqID []byte) sdk.Error {

	// TODO: CreateRequestContext

	// TODO: StartRequestContext

	return nil
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
