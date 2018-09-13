package parameter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

type Parameter interface {
	InitGenesis()

	GetStoreKey() string

	SetReadWriter(setter params.Setter)

	SaveValue(ctx sdk.Context)

	LoadValue(ctx sdk.Context) bool
}

type SignalParameter interface {
	Parameter
}

type GovParameter interface {
	Parameter

	Valid(json string) error

	Update(ctx sdk.Context, json string)

	ToJson() string
}

type GovArrayParameter interface {
	GovParameter

	LoadValueByKey(ctx sdk.Context, key string) bool

	Insert(ctx sdk.Context, json string)
}