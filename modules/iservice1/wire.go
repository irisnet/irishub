package iservice

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgSvcDef{}, "iris-hub/iservice1/MsgSvcDef", nil)
	cdc.RegisterConcrete(MsgSvcBind{}, "iris-hub/iservice1/MsgSvcBinding", nil)
	cdc.RegisterConcrete(MsgSvcBindingUpdate{}, "iris-hub/iservice1/MsgSvcBindingUpdate", nil)
	cdc.RegisterConcrete(MsgSvcRefundDeposit{}, "iris-hub/iservice1/MsgSvcRefundDeposit", nil)

	cdc.RegisterConcrete(SvcDef{}, "iris-hub/iservice1/SvcDef", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	RegisterWire(msgCdc)
}
