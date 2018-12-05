package auth

import (
	sdk "github.com/irisnet/irishub/types"
	banktype "github.com/irisnet/irishub/types/bank"
)

type contextKey int // local to the auth module

const (
	contextKeySigners contextKey = iota
)

// add the signers to the context
func WithSigners(ctx sdk.Context, accounts []banktype.Account) sdk.Context {
	return ctx.WithValue(contextKeySigners, accounts)
}

// get the signers from the context
func GetSigners(ctx sdk.Context) []banktype.Account {
	v := ctx.Value(contextKeySigners)
	if v == nil {
		return []banktype.Account{}
	}
	return v.([]banktype.Account)
}

