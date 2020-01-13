package types

// GenesisState contains all rand state that must be provided at genesis
type GenesisState struct {
	PendingRandRequests map[string][]Request `json:"pending_rand_requests"` // pending rand requests: height->[]Request
}
