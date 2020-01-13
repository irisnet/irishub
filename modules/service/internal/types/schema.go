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
func ValidateRequestInput(schemas string, input string) error {
	inputSchemaBz, err := parseInputSchema(schemas)
	if err != nil {
		return err
	}

	if !validDocument(inputSchemaBz, input) {
		return sdkerrors.Wrap(ErrInvalidRequestInput, "invalid request input")
	}

	return nil
}

// ValidateResponseOutput validates the response output against the output schema
func ValidateResponseOutput(schemas string, output string) error {
	outputSchemaBz, err := parseOutputSchema(schemas)
	if err != nil {
		return err
	}

	if !validDocument(outputSchemaBz, output) {
		return sdkerrors.Wrap(ErrInvalidResponseOutput, "invalid response output")
	}

	return nil
}

// ValidateResponseError validates the response err against the error schema.
func ValidateResponseError(schemas string, errResp string) error {
	errSchemaBz, err := parseErrorSchema(schemas)
	if err != nil {
		return err
	}

	if !validDocument(errSchemaBz, errResp) {
		return sdkerrors.Wrap(ErrInvalidResponseErr, "invalid response err")
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
func parseInputSchema(schemas string) ([]byte, error) {
	var svcSchemas ServiceSchemas
	if err := json.Unmarshal([]byte(schemas), &svcSchemas); err != nil {
		return nil, err
	}

	inputSchemaBz, err := json.Marshal(svcSchemas.Input)
	if err != nil {
		return nil, err
	}

	return inputSchemaBz, nil
}

// parseOutputSchema parses the output schema from the given schemas
func parseOutputSchema(schemas string) ([]byte, error) {
	var svcSchemas ServiceSchemas
	if err := json.Unmarshal([]byte(schemas), &svcSchemas); err != nil {
		return nil, err
	}

	outputSchemaBz, err := json.Marshal(svcSchemas.Output)
	if err != nil {
		return nil, err
	}

	return outputSchemaBz, nil
}

// parseErrorSchema parses the error schema from the given schemas
func parseErrorSchema(schemas string) ([]byte, error) {
	var svcSchemas ServiceSchemas
	if err := json.Unmarshal([]byte(schemas), &svcSchemas); err != nil {
		return nil, err
	}

	errSchemaBz, err := json.Marshal(svcSchemas.Error)
	if err != nil {
		return nil, err
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
