package types

import (
	"github.com/irisnet/irishub/codec"
)

// Register concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgIssueToken{}, "irishub/asset/MsgIssueToken", nil)
	cdc.RegisterConcrete(MsgEditToken{}, "irishub/asset/MsgEditToken", nil)
	cdc.RegisterConcrete(MsgMintToken{}, "irishub/asset/MsgMintToken", nil)
	cdc.RegisterConcrete(MsgTransferTokenOwner{}, "irishub/asset/MsgTransferTokenOwner", nil)

	cdc.RegisterConcrete(FungibleToken{}, "irishub/asset/FungibleToken", nil)

	cdc.RegisterConcrete(&Params{}, "irishub/asset/Params", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
