package types

// DONTCOVER

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	proto "github.com/cosmos/gogoproto/proto"

	"github.com/irisnet/irismod/nft/exported"
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

// RegisterLegacyAminoCodec concrete types on codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgIssueDenom{}, "irismod/nft/MsgIssueDenom", nil)
	cdc.RegisterConcrete(&MsgTransferNFT{}, "irismod/nft/MsgTransferNFT", nil)
	cdc.RegisterConcrete(&MsgEditNFT{}, "irismod/nft/MsgEditNFT", nil)
	cdc.RegisterConcrete(&MsgMintNFT{}, "irismod/nft/MsgMintNFT", nil)
	cdc.RegisterConcrete(&MsgBurnNFT{}, "irismod/nft/MsgBurnNFT", nil)
	cdc.RegisterConcrete(&MsgTransferDenom{}, "irismod/nft/MsgTransferDenom", nil)

	cdc.RegisterInterface((*exported.NFT)(nil), nil)
	cdc.RegisterConcrete(&BaseNFT{}, "irismod/nft/BaseNFT", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgIssueDenom{},
		&MsgTransferNFT{},
		&MsgEditNFT{},
		&MsgMintNFT{},
		&MsgBurnNFT{},
		&MsgTransferDenom{},
	)

	registry.RegisterImplementations(
		(*exported.NFT)(nil),
		&BaseNFT{},
	)

	registry.RegisterImplementations(
		(*proto.Message)(nil),
		&DenomMetadata{},
		&NFTMetadata{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
