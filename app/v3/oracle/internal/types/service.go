package types

import (
	"github.com/tendermint/tendermint/crypto"

	"github.com/irisnet/irishub/app/v3/service"
	sdk "github.com/irisnet/irishub/types"
)

const (
	ServiceName       = ModuleName
	ServiceDesc       = "system service definition of oracle moudle"
	ServiceSchemas    = `{"input":{"type":"object"},"output":{"type":"object"},"error":{"type":"object"}}`
	AuthorDescription = "oracle module account"
)

var (
	ServiceTags = []string{ModuleName}
	Auther      = sdk.AccAddress(crypto.AddressHash([]byte("oracle")))
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
