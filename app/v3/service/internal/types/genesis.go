package types

import (
	"encoding/hex"

	cmn "github.com/tendermint/tendermint/libs/common"

	sdk "github.com/irisnet/irishub/types"
)

// GenesisState - all service state that must be provided at genesis
type GenesisState struct {
	Params              Params                      `json:"params"`                // service params
	Definitions         []ServiceDefinition         `json:"definetions"`           // service definitions
	Bindings            map[string][]ServiceBinding `json:"bindings"`              // service bindings
	WithdrawAddresses   []sdk.AccAddress            `json:"withdraw_addresses"`    // withdraw addresses
	RequestContexts     map[string]RequestContext   `json:"request_contexts"`      // request contexts
	NewRequestBatch     map[string][]cmn.HexBytes   `json:"new_request_batch"`     // new request batch
	ExpiredRequestBatch map[string][]cmn.HexBytes   `json:"expired_request_batch"` // expired request batch
	Requests            map[string]CompactRequest   `json:"requests"`              // requests
	Responses           map[string]Response         `json:"responses"`             // responses
}

// NewGenesisState constructs a GenesisState
func NewGenesisState(
	params Params,
	definitions []ServiceDefinition,
	bindings map[string][]ServiceBinding,
	withdrawAddresses []sdk.AccAddress,
	requestContexts map[string]RequestContext,
	newRequestBatch map[string][]cmn.HexBytes,
	expiredRequestBatch map[string][]cmn.HexBytes,
	requests map[string]CompactRequest,
	reponses map[string]Response,
) GenesisState {
	return GenesisState{
		Params:              params,
		Definitions:         definitions,
		Bindings:            bindings,
		WithdrawAddresses:   withdrawAddresses,
		RequestContexts:     requestContexts,
		NewRequestBatch:     newRequestBatch,
		ExpiredRequestBatch: expiredRequestBatch,
		Requests:            requests,
		Responses:           reponses,
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
	if err := validateParams(data.Params); err != nil {
		return err
	}

	for _, _ = range data.Definitions {
		// TODO: validate Definitions
	}

	for _, _ = range data.Bindings {
		// TODO: validate Bindings
	}

	for _, _ = range data.WithdrawAddresses {
		// TODO: validate WithdrawAddresses
	}

	for requestContextID := range data.RequestContexts {
		if _, err := hex.DecodeString(requestContextID); err != nil {
			return err
		}
	}

	for _, _ = range data.NewRequestBatch {
		// TODO: validate NewRequestBatch
	}

	for _, _ = range data.ExpiredRequestBatch {
		// TODO: validate ExpiredRequestBatch
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
