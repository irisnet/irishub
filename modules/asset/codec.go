package asset

import (
	"github.com/cosmos/cosmos-sdk/codec"

	token "github.com/irisnet/irishub/modules/asset/01-token"
)

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	token.RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}

// RegisterCodec registers concrete types on the codec.
func RegisterCodec(cdc *codec.Codec) {
	token.RegisterCodec(cdc)
}
