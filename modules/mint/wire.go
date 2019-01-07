package mint

import (
	"github.com/irisnet/irishub/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(Minter{}, "irishub/mint/Minter", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
