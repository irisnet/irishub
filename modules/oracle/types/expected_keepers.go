package types

import (
	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"

	service "github.com/irismod/service/exported"

	guardiantypes "github.com/irisnet/irishub/modules/guardian/types"
)

// expected Service keeper
type ServiceKeeper interface {
	RegisterResponseCallback(
		moduleName string, respCallback service.ResponseCallback,
	) error

	RegisterStateCallback(
		moduleName string, stateCallback service.StateCallback,
	) error

	RegisterModuleService(
		moduleName string, moduleService *service.ModuleService,
	) error

	GetRequestContext(
		ctx sdk.Context, requestContextID tmbytes.HexBytes,
	) (service.RequestContext, bool)

	CreateRequestContext(
		ctx sdk.Context,
		serviceName string,
		providers []sdk.AccAddress,
		consumer sdk.AccAddress,
		input string,
		serviceFeeCap sdk.Coins,
		timeout int64,
		superMode bool,
		repeated bool,
		repeatedFrequency uint64,
		repeatedTotal int64,
		state service.RequestContextState,
		responseThreshold uint32,
		moduleName string,
	) (tmbytes.HexBytes, error)

	UpdateRequestContext(
		ctx sdk.Context,
		requestContextID tmbytes.HexBytes,
		providers []sdk.AccAddress,
		respThreshold uint32,
		serviceFeeCap sdk.Coins,
		timeout int64,
		repeatedFreq uint64,
		repeatedTotal int64,
		consumer sdk.AccAddress,
	) error

	StartRequestContext(
		ctx sdk.Context,
		requestContextID tmbytes.HexBytes,
		consumer sdk.AccAddress,
	) error

	PauseRequestContext(
		ctx sdk.Context,
		requestContextID tmbytes.HexBytes,
		consumer sdk.AccAddress,
	) error
}

// GuardianKeeper defines the expected guardian keeper (noalias)
type GuardianKeeper interface {
	GetProfiler(ctx sdk.Context, addr sdk.AccAddress) (guardian guardiantypes.Guardian, found bool)
}

var (
	RequestContextStateFromString = service.RequestContextStateFromString
)
