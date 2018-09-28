package record

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {

	////////////////////  iris begin  ///////////////////////////
	/*cdc.RegisterConcrete(govparams.DepositProcedure{}, "cosmos-sdk/DepositProcedure", nil)
	cdc.RegisterConcrete(TallyingProcedure{}, "cosmos-sdk/TallyingProcedure", nil)
	cdc.RegisterConcrete(govparams.VotingProcedure{}, "cosmos-sdk/VotingProcedure", nil)
	cdc.RegisterConcrete(&ParameterProposal{}, "gov/ParameterProposal", nil)
	cdc.RegisterConcrete(&SoftwareUpgradeProposal{}, "gov/SoftwareUpgradeProposal", nil)*/
	////////////////////  iris end  ///////////////////////////
}

var msgCdc = wire.NewCodec()
