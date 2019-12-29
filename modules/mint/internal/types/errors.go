package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Rand errors reserve 100 ~ 199.
// nolint
const (
	CodeInvalidMintInflation sdk.CodeType = 400
	CodeInvalidMintDenom     sdk.CodeType = 401
)
