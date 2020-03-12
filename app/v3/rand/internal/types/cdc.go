package types

import (
	"github.com/irisnet/irishub/codec"
)

var msgCdc = codec.New()

// Register concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgRequestRand{}, "irishub/rand/MsgRequestRand", nil)

	cdc.RegisterConcrete(&Rand{}, "irishub/rand/Rand", nil)
	cdc.RegisterConcrete(&Request{}, "irishub/rand/Request", nil)
}

func init() {
	RegisterCodec(msgCdc)
}
