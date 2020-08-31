package keeper

import (
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/tidwall/gjson"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	servicetypes "github.com/irismod/service/types"

	"github.com/irismod/service/exported"

	"github.com/irisnet/irishub/modules/random/types"
)

// RequestService request the service for oracle seed
func (k Keeper) RequestService(ctx sdk.Context, consumer sdk.AccAddress, serviceFeeCap sdk.Coins) (tmbytes.HexBytes, error) {
	iterator := k.serviceKeeper.ServiceBindingsIterator(ctx, types.ServiceName)
	defer iterator.Close()

	var bindings []servicetypes.ServiceBinding
	for ; iterator.Valid(); iterator.Next() {
		var binding servicetypes.ServiceBinding
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &binding)

		bindings = append(bindings, binding)
	}

	if len(bindings) < 1 {
		return nil, types.ErrInvalidServiceBindings
	}

	coins := k.bankKeeper.SpendableCoins(ctx, consumer)
	if !coins.IsAllGTE(serviceFeeCap) {
		return nil, sdkerrors.ErrInsufficientFee
	}

	rand.Seed(time.Now().UnixNano())
	provider := []sdk.AccAddress{bindings[rand.Intn(len(bindings))].Provider}
	timeout := k.serviceKeeper.GetParams(ctx).MaxRequestTimeout

	requestContextID, err := k.serviceKeeper.CreateRequestContext(
		ctx,
		types.ServiceName,
		provider,
		consumer,
		"{}",
		serviceFeeCap,
		timeout,
		false,
		false,
		0,
		0,
		exported.PAUSED,
		1,
		types.ModuleName,
	)
	if err != nil {
		return nil, err
	}

	return requestContextID, nil
}

// StartRequestContext starts the service context
func (k Keeper) StartRequestContext(
	ctx sdk.Context,
	serviceContextID tmbytes.HexBytes,
	consumer sdk.AccAddress,
) error {
	return k.serviceKeeper.StartRequestContext(ctx, serviceContextID, consumer)
}

func (k Keeper) HandlerStateChanged(ctx sdk.Context, requestContextID tmbytes.HexBytes, err string) {
	reqCtx, existed := k.serviceKeeper.GetRequestContext(ctx, requestContextID)
	if !existed {
		ctx.Logger().Error(
			"Not existed requestContext",
			"requestContextID", requestContextID.String(),
		)
		return
	}
	ctx.Logger().Error(
		"Oracle state invalid", "requestContextID",
		requestContextID.String(), "state", reqCtx.State.String(),
	)
	k.DeleteOracleRandRequest(ctx, requestContextID)
	return
}

// HandlerResponse is responsible for processing the data returned from the service module
func (k Keeper) HandlerResponse(ctx sdk.Context, requestContextID tmbytes.HexBytes, responseOutput []string, err error) {
	if len(responseOutput) == 0 || err != nil {
		ctx.Logger().Error(
			"respond service failed",
			"requestContextID",
			requestContextID.String(),
			"err",
			err.Error(),
		)
		k.DeleteOracleRandRequest(ctx, requestContextID)
		return
	}

	_, existed := k.serviceKeeper.GetRequestContext(ctx, requestContextID)
	if !existed {
		k.DeleteOracleRandRequest(ctx, requestContextID)
		return
	}

	request, err := k.GetOracleRandRequest(ctx, requestContextID)
	if err != nil {
		ctx.Logger().Error(
			"can not find request",
			"requestContextID",
			requestContextID.String(),
			"err",
			err.Error(),
		)
		k.DeleteOracleRandRequest(ctx, requestContextID)
		return
	}

	result := gjson.Get(responseOutput[0], types.ServiceValueJsonPath)

	seed, err := hex.DecodeString(result.String())
	if err != nil || len(seed) != types.SeedBytesLength {
		ctx.Logger().Error(
			"invalid seed",
			"seed",
			hex.EncodeToString(seed),
			"err",
			err.Error(),
		)
		k.DeleteOracleRandRequest(ctx, requestContextID)
		return
	}

	currentTimestamp := ctx.BlockHeader().Time.Unix()
	lastBlockHeight := ctx.BlockHeight() - 1
	lastBlockHash := []byte(ctx.BlockHeader().LastBlockId.Hash)

	// get the request id
	reqID := types.GenerateRequestID(request)

	// generate a random number
	random := types.MakePRNG(lastBlockHash, currentTimestamp, request.Consumer, seed, true).GetRand()
	k.SetRandom(ctx, reqID, types.NewRandom(request.TxHash, lastBlockHeight, random.FloatString(types.RandPrec)))

	k.DeleteOracleRandRequest(ctx, requestContextID)
}

// GetRequestContext retrieves the request context by the specified request context id
func (k Keeper) GetRequestContext(ctx sdk.Context, requestContextID tmbytes.HexBytes) (exported.RequestContext, bool) {
	return k.serviceKeeper.GetRequestContext(ctx, requestContextID)
}

// GetMaxServiceRequestTimeout returns MaxServiceRequestTimeout
func (k Keeper) GetMaxServiceRequestTimeout(ctx sdk.Context) int64 {
	return k.serviceKeeper.GetParams(ctx).MaxRequestTimeout
}
