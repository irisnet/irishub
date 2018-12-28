package gov

import (
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/types/gov/params"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {

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
