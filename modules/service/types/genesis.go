package types

import (
	"encoding/hex"
	"fmt"
)

// NewGenesisState constructs a GenesisState
func NewGenesisState(
	params Params,
	definitions []ServiceDefinition,
	bindings []ServiceBinding,
	withdrawAddresses map[string][]byte,
	requestContexts map[string]*RequestContext,
) *GenesisState {
	return &GenesisState{
		Params:            params,
		Definitions:       definitions,
		Bindings:          bindings,
		WithdrawAddresses: withdrawAddresses,
		RequestContexts:   requestContexts,
	}
}

// DefaultGenesisState gets the raw genesis raw message for testing
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// ValidateGenesis validates the provided service genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
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
