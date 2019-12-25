package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

	DefaultCodespace sdk.CodespaceType = ModuleName

	// DefaultParamspace default name for parameter store
	DefaultParamspace = ModuleName
)

// KVStore key prefixes for Asset
const (
	KeyTokenPrefix int = iota + 1
	KeyTokenSymbolPrefix
	KeyTokenMinUnitPrefix

	KeyNFTCollectionPrefix
	KeyNFTIDCollectionPrefix
)

// KeyPrefixBytes return the key prefix bytes from a URL string format
func KeyPrefixBytes(prefix int) []byte {
	return []byte(fmt.Sprintf("%d/", prefix))
}
