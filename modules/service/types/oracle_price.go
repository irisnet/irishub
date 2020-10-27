package types

import (
	"fmt"
	"time"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	RegisterModuleName           = "oracle"
	OraclePriceServiceName       = "oracle-price"
	OraclePriceValueJSONPath     = "rate"
	OraclePriceServiceDesc       = "system service definition of oracle module"
	OraclePriceAuthorDescription = "oracle module account"
	OraclePriceSchemas           = `
	{
		"input": {
			"$schema": "http://json-schema.org/draft-04/schema#",
			"title": "iservice-oracle-price-input-body",
			"description": "Interchain Service Oracle Price Input Body Schema",
			"type": "object",
			"properties": {
				"pair": {
					"description": "exchange pair",
					"type": "string",
					"pattern": "^[0-9a-zA-Z]+-[0-9a-zA-Z]+$"
				}
			},
			"additionalProperties": false,
			"required": [
				"pair"
			]
		},
		"output": {
			"$schema": "http://json-schema.org/draft-04/schema#",
			"title": "iservice-oracle-price-output-body",
			"description": "Interchain Service Oracle Price Output Body Schema",
			"type": "object",
			"properties": {
				"rate": {
					"description": "exchange rate",
					"type": "string",
					"pattern": "^(?:[1-9]+\\d*|0)(\\.\\d+)?$"
				}
			},
			"additionalProperties": false,
			"required": [
				"rate"
			]
		}
	}`
)

var (
	OraclePriceServiceTags     = []string{"oracle"}
	OraclePriceServiceAuthor   = sdk.AccAddress(crypto.AddressHash([]byte(RegisterModuleName)))
	OraclePriceServiceProvider = sdk.AccAddress(crypto.AddressHash([]byte(RegisterModuleName)))
)

func GenOraclePriceSvcDefinition() ServiceDefinition {
	return ServiceDefinition{
		Name:              OraclePriceServiceName,
		Description:       OraclePriceAuthorDescription,
		Tags:              OraclePriceServiceTags,
		Author:            OraclePriceServiceAuthor,
		AuthorDescription: OraclePriceAuthorDescription,
		Schemas:           OraclePriceSchemas,
	}
}

func GenOraclePriceSvcBinding(baseDenom string) ServiceBinding {
	return ServiceBinding{
		ServiceName:  OraclePriceServiceName,
		Provider:     OraclePriceServiceProvider,
		Deposit:      sdk.NewCoins(sdk.NewCoin(baseDenom, sdk.NewInt(0))),
		Pricing:      fmt.Sprintf(`{"price": "0%s"}`, baseDenom),
		QoS:          1,
		Options:      `{}`,
		Available:    true,
		DisabledTime: time.Time{},
		Owner:        OraclePriceServiceProvider,
	}
}
