package types

import "github.com/irisnet/irishub/codec"

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateFeed{}, "irishub/oracle/MsgCreateFeed", nil)
	cdc.RegisterConcrete(MsgStartFeed{}, "irishub/oracle/MsgStartFeed", nil)
	cdc.RegisterConcrete(MsgPauseFeed{}, "irishub/oracle/MsgPauseFeed", nil)
	cdc.RegisterConcrete(MsgEditFeed{}, "irishub/oracle/MsgEditFeed", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
