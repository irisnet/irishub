package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/service/internal/types"
)

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
) (req types.SvcRequest, err sdk.Error) {
	binding, found := k.GetServiceBinding(ctx, defChainID, defName, bindChainID, provider)
	if !found {
		return req, types.ErrSvcBindingNotExists(k.codespace)
	}

	if !binding.Available {
		return req, types.ErrSvcBindingNotAvailable(k.codespace)
	}

	_, found = k.GetMethod(ctx, defChainID, defName, methodID)
	if !found {
		return req, types.ErrMethodNotExists(k.codespace, methodID)
	}

	if profiling {
		if _, found := k.gk.GetProfiler(ctx, consumer); !found {
			return req, types.ErrNotProfiler(k.codespace, consumer)
		}
	}

	//Method id start at 1
	if len(binding.Prices) >= int(methodID) && !serviceFee.IsAllGTE(sdk.Coins{binding.Prices[methodID-1]}) {
		return req, types.ErrLtServiceFee(k.codespace, sdk.Coins{binding.Prices[methodID-1]})
	}

	// request service fee is equal to service binding service fee if not profiling
	if len(binding.Prices) >= int(methodID) && !profiling {
		serviceFee = sdk.Coins{binding.Prices[methodID-1]}
	} else {
		serviceFee = nil
	}

	req = types.NewSvcRequest(defChainID, defName, bindChainID, reqChainID,
		consumer, provider, methodID, input, serviceFee, profiling)

	counter := k.GetIntraTxCounter(ctx)
	k.SetIntraTxCounter(ctx, counter+1)

	req.RequestIntraTxCounter = counter
	req.RequestHeight = ctx.BlockHeight()

	params := k.GetParams(ctx)
	req.ExpirationHeight = req.RequestHeight + params.MaxRequestTimeout

	err = k.sk.SendCoinsFromAccountToModule(ctx, req.Consumer, types.RequestAccName, req.ServiceFee)
	if err != nil {
		return req, err
	}

	k.SetRequest(ctx, req)
	k.AddActiveRequest(ctx, req)
	k.AddRequestExpiration(ctx, req)

	k.metrics.ActiveRequests.Add(1)

	return req, nil
}

func (k Keeper) SetRequest(ctx sdk.Context, req types.SvcRequest) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(req)
	store.Set(types.GetRequestKey(req.DefChainID, req.DefName, req.BindChainID,
		req.Provider, req.RequestHeight, req.RequestIntraTxCounter), bz)
}

func (k Keeper) AddActiveRequest(ctx sdk.Context, req types.SvcRequest) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(req)
	store.Set(types.GetActiveRequestKey(req.DefChainID, req.DefName, req.BindChainID,
		req.Provider, req.RequestHeight, req.RequestIntraTxCounter), bz)
}

func (k Keeper) DeleteActiveRequest(ctx sdk.Context, req types.SvcRequest) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetActiveRequestKey(req.DefChainID, req.DefName, req.BindChainID, req.Provider,
		req.RequestHeight, req.RequestIntraTxCounter))
}

func (k Keeper) AddRequestExpiration(ctx sdk.Context, req types.SvcRequest) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(req)
	store.Set(types.GetRequestsByExpirationIndexKeyByReq(req), bz)
}

func (k Keeper) DeleteRequestExpiration(ctx sdk.Context, req types.SvcRequest) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetRequestsByExpirationIndexKeyByReq(req))
}

func (k Keeper) GetActiveRequest(ctx sdk.Context, expHeight, reqHeight int64, counter int16) (req types.SvcRequest, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetRequestsByExpirationIndexKey(expHeight, reqHeight, counter))
	if bz == nil {
		return req, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &req)
	return req, true
}

// Returns an iterator for all the request in the Active Queue of specified service binding
func (k Keeper) ActiveBindRequestsIterator(ctx sdk.Context, defChainID, defName, bindChainID string, provider sdk.AccAddress) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.GetSubActiveRequestKey(defChainID, defName, bindChainID, provider))
}

// Returns an iterator for all the request in the Active Queue that expire by block height
func (k Keeper) ActiveRequestQueueIterator(ctx sdk.Context, height int64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.GetRequestsByExpirationPrefix(height))
}

// Returns an iterator for all the request in the Active Queue
func (k Keeper) ActiveAllRequestQueueIterator(store sdk.KVStore) sdk.Iterator {
	return sdk.KVStorePrefixIterator(store, types.ActiveRequestKey)
}

func (k Keeper) AddResponse(
	ctx sdk.Context,
	reqChainID string,
	requestID string,
	provider sdk.AccAddress,
	output,
	errorMsg []byte,
) (resp types.SvcResponse, err sdk.Error) {
	expHeight, reqHeight, counter, _ := types.ConvertRequestID(requestID)

	req, found := k.GetActiveRequest(ctx, expHeight, reqHeight, counter)
	if !found {
		req.ExpirationHeight = expHeight
		req.RequestHeight = reqHeight
		req.RequestIntraTxCounter = counter

		return resp, types.ErrRequestNotActive(k.codespace, req.RequestID())
	}

	if !(provider.Equals(req.Provider)) {
		return resp, types.ErrNotMatchingProvider(k.codespace, provider)
	}

	if reqChainID != req.ReqChainID {
		return resp, types.ErrNotMatchingReqChainID(k.codespace, reqChainID)
	}

	err = k.AddIncomingFee(ctx, resp.Provider, req.ServiceFee)
	if err != nil {
		return resp, err
	}

	resp = types.NewSvcResponse(reqChainID, expHeight, reqHeight, counter, provider, req.Consumer, output, errorMsg)
	k.SetResponse(ctx, resp)

	// delete request from active request list and expiration list
	k.DeleteActiveRequest(ctx, req)
	k.DeleteRequestExpiration(ctx, req)

	return resp, nil
}

func (k Keeper) SetResponse(ctx sdk.Context, resp types.SvcResponse) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(resp)
	store.Set(types.GetResponseKey(resp.ReqChainID, resp.ExpirationHeight, resp.RequestHeight, resp.RequestIntraTxCounter), bz)
}

func (k Keeper) GetResponse(ctx sdk.Context, reqChainID string, eHeight, rHeight int64, counter int16) (resp types.SvcResponse, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetResponseKey(reqChainID, eHeight, rHeight, counter))
	if bz == nil {
		return resp, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &resp)
	return resp, true
}

func (k Keeper) Slash(ctx sdk.Context, binding types.SvcBinding, slashCoins sdk.Coins) sdk.Error {
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

	ctx.Logger().Info("Slash service provider", "provider", binding.Provider.String(), "slash_amount", slashCoins.String())

	k.SetServiceBinding(ctx, binding)

	return nil
}

// Add return fee for a particular consumer, if it is not existed will create a new
func (k Keeper) AddReturnFee(ctx sdk.Context, address sdk.AccAddress, coins sdk.Coins) {
	fee, found := k.GetReturnFee(ctx, address)
	if !found {
		k.SetReturnFee(ctx, address, coins)
		return
	}

	k.SetReturnFee(ctx, address, fee.Coins.Add(coins))
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
func (k Keeper) RefundFee(ctx sdk.Context, address sdk.AccAddress) sdk.Error {
	fee, found := k.GetReturnFee(ctx, address)
	if !found {
		return types.ErrReturnFeeNotExists(k.codespace, address)
	}

	err := k.sk.SendCoinsFromModuleToAccount(ctx, types.RequestAccName, address, fee.Coins)
	if err != nil {
		return err
	}

	ctx.Logger().Info("Refund fees", "address", address.String(), "amount", fee.Coins.String())
	k.DeleteReturnFee(ctx, address)

	return nil
}

// Add incoming fee for a particular provider, if it is not existed will create a new
func (k Keeper) AddIncomingFee(ctx sdk.Context, address sdk.AccAddress, coins sdk.Coins) sdk.Error {
	params := k.GetParams(ctx)
	feeTax := params.ServiceFeeTax

	taxCoins := sdk.Coins{}
	for _, coin := range coins {
		taxAmount := sdk.NewDecFromInt(coin.Amount).Mul(feeTax).TruncateInt()
		taxCoins = append(taxCoins, sdk.NewCoin(coin.Denom, taxAmount))
	}

	taxCoins = taxCoins.Sort()

	err := k.sk.SendCoinsFromModuleToModule(ctx, types.RequestAccName, types.TaxAccName, taxCoins)
	if err != nil {
		return err
	}

	incomingFee, hasNeg := coins.SafeSub(taxCoins)
	if hasNeg {
		errMsg := fmt.Sprintf("%s is less than %s", coins, taxCoins)
		return sdk.ErrInsufficientFunds(errMsg)
	}

	fee, found := k.GetIncomingFee(ctx, address)
	if !found {
		k.SetIncomingFee(ctx, address, coins)
	}

	k.SetIncomingFee(ctx, address, fee.Coins.Add(incomingFee))
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
func (k Keeper) WithdrawFee(ctx sdk.Context, address sdk.AccAddress) sdk.Error {
	fee, found := k.GetIncomingFee(ctx, address)
	if !found {
		return types.ErrWithdrawFeeNotExists(k.codespace, address)
	}

	err := k.sk.SendCoinsFromModuleToAccount(ctx, types.RequestAccName, address, fee.Coins)
	if err != nil {
		return err
	}

	ctx.Logger().Info("Withdraw fees", "address", address.String(), "amount", fee.Coins.String())
	k.DeleteIncomingFee(ctx, address)

	return nil
}

func (k Keeper) WithdrawTax(ctx sdk.Context, trustee sdk.AccAddress, destAddress sdk.AccAddress, amt sdk.Coins) sdk.Error {
	_, found := k.gk.GetTrustee(ctx, trustee)
	if !found {
		return types.ErrNotTrustee(k.codespace, trustee)
	}

	err := k.sk.SendCoinsFromModuleToAccount(ctx, types.ServiceTaxAccName, destAddress, amt)
	return err
}

// get the current in-block request operation counter
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

// set the current in-block request counter
func (k Keeper) SetIntraTxCounter(ctx sdk.Context, counter int16) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(counter)
	store.Set(types.IntraTxCounterKey, bz)
}
