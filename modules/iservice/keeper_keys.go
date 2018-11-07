package iservice

var (
	// the separator for string key
	emptyByte = []byte{0x00}

	// Keys for store prefixes
	serviceDefinitionKey = []byte{0x01}
	methodPropertyKey    = []byte{0x02}
)

func GetServiceDefinitionKey(chainId, name string) []byte {
	return append(append(append(
		serviceDefinitionKey,
		[]byte(chainId)...),
		emptyByte...),
		[]byte(name)...)
}

func GetMethodPropertyKey(chainId, serviceName, methodName string) []byte {
	return append(append(append(append(append(
		methodPropertyKey,
		[]byte(chainId)...),
		emptyByte...),
		[]byte(serviceName)...),
		emptyByte...),
		[]byte(methodName)...)
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
