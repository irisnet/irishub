package upgrade

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"reflect"
	"fmt"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSwitch:
			return handlerSwitch(ctx, msg, k)
		default:
			errMsg := "Unrecognized Upgrade Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handlerSwitch(ctx sdk.Context, msg sdk.Msg, k Keeper) sdk.Result {
	msgSwitch, ok := msg.(MsgSwitch)
	if !ok {
		return NewError(DefaultCodespace, CodeInvalidMsgType, "Handler should only receive MsgSwitch").Result()
	}

	return sdk.Result{
		Code: 0,
		Log:  fmt.Sprintf("Switch %s by %s", msgSwitch.Title, msgSwitch.Voter.String()),
	}
}
