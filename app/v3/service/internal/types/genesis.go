package types

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

// GenesisState - all service state that must be provided at genesis
type GenesisState struct {
	Params            Params                    `json:"params"`             // service params
	Definitions       []ServiceDefinition       `json:"definitions"`        // service definitions
	Bindings          []ServiceBinding          `json:"bindings"`           // service bindings
	WithdrawAddresses map[string]sdk.AccAddress `json:"withdraw_addresses"` // withdrawal addresses
	RequestContexts   map[string]RequestContext `json:"request_contexts"`   // request contexts
}

// NewGenesisState constructs a GenesisState
func NewGenesisState(
	params Params,
	definitions []ServiceDefinition,
	bindings []ServiceBinding,
	withdrawAddresses map[string]sdk.AccAddress,
	requestContexts map[string]RequestContext,
) GenesisState {
	return GenesisState{
		Params:            params,
		Definitions:       definitions,
		Bindings:          bindings,
		WithdrawAddresses: withdrawAddresses,
		RequestContexts:   requestContexts,
	}
}

// DefaultGenesisState returns the default genesis state
func DefaultGenesisState(moduleSvcDefinitions []ServiceDefinition) GenesisState {
	return GenesisState{
		Params:      DefaultParams(),
		Definitions: moduleSvcDefinitions,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisStateForTest(moduleSvcDefinitions []ServiceDefinition) GenesisState {
	return GenesisState{
		Params:      DefaultParamsForTest(),
		Definitions: moduleSvcDefinitions,
	}
}

// ValidateGenesis validates the provided service genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	if err := validateParams(data.Params); err != nil {
		return err
	}

	for _, definition := range data.Definitions {
		if err := definition.Validate(); err != nil {
			return err
		}
	}

	for _, binding := range data.Bindings {
		if err := binding.Validate(); err != nil {
			return err
		}
	}

	for providerAddressStr := range data.WithdrawAddresses {
		if _, err := hex.DecodeString(providerAddressStr); err != nil {
			return err
		}
	}

	for requestContextID, requestContext := range data.RequestContexts {
		if _, err := hex.DecodeString(requestContextID); err != nil {
			return err
		}
		if err := requestContext.Validate(); err != nil {
			return err
		}
		if requestContext.State != PAUSED {
			return fmt.Errorf("invalid request context state, ID:%s, State:%s", requestContextID, requestContext.State)
		}
		if requestContext.BatchState != BATCHCOMPLETED {
			return fmt.Errorf("invalid request context batch state, ID:%s, BatchState:%s", requestContextID, requestContext.BatchState)
		}
	}

	return nil
}
