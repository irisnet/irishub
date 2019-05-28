package asset

import (
	"github.com/irisnet/irishub/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {

	// TODO
	cdc.RegisterConcrete(&Params{}, "irishub/asset/Params", nil)
	cdc.RegisterConcrete(&MsgCreateGateway{}, "irishub/asset/MsgCreateGateway", nil)
	cdc.RegisterConcrete(&MsgEditGateway{}, "irishub/asset/MsgEditGateway", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
