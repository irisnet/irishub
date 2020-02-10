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
	cdc.RegisterConcrete(MsgRequestService{}, "irishub/service/MsgRequestService", nil)
	cdc.RegisterConcrete(MsgRespondService{}, "irishub/service/MsgRespondService", nil)
	cdc.RegisterConcrete(MsgStopRepeated{}, "irishub/service/MsgStopRepeated", nil)
	cdc.RegisterConcrete(MsgUpdateRequestContext{}, "irishub/service/MsgUpdateRequestContext", nil)
	cdc.RegisterConcrete(MsgWithdrawEarnedFees{}, "irishub/service/MsgWithdrawEarnedFees", nil)
	cdc.RegisterConcrete(MsgWithdrawTax{}, "irishub/service/MsgWithdrawTax", nil)

	cdc.RegisterConcrete(ServiceDefinition{}, "irishub/service/ServiceDefinition", nil)
	cdc.RegisterConcrete(ServiceBinding{}, "irishub/service/ServiceBinding", nil)
	cdc.RegisterConcrete(RequestContext{}, "irishub/service/RequestContext", nil)
	cdc.RegisterConcrete(CompactRequest{}, "irishub/service/CompactRequest", nil)
	cdc.RegisterConcrete(Request{}, "irishub/service/Request", nil)
	cdc.RegisterConcrete(Response{}, "irishub/service/Response", nil)
	cdc.RegisterConcrete(EarnedFees{}, "irishub/service/EarnedFees", nil)

	cdc.RegisterConcrete(&Params{}, "irishub/service/Params", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
