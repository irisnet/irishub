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

// ValidateBindingPricing validates the given pricing against the Pricing JSON Schema
func ValidateBindingPricing(pricing string) sdk.Error {
	if !validDocument([]byte(PricingSchema), pricing) {
		return ErrInvalidPricing(DefaultCodespace, "invalid pricing")
	}

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
		return ErrInvalidSchemas(DefaultCodespace, fmt.Sprintf("invalid input schema: %s", err))
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
		return ErrInvalidSchemas(DefaultCodespace, fmt.Sprintf("invalid output schema: %s", err))
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
		return ErrInvalidSchemas(DefaultCodespace, fmt.Sprintf("invalid error schema: %s", err))
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

// PricingSchema is the Pricing JSON Schema
const PricingSchema = `
{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "title": "Irishub Service Pricing",
    "description": "The Irishub Service Pricing specification",
    "type": "object",

    "definitions": {
		"coin": {
			"type": "object",
			"description": "price coin",

			"properties": {
				"denom": {
					"type": "string",
					"description": "the denomination of the coin",
					"minLength": 1
				},
				"amount": {
					"type": "string",
					"description": "the amount of the coin",
					"pattern": "^[1-9]\\d+$"
				}
			},
			
			"additionalProperties": false,
			"required": ["denom", "amount"]
		},
		"discount": {
			"type": "number",
			"description": "promotion discount",
			"minimum": 0,
			"exclusiveMinimum": true,
			"maximum": 1,
			"exclusiveMaximum": true
		},
        "promotion_by_time": {
            "type": "object",
            "description": "promotion activity by time",

            "properties": {
                "start_time": {
                    "type": "integer",
					"description": "starting time of the promotion",
					"minimum": 0,
					"exclusiveMinimum": true
				},
                "end_time": {
                    "type": "integer",
					"description": "ending time of the promotion",
					"minimum": 0,
					"exclusiveMinimum": true
				},
				"discount": {
                    "$ref": "#/definitions/discount"
				}
			},
			
			"additionalProperties": false,
			"required": ["start_time", "end_time", "discount"]
		},
		"promotion_by_volume": {
            "type": "object",
            "description": "promotion activity by volume",

            "properties": {
                "volume": {
                    "type": "integer",
					"description": "minimal volume for the promotion",
					"minimum": 1
				},
				"discount": {
                    "$ref": "#/definitions/discount"
				}
			},
			
			"additionalProperties": false,
			"required": ["volume", "discount"]
        }
    },

    "properties": {
        "price": {
            "type": "array",
			"description": "normal service price",

			"item": {
				"$ref": "#/definitions/coin"
			}
        },
        "promotions_by_time": {
            "type": "array",
            "description": "promotion activities by time",

            "item": {
                "$ref": "#/definitions/promotion_by_time"
            }
		},
		"promotions_by_volume": {
            "type": "array",
            "description": "promotion activities by volume",

            "item": {
                "$ref": "#/definitions/promotion_by_volume"
            }
		}
    },
	
	"additionalProperties": false,
    "required": ["price"]
}
`
