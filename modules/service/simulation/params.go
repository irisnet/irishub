package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/irisnet/irishub/modules/service/internal/types"
)

const (
	keyMaxRequestTimeout    = "MaxRequestTimeout"
	keyMinDepositMultiple   = "MinDepositMultiple"
	keyServiceFeeTax        = "ServiceFeeTax"
	keySlashFraction        = "SlashFraction"
	keyComplaintRetrospect  = "ComplaintRetrospect"
	keyArbitrationTimeLimit = "ArbitrationTimeLimit"
	keyTxSizeLimit          = "TxSizeLimit"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simulation.ParamChange {
	return []simulation.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, keyMaxRequestTimeout,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenMaxRequestTimeout(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, keyMinDepositMultiple,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenMinDepositMultiple(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, keyServiceFeeTax,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenServiceFeeTax(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, keySlashFraction,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenSlashFraction(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, keyComplaintRetrospect,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenComplaintRetrospect(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, keyArbitrationTimeLimit,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenArbitrationTimeLimit(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, keyTxSizeLimit,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenTxSizeLimit(r))
			},
		),
	}
}
