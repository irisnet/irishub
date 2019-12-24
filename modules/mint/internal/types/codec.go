package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

// RegisterCodec registers concrete types on the codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&Params{}, "irishub/mint/Params", nil)
}

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
