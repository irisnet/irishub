package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	iristypes "github.com/irisnet/irishub/types"
)

var (
	PrefixToken = []byte("token:") // prefix for the token store
)

// KeyToken returns the key of the specified token source and id
func KeyToken(tokenID string) []byte {
	keyID, _ := iristypes.ConvertIDToTokenKeyID(tokenID)
	return []byte(fmt.Sprintf("token:%s", keyID))
}

// KeyTokens returns the key of the specified owner and token id. Intended for querying all tokens of an owner
func KeyTokens(owner sdk.AccAddress, tokenID string) []byte {
	if owner.Empty() {
		return []byte(fmt.Sprintf("tokens:%s", tokenID))
	}
	return []byte(fmt.Sprintf("ownerTokens:%s:%s", owner, tokenID))
}
