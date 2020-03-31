package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irishub/modules/service/internal/types"
)

// AddRequest
func (k Keeper) AddRequest(
	ctx sdk.Context,
	defChainID,
	defName,
	bindChainID,
	reqChainID string,
	consumer,
	provider sdk.AccAddress,
	methodID int16,
	input []byte,
	serviceFee sdk.Coins,
	profiling bool,
) (req types.SvcRequest, err error) {
	binding, found := k.GetServiceBinding(ctx, defChainID, defName, bindChainID, provider)
	if !found {
		return req, types.ErrUnknownSvcBinding
	}

	if !binding.Available {
		return req, types.ErrUnavailable
	}

	if profiling {
		if _, found := k.gk.GetProfiler(ctx, consumer); !found {
			return req, sdkerrors.Wrap(types.ErrUnknownMethod, consumer.String())
		}
	}

	//Method id start at 1
	if len(binding.Prices) >= int(methodID) && !serviceFee.IsAllGTE(sdk.Coins{binding.Prices[methodID-1]}) {
		return req, sdkerrors.Wrapf(types.ErrLtServiceFee, "service fee: %s, price: %s",
			serviceFee.String(), sdk.Coins{binding.Prices[methodID-1]}.String())
	}

	// request service fee is equal to service binding service fee if not profiling
	if len(binding.Prices) >= int(methodID) && !profiling {
		serviceFee = sdk.Coins{binding.Prices[methodID-1]}
	} else {
		serviceFee = nil
	}

	req = types.NewSvcRequest(
		defChainID, defName, bindChainID, reqChainID, consumer,
		provider, methodID, input, serviceFee, profiling,
	)

	counter := k.GetIntraTxCounter(ctx)
	k.SetIntraTxCounter(ctx, counter+1)

	req.RequestIntraTxCounter = counter
	req.RequestHeight = ctx.BlockHeight()

	params := k.GetParams(ctx)
	req.ExpirationHeight = req.RequestHeight + params.MaxRequestTimeout

	if err := k.sk.SendCoinsFromAccountToModule(
		ctx, req.Consumer, types.RequestAccName, req.ServiceFee,
	); err != nil {
		return req, err
	}

	k.SetRequest(ctx, req)
	k.AddActiveRequest(ctx, req)
	k.AddRequestExpiration(ctx, req)

	return req, nil
}

// SetRequest
func (k Keeper) SetRequest(ctx sdk.Context, req types.SvcRequest) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(req)
	store.Set(types.GetRequestKey(req.DefChainID, req.DefName, req.BindChainID,
		req.Provider, req.RequestHeight, req.RequestIntraTxCounter), bz)
}

// AddActiveRequest
func (k Keeper) AddActiveRequest(ctx sdk.Context, req types.SvcRequest) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(req)
	store.Set(types.GetActiveRequestKey(req.DefChainID, req.DefName, req.BindChainID,
		req.Provider, req.RequestHeight, req.RequestIntraTxCounter), bz)
}

// DeleteActiveRequest
func (k Keeper) DeleteActiveRequest(ctx sdk.Context, req types.SvcRequest) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetActiveRequestKey(req.DefChainID, req.DefName, req.BindChainID, req.Provider,
		req.RequestHeight, req.RequestIntraTxCounter))
}

// AddRequestExpiration
func (k Keeper) AddRequestExpiration(ctx sdk.Context, req types.SvcRequest) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(req)
	store.Set(types.GetRequestsByExpirationIndexKeyByReq(req), bz)
}

// DeleteRequestExpiration
func (k Keeper) DeleteRequestExpiration(ctx sdk.Context, req types.SvcRequest) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetRequestsByExpirationIndexKeyByReq(req))
}

// GetActiveRequest
func (k Keeper) GetActiveRequest(ctx sdk.Context, expHeight, reqHeight int64, counter int16) (req types.SvcRequest, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetRequestsByExpirationIndexKey(expHeight, reqHeight, counter))
	if bz == nil {
		return req, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &req)
	return req, true
}

// ActiveBindRequestsIterator returns an iterator for all the request in the Active Queue of specified service binding
func (k Keeper) ActiveBindRequestsIterator(ctx sdk.Context, defChainID, defName, bindChainID string, provider sdk.AccAddress) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.GetSubActiveRequestKey(defChainID, defName, bindChainID, provider))
}

// ActiveRequestQueueIterator returns an iterator for all the request in the Active Queue that expire by block height
func (k Keeper) ActiveRequestQueueIterator(ctx sdk.Context, height int64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.GetRequestsByExpirationPrefix(height))
}

// ActiveAllRequestQueueIterator returns an iterator for all the requests in the active queue
func (k Keeper) ActiveAllRequestQueueIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.ActiveRequestKey)
}

// AddResponse
func (k Keeper) AddResponse(
	ctx sdk.Context,
	reqChainID string,
	requestID string,
	provider sdk.AccAddress,
	output,
	errorMsg []byte,
) (resp types.SvcResponse, err error) {
	expHeight, reqHeight, counter, _ := types.ConvertRequestID(requestID)

	req, found := k.GetActiveRequest(ctx, expHeight, reqHeight, counter)
	if !found {
		req.ExpirationHeight = expHeight
		req.RequestHeight = reqHeight
		req.RequestIntraTxCounter = counter
		return resp, sdkerrors.Wrap(types.ErrUnknownActiveRequest, req.RequestID())
	}

	if !(provider.Equals(req.Provider)) {
		return resp, sdkerrors.Wrapf(types.ErrNotMatchingProvider,
			"expected: %s, got: %s", provider.String(), req.Provider.String())
	}

	if reqChainID != req.ReqChainID {
		return resp, sdkerrors.Wrapf(types.ErrNotMatchingReqChainID,
			"expected: %s, got: %s", reqChainID, req.ReqChainID)
	}

	if err := k.AddIncomingFee(ctx, provider, req.ServiceFee); err != nil {
		return resp, err
	}

	resp = types.NewSvcResponse(reqChainID, expHeight, reqHeight, counter, provider, req.Consumer, output, errorMsg)
	k.SetResponse(ctx, resp)

	// delete request from active request list and expiration list
	k.DeleteActiveRequest(ctx, req)
	k.DeleteRequestExpiration(ctx, req)

	return resp, nil
}

// SetResponse
func (k Keeper) SetResponse(ctx sdk.Context, resp types.SvcResponse) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(resp)
	store.Set(types.GetResponseKey(resp.ReqChainID, resp.ExpirationHeight, resp.RequestHeight, resp.RequestIntraTxCounter), bz)
}

// GetResponse
func (k Keeper) GetResponse(ctx sdk.Context, reqChainID string, eHeight, rHeight int64, counter int16) (resp types.SvcResponse, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetResponseKey(reqChainID, eHeight, rHeight, counter))
	if bz == nil {
		return resp, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &resp)
	return resp, true
}

// Slash
func (k Keeper) Slash(ctx sdk.Context, binding types.SvcBinding, slashCoins sdk.Coins) error {
	deposit, hasNeg := binding.Deposit.SafeSub(slashCoins)
	if hasNeg {
		errMsg := fmt.Sprintf("%s is less than %s", binding.Deposit, slashCoins)
		panic(errMsg)
	}

	binding.Deposit = deposit
	minDeposit, err := k.getMinDeposit(ctx, binding.Prices)
	if err != nil {
		return err
	}

	if !binding.Deposit.IsAllGTE(minDeposit) {
		binding.Available = false
		binding.DisableTime = ctx.BlockHeader().Time
	}

	if err := k.sk.BurnCoins(ctx, types.DepositAccName, slashCoins); err != nil {
		return err
	}

	k.Logger(ctx).Info("Slash service provider", "provider", binding.Provider.String(), "slash_amount", slashCoins.String())
	k.SetServiceBinding(ctx, binding)

	return nil
}

// AddReturnFee add return fee for a particular consumer, if it is not existed will create a new
func (k Keeper) AddReturnFee(ctx sdk.Context, address sdk.AccAddress, coins sdk.Coins) {
	fee, found := k.GetReturnFee(ctx, address)
	if !found {
		k.SetReturnFee(ctx, address, coins)
		return
	}

	k.SetReturnFee(ctx, address, fee.Coins.Add(coins...))
}

func (k Keeper) SetReturnFee(ctx sdk.Context, address sdk.AccAddress, coins sdk.Coins) {
	store := ctx.KVStore(k.storeKey)

	fee := types.NewReturnedFee(address, coins)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(fee)

	store.Set(types.GetReturnedFeeKey(address), bz)
}

func (k Keeper) DeleteReturnFee(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetReturnedFeeKey(address))
}

func (k Keeper) GetReturnFee(ctx sdk.Context, address sdk.AccAddress) (fee types.ReturnedFee, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetReturnedFeeKey(address))
	if bz == nil {
		return fee, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &fee)
	return fee, true
}

// refund fees from a particular consumer, and delete it
func (k Keeper) RefundFee(ctx sdk.Context, address sdk.AccAddress) error {
	fee, found := k.GetReturnFee(ctx, address)
	if !found {
		return types.ErrUnknownReturnFee
	}

	if err := k.sk.SendCoinsFromModuleToAccount(ctx, types.RequestAccName, address, fee.Coins); err != nil {
		return err
	}

	k.Logger(ctx).Info("Refund fees", "address", address.String(), "amount", fee.Coins.String())
	k.DeleteReturnFee(ctx, address)

	return nil
}

// Add incoming fee for a particular provider, if it is not existed will create a new
func (k Keeper) AddIncomingFee(ctx sdk.Context, address sdk.AccAddress, coins sdk.Coins) error {
	params := k.GetParams(ctx)
	feeTax := params.ServiceFeeTax

	taxCoins := sdk.Coins{}
	for _, coin := range coins {
		taxAmount := sdk.NewDecFromInt(coin.Amount).Mul(feeTax).TruncateInt()
		taxCoins = taxCoins.Add(sdk.NewCoin(coin.Denom, taxAmount))
	}

	if err := k.sk.SendCoinsFromModuleToModule(ctx, types.RequestAccName, types.TaxAccName, taxCoins); err != nil {
		return err
	}

	incomingFee, hasNeg := coins.SafeSub(taxCoins)
	if hasNeg {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInsufficientFunds, "insufficient account funds; %s < %s", coins, taxCoins,
		)
	}

	fee, _ := k.GetIncomingFee(ctx, address)
	k.SetIncomingFee(ctx, address, fee.Coins.Add(incomingFee...))
	return nil
}

func (k Keeper) SetIncomingFee(ctx sdk.Context, address sdk.AccAddress, coins sdk.Coins) {
	store := ctx.KVStore(k.storeKey)

	fee := types.NewIncomingFee(address, coins)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(fee)

	store.Set(types.GetIncomingFeeKey(address), bz)
}

func (k Keeper) DeleteIncomingFee(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetIncomingFeeKey(address))
}

func (k Keeper) GetIncomingFee(ctx sdk.Context, address sdk.AccAddress) (fee types.IncomingFee, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetIncomingFeeKey(address))
	if bz == nil {
		return fee, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &fee)
	return fee, true
}

// withdraw fees from a particular provider, and delete it
func (k Keeper) WithdrawFee(ctx sdk.Context, address sdk.AccAddress) error {
	fee, found := k.GetIncomingFee(ctx, address)
	if !found {
		return types.ErrUnknownWithdrawFee
	}

	if err := k.sk.SendCoinsFromModuleToAccount(ctx, types.RequestAccName, address, fee.Coins); err != nil {
		return err
	}

	k.Logger(ctx).Info("Withdraw fees", "address", address.String(), "amount", fee.Coins.String())
	k.DeleteIncomingFee(ctx, address)

	return nil
}

func (k Keeper) WithdrawTax(ctx sdk.Context, trustee sdk.AccAddress, destAddress sdk.AccAddress, amt sdk.Coins) error {
	if _, found := k.gk.GetTrustee(ctx, trustee); !found {
		return types.ErrUnknownTrustee
	}
	return k.sk.SendCoinsFromModuleToAccount(ctx, types.TaxAccName, destAddress, amt)
}

// AllReturnedFeesIterator returns an iterator for all the returned fees
func (k Keeper) AllReturnedFeesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.ReturnedFeeKey)
}

// AllIncomingFeesIterator returns an iterator for all the incoming fees
func (k Keeper) AllIncomingFeesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.IncomingFeeKey)
}

// RefundReturnedFees refunds all the returned fees
func (k Keeper) RefundReturnedFees(ctx sdk.Context) error {
	iterator := k.AllReturnedFeesIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var returnedFee types.ReturnedFee
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &returnedFee)

		if err := k.sk.SendCoinsFromModuleToAccount(
			ctx, types.RequestAccName, returnedFee.Address, returnedFee.Coins,
		); err != nil {
			return err
		}
	}

	return nil
}

// RefundIncomingFees refunds all the incoming fees
func (k Keeper) RefundIncomingFees(ctx sdk.Context) error {
	iterator := k.AllIncomingFeesIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var incomingFee types.IncomingFee
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &incomingFee)

		if err := k.sk.SendCoinsFromModuleToAccount(
			ctx, types.RequestAccName, incomingFee.Address, incomingFee.Coins,
		); err != nil {
			return err
		}
	}

	return nil
}

// RefundServiceFees refunds the service fees of all the active requests
func (k Keeper) RefundServiceFees(ctx sdk.Context) error {
	iterator := k.ActiveAllRequestQueueIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var request types.SvcRequest
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &request)

		if err := k.sk.SendCoinsFromModuleToAccount(
			ctx, types.RequestAccName, request.Consumer, request.ServiceFee,
		); err != nil {
			return err
		}
	}

	return nil
}

// GetIntraTxCounter get the current in-block request operation counter
func (k Keeper) GetIntraTxCounter(ctx sdk.Context) int16 {
	store := ctx.KVStore(k.storeKey)

	b := store.Get(types.IntraTxCounterKey)
	if b == nil {
		return 0
	}

	var counter int16
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &counter)

	return counter
}

// SetIntraTxCounter set the current in-block request counter
func (k Keeper) SetIntraTxCounter(ctx sdk.Context, counter int16) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(counter)
	store.Set(types.IntraTxCounterKey, bz)
}
