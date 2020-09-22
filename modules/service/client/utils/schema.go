package utils

import (
	"encoding/json"
	"regexp"
)

var (
	reSchemaReplace = regexp.MustCompile(`[\n\t]`)
)

// SchemaType defines the schema type
type SchemaType string

// String implements fmt.Stringer
func (schema SchemaType) String() string {
	return string(schema)
}

// MarshalJSON returns the JSON representation
func (schema SchemaType) MarshalJSON() ([]byte, error) {
	return []byte(schema.String()), nil
}

// MarshalYAML returns the YAML representation
func (schema SchemaType) MarshalYAML() (interface{}, error) {
	rawStr := schema.String()
	prettyStr := reSchemaReplace.ReplaceAllString(rawStr, "")

	return prettyStr, nil
}

// UnmarshalJSON unmarshals raw JSON bytes into a SchemaType.
func (schema *SchemaType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return nil
	}

	*schema = SchemaType(s)
	return nil
}
