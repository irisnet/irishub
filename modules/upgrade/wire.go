package upgrade

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgSwitch{}, "iris-hub/MsgSwitch", nil)
	cdc.RegisterConcrete(&ModuleLifeTime{}, "iris-hub/", nil)
	cdc.RegisterConcrete(&Version{}, "iris-hub/", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	RegisterWire(msgCdc)
}
