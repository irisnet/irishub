package types

import (
	gogotypes "github.com/gogo/protobuf/types"

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
	cdc.RegisterConcrete(&MsgCreatePool{}, "irismod/farm/MsgCreatePool", nil)
	cdc.RegisterConcrete(&MsgDestroyPool{}, "irismod/farm/MsgDestroyPool", nil)
	cdc.RegisterConcrete(&MsgAdjustPool{}, "irismod/farm/MsgAdjustPool", nil)
	cdc.RegisterConcrete(&MsgStake{}, "irismod/farm/MsgStake", nil)
	cdc.RegisterConcrete(&MsgUnstake{}, "irismod/farm/MsgUnstake", nil)
	cdc.RegisterConcrete(&MsgHarvest{}, "irismod/farm/MsgHarvest", nil)
}

// RegisterInterfaces registers the interface
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreatePool{},
		&MsgDestroyPool{},
		&MsgAdjustPool{},
		&MsgStake{},
		&MsgUnstake{},
		&MsgHarvest{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// MustUnMarshalPoolName return the poolName protobuf code
func MustMarshalPoolName(cdc codec.Codec, poolName string) []byte {
	poolNameWrap := gogotypes.StringValue{Value: poolName}
	return cdc.MustMarshal(&poolNameWrap)
}

// MustUnMarshalPoolName return the poolName
func MustUnMarshalPoolName(cdc codec.Codec, poolName []byte) string {
	var poolNameWrap gogotypes.StringValue
	cdc.MustUnmarshal(poolName, &poolNameWrap)
	return poolNameWrap.Value
}
