package guardian

import (
	"github.com/irisnet/irishub/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgAddProfiler{}, "iris-hub/profiling/MsgAddProfiler", nil)
	cdc.RegisterConcrete(Profiler{}, "iris-hub/profiling/Profiler", nil)
	cdc.RegisterConcrete(Trustee{}, "iris-hub/profiling/Trustee", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
