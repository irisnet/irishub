package keeper

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/tidwall/gjson"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/irisnet/irismod/modules/oracle/types"
	serviceexported "github.com/irisnet/irismod/modules/service/exported"
)

// Keeper defines a struct for the oracle keeper
type Keeper struct {
	cdc        codec.Marshaler
	storeKey   sdk.StoreKey
	sk         types.ServiceKeeper
	paramSpace paramtypes.Subspace
}

// NewKeeper returns an instance of the oracle Keeper
func NewKeeper(
	cdc codec.Marshaler,
	storeKey sdk.StoreKey,
	paramSpace paramtypes.Subspace,
	sk types.ServiceKeeper,
) Keeper {
	keeper := Keeper{
		storeKey:   storeKey,
		cdc:        cdc,
		sk:         sk,
		paramSpace: paramSpace,
	}

	_ = sk.RegisterResponseCallback(types.ModuleName, keeper.HandlerResponse)
	_ = sk.RegisterStateCallback(types.ModuleName, keeper.HandlerStateChanged)
	_ = sk.RegisterModuleService(
		serviceexported.RegisterModuleName,
		&serviceexported.ModuleService{
			ServiceName:     serviceexported.OraclePriceServiceName,
			Provider:        serviceexported.OraclePriceServiceProvider,
			ReuquestService: keeper.ModuleServiceRequest,
		},
	)

	return keeper
}

// CreateFeed creates a stopped feed
func (k Keeper) CreateFeed(ctx sdk.Context, msg *types.MsgCreateFeed) error {
	if _, found := k.GetFeed(ctx, msg.FeedName); found {
		return sdkerrors.Wrapf(types.ErrExistedFeedName, msg.FeedName)
	}

	providers := make([]sdk.AccAddress, len(msg.Providers))
	for i, provider := range msg.Providers {
		pd, _ := sdk.AccAddressFromBech32(provider)
		providers[i] = pd
	}

	creator, _ := sdk.AccAddressFromBech32(msg.Creator)
	requestContextID, err := k.sk.CreateRequestContext(
		ctx,
		msg.ServiceName,
		providers,
		creator,
		msg.Input,
		msg.ServiceFeeCap,
		msg.Timeout,
		true,
		msg.RepeatedFrequency,
		-1,
		serviceexported.PAUSED,
		msg.ResponseThreshold,
		types.ModuleName,
	)
	if err != nil {
		return err
	}

	k.SetFeed(ctx, types.Feed{
		FeedName:         msg.FeedName,
		AggregateFunc:    msg.AggregateFunc,
		ValueJsonPath:    msg.ValueJsonPath,
		LatestHistory:    msg.LatestHistory,
		RequestContextID: requestContextID.String(),
		Description:      msg.Description,
		Creator:          msg.Creator,
	})
	k.Enqueue(ctx, msg.FeedName, serviceexported.PAUSED)

	return nil
}

// StartFeed starts a stopped feed
func (k Keeper) StartFeed(ctx sdk.Context, msg *types.MsgStartFeed) error {
	feed, found := k.GetFeed(ctx, msg.FeedName)
	if !found {
		return sdkerrors.Wrapf(types.ErrUnknownFeedName, msg.FeedName)
	}

	requestContextID, _ := hex.DecodeString(feed.RequestContextID)
	creator, _ := sdk.AccAddressFromBech32(msg.Creator)

	if msg.Creator != feed.Creator {
		return sdkerrors.Wrapf(types.ErrUnauthorized, msg.Creator)
	}

	reqCtx, existed := k.sk.GetRequestContext(ctx, requestContextID)
	if !existed {
		return sdkerrors.Wrapf(types.ErrUnknownFeedName, msg.FeedName)
	}

	// Can not start feed in "running" state
	if reqCtx.State == serviceexported.RUNNING {
		return sdkerrors.Wrapf(types.ErrInvalidFeedState, msg.FeedName)
	}

	if err := k.sk.StartRequestContext(ctx, requestContextID, creator); err != nil {
		return err
	}

	k.dequeueAndEnqueue(ctx, msg.FeedName, serviceexported.PAUSED, serviceexported.RUNNING)
	return nil
}

// PauseFeed pauses a running feed
func (k Keeper) PauseFeed(ctx sdk.Context, msg *types.MsgPauseFeed) error {
	feed, found := k.GetFeed(ctx, msg.FeedName)
	if !found {
		return sdkerrors.Wrapf(types.ErrUnknownFeedName, msg.FeedName)
	}

	requestContextID, _ := hex.DecodeString(feed.RequestContextID)
	creator, _ := sdk.AccAddressFromBech32(msg.Creator)

	if msg.Creator != feed.Creator {
		return sdkerrors.Wrapf(types.ErrUnauthorized, msg.Creator)
	}

	reqCtx, existed := k.sk.GetRequestContext(ctx, requestContextID)
	if !existed {
		return sdkerrors.Wrapf(types.ErrUnknownFeedName, msg.FeedName)
	}

	// Can only pause feed in "running" state
	if reqCtx.State != serviceexported.RUNNING {
		return sdkerrors.Wrapf(types.ErrInvalidFeedState, msg.FeedName)
	}

	if err := k.sk.PauseRequestContext(ctx, requestContextID, creator); err != nil {
		return err
	}

	k.dequeueAndEnqueue(ctx, msg.FeedName, serviceexported.RUNNING, serviceexported.PAUSED)
	return nil
}

// EditFeed edits a feed
func (k Keeper) EditFeed(ctx sdk.Context, msg *types.MsgEditFeed) error {
	feed, found := k.GetFeed(ctx, msg.FeedName)
	if !found {
		return sdkerrors.Wrapf(types.ErrUnknownFeedName, msg.FeedName)
	}

	if msg.Creator != feed.Creator {
		return sdkerrors.Wrapf(types.ErrUnauthorized, msg.Creator)
	}

	requestContextID, _ := hex.DecodeString(feed.RequestContextID)
	creator, _ := sdk.AccAddressFromBech32(msg.Creator)

	providers := make([]sdk.AccAddress, len(msg.Providers))
	for i, provider := range msg.Providers {
		pd, _ := sdk.AccAddressFromBech32(provider)
		providers[i] = pd
	}

	if err := k.sk.UpdateRequestContext(
		ctx,
		requestContextID,
		providers,
		msg.ResponseThreshold,
		msg.ServiceFeeCap,
		msg.Timeout,
		msg.RepeatedFrequency,
		-1,
		creator,
	); err != nil {
		return err
	}

	if msg.LatestHistory > 0 {
		cnt := k.getFeedValuesCnt(ctx, feed.FeedName)
		if expectCnt := int(msg.LatestHistory); expectCnt < cnt {
			k.deleteOldestFeedValue(ctx, feed.FeedName, cnt-expectCnt)
		}
		feed.LatestHistory = msg.LatestHistory
	}

	if types.Modified(msg.Description) {
		feed.Description = msg.Description
	}

	k.SetFeed(ctx, feed)
	return nil
}

// HandlerResponse is responsible for processing the data returned from the service module,
// processed by the aggregation function, and then saved
func (k Keeper) HandlerResponse(
	ctx sdk.Context,
	requestContextID tmbytes.HexBytes,
	responseOutput []string,
	err error,
) {
	if len(responseOutput) == 0 || err != nil {
		ctx.Logger().Info(
			"Oracle feed failed",
			"requestContextID", requestContextID.String(),
			"err", err.Error(),
		)
		return
	}

	feed, found := k.GetFeedByReqCtxID(ctx, requestContextID)
	if !found {
		ctx.Logger().Error(
			"Not existed requestContext", "requestContextID", requestContextID.String(),
		)
		return
	}

	reqCtx, existed := k.sk.GetRequestContext(ctx, requestContextID)
	if !existed {
		ctx.Logger().Error(
			"Not existed requestContext", "requestContextID", requestContextID.String(),
		)
		return
	}

	aggregate, err := types.GetAggregateFunc(feed.AggregateFunc)
	if err != nil {
		ctx.Logger().Error(
			"Not existed aggregateFunc", "aggregateFunc", feed.AggregateFunc,
		)
		return
	}

	var data []types.ArgsType
	for _, jsonStr := range responseOutput {
		result := gjson.Get(jsonStr, serviceexported.PATH_BODY).Get(feed.ValueJsonPath)
		data = append(data, result)
	}

	result := aggregate(data)
	value := types.FeedValue{
		Data:      result,
		Timestamp: ctx.BlockTime(),
	}
	k.SetFeedValue(ctx, feed.FeedName, reqCtx.BatchCounter, feed.LatestHistory, value)

	bz, _ := json.Marshal(value)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSetFeedValue,
			sdk.NewAttribute(types.AttributeKeyFeedName, feed.FeedName),
			sdk.NewAttribute(types.AttributeKeyFeedValue, string(bz)),
		),
	)
}

// HandlerStateChanged is responsible for update feed state
func (k Keeper) HandlerStateChanged(ctx sdk.Context, requestContextID tmbytes.HexBytes, _ string) {
	reqCtx, existed := k.sk.GetRequestContext(ctx, requestContextID)
	if !existed {
		ctx.Logger().Error(
			"Not existed requestContext", "requestContextID", requestContextID.String(),
		)
		return
	}

	feed, found := k.GetFeedByReqCtxID(ctx, requestContextID)
	if !found {
		ctx.Logger().Error(
			"Not existed feed", "requestContextID", requestContextID.String(),
		)
		return
	}

	var oldState serviceexported.RequestContextState
	switch reqCtx.State {
	case serviceexported.PAUSED:
		oldState = serviceexported.RUNNING
	case serviceexported.RUNNING:
		oldState = serviceexported.PAUSED
	case serviceexported.COMPLETED:
		ctx.Logger().Error(
			"Feed state invalid",
			"requestContextID", requestContextID.String(),
			"state", reqCtx.State.String(),
		)
		return
	}

	ctx.Logger().Info(
		"Feed state changed",
		"feed", feed.FeedName,
		"srcState", oldState,
		"dstState", reqCtx.State.String(),
	)
	k.dequeueAndEnqueue(ctx, feed.FeedName, oldState, reqCtx.State)
}

func (k Keeper) GetRequestContext(ctx sdk.Context, requestContextID tmbytes.HexBytes) (serviceexported.RequestContext, bool) {
	return k.sk.GetRequestContext(ctx, requestContextID)
}

func (k Keeper) ModuleServiceRequest(ctx sdk.Context, input string) (result string, output string) {
	feedName := gjson.Get(input, serviceexported.PATH_BODY).Get("pair").String()
	if _, found := k.GetFeed(ctx, feedName); !found {
		result = `{"code":"400","message":"feed not found"}`
		return
	}

	feedValues := k.GetFeedValues(ctx, feedName)
	if len(feedValues) == 0 {
		result = `{"code":"401","message":"no value"}`
		return
	}

	value := feedValues[0]
	valueData := value.Data
	valueTime := value.Timestamp

	if time.Since(valueTime) > time.Duration(time.Minute*5) {
		result = `{"code":"402","message":"all values expired"}`
		return
	}

	result = `{"code":"200","message":""}`
	output = fmt.Sprintf(`{"header":{},"body":{"rate":"%s"}}`, valueData)

	return
}
