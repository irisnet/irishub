package types

// nolint
const (
	// module name
	ModuleName = "htlc"

	// RouterKey is the message route for the htlc module.
	RouterKey = ModuleName

	// StoreKey is the default store key for the htlc module.
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the htlc module.
	QuerierRoute = StoreKey
)
