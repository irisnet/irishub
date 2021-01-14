package types

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/xeipuuv/gojsonschema"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const PATH_BODY = "body"

// ServiceSchemas defines the service schemas
type ServiceSchemas struct {
	Input  map[string]interface{} `json:"input"`
	Output map[string]interface{} `json:"output"`
}

// ValidateServiceSchemas validates the given service schemas
func ValidateServiceSchemas(schemas string) error {
	if len(schemas) == 0 {
		return sdkerrors.Wrap(ErrInvalidSchemas, "schemas missing")
	}

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

	return nil
}

// ValidateBindingPricing validates the given pricing against the Pricing JSON Schema
func ValidateBindingPricing(pricing string) error {
	if len(pricing) == 0 {
		return sdkerrors.Wrap(ErrInvalidPricing, "pricing missing")
	}

	if err := validateDocument([]byte(PricingSchema), pricing); err != nil {
		return sdkerrors.Wrap(ErrInvalidPricing, err.Error())
	}

	return nil
}

// ValidateRequestInput validates the request input against the input schema
func ValidateRequestInput(input string) error {
	if err := validateDocument([]byte(InputSchema), input); err != nil {
		return sdkerrors.Wrap(ErrInvalidRequestInput, err.Error())
	}
	return nil
}

func ValidateRequestInputBody(schemas string, inputBody string) error {
	inputSchemaBz, err := parseInputSchema(schemas)
	if err != nil {
		return err
	}

	if err := validateDocument(inputSchemaBz, inputBody); err != nil {
		return sdkerrors.Wrap(ErrInvalidRequestInputBody, err.Error())
	}

	return nil
}

// ValidateResponseResult validates the response result against the result schema
func ValidateResponseResult(result string) error {
	if len(result) == 0 {
		return sdkerrors.Wrap(ErrInvalidResponseResult, "result missing")
	}

	if err := validateDocument([]byte(ResultSchema), result); err != nil {
		return sdkerrors.Wrap(ErrInvalidResponseResult, err.Error())
	}

	return nil
}

// ValidateResponseOutput validates the response output against the output schema
func ValidateResponseOutput(output string) error {
	if err := validateDocument([]byte(OutputSchema), output); err != nil {
		return sdkerrors.Wrap(ErrInvalidResponseOutput, err.Error())
	}
	return nil
}

func ValidateResponseOutputBody(schemas string, outputBody string) error {
	outputSchemaBz, err := parseOutputSchema(schemas)
	if err != nil {
		return err
	}

	if err := validateDocument(outputSchemaBz, outputBody); err != nil {
		return sdkerrors.Wrap(ErrInvalidResponseOutputBody, err.Error())
	}

	return nil
}

func validateInputSchema(inputSchema map[string]interface{}) error {
	inputSchemaBz, err := json.Marshal(inputSchema)
	if err != nil {
		return sdkerrors.Wrap(ErrInvalidSchemas, fmt.Sprintf("failed to marshal the input schema: %s", err))
	}

	if _, err = gojsonschema.NewSchema(gojsonschema.NewBytesLoader(inputSchemaBz)); err != nil {
		return sdkerrors.Wrap(ErrInvalidSchemas, fmt.Sprintf("invalid input schema: %s", err))
	}

	return nil
}

func validateOutputSchema(outputSchema map[string]interface{}) error {
	outputSchemaBz, err := json.Marshal(outputSchema)
	if err != nil {
		return sdkerrors.Wrap(ErrInvalidSchemas, fmt.Sprintf("failed to marshal the output schema: %s", err))
	}

	if _, err = gojsonschema.NewSchema(gojsonschema.NewBytesLoader(outputSchemaBz)); err != nil {
		return sdkerrors.Wrap(ErrInvalidSchemas, fmt.Sprintf("invalid output schema: %s", err))
	}

	return nil
}

// parseServiceSchemas parses the given schemas to ServiceSchemas
func parseServiceSchemas(schemas string) (ServiceSchemas, error) {
	var svcSchemas ServiceSchemas
	if err := json.Unmarshal([]byte(schemas), &svcSchemas); err != nil {
		return svcSchemas, sdkerrors.Wrap(ErrInvalidSchemas, fmt.Sprintf("failed to unmarshal the schemas: %s", err))
	}

	return svcSchemas, nil
}

// parseInputSchema parses the input schema from the given schemas
func parseInputSchema(schemas string) ([]byte, error) {
	svcSchemas, err := parseServiceSchemas(schemas)
	if err != nil {
		return nil, err
	}

	inputSchemaBz, err := json.Marshal(svcSchemas.Input)
	if err != nil {
		return nil, sdkerrors.Wrap(ErrInvalidSchemas, fmt.Sprintf("failed to marshal the input schema: %s", err))
	}

	return inputSchemaBz, nil
}

// parseOutputSchema parses the output schema from the given schemas
func parseOutputSchema(schemas string) ([]byte, error) {
	svcSchemas, err := parseServiceSchemas(schemas)
	if err != nil {
		return nil, err
	}

	outputSchemaBz, err := json.Marshal(svcSchemas.Output)
	if err != nil {
		return nil, sdkerrors.Wrap(ErrInvalidSchemas, fmt.Sprintf("failed to marshal the output schema: %s", err))
	}

	return outputSchemaBz, nil
}

// validateDocument wraps the gojsonschema validation
func validateDocument(schema []byte, document string) error {
	schemaLoader := gojsonschema.NewBytesLoader(schema)
	docLoader := gojsonschema.NewStringLoader(document)

	res, err := gojsonschema.Validate(schemaLoader, docLoader)
	if err != nil {
		return err
	}

	if !res.Valid() {
		for _, e := range res.Errors() {
			return errors.New(e.String()) // only return the first error
		}
	}

	return nil
}

const (
	// PricingSchema is the Pricing JSON Schema
	PricingSchema = `
	{
		"$schema": "http://json-schema.org/draft-04/schema#",
		"title": "iservice-pricing",
		"description": "Interchain Service Pricing Schema",
		"type": "object",
		"definitions": {
			"discount": {
				"description": "promotion discount, greater than 0 and less than 1",
				"type": "string",
				"pattern": "^0\\.\\d*[1-9]$"
			},
			"promotion_by_time": {
				"description": "promotion by time",
				"type": "object",
				"properties": {
					"start_time": {
						"description": "starting time of the promotion",
						"type": "string",
						"format": "date-time"
					},
					"end_time": {
						"description": "ending time of the promotion",
						"type": "string",
						"format": "date-time"
					},
					"discount": {
						"$ref": "#/definitions/discount"
					}
				},
				"additionalProperties": false,
				"required": [
					"start_time",
					"end_time",
					"discount"
				]
			},
			"promotion_by_volume": {
				"description": "promotion by volume",
				"type": "object",
				"properties": {
					"volume": {
						"description": "minimal volume for the promotion",
						"type": "integer",
						"minimum": 1
					},
					"discount": {
						"$ref": "#/definitions/discount"
					}
				},
				"additionalProperties": false,
				"required": [
					"volume",
					"discount"
				]
			}
		},
		"properties": {
			"price": {
				"description": "base price in min unit, e.g. 5uiris",
				"type": "string",
				"pattern": "^\\d+[a-z][a-z0-9/]{2,127}$"
			},
			"promotions_by_time": {
				"description": "promotions by time, in ascending order",
				"type": "array",
				"items": {
					"$ref": "#/definitions/promotion_by_time"
				},
				"maxItems": 5,
				"uniqueItems": true
			},
			"promotions_by_volume": {
				"description": "promotions by volume, in ascending order",
				"type": "array",
				"items": {
					"$ref": "#/definitions/promotion_by_volume"
				},
				"maxItems": 5,
				"uniqueItems": true
			}
		},
		"additionalProperties": false,
		"required": [
			"price"
		]
	}`

	// ResultSchema is the JSON Schema for the response result
	ResultSchema = `
	{
		"$schema": "http://json-schema.org/draft-04/schema#",
		"title": "iservice-result",
		"description": "Interchain Service Result Schema",
		"type": "object",
		"properties": {
			"code": {
				"description": "result code",
				"type": "integer",
				"enum": [
					200,
					400,
					500
				]
			},
			"message": {
				"description": "result message",
				"type": "string"
			}
		},
		"additionalProperties": false,
		"required": [
			"code",
			"message"
		]
	}`

	// InputSchema is the JSON Schema for the request input
	InputSchema = `
	{
		"$schema": "http://json-schema.org/draft-04/schema#",
		"title": "iservice-input",
		"description": "Interchain Service Input schema",
		"type": "object",
		"properties": {
			"header": {
				"description": "header",
				"type": "object"
			},
			"body": {
				"description": "body",
				"type": "object"
			}
		},
		"required": [
			"header"
		]
	}`

	// OutputSchema is the JSON Schema for the response output
	OutputSchema = `
	{
		"$schema": "http://json-schema.org/draft-04/schema#",
		"title": "iservice-output",
		"description": "Interchain Service Output Schema",
		"type": "object",
		"properties": {
			"header": {
				"description": "header",
				"type": "object"
			},
			"body": {
				"description": "body",
				"type": "object"
			}
		},
		"required": [
			"header"
		]
	}`
)
