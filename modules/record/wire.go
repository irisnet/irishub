package record

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {

	cdc.RegisterConcrete(MsgSubmitRecord{}, "cosmos-sdk/MsgSubmitRecord", nil)
}

var msgCdc = codec.New()
