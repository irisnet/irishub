package iservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	// the separator for string key
	emptyByte = []byte{0x00}

	// Keys for store prefixes
	serviceDefinitionKey = []byte{0x01}
	methodPropertyKey    = []byte{0x02}
	bindingPropertyKey   = []byte{0x03}
)

func GetServiceDefinitionKey(chainId, name string) []byte {
	return append(append(append(
		serviceDefinitionKey,
		[]byte(chainId)...),
		emptyByte...),
		[]byte(name)...)
}

// id can not zero
func GetMethodPropertyKey(chainId, serviceName string, id int) []byte {
	return append(append(append(append(append(
		methodPropertyKey,
		[]byte(chainId)...),
		emptyByte...),
		[]byte(serviceName)...),
		emptyByte...),
		[]byte(string(id))...)
}

// Key for getting all methods on a service from the store
func GetMethodsSubspaceKey(chainId, serviceName string) []byte {
	return append(append(append(append(
		methodPropertyKey,
		[]byte(chainId)...),
		emptyByte...),
		[]byte(serviceName)...),
		emptyByte...)
}

func GetServiceBindingKey(defChainId, name, bindChainId string, provider sdk.AccAddress) []byte {
	return append(append(append(append(append(append(append(
		bindingPropertyKey,
		[]byte(defChainId)...),
		emptyByte...),
		[]byte(name)...),
		emptyByte...),
		[]byte(bindChainId)...),
		emptyByte...),
		[]byte(provider.String())...)
}
