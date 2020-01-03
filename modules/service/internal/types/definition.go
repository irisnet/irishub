package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ServiceDefinition represents a service definition
type ServiceDefinition struct {
	Name              string         `json:"name" yaml:"name"`
	Description       string         `json:"description" yaml:"description"`
	Tags              []string       `json:"tags" yaml:"tags"`
	Author            sdk.AccAddress `json:"author" yaml:"author"`
	AuthorDescription string         `json:"author_description" yaml:"author_description"`
	Schema            string         `json:"schema" yaml:"schema"`
}

// NewServiceDefinition constructs a new ServiceDefinition
func NewServiceDefinition(name, description string, tags []string, author sdk.AccAddress, authorDescription, schema string) ServiceDefinition {
	return ServiceDefinition{
		Name:              name,
		Description:       description,
		Tags:              tags,
		Author:            author,
		AuthorDescription: authorDescription,
		Schema:            schema,
	}
}
