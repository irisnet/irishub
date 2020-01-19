package types

import (
	"github.com/irisnet/irishub/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgDefineService{}, "irishub/service/MsgDefineService", nil)
	cdc.RegisterConcrete(MsgBindService{}, "irishub/service/MsgBindService", nil)
	cdc.RegisterConcrete(MsgUpdateServiceBinding{}, "irishub/service/MsgUpdateServiceBinding", nil)
	cdc.RegisterConcrete(MsgSetWithdrawAddress{}, "irishub/service/MsgSetWithdrawAddress", nil)
	cdc.RegisterConcrete(MsgDisableService{}, "irishub/service/MsgDisableService", nil)
	cdc.RegisterConcrete(MsgEnableService{}, "irishub/service/MsgEnableService", nil)
	cdc.RegisterConcrete(MsgRefundServiceDeposit{}, "irishub/service/MsgRefundServiceDeposit", nil)
	cdc.RegisterConcrete(MsgSvcRequest{}, "irishub/service/MsgSvcRequest", nil)
	cdc.RegisterConcrete(MsgSvcResponse{}, "irishub/service/MsgSvcResponse", nil)
	cdc.RegisterConcrete(MsgSvcRefundFees{}, "irishub/service/MsgSvcRefundFees", nil)
	cdc.RegisterConcrete(MsgSvcWithdrawFees{}, "irishub/service/MsgSvcWithdrawFees", nil)
	cdc.RegisterConcrete(MsgSvcWithdrawTax{}, "irishub/service/MsgSvcWithdrawTax", nil)

	cdc.RegisterConcrete(ServiceDefinition{}, "irishub/service/ServiceDefinition", nil)
	cdc.RegisterConcrete(ServiceBinding{}, "irishub/service/ServiceBinding", nil)
	cdc.RegisterConcrete(SvcRequest{}, "irishub/service/SvcRequest", nil)
	cdc.RegisterConcrete(SvcResponse{}, "irishub/service/SvcResponse", nil)
	cdc.RegisterConcrete(IncomingFee{}, "irishub/service/IncomingFee", nil)
	cdc.RegisterConcrete(ReturnedFee{}, "irishub/service/ReturnedFee", nil)

	cdc.RegisterConcrete(&Params{}, "irishub/service/Params", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
