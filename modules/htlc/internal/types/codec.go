package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

// RegisterCodec registers concrete types on the codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateHTLC{}, "irishub/htlc/MsgCreateHTLC", nil)
	cdc.RegisterConcrete(MsgClaimHTLC{}, "irishub/htlc/MsgClaimHTLC", nil)
	cdc.RegisterConcrete(MsgRefundHTLC{}, "irishub/htlc/MsgRefundHTLC", nil)
	cdc.RegisterConcrete(&HTLC{}, "irishub/htlc/HTLC", nil)
}

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
