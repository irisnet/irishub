package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Handler defines the core of the state transition function of an application.
type Handler func(ctx sdk.Context, msg Msg) Result

// AnteHandler authenticates transactions, before their internal messages are handled.
// If newCtx.IsZero(), ctx is used instead.
type AnteHandler func(ctx sdk.Context, tx Tx) (newCtx sdk.Context, result Result, abort bool)
type FeeRefundHandler func(ctx sdk.Context, tx Tx, result Result) (Result, error)
