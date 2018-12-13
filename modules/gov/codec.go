package gov

import (
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/gov/params"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {

	cdc.RegisterConcrete(MsgSubmitProposal{}, "cosmos-sdk/MsgSubmitProposal", nil)
	cdc.RegisterConcrete(MsgSubmitTxTaxUsageProposal{}, "gov/MsgSubmitTxTaxUsageProposal", nil)
	cdc.RegisterConcrete(MsgSubmitSoftwareUpgradeProposal{}, "gov/MsgSubmitSoftwareUpgradeProposal", nil)
	cdc.RegisterConcrete(MsgDeposit{}, "cosmos-sdk/MsgDeposit", nil)
	cdc.RegisterConcrete(MsgVote{}, "cosmos-sdk/MsgVote", nil)

	cdc.RegisterInterface((*Proposal)(nil), nil)
	cdc.RegisterConcrete(&TextProposal{}, "gov/TextProposal", nil)

	////////////////////  iris begin  ///////////////////////////
	cdc.RegisterConcrete(&govparams.DepositProcedure{}, "cosmos-sdk/DepositProcedure", nil)
	cdc.RegisterConcrete(&govparams.TallyingProcedure{}, "cosmos-sdk/TallyingProcedure", nil)
	cdc.RegisterConcrete(&govparams.VotingProcedure{}, "cosmos-sdk/VotingProcedure", nil)
	cdc.RegisterConcrete(&ParameterProposal{}, "gov/ParameterProposal", nil)
	cdc.RegisterConcrete(&SoftwareUpgradeProposal{}, "gov/SoftwareUpgradeProposal", nil)
	cdc.RegisterConcrete(&HaltProposal{}, "gov/TerminatorProposal", nil)
	cdc.RegisterConcrete(&TaxUsageProposal{}, "gov/TaxUsageProposal", nil)
	////////////////////  iris end  ///////////////////////////
}

var msgCdc = codec.New()
