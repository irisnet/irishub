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
	requestContextKey      = []byte{0x03}
	expiredRequestBatchKey = []byte{0x04}
	newRequestBatchKey     = []byte{0x05}
	requestKey             = []byte{0x06}
	activeRequestKey       = []byte{0x07}
	responseKey            = []byte{0x08}
	returnedFeeKey         = []byte{0x09}
	incomingFeeKey         = []byte{0x10}
)

const (
	// tx counter key in the context
	contextKeyIntraTxCounter = 0
)

// GetServiceDefinitionKey returns the key for the service definition with the specified name
func GetServiceDefinitionKey(name string) []byte {
	return append(serviceDefinitionKey, []byte(name)...)
}

// GetServiceBindingKey returns the key for the service binding with the specified name and provider
func GetServiceBindingKey(serviceName string, provider sdk.AccAddress) []byte {
	return append(serviceBindingKey, getStringsKey([]string{serviceName, provider.String()})...)
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
	reqBatchExpiration := append(append(sdk.Uint64ToBigEndian(uint64(batchExpirationHeight)), emptyByte...), requestContextID...)
	return append(expiredRequestBatchKey, reqBatchExpiration...)
}

// GetNewRequestBatchKey returns the key for the new batch request of the specified request context in the given height
func GetNewRequestBatchKey(requestContextID []byte, batchRequestHeight int64) []byte {
	newBatchRequest := append(append(sdk.Uint64ToBigEndian(uint64(batchRequestHeight)), emptyByte...), requestContextID...)
	return append(expiredRequestBatchKey, newBatchRequest...)
}

// GetRequestKey returns the key for the request with the specified request ID
func GetRequestKey(requestID []byte) []byte {
	return append(requestKey, requestID...)
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

// GetResponseKey returns the key for the response for the given request ID
func GetResponseKey(requestID []byte) []byte {
	return append(responseKey, requestID...)
}

func GetReturnedFeeKey(address sdk.AccAddress) []byte {
	return append(returnedFeeKey, address.Bytes()...)
}

func GetIncomingFeeKey(address sdk.AccAddress) []byte {
	return append(incomingFeeKey, address.Bytes()...)
}

func GetIntraTxCounterKey() int {
	return contextKeyIntraTxCounter
}

func getStringsKey(ss []string) (result []byte) {
	for _, s := range ss {
		result = append(append(
			result,
			[]byte(s)...),
			emptyByte...)
	}
	if len(result) > 0 {
		return result[0 : len(result)-1]
	}
	return
}
