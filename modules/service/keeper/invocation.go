package keeper

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	gogotypes "github.com/gogo/protobuf/types"

	"github.com/tendermint/tendermint/crypto/tmhash"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/service/types"
)

// RegisterResponseCallback registers a module callback for response handling
func (k Keeper) RegisterResponseCallback(moduleName string, respCallback types.ResponseCallback) error {
	if _, ok := k.respCallbacks[moduleName]; ok {
		return sdkerrors.Wrapf(types.ErrCallbackRegistered, "%s already registered for module %s", "response callback", moduleName)
	}

	k.respCallbacks[moduleName] = respCallback

	return nil
}

// RegisterStateCallback registers a module callback for state handling
func (k Keeper) RegisterStateCallback(moduleName string, stateCallback types.StateCallback) error {
	if _, ok := k.stateCallbacks[moduleName]; ok {
		return sdkerrors.Wrapf(types.ErrCallbackRegistered, "%s already registered for module %s", "state callback", moduleName)
	}

	k.stateCallbacks[moduleName] = stateCallback

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
	responseThreshold uint32,
	moduleName string,
) (
	tmbytes.HexBytes, error,
) {
	if len(moduleName) != 0 {
		if _, err := k.GetResponseCallback(moduleName); err != nil {
			return nil, err
		}

		if _, err := k.GetStateCallback(moduleName); err != nil {
			return nil, err
		}

		if err := types.ValidateRequest(
			serviceName, serviceFeeCap, providers, input,
			timeout, repeated, repeatedFrequency, repeatedTotal,
		); err != nil {
			return nil, err
		}

		if responseThreshold < 1 || int(responseThreshold) > len(providers) {
			return nil, sdkerrors.Wrapf(types.ErrInvalidResponseThreshold, "response threshold [%d] must be between [1,%d]", responseThreshold, len(providers))
		}
	}

	_, found := k.GetServiceDefinition(ctx, serviceName)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrUnknownServiceDefinition, serviceName)
	}

	if err := types.ValidateRequestInput(input); err != nil {
		return nil, err
	}

	if err := k.validateServiceFeeCap(ctx, serviceFeeCap); err != nil {
		return nil, err
	}

	maxRequestTimeout := k.MaxRequestTimeout(ctx)
	if timeout > maxRequestTimeout {
		return nil, sdkerrors.Wrapf(types.ErrInvalidTimeout, "timeout [%d] must not be greater than the max request timeout [%d]", timeout, maxRequestTimeout)
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
	batchRequestCount := uint32(0)
	batchResponseCount := uint32(0)
	batchResponseThreshold := responseThreshold
	batchState := types.BATCHCOMPLETED

	requestContext := types.NewRequestContext(
		serviceName, providers, consumer, input, serviceFeeCap, timeout,
		superMode, repeated, repeatedFrequency, repeatedTotal, batchCounter,
		batchRequestCount, batchResponseCount, batchResponseThreshold,
		batchState, state, responseThreshold, moduleName,
	)

	requestContextID := types.GenerateRequestContextID(TxHash(ctx), k.GetInternalIndex(ctx))
	k.SetRequestContext(ctx, requestContextID, requestContext)

	if requestContext.State == types.RUNNING {
		k.AddNewRequestBatch(ctx, requestContextID, ctx.BlockHeight())
	}

	return requestContextID, nil
}

// UpdateRequestContext updates the specified request context
func (k Keeper) UpdateRequestContext(
	ctx sdk.Context,
	requestContextID tmbytes.HexBytes,
	providers []sdk.AccAddress,
	respThreshold uint32,
	serviceFeeCap sdk.Coins,
	timeout int64,
	repeatedFreq uint64,
	repeatedTotal int64,
	consumer sdk.AccAddress,
) error {
	pds := make([]string, len(providers))
	for i, provider := range providers {
		pds[i] = provider.String()
	}
	requestContext, found := k.GetRequestContext(ctx, requestContextID)
	if !found {
		return sdkerrors.Wrap(types.ErrUnknownRequestContext, requestContextID.String())
	}

	// check authority when called by module
	if len(requestContext.ModuleName) > 0 {
		if err := k.CheckAuthority(ctx, consumer, requestContextID, false); err != nil {
			return err
		}
	}

	if requestContext.State == types.COMPLETED {
		return types.ErrRequestContextCompleted
	}

	if len(requestContext.ModuleName) > 0 {
		if err := types.ValidateRequestContextUpdating(providers, serviceFeeCap, timeout, repeatedFreq, repeatedTotal); err != nil {
			return err
		}

		if respThreshold == 0 {
			respThreshold = requestContext.ResponseThreshold
		}

		if len(pds) == 0 {
			pds = requestContext.Providers
		}

		if respThreshold > uint32(len(pds)) {
			return sdkerrors.Wrapf(types.ErrInvalidResponseThreshold, "response threshold [%d] must be between [1,%d]", respThreshold, len(pds))
		}

		if respThreshold > 0 {
			requestContext.ResponseThreshold = respThreshold
		}
	}

	if !serviceFeeCap.Empty() {
		if err := k.validateServiceFeeCap(ctx, serviceFeeCap); err != nil {
			return err
		}

		requestContext.ServiceFeeCap = serviceFeeCap
	}

	maxRequestTimeout := k.MaxRequestTimeout(ctx)
	if timeout > maxRequestTimeout {
		return sdkerrors.Wrapf(types.ErrInvalidTimeout, "timeout [%d] must not be greater than the max request timeout [%d]", timeout, maxRequestTimeout)
	}

	if timeout == 0 {
		timeout = requestContext.Timeout
	}

	if repeatedFreq == 0 {
		repeatedFreq = requestContext.RepeatedFrequency
	}

	if repeatedFreq < uint64(timeout) {
		return sdkerrors.Wrapf(types.ErrInvalidRepeatedFreq, "repeated frequency [%d] must not be less than the timeout [%d]", repeatedFreq, requestContext.Timeout)
	}

	if repeatedTotal >= 1 && repeatedTotal < int64(requestContext.BatchCounter) {
		return sdkerrors.Wrapf(types.ErrInvalidRepeatedTotal, "updated repeated total [%d] must not be less than the current batch counter [%d]", repeatedTotal, requestContext.BatchCounter)
	}

	if len(pds) > 0 {
		requestContext.Providers = pds
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
	requestContextID tmbytes.HexBytes,
	consumer sdk.AccAddress,
) error {
	requestContext, found := k.GetRequestContext(ctx, requestContextID)
	if !found {
		return sdkerrors.Wrap(types.ErrUnknownRequestContext, requestContextID.String())
	}

	// check authority when called by module
	if len(requestContext.ModuleName) > 0 {
		if err := k.CheckAuthority(ctx, consumer, requestContextID, false); err != nil {
			return err
		}
	}

	if !requestContext.Repeated {
		return types.ErrRequestContextNonRepeated
	}

	if requestContext.State != types.RUNNING {
		return types.ErrRequestContextNotRunning
	}

	requestContext.State = types.PAUSED
	k.SetRequestContext(ctx, requestContextID, requestContext)

	return nil
}

// StartRequestContext starts the specified request context
func (k Keeper) StartRequestContext(
	ctx sdk.Context,
	requestContextID tmbytes.HexBytes,
	consumer sdk.AccAddress,
) error {
	requestContext, found := k.GetRequestContext(ctx, requestContextID)
	if !found {
		return sdkerrors.Wrap(types.ErrUnknownRequestContext, requestContextID.String())
	}

	// check authority when called by module
	if len(requestContext.ModuleName) > 0 {
		if err := k.CheckAuthority(ctx, consumer, requestContextID, false); err != nil {
			return err
		}
	}

	if requestContext.State != types.PAUSED {
		return types.ErrRequestContextNotPaused
	}

	requestContext.State = types.RUNNING
	k.SetRequestContext(ctx, requestContextID, requestContext)

	// add to the new request batch queue if existing in neither expired nor new request batch queue
	if !k.HasRequestBatchExpiration(ctx, requestContextID) && !k.HasNewRequestBatch(ctx, requestContextID) {
		k.AddNewRequestBatch(ctx, requestContextID, ctx.BlockHeight())
	}

	return nil
}

// KillRequestContext terminates the specified request context
func (k Keeper) KillRequestContext(
	ctx sdk.Context,
	requestContextID tmbytes.HexBytes,
	consumer sdk.AccAddress,
) error {
	requestContext, found := k.GetRequestContext(ctx, requestContextID)
	if !found {
		return sdkerrors.Wrap(types.ErrUnknownRequestContext, requestContextID.String())
	}

	// check authority when called by module
	if len(requestContext.ModuleName) > 0 {
		if err := k.CheckAuthority(ctx, consumer, requestContextID, false); err != nil {
			return err
		}
	}

	if !requestContext.Repeated {
		return types.ErrRequestContextNonRepeated
	}

	requestContext.State = types.COMPLETED
	k.SetRequestContext(ctx, requestContextID, requestContext)

	return nil
}

// SetRequestContext sets the specified request context
func (k Keeper) SetRequestContext(ctx sdk.Context, requestContextID tmbytes.HexBytes, requestContext types.RequestContext) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(&requestContext)
	store.Set(types.GetRequestContextKey(requestContextID), bz)
}

// DeleteRequestContext deletes the specified request context
func (k Keeper) DeleteRequestContext(ctx sdk.Context, requestContextID tmbytes.HexBytes) {
	store := ctx.KVStore(k.storeKey)

	store.Delete(types.GetRequestContextKey(requestContextID))
}

// GetRequestContext retrieves the specified request context
func (k Keeper) GetRequestContext(ctx sdk.Context, requestContextID tmbytes.HexBytes) (requestContext types.RequestContext, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetRequestContextKey(requestContextID))
	if bz == nil {
		return requestContext, false
	}

	k.cdc.MustUnmarshalBinaryBare(bz, &requestContext)
	return requestContext, true
}

// IterateRequestContexts iterates through all request contexts
func (k Keeper) IterateRequestContexts(
	ctx sdk.Context,
	op func(requestContextID tmbytes.HexBytes, requestContext types.RequestContext) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.RequestContextKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		requestContextID := iterator.Key()[1:]

		var requestContext types.RequestContext
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &requestContext)

		if stop := op(requestContextID, requestContext); stop {
			break
		}
	}
}

// InitiateRequests creates requests for the given providers from the specified request context
// Note: make sure that request context is valid and running, and providers are valid
func (k Keeper) InitiateRequests(
	ctx sdk.Context,
	requestContextID tmbytes.HexBytes,
	providers []sdk.AccAddress,
	providerRequests map[string][]string,
) []tmbytes.HexBytes {
	requestContext, _ := k.GetRequestContext(ctx, requestContextID)
	requestContext.BatchCounter++

	var requests []types.CompactRequest
	var requestIDs []tmbytes.HexBytes

	consumer, _ := sdk.AccAddressFromBech32(requestContext.Consumer)

	for providerIndex, provider := range providers {
		request := k.buildRequest(
			ctx, requestContextID, requestContext.BatchCounter,
			requestContext.ServiceName, provider, requestContext.SuperMode,
			consumer, requestContext.Timeout,
		)

		requestID := types.GenerateRequestID(requestContextID, requestContext.BatchCounter, ctx.BlockHeight(), int16(providerIndex))
		k.SetCompactRequest(ctx, requestID, request)

		k.AddActiveRequest(ctx, requestContext.ServiceName, provider, ctx.BlockHeight()+requestContext.Timeout, requestID)

		requests = append(requests, request)

		// tags for provider by one service
		providerRequests[types.ActionTag(requestContext.ServiceName, provider.String())] = append(
			providerRequests[types.ActionTag(requestContext.ServiceName, provider.String())],
			requestID.String(),
		)

		requestIDs = append(requestIDs, requestID)
	}

	requestContext.BatchState = types.BATCHRUNNING
	requestContext.BatchResponseCount = 0
	requestContext.BatchRequestCount = uint32(len(providers))
	requestContext.BatchResponseThreshold = requestContext.ResponseThreshold

	k.SetRequestContext(ctx, requestContextID, requestContext)

	if len(requests) > 0 {
		requestsJSON, _ := json.Marshal(requests)

		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeNewBatchRequest,
				sdk.NewAttribute(types.AttributeKeyServiceName, requestContext.ServiceName),
				sdk.NewAttribute(types.AttributeKeyRequestContextID, requestContextID.String()),
				sdk.NewAttribute(types.AttributeKeyRequests, string(requestsJSON)),
			),
		})
	}

	return requestIDs
}

// SkipCurrentRequestBatch skips the current request batch
func (k Keeper) SkipCurrentRequestBatch(ctx sdk.Context, requestContextID tmbytes.HexBytes, requestContext types.RequestContext) {
	requestContext.BatchCounter++
	requestContext.BatchState = types.BATCHRUNNING
	requestContext.BatchRequestCount = 0
	requestContext.BatchResponseCount = 0
	requestContext.BatchResponseThreshold = requestContext.ResponseThreshold

	k.SetRequestContext(ctx, requestContextID, requestContext)
	k.AddRequestBatchExpiration(ctx, requestContextID, ctx.BlockHeight()+requestContext.Timeout)
}

// buildRequest builds a request to the given provider from the specified request context
// Note: make sure that the binding exists
func (k Keeper) buildRequest(
	ctx sdk.Context,
	requestContextID tmbytes.HexBytes,
	batchCounter uint64,
	serviceName string,
	provider sdk.AccAddress,
	superMode bool,
	consumer sdk.AccAddress,
	timeout int64,
) types.CompactRequest {
	var serviceFee sdk.Coins

	if !superMode {
		binding, _ := k.GetServiceBinding(ctx, serviceName, provider)
		serviceFee = k.GetPrice(ctx, consumer, binding)
	}

	return types.NewCompactRequest(
		requestContextID,
		batchCounter,
		provider,
		serviceFee,
		ctx.BlockHeight(),
		ctx.BlockHeight()+timeout,
	)
}

// SetCompactRequest sets the specified compact request
func (k Keeper) SetCompactRequest(ctx sdk.Context, requestID tmbytes.HexBytes, request types.CompactRequest) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(&request)
	store.Set(types.GetRequestKey(requestID), bz)
}

// GetCompactRequest retrieves the specified compact request
func (k Keeper) GetCompactRequest(ctx sdk.Context, requestID tmbytes.HexBytes) (request types.CompactRequest, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetRequestKey(requestID))
	if bz == nil {
		return request, false
	}

	k.cdc.MustUnmarshalBinaryBare(bz, &request)
	return request, true
}

// DeleteCompactRequest deletes the specified compact request
func (k Keeper) DeleteCompactRequest(ctx sdk.Context, requestID tmbytes.HexBytes) {
	store := ctx.KVStore(k.storeKey)

	store.Delete(types.GetRequestKey(requestID))
}

// GetRequest returns the specified request
func (k Keeper) GetRequest(ctx sdk.Context, requestID tmbytes.HexBytes) (request types.Request, found bool) {
	compactRequest, found := k.GetCompactRequest(ctx, requestID)
	if !found {
		return request, false
	}

	provider, err := sdk.AccAddressFromBech32(compactRequest.Provider)
	if err != nil {
		return request, false
	}

	requestContextId, err := hex.DecodeString(compactRequest.RequestContextId)
	if err != nil {
		return request, false
	}

	requestContext, found := k.GetRequestContext(ctx, requestContextId)
	if !found {
		return request, false
	}

	consumer, err := sdk.AccAddressFromBech32(requestContext.Consumer)
	if err != nil {
		return request, false
	}

	return types.NewRequest(
		requestID,
		requestContext.ServiceName,
		provider,
		consumer,
		requestContext.Input,
		compactRequest.ServiceFee,
		requestContext.SuperMode,
		compactRequest.RequestHeight,
		compactRequest.ExpirationHeight,
		requestContextId,
		compactRequest.RequestContextBatchCounter,
	), true
}

// IterateRequests iterates through all compact requests
func (k Keeper) IterateRequests(
	ctx sdk.Context,
	op func(requestID tmbytes.HexBytes, request types.CompactRequest) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.RequestKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		requestID := iterator.Key()[1:]

		var request types.CompactRequest
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &request)

		if stop := op(requestID, request); stop {
			break
		}
	}
}

// RequestsIteratorByReqCtx returns an iterator for all requests of the specified request context ID and batch counter
func (k Keeper) RequestsIteratorByReqCtx(ctx sdk.Context, requestContextID tmbytes.HexBytes, batchCounter uint64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.GetRequestSubspaceByReqCtx(requestContextID, batchCounter))
}

// AddActiveRequest adds the specified active request
func (k Keeper) AddActiveRequest(
	ctx sdk.Context,
	serviceName string,
	provider sdk.AccAddress,
	expirationHeight int64,
	requestID tmbytes.HexBytes,
) {
	k.AddActiveRequestByBinding(ctx, serviceName, provider, expirationHeight, requestID)
	k.AddActiveRequestByID(ctx, requestID)
}

// DeleteActiveRequest deletes the specified active request
func (k Keeper) DeleteActiveRequest(
	ctx sdk.Context,
	serviceName string,
	provider sdk.AccAddress,
	expirationHeight int64,
	requestID tmbytes.HexBytes,
) {
	k.DeleteActiveRequestByBinding(ctx, serviceName, provider, expirationHeight, requestID)
	k.DeleteActiveRequestByID(ctx, requestID)
}

// AddActiveRequestByBinding adds the specified active request by the binding
func (k Keeper) AddActiveRequestByBinding(
	ctx sdk.Context,
	serviceName string,
	provider sdk.AccAddress,
	expirationHeight int64,
	requestID tmbytes.HexBytes,
) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(&gogotypes.BytesValue{Value: requestID})
	store.Set(types.GetActiveRequestKey(serviceName, provider, expirationHeight, requestID), bz)
}

// DeleteActiveRequestByBinding deletes the specified active request by the binding
func (k Keeper) DeleteActiveRequestByBinding(
	ctx sdk.Context,
	serviceName string,
	provider sdk.AccAddress,
	expirationHeight int64,
	requestID tmbytes.HexBytes,
) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetActiveRequestKey(serviceName, provider, expirationHeight, requestID))
}

// AddActiveRequestByID adds the specified active request by request ID
func (k Keeper) AddActiveRequestByID(ctx sdk.Context, requestID tmbytes.HexBytes) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(&gogotypes.BytesValue{Value: requestID})
	store.Set(types.GetActiveRequestKeyByID(requestID), bz)
}

// DeleteActiveRequestByID deletes the specified active request by request ID
func (k Keeper) DeleteActiveRequestByID(ctx sdk.Context, requestID tmbytes.HexBytes) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetActiveRequestKeyByID(requestID))
}

// IsRequestActive checks if the specified request is active
func (k Keeper) IsRequestActive(ctx sdk.Context, requestID tmbytes.HexBytes) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetActiveRequestKeyByID(requestID))
}

// AddRequestBatchExpiration adds a request batch to the expiration queue
func (k Keeper) AddRequestBatchExpiration(ctx sdk.Context, requestContextID tmbytes.HexBytes, expirationHeight int64) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(&gogotypes.BytesValue{Value: requestContextID})
	store.Set(types.GetExpiredRequestBatchKey(requestContextID, expirationHeight), bz)

	k.SetRequestBatchExpirationHeight(ctx, requestContextID, expirationHeight)
}

// DeleteRequestBatchExpiration deletes the request batch from the expiration queue
func (k Keeper) DeleteRequestBatchExpiration(ctx sdk.Context, requestContextID tmbytes.HexBytes, expirationHeight int64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetExpiredRequestBatchKey(requestContextID, expirationHeight))

	k.DeleteRequestBatchExpirationHeight(ctx, requestContextID)
}

// HasRequestBatchExpiration checks if the request batch expiration of the specified request context exists
func (k Keeper) HasRequestBatchExpiration(ctx sdk.Context, requestContextID tmbytes.HexBytes) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetExpiredRequestBatchHeightKey(requestContextID))
}

// AddNewRequestBatch adds a request batch to the new request batch queue
func (k Keeper) AddNewRequestBatch(ctx sdk.Context, requestContextID tmbytes.HexBytes, requestBatchHeight int64) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(&gogotypes.BytesValue{Value: requestContextID})
	store.Set(types.GetNewRequestBatchKey(requestContextID, requestBatchHeight), bz)

	k.SetNewRequestBatchHeight(ctx, requestContextID, requestBatchHeight)
}

// DeleteNewRequestBatch deletes the request batch in the given height from the new request batch queue
func (k Keeper) DeleteNewRequestBatch(ctx sdk.Context, requestContextID tmbytes.HexBytes, requestBatchHeight int64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetNewRequestBatchKey(requestContextID, requestBatchHeight))

	k.DeleteNewRequestBatchHeight(ctx, requestContextID)
}

// HasNewRequestBatch checks if the new request batch of the specified request context exists
func (k Keeper) HasNewRequestBatch(ctx sdk.Context, requestContextID tmbytes.HexBytes) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetNewRequestBatchHeightKey(requestContextID))
}

// SetRequestBatchExpirationHeight sets the request batch expiration height for the specified request context
func (k Keeper) SetRequestBatchExpirationHeight(ctx sdk.Context, requestContextID tmbytes.HexBytes, expirationHeight int64) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(&gogotypes.Int64Value{Value: expirationHeight})
	store.Set(types.GetExpiredRequestBatchHeightKey(requestContextID), bz)
}

// DeleteRequestBatchExpirationHeight deletes the request batch expiration height for the specified request context
func (k Keeper) DeleteRequestBatchExpirationHeight(ctx sdk.Context, requestContextID tmbytes.HexBytes) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetExpiredRequestBatchHeightKey(requestContextID))
}

// SetNewRequestBatchHeight sets the new request batch height for the specified request context
func (k Keeper) SetNewRequestBatchHeight(ctx sdk.Context, requestContextID tmbytes.HexBytes, requestBatchHeight int64) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(&gogotypes.Int64Value{Value: requestBatchHeight})
	store.Set(types.GetNewRequestBatchHeightKey(requestContextID), bz)
}

// DeleteNewRequestBatchHeight deletes the new request batch height for the specified request context
func (k Keeper) DeleteNewRequestBatchHeight(ctx sdk.Context, requestContextID tmbytes.HexBytes) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetNewRequestBatchHeightKey(requestContextID))
}

// IterateExpiredRequestBatch iterates through the expired request batch queue in the specified height
func (k Keeper) IterateExpiredRequestBatch(
	ctx sdk.Context,
	expirationHeight int64,
	op func(requestContextID tmbytes.HexBytes, requestContext types.RequestContext),
) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.GetExpiredRequestBatchSubspace(expirationHeight))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var requestContextID gogotypes.BytesValue
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &requestContextID)

		requestContext, _ := k.GetRequestContext(ctx, requestContextID.Value)

		op(requestContextID.Value, requestContext)
	}
}

// IterateNewRequestBatch iterates through the new request batch queue in the specified height
func (k Keeper) IterateNewRequestBatch(
	ctx sdk.Context,
	requestBatchHeight int64,
	op func(requestContextID tmbytes.HexBytes, requestContext *types.RequestContext),
) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.GetNewRequestBatchSubspace(requestBatchHeight))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var requestContextID gogotypes.BytesValue
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &requestContextID)

		requestContext, _ := k.GetRequestContext(ctx, requestContextID.Value)

		op(requestContextID.Value, &requestContext)
	}
}

// ActiveRequestsIterator returns an iterator for all the active requests of the specified service binding
func (k Keeper) ActiveRequestsIterator(ctx sdk.Context, serviceName string, provider sdk.AccAddress) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.GetActiveRequestSubspace(serviceName, provider))
}

// ActiveRequestsIteratorByReqCtx returns an iterator for all the active requests of the specified service binding
func (k Keeper) ActiveRequestsIteratorByReqCtx(ctx sdk.Context, requestContextID tmbytes.HexBytes, batchCounter uint64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.GetActiveRequestSubspaceByReqCtx(requestContextID, batchCounter))
}

// AllActiveRequestsIterator returns an iterator for all the active requests
func (k Keeper) AllActiveRequestsIterator(store sdk.KVStore) sdk.Iterator {
	return sdk.KVStorePrefixIterator(store, types.ActiveRequestKey)
}

// IterateActiveRequests iterates through the active requests for the specified request context ID and batch counter
func (k Keeper) IterateActiveRequests(
	ctx sdk.Context,
	requestContextID tmbytes.HexBytes,
	batchCounter uint64,
	op func(requestID tmbytes.HexBytes, request types.Request),
) {
	iterator := k.ActiveRequestsIteratorByReqCtx(ctx, requestContextID, batchCounter)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var requestID gogotypes.BytesValue
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &requestID)

		request, _ := k.GetRequest(ctx, requestID.Value)

		op(requestID.Value, request)
	}
}

// FilterServiceProviders gets the providers which satisfy the specified requirement
func (k Keeper) FilterServiceProviders(
	ctx sdk.Context,
	serviceName string,
	providers []sdk.AccAddress,
	timeout int64,
	serviceFeeCap sdk.Coins,
	consumer sdk.AccAddress,
) (
	[]sdk.AccAddress, sdk.Coins, string, error,
) {
	var newProviders []sdk.AccAddress
	var totalPrices sdk.Coins

	for _, provider := range providers {
		binding, found := k.GetServiceBinding(ctx, serviceName, provider)

		if found && binding.Available {
			if binding.QoS <= uint64(timeout) {
				exchangedPrice, rawDenom, err := k.GetExchangedPrice(ctx, consumer, binding)
				if err != nil {
					return nil, nil, rawDenom, err
				}

				provider, err := sdk.AccAddressFromBech32(binding.Provider)
				if err != nil {
					return nil, nil, rawDenom, fmt.Errorf("invalid provider address: %s", provider)
				}

				price := k.GetPricing(ctx, binding.ServiceName, provider).Price
				if exchangedPrice.IsAllLTE(serviceFeeCap) {
					newProviders = append(newProviders, provider)
					totalPrices = totalPrices.Add(price...)
				}
			}
		}
	}

	return newProviders, totalPrices, "", nil
}

// DeductServiceFees deducts the given service fees from the specified consumer
func (k Keeper) DeductServiceFees(ctx sdk.Context, consumer sdk.AccAddress, serviceFees sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromAccountToModule(ctx, consumer, types.RequestAccName, serviceFees)
}

func (k Keeper) GetPrice(
	ctx sdk.Context,
	consumer sdk.AccAddress,
	binding types.ServiceBinding,
) sdk.Coins {
	provider, _ := sdk.AccAddressFromBech32(binding.Provider)
	pricing := k.GetPricing(ctx, binding.ServiceName, provider)

	// get discounts
	discountByTime := types.GetDiscountByTime(pricing, ctx.BlockTime())
	discountByVolume := types.GetDiscountByVolume(
		pricing, k.GetRequestVolume(ctx, consumer, binding.ServiceName, provider),
	)

	var price []sdk.Coin
	for _, token := range pricing.Price {
		priceAmount := sdk.NewDecFromInt(token.Amount).Mul(discountByTime).Mul(discountByVolume)
		price = append(price, sdk.NewCoin(token.Denom, priceAmount.TruncateInt()))
	}

	return sdk.NewCoins(price...)
}

// AddResponse adds the response for the specified request ID
func (k Keeper) AddResponse(
	ctx sdk.Context,
	requestID tmbytes.HexBytes,
	provider sdk.AccAddress,
	result string,
	output string,
) (
	request types.Request,
	response types.Response,
	err error,
) {
	request, found := k.GetRequest(ctx, requestID)
	if !found {
		return request, response, sdkerrors.Wrap(types.ErrUnknownRequest, requestID.String())
	}

	requestProvider, _ := sdk.AccAddressFromBech32(request.Provider)
	if !provider.Equals(requestProvider) {
		return request, response, sdkerrors.Wrap(types.ErrInvalidResponse, "provider does not match")
	}

	if !k.IsRequestActive(ctx, requestID) {
		return request, response, sdkerrors.Wrap(types.ErrInvalidResponse, "request is not active")
	}

	if err := types.ValidateResponseOutput(output); err != nil {
		return request, response, sdkerrors.Wrap(types.ErrInvalidResponseOutput, err.Error())
	}

	if err := k.AddEarnedFee(ctx, provider, request.ServiceFee); err != nil {
		return request, response, err
	}

	requestContextID, _ := hex.DecodeString(request.RequestContextId)
	consumer, _ := sdk.AccAddressFromBech32(request.Consumer)

	response = types.NewResponse(provider, consumer, result, output, requestContextID, request.RequestContextBatchCounter)
	k.SetResponse(ctx, requestID, response)

	k.DeleteActiveRequest(ctx, request.ServiceName, provider, request.ExpirationHeight, requestID)
	k.IncreaseRequestVolume(ctx, consumer, request.ServiceName, provider)

	requestContext, _ := k.GetRequestContext(ctx, requestContextID)
	requestContext.BatchResponseCount++

	if requestContext.BatchResponseCount == requestContext.BatchRequestCount {
		requestContext = k.CompleteBatch(ctx, requestContext, requestContextID)
	}

	k.SetRequestContext(ctx, requestContextID, requestContext)

	return request, response, nil
}

// Callback callbacks the corresponding response callback handler
func (k Keeper) Callback(ctx sdk.Context, requestContextID tmbytes.HexBytes) {
	requestContext, _ := k.GetRequestContext(ctx, requestContextID)

	respCallback, _ := k.GetResponseCallback(requestContext.ModuleName)
	outputs := k.GetResponseOutputs(ctx, requestContextID, requestContext.BatchCounter)

	if len(outputs) >= int(requestContext.BatchResponseThreshold) {
		respCallback(ctx, requestContextID, outputs, nil)
	} else {
		respCallback(
			ctx, requestContextID, outputs,
			fmt.Errorf(
				"batch %d at least %d valid outputs required, but %d received",
				requestContext.BatchCounter, requestContext.BatchResponseThreshold, len(outputs),
			),
		)
	}
}

// SetResponse sets the specified response
func (k Keeper) SetResponse(ctx sdk.Context, requestID tmbytes.HexBytes, response types.Response) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(&response)
	store.Set(types.GetResponseKey(requestID), bz)
}

// GetResponse returns a response with the speicified request ID
func (k Keeper) GetResponse(ctx sdk.Context, requestID tmbytes.HexBytes) (response types.Response, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetResponseKey(requestID))
	if bz == nil {
		return response, false
	}

	k.cdc.MustUnmarshalBinaryBare(bz, &response)
	return response, true
}

// DeleteResponse deletes a response with the speicified request ID
func (k Keeper) DeleteResponse(ctx sdk.Context, requestID tmbytes.HexBytes) {
	store := ctx.KVStore(k.storeKey)

	store.Delete(types.GetResponseKey(requestID))
}

// IterateResponses iterates through all responses
func (k Keeper) IterateResponses(
	ctx sdk.Context,
	op func(requestID tmbytes.HexBytes, response types.Response) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.ResponseKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		requestID := iterator.Key()[1:]

		var response types.Response
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &response)

		if stop := op(requestID, response); stop {
			break
		}
	}
}

// ResponsesIteratorByReqCtx returns an iterator for all responses of the specified request context and batch counter
func (k Keeper) ResponsesIteratorByReqCtx(ctx sdk.Context, requestContextID tmbytes.HexBytes, batchCounter uint64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.GetResponseSubspaceByReqCtx(requestContextID, batchCounter))
}

// GetResponseOutputs retrieves all response outputs of the specified request context and batch counter
func (k Keeper) GetResponseOutputs(ctx sdk.Context, requestContextID tmbytes.HexBytes, batchCounter uint64) []string {
	iterator := k.ResponsesIteratorByReqCtx(ctx, requestContextID, batchCounter)
	defer iterator.Close()

	var outputs []string
	for ; iterator.Valid(); iterator.Next() {
		var response types.Response
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &response)

		if len(response.Output) > 0 {
			outputs = append(outputs, response.Output)
		}
	}

	return outputs
}

// IncreaseRequestVolume increases the request volume by 1
func (k Keeper) IncreaseRequestVolume(
	ctx sdk.Context,
	consumer sdk.AccAddress,
	serviceName string,
	provider sdk.AccAddress,
) {
	currentVolume := k.GetRequestVolume(ctx, consumer, serviceName, provider)
	k.SetRequestVolume(ctx, consumer, serviceName, provider, currentVolume+1)
}

// SetRequestVolume sets the request volume for the specified consumer and binding
func (k Keeper) SetRequestVolume(
	ctx sdk.Context,
	consumer sdk.AccAddress,
	serviceName string,
	provider sdk.AccAddress,
	volume uint64,
) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(&gogotypes.UInt64Value{Value: volume})
	store.Set(types.GetRequestVolumeKey(consumer, serviceName, provider), bz)
}

// GetRequestVolume gets the current request volume for the specified consumer and binding
func (k Keeper) GetRequestVolume(
	ctx sdk.Context,
	consumer sdk.AccAddress,
	serviceName string,
	provider sdk.AccAddress,
) uint64 {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetRequestVolumeKey(consumer, serviceName, provider))
	if bz == nil {
		return 0
	}

	var volume gogotypes.UInt64Value
	k.cdc.MustUnmarshalBinaryBare(bz, &volume)

	return volume.Value
}

// Slash slashes the provider from the specified request
// Note: ensure that the request is valid
func (k Keeper) Slash(ctx sdk.Context, requestID tmbytes.HexBytes) error {
	request, _ := k.GetRequest(ctx, requestID)
	provider, _ := sdk.AccAddressFromBech32(request.Provider)
	binding, _ := k.GetServiceBinding(ctx, request.ServiceName, provider)

	slashFraction := k.SlashFraction(ctx)
	baseDenom := k.BaseDenom(ctx)

	depositAmt := binding.Deposit.AmountOf(baseDenom)
	slashedAmt := sdk.NewDecFromInt(depositAmt).Mul(slashFraction).TruncateInt()
	slashedCoins := sdk.NewCoins(sdk.NewCoin(baseDenom, slashedAmt))

	deposit, hasNeg := binding.Deposit.SafeSub(slashedCoins)
	if hasNeg {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "%s is less than %s", binding.Deposit.String(), slashedCoins.String())
	}

	if err := k.bankKeeper.BurnCoins(ctx, types.DepositAccName, slashedCoins); err != nil {
		return err
	}

	binding.Deposit = deposit
	if binding.Available {
		minDeposit := k.getMinDeposit(ctx, k.GetPricing(ctx, binding.ServiceName, provider))
		if !binding.Deposit.IsAllGTE(minDeposit) {
			binding.Available = false
			binding.DisabledTime = ctx.BlockHeader().Time
		}
	}

	k.SetServiceBinding(ctx, binding)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeServiceSlash,
			sdk.NewAttribute(types.AttributeKeyRequestID, requestID.String()),
			sdk.NewAttribute(types.AttributeKeyProvider, request.Provider),
			sdk.NewAttribute(types.AttributeKeySlashedCoins, slashedCoins.String()),
		),
	})

	return nil
}

// CheckAuthority checks if the operation on the specified request context is authorized
func (k Keeper) CheckAuthority(
	ctx sdk.Context,
	consumer sdk.AccAddress,
	requestContextID tmbytes.HexBytes,
	checkModule bool,
) error {
	requestContext, found := k.GetRequestContext(ctx, requestContextID)
	if !found {
		return sdkerrors.Wrap(types.ErrUnknownRequestContext, requestContextID.String())
	}

	if !consumer.Equals(consumer) {
		return sdkerrors.Wrapf(types.ErrNotAuthorized, "consumer not matching")
	}

	if checkModule && len(requestContext.ModuleName) > 0 {
		return sdkerrors.Wrapf(types.ErrNotAuthorized, "not authorized operation")
	}

	return nil
}

// GetResponseCallback gets the registered module callback for response handling
func (k Keeper) GetResponseCallback(moduleName string) (types.ResponseCallback, error) {
	respCallback, ok := k.respCallbacks[moduleName]
	if !ok {
		return nil, sdkerrors.Wrapf(types.ErrCallbackNotRegistered, "%s not registered for module %s", "response callback", moduleName)
	}

	return respCallback, nil
}

// GetStateCallback gets the registered module callback for state handling
func (k Keeper) GetStateCallback(moduleName string) (types.StateCallback, error) {
	stateCallback, ok := k.stateCallbacks[moduleName]
	if !ok {
		return nil, sdkerrors.Wrapf(types.ErrCallbackNotRegistered, "%s not registered for module %s", "state callback", moduleName)
	}

	return stateCallback, nil
}

// ResetRequestContextsStateAndBatch reset request contexts state and batch
func (k Keeper) ResetRequestContextsStateAndBatch(ctx sdk.Context) error {
	k.IterateRequestContexts(
		ctx,
		func(requestContextID tmbytes.HexBytes, requestContext types.RequestContext) bool {
			requestContext.State = types.PAUSED

			requestContext.BatchState = types.BATCHCOMPLETED
			requestContext.BatchRequestCount = 0
			requestContext.BatchResponseCount = 0

			k.SetRequestContext(ctx, requestContextID, requestContext)
			return false
		},
	)

	return nil
}

// validateServiceFeeCap validates the given service fee cap
func (k Keeper) validateServiceFeeCap(ctx sdk.Context, serviceFeeCap sdk.Coins) error {
	baseDenom := k.BaseDenom(ctx)

	if len(serviceFeeCap) != 1 || serviceFeeCap[0].Denom != baseDenom {
		return sdkerrors.Wrapf(types.ErrInvalidServiceFeeCap, "service fee cap only accepts %s", baseDenom)
	}

	return nil
}

func TxHash(ctx sdk.Context) []byte {
	return tmhash.Sum(ctx.TxBytes())
}

//  GetInternalIndex sets the internal index
func (k Keeper) SetInternalIndex(ctx sdk.Context, index int64) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(&gogotypes.Int64Value{
		Value: index,
	})
	store.Set(types.InternalCounterKey, bz)
}

// GetInternalIndex returns the internal index and increases the internal index + 1
func (k Keeper) GetInternalIndex(ctx sdk.Context) int64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.InternalCounterKey)
	if bz == nil {
		return 0
	}

	var index gogotypes.Int64Value
	k.cdc.MustUnmarshalBinaryBare(bz, &index)
	k.SetInternalIndex(ctx, index.Value+1)
	return index.Value
}
