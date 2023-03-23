package v1

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*TokenI)(nil), nil)
	cdc.RegisterConcrete(&Token{}, "irismod/token/v1/Token", nil)

	cdc.RegisterConcrete(&MsgIssueToken{}, "irismod/token/v1/MsgIssueToken", nil)
	cdc.RegisterConcrete(&MsgEditToken{}, "irismod/token/v1/MsgEditToken", nil)
	cdc.RegisterConcrete(&MsgMintToken{}, "irismod/token/v1/MsgMintToken", nil)
	cdc.RegisterConcrete(&MsgBurnToken{}, "irismod/token/v1/MsgBurnToken", nil)
	cdc.RegisterConcrete(&MsgTransferTokenOwner{}, "irismod/token/v1/MsgTransferTokenOwner", nil)
	cdc.RegisterConcrete(&MsgSwapFeeToken{}, "irismod/token/v1/MsgSwapFeeToken", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgIssueToken{},
		&MsgEditToken{},
		&MsgMintToken{},
		&MsgBurnToken{},
		&MsgTransferTokenOwner{},
		&MsgSwapFeeToken{},
	)
	registry.RegisterInterface(
		"irismod.token.v1.TokenI",
		(*TokenI)(nil),
		&Token{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
