package upgrade

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgSwitch{}, "iris-hub/upgrade/MsgSwitch", nil)
	cdc.RegisterConcrete(&ModuleLifeTime{}, "iris-hub/upgrade/ModuleLifeTime", nil)
	cdc.RegisterConcrete(&Version{}, "iris-hub/upgrade/Version", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	RegisterWire(msgCdc)
}
