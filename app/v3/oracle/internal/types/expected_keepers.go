package types

import (
	service "github.com/irisnet/irishub/app/v3/service/exported"
	"github.com/irisnet/irishub/modules/guardian"
	sdk "github.com/irisnet/irishub/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

//expected Service keeper
type ServiceKeeper interface {
	RegisterResponseCallback(moduleName string,
		respCallback service.ResponseCallback) sdk.Error

	RegisterStateCallback(moduleName string,
		stateCallback service.StateCallback) sdk.Error

	GetRequestContext(ctx sdk.Context,
		requestContextID cmn.HexBytes) (service.RequestContext, bool)

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
		respHandler string) (cmn.HexBytes, sdk.Tags, sdk.Error)

	UpdateRequestContext(ctx sdk.Context,
		requestContextID cmn.HexBytes,
		providers []sdk.AccAddress,
		serviceFeeCap sdk.Coins,
		timeout int64,
		repeatedFreq uint64,
		repeatedTotal int64,
		consumer sdk.AccAddress) sdk.Error

	StartRequestContext(ctx sdk.Context,
		requestContextID cmn.HexBytes,
		consumer sdk.AccAddress) sdk.Error

	PauseRequestContext(ctx sdk.Context,
		requestContextID cmn.HexBytes,
		consumer sdk.AccAddress) sdk.Error
}

// GuardianKeeper defines the expected guardian keeper (noalias)
type GuardianKeeper interface {
	GetProfiler(ctx sdk.Context, addr sdk.AccAddress) (guardian guardian.Guardian, found bool)
}

var (
	RequestContextStateFromString = service.RequestContextStateFromString
)
