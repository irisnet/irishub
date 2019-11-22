package types

const (
	QueryHTLC = "htlc"
)

// QueryHTLCParams is the query parameters for 'custom/htlc/htlc'
type QueryHTLCParams struct {
	HashLock []byte
}
