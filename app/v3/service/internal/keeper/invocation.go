package keeper

import (
	"fmt"

	cmn "github.com/tendermint/tendermint/libs/common"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v3/service/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

// RegisterResponseCallback registers a module callback for response handling
func (k Keeper) RegisterResponseCallback(moduleName string, respCallback types.ResponseCallback) sdk.Error {
	if _, ok := k.respCallbacks[moduleName]; ok {
		return types.ErrModuleNameRegistered(k.codespace, moduleName)
	}

	k.respCallbacks[moduleName] = respCallback

	return nil
}

// CreateRequestContext creates a request context with the specified params
func (k Keeper) CreateRequestContext(
	ctx sdk.Context,
	serviceName string,
	providers []sdk.AccAddress,
	consumer sdk.AccAddress,
	input string,
	serviceFeeCap sdk.Coins,
	timeout int64,
	superMode bool,
	repeated bool,
	repeatedFrequency uint64,
	repeatedTotal int64,
	state types.RequestContextState,
	responseThreshold uint16,
	moduleName string,
) (cmn.HexBytes, sdk.Error) {
	if superMode {
		_, found := k.gk.GetProfiler(ctx, consumer)
		if !found {
			return nil, types.ErrInvalidProfiler(k.codespace, consumer)
		}
	}

	if len(moduleName) != 0 {
		if _, err := k.GetResponseCallback(moduleName); err != nil {
			return nil, err
		}

		if err := types.ValidateRequest(
			serviceName, serviceFeeCap, providers, input,
			timeout, repeated, repeatedFrequency, repeatedTotal,
		); err != nil {
			return nil, err
		}

		if responseThreshold < 1 || int(responseThreshold) > len(providers) {
			return nil, types.ErrInvalidThreshold(k.codespace, fmt.Sprintf("response threshold must be between [1,%d]", len(providers)))
		}
	}

	svcDef, found := k.GetServiceDefinition(ctx, serviceName)
	if !found {
		return nil, types.ErrUnknownServiceDefinition(k.codespace, serviceName)
	}

	if err := types.ValidateRequestInput(svcDef.Schemas, input); err != nil {
		return nil, err
	}

	params := k.GetParamSet(ctx)
	if timeout > params.MaxRequestTimeout {
		return nil, types.ErrInvalidTimeout(k.codespace, fmt.Sprintf("timeout [%d] must not be greater than the max request timeout [%d]", timeout, params.MaxRequestTimeout))
	}

	if repeated {
		if repeatedFrequency == 0 {
			repeatedFrequency = uint64(timeout)
		}
	} else {
		repeatedFrequency = 0
		repeatedTotal = 0
	}

	batchCounter := uint64(0)
	batchRequestCount := uint16(0)
	batchResponseCount := uint16(0)
	batchState := types.BATCHCOMPLETED

	requestContext := types.NewRequestContext(
		serviceName, providers, consumer, input, serviceFeeCap, timeout,
		superMode, repeated, repeatedFrequency, repeatedTotal, batchCounter,
		batchRequestCount, batchResponseCount, batchState, state,
		responseThreshold, moduleName,
	)

	requestContextID := types.GenerateRequestContextID(ctx.BlockHeight(), k.GetIntraTxCounter(ctx))
	k.SetRequestContext(ctx, requestContextID, requestContext)

	if requestContext.State == types.RUNNING {
		k.AddNewRequestBatch(ctx, requestContextID, ctx.BlockHeight())
	}

	return requestContextID, nil
}

// UpdateRequestContext updates the specified request context
func (k Keeper) UpdateRequestContext(
	ctx sdk.Context,
	requestContextID cmn.HexBytes,
	providers []sdk.AccAddress,
	serviceFeeCap sdk.Coins,
	timeout int64,
	repeatedFreq uint64,
	repeatedTotal int64,
	consumer sdk.AccAddress,
) sdk.Error {
	requestContext, found := k.GetRequestContext(ctx, requestContextID)
	if !found {
		return types.ErrUnknownRequestContext(k.codespace, requestContextID)
	}

	// check authority when called by module
	if len(requestContext.ModuleName) > 0 {
		if err := k.CheckAuthority(ctx, consumer, requestContextID, false); err != nil {
			return err
		}
	}

	if !requestContext.Repeated {
		return types.ErrRequestContextNonRepeated(k.codespace)
	}

	if requestContext.State == types.COMPLETED {
		return types.ErrRequestContextCompleted(k.codespace)
	}

	if len(requestContext.ModuleName) != 0 {
		if err := types.ValidateRequestContextUpdating(providers, serviceFeeCap, timeout, repeatedFreq, repeatedTotal); err != nil {
			return err
		}
	}

	if len(providers) > 0 && requestContext.ResponseThreshold > 0 && len(providers) < int(requestContext.ResponseThreshold) {
		return types.ErrInvalidProviders(k.codespace, fmt.Sprintf("length [%d] of providers must not be less than the response threshold [%d]", len(providers), requestContext.ResponseThreshold))
	}

	params := k.GetParamSet(ctx)
	if timeout > params.MaxRequestTimeout {
		return types.ErrInvalidTimeout(k.codespace, fmt.Sprintf("timeout [%d] must not be greater than the max request timeout [%d]", timeout, params.MaxRequestTimeout))
	}

	if timeout == 0 {
		timeout = requestContext.Timeout
	}

	if repeatedFreq == 0 {
		repeatedFreq = requestContext.RepeatedFrequency
	}

	if repeatedFreq < uint64(timeout) {
		return types.ErrInvalidRepeatedFreq(k.codespace, fmt.Sprintf("repeated frequency [%d] must not be less than the timeout [%d]", repeatedFreq, requestContext.Timeout))
	}

	if repeatedTotal >= 1 && repeatedTotal < int64(requestContext.BatchCounter) {
		return types.ErrInvalidRepeatedTotal(k.codespace, fmt.Sprintf("updated repeated total [%d] must not be less than the current batch counter [%d]", repeatedTotal, requestContext.BatchCounter))
	}

	if len(providers) > 0 {
		requestContext.Providers = providers
	}

	if !serviceFeeCap.Empty() {
		requestContext.ServiceFeeCap = serviceFeeCap
	}

	if timeout > 0 {
		requestContext.Timeout = timeout
	}

	if repeatedFreq > 0 {
		requestContext.RepeatedFrequency = repeatedFreq
	}

	if repeatedTotal != 0 {
		requestContext.RepeatedTotal = repeatedTotal
	}

	k.SetRequestContext(ctx, requestContextID, requestContext)

	return nil
}

// PauseRequestContext suspends the specified request context
func (k Keeper) PauseRequestContext(
	ctx sdk.Context,
	requestContextID cmn.HexBytes,
	consumer sdk.AccAddress,
) sdk.Error {
	requestContext, found := k.GetRequestContext(ctx, requestContextID)
	if !found {
		return types.ErrUnknownRequestContext(k.codespace, requestContextID)
	}

	// check authority when called by module
	if len(requestContext.ModuleName) > 0 {
		if err := k.CheckAuthority(ctx, consumer, requestContextID, false); err != nil {
			return err
		}
	}

	if !requestContext.Repeated {
		return types.ErrRequestContextNonRepeated(k.codespace)
	}

	if requestContext.State != types.RUNNING {
		return types.ErrRequestContextNotStarted(k.codespace)
	}

	requestContext.State = types.PAUSED
	k.SetRequestContext(ctx, requestContextID, requestContext)

	return nil
}

// StartRequestContext starts the specified request context
func (k Keeper) StartRequestContext(
	ctx sdk.Context,
	requestContextID cmn.HexBytes,
	consumer sdk.AccAddress,
) sdk.Error {
	requestContext, found := k.GetRequestContext(ctx, requestContextID)
	if !found {
		return types.ErrUnknownRequestContext(k.codespace, requestContextID)
	}

	// check authority when called by module
	if len(requestContext.ModuleName) > 0 {
		if err := k.CheckAuthority(ctx, consumer, requestContextID, false); err != nil {
			return err
		}
	}

	if !requestContext.Repeated {
		return types.ErrRequestContextNonRepeated(k.codespace)
	}

	if requestContext.State != types.PAUSED {
		return types.ErrRequestContextNotPaused(k.codespace)
	}

	requestContext.State = types.RUNNING
	k.SetRequestContext(ctx, requestContextID, requestContext)

	if requestContext.BatchState == types.BATCHCOMPLETED &&
		(requestContext.BatchRequestCount == 0 ||
			requestContext.BatchRequestCount != requestContext.BatchResponseCount) {
		k.AddNewRequestBatch(ctx, requestContextID, ctx.BlockHeight())
	}

	return nil
}

// KillRequestContext terminates the specified request context
func (k Keeper) KillRequestContext(
	ctx sdk.Context,
	requestContextID cmn.HexBytes,
	consumer sdk.AccAddress,
) sdk.Error {
	requestContext, found := k.GetRequestContext(ctx, requestContextID)
	if !found {
		return types.ErrUnknownRequestContext(k.codespace, requestContextID)
	}

	// check authority when called by module
	if len(requestContext.ModuleName) > 0 {
		if err := k.CheckAuthority(ctx, consumer, requestContextID, false); err != nil {
			return err
		}
	}

	if !requestContext.Repeated {
		return types.ErrRequestContextNonRepeated(k.codespace)
	}

	requestContext.State = types.COMPLETED
	k.SetRequestContext(ctx, requestContextID, requestContext)

	return nil
}

// SetRequestContext sets the specified request context
func (k Keeper) SetRequestContext(ctx sdk.Context, requestContextID cmn.HexBytes, requestContext types.RequestContext) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(requestContext)
	store.Set(GetRequestContextKey(requestContextID), bz)
}

// GetRequestContext retrieves the specified request context
func (k Keeper) GetRequestContext(ctx sdk.Context, requestContextID cmn.HexBytes) (requestContext types.RequestContext, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(GetRequestContextKey(requestContextID))
	if bz == nil {
		return requestContext, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &requestContext)
	return requestContext, true
}

// IterateRequestContexts iterates through all request contexts
func (k Keeper) IterateRequestContexts(
	ctx sdk.Context,
	op func(requestContextID cmn.HexBytes, requestContext types.RequestContext) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, requestContextKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		requestContextID := iterator.Key()[1:]

		var requestContext types.RequestContext
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &requestContext)

		if stop := op(requestContextID, requestContext); stop {
			break
		}
	}
}

// InitiateRequests creates requests for the given providers from the specified request context
// Note: make sure that request context is valid and running, and providers are valid
func (k Keeper) InitiateRequests(
	ctx sdk.Context,
	requestContextID cmn.HexBytes,
	providers []sdk.AccAddress,
) (tags sdk.Tags) {
	requestContext, _ := k.GetRequestContext(ctx, requestContextID)
	requestContext.BatchCounter++

	tags = sdk.NewTags()

	for providerIndex, provider := range providers {
		request := k.buildRequest(
			ctx, requestContextID, requestContext.BatchCounter,
			requestContext.ServiceName, provider, requestContext.SuperMode,
		)

		requestID := types.GenerateRequestID(requestContextID, requestContext.BatchCounter, int16(providerIndex))
		k.SetCompactRequest(ctx, requestID, request)

		k.AddActiveRequest(ctx, requestContext.ServiceName, provider, ctx.BlockHeight()+requestContext.Timeout, requestID)

		tags = tags.AppendTags(sdk.NewTags(
			types.TagRequestID, []byte(requestID.String()),
			types.TagProvider, []byte(provider.String()),
			types.TagConsumer, []byte(requestContext.Consumer.String()),
			types.TagServiceName, []byte(requestContext.ServiceName),
			types.TagServiceFee, []byte(request.ServiceFee.String()),
			types.TagRequestHeight, []byte(fmt.Sprintf("%d", request.RequestHeight)),
			types.TagExpirationHeight, []byte(fmt.Sprintf("%d", request.RequestHeight+requestContext.Timeout)),
		))
	}

	requestContext.BatchState = types.BATCHRUNNING
	requestContext.BatchResponseCount = 0
	requestContext.BatchRequestCount = uint16(len(providers))

	k.SetRequestContext(ctx, requestContextID, requestContext)

	return tags
}

// SkipCurrentRequestBatch skips the current request batch
func (k Keeper) SkipCurrentRequestBatch(ctx sdk.Context, requestContextID cmn.HexBytes, requestContext types.RequestContext) {
	requestContext.BatchCounter++
	requestContext.BatchState = types.BATCHRUNNING
	requestContext.BatchRequestCount = 0
	requestContext.BatchResponseCount = 0

	k.SetRequestContext(ctx, requestContextID, requestContext)
	k.AddRequestBatchExpiration(ctx, requestContextID, ctx.BlockHeight()+requestContext.Timeout)
}

// buildRequest builds a request for the given provider from the specified request context
// Note: make that the binding exists
func (k Keeper) buildRequest(
	ctx sdk.Context,
	requestContextID cmn.HexBytes,
	batchCounter uint64,
	serviceName string,
	provider sdk.AccAddress,
	superMode bool,
) types.CompactRequest {
	var serviceFee sdk.Coins

	if !superMode {
		binding, _ := k.GetServiceBinding(ctx, serviceName, provider)
		serviceFee = k.GetBasePrice(ctx, binding)
	}

	request := types.NewCompactRequest(
		requestContextID,
		batchCounter,
		provider,
		serviceFee,
		ctx.BlockHeight(),
	)

	return request
}

// SetCompactRequest sets the specified compact request
func (k Keeper) SetCompactRequest(ctx sdk.Context, requestID cmn.HexBytes, request types.CompactRequest) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(request)
	store.Set(GetRequestKey(requestID), bz)
}

// GetCompactRequest retrieves the specified compact request
func (k Keeper) GetCompactRequest(ctx sdk.Context, requestID cmn.HexBytes) (request types.CompactRequest, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(GetRequestKey(requestID))
	if bz == nil {
		return request, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &request)
	return request, true
}

// GetRequest returns the specified request
func (k Keeper) GetRequest(ctx sdk.Context, requestID cmn.HexBytes) (request types.Request, found bool) {
	compactRequest, found := k.GetCompactRequest(ctx, requestID)
	if !found {
		return request, false
	}

	requestContext, found := k.GetRequestContext(ctx, compactRequest.RequestContextID)
	if !found {
		return request, false
	}

	request = types.NewRequest(
		requestID,
		requestContext.ServiceName,
		compactRequest.Provider,
		requestContext.Consumer,
		requestContext.Input,
		compactRequest.ServiceFee,
		requestContext.SuperMode,
		compactRequest.RequestHeight,
		compactRequest.RequestHeight+requestContext.Timeout,
		compactRequest.RequestContextID,
		compactRequest.RequestContextBatchCounter,
	)

	return request, true
}

// IterateRequests iterates through all compact requests
func (k Keeper) IterateRequests(
	ctx sdk.Context,
	op func(requestID cmn.HexBytes, request types.CompactRequest) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, requestKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		requestID := iterator.Key()[1:]

		var request types.CompactRequest
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &request)

		if stop := op(requestID, request); stop {
			break
		}
	}
}

// RequestsIteratorByReqCtx returns an iterator for all requests of the specified request context ID and batch counter
func (k Keeper) RequestsIteratorByReqCtx(ctx sdk.Context, requestContextID cmn.HexBytes, batchCounter uint64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, GetRequestSubspaceByReqCtx(requestContextID, batchCounter))
}

// AddActiveRequest adds the specified active request
func (k Keeper) AddActiveRequest(
	ctx sdk.Context,
	serviceName string,
	provider sdk.AccAddress,
	expirationHeight int64,
	requestID cmn.HexBytes,
) {
	k.AddActiveRequestByBinding(ctx, serviceName, provider, expirationHeight, requestID)
	k.AddActiveRequestByID(ctx, requestID)

	k.GetMetrics().ActiveRequests.Add(1)
}

// DeleteActiveRequest deletes the specified active request
func (k Keeper) DeleteActiveRequest(
	ctx sdk.Context,
	serviceName string,
	provider sdk.AccAddress,
	expirationHeight int64,
	requestID cmn.HexBytes,
) {
	k.DeleteActiveRequestByBinding(ctx, serviceName, provider, expirationHeight, requestID)
	k.DeleteActiveRequestByID(ctx, requestID)

	k.GetMetrics().ActiveRequests.Add(-1)
}

// AddActiveRequestByBinding adds the specified active request by the binding
func (k Keeper) AddActiveRequestByBinding(
	ctx sdk.Context,
	serviceName string,
	provider sdk.AccAddress,
	expirationHeight int64,
	requestID cmn.HexBytes,
) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(requestID)
	store.Set(GetActiveRequestKey(serviceName, provider, expirationHeight, requestID), bz)
}

// DeleteActiveRequestByBinding deletes the specified active request by the binding
func (k Keeper) DeleteActiveRequestByBinding(
	ctx sdk.Context,
	serviceName string,
	provider sdk.AccAddress,
	expirationHeight int64,
	requestID cmn.HexBytes,
) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetActiveRequestKey(serviceName, provider, expirationHeight, requestID))
}

// AddActiveRequestByID adds the specified active request by request ID
func (k Keeper) AddActiveRequestByID(
	ctx sdk.Context,
	requestID cmn.HexBytes,
) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(requestID)
	store.Set(GetActiveRequestKeyByID(requestID), bz)
}

// DeleteActiveRequestByID deletes the specified active request by request ID
func (k Keeper) DeleteActiveRequestByID(
	ctx sdk.Context,
	requestID cmn.HexBytes,
) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetActiveRequestKeyByID(requestID))
}

// IsRequestActive checks if the specified request is active
func (k Keeper) IsRequestActive(
	ctx sdk.Context,
	requestID cmn.HexBytes,
) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(GetActiveRequestKeyByID(requestID))
}

// AddRequestBatchExpiration adds a request batch to the expiration queue
func (k Keeper) AddRequestBatchExpiration(ctx sdk.Context, requestContextID cmn.HexBytes, expirationHeight int64) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(requestContextID)
	store.Set(GetExpiredRequestBatchKey(requestContextID, expirationHeight), bz)
}

// DeleteRequestBatchExpiration deletes the request batch from the expiration queue
func (k Keeper) DeleteRequestBatchExpiration(ctx sdk.Context, requestContextID cmn.HexBytes, expirationHeight int64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetExpiredRequestBatchKey(requestContextID, expirationHeight))
}

// AddNewRequestBatch adds a request batch to the new request batch queue
func (k Keeper) AddNewRequestBatch(ctx sdk.Context, requestContextID cmn.HexBytes, requestBatchHeight int64) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(requestContextID)
	store.Set(GetNewRequestBatchKey(requestContextID, requestBatchHeight), bz)
}

// DeleteNewRequestBatch deletes the request batch in the given height from the new request batch queue
func (k Keeper) DeleteNewRequestBatch(ctx sdk.Context, requestContextID cmn.HexBytes, requestBatchHeight int64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetNewRequestBatchKey(requestContextID, requestBatchHeight))
}

// HasNewRequestBatch checks if the new request batch from the specified request context exists in the given height
func (k Keeper) HasNewRequestBatch(ctx sdk.Context, requestContextID cmn.HexBytes, requestBatchHeight int64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(GetNewRequestBatchKey(requestContextID, requestBatchHeight))
}

// IterateExpiredRequestBatch iterates through the expired request batch queue in the specified height
func (k Keeper) IterateExpiredRequestBatch(
	ctx sdk.Context,
	expirationHeight int64,
	op func(requestContextID cmn.HexBytes, requestContext types.RequestContext),
) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, GetExpiredRequestBatchSubspace(expirationHeight))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var requestContextID cmn.HexBytes
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &requestContextID)

		requestContext, _ := k.GetRequestContext(ctx, requestContextID)

		op(requestContextID, requestContext)
	}
}

// IterateNewRequestBatch iterates through the new request batch queue in the specified height
func (k Keeper) IterateNewRequestBatch(
	ctx sdk.Context,
	requestBatchHeight int64,
	op func(requestContextID cmn.HexBytes, requestContext types.RequestContext),
) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, GetNewRequestBatchSubspace(requestBatchHeight))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var requestContextID cmn.HexBytes
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &requestContextID)

		requestContext, _ := k.GetRequestContext(ctx, requestContextID)

		op(requestContextID, requestContext)
	}
}

// ActiveRequestsIterator returns an iterator for all the active requests of the specified service binding
func (k Keeper) ActiveRequestsIterator(ctx sdk.Context, serviceName string, provider sdk.AccAddress) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, GetActiveRequestSubspace(serviceName, provider))
}

// ActiveRequestsIteratorByReqCtx returns an iterator for all the active requests of the specified service binding
func (k Keeper) ActiveRequestsIteratorByReqCtx(ctx sdk.Context, requestContextID cmn.HexBytes, batchCounter uint64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, GetActiveRequestSubspaceByReqCtx(requestContextID, batchCounter))
}

// AllActiveRequestsIterator returns an iterator for all the active requests
func (k Keeper) AllActiveRequestsIterator(store sdk.KVStore) sdk.Iterator {
	return sdk.KVStorePrefixIterator(store, activeRequestKey)
}

// IterateActiveRequests iterates through the active requests for the specified request context ID and batch counter
func (k Keeper) IterateActiveRequests(
	ctx sdk.Context,
	requestContextID cmn.HexBytes,
	batchCounter uint64,
	op func(requestID cmn.HexBytes, request types.Request),
) {
	iterator := k.ActiveRequestsIteratorByReqCtx(ctx, requestContextID, batchCounter)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var requestID cmn.HexBytes
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &requestID)

		request, _ := k.GetRequest(ctx, requestID)

		op(requestID, request)
	}
}

// FilterServiceProviders gets the providers which satisfy the specified service fee requirement
func (k Keeper) FilterServiceProviders(
	ctx sdk.Context,
	serviceName string,
	providers []sdk.AccAddress,
	serviceFeeCap sdk.Coins,
) ([]sdk.AccAddress, sdk.Coins) {
	var newProviders []sdk.AccAddress
	var totalPrices sdk.Coins

	for _, provider := range providers {
		binding, found := k.GetServiceBinding(ctx, serviceName, provider)

		if found && binding.Available {
			price := k.GetBasePrice(ctx, binding)

			if price.IsAllLTE(serviceFeeCap) {
				newProviders = append(newProviders, provider)
				totalPrices = totalPrices.Add(price)
			}
		}
	}

	return newProviders, totalPrices
}

// DeductServiceFees deducts the given service fees from the specified consumer
func (k Keeper) DeductServiceFees(ctx sdk.Context, consumer sdk.AccAddress, serviceFees sdk.Coins) sdk.Error {
	_, err := k.bk.SendCoins(ctx, consumer, auth.ServiceRequestCoinsAccAddr, serviceFees)

	if !serviceFees.IsZero() {
		ctx.CoinFlowTags().AppendCoinFlowTag(ctx, consumer.String(),
			auth.ServiceDepositCoinsAccAddr.String(), serviceFees.String(), sdk.ServiceFeeDeductFlow, "")
	}

	if err != nil {
		return err
	}

	return nil
}

// AddResponse adds the response for the specified request ID
func (k Keeper) AddResponse(
	ctx sdk.Context,
	requestID string,
	provider sdk.AccAddress,
	result,
	output string,
) (request types.Request, response types.Response, tags sdk.Tags, err sdk.Error) {
	reqID, _ := types.ConvertRequestID(requestID)

	request, found := k.GetRequest(ctx, reqID)
	if !found {
		return request, response, tags, types.ErrUnknownRequest(k.codespace, reqID)
	}

	if !provider.Equals(request.Provider) {
		return request, response, tags, types.ErrInvalidResponse(k.codespace, "provider does not match")
	}

	if !k.IsRequestActive(ctx, reqID) {
		return request, response, tags, types.ErrInvalidResponse(k.codespace, "request is not active")
	}

	svcDef, _ := k.GetServiceDefinition(ctx, request.ServiceName)

	if len(output) > 0 && types.ValidateResponseOutput(svcDef.Schemas, output) != nil {
		tags, err = k.Slash(ctx, reqID)
		if err != nil {
			panic(err)
		}

		if err := k.RefundServiceFee(ctx, request.Consumer, request.ServiceFee); err != nil {
			panic(err)
		}
	} else {
		if err := k.AddEarnedFee(ctx, provider, request.ServiceFee); err != nil {
			return request, response, tags, err
		}
	}

	requestContextID := request.RequestContextID

	response = types.NewResponse(provider, request.Consumer, result, output, requestContextID, request.RequestContextBatchCounter)
	k.SetResponse(ctx, reqID, response)

	k.DeleteActiveRequest(ctx, request.ServiceName, provider, request.ExpirationHeight, reqID)

	requestContext, _ := k.GetRequestContext(ctx, requestContextID)
	requestContext.BatchResponseCount++

	if requestContext.BatchResponseCount == requestContext.BatchRequestCount {
		requestContext.BatchState = types.BATCHCOMPLETED

		if len(requestContext.ModuleName) != 0 {
			tags = tags.AppendTags(k.Callback(ctx, requestContextID))
		}
	}

	k.SetRequestContext(ctx, requestContextID, requestContext)

	return request, response, tags, nil
}

// Callback callbacks the corresponding response callback handler
func (k Keeper) Callback(ctx sdk.Context, requestContextID cmn.HexBytes) sdk.Tags {
	requestContext, _ := k.GetRequestContext(ctx, requestContextID)

	respCallback, _ := k.GetResponseCallback(requestContext.ModuleName)
	outputs := k.GetResponseOutputs(ctx, requestContextID, requestContext.BatchCounter)

	if len(outputs) >= int(requestContext.ResponseThreshold) {
		return respCallback(ctx, requestContextID, outputs, nil)
	} else {
		return respCallback(
			ctx,
			requestContextID,
			outputs,
			fmt.Errorf(
				"batch %d at least %d valid outputs required, but %d received",
				requestContext.BatchCounter, requestContext.ResponseThreshold, len(outputs),
			),
		)
	}
}

// SetResponse sets the specified response
func (k Keeper) SetResponse(ctx sdk.Context, requestID cmn.HexBytes, response types.Response) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(response)
	store.Set(GetResponseKey(requestID), bz)
}

// GetResponse returns a response with the speicified request ID
func (k Keeper) GetResponse(ctx sdk.Context, requestID cmn.HexBytes) (response types.Response, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(GetResponseKey(requestID))
	if bz == nil {
		return response, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &response)
	return response, true
}

// IterateResponses iterates through all responses
func (k Keeper) IterateResponses(
	ctx sdk.Context,
	op func(requestID cmn.HexBytes, response types.Response) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, responseKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		requestID := iterator.Key()[1:]

		var response types.Response
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &response)

		if stop := op(requestID, response); stop {
			break
		}
	}
}

// ResponsesIteratorByReqCtx returns an iterator for all responses of the specified request context and batch counter
func (k Keeper) ResponsesIteratorByReqCtx(ctx sdk.Context, requestContextID cmn.HexBytes, batchCounter uint64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, GetResponseSubspaceByReqCtx(requestContextID, batchCounter))
}

// GetResponseOutputs retrieves all response outputs of the specified request context and batch counter
func (k Keeper) GetResponseOutputs(ctx sdk.Context, requestContextID cmn.HexBytes, batchCounter uint64) []string {
	iterator := k.ResponsesIteratorByReqCtx(ctx, requestContextID, batchCounter)
	defer iterator.Close()

	var outputs []string
	for ; iterator.Valid(); iterator.Next() {
		var response types.Response
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &response)

		if len(response.Output) > 0 {
			outputs = append(outputs, response.Output)
		}
	}

	return outputs
}

// Slash slashes the provider from the specified request
// Ensure that the request is valid
func (k Keeper) Slash(ctx sdk.Context, requestID cmn.HexBytes) (tags sdk.Tags, err sdk.Error) {
	request, _ := k.GetRequest(ctx, requestID)
	binding, _ := k.GetServiceBinding(ctx, request.ServiceName, request.Provider)

	slashFraction := k.GetParamSet(ctx).SlashFraction

	depositAmt := binding.Deposit.AmountOf(sdk.IrisAtto)
	slashedAmt := sdk.NewDecFromInt(depositAmt).Mul(slashFraction).TruncateInt()
	slashedCoins := sdk.NewCoins(sdk.NewCoin(sdk.IrisAtto, slashedAmt))

	deposit, hasNeg := binding.Deposit.SafeSub(slashedCoins)
	if hasNeg {
		return tags, sdk.ErrInsufficientCoins(fmt.Sprintf("%s is less than %s", binding.Deposit, slashedCoins))
	}

	_, err = k.bk.BurnCoins(ctx, auth.ServiceDepositCoinsAccAddr, slashedCoins)
	if err != nil {
		return tags, err
	}

	binding.Deposit = deposit
	if binding.Available {
		minDeposit := k.getMinDeposit(ctx, binding.Pricing)

		if !binding.Deposit.IsAllGTE(minDeposit) {
			binding.Available = false
			binding.DisabledTime = ctx.BlockHeader().Time
		}
	}

	k.SetServiceBinding(ctx, binding)

	slashTags := sdk.NewTags(
		types.TagRequestID, []byte(requestID.String()),
		types.TagProvider, []byte(request.Provider.String()),
		types.TagConsumer, []byte(request.Consumer.String()),
		types.TagSlashedCoins, []byte(slashedCoins.String()),
	)

	return slashTags, nil
}

// RefundServiceFee refunds the service fee to the specified consumer
func (k Keeper) RefundServiceFee(ctx sdk.Context, consumer sdk.AccAddress, serviceFee sdk.Coins) sdk.Error {
	_, err := k.bk.SendCoins(ctx, auth.ServiceRequestCoinsAccAddr, consumer, serviceFee)

	if !serviceFee.IsZero() {
		ctx.CoinFlowTags().AppendCoinFlowTag(ctx, auth.ServiceRequestCoinsAccAddr.String(),
			consumer.String(), serviceFee.String(), sdk.ServiceFeeRefundFlow, "")
	}

	if err != nil {
		return err
	}

	return nil
}

// AddEarnedFee adds the earned fee for the given provider
func (k Keeper) AddEarnedFee(ctx sdk.Context, provider sdk.AccAddress, fee sdk.Coins) sdk.Error {
	params := k.GetParamSet(ctx)
	taxRate := params.ServiceFeeTax

	taxCoins := sdk.Coins{}
	for _, coin := range fee {
		taxAmount := sdk.NewDecFromInt(coin.Amount).Mul(taxRate).TruncateInt()
		taxCoins = taxCoins.Add(sdk.NewCoins(sdk.NewCoin(coin.Denom, taxAmount)))
	}

	_, err := k.bk.SendCoins(ctx, auth.ServiceRequestCoinsAccAddr, auth.ServiceTaxCoinsAccAddr, taxCoins)
	if err != nil {
		return err
	}

	earnedFee, hasNeg := fee.SafeSub(taxCoins)
	if hasNeg {
		errMsg := fmt.Sprintf("%s is less than %s", fee, taxCoins)
		return sdk.ErrInsufficientFunds(errMsg)
	}

	fees, _ := k.GetEarnedFees(ctx, provider)
	k.SetEarnedFees(ctx, provider, fees.Coins.Add(earnedFee))

	return nil
}

func (k Keeper) SetEarnedFees(ctx sdk.Context, provider sdk.AccAddress, fees sdk.Coins) {
	store := ctx.KVStore(k.storeKey)

	earnedFees := types.NewEarnedFees(provider, fees)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(earnedFees)

	store.Set(GetEarnedFeesKey(provider), bz)
}

func (k Keeper) DeleteEarnedFees(ctx sdk.Context, provider sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetEarnedFeesKey(provider))
}

func (k Keeper) GetEarnedFees(ctx sdk.Context, provider sdk.AccAddress) (fees types.EarnedFees, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(GetEarnedFeesKey(provider))
	if bz == nil {
		return fees, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &fees)
	return fees, true
}

// WithdrawEarnedFees withdraws the earned fees of the specified provider
func (k Keeper) WithdrawEarnedFees(ctx sdk.Context, provider sdk.AccAddress) sdk.Error {
	fees, found := k.GetEarnedFees(ctx, provider)
	if !found {
		return types.ErrNoEarnedFees(k.codespace, provider)
	}

	withdrawAddr := k.GetWithdrawAddress(ctx, provider)

	_, err := k.bk.SendCoins(ctx, auth.ServiceRequestCoinsAccAddr, withdrawAddr, fees.Coins)
	if err != nil {
		return err
	}

	k.DeleteEarnedFees(ctx, provider)

	return nil
}

// WithdrawTax withdraws the service tax to the speicified destination address by the trustee
func (k Keeper) WithdrawTax(ctx sdk.Context, trustee sdk.AccAddress, destAddress sdk.AccAddress, amt sdk.Coins) sdk.Error {
	if _, found := k.gk.GetTrustee(ctx, trustee); !found {
		return types.ErrInvalidTrustee(k.codespace, trustee)
	}

	_, err := k.bk.SendCoins(ctx, auth.ServiceTaxCoinsAccAddr, destAddress, amt)
	if err != nil {
		return err
	}

	return nil
}

// AllEarnedFeesIterator returns an iterator for all the earned fees
func (k Keeper) AllEarnedFeesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, earnedFeesKey)
}

// RefundEarnedFees refunds all the incoming fees
func (k Keeper) RefundEarnedFees(ctx sdk.Context) sdk.Error {
	iterator := k.AllEarnedFeesIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var earnedFees types.EarnedFees
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &earnedFees)

		_, err := k.bk.SendCoins(ctx, auth.ServiceRequestCoinsAccAddr, earnedFees.Address, earnedFees.Coins)
		if err != nil {
			return err
		}
	}

	return nil
}

// RefundServiceFees refunds the service fees of all the active requests
func (k Keeper) RefundServiceFees(ctx sdk.Context) sdk.Error {
	iterator := k.AllActiveRequestsIterator(ctx.KVStore(k.storeKey))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var requestID cmn.HexBytes
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &requestID)

		request, _ := k.GetRequest(ctx, requestID)

		_, err := k.bk.SendCoins(ctx, auth.ServiceRequestCoinsAccAddr, request.Consumer, request.ServiceFee)
		if err != nil {
			return err
		}
	}

	return nil
}

// CheckAuthority checks if the operation on the specified request context is authorized
func (k Keeper) CheckAuthority(
	ctx sdk.Context,
	consumer sdk.AccAddress,
	requestContextID cmn.HexBytes,
	checkModule bool,
) sdk.Error {
	requestContext, found := k.GetRequestContext(ctx, requestContextID)
	if !found {
		return types.ErrUnknownRequestContext(k.codespace, requestContextID)
	}

	if !consumer.Equals(requestContext.Consumer) {
		return types.ErrNotAuthorized(k.codespace, "consumer not matching")
	}

	if checkModule && len(requestContext.ModuleName) > 0 {
		return types.ErrNotAuthorized(k.codespace, "not authorized operation")
	}

	return nil
}

// GetResponseCallback gets the registered module callback for response handling
func (k Keeper) GetResponseCallback(moduleName string) (types.ResponseCallback, sdk.Error) {
	respCallback, ok := k.respCallbacks[moduleName]
	if !ok {
		return nil, types.ErrModuleNameNotRegistered(k.codespace, moduleName)
	}

	return respCallback, nil
}

// GetIntraTxCounter returns the current tx counter and increases it by 1
func (k Keeper) GetIntraTxCounter(ctx sdk.Context) int16 {
	store := ctx.KVStore(k.storeKey)
	var counter int16

	bz := store.Get(GetIntraTxCounterKey())
	if bz != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &counter)
	}

	k.SetIntraTxCounter(ctx, counter+1)
	return counter
}

// SetIntraTxCounter sets the tx counter
func (k Keeper) SetIntraTxCounter(ctx sdk.Context, counter int16) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(counter)
	store.Set(GetIntraTxCounterKey(), bz)
}
