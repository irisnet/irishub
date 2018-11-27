package ibc

import (
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/bank"
)

func NewHandler(ibcm Mapper, ck bank.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case IBCGetMsg:
			return handleIBCGetMsg(ctx, ibcm, ck, msg)
		case IBCSetMsg:
			return handleIBCSetMsg(ctx, ibcm, ck, msg)
		default:
			errMsg := "Unrecognized IBC Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// IBCTransferMsg deducts coins from the account and creates an egress IBC packet.
func handleIBCSetMsg(ctx sdk.Context, ibcm Mapper, ck bank.Keeper, msg IBCSetMsg) sdk.Result {
	ibcm.Set(ctx,msg.Addr.String())
	return sdk.Result{Log:"This is new module - ibc1 !!"}
}

// IBCReceiveMsg adds coins to the destination address and creates an ingress IBC packet.
func handleIBCGetMsg(ctx sdk.Context, ibcm Mapper, ck bank.Keeper, msg IBCGetMsg) sdk.Result {
	AddrString,_:=ibcm.Get(ctx)
	return sdk.Result{Log:AddrString}
}
