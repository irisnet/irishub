package types

import (
	"github.com/irisnet/irishub/codec"
)

// ModuleCdc generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

// RegisterCodec registers concrete types on the codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSwapOrder{}, "irishub/coinswap/MsgSwapOrder", nil)
	cdc.RegisterConcrete(MsgAddLiquidity{}, "irishub/coinswap/MsgAddLiquidity", nil)
	cdc.RegisterConcrete(MsgRemoveLiquidity{}, "irishub/coinswap/MsgRemoveLiquidity", nil)
	cdc.RegisterConcrete(&Params{}, "irishub/coinswap/Params", nil)
}

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
