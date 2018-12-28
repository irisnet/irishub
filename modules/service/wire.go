package service

import (
	"github.com/irisnet/irishub/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSvcDef{}, "irishub/service/MsgSvcDef", nil)
	cdc.RegisterConcrete(MsgSvcBind{}, "irishub/service/MsgSvcBinding", nil)
	cdc.RegisterConcrete(MsgSvcBindingUpdate{}, "irishub/service/MsgSvcBindingUpdate", nil)
	cdc.RegisterConcrete(MsgSvcDisable{}, "irishub/service/MsgSvcDisable", nil)
	cdc.RegisterConcrete(MsgSvcEnable{}, "irishub/service/MsgSvcEnable", nil)
	cdc.RegisterConcrete(MsgSvcRefundDeposit{}, "irishub/service/MsgSvcRefundDeposit", nil)
	cdc.RegisterConcrete(MsgSvcRequest{}, "irishub/service/MsgSvcRequest", nil)
	cdc.RegisterConcrete(MsgSvcResponse{}, "irishub/service/MsgSvcResponse", nil)
	cdc.RegisterConcrete(MsgSvcRefundFees{}, "irishub/service/MsgSvcRefundFees", nil)
	cdc.RegisterConcrete(MsgSvcWithdrawFees{}, "irishub/service/MsgSvcWithdrawFees", nil)
	cdc.RegisterConcrete(MsgSvcWithdrawTax{}, "irishub/service/MsgSvcWithdrawTax", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
