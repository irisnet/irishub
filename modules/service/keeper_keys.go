package service

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	incomingFeeKey               = []byte{0x10}
)

func GetServiceDefinitionKey(chainId, name string) []byte {
	return append(serviceDefinitionKey, getStringsKey([]string{chainId, name})...)
}

// id can not zero
func GetMethodPropertyKey(chainId, serviceName string, id int16) []byte {
	return append(methodPropertyKey, getStringsKey([]string{chainId, serviceName, string(id)})...)
}

// Key for getting all methods on a service from the store
func GetMethodsSubspaceKey(chainId, serviceName string) []byte {
	return append(methodPropertyKey, getStringsKey([]string{chainId, serviceName})...)
}

func GetServiceBindingKey(defChainId, name, bindChainId string, provider sdk.AccAddress) []byte {
	return append(bindingPropertyKey, getStringsKey([]string{defChainId, name, bindChainId, provider.String()})...)
}

// Key for getting all methods on a service from the store
func GetBindingsSubspaceKey(chainId, serviceName string) []byte {
	return append(bindingPropertyKey, getStringsKey([]string{chainId, serviceName})...)
}

func GetRequestKey(defChainId, serviceName, bindChainId string, provider sdk.AccAddress, height int64, counter int16) []byte {
	return append(requestKey, getStringsKey([]string{defChainId, serviceName,
		bindChainId, provider.String(), string(height), string(counter)})...)
}

func GetActiveRequestKey(defChainId, serviceName, bindChainId string, provider sdk.AccAddress, height int64, counter int16) []byte {
	return append(activeRequestKey, getStringsKey([]string{defChainId, serviceName,
		bindChainId, provider.String(), string(height), string(counter)})...)
}

func GetResponseKey(reqChainId string, consumer sdk.AccAddress, height int64, counter int16) []byte {
	return append(responseKey, getStringsKey([]string{reqChainId, consumer.String(),
		string(height), string(counter)})...)
}

// get the expiration index of a request
func GetRequstsByExpirationIndexKey(height int64, counter int16) []byte {
	// key is of format prefix || expirationHeight || counterBytes
	key := make([]byte, 1+8+2)
	key[0] = requestsByExpirationIndexKey[0]
	binary.BigEndian.PutUint64(key[1:9], uint64(height))
	binary.BigEndian.PutUint16(key[9:11], uint16(counter))
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
