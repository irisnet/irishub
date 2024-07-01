package types

import (
	"github.com/cometbft/cometbft/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	RandomServiceName          = "random"
	RandomServiceDesc          = "system service definition of random module"
	RandomServiceValueJSONPath = "seed"
	RandomAuthorDescription    = "random module account"
	RandomServiceSchemas       = `
	{
		"input": {
			"$schema": "http://json-schema.org/draft-04/schema#",
			"title": "random-seed-input-body",
			"description": "IRIS Hub Random Seed Input Body Schema",
			"type": "object",
			"additionalProperties": false
		},
		"output": {
			"$schema": "http://json-schema.org/draft-04/schema#",
			"title": "random-seed-output-body",
			"description": "IRIS Hub Random Seed Output Body Schema",
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
	}`
)

var (
	RandomServiceTags = []string{"oracle"}
	RandomAuthor      = sdk.AccAddress(crypto.AddressHash([]byte("oracle")))
)

func GetRandomSvcDefinition() ServiceDefinition {
	return NewServiceDefinition(
		RandomServiceName,
		RandomServiceDesc,
		RandomServiceTags,
		RandomAuthor,
		RandomAuthorDescription,
		RandomServiceSchemas,
	)
}