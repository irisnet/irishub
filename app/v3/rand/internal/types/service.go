package types

import (
	"github.com/tendermint/tendermint/crypto"

	"github.com/irisnet/irishub/app/v3/service"
	sdk "github.com/irisnet/irishub/types"
)

const (
	ServiceName          = "random"
	ServiceDesc          = "system service definition of rand module"
	ServiceValueJsonPath = "seed"
	AuthorDescription    = "rand module account"
	ServiceSchemas       = `
	{
		"$schema": "http://json-schema.org/draft-04/schema#",
		"title": "irishub-random-seed",
		"description": "IRIS Hub Random Seed Schema",
		"type": "object",
		"input": {
			"type": "object",
			"additionalProperties": false
		},
		"output": {
			"type": "object",
			"properties": {
				"seed": {
					"description": "random seed",
					"type": "string",
					"pattern": "^[0-9a-fA-F]{64}$"
				}
			},
			"additionalProperties": false,
			"required": [
				"seed"
			]
		}
	}
	`
)

var (
	ServiceTags = []string{ModuleName}
	Auther      = sdk.AccAddress(crypto.AddressHash([]byte(ModuleName)))
)

func GetSvcDefinitions() []service.ServiceDefinition {
	return []service.ServiceDefinition{
		service.NewServiceDefinition(
			ServiceName,
			ServiceDesc,
			ServiceTags,
			Auther,
			AuthorDescription,
			ServiceSchemas,
		),
	}
}
