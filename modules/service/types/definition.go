package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewServiceDefinition creates a new ServiceDefinition instance
func NewServiceDefinition(
	name string,
	description string,
	tags []string,
	author sdk.AccAddress,
	authorDescription string,
	schemas string,
) ServiceDefinition {
	return ServiceDefinition{
		Name:              name,
		Description:       description,
		Tags:              tags,
		Author:            author,
		AuthorDescription: authorDescription,
		Schemas:           schemas,
	}
}

// Validate validates the service definition
func (svcDef ServiceDefinition) Validate() error {
	if err := ValidateAuthor(svcDef.Author); err != nil {
		return err
	}

	if err := ValidateServiceName(svcDef.Name); err != nil {
		return err
	}

	if err := ValidateTags(svcDef.Tags); err != nil {
		return err
	}

	if err := ValidateServiceDescription(svcDef.Description); err != nil {
		return err
	}

	if err := ValidateAuthorDescription(svcDef.AuthorDescription); err != nil {
		return err
	}

	return ValidateServiceSchemas(svcDef.Schemas)
}
