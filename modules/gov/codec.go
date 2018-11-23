package gov

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/irisnet/irishub/modules/gov/params"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {

	cdc.RegisterConcrete(MsgSubmitProposal{}, "cosmos-sdk/MsgSubmitProposal", nil)
	cdc.RegisterConcrete(MsgDeposit{}, "cosmos-sdk/MsgDeposit", nil)
	cdc.RegisterConcrete(MsgVote{}, "cosmos-sdk/MsgVote", nil)

	cdc.RegisterInterface((*Proposal)(nil), nil)
	cdc.RegisterConcrete(&TextProposal{}, "gov/TextProposal", nil)

	////////////////////  iris begin  ///////////////////////////
	cdc.RegisterConcrete(&govparams.DepositProcedure{},"cosmos-sdk/DepositProcedure",nil)
	cdc.RegisterConcrete(&govparams.TallyingProcedure{},"cosmos-sdk/TallyingProcedure",nil)
	cdc.RegisterConcrete(&govparams.VotingProcedure{},"cosmos-sdk/VotingProcedure",nil)
	cdc.RegisterConcrete(&ParameterProposal{}, "gov/ParameterProposal", nil)
	cdc.RegisterConcrete(&SoftwareUpgradeProposal{}, "gov/SoftwareUpgradeProposal", nil)
	////////////////////  iris end  ///////////////////////////
}

var msgCdc = codec.New()
