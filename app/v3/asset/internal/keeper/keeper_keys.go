package keeper

import (
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

var (
	PrefixToken = []byte("token:") // prefix for the token store
)

// KeyToken returns the key of the token with the specified symbol
func KeyToken(symbol string) []byte {
	return []byte(fmt.Sprintf("token:%s", symbol))
}

// KeyTokens returns the key of the specifed owner and symbol. Intended for querying all tokens of an owner
func KeyTokens(owner sdk.AccAddress, symbol string) []byte {
	return []byte(fmt.Sprintf("ownerTokens:%s:%s", owner, symbol))
}
