package keeper

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"github.com/tidwall/gjson"

	cmn "github.com/tendermint/tendermint/libs/common"

	"github.com/irisnet/irishub/app/v3/rand/internal/types"
	"github.com/irisnet/irishub/app/v3/service"
	"github.com/irisnet/irishub/app/v3/service/exported"
	sdk "github.com/irisnet/irishub/types"
)

// RequestService request the service for oracle seed
func (k Keeper) RequestService(ctx sdk.Context, reqID []byte, consumer sdk.AccAddress, serviceFeeCap sdk.Coins) ([]byte, sdk.Error) {
	iterator := k.sk.ServiceBindingsIterator(ctx, types.ServiceName)
	defer iterator.Close()

	var bindings []service.ServiceBinding
	for ; iterator.Valid(); iterator.Next() {
		var binding service.ServiceBinding
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &binding)

		bindings = append(bindings, binding)
	}

	if len(bindings) < 1 {
		return nil, types.ErrInvalidServiceBindings(types.DefaultCodespace, fmt.Sprintf("no service bindings available"))
	}

	coins := k.bk.GetCoins(ctx, consumer)
	if !coins.IsAllGTE(serviceFeeCap) {
		return nil, types.ErrInsufficientBalance(types.DefaultCodespace, fmt.Sprintf("insufficient balance"))
	}

	rand.Seed(time.Now().UnixNano())
	provider := []sdk.AccAddress{bindings[rand.Intn(len(bindings))].Provider}
	timeout := k.sk.GetParamSet(ctx).MaxRequestTimeout

	requestContextID, err := k.sk.CreateRequestContext(
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
		exported.RUNNING,
		1,
		types.ModuleName,
	)
	if err != nil {
		return nil, err
	}

	return requestContextID, nil
}

// HandlerResponse is responsible for processing the data returned from the service module
func (k Keeper) HandlerResponse(ctx sdk.Context, requestContextID cmn.HexBytes, responseOutput []string, err error) {
	if len(responseOutput) == 0 || err != nil {
		ctx = ctx.WithLogger(ctx.Logger().With("handler", "HandlerResponse"))
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

	_, existed := k.sk.GetRequestContext(ctx, requestContextID)
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

	result := gjson.Get(responseOutput[0], types.ValueJsonPath)

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
	rand := types.MakePRNG(lastBlockHash, currentTimestamp, request.Consumer, seed, true).GetRand()
	k.SetRand(ctx, reqID, types.NewRand(request.TxHash, lastBlockHeight, rand))

	k.DeleteOracleRandRequest(ctx, requestContextID)
}

// GetRequestContext retrieves the request context by the specified request context id
func (k Keeper) GetRequestContext(ctx sdk.Context, requestContextID []byte) (exported.RequestContext, bool) {
	return k.sk.GetRequestContext(ctx, requestContextID)
}

// GetMaxServiceRequestTimeout returns MaxServiceRequestTimeout
func (k Keeper) GetMaxServiceRequestTimeout(ctx sdk.Context) int64 {
	return k.sk.GetParamSet(ctx).MaxRequestTimeout
}
