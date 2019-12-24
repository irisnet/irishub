package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

// RegisterCodec registers concrete types on the codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgIssueToken{}, "irishub/asset/MsgIssueToken", nil)
	cdc.RegisterConcrete(MsgEditToken{}, "irishub/asset/MsgEditToken", nil)
	cdc.RegisterConcrete(MsgMintToken{}, "irishub/asset/MsgMintToken", nil)
	cdc.RegisterConcrete(MsgTransferTokenOwner{}, "irishub/asset/MsgTransferTokenOwner", nil)
	cdc.RegisterConcrete(BaseToken{}, "irishub/asset/BaseToken", nil)
	cdc.RegisterConcrete(FungibleToken{}, "irishub/asset/FungibleToken", nil)
	cdc.RegisterConcrete(&Params{}, "irishub/asset/Params", nil)
}

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
