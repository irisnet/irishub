package types

import (
	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	service "github.com/irisnet/irismod/modules/service/exported"
	servicetypes "github.com/irisnet/irismod/modules/service/types"
)

// accountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
}

// bankKeeper defines the expected bank keeper for module accounts (noalias)
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}

//expected Service keeper
type ServiceKeeper interface {
	RegisterResponseCallback(
		moduleName string, respCallback service.ResponseCallback,
	) error

	RegisterStateCallback(
		moduleName string, stateCallback service.StateCallback,
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

	ServiceBindingsIterator(ctx sdk.Context, serviceName string) sdk.Iterator

	GetParams(ctx sdk.Context) servicetypes.Params
}
