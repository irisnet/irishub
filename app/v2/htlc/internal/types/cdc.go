package types

import (
	"github.com/irisnet/irishub/codec"
)

// Register concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateHTLC{}, "irishub/htlc/MsgCreateHTLC", nil)
	cdc.RegisterConcrete(MsgClaimHTLC{}, "irishub/htlc/MsgClaimHTLC", nil)
	cdc.RegisterConcrete(MsgRefundHTLC{}, "irishub/htlc/MsgRefundHTLC", nil)

	cdc.RegisterConcrete(&HTLC{}, "irishub/htlc/HTLC", nil)
	cdc.RegisterConcrete(&Params{}, "irishub/htlc/Params", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
