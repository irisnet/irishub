package types

const (
	// QueryHTLC - HTLC query endpoint supported by the HTLC querier
	QueryHTLC = "htlc"
)

// QueryHTLCParams is the query parameters for 'custom/htlc/htlc'
type QueryHTLCParams struct {
	HashLock HTLCHashLock
}
