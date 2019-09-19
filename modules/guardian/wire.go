package guardian

import (
	"github.com/irisnet/irishub/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgAddProfiler{}, "irishub/guardian/MsgAddProfiler", nil)
	cdc.RegisterConcrete(MsgAddTrustee{}, "irishub/guardian/MsgAddTrustee", nil)
	cdc.RegisterConcrete(MsgDeleteProfiler{}, "irishub/guardian/MsgDeleteProfiler", nil)
	cdc.RegisterConcrete(MsgDeleteTrustee{}, "irishub/guardian/MsgDeleteTrustee", nil)

	cdc.RegisterConcrete(Guardian{}, "irishub/guardian/Guardian", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
