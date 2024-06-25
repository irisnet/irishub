package keeper

import (
	"encoding/hex"

	"github.com/tidwall/gjson"

	tmbytes "github.com/cometbft/cometbft/libs/bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"mods.irisnet.org/modules/random/types"
	"mods.irisnet.org/modules/service/exported"
	servicetypes "mods.irisnet.org/modules/service/types"
)

// RequestService requests the service for the oracle seed
func (k Keeper) RequestService(
	ctx sdk.Context,
	consumer sdk.AccAddress,
	serviceFeeCap sdk.Coins,
) (tmbytes.HexBytes, error) {
	iterator := k.serviceKeeper.ServiceBindingsIterator(ctx, types.ServiceName)
	defer iterator.Close()

	var bindings []servicetypes.ServiceBinding
	for ; iterator.Valid(); iterator.Next() {
		var binding servicetypes.ServiceBinding
		k.cdc.MustUnmarshal(iterator.Value(), &binding)

		bindings = append(bindings, binding)
	}

	if len(bindings) < 1 {
		return nil, types.ErrInvalidServiceBindings
	}

	if coins := k.bankKeeper.SpendableCoins(ctx, consumer); !coins.IsAllGTE(serviceFeeCap) {
		return nil, sdkerrors.ErrInsufficientFee
	}

	prng := types.MakePRNG(
		ctx.BlockHeader().LastBlockId.Hash,
		ctx.BlockHeader().Time.UnixNano(),
		consumer, nil, true)
	provider, err := sdk.AccAddressFromBech32(bindings[prng.Intn(len(bindings))].Provider)
	if err != nil {
		return nil, err
	}

	timeout := k.serviceKeeper.GetParams(ctx).MaxRequestTimeout

	return k.serviceKeeper.CreateRequestContext(
		ctx,
		types.ServiceName,
		[]sdk.AccAddress{provider},
		consumer,
		`{"header":{}}`,
		serviceFeeCap,
		timeout,
		false,
		0,
		0,
		exported.PAUSED,
		1,
		types.ModuleName,
	)
}

// StartRequestContext starts the service context
func (k Keeper) StartRequestContext(
	ctx sdk.Context,
	serviceContextID tmbytes.HexBytes,
	consumer sdk.AccAddress,
) error {
	return k.serviceKeeper.StartRequestContext(ctx, serviceContextID, consumer)
}

func (k Keeper) HandlerStateChanged(
	ctx sdk.Context,
	requestContextID tmbytes.HexBytes,
	err string,
) {
	reqCtx, existed := k.serviceKeeper.GetRequestContext(ctx, requestContextID)
	if !existed {
		ctx.Logger().Error(
			"Request context not found",
			"requestContextID", requestContextID.String(),
		)
		return
	}
	ctx.Logger().Error(
		"Oracle state invalid",
		"requestContextID", requestContextID.String(),
		"state", reqCtx.State.String(),
	)
	k.DeleteOracleRandRequest(ctx, requestContextID)
}

// HandlerResponse is responsible for processing the data returned from the service module
func (k Keeper) HandlerResponse(
	ctx sdk.Context,
	requestContextID tmbytes.HexBytes,
	responseOutput []string,
	err error,
) {
	if len(responseOutput) == 0 || err != nil {
		ctx.Logger().Error(
			"respond service failed",
			"requestContextID", requestContextID.String(),
			"err", err.Error(),
		)
		k.DeleteOracleRandRequest(ctx, requestContextID)
		return
	}

	if _, existed := k.serviceKeeper.GetRequestContext(ctx, requestContextID); !existed {
		k.DeleteOracleRandRequest(ctx, requestContextID)
		return
	}

	request, err := k.GetOracleRandRequest(ctx, requestContextID)
	if err != nil {
		ctx.Logger().Error(
			"request not found",
			"requestContextID", requestContextID.String(),
			"err", err.Error(),
		)
		k.DeleteOracleRandRequest(ctx, requestContextID)
		return
	}

	outputBody := gjson.Get(responseOutput[0], servicetypes.PATH_BODY).String()
	if err := servicetypes.ValidateResponseOutputBody(servicetypes.RandomServiceSchemas, outputBody); err != nil {
		ctx.Logger().Error(
			"invalid output body",
			"body", outputBody,
			"err", err.Error(),
		)
		return
	}

	seedStr := gjson.Get(outputBody, servicetypes.RandomServiceValueJSONPath).String()
	seed, err := hex.DecodeString(seedStr)
	if err != nil || len(seed) != types.SeedBytesLength {
		ctx.Logger().Error(
			"invalid seed",
			"seed", hex.EncodeToString(seed),
			"err", err.Error(),
		)
		k.DeleteOracleRandRequest(ctx, requestContextID)
		return
	}

	currentTimestamp := ctx.BlockHeader().Time.Unix()
	lastBlockHeight := ctx.BlockHeight() - 1
	lastBlockHash := ctx.BlockHeader().LastBlockId.Hash

	// get the request id
	reqID := types.GenerateRequestID(request)

	// generate a random number
	consumer, _ := sdk.AccAddressFromBech32(request.Consumer)
	random := types.MakePRNG(lastBlockHash, currentTimestamp, consumer, seed, true).GetRand()
	k.SetRandom(
		ctx,
		reqID,
		types.NewRandom(request.TxHash, lastBlockHeight, random.FloatString(types.RandPrec)),
	)

	k.DeleteOracleRandRequest(ctx, requestContextID)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeGenerateRandom,
			sdk.NewAttribute(types.AttributeKeyRequestID, hex.EncodeToString(reqID)),
			sdk.NewAttribute(types.AttributeKeyRandom, random.String()),
		),
	)
}

// GetRequestContext retrieves the request context by the specified request context id
func (k Keeper) GetRequestContext(
	ctx sdk.Context,
	requestContextID tmbytes.HexBytes,
) (exported.RequestContext, bool) {
	return k.serviceKeeper.GetRequestContext(ctx, requestContextID)
}

// GetMaxServiceRequestTimeout returns MaxServiceRequestTimeout
func (k Keeper) GetMaxServiceRequestTimeout(ctx sdk.Context) int64 {
	return k.serviceKeeper.GetParams(ctx).MaxRequestTimeout
}
