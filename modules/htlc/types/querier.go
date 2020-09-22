package types

import (
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

const (
	QueryHTLC = "htlc" // query an HTLC
)

// QueryHTLCParams defines the params to query an HTLC
type QueryHTLCParams struct {
	HashLock tmbytes.HexBytes
}
