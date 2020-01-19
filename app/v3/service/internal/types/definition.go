package types

import (
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

// ServiceDefinition represents a service definition
type ServiceDefinition struct {
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	Tags              []string       `json:"tags"`
	Author            sdk.AccAddress `json:"author"`
	AuthorDescription string         `json:"author_description"`
	Schemas           string         `json:"schemas"`
}

// NewServiceDefinition creates a new ServiceDefinition instance
func NewServiceDefinition(name, description string, tags []string, author sdk.AccAddress, authorDescription, schemas string) ServiceDefinition {
	return ServiceDefinition{
		Name:              name,
		Description:       description,
		Tags:              tags,
		Author:            author,
		AuthorDescription: authorDescription,
		Schemas:           schemas,
	}
}

// String implements fmt.Stringer
func (svcDef ServiceDefinition) String() string {
	return fmt.Sprintf(`ServiceDefinition:
		Name:                  %s
		Description:           %s
		Tags:                  %v
		Author:                %s
		AuthorDescription:     %s
		Schemas:               %s`,
		svcDef.Name, svcDef.Description, svcDef.Tags,
		svcDef.Author, svcDef.AuthorDescription, svcDef.Schemas,
	)
}
