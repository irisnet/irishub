package types

import (
	"github.com/irisnet/irishub/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgWithdrawDelegatorRewardsAll{}, "irishub/distr/MsgWithdrawDelegationRewardsAll", nil)
	cdc.RegisterConcrete(MsgWithdrawDelegatorReward{}, "irishub/distr/MsgWithdrawDelegationReward", nil)
	cdc.RegisterConcrete(MsgWithdrawValidatorRewardsAll{}, "irishub/distr/MsgWithdrawValidatorRewardsAll", nil)
	cdc.RegisterConcrete(MsgSetWithdrawAddress{}, "irishub/distr/MsgModifyWithdrawAddress", nil)

	cdc.RegisterConcrete(DelegationDistInfo{}, "irishub/distr/DelegationDistInfo", nil)
	cdc.RegisterConcrete(FeePool{}, "irishub/distr/FeePool", nil)

	cdc.RegisterConcrete(&Params{}, "irishub/distr/Params", nil)
}

// generic sealed codec to be used throughout module
var MsgCdc *codec.Codec

func init() {
	cdc := codec.New()
	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	MsgCdc = cdc.Seal()
}
