package ibc

import (
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

func NewHandler(ibcm Mapper, ck bank.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case IBCTransferMsg:
			return handleIBCTransferMsg(ctx, ibcm, ck, msg)
		case IBCReceiveMsg:
			return handleIBCReceiveMsg(ctx, ibcm, ck, msg)
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
func handleIBCTransferMsg(ctx sdk.Context, ibcm Mapper, ck bank.Keeper, msg IBCTransferMsg) sdk.Result {
	packet := msg.IBCPacket

	_, _, err := ck.SubtractCoins(ctx, packet.SrcAddr, packet.Coins)
	if err != nil {
		return err.Result()
	}

	err = ibcm.PostIBCPacket(ctx, packet)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

// IBCReceiveMsg adds coins to the destination address and creates an ingress IBC packet.
func handleIBCReceiveMsg(ctx sdk.Context, ibcm Mapper, ck bank.Keeper, msg IBCReceiveMsg) sdk.Result {
	packet := msg.IBCPacket

	seq := ibcm.GetIngressSequence(ctx, packet.SrcChain)
	if msg.Sequence != seq {
		return ErrInvalidSequence(ibcm.codespace).Result()
	}

	_, _, err := ck.AddCoins(ctx, packet.DestAddr, packet.Coins)
	if err != nil {
		return err.Result()
	}

	ibcm.SetIngressSequence(ctx, packet.SrcChain, seq+1)

	return sdk.Result{}
}



// IBCTransferMsg deducts coins from the account and creates an egress IBC packet.
func handleIBCSetMsg(ctx sdk.Context, ibcm Mapper, ck bank.Keeper, msg IBCSetMsg) sdk.Result {
	ibcm.Set(ctx,msg.Addr)
	return sdk.Result{Log:"This is new module!!"}
}

// IBCReceiveMsg adds coins to the destination address and creates an ingress IBC packet.
func handleIBCGetMsg(ctx sdk.Context, ibcm Mapper, ck bank.Keeper, msg IBCGetMsg) sdk.Result {
	Addr,_:=ibcm.Get(ctx)
	return sdk.Result{Log:Addr.String()}
}
