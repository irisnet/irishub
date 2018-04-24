package iservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/irisnet/iris-hub/modules/iservice/bind"
	"github.com/irisnet/iris-hub/modules/iservice/call"
	"github.com/irisnet/iris-hub/modules/iservice/def"
	"reflect"
)

func NewIServiceHander(route IServiceRoute) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		service := route[reflect.TypeOf(msg).String()]
		if ctx.IsCheckTx() {
			service.CheckTx(ctx, msg)
		}
		return service.DeliverTx(ctx, msg)
	}
}

func NewIServiceRoute(key sdk.StoreKey, coinKeeper bank.CoinKeeper) IServiceRoute {

	svcDefKeeper := def.NewSvcDefKeeper(key, coinKeeper)
	svcDefService := def.NewSvcDefService(svcDefKeeper)

	svcBindKeeper := bind.NewSvcBindKeeper(key, coinKeeper)
	svcBindService := bind.NewSvcBindService(svcBindKeeper)

	svcCallKeeper := call.NewSvcCallKeeper(key, coinKeeper)
	svcCallService := call.NewSvcCallService(svcCallKeeper)

	return IServiceRoute{
		reflect.TypeOf(def.SvcDefMsg{}).String():   svcDefService,
		reflect.TypeOf(bind.SvcBindMsg{}).String(): svcBindService,
		reflect.TypeOf(call.SvcReqMsg{}).String():  svcCallService,
	}
}
