package types

import (
	"strconv"
)

// GenesisState contains all rand state that must be provided at genesis
type GenesisState struct {
	PendingRandRequests map[string][]Request // pending rand requests: height->[]Request
}

// ValidateGenesis validates the given rand genesis state
func ValidateGenesis(data GenesisState) error {
	for height, requests := range data.PendingRandRequests {
		if _, err := strconv.ParseUint(height, 10, 64); err != nil {
			return err
		}
		for _, request := range requests {
			if err := ValidateConsumer(request.Consumer); err != nil {
				return err
			}

			// check request context exists
			if request.Oracle {
				if err := ValidateServiceFeeCap(request.ServiceFeeCap); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
