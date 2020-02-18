package keeper

import (
	"fmt"

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
	repeated bool,
	repeatedFrequency uint64,
	repeatedTotal int64,
	state types.RequestContextState,
	respThreshold uint16,
	respHandler string,
) ([]byte, sdk.Error) {
	svcDef, found := k.GetServiceDefinition(ctx, serviceName)
	if !found {
		return nil, types.ErrUnknownServiceDefinition(k.codespace, serviceName)
	}

	if err := types.ValidateRequestInput(svcDef.Schemas, input); err != nil {
		return nil, err
	}

	params := k.GetParamSet(ctx)
	if timeout > params.MaxRequestTimeout {
		return nil, types.ErrInvalidRequest(k.codespace, fmt.Sprintf("timeout must not be greater than %d: %d", params.MaxRequestTimeout, timeout))
	}

	if timeout == 0 {
		timeout = params.MaxRequestTimeout
	}

	if repeated {
		if repeatedFrequency == 0 {
			repeatedFrequency = uint64(timeout)
		}

		if repeatedFrequency < uint64(timeout) {
			return nil, types.ErrInvalidRequest(k.codespace, fmt.Sprintf("repeated frequency [%d] must not be less than timeout [%d]", repeatedFrequency, timeout))
		}
	} else {
		repeatedFrequency = 0
		repeatedTotal = 0
	}

	batchCounter := uint64(0)

	requestContext := types.NewRequestContext(
		serviceName, providers, consumer, input, serviceFeeCap,
		timeout, repeated, repeatedFrequency, repeatedTotal,
		batchCounter, state, respThreshold, respHandler,
	)

	requestContextID := types.GenerateRequestContextID(ctx.BlockHeight(), k.GetIntraTxCounter(ctx))
	k.SetRequestContext(ctx, requestContextID, requestContext)

	return requestContextID, nil
}

// UpdateRequestContext updates the specified request context
func (k Keeper) UpdateRequestContext(
	ctx sdk.Context,
	requestContextID []byte,
	providers []sdk.AccAddress,
	serviceFeeCap sdk.Coins,
	repeatedFreq uint64,
	repeatedTotal int64,
) sdk.Error {
	requestContext, found := k.GetRequestContext(ctx, requestContextID)
	if !found {
		return types.ErrInvalidRequestContextID(k.codespace, "invalid request context ID")
	}

	if !requestContext.Repeated {
		return types.ErrRequestContextNonRepeated(k.codespace)
	}

	if len(providers) > 0 && requestContext.ResponseThreshold > 0 && len(providers) < int(requestContext.ResponseThreshold) {
		return types.ErrInvalidProviders(k.codespace, "length of providers must not be less than the response threshold")
	}

	if repeatedFreq > 0 && repeatedFreq < uint64(requestContext.Timeout) {
		return types.ErrInvalidRepeatedFreq(k.codespace, "repeated frequency must not be less than the timeout")
	}

	if repeatedTotal >= 1 && repeatedTotal <= int64(requestContext.BatchCounter) {
		return types.ErrInvalidRepeatedTotal(k.codespace, "updated repeated total must be greater than the current batch counter")
	}

	if len(providers) > 0 {
		requestContext.Providers = providers
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
	requestContextID []byte,
) sdk.Error {
	requestContext, found := k.GetRequestContext(ctx, requestContextID)
	if !found {
		return types.ErrInvalidRequestContextID(k.codespace, "invalid request context ID")
	}

	if !requestContext.Repeated {
		return types.ErrRequestContextNonRepeated(k.codespace)
	}

	if requestContext.State != types.RequestContextState(0x00) {
		return types.ErrRequestContextNotStarted(k.codespace)
	}

	requestContext.State = types.RequestContextState(0x01)
	k.DeleteNewBatchRequest(ctx, requestContextID)

	k.SetRequestContext(ctx, requestContextID, requestContext)

	return nil
}

// StartRequestContext starts the specified request context
func (k Keeper) StartRequestContext(
	ctx sdk.Context,
	requestContextID []byte,
) sdk.Error {
	requestContext, found := k.GetRequestContext(ctx, requestContextID)
	if !found {
		return types.ErrInvalidRequestContextID(k.codespace, "invalid request context ID")
	}

	if !requestContext.Repeated {
		return types.ErrRequestContextNonRepeated(k.codespace)
	}

	if requestContext.State != types.RequestContextState(0x01) {
		return types.ErrRequestContextNotPaused(k.codespace)
	}

	requestContext.State = types.RequestContextState(0x00)
	k.AddNewBatchRequest(ctx, requestContextID)

	k.SetRequestContext(ctx, requestContextID, requestContext)

	return nil
}

// KillRequestContext terminates the specified request context
func (k Keeper) KillRequestContext(
	ctx sdk.Context,
	requestContextID []byte,
) sdk.Error {
	requestContext, found := k.GetRequestContext(ctx, requestContextID)
	if !found {
		return types.ErrInvalidRequestContextID(k.codespace, "invalid request context ID")
	}

	if !requestContext.Repeated {
		return types.ErrRequestContextNonRepeated(k.codespace)
	}

	requestContext.State = types.RequestContextState(0x01)
	requestContext.RepeatedTotal = int64(requestContext.BatchCounter)

	k.SetRequestContext(ctx, requestContextID, requestContext)

	return nil
}

// SetRequestContext sets the specified request context
func (k Keeper) SetRequestContext(ctx sdk.Context, requestContextID []byte, requestContext types.RequestContext) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(requestContext)
	store.Set(GetRequestContextKey(requestContextID), bz)
}

// GetRequestContext retrieves the specified request context
func (k Keeper) GetRequestContext(ctx sdk.Context, requestContextID []byte) (requestContext types.RequestContext, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(GetRequestContextKey(requestContextID))
	if bz == nil {
		return requestContext, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &requestContext)
	return requestContext, true
}

// InitiateRequests creates requests for the given providers from the specified request context
func (k Keeper) InitiateRequests(
	ctx sdk.Context,
	requestContextID []byte,
	providers []sdk.AccAddress,
) sdk.Error {
	requestContext, found := k.GetRequestContext(ctx, requestContextID)
	if !found {
		return types.ErrInvalidRequestContextID(k.codespace, "invalid request context ID")
	}

	if requestContext.State != types.RequestContextState(0x00) {
		return types.ErrRequestContextNotStarted(k.codespace)
	}

	for providerIndex, provider := range providers {
		request, err := k.buildRequest(ctx, requestContextID, requestContext.BatchCounter, requestContext.ServiceName, provider)
		if err != nil {
			return err
		}

		requestID := types.GenerateRequestID(requestContextID, requestContext.BatchCounter, int16(providerIndex))

		k.SetCompactRequest(ctx, requestID, request)
		k.AddActiveRequest(ctx, requestContext.ServiceName, ctx.BlockHeight()+requestContext.Timeout, requestID, request)
	}

	return nil
}

// buildRequest builds a request for the given provider from the specified request context
func (k Keeper) buildRequest(
	ctx sdk.Context,
	requestContextID []byte,
	batchCounter uint64,
	serviceName string,
	provider sdk.AccAddress,
) (request types.CompactRequest, err sdk.Error) {
	_, found := k.GetServiceBinding(ctx, serviceName, provider)
	if !found {
		return request, types.ErrUnknownServiceBinding(k.codespace)
	}

	// TODO: extract price from binding.Pricing
	serviceFee := sdk.Coins{}

	request = types.NewCompactRequest(
		requestContextID,
		batchCounter,
		provider,
		serviceFee,
		ctx.BlockHeight(),
	)

	return request, nil
}

// SetCompactRequest sets the specified compact request
func (k Keeper) SetCompactRequest(ctx sdk.Context, requestID []byte, request types.CompactRequest) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(request)
	store.Set(GetRequestKey(requestID), bz)
}

// GetCompactRequest retrieves the specified compact request
func (k Keeper) GetCompactRequest(ctx sdk.Context, requestID []byte) (request types.CompactRequest, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(GetRequestKey(requestID))
	if bz == nil {
		return request, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &request)
	return request, true
}

// GetRequest returns the specified request
func (k Keeper) GetRequest(ctx sdk.Context, requestID []byte) (request types.Request, found bool) {
	compactRequest, found := k.GetCompactRequest(ctx, requestID)
	if !found {
		return request, false
	}

	requestContext, found := k.GetRequestContext(ctx, compactRequest.RequestContextID)
	if !found {
		return request, false
	}

	request = types.NewRequest(
		requestContext.ServiceName,
		compactRequest.Provider,
		requestContext.Consumer,
		requestContext.Input,
		compactRequest.ServiceFee,
		requestContext.Profiling,
		compactRequest.RequestHeight,
		compactRequest.RequestHeight+requestContext.Timeout,
		compactRequest.RequestContextID,
		compactRequest.RequestContextBatchCounter,
	)

	return request, true
}

// AddActiveRequest adds the specified active request
func (k Keeper) AddActiveRequest(
	ctx sdk.Context,
	serviceName string,
	expirationHeight int64,
	requestID []byte,
	request types.CompactRequest,
) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(request)
	store.Set(GetActiveRequestKey(serviceName, request.Provider, expirationHeight, requestID), bz)
}

// DeleteActiveRequest deletes the specified active request
func (k Keeper) DeleteActiveRequest(
	ctx sdk.Context,
	serviceName string,
	provider sdk.AccAddress,
	expirationHeight int64,
	requestID []byte,
) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetActiveRequestKey(serviceName, provider, expirationHeight, requestID))
}

// IsRequestActive checks if the specified request is active
func (k Keeper) IsRequestActive(
	ctx sdk.Context,
	serviceName string,
	provider sdk.AccAddress,
	expirationHeight int64,
	requestID []byte,
) bool {
	store := ctx.KVStore(k.storeKey)
	store.Has(GetActiveRequestKey(serviceName, provider, expirationHeight, requestID))
}

// AddRequestBatchExpiration adds a request batch to the expiration queue
func (k Keeper) AddRequestBatchExpiration(ctx sdk.Context, requestContextID []byte, expirationHeight int64) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(requestContextID)
	store.Set(GetExpiredRequestBatchKey(requestContextID, expirationHeight), bz)
}

// DeleteRequestBatchExpiration deletes the request batch from the expiration queue
func (k Keeper) DeleteRequestBatchExpiration(ctx sdk.Context, requestContextID []byte, expirationHeight int64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetExpiredRequestBatchKey(requestContextID, expirationHeight))
}

// AddNewRequestBatch adds a request batch to the new request batch queue
func (k Keeper) AddNewRequestBatch(ctx sdk.Context, requestContextID []byte, requestBatchHeight int64) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(requestContextID)
	store.Set(GetNewRequestBatchKey(requestContextID, requestBatchHeight), bz)
}

// DeleteNewRequestBatch deletes the request batch in the given height from the new request batch queue
func (k Keeper) DeleteNewRequestBatch(ctx sdk.Context, requestContextID []byte, requestBatchHeight int64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetNewRequestBatchKey(requestContextID, requestBatchHeight))
}

// ExpiredRequestBatchIterator returns an iterator for the request batch expiration queue
func (k Keeper) ExpiredRequestBatchIterator(ctx sdk.Context, expirationHeight int64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, GetExpiredRequestBatchKey(expirationHeight))
}

// NewRequestBatchIterator returns an iterator for the new request batch queue
func (k Keeper) NewRequestBatchIterator(ctx sdk.Context, requestBatchHeight int64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, GetNewRequestBatchKey(requestBatchHeight))
}

// ActiveRequestsIterator returns an iterator for all the request in the Active Queue of specified service binding
func (k Keeper) ActiveRequestsIterator(ctx sdk.Context, serviceName string, provider sdk.AccAddress) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, GetActiveRequestSubspace(serviceName, provider))
}

// FilterServiceProviders gets the providers which satisfy the specified service fee requirement
// Note: make sure that the binding exists for every provider
func (k Keeper) FilterServiceProviders(
	ctx sdk.Context,
	serviceName string,
	providers []sdk.AccAddress,
	serviceFeeCap sdk.Coins,
) ([]sdk.AccAddress, sdk.Coins) {
	var newProviders []sdk.AccAddress
	var totalPrice sdk.Coins

	for _, provider := range providers {
		binding, _ := k.GetServiceBinding(ctx, serviceName, provider)
		if binding.Pricing <= serviceFeeCap {
			newProviders = append(newProviders, provider)
			totalPrice = totalPrice.Add(binding.Pricing)
		}
	}

	return newProviders, totalPrice
}

// DeductServiceFees deducts the given service fees from the specified consumer
func (k Keeper) DeductServiceFees(ctx sdk.Context, consumer sdk.AccAddress, serviceFees sdk.Coins) sdk.Error {
	_, err := k.bk.SendCoins(ctx, consumer, auth.ServiceRequestCoinsAccAddr, serviceFees)
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
	output,
	errMsg string,
) (response types.Response, err sdk.Error) {
	reqID, _ := types.ConvertRequestID(requestID)

	request, found := k.GetRequest(ctx, reqID)
	if !found {
		return response, types.ErrInvalidRequestID(k.codespace, requestID)
	}

	if !provider.Equals(request.Provider) {
		return response, types.ErrInvalidResponse(k.codespace, "provider does not match")
	}

	if !k.IsRequestActive(ctx, request.ServiceName, provider, request.ExpirationHeight, reqID) {
		return response, types.ErrInvalidResponse(k.codespace, "request is not active")
	}

	svcDef, _ := k.GetServiceDefinition(ctx, request.ServiceName)

	if len(output) > 0 {
		if err := types.ValidateResponseOutput(svcDef.Schemas, output); err != nil {
			return response, err
		}
	} else {
		if err := types.ValidateResponseError(svcDef.Schemas, errMsg); err != nil {
			return response, err
		}
	}

	if err := k.AddEarnedFee(ctx, provider, request.ServiceFee); err != nil {
		return response, err
	}

	response = types.NewResponse(provider, request.Consumer, output, errMsg, request.RequestContextID, request.RequestContextBatchCounter)
	k.SetResponse(ctx, reqID, response)

	k.DeleteActiveRequest(ctx, request.ServiceName, provider, request.ExpirationHeight, reqID)

	return response, nil
}

// SetResponse sets the specified response
func (k Keeper) SetResponse(ctx sdk.Context, requestID []byte, response types.Response) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(response)
	store.Set(GetResponseKey(requestID), bz)
}

// GetResponse returns a response with the speicified request ID
func (k Keeper) GetResponse(ctx sdk.Context, requestID []byte) (response types.Response, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(GetResponseKey(requestID))
	if bz == nil {
		return response, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &response)
	return response, true
}

// Slash
func (k Keeper) Slash(ctx sdk.Context, binding types.ServiceBinding, slashCoins sdk.Coins) sdk.Error {
	deposit, hasNeg := binding.Deposit.SafeSub(slashCoins)
	if hasNeg {
		errMsg := fmt.Sprintf("%s is less than %s", binding.Deposit, slashCoins)
		panic(errMsg)
	}

	binding.Deposit = deposit
	minDeposit := k.getMinDeposit(ctx, binding.Pricing)

	if !binding.Deposit.IsAllGTE(minDeposit) {
		binding.Available = false
		binding.DisabledTime = ctx.BlockHeader().Time
	}

	_, err := k.bk.BurnCoins(ctx, auth.ServiceDepositCoinsAccAddr, slashCoins)
	if err != nil {
		return err
	}

	ctx.Logger().Info("Slash service provider", "provider", binding.Provider.String(), "slash_amount", slashCoins.String())
	k.SetServiceBinding(ctx, binding)

	return nil
}

// RefundServiceFee refunds the service fee to the specified consumer
func (k Keeper) RefundServiceFee(ctx sdk.Context, consumer sdk.AccAddress, serviceFee sdk.Coins) sdk.Error {
	_, err := k.bk.SendCoins(ctx, auth.ServiceRequestCoinsAccAddr, consumer, serviceFee)
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
		errMsg := fmt.Sprintf("%s is less than %s", coins, taxCoins)
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

// WithdrawEarnedFees withdraws the earned fees to the specified provider
func (k Keeper) WithdrawEarnedFees(ctx sdk.Context, provider sdk.AccAddress) sdk.Error {
	fees, found := k.GetEarnedFees(ctx, provider)
	if !found {
		return types.ErrNoEarnedFees(k.codespace, provider)
	}

	_, err := k.bk.SendCoins(ctx, auth.ServiceRequestCoinsAccAddr, provider, fees.Coins)
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
	store := ctx.KVStore(k.storeKey)

	iterator := k.ActiveAllRequestQueueIterator(store)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var request types.CompactRequest
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &request)

		_, err := k.bk.SendCoins(ctx, auth.ServiceRequestCoinsAccAddr, request.Consumer, request.ServiceFee)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetResponseCallback gets the registered module callback for response handling
func (k Keeper) GetResponseCallback(moduleName string) (types.ResponseCallback, sdk.Error) {
	respCallback, ok := k.respCallbacks[moduleName]
	if !ok {
		return nil, types.ErrModuleNameNotRegistered(k.Codespace(), moduleName)
	}

	return respCallback, nil
}

// GetIntraTxCounter gets the current tx counter
func (k Keeper) GetIntraTxCounter(ctx sdk.Context) int16 {
	v := ctx.Value(GetIntraTxCounterKey())
	if v == nil {
		return 0
	}

	return v.(int16)
}

// SetIntraTxCounter sets the tx counter to the context
func (k Keeper) SetIntraTxCounter(ctx sdk.Context, counter int16) {
	ctx = ctx.WithValue(GetIntraTxCounterKey(), counter)
}
