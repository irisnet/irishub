package upgrade

import (
	"github.com/irisnet/irishub/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSwitch{}, "iris-hub/upgrade/MsgSwitch", nil)
	cdc.RegisterConcrete(&ModuleLifeTime{}, "iris-hub/upgrade/ModuleLifeTime", nil)
	cdc.RegisterConcrete(&Version{}, "iris-hub/upgrade/Version", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
