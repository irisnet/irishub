package asset

import (
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

var (
	PrefixGateway = []byte("gateways:") // prefix for the gateway store
)

// KeyToken returns the key of the specified token source and id
func KeyToken(tokenId string) []byte {
	keyId, _ := sdk.ConvertIdToTokenKeyId(tokenId)
	return []byte(fmt.Sprintf("token:%s", keyId))
}

// KeyTokens returns the key of the specifed owner and token id. Intended for querying all tokens of an owner
func KeyTokens(owner sdk.AccAddress, tokenId string) []byte {
	if owner.Empty() {
		return []byte(fmt.Sprintf("tokens:%s", tokenId))
	}
	return []byte(fmt.Sprintf("ownerTokens:%s:%s", owner, tokenId))
}

// KeyGateway returns the key of the specified moniker
func KeyGateway(moniker string) []byte {
	return []byte(fmt.Sprintf("gateways:%s", moniker))
}

// KeyOwnerGateway returns the key of the specifed owner and moniker. Intended for querying all gateways of an owner
func KeyOwnerGateway(owner sdk.AccAddress, moniker string) []byte {
	return []byte(fmt.Sprintf("ownerGateways:%d:%s", owner, moniker))
}

// KeyGatewaysSubspace returns the key prefix for iterating on all gateways of an owner
func KeyGatewaysSubspace(owner sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("ownerGateways:%d:", owner))
}
