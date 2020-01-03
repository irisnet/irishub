package types

// ServiceSchema defines the service schema fields which are required
type ServiceSchema struct {
	Input  map[string]interface{}
	Output map[string]interface{}
	Error  map[string]interface{}
}

// ValidateServiceSchema checks if the given schema is valid
func ValidateServiceSchema(schema string) error {
	return nil
}

// ParseServiceSchema parses the given schema
func ParseServiceSchema(schema string) (ServiceSchema, error) {
	// TODO
	return ServiceSchema{}, nil
}

// ValidateBindingPricing checks if the given pricing is valid
func ValidateBindingPricing(pricing string) error {
	// TODO
	return nil
}

// ValidateRequestInput validates the request input against the given schema
func ValidateRequestInput(schema string, input string) error {
	// TODO
	return nil
}

// ValidateResponseOutput validates the response output against the given schema
func ValidateResponseOutput(schema string, output string) error {
	// TODO
	return nil
}

// ValidateResponseError validates the response err against the given schema
func ValidateResponseError(schema string, err string) error {
	// TODO
	return nil
}
