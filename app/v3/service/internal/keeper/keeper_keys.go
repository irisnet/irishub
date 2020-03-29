package keeper

import (
	sdk "github.com/irisnet/irishub/types"
)

var (
	// the separator for string key
	emptyByte = []byte{0x00}

	// Keys for store prefixes
	serviceDefinitionKey   = []byte{0x01}
	serviceBindingKey      = []byte{0x02}
	pricingKey             = []byte{0x03}
	withdrawAddrKey        = []byte{0x04}
	requestContextKey      = []byte{0x05}
	expiredRequestBatchKey = []byte{0x06}
	newRequestBatchKey     = []byte{0x07}
	requestKey             = []byte{0x08}
	activeRequestKey       = []byte{0x09}
	activeRequestByIDKey   = []byte{0x10}
	responseKey            = []byte{0x11}
	requestVolumeKey       = []byte{0x12}
	earnedFeesKey          = []byte{0x13}
)

// GetServiceDefinitionKey returns the key for the service definition with the specified name
func GetServiceDefinitionKey(name string) []byte {
	return append(serviceDefinitionKey, []byte(name)...)
}

// GetServiceBindingKey returns the key for the service binding with the specified name and provider
func GetServiceBindingKey(serviceName string, provider sdk.AccAddress) []byte {
	return append(serviceBindingKey, getStringsKey([]string{serviceName, provider.String()})...)
}

// GetPricingKey returns the key for the pricing of the specified binding
func GetPricingKey(serviceName string, provider sdk.AccAddress) []byte {
	return append(pricingKey, getStringsKey([]string{serviceName, provider.String()})...)
}

// GetWithdrawAddrKey returns the key for the withdrawal address of the specified provider
func GetWithdrawAddrKey(provider sdk.AccAddress) []byte {
	return append(withdrawAddrKey, provider.Bytes()...)
}

// GetBindingsSubspace returns the key for retrieving all bindings of the specified service
func GetBindingsSubspace(serviceName string) []byte {
	return append(append(serviceBindingKey, []byte(serviceName)...), emptyByte...)
}

// GetRequestContextKey returns the key for the request context with the specified ID
func GetRequestContextKey(requestContextID []byte) []byte {
	return append(requestContextKey, requestContextID...)
}

// GetExpiredRequestBatchKey returns the key for the request batch expiration of the specified request context
func GetExpiredRequestBatchKey(requestContextID []byte, batchExpirationHeight int64) []byte {
	reqBatchExpiration := append(sdk.Uint64ToBigEndian(uint64(batchExpirationHeight)), requestContextID...)
	return append(expiredRequestBatchKey, reqBatchExpiration...)
}

// GetNewRequestBatchKey returns the key for the new batch request of the specified request context in the given height
func GetNewRequestBatchKey(requestContextID []byte, requestBatchHeight int64) []byte {
	newBatchRequest := append(sdk.Uint64ToBigEndian(uint64(requestBatchHeight)), requestContextID...)
	return append(newRequestBatchKey, newBatchRequest...)
}

// GetExpiredRequestBatchSubspace returns the key for iterating through the expired request batch queue in the specified height
func GetExpiredRequestBatchSubspace(batchExpirationHeight int64) []byte {
	return append(expiredRequestBatchKey, sdk.Uint64ToBigEndian(uint64(batchExpirationHeight))...)
}

// GetNewRequestBatchSubspace returns the key for iterating through the new request batch queue in the specified height
func GetNewRequestBatchSubspace(requestBatchHeight int64) []byte {
	return append(newRequestBatchKey, sdk.Uint64ToBigEndian(uint64(requestBatchHeight))...)
}

// GetRequestKey returns the key for the request with the specified request ID
func GetRequestKey(requestID []byte) []byte {
	return append(requestKey, requestID...)
}

// GetRequestSubspaceByReqCtx returns the key for the requests of the specified request context
func GetRequestSubspaceByReqCtx(requestContextID []byte, batchCounter uint64) []byte {
	return append(append(requestKey, requestContextID...), sdk.Uint64ToBigEndian(batchCounter)...)
}

// GetActiveRequestKey returns the key for the active request with the specified request ID in the given height
func GetActiveRequestKey(serviceName string, provider sdk.AccAddress, expirationHeight int64, requestID []byte) []byte {
	activeRequest := append(append(append(getStringsKey([]string{serviceName, provider.String()}), emptyByte...), sdk.Uint64ToBigEndian(uint64(expirationHeight))...), requestID...)
	return append(activeRequestKey, activeRequest...)
}

// GetActiveRequestSubspace returns the key for the active requests for the specified provider
func GetActiveRequestSubspace(serviceName string, provider sdk.AccAddress) []byte {
	return append(append(activeRequestKey, getStringsKey([]string{serviceName, provider.String()})...), emptyByte...)
}

// GetActiveRequestKeyByID returns the key for the active request with the specified request ID
func GetActiveRequestKeyByID(requestID []byte) []byte {
	return append(activeRequestByIDKey, requestID...)
}

// GetActiveRequestSubspaceByReqCtx returns the key for the active requests for the specified request context
func GetActiveRequestSubspaceByReqCtx(requestContextID []byte, batchCounter uint64) []byte {
	return append(append(activeRequestByIDKey, requestContextID...), sdk.Uint64ToBigEndian(batchCounter)...)
}

// GetRequestVolumeKey returns the key for the request volume for the specified consumer and binding
func GetRequestVolumeKey(consumer sdk.AccAddress, serviceName string, provider sdk.AccAddress) []byte {
	return append(append(requestVolumeKey, getStringsKey([]string{consumer.String(), serviceName, provider.String()})...), emptyByte...)
}

// GetResponseKey returns the key for the response for the given request ID
func GetResponseKey(requestID []byte) []byte {
	return append(responseKey, requestID...)
}

// GetResponseSubspaceByReqCtx returns the key for responses for the specified request context and batch counter
func GetResponseSubspaceByReqCtx(requestContextID []byte, batchCounter uint64) []byte {
	return append(append(responseKey, requestContextID...), sdk.Uint64ToBigEndian(batchCounter)...)
}

// GetEarnedFeesKey returns the key for the earned fees of the specified provider
func GetEarnedFeesKey(provider sdk.AccAddress) []byte {
	return append(earnedFeesKey, provider.Bytes()...)
}

func getStringsKey(ss []string) (result []byte) {
	for _, s := range ss {
		result = append(append(result, []byte(s)...), emptyByte...)
	}

	if len(result) > 0 {
		return result[0 : len(result)-1]
	}

	return
}
