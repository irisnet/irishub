package types

// GenesisState contains all HTLC state that must be provided at genesis
type GenesisState struct {
	PendingHTLCs map[string]HTLC `json:"pending_htlcs"` // claimable HTLCs
}
