package types

import "encoding/hex"

// GenesisState - all service state that must be provided at genesis
type GenesisState struct {
	Params          Params                    `json:"params"`           // service params
	RequestContexts map[string]RequestContext `json:"request_contexts"` // request contexts
	Requests        map[string]CompactRequest `json:"requests"`         // requests
	Responses       map[string]Response       `json:"responses"`        // responses
}

// NewGenesisState constructs a GenesisState
func NewGenesisState(
	params Params,
	requestContexts map[string]RequestContext,
	requests map[string]CompactRequest,
	reponses map[string]Response,
) GenesisState {
	return GenesisState{
		Params:          params,
		RequestContexts: requestContexts,
		Requests:        requests,
		Responses:       reponses,
	}
}

// DefaultGenesisState returns the default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

// get raw genesis raw message for testing
func DefaultGenesisStateForTest() GenesisState {
	return GenesisState{
		Params: DefaultParamsForTest(),
	}
}

// ValidateGenesis validates the provided service genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	err := validateParams(data.Params)
	if err != nil {
		return err
	}

	for requestContextID := range data.RequestContexts {
		if _, err := hex.DecodeString(requestContextID); err != nil {
			return err
		}
	}

	for requestID := range data.Requests {
		if _, err := ConvertRequestID(requestID); err != nil {
			return err
		}
	}

	for requestID := range data.Responses {
		if _, err := ConvertRequestID(requestID); err != nil {
			return err
		}
	}

	return nil
}
