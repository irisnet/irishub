package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/irisnet/irishub/modules/service/internal/types"
)

// Simulation parameter constants
const (
	MaxRequestTimeout    = "max_request_timeout"
	MinDepositMultiple   = "min_deposit_multiple"
	ServiceFeeTax        = "service_fee_tax"
	SlashFraction        = "slash_fraction"
	ComplaintRetrospect  = "complaint_retrospect"
	ArbitrationTimeLimit = "arbitration_time_limit"
	TxSizeLimit          = "tx_size_limit"
)

// GenMaxRequestTimeout returns the randomized MaxRequestTimeout
func GenMaxRequestTimeout(r *rand.Rand) int64 {
	return int64(simulation.RandIntBetween(r, 20, 50))
}

// GenMinDepositMultiple returns the randomized MinDepositMultiple
func GenMinDepositMultiple(r *rand.Rand) int64 {
	return int64(simulation.RandIntBetween(r, 500, 5000))
}

// GenServiceFeeTax returns the randomized ServiceFeeTax
func GenServiceFeeTax(r *rand.Rand) sdk.Dec {
	return simulation.RandomDecAmount(r, sdk.NewDecWithPrec(2, 1))
}

// GenSlashFraction returns the randomized SlashFraction
func GenSlashFraction(r *rand.Rand) sdk.Dec {
	return simulation.RandomDecAmount(r, sdk.NewDecWithPrec(1, 2))
}

// GenComplaintRetrospect returns the randomized ComplaintRetrospect
func GenComplaintRetrospect(r *rand.Rand) time.Duration {
	return time.Duration(simulation.RandIntBetween(r, 15, 30)) * 24 * time.Hour
}

// GenArbitrationTimeLimit returns the randomized ArbitrationTimeLimit
func GenArbitrationTimeLimit(r *rand.Rand) time.Duration {
	return time.Duration(simulation.RandIntBetween(r, 5, 10)) * 24 * time.Hour
}

// GenTxSizeLimit returns the randomized TxSizeLimit
func GenTxSizeLimit(r *rand.Rand) uint64 {
	return uint64(simulation.RandIntBetween(r, 2000, 6000))
}

// RandomizedGenState generates a random GenesisState for service
func RandomizedGenState(simState *module.SimulationState) {
	var maxRequestTimeout int64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, MaxRequestTimeout, &maxRequestTimeout, simState.Rand,
		func(r *rand.Rand) { maxRequestTimeout = GenMaxRequestTimeout(r) },
	)

	var minDepositMultiple int64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, MinDepositMultiple, &minDepositMultiple, simState.Rand,
		func(r *rand.Rand) { minDepositMultiple = GenMinDepositMultiple(r) },
	)

	var serviceFeeTax sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, ServiceFeeTax, &serviceFeeTax, simState.Rand,
		func(r *rand.Rand) { serviceFeeTax = GenServiceFeeTax(r) },
	)

	var slashFraction sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, SlashFraction, &slashFraction, simState.Rand,
		func(r *rand.Rand) { slashFraction = GenSlashFraction(r) },
	)

	var complaintRetrospect time.Duration
	simState.AppParams.GetOrGenerate(
		simState.Cdc, ComplaintRetrospect, &complaintRetrospect, simState.Rand,
		func(r *rand.Rand) { complaintRetrospect = GenComplaintRetrospect(r) },
	)

	var arbitrationTimeLimit time.Duration
	simState.AppParams.GetOrGenerate(
		simState.Cdc, ArbitrationTimeLimit, &arbitrationTimeLimit, simState.Rand,
		func(r *rand.Rand) { arbitrationTimeLimit = GenArbitrationTimeLimit(r) },
	)

	var txSizeLimit uint64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, TxSizeLimit, &txSizeLimit, simState.Rand,
		func(r *rand.Rand) { txSizeLimit = GenTxSizeLimit(r) },
	)

	params := types.NewParams(
		maxRequestTimeout, minDepositMultiple, serviceFeeTax,
		slashFraction, complaintRetrospect, arbitrationTimeLimit, txSizeLimit,
	)

	serviceGenesis := types.NewGenesisState(params)

	fmt.Printf("Selected randomly generated service parameters:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, serviceGenesis.Params))
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(serviceGenesis)
}
