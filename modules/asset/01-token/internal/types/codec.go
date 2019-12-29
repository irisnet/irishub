package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

// RegisterCodec registers concrete types on the codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgIssueToken{}, "irishub/asset/token/MsgIssueToken", nil)
	cdc.RegisterConcrete(MsgEditToken{}, "irishub/asset/token/MsgEditToken", nil)
	cdc.RegisterConcrete(MsgMintToken{}, "irishub/asset/token/MsgMintToken", nil)
	cdc.RegisterConcrete(MsgTransferToken{}, "irishub/asset/token/MsgTransferToken", nil)
	cdc.RegisterConcrete(MsgBurnToken{}, "irishub/asset/token/MsgBurnToken", nil)
	cdc.RegisterConcrete(FungibleToken{}, "irishub/asset/token/FungibleToken", nil)
	cdc.RegisterConcrete(&Params{}, "irishub/asset/token/Params", nil)
}

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
