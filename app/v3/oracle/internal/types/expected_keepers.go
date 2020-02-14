package types

import sdk "github.com/irisnet/irishub/types"

//expected Service keeper
type ServiceKeeper interface {
	GetRequestContext(ctx sdk.Context, requestContextID []byte) //TODO

	CreateRequestContext(ctx sdk.Context) ([]byte, error)

	UpdateRequestContext(ctx sdk.Context, requestContextID []byte) error

	StartRequestContext(ctx sdk.Context, requestContextID []byte) error

	PauseRequestContext(ctx sdk.Context, requestContextID []byte) error

	KillRequestContext(ctx sdk.Context, requestContextID []byte) error
}
