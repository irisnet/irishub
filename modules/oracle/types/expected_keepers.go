package types

import (
	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	service "github.com/irisnet/irismod/modules/service/exported"
)

// ServiceKeeper defines the expected service keeper (noalias)
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
	AddServiceBinding(
		ctx sdk.Context,
		serviceName string,
		provider sdk.AccAddress,
		deposit sdk.Coins,
		pricing string,
		qos uint64,
		options string,
		owner sdk.AccAddress,
	) error
	AddServiceDefinition(
		ctx sdk.Context,
		name,
		description string,
		tags []string,
		author sdk.AccAddress,
		authorDescription,
		schemas string,
	) error
}

// AuthKeeper defines the expected auth keeper (noalias)
type AuthKeeper interface {
	Authorized(ctx sdk.Context, addr sdk.AccAddress) bool
}

var (
	RequestContextStateFromString = service.RequestContextStateFromString
)

type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
}
