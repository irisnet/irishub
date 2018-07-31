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

	proposalID := msgSwitch.ProposalID
	CurrentProposalID := k.GetCurrentProposalID(ctx)

	if proposalID != CurrentProposalID {

		return NewError(DefaultCodespace, CodeNotCurrentProposal, "It isn't the current SoftwareUpgradeProposal").Result()

	}

	voter := msgSwitch.Voter

	if _,ok := k.sk.GetValidator(ctx,voter); !ok {
		return NewError(DefaultCodespace, CodeNotValidator, "Not a validator").Result()
	}

	if _,ok := k.GetSwitch(ctx,proposalID,voter); ok {
		return NewError(DefaultCodespace, CodeDoubleSwitch, "You have sent the switch msg").Result()
	}

	k.SetSwitch(ctx,proposalID,voter,msgSwitch)

	return sdk.Result{
		Code: 0,
		Log:  fmt.Sprintf("Switch %s by %s", msgSwitch.Title, msgSwitch.Voter.String()),
	}
}


// do switch
func EndBlocker(ctx sdk.Context, keeper Keeper) (tags sdk.Tags) {
	tags = sdk.NewTags()

    if (keeper.GetCurrentProposalID(ctx) != -1) && (ctx.BlockHeight() == keeper.GetCurrentProposalAcceptHeight(ctx) + defaultSwichPeriod) {
		passes := tally(ctx, keeper)
		if passes {
			tags.AppendTag("action", []byte("switchPassed"))

			keeper.DoSwitchBegin(ctx)
		}else{
			tags.AppendTag("action", []byte("switchDropped"))

			keeper.SetCurrentProposalID(ctx, -1)
		}
	}

	if keeper.GetDoingSwitch(ctx) && ctx.BlockHeader().GetNumTxs() == 0 {
		tags.AppendTag("action", []byte("readyToDoSwitch"))

		keeper.DoSwitchEnd(ctx)
	}

	return tags
}