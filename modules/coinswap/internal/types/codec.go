package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// module codec
var ModuleCdc = codec.New()

// RegisterCodec registers concrete types on the codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSwapOrder{}, "irishub/coinswap/MsgSwapOrder", nil)
	cdc.RegisterConcrete(MsgAddLiquidity{}, "irishub/coinswap/MsgAddLiquidity", nil)
	cdc.RegisterConcrete(MsgRemoveLiquidity{}, "irishub/coinswap/MsgRemoveLiquidity", nil)
	cdc.RegisterConcrete(&Params{}, "irishub/coinswap/Params", nil)
}

func init() {
	RegisterCodec(ModuleCdc)
}
