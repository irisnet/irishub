package rand

import (
	"github.com/irisnet/irishub/codec"
)

// Register concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgRequestRand{}, "irishub/rand/MsgRequestRand", nil)

	cdc.RegisterConcrete(&Rand{}, "irishub/rand/Rand", nil)
	cdc.RegisterConcrete(&Request{}, "irishub/rand/Request", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
