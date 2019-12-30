package types

import (
	"fmt"
)

const (
	// ModuleName is the name of the Asset module
	ModuleName = "asset"

	// StoreKey is the string store representation
	StoreKey string = ModuleName

	// QuerierRoute is the querier route for the Asset module
	QuerierRoute string = ModuleName

	// RouterKey is the msg router key for the Asset module
	RouterKey string = ModuleName

	// DefaultParamspace default name for parameter store
	DefaultParamspace = ModuleName
)

// KVStore key prefixes for Asset
const (
	//for token submodule
	KeyTokenPrefix int = iota + 1
	KeyTokenSymbolPrefix
	KeyTokenMinUnitPrefix

	//for NFT submodule
	KeyNFTCollectionPrefix
	KeyNFTIDCollectionPrefix

	//for IBC-Token submodule
	KeyIBCTokenCanonicalSymbol
)

// KeyPrefixBytes return the key prefix bytes from a URL string format
func KeyPrefixBytes(submodule string, prefix int) []byte {
	return []byte(fmt.Sprintf("%s/%d", submodule, prefix))
}
