package types

import (
	"encoding/json"
	"errors"
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
	if err := validateDocument([]byte(PricingSchema), pricing); err != nil {
		return ErrInvalidPricing(DefaultCodespace, err.Error())
	}

	return nil
}

// ValidateRequestInput validates the request input against the input schema
func ValidateRequestInput(schemas string, input string) sdk.Error {
	inputSchemaBz, err := parseInputSchema(schemas)
	if err != nil {
		return err
	}

	if err := validateDocument(inputSchemaBz, input); err != nil {
		return ErrInvalidRequestInput(DefaultCodespace, err.Error())
	}

	return nil
}

// ValidateResponseOutput validates the response output against the output schema
func ValidateResponseOutput(schemas string, output string) sdk.Error {
	outputSchemaBz, err := parseOutputSchema(schemas)
	if err != nil {
		return err
	}

	if err := validateDocument(outputSchemaBz, output); err != nil {
		return ErrInvalidResponseOutput(DefaultCodespace, err.Error())
	}

	return nil
}

// ValidateResponseError validates the response err against the error schema
func ValidateResponseError(schemas string, errResp string) sdk.Error {
	errSchemaBz, err := parseErrorSchema(schemas)
	if err != nil {
		return err
	}

	if err := validateDocument(errSchemaBz, errResp); err != nil {
		return ErrInvalidResponseErr(DefaultCodespace, err.Error())
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

// PricingSchema is the Pricing JSON Schema
const PricingSchema = `
{
	"$schema": "http://json-schema.org/draft-04/schema#",
	"title": "Irishub Service Pricing",
	"description": "The Irishub Service Pricing specification",
	"type": "object",
	"definitions": {
	  "coin": {
		"description": "price coin",
		"type": "object",
		"properties": {
		  "denom": {
			"description": "the denomination of the coin",
			"type": "string",
			"pattern": "^([a-z][0-9a-z]{2}[:])?(([a-z][a-z0-9]{2,7}|x)\\.)?([a-z][a-z0-9]{2,7})(-[a-z]{3,5})?$"
		  },
		  "amount": {
			"description": "the amount of the coin",
			"type": "string",
			"pattern": "^[0-9]+(\\.[0-9]+)?$"
		  }
		},
		"additionalProperties": false,
		"required": [
		  "denom",
		  "amount"
		]
	  },
	  "discount": {
		"description": "promotion discount",
		"type": "number",
		"minimum": 0,
		"exclusiveMinimum": true,
		"maximum": 1,
		"exclusiveMaximum": true
	  },
	  "promotion_by_time": {
		"description": "promotion by time",
		"type": "object",
		"properties": {
		  "start_time": {
			"description": "starting time of the promotion",
			"type": "integer",
			"minimum": 0,
			"exclusiveMinimum": true
		  },
		  "end_time": {
			"description": "ending time of the promotion",
			"type": "integer",
			"minimum": 0,
			"exclusiveMinimum": true
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
		"description": "base price",
		"type": "array",
		"items": {
		  "$ref": "#/definitions/coin"
		},
		"uniqueItems": true
	  },
	  "promotions_by_time": {
		"description": "promotions by time",
		"type": "array",
		"items": {
		  "$ref": "#/definitions/promotion_by_time"
		},
		"uniqueItems": true
	  },
	  "promotions_by_volume": {
		"description": "promotions by volume",
		"type": "array",
		"items": {
		  "$ref": "#/definitions/promotion_by_volume"
		},
		"uniqueItems": true
	  }
	},
	"additionalProperties": false,
	"required": [
	  "price"
	]
  }
`
