package keeper

import (
	"github.com/irisnet/irishub/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&UpgradeConfig{}, "iris-hub/protocol/upgradeConfig", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
