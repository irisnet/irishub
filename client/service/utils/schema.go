package utils

import (
	"encoding/json"
)

// SchemaType defines the schema type
type SchemaType string

// String implements fmt.Stringer
func (schema SchemaType) String() string {
	return string(schema)
}

// MarshalJSON marshals the schema to JSON
func (schema SchemaType) MarshalJSON() ([]byte, error) {
	return []byte(schema.String()), nil
}

// UnmarshalJSON unmarshals the data to the schema
func (schema *SchemaType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return nil
	}

	*schema = SchemaType(s)
	return nil
}
