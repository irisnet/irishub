package ibc

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(IBCSetMsg{},"cosmos-sdk/IBCSetMsg/1",nil)
	cdc.RegisterConcrete(IBCGetMsg{},"cosmos-sdk/IBCGetMsg/1",nil)
}
