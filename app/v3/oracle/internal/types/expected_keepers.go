package types

import (
	service "github.com/irisnet/irishub/app/v3/service/exported"
	"github.com/irisnet/irishub/modules/guardian"
	sdk "github.com/irisnet/irishub/types"
)

//expected Service keeper
type ServiceKeeper interface {
	RegisterResponseCallback(moduleName string,
		respCallback service.ResponseCallback) sdk.Error

	GetRequestContext(ctx sdk.Context,
		requestContextID []byte) (service.RequestContext, bool)

	CreateRequestContext(ctx sdk.Context,
		serviceName string,
		providers []sdk.AccAddress,
		consumer sdk.AccAddress,
		input string,
		serviceFeeCap sdk.Coins,
		timeout int64,
		superMode,
		repeated bool,
		repeatedFrequency uint64,
		repeatedTotal int64,
		state service.RequestContextState,
		respThreshold uint16,
		respHandler string) ([]byte, sdk.Error)

	UpdateRequestContext(ctx sdk.Context,
		requestContextID []byte,
		providers []sdk.AccAddress,
		serviceFeeCap sdk.Coins,
		timeout int64,
		repeatedFreq uint64,
		repeatedTotal int64,
		consumer sdk.AccAddress) sdk.Error

	StartRequestContext(ctx sdk.Context,
		requestContextID []byte,
		consumer sdk.AccAddress) sdk.Error

	PauseRequestContext(ctx sdk.Context,
		requestContextID []byte,
		consumer sdk.AccAddress) sdk.Error
}

// GuardianKeeper defines the expected guardian keeper (noalias)
type GuardianKeeper interface {
	GetProfiler(ctx sdk.Context, addr sdk.AccAddress) (guardian guardian.Guardian, found bool)
}

var (
	RequestContextStateFromString = service.RequestContextStateFromString
)
