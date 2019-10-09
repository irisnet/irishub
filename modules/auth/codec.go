package auth

import (
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

// Register concrete types on codec codec for default AppAccount
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*Account)(nil), nil)
	cdc.RegisterConcrete(&sdk.BaseAccount{}, "irishub/bank/Account", nil)
	cdc.RegisterConcrete(StdTx{}, "irishub/bank/StdTx", nil)
	cdc.RegisterConcrete(&Params{}, "irishub/Auth/Params", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
	codec.RegisterCrypto(msgCdc)
}
