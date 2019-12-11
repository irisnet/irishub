package types

import (
	"encoding/hex"
)

// GenesisState contains all HTLC state that must be provided at genesis
type GenesisState struct {
	PendingHTLCs map[string]HTLC // claimable HTLCs
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
