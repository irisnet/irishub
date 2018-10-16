package record

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {

	cdc.RegisterConcrete(MsgSubmitFile{}, "cosmos-sdk/MsgSubmitFile", nil)
}

var msgCdc = wire.NewCodec()
