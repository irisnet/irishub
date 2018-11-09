package service

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSvcDef{}, "iris-hub/service/MsgSvcDef", nil)
	cdc.RegisterConcrete(MsgSvcBind{}, "iris-hub/service/MsgSvcBinding", nil)
	cdc.RegisterConcrete(MsgSvcBindingUpdate{}, "iris-hub/service/MsgSvcBindingUpdate", nil)
	cdc.RegisterConcrete(MsgSvcRefundDeposit{}, "iris-hub/service/MsgSvcRefundDeposit", nil)

	cdc.RegisterConcrete(SvcDef{}, "iris-hub/service/SvcDef", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
