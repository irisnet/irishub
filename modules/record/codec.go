package record

import (
	"github.com/irisnet/irishub/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSubmitRecord{}, "irishub/record/MsgSubmitRecord", nil)
}

var msgCdc = codec.New()
