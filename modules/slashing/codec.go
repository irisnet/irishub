package slashing

import (
	"github.com/irisnet/irishub/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgUnjail{}, "irishub/slashing/MsgUnjail", nil)
	cdc.RegisterConcrete(&Params{}, "irishub/slashing/Params", nil)
}

var cdcEmpty = codec.New()
