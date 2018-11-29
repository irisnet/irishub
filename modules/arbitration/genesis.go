package arbitration

import (
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/iparam"
	"github.com/irisnet/irishub/modules/arbitration/params"
	"time"
)

// GenesisState - all arbitration state that must be provided at genesis
type GenesisState struct {
	ComplaintRetrospect  time.Duration `json:"complaint_retrospect"`
	ArbitrationTimelimit time.Duration `json:"arbitration_timelimit"`
}

func NewGenesisState(complaintRetrospect, arbitrationTimelimit time.Duration) GenesisState {
	return GenesisState{
		ComplaintRetrospect:  complaintRetrospect,
		ArbitrationTimelimit: arbitrationTimelimit,
	}
}

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, data GenesisState) {
	iparam.InitGenesisParameter(&arbitrationparams.ComplaintRetrospectParameter, ctx, data.ComplaintRetrospect)
	iparam.InitGenesisParameter(&arbitrationparams.ArbitrationTimelimitParameter, ctx, data.ArbitrationTimelimit)
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context) GenesisState {
	complaintRetrospect := arbitrationparams.GetComplaintRetrospect(ctx)
	arbitrationTimelimit := arbitrationparams.GetArbitrationTimelimit(ctx)

	return GenesisState{
		ComplaintRetrospect:  complaintRetrospect,
		ArbitrationTimelimit: arbitrationTimelimit,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		ComplaintRetrospect:  15 * 24 * time.Hour,
		ArbitrationTimelimit: 5 * 24 * time.Hour,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisStateForTest() GenesisState {
	return GenesisState{
		ComplaintRetrospect:  20 * time.Second,
		ArbitrationTimelimit: 20 * time.Second,
	}
}
