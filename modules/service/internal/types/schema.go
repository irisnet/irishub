package types

import (
	"encoding/json"

	"github.com/xeipuuv/gojsonschema"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ServiceSchemas defines the service schemas
type ServiceSchemas struct {
	Input  map[string]interface{} `json:"input"`
	Output map[string]interface{} `json:"output"`
	Error  map[string]interface{} `json:"error"`
}

// ValidateServiceSchemas validates the given service schemas
func ValidateServiceSchemas(schemas string) error {
	var svcSchemas ServiceSchemas
	if err := json.Unmarshal([]byte(schemas), &svcSchemas); err != nil {
		return sdkerrors.Wrapf(ErrInvalidSchemas, "failed to unmarshal the schemas: %s", err)
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
func ValidateBindingPricing(pricing string) error {
	// TODO
	return nil
}

// ValidateRequestInput validates the request input against the input schema.
// assume that the schemas is valid.
func ValidateRequestInput(schemas string, input string) error {
	inputSchemaBz := parseInputSchema(schemas)
	inputSchemaLoader := gojsonschema.NewBytesLoader(inputSchemaBz)

	inputLoader := gojsonschema.NewStringLoader(input)

	if _, err := gojsonschema.Validate(inputSchemaLoader, inputLoader); err != nil {
		return err
	}

	return nil
}

// ValidateResponseOutput validates the response output against the output schema.
// assume that the schemas is valid.
func ValidateResponseOutput(schemas string, output string) error {
	outputSchemaBz := parseOutputSchema(schemas)
	outputSchemaLoader := gojsonschema.NewBytesLoader(outputSchemaBz)

	outputLoader := gojsonschema.NewStringLoader(output)

	if _, err := gojsonschema.Validate(outputSchemaLoader, outputLoader); err != nil {
		return err
	}

	return nil
}

// ValidateResponseError validates the response err against the error schema.
// assume that the schemas is valid.
func ValidateResponseError(schemas string, err string) error {
	errSchemaBz := parseErrorSchema(schemas)
	errSchemaLoader := gojsonschema.NewBytesLoader(errSchemaBz)

	errLoader := gojsonschema.NewStringLoader(err)

	if _, err := gojsonschema.Validate(errSchemaLoader, errLoader); err != nil {
		return err
	}

	return nil
}

func validateInputSchema(inputSchema map[string]interface{}) error {
	inputSchemaBz, err := json.Marshal(inputSchema)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidSchemas, "invalid input schema: %s", err)
	}

	_, err = gojsonschema.NewSchema(gojsonschema.NewBytesLoader(inputSchemaBz))
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidSchemas, "invalid input schema: %s", err)
	}

	return nil
}

func validateOutputSchema(outputSchema map[string]interface{}) error {
	outputSchemaBz, err := json.Marshal(outputSchema)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidSchemas, "invalid output schema: %s", err)
	}

	_, err = gojsonschema.NewSchema(gojsonschema.NewBytesLoader(outputSchemaBz))
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidSchemas, "invalid output schema: %s", err)
	}

	return nil
}

func validateErrorSchema(errSchema map[string]interface{}) error {
	errSchemaBz, err := json.Marshal(errSchema)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidSchemas, "invalid error schema: %s", err)
	}

	_, err = gojsonschema.NewSchema(gojsonschema.NewBytesLoader(errSchemaBz))
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidSchemas, "invalid error schema: %s", err)
	}

	return nil
}

// parseInputSchema parses the input schema from the given schemas
// assume that the schemas is valid. Panic if invalid
func parseInputSchema(schemas string) []byte {
	var svcSchemas ServiceSchemas
	if err := json.Unmarshal([]byte(schemas), &svcSchemas); err != nil {
		panic(err)
	}

	inputSchemaBz, err := json.Marshal(svcSchemas.Input)
	if err != nil {
		panic(err)
	}

	return inputSchemaBz
}

// parseOutputSchema parses the output schema from the given schemas
// assume that the schemas is valid. Panic if invalid
func parseOutputSchema(schemas string) []byte {
	var svcSchemas ServiceSchemas
	if err := json.Unmarshal([]byte(schemas), &svcSchemas); err != nil {
		panic(err)
	}

	outputSchemaBz, err := json.Marshal(svcSchemas.Output)
	if err != nil {
		panic(err)
	}

	return outputSchemaBz
}

// parseErrorSchema parses the error schema from the given schemas
// assume that the schemas is valid. Panic if invalid
func parseErrorSchema(schemas string) []byte {
	var svcSchemas ServiceSchemas
	if err := json.Unmarshal([]byte(schemas), &svcSchemas); err != nil {
		panic(err)
	}

	errSchemaBz, err := json.Marshal(svcSchemas.Error)
	if err != nil {
		panic(err)
	}

	return errSchemaBz
}
