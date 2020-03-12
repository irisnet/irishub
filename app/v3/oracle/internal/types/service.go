package types

import (
	"github.com/tendermint/tendermint/crypto"

	"github.com/irisnet/irishub/app/v3/service"
	sdk "github.com/irisnet/irishub/types"
)

const (
	ServiceName       = "feed"
	ServiceDesc       = "system moudle service for oracle"
	ServiceSchemas    = `{"input":{"type":"object"},"output":{"type":"object"},"error":{"type":"object"}}`
	AuthorDescription = "oracle module account"
)

var (
	ServiceTags = []string{"oracle"}
	Auther      = sdk.AccAddress(crypto.AddressHash([]byte("oracle")))
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
