package mint

import (
	"github.com/irisnet/irishub/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	// Not Register mint codec in app, deprecated now
	//cdc.RegisterConcrete(Minter{}, "irishub/mint/Minter", nil)
	cdc.RegisterConcrete(&Params{}, "irishub/mint/Params", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
