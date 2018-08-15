package ibc

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(IBCTransferMsg{}, "cosmos-sdk/IBCTransferMsg/2", nil)
	cdc.RegisterConcrete(IBCReceiveMsg{}, "cosmos-sdk/IBCReceiveMsg/2", nil)
	cdc.RegisterConcrete(IBCSetMsg{},"cosmos-sdk/IBCSetMsg/2",nil)
	cdc.RegisterConcrete(IBCGetMsg{},"cosmos-sdk/IBCGetMsg/2",nil)
}
