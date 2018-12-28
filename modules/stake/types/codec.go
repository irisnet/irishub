package types

import (
	"github.com/irisnet/irishub/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateValidator{}, "irishub/stake/MsgCreateValidator", nil)
	cdc.RegisterConcrete(MsgEditValidator{}, "irishub/stake/MsgEditValidator", nil)
	cdc.RegisterConcrete(MsgDelegate{}, "irishub/stake/MsgDelegate", nil)
	cdc.RegisterConcrete(MsgBeginUnbonding{}, "irishub/stake/BeginUnbonding", nil)
	cdc.RegisterConcrete(MsgBeginRedelegate{}, "irishub/stake/BeginRedelegate", nil)
}

// generic sealed codec to be used throughout sdk
var MsgCdc *codec.Codec

func init() {
	cdc := codec.New()
	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	MsgCdc = cdc.Seal()
}
