package types

import (
	"github.com/tendermint/tendermint/crypto"

	"github.com/irisnet/irishub/app/v3/service"
	sdk "github.com/irisnet/irishub/types"
)

const (
	ServiceName          = ModuleName
	ServiceDesc          = "system service definition of rand module"
	ServiceSchemas       = `{"input":{"type":"object","properties":{}},"output":{"type":"object","properties":{"seed":{"description":"seed","type":"string","pattern":"^[0-9a-fA-F]{64}$"}}},"error":{"type":"string"}}`
	ServiceValueJsonPath = "seed"
	AuthorDescription    = "rand module account"
)

var (
	ServiceTags = []string{ModuleName}
	Auther      = sdk.AccAddress(crypto.AddressHash([]byte(ModuleName)))
)

func GetSvcDefinition() service.ServiceDefinition {
	return service.NewServiceDefinition(
		ServiceName,
		ServiceDesc,
		ServiceTags,
		Auther,
		AuthorDescription,
		ServiceSchemas,
	)
}
