package types

import (
	"github.com/irisnet/irishub/codec"
)

var msgCdc = codec.New()

// Register concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateHTLC{}, "irishub/htlc/MsgCreateHTLC", nil)
	cdc.RegisterConcrete(MsgClaimHTLC{}, "irishub/htlc/MsgClaimHTLC", nil)
	cdc.RegisterConcrete(MsgRefundHTLC{}, "irishub/htlc/MsgRefundHTLC", nil)
	cdc.RegisterConcrete(&HTLC{}, "irishub/htlc/HTLC", nil)
}

func init() {
	RegisterCodec(msgCdc)
}
