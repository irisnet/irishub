package iservice

var (
	// Keys for store prefixes
	serviceDefinitionKey = []byte{0x00}
	methodPropertyKey    = []byte{0x01}
)

func GetServiceDefinitionKey(chainId, name string) []byte {
	return append(append(
		serviceDefinitionKey,
		[]byte(chainId)...),
		[]byte(name)...)
}

func GetMethodPropertyKey(chainId, serviceName, methodName string) []byte {
	return append(append(append(
		methodPropertyKey,
		[]byte(chainId)...),
		[]byte(serviceName)...),
		[]byte(methodName)...)
}
