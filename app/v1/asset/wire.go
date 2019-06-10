package asset

import (
	"github.com/irisnet/irishub/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&MsgIssueAsset{}, "irishub/asset/MsgIssueAsset", nil)
	cdc.RegisterConcrete(&MsgCreateGateway{}, "irishub/asset/MsgCreateGateway", nil)
	cdc.RegisterConcrete(&MsgEditGateway{}, "irishub/asset/MsgEditGateway", nil)

	cdc.RegisterInterface((*Asset)(nil), nil)
	cdc.RegisterConcrete(&BaseAsset{}, "irishub/asset/BaseAsset", nil)
	cdc.RegisterConcrete(&FungibleToken{}, "irishub/asset/FungibleToken", nil)
	cdc.RegisterConcrete(&NonFungibleToken{}, "irishub/asset/NonFungibleToken", nil)
	cdc.RegisterConcrete(&Gateway{}, "irishub/asset/Gateway", nil)

	cdc.RegisterConcrete(&Params{}, "irishub/asset/Params", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
