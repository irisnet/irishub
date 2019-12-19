package types

import (
	"encoding/hex"
)

// GenesisState contains all HTLC state that must be provided at genesis
type GenesisState struct {
	PendingHTLCs map[string]HTLC `json:"pending_htlcs" yaml:"pending_htlcs"` // claimable HTLCs
}

// DefaultGenesisState gets the default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		PendingHTLCs: map[string]HTLC{},
	}
}

// ValidateGenesis checks if parameters are within valid ranges
func ValidateGenesis(data GenesisState) error {
	for hashLockStr, htlc := range data.PendingHTLCs {
		hashLock, err := hex.DecodeString(hashLockStr)
		if err != nil {
			return err
		}
		if err := htlc.Validate(HTLCHashLock(hashLock)); err != nil {
			return err
		}
	}
	return nil
}
