package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgRequestRand{}, "irishub/rand/MsgRequestRand", nil)

	cdc.RegisterConcrete(&Rand{}, "irishub/rand/Rand", nil)
	cdc.RegisterConcrete(&Request{}, "irishub/rand/Request", nil)
}
// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
