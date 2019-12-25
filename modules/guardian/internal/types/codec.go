package types

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/irisnet/irishub/modules/guardian/exported"
)

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

// RegisterCodec registers concrete types on the codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*exported.GuardianI)(nil), nil)
	cdc.RegisterInterface((*exported.AccountTypeI)(nil), nil)
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
