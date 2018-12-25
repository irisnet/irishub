package guardian

import (
	"github.com/irisnet/irishub/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgAddProfiler{}, "iris-hub/guardian/MsgAddProfiler", nil)
	cdc.RegisterConcrete(MsgAddTrustee{}, "iris-hub/guardian/MsgAddTrustee", nil)
	cdc.RegisterConcrete(MsgDeleteProfiler{}, "iris-hub/guardian/MsgDeleteProfiler", nil)
	cdc.RegisterConcrete(MsgDeleteTrustee{}, "iris-hub/guardian/MsgDeleteTrustee", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
