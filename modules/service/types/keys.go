package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the service module
	ModuleName = "service"

	// StoreKey is the string store representation
	StoreKey string = ModuleName

	// QuerierRoute is the querier route for the service module
	QuerierRoute string = ModuleName

	// RouterKey is the msg router key for the service module
	RouterKey string = ModuleName

	// DefaultParamspace is the default name for parameter store
	DefaultParamspace = ModuleName

	// DepositAccName is the root string for the service deposit account address
	DepositAccName = "service_deposit_account"

	// RequestAccName is the root string for the service request account address
	RequestAccName = "service_request_account"

	// TaxAccName is the root string for the service tax account address
	TaxAccName = "service_tax_account"
)

var (
	// Separator for string key
	Delimiter = []byte{0x00}

	// Keys for store prefixes
	ServiceDefinitionKey         = []byte{0x01} // prefix for service definition
	ServiceBindingKey            = []byte{0x02} // prefix for service binding
	OwnerServiceBindingKey       = []byte{0x03} // prefix for owner service binding
	OwnerKey                     = []byte{0x04} // prefix for the owner of a provider
	OwnerProviderKey             = []byte{0x05} // prefix for the provider with an owner
	PricingKey                   = []byte{0x06} // prefix for pricing
	WithdrawAddrKey              = []byte{0x07} // prefix for withdrawal address
	RequestContextKey            = []byte{0x08} // prefix for request context
	ExpiredRequestBatchKey       = []byte{0x09} // prefix for expired request batch
	NewRequestBatchKey           = []byte{0x10} // prefix for new request batch
	ExpiredRequestBatchHeightKey = []byte{0x11} // prefix for expired request batch height
	NewRequestBatchHeightKey     = []byte{0x12} // prefix for new request batch height
	RequestKey                   = []byte{0x13} // prefix for request
	ActiveRequestKey             = []byte{0x14} // prefix for active request
	ActiveRequestByIDKey         = []byte{0x15} // prefix for active requests by ID
	ResponseKey                  = []byte{0x16} // prefix for response
	RequestVolumeKey             = []byte{0x17} // prefix for request volume
	EarnedFeesKey                = []byte{0x18} // prefix for provider earned fees
	OwnerEarnedFeesKey           = []byte{0x19} // prefix for owner earned fees
	InternalCounterKey           = []byte{0x20} // prefix for internal counter key
)

// GetServiceDefinitionKey gets the key for the service definition with the specified service name
// VALUE: service/ServiceDefinition
func GetServiceDefinitionKey(serviceName string) []byte {
	return append(ServiceDefinitionKey, []byte(serviceName)...)
}

// GetServiceBindingKey gets the key for the service binding with the specified service name and provider
// VALUE: service/ServiceBinding
func GetServiceBindingKey(serviceName string, provider sdk.AccAddress) []byte {
	return append(ServiceBindingKey, getStringsKey([]string{serviceName, provider.String()})...)
}

// GetOwnerServiceBindingKey gets the key for the service binding with the specified owner
// VALUE: []byte{}
func GetOwnerServiceBindingKey(owner sdk.AccAddress, serviceName string, provider sdk.AccAddress) []byte {
	return append(append(append(append(OwnerServiceBindingKey, owner.Bytes()...), []byte(serviceName)...), Delimiter...), provider.Bytes()...)
}

// GetOwnerKey gets the key for the specified provider
// VALUE: sdk.AccAddress
func GetOwnerKey(provider sdk.AccAddress) []byte {
	return append(OwnerKey, provider.Bytes()...)
}

// GetOwnerProviderKey gets the key for the specified owner and provider
// VALUE: []byte{}
func GetOwnerProviderKey(owner, provider sdk.AccAddress) []byte {
	return append(append(OwnerProviderKey, owner.Bytes()...), provider.Bytes()...)
}

// GetPricingKey gets the key for the pricing of the specified binding
// VALUE: service/Pricing
func GetPricingKey(serviceName string, provider sdk.AccAddress) []byte {
	return append(PricingKey, getStringsKey([]string{serviceName, provider.String()})...)
}

// GetWithdrawAddrKey gets the key for the withdrawal address of the specified provider
// VALUE: withdrawal address ([]byte)
func GetWithdrawAddrKey(provider sdk.AccAddress) []byte {
	return append(WithdrawAddrKey, provider.Bytes()...)
}

// GetBindingsSubspace gets the key prefix for iterating through all bindings of the specified service name
func GetBindingsSubspace(serviceName string) []byte {
	return append(append(ServiceBindingKey, []byte(serviceName)...), Delimiter...)
}

// GetOwnerBindingsSubspace gets the key prefix for iterating through all bindings of the specified service name and owner
func GetOwnerBindingsSubspace(owner sdk.AccAddress, serviceName string) []byte {
	return append(append(append(OwnerServiceBindingKey, owner.Bytes()...), []byte(serviceName)...), Delimiter...)
}

// GetOwnerProvidersSubspace gets the key prefix for iterating through providers of the specified owner
func GetOwnerProvidersSubspace(owner sdk.AccAddress) []byte {
	return append(OwnerProviderKey, owner.Bytes()...)
}

// GetRequestContextKey returns the key for the request context with the specified ID
func GetRequestContextKey(requestContextID []byte) []byte {
	return append(RequestContextKey, requestContextID...)
}

// GetExpiredRequestBatchKey returns the key for the request batch expiration of the specified request context
func GetExpiredRequestBatchKey(requestContextID []byte, batchExpirationHeight int64) []byte {
	reqBatchExpiration := append(sdk.Uint64ToBigEndian(uint64(batchExpirationHeight)), requestContextID...)
	return append(ExpiredRequestBatchKey, reqBatchExpiration...)
}

// GetNewRequestBatchKey returns the key for the new batch request of the specified request context in the given height
func GetNewRequestBatchKey(requestContextID []byte, requestBatchHeight int64) []byte {
	newBatchRequest := append(sdk.Uint64ToBigEndian(uint64(requestBatchHeight)), requestContextID...)
	return append(NewRequestBatchKey, newBatchRequest...)
}

// GetExpiredRequestBatchSubspace returns the key prefix for iterating through the expired request batch queue in the specified height
func GetExpiredRequestBatchSubspace(batchExpirationHeight int64) []byte {
	return append(ExpiredRequestBatchKey, sdk.Uint64ToBigEndian(uint64(batchExpirationHeight))...)
}

// GetNewRequestBatchSubspace returns the key prefix for iterating through the new request batch queue in the specified height
func GetNewRequestBatchSubspace(requestBatchHeight int64) []byte {
	return append(NewRequestBatchKey, sdk.Uint64ToBigEndian(uint64(requestBatchHeight))...)
}

// GetExpiredRequestBatchHeightKey returns the key for the current request batch expiration height of the specified request context
func GetExpiredRequestBatchHeightKey(requestContextID []byte) []byte {
	return append(ExpiredRequestBatchHeightKey, requestContextID...)
}

// GetNewRequestBatchHeightKey returns the key for the new request batch height of the specified request context
func GetNewRequestBatchHeightKey(requestContextID []byte) []byte {
	return append(NewRequestBatchHeightKey, requestContextID...)
}

// GetRequestKey returns the key for the request with the specified request ID
func GetRequestKey(requestID []byte) []byte {
	return append(RequestKey, requestID...)
}

// GetRequestSubspaceByReqCtx returns the key prefix for iterating through the requests of the specified request context
func GetRequestSubspaceByReqCtx(requestContextID []byte, batchCounter uint64) []byte {
	return append(append(RequestKey, requestContextID...), sdk.Uint64ToBigEndian(batchCounter)...)
}

// GetActiveRequestKey returns the key for the active request with the specified request ID in the given height
func GetActiveRequestKey(serviceName string, provider sdk.AccAddress, expirationHeight int64, requestID []byte) []byte {
	activeRequest := append(append(append(getStringsKey([]string{serviceName, provider.String()}), Delimiter...), sdk.Uint64ToBigEndian(uint64(expirationHeight))...), requestID...)
	return append(ActiveRequestKey, activeRequest...)
}

// GetActiveRequestSubspace returns the key prefix for iterating through the active requests for the specified provider
func GetActiveRequestSubspace(serviceName string, provider sdk.AccAddress) []byte {
	return append(append(ActiveRequestKey, getStringsKey([]string{serviceName, provider.String()})...), Delimiter...)
}

// GetActiveRequestKeyByID returns the key for the active request with the specified request ID
func GetActiveRequestKeyByID(requestID []byte) []byte {
	return append(ActiveRequestByIDKey, requestID...)
}

// GetActiveRequestSubspaceByReqCtx returns the key prefix for iterating through the active requests for the specified request context
func GetActiveRequestSubspaceByReqCtx(requestContextID []byte, batchCounter uint64) []byte {
	return append(append(ActiveRequestByIDKey, requestContextID...), sdk.Uint64ToBigEndian(batchCounter)...)
}

// GetRequestVolumeKey returns the key for the request volume for the specified consumer and binding
func GetRequestVolumeKey(consumer sdk.AccAddress, serviceName string, provider sdk.AccAddress) []byte {
	return append(append(RequestVolumeKey, getStringsKey([]string{consumer.String(), serviceName, provider.String()})...), Delimiter...)
}

// GetResponseKey returns the key for the response for the given request ID
func GetResponseKey(requestID []byte) []byte {
	return append(ResponseKey, requestID...)
}

// GetResponseSubspaceByReqCtx returns the key prefix for iterating throuth responses for the specified request context and batch counter
func GetResponseSubspaceByReqCtx(requestContextID []byte, batchCounter uint64) []byte {
	return append(append(ResponseKey, requestContextID...), sdk.Uint64ToBigEndian(batchCounter)...)
}

// GetEarnedFeesKey gets the key for the earned fees of the specified provider and denom
func GetEarnedFeesKey(provider sdk.AccAddress, denom string) []byte {
	return append(append(EarnedFeesKey, provider.Bytes()...), []byte(denom)...)
}

// GetEarnedFeesSubspace gets the subspace for the earned fees of the specified provider
func GetEarnedFeesSubspace(provider sdk.AccAddress) []byte {
	return append(EarnedFeesKey, provider.Bytes()...)
}

// GetOwnerEarnedFeesKey returns the key for the earned fees of the specified owner and denom
func GetOwnerEarnedFeesKey(owner sdk.AccAddress, denom string) []byte {
	return append(OwnerEarnedFeesKey, owner.Bytes()...)
}

// GetEarnedFeesSubspace gets the subspace for the earned fees of the specified provider
func GetOwnerEarnedFeesSubspace(owner sdk.AccAddress) []byte {
	return append(OwnerEarnedFeesKey, owner.Bytes()...)
}

func getStringsKey(ss []string) (result []byte) {
	for _, s := range ss {
		result = append(append(result, []byte(s)...), Delimiter...)
	}

	if len(result) > 0 {
		return result[0 : len(result)-1]
	}

	return
}
