package bank

import (
	"github.com/irisnet/irishub/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "cosmos-sdk/bank/Send", nil)
	cdc.RegisterConcrete(MsgIssue{}, "cosmos-sdk/bank/Issue", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
