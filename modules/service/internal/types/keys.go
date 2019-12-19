package types

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the service module
	ModuleName = "service"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the service module
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the service module
	RouterKey = ModuleName

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
	// the separator for string key
	emptyByte = []byte{0x00}

	// Keys for store prefixes
	ServiceDefinitionKey         = []byte{0x01}
	MethodPropertyKey            = []byte{0x02}
	BindingPropertyKey           = []byte{0x03}
	RequestKey                   = []byte{0x05}
	ResponseKey                  = []byte{0x06}
	RequestsByExpirationIndexKey = []byte{0x07}
	IntraTxCounterKey            = []byte{0x08} // key for intra-block tx index
	ActiveRequestKey             = []byte{0x09} // key for active request
	ReturnedFeeKey               = []byte{0x10}
	IncomingFeeKey               = []byte{0x11}

	ServiceFeeTaxKey        = []byte{0x12}
	ServiceSlashFractionKey = []byte{0x13}
)

func GetServiceDefinitionKey(chainId, name string) []byte {
	return append(ServiceDefinitionKey, getStringsKey([]string{chainId, name})...)
}

// id can not be zero
func GetMethodPropertyKey(chainId, serviceName string, id int16) []byte {
	return append(MethodPropertyKey, getStringsKey([]string{chainId, serviceName, string(id)})...)
}

// Key for getting all methods on a service from the store
func GetMethodsSubspaceKey(chainId, serviceName string) []byte {
	return append(append(MethodPropertyKey, getStringsKey([]string{chainId, serviceName})...), emptyByte...)
}

func GetServiceBindingKey(defChainId, name, bindChainId string, provider sdk.AccAddress) []byte {
	return append(BindingPropertyKey, getStringsKey([]string{defChainId, name, bindChainId, provider.String()})...)
}

// Key for getting all methods on a service from the store
func GetBindingsSubspaceKey(chainId, serviceName string) []byte {
	return append(append(BindingPropertyKey, getStringsKey([]string{chainId, serviceName})...), emptyByte...)
}

func GetRequestKey(defChainId, serviceName, bindChainId string, provider sdk.AccAddress, height int64, counter int16) []byte {
	return append(RequestKey, getStringsKey([]string{defChainId, serviceName,
		bindChainId, provider.String(), string(height), string(counter)})...)
}

func GetActiveRequestKey(defChainId, serviceName, bindChainId string, provider sdk.AccAddress, height int64, counter int16) []byte {
	return append(ActiveRequestKey, getStringsKey([]string{defChainId, serviceName,
		bindChainId, provider.String(), string(height), string(counter)})...)
}

func GetSubActiveRequestKey(defChainId, serviceName, bindChainId string, provider sdk.AccAddress) []byte {
	return append(append(
		ActiveRequestKey, getStringsKey([]string{defChainId, serviceName,
			bindChainId, provider.String()})...),
		emptyByte...)
}

func GetResponseKey(reqChainId string, eHeight, rHeight int64, counter int16) []byte {
	return append(ResponseKey, getStringsKey([]string{reqChainId,
		string(eHeight), string(rHeight), string(counter)})...)
}

// get the expiration index of a request
func GetRequestsByExpirationIndexKeyByReq(req SvcRequest) []byte {
	return GetRequestsByExpirationIndexKey(req.ExpirationHeight, req.RequestHeight, req.RequestIntraTxCounter)
}

func GetRequestsByExpirationIndexKey(eHeight, rHeight int64, counter int16) []byte {
	// key is of format prefix(1) || expirationHeight(8) || requestHeight(8) || counterBytes(2)
	key := make([]byte, 1+8+8+2)
	key[0] = RequestsByExpirationIndexKey[0]
	binary.BigEndian.PutUint64(key[1:9], uint64(eHeight))
	binary.BigEndian.PutUint64(key[9:17], uint64(rHeight))
	binary.BigEndian.PutUint16(key[17:19], uint16(counter))
	return key
}

// get the expiration prefix for all request of a block height
func GetRequestsByExpirationPrefix(height int64) []byte {
	// key is of format prefix || expirationHeight
	key := make([]byte, 1+8)
	key[0] = RequestsByExpirationIndexKey[0]
	binary.BigEndian.PutUint64(key[1:9], uint64(height))
	return key
}

func GetReturnedFeeKey(address sdk.AccAddress) []byte {
	return append(ReturnedFeeKey, address.Bytes()...)
}

func GetIncomingFeeKey(address sdk.AccAddress) []byte {
	return append(IncomingFeeKey, address.Bytes()...)
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
