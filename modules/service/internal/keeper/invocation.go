package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/irisnet/irishub/modules/service/internal/types"
)

func (k Keeper) AddRequest(ctx sdk.Context, req types.SvcRequest) (types.SvcRequest, sdk.Error) {
	counter := k.GetIntraTxCounter(ctx)
	k.SetIntraTxCounter(ctx, counter+1)

	req.RequestIntraTxCounter = counter
	req.RequestHeight = ctx.BlockHeight()

	params := k.GetParams(ctx)
	req.ExpirationHeight = req.RequestHeight + params.MaxRequestTimeout

	err := k.bk.SendCoins(ctx, req.Consumer, auth.ServiceRequestCoinsAccAddr, req.ServiceFee)
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

func (k Keeper) GetActiveRequest(ctx sdk.Context, expirationHeight, requestHeight int64, counter int16) (req types.SvcRequest, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetRequestsByExpirationIndexKey(expirationHeight, requestHeight, counter))
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

func (k Keeper) AddResponse(ctx sdk.Context, resp types.SvcResponse) {
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

	err := k.bk.SendCoins(ctx, auth.ServiceRequestCoinsAccAddr, address, fee.Coins)
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

	err := k.bk.SendCoins(ctx, auth.ServiceRequestCoinsAccAddr, auth.ServiceTaxCoinsAccAddr, taxCoins)
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

	err := k.bk.SendCoins(ctx, auth.ServiceRequestCoinsAccAddr, address, fee.Coins)
	if err != nil {
		return err
	}

	ctx.Logger().Info("Withdraw fees", "address", address.String(), "amount", fee.Coins.String())
	k.DeleteIncomingFee(ctx, address)

	return nil
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
