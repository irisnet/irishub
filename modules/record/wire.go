package record

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgRecord{}, "iris-hub/record/MsgRecord", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	RegisterWire(msgCdc)
}
