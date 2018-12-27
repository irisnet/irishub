package gov

import (
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/gov/params"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {

	cdc.RegisterConcrete(MsgSubmitProposal{}, "irishub/gov/MsgSubmitProposal", nil)
	cdc.RegisterConcrete(MsgSubmitTxTaxUsageProposal{}, "irishub/gov/MsgSubmitTxTaxUsageProposal", nil)
	cdc.RegisterConcrete(MsgSubmitSoftwareUpgradeProposal{}, "irishub/gov/MsgSubmitSoftwareUpgradeProposal", nil)
	cdc.RegisterConcrete(MsgDeposit{}, "irishub/gov/MsgDeposit", nil)
	cdc.RegisterConcrete(MsgVote{}, "irishub/gov/MsgVote", nil)
	cdc.RegisterInterface((*Proposal)(nil), nil)
	cdc.RegisterConcrete(&TextProposal{}, "irishub/gov/TextProposal", nil)
	cdc.RegisterConcrete(&govparams.DepositProcedure{}, "irishub/gov/DepositProcedure", nil)
	cdc.RegisterConcrete(&govparams.TallyingProcedure{}, "irishub/gov/TallyingProcedure", nil)
	cdc.RegisterConcrete(&govparams.VotingProcedure{}, "irishub/gov/VotingProcedure", nil)
	cdc.RegisterConcrete(&ParameterProposal{}, "irishub/gov/ParameterProposal", nil)
	cdc.RegisterConcrete(&SoftwareUpgradeProposal{}, "irishub/gov/SoftwareUpgradeProposal", nil)
	cdc.RegisterConcrete(&SystemHaltProposal{}, "irishub/gov/SystemHaltProposal", nil)
	cdc.RegisterConcrete(&TaxUsageProposal{}, "irishub/gov/TaxUsageProposal", nil)
}

var msgCdc = codec.New()
