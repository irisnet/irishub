package service

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSvcDef{}, "iris-hub/service/MsgSvcDef", nil)
	cdc.RegisterConcrete(MsgSvcBind{}, "iris-hub/service/MsgSvcBinding", nil)
	cdc.RegisterConcrete(MsgSvcBindingUpdate{}, "iris-hub/service/MsgSvcBindingUpdate", nil)
	cdc.RegisterConcrete(MsgSvcDisable{}, "iris-hub/service/MsgSvcDisable", nil)
	cdc.RegisterConcrete(MsgSvcEnable{}, "iris-hub/service/MsgSvcEnable", nil)
	cdc.RegisterConcrete(MsgSvcRefundDeposit{}, "iris-hub/service/MsgSvcRefundDeposit", nil)

	cdc.RegisterConcrete(SvcDef{}, "iris-hub/service/SvcDef", nil)
	cdc.RegisterConcrete(SvcBinding{}, "iris-hub/service/SvcBinding", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
