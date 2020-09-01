package types

import (
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	servicetypes "github.com/irismod/service/types"

	"github.com/irisnet/irishub/modules/oracle/types"
)

const (
	ServiceName          = "random"
	ServiceDesc          = "system service definition of random module"
	ServiceValueJsonPath = "seed"
	AuthorDescription    = "random module account"
	ServiceSchemas       = `
	{
		"input": {
			"$schema": "http://json-schema.org/draft-04/schema#",
			"title": "irishub-random-seed-input",
			"description": "IRIS Hub Random Seed Input Schema",
			"type": "object",
			"additionalProperties": false
		},
		"output": {
			"$schema": "http://json-schema.org/draft-04/schema#",
			"title": "irishub-random-seed-output",
			"description": "IRIS Hub Random Seed Output Schema",
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
	ServiceTags = []string{types.ModuleName}
	Author      = sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))
)

func GetSvcDefinitions() []servicetypes.ServiceDefinition {
	return []servicetypes.ServiceDefinition{
		servicetypes.NewServiceDefinition(
			ServiceName,
			ServiceDesc,
			ServiceTags,
			Author,
			AuthorDescription,
			ServiceSchemas,
		),
	}
}
