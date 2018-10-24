package record

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {

	cdc.RegisterConcrete(MsgSubmitRecord{}, "cosmos-sdk/MsgSubmitRecord", nil)
}

var msgCdc = wire.NewCodec()
