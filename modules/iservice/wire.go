package iservice

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSvcDef{}, "iris-hub/iservice/MsgSvcDef", nil)
	cdc.RegisterConcrete(MsgSvcBind{}, "iris-hub/iservice/MsgSvcBinding", nil)
	cdc.RegisterConcrete(MsgSvcBindingUpdate{}, "iris-hub/iservice/MsgSvcBindingUpdate", nil)
	cdc.RegisterConcrete(MsgSvcRefundDeposit{}, "iris-hub/iservice/MsgSvcRefundDeposit", nil)

	cdc.RegisterConcrete(SvcDef{}, "iris-hub/iservice/SvcDef", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
