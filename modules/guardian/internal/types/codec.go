package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// module codec
var ModuleCdc *codec.Codec

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgAddProfiler{}, "irishub/guardian/MsgAddProfiler", nil)
	cdc.RegisterConcrete(MsgAddTrustee{}, "irishub/guardian/MsgAddTrustee", nil)
	cdc.RegisterConcrete(MsgDeleteProfiler{}, "irishub/guardian/MsgDeleteProfiler", nil)
	cdc.RegisterConcrete(MsgDeleteTrustee{}, "irishub/guardian/MsgDeleteTrustee", nil)
	cdc.RegisterConcrete(Guardian{}, "irishub/guardian/Guardian", nil)
}

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
