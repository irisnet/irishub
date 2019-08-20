package bank

import (
	"github.com/irisnet/irishub/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "irishub/bank/Send", nil)
	cdc.RegisterConcrete(MsgBurn{}, "irishub/bank/Burn", nil)
	cdc.RegisterConcrete(MsgSetMemoRegexp{}, "irishub/bank/SetMemoRegexp", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
