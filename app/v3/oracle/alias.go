package oracle

import (
	"github.com/irisnet/irishub/app/v3/oracle/internal/keeper"
	"github.com/irisnet/irishub/app/v3/oracle/internal/types"
)

// nolint
type (
	Keeper = keeper.Keeper

	MsgCreateFeed = types.MsgCreateFeed
	MsgStartFeed  = types.MsgStartFeed
	MsgPauseFeed  = types.MsgPauseFeed
	MsgEditFeed   = types.MsgEditFeed
)
