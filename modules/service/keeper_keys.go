package service

import (
	sdk "github.com/irisnet/irishub/types"
	"encoding/binary"
)

var (
	// the separator for string key
	emptyByte = []byte{0x00}

	// Keys for store prefixes
	serviceDefinitionKey         = []byte{0x01}
	methodPropertyKey            = []byte{0x02}
	bindingPropertyKey           = []byte{0x03}
	requestKey                   = []byte{0x05}
	responseKey                  = []byte{0x06}
	requestsByExpirationIndexKey = []byte{0x07}
	intraTxCounterKey            = []byte{0x08} // key for intra-block tx index
	activeRequestKey             = []byte{0x09} // key for active request
	returnedFeeKey               = []byte{0x10}
	incomingFeeKey               = []byte{0x11}

	serviceFeeTaxKey     = []byte{0x12}
	serviceFeeTaxPoolKey = []byte{0x13}
)

func GetServiceDefinitionKey(chainId, name string) []byte {
	return append(serviceDefinitionKey, getStringsKey([]string{chainId, name})...)
}

// id can not be zero
func GetMethodPropertyKey(chainId, serviceName string, id int16) []byte {
	return append(methodPropertyKey, getStringsKey([]string{chainId, serviceName, string(id)})...)
}

// Key for getting all methods on a service from the store
func GetMethodsSubspaceKey(chainId, serviceName string) []byte {
	return append(append(methodPropertyKey, getStringsKey([]string{chainId, serviceName})...), emptyByte...)
}

func GetServiceBindingKey(defChainId, name, bindChainId string, provider sdk.AccAddress) []byte {
	return append(bindingPropertyKey, getStringsKey([]string{defChainId, name, bindChainId, provider.String()})...)
}

// Key for getting all methods on a service from the store
func GetBindingsSubspaceKey(chainId, serviceName string) []byte {
	return append(append(bindingPropertyKey, getStringsKey([]string{chainId, serviceName})...), emptyByte...)
}

func GetRequestKey(defChainId, serviceName, bindChainId string, provider sdk.AccAddress, height int64, counter int16) []byte {
	return append(requestKey, getStringsKey([]string{defChainId, serviceName,
		bindChainId, provider.String(), string(height), string(counter)})...)
}

func GetActiveRequestKey(defChainId, serviceName, bindChainId string, provider sdk.AccAddress, height int64, counter int16) []byte {
	return append(activeRequestKey, getStringsKey([]string{defChainId, serviceName,
		bindChainId, provider.String(), string(height), string(counter)})...)
}

func GetSubActiveRequestKey(defChainId, serviceName, bindChainId string, provider sdk.AccAddress) []byte {
	return append(append(
		activeRequestKey, getStringsKey([]string{defChainId, serviceName,
			bindChainId, provider.String()})...),
		emptyByte...)
}

func GetResponseKey(reqChainId string, eHeight, rHeight int64, counter int16) []byte {
	return append(responseKey, getStringsKey([]string{reqChainId,
		string(eHeight), string(rHeight), string(counter)})...)
}

// get the expiration index of a request
func GetRequestsByExpirationIndexKeyByReq(req SvcRequest) []byte {
	return GetRequestsByExpirationIndexKey(req.ExpirationHeight, req.RequestHeight, req.RequestIntraTxCounter)
}

func GetRequestsByExpirationIndexKey(eHeight, rHeight int64, counter int16) []byte {
	// key is of format prefix(1) || expirationHeight(8) || requestHeight(8) || counterBytes(2)
	key := make([]byte, 1+8+8+2)
	key[0] = requestsByExpirationIndexKey[0]
	binary.BigEndian.PutUint64(key[1:9], uint64(eHeight))
	binary.BigEndian.PutUint64(key[9:17], uint64(rHeight))
	binary.BigEndian.PutUint16(key[17:19], uint16(counter))
	return key
}

// get the expiration prefix for all request of a block height
func GetRequestsByExpirationPrefix(height int64) []byte {
	// key is of format prefix || expirationHeight
	key := make([]byte, 1+8)
	key[0] = requestsByExpirationIndexKey[0]
	binary.BigEndian.PutUint64(key[1:9], uint64(height))
	return key
}

func GetReturnedFeeKey(address sdk.AccAddress) []byte {
	return append(returnedFeeKey, address.Bytes()...)
}

func GetIncomingFeeKey(address sdk.AccAddress) []byte {
	return append(incomingFeeKey, address.Bytes()...)
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
