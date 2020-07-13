package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RegisterCodec registers the necessary x/bank interfaces and concrete types
// on the provided Amino codec. These types are used for Amino JSON serialization.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateFeed{}, "irishub/oracle/MsgCreateFeed", nil)
	cdc.RegisterConcrete(MsgStartFeed{}, "irishub/oracle/MsgStartFeed", nil)
	cdc.RegisterConcrete(MsgPauseFeed{}, "irishub/oracle/MsgPauseFeed", nil)
	cdc.RegisterConcrete(MsgEditFeed{}, "irishub/oracle/MsgEditFeed", nil)

	cdc.RegisterConcrete(Feed{}, "irishub/oracle/Feed", nil)
	cdc.RegisterConcrete(FeedContext{}, "irishub/oracle/FeedContext", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateFeed{},
		&MsgStartFeed{},
		&MsgPauseFeed{},
		&MsgEditFeed{},
	)
}

var (
	amino = codec.New()

	ModuleCdc = codec.NewHybridCodec(amino, types.NewInterfaceRegistry())
)

func init() {
	RegisterCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
