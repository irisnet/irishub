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

// ValidateResponseResult validates the response result against the result schema
func ValidateResponseResult(result string) sdk.Error {
	if err := validateDocument([]byte(ResultSchema), result); err != nil {
		return ErrInvalidResponseResult(DefaultCodespace, err.Error())
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
	"title": "irishub-service-pricing",
	"description": "IRIS Hub Service Pricing Schema",
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
		"description": "base price in main unit, e.g. 0.5iris",
		"type": "string",
		"pattern": "^\\d+(\\.\\d+)?[a-z][a-z0-9]{2,7}(,\\d+(\\.\\d+)?[a-z][a-z0-9]{2,7})*$"
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
}
`

	// ResultSchema is the JSON Schema for the response result
	ResultSchema = `
{
	"$schema": "http://json-schema.org/draft-04/schema#",
	"title": "irishub-service-result",
	"description": "IRIS Hub Service Result Schema",
	"type": "object",
	"properties": {
	  "code": {
		"description": "result code",
		"type": "integer",
		"enum": [200, 400, 500]
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
}
`
)
