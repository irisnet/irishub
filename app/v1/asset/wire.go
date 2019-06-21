package asset

import (
	"github.com/irisnet/irishub/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgIssueToken{}, "irishub/asset/MsgIssueToken", nil)
	cdc.RegisterConcrete(MsgCreateGateway{}, "irishub/asset/MsgCreateGateway", nil)
	cdc.RegisterConcrete(MsgEditGateway{}, "irishub/asset/MsgEditGateway", nil)
	cdc.RegisterConcrete(MsgEditToken{}, "irishub/asset/MsgEditToken", nil)

	cdc.RegisterConcrete(BaseToken{}, "irishub/asset/BaseToken", nil)
	cdc.RegisterConcrete(FungibleToken{}, "irishub/asset/FungibleToken", nil)

	cdc.RegisterConcrete(&Params{}, "irishub/asset/Params", nil)
	cdc.RegisterConcrete(&Gateway{}, "irishub/asset/Gateway", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
