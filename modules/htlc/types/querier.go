package types

import (
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

const (
	QueryHTLC          = "htlc"          // query an HTLC
	QueryAssetSupply   = "assetSupply"   // query an asset supply
	QueryAssetSupplies = "assetSupplies" // query all asset supplies
	QueryParameters    = "parameters"    // query parameters
)

// QueryHTLCParams defines the params to query an HTLC
type QueryHTLCParams struct {
	ID tmbytes.HexBytes
}

type QueryAssetSupplyParams struct {
	Denom string
}
