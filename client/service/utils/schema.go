package utils

import (
	"bytes"
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
	buf := bytes.NewBuffer([]byte{})
	if err := json.Compact(buf, []byte(schema)); err != nil {
		return nil, err
	}

	return json.Marshal(buf.String())
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
