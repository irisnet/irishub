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

// Validate validates the service definition
func (svcDef ServiceDefinition) Validate() sdk.Error {
	if len(svcDef.Author) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "author missing")
	}

	if !validServiceName(svcDef.Name) {
		return ErrInvalidServiceName(DefaultCodespace, svcDef.Name)
	}

	if err := validateServiceDefLength(svcDef); err != nil {
		return err
	}

	if sdk.HasDuplicate(svcDef.Tags) {
		return ErrDuplicateTags(DefaultCodespace)
	}

	if len(svcDef.Schemas) == 0 {
		return ErrInvalidSchemas(DefaultCodespace, "schemas missing")
	}

	return ValidateServiceSchemas(svcDef.Schemas)
}

func validateServiceDefLength(svcDef ServiceDefinition) sdk.Error {
	if err := ensureServiceNameLength(svcDef.Name); err != nil {
		return err
	}

	if len(svcDef.Description) > MaxDescriptionLength {
		return ErrInvalidLength(DefaultCodespace, fmt.Sprintf("invalid description length; got: %d, max: %d", len(svcDef.Description), MaxDescriptionLength))
	}

	if len(svcDef.Tags) > MaxTagsNum {
		return ErrInvalidLength(DefaultCodespace, fmt.Sprintf("invalid tags size; got: %d, max: %d", len(svcDef.Tags), MaxTagsNum))
	}

	for i, tag := range svcDef.Tags {
		if len(tag) == 0 {
			return ErrInvalidLength(DefaultCodespace, fmt.Sprintf("invalid tag[%d] length: tag must not be empty", i))
		}

		if len(tag) > MaxTagLength {
			return ErrInvalidLength(DefaultCodespace, fmt.Sprintf("invalid tag[%d] length; got: %d, max: %d", i, len(tag), MaxTagLength))
		}
	}

	if len(svcDef.AuthorDescription) > MaxDescriptionLength {
		return ErrInvalidLength(DefaultCodespace, fmt.Sprintf("invalid author description length; got: %d, max: %d", len(svcDef.AuthorDescription), MaxDescriptionLength))
	}

	return nil
}
