package types

const (
	QueryHTLC = "htlc" // HTLC query endpoint supported by the HTLC querier
)

// QueryHTLCParams is the query parameters for 'custom/htlc/htlc'
type QueryHTLCParams struct {
	HashLock HTLCHashLock `json:"hash_lock" yaml:"hash_lock"` // hash lock of an HTLC
}
