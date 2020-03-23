package types

import (
	cmn "github.com/tendermint/tendermint/libs/common"

	"github.com/irisnet/irishub/app/v3/service"
	"github.com/irisnet/irishub/app/v3/service/exported"
	sdk "github.com/irisnet/irishub/types"
)

//expected Service keeper
type ServiceKeeper interface {
	RegisterResponseCallback(
		moduleName string,
		respCallback exported.ResponseCallback,
	) sdk.Error

	GetRequestContext(
		ctx sdk.Context,
		requestContextID cmn.HexBytes,
	) (exported.RequestContext, bool)

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
		state exported.RequestContextState,
		respThreshold uint16,
		respHandler string,
	) (cmn.HexBytes, sdk.Tags, sdk.Error)

	StartRequestContext(
		ctx sdk.Context,
		requestContextID cmn.HexBytes,
		consumer sdk.AccAddress,
	) sdk.Error

	ServiceBindingsIterator(ctx sdk.Context, serviceName string) sdk.Iterator

	GetParamSet(ctx sdk.Context) service.Params
}

//expected Bank keeper
type BankKeeper interface {
	GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}
