package types

// nolint
const (
	// module name
	ModuleName = "record"

	// StoreKey is the default store key for record
	StoreKey = ModuleName

	// RouterKey is the message route for record
	RouterKey = ModuleName

	// QuerierRoute is the querier route for the record store.
	QuerierRoute = StoreKey

	// Query endpoints supported by the record querier
	QueryRecord = "record"
)

var (
	RecordKey         = []byte{0x01} // record key
	IntraTxCounterKey = []byte{0x02} // key for intra-block tx index
)

// GetRecordKey returns record key bytes
func GetRecordKey(recordID []byte) []byte {
	return append(RecordKey, recordID...)
}
