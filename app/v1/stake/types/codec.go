package types

import (
	"github.com/irisnet/irishub/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateValidator{}, "irishub/stake/MsgCreateValidator", nil)
	cdc.RegisterConcrete(MsgEditValidator{}, "irishub/stake/MsgEditValidator", nil)
	cdc.RegisterConcrete(MsgDelegate{}, "irishub/stake/MsgDelegate", nil)
	cdc.RegisterConcrete(MsgBeginUnbonding{}, "irishub/stake/BeginUnbonding", nil)
	cdc.RegisterConcrete(MsgBeginRedelegate{}, "irishub/stake/BeginRedelegate", nil)

	cdc.RegisterConcrete(Pool{}, "irishub/stake/Pool", nil)
	cdc.RegisterConcrete(BondedPool{}, "irishub/stake/BondedPool", nil)
	cdc.RegisterConcrete(Validator{}, "irishub/stake/Validator", nil)
	cdc.RegisterConcrete(Delegation{}, "irishub/stake/Delegation", nil)
	cdc.RegisterConcrete(UnbondingDelegation{}, "irishub/stake/UnbondingDelegation", nil)
	cdc.RegisterConcrete(Redelegation{}, "irishub/stake/Redelegation", nil)

	cdc.RegisterConcrete(&Params{}, "irishub/stake/Params", nil)
}

// generic sealed codec to be used throughout sdk
var MsgCdc *codec.Codec

func init() {
	cdc := codec.New()
	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	MsgCdc = cdc.Seal()
}
