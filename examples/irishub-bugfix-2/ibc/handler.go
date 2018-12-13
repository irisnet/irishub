package ibc

import (
	"reflect"

	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/bank"
	"github.com/irisnet/irishub/modules/upgrade"
)

func NewHandler(ibcm Mapper, ck bank.Keeper, uk upgrade.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case IBCGetMsg:
			return handleIBCGetMsg(ctx, ibcm, ck, uk, msg)
		case IBCSetMsg:
			return handleIBCSetMsg(ctx, ibcm, ck, uk, msg)
		default:
			errMsg := "Unrecognized IBC Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}


// IBCTransferMsg deducts coins from the account and creates an egress IBC packet.
func handleIBCSetMsg(ctx sdk.Context, ibcm Mapper, ck bank.Keeper, uk upgrade.Keeper, msg IBCSetMsg) sdk.Result {

	if uk.OnlyRunAfterVersionId(ctx, 2) {
		ibcm.Set(ctx,msg.Addr.String()+":ibc-2")
		return sdk.Result{Log:"This is new module - ibc2 !!"}
	} else if uk.OnlyRunAfterVersionId(ctx, 1) {
		ibcm.Set(ctx,msg.Addr.String()+":ibc-1")
		return sdk.Result{Log:"This is new module - ibc1 !!"}
	} else {
		ibcm.Set(ctx,msg.Addr.String()+":ibc-0")
		return sdk.Result{Log:"This is new module - ibc0 !!"}
	}
}

// IBCReceiveMsg adds coins to the destination address and creates an ingress IBC packet.
func handleIBCGetMsg(ctx sdk.Context, ibcm Mapper, ck bank.Keeper, uk upgrade.Keeper, msg IBCGetMsg) sdk.Result {
	AddrString,_:=ibcm.Get(ctx)
	return sdk.Result{Log:AddrString}
}
