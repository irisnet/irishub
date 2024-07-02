package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	gogotypes "github.com/cosmos/gogoproto/types"
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
	cdc.RegisterConcrete(&MsgCreatePool{}, "irismod/farm/MsgCreatePool", nil)
	cdc.RegisterConcrete(
		&MsgCreatePoolWithCommunityPool{},
		"irismod/farm/MsgCreatePoolWithCommunityPool",
		nil,
	)
	cdc.RegisterConcrete(&MsgDestroyPool{}, "irismod/farm/MsgDestroyPool", nil)
	cdc.RegisterConcrete(&MsgAdjustPool{}, "irismod/farm/MsgAdjustPool", nil)
	cdc.RegisterConcrete(&MsgStake{}, "irismod/farm/MsgStake", nil)
	cdc.RegisterConcrete(&MsgUnstake{}, "irismod/farm/MsgUnstake", nil)
	cdc.RegisterConcrete(&MsgHarvest{}, "irismod/farm/MsgHarvest", nil)
	cdc.RegisterConcrete(&MsgUpdateParams{}, "irismod/farm/MsgUpdateParams", nil)
}

// RegisterInterfaces registers the interface
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreatePool{},
		&MsgCreatePoolWithCommunityPool{},
		&MsgDestroyPool{},
		&MsgAdjustPool{},
		&MsgStake{},
		&MsgUnstake{},
		&MsgHarvest{},
		&MsgUpdateParams{},
	)

	registry.RegisterImplementations(
		(*govv1beta1.Content)(nil),
		&CommunityPoolCreateFarmProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// MustMarshalPoolId return the poolId protobuf code
func MustMarshalPoolId(cdc codec.Codec, poolId string) []byte {
	poolIdWrap := gogotypes.StringValue{Value: poolId}
	return cdc.MustMarshal(&poolIdWrap)
}

// MustUnMarshalPoolId return the poolId
func MustUnMarshalPoolId(cdc codec.Codec, poolId []byte) string {
	var poolIdWrap gogotypes.StringValue
	cdc.MustUnmarshal(poolId, &poolIdWrap)
	return poolIdWrap.Value
}
