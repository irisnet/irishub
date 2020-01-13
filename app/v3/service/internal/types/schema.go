package types

import (
	"encoding/json"
	"fmt"

	"github.com/xeipuuv/gojsonschema"

	sdk "github.com/irisnet/irishub/types"
)

// ServiceSchemas defines the service schemas
type ServiceSchemas struct {
	Input  map[string]interface{} `json:"input"`
	Output map[string]interface{} `json:"output"`
	Error  map[string]interface{} `json:"error"`
}

// ValidateServiceSchemas validates the given service schemas
func ValidateServiceSchemas(schemas string) sdk.Error {
	svcSchemas, err := parseServiceSchemas(schemas)
	if err != nil {
		return err
	}

	if err := validateInputSchema(svcSchemas.Input); err != nil {
		return err
	}

	if err := validateOutputSchema(svcSchemas.Output); err != nil {
		return err
	}

	if err := validateErrorSchema(svcSchemas.Error); err != nil {
		return err
	}

	return nil
}

// ValidateBindingPricing checks if the given pricing is valid
func ValidateBindingPricing(pricing string) sdk.Error {
	// TODO
	return nil
}

// ValidateRequestInput validates the request input against the input schema
func ValidateRequestInput(schemas string, input string) sdk.Error {
	inputSchemaBz, err := parseInputSchema(schemas)
	if err != nil {
		return err
	}

	if !validDocument(inputSchemaBz, input) {
		return ErrInvalidRequestInput(DefaultCodespace, "invalid request input")
	}

	return nil
}

// ValidateResponseOutput validates the response output against the output schema
func ValidateResponseOutput(schemas string, output string) sdk.Error {
	outputSchemaBz, err := parseOutputSchema(schemas)
	if err != nil {
		return err
	}

	if !validDocument(outputSchemaBz, output) {
		return ErrInvalidResponseOutput(DefaultCodespace, "invalid response output")
	}

	return nil
}

// ValidateResponseError validates the response err against the error schema
func ValidateResponseError(schemas string, errResp string) sdk.Error {
	errSchemaBz, err := parseErrorSchema(schemas)
	if err != nil {
		return err
	}

	if !validDocument(errSchemaBz, errResp) {
		return ErrInvalidResponseErr(DefaultCodespace, "invalid response err")
	}

	return nil
}

func validateInputSchema(inputSchema map[string]interface{}) sdk.Error {
	inputSchemaBz, err := json.Marshal(inputSchema)
	if err != nil {
		return ErrInvalidSchemas(DefaultCodespace, fmt.Sprintf("failed to marshal the input schema: %s", err))
	}

	_, err = gojsonschema.NewSchema(gojsonschema.NewBytesLoader(inputSchemaBz))
	if err != nil {
		return ErrInvalidSchemas(DefaultCodespace, fmt.Sprintf("invalid input schema: %s", err))
	}

	return nil
}

func validateOutputSchema(outputSchema map[string]interface{}) sdk.Error {
	outputSchemaBz, err := json.Marshal(outputSchema)
	if err != nil {
		return ErrInvalidSchemas(DefaultCodespace, fmt.Sprintf("failed to marshal the output schema: %s", err))
	}

	_, err = gojsonschema.NewSchema(gojsonschema.NewBytesLoader(outputSchemaBz))
	if err != nil {
		return ErrInvalidSchemas(DefaultCodespace, fmt.Sprintf("invalid output schema: %s", err))
	}

	return nil
}

func validateErrorSchema(errSchema map[string]interface{}) sdk.Error {
	errSchemaBz, err := json.Marshal(errSchema)
	if err != nil {
		return ErrInvalidSchemas(DefaultCodespace, fmt.Sprintf("failed to marshal the err schema: %s", err))
	}

	_, err = gojsonschema.NewSchema(gojsonschema.NewBytesLoader(errSchemaBz))
	if err != nil {
		return ErrInvalidSchemas(DefaultCodespace, fmt.Sprintf("invalid error schema: %s", err))
	}

	return nil
}

// parseServiceSchemas parses the given schemas to ServiceSchemas
func parseServiceSchemas(schemas string) (ServiceSchemas, sdk.Error) {
	var svcSchemas ServiceSchemas
	if err := json.Unmarshal([]byte(schemas), &svcSchemas); err != nil {
		return svcSchemas, ErrInvalidSchemas(DefaultCodespace, fmt.Sprintf("failed to unmarshal the schemas: %s", err))
	}

	return svcSchemas, nil
}

// parseInputSchema parses the input schema from the given schemas
func parseInputSchema(schemas string) ([]byte, sdk.Error) {
	svcSchemas, err := parseServiceSchemas(schemas)
	if err != nil {
		return nil, err
	}

	inputSchemaBz, err2 := json.Marshal(svcSchemas.Input)
	if err2 != nil {
		return nil, ErrInvalidSchemas(DefaultCodespace, fmt.Sprintf("failed to marshal the input schema: %s", err2))
	}

	return inputSchemaBz, nil
}

// parseOutputSchema parses the output schema from the given schemas
func parseOutputSchema(schemas string) ([]byte, sdk.Error) {
	svcSchemas, err := parseServiceSchemas(schemas)
	if err != nil {
		return nil, err
	}

	outputSchemaBz, err2 := json.Marshal(svcSchemas.Output)
	if err != nil {
		return nil, ErrInvalidSchemas(DefaultCodespace, fmt.Sprintf("failed to marshal the output schema: %s", err2))
	}

	return outputSchemaBz, nil
}

// parseErrorSchema parses the error schema from the given schemas
func parseErrorSchema(schemas string) ([]byte, sdk.Error) {
	svcSchemas, err := parseServiceSchemas(schemas)
	if err != nil {
		return nil, err
	}

	errSchemaBz, err2 := json.Marshal(svcSchemas.Error)
	if err != nil {
		return nil, ErrInvalidSchemas(DefaultCodespace, fmt.Sprintf("failed to marshal the err schema: %s", err2))
	}

	return errSchemaBz, nil
}

// validDocument wraps the gojsonschema validation
func validDocument(schema []byte, document string) bool {
	schemaLoader := gojsonschema.NewBytesLoader(schema)
	docLoader := gojsonschema.NewStringLoader(document)

	res, err := gojsonschema.Validate(schemaLoader, docLoader)
	if err != nil {
		return false
	}

	return res.Valid()
}
