package profiling

import (
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/profiling/tags"
)

// handle all "profiling" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgAddProfiler:
			return handleMsgAddProfiler(ctx, k, msg)
		default:
			return sdk.ErrTxDecode("invalid message parse in profiling module").Result()
		}
	}
}

func handleMsgAddProfiler(ctx sdk.Context, k Keeper, msg MsgAddProfiler) sdk.Result {
	if _, found := k.GetProfiler(ctx, msg.AddedAddr); !found {
		return ErrProfilerNotExists(DefaultCodespace, msg.AddedAddr).Result()
	}
	if _, found := k.GetProfiler(ctx, msg.Addr); found {
		return ErrProfilerExists(DefaultCodespace, msg.Addr).Result()
	}
	err := k.AddProfiler(ctx, msg.Profiler)
	if err != nil {
		return err.Result()
	}
	resTags := sdk.NewTags(
		tags.Action, tags.ActionAddProfiler,
	)
	return sdk.Result{
		Tags: resTags,
	}
}
