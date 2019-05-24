package upgrade

import (
	"github.com/irisnet/irishub/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&VersionInfo{}, "irishub/upgrade/VersionInfo", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
