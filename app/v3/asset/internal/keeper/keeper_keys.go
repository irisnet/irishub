package keeper

import (
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

var (
	PrefixToken = []byte("token:") // prefix for the token store
)

// KeyToken returns the key of the specified token source and id
func KeyToken(tokenId string) []byte {
	keyId, _ := sdk.ConvertIdToTokenKeyId(tokenId)
	return []byte(fmt.Sprintf("token:%s", keyId))
}

// KeyTokens returns the key of the specified owner and token id. Intended for querying all tokens of an owner
func KeyTokens(owner sdk.AccAddress, tokenId string) []byte {
	if owner.Empty() {
		return []byte(fmt.Sprintf("tokens:%s", tokenId))
	}

	return []byte(fmt.Sprintf("ownerTokens:%s:%s", owner, tokenId))
}
