package iservice

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgSvcDef{}, "iris-hub/iservice/MsgSvcDef", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	RegisterWire(msgCdc)
}
