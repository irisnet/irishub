package service

import (
	"fmt"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/bank"
	"github.com/irisnet/irishub/modules/guardian"
	"github.com/irisnet/irishub/modules/params"
	"github.com/irisnet/irishub/tools/protoidl"
	sdk "github.com/irisnet/irishub/types"
	"github.com/tendermint/tendermint/crypto"
	"time"
)

var DepositedCoinsAccAddr = sdk.AccAddress(crypto.AddressHash([]byte("serviceDepositedCoins")))
var RequestCoinsAccAddr = sdk.AccAddress(crypto.AddressHash([]byte("serviceRequestCoins")))
var TaxCoinsAccAddr = sdk.AccAddress(crypto.AddressHash([]byte("serviceTaxCoins")))

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
	ck       bank.Keeper
	gk       guardian.Keeper

	// codespace
	codespace sdk.CodespaceType
	// params subspace
	paramSpace params.Subspace
	// metrics
	metrics *Metrics
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, ck bank.Keeper, gk guardian.Keeper, codespace sdk.CodespaceType, paramSpace params.Subspace, metrics *Metrics) Keeper {
	keeper := Keeper{
		storeKey:   key,
		cdc:        cdc,
		ck:         ck,
		gk:         gk,
		codespace:  codespace,
		paramSpace: paramSpace.WithTypeTable(ParamTypeTable()),
		metrics:    metrics,
	}
	return keeper
}

// return the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

func (k Keeper) AddServiceDefinition(ctx sdk.Context, svcDef SvcDef) {
	kvStore := ctx.KVStore(k.storeKey)

	svcDefBytes, err := k.cdc.MarshalBinaryLengthPrefixed(svcDef)
	if err != nil {
		panic(err)
	}

	kvStore.Set(GetServiceDefinitionKey(svcDef.ChainId, svcDef.Name), svcDefBytes)
}

func (k Keeper) AddMethods(ctx sdk.Context, svcDef SvcDef) sdk.Error {
	methods, err := protoidl.GetMethods(svcDef.IDLContent)
	if err != nil {
		panic(err)
	}
	kvStore := ctx.KVStore(k.storeKey)
	for index, method := range methods {
		methodProperty, err := methodToMethodProperty(index+1, method)
		if err != nil {
			return err
		}
		methodBytes := k.cdc.MustMarshalBinaryLengthPrefixed(methodProperty)
		kvStore.Set(GetMethodPropertyKey(svcDef.ChainId, svcDef.Name, methodProperty.ID), methodBytes)
	}
	return nil
}

func (k Keeper) GetServiceDefinition(ctx sdk.Context, chainId, name string) (svcDef SvcDef, found bool) {
	kvStore := ctx.KVStore(k.storeKey)

	serviceDefBytes := kvStore.Get(GetServiceDefinitionKey(chainId, name))
	if serviceDefBytes != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(serviceDefBytes, &svcDef)
		return svcDef, true
	}
	return svcDef, false
}

// Gets the method in a specific service and methodID
func (k Keeper) GetMethod(ctx sdk.Context, chainId, name string, id int16) (method MethodProperty, found bool) {
	store := ctx.KVStore(k.storeKey)
	methodBytes := store.Get(GetMethodPropertyKey(chainId, name, id))
	if methodBytes != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(methodBytes, &method)
		return method, true
	}
	return method, false
}

// Gets all the methods in a specific service
func (k Keeper) GetMethods(ctx sdk.Context, chainId, name string) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, GetMethodsSubspaceKey(chainId, name))
}

func (k Keeper) AddServiceBinding(ctx sdk.Context, svcBinding SvcBinding) sdk.Error {
	kvStore := ctx.KVStore(k.storeKey)
	_, found := k.GetServiceDefinition(ctx, svcBinding.DefChainID, svcBinding.DefName)
	if !found {
		return ErrSvcDefNotExists(k.Codespace(), svcBinding.DefChainID, svcBinding.DefName)
	}

	_, found = k.GetServiceBinding(ctx, svcBinding.DefChainID, svcBinding.DefName, svcBinding.BindChainID, svcBinding.Provider)
	if found {
		return ErrSvcBindingExists(k.Codespace())
	}

	minDeposit, err := k.getMinDeposit(ctx, svcBinding.Prices)
	if err != nil {
		return err
	}

	if !svcBinding.Deposit.IsAllGTE(minDeposit) {
		return ErrLtMinProviderDeposit(k.Codespace(), minDeposit)
	}

	err = k.validateMethodPrices(ctx, svcBinding)
	if err != nil {
		return err
	}

	// Subtract coins from provider's account
	_, err = k.ck.SendCoins(ctx, svcBinding.Provider, DepositedCoinsAccAddr, svcBinding.Deposit)
	if err != nil {
		return err
	}

	svcBinding.DisableTime = time.Time{}
	svcBindingBytes := k.cdc.MustMarshalBinaryLengthPrefixed(svcBinding)
	kvStore.Set(GetServiceBindingKey(svcBinding.DefChainID, svcBinding.DefName, svcBinding.BindChainID, svcBinding.Provider), svcBindingBytes)
	return nil
}

func (k Keeper) GetServiceBinding(ctx sdk.Context, defChainID, defName, bindChainID string, provider sdk.AccAddress) (svcBinding SvcBinding, found bool) {
	kvStore := ctx.KVStore(k.storeKey)

	svcBindingBytes := kvStore.Get(GetServiceBindingKey(defChainID, defName, bindChainID, provider))
	if svcBindingBytes != nil {
		var svcBinding SvcBinding
		k.cdc.MustUnmarshalBinaryLengthPrefixed(svcBindingBytes, &svcBinding)
		return svcBinding, true
	}
	return svcBinding, false
}

func (k Keeper) ServiceBindingsIterator(ctx sdk.Context, defChainID, defName string) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, GetBindingsSubspaceKey(defChainID, defName))
}

func (k Keeper) UpdateServiceBinding(ctx sdk.Context, svcBinding SvcBinding) sdk.Error {
	kvStore := ctx.KVStore(k.storeKey)
	oldBinding, found := k.GetServiceBinding(ctx, svcBinding.DefChainID, svcBinding.DefName, svcBinding.BindChainID, svcBinding.Provider)
	if !found {
		return ErrSvcBindingNotExists(k.Codespace())
	}

	if len(svcBinding.Prices) > 0 {
		err := k.validateMethodPrices(ctx, svcBinding)
		if err != nil {
			return err
		}
		oldBinding.Prices = svcBinding.Prices
	}

	if svcBinding.BindingType != 0x00 {
		oldBinding.BindingType = svcBinding.BindingType
	}

	// Add coins to svcBinding deposit
	if !svcBinding.Deposit.IsAnyNegative() {
		oldBinding.Deposit = oldBinding.Deposit.Add(svcBinding.Deposit)
	}

	// Subtract coins from provider's account
	_, err := k.ck.SendCoins(ctx, svcBinding.Provider, DepositedCoinsAccAddr, svcBinding.Deposit)
	if err != nil {
		return err
	}

	if svcBinding.Level.UsableTime != 0 {
		oldBinding.Level.UsableTime = svcBinding.Level.UsableTime
	}
	if svcBinding.Level.AvgRspTime != 0 {
		oldBinding.Level.AvgRspTime = svcBinding.Level.AvgRspTime
	}

	// only check deposit if binding is available
	if oldBinding.Available {
		minDeposit, err := k.getMinDeposit(ctx, oldBinding.Prices)
		if err != nil {
			return err
		}

		if !oldBinding.Deposit.IsAllGTE(minDeposit) {
			return ErrLtMinProviderDeposit(k.Codespace(), minDeposit.Sub(oldBinding.Deposit).Add(svcBinding.Deposit))
		}
	}

	svcBindingBytes := k.cdc.MustMarshalBinaryLengthPrefixed(oldBinding)
	kvStore.Set(GetServiceBindingKey(svcBinding.DefChainID, svcBinding.DefName, svcBinding.BindChainID, svcBinding.Provider), svcBindingBytes)
	return nil
}

func (k Keeper) Disable(ctx sdk.Context, defChainID, defName, bindChainID string, provider sdk.AccAddress) sdk.Error {
	kvStore := ctx.KVStore(k.storeKey)
	binding, found := k.GetServiceBinding(ctx, defChainID, defName, bindChainID, provider)
	if !found {
		return ErrSvcBindingNotExists(k.Codespace())
	}

	if !binding.Available {
		return ErrDisable(k.Codespace(), "service binding is unavailable")
	}
	binding.Available = false
	binding.DisableTime = ctx.BlockHeader().Time
	svcBindingBytes := k.cdc.MustMarshalBinaryLengthPrefixed(binding)
	kvStore.Set(GetServiceBindingKey(binding.DefChainID, binding.DefName, binding.BindChainID, binding.Provider), svcBindingBytes)
	return nil
}

func (k Keeper) Enable(ctx sdk.Context, defChainID, defName, bindChainID string, provider sdk.AccAddress, deposit sdk.Coins) sdk.Error {
	kvStore := ctx.KVStore(k.storeKey)
	binding, found := k.GetServiceBinding(ctx, defChainID, defName, bindChainID, provider)
	if !found {
		return ErrSvcBindingNotExists(k.Codespace())
	}

	if binding.Available {
		return ErrEnable(k.Codespace(), "service binding is available")
	}

	// Add coins to svcBinding deposit
	if !deposit.IsAnyNegative() {
		binding.Deposit = binding.Deposit.Add(deposit)
	}

	minDeposit, err := k.getMinDeposit(ctx, binding.Prices)
	if err != nil {
		return err
	}

	if !binding.Deposit.IsAllGTE(minDeposit) {
		return ErrLtMinProviderDeposit(k.Codespace(), minDeposit.Sub(binding.Deposit).Add(deposit))
	}

	// Subtract coins from provider's account
	_, err = k.ck.SendCoins(ctx, binding.Provider, DepositedCoinsAccAddr, deposit)
	if err != nil {
		return err
	}

	binding.Available = true
	binding.DisableTime = time.Time{}
	svcBindingBytes := k.cdc.MustMarshalBinaryLengthPrefixed(binding)
	kvStore.Set(GetServiceBindingKey(binding.DefChainID, binding.DefName, binding.BindChainID, binding.Provider), svcBindingBytes)
	return nil
}

func (k Keeper) RefundDeposit(ctx sdk.Context, defChainID, defName, bindChainID string, provider sdk.AccAddress) sdk.Error {
	kvStore := ctx.KVStore(k.storeKey)
	binding, found := k.GetServiceBinding(ctx, defChainID, defName, bindChainID, provider)
	if !found {
		return ErrSvcBindingNotExists(k.Codespace())
	}

	if binding.Available {
		return ErrRefundDeposit(k.Codespace(), "can't refund from a available service binding")
	}

	if binding.Deposit.IsZero() {
		return ErrRefundDeposit(k.Codespace(), "service binding deposit is zero")
	}

	blockTime := ctx.BlockHeader().Time
	params := k.GetParamSet(ctx)
	refundTime := binding.DisableTime.Add(params.ArbitrationTimeLimit).Add(params.ComplaintRetrospect)
	if blockTime.Before(refundTime) {
		return ErrRefundDeposit(k.Codespace(), fmt.Sprintf("can not refund deposit before %s", refundTime.Format("2006-01-02 15:04:05")))
	}

	// Add coins to provider's account
	_, err := k.ck.SendCoins(ctx, DepositedCoinsAccAddr, binding.Provider, binding.Deposit)
	if err != nil {
		return err
	}

	binding.Deposit = sdk.Coins{}

	svcBindingBytes := k.cdc.MustMarshalBinaryLengthPrefixed(binding)
	kvStore.Set(GetServiceBindingKey(binding.DefChainID, binding.DefName, binding.BindChainID, binding.Provider), svcBindingBytes)
	return nil
}

func (k Keeper) validateMethodPrices(ctx sdk.Context, svcBinding SvcBinding) sdk.Error {
	iterator := k.GetMethods(ctx, svcBinding.DefChainID, svcBinding.DefName)
	defer iterator.Close()
	var methods []MethodProperty
	for ; iterator.Valid(); iterator.Next() {
		var method MethodProperty
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &method)
		methods = append(methods, method)
	}

	if len(methods) != len(svcBinding.Prices) {
		return ErrInvalidPriceCount(k.Codespace(), len(svcBinding.Prices), len(methods))
	}
	return nil
}

//__________________________________________________________________________

func (k Keeper) AddRequest(ctx sdk.Context, req SvcRequest) (SvcRequest, sdk.Error) {
	store := ctx.KVStore(k.storeKey)

	counter := k.GetIntraTxCounter(ctx)
	req.RequestHeight = ctx.BlockHeight()
	req.RequestIntraTxCounter = counter
	k.SetIntraTxCounter(ctx, counter+1)

	params := k.GetParamSet(ctx)
	req.ExpirationHeight = req.RequestHeight + params.MaxRequestTimeout

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(req)

	store.Set(GetRequestKey(req.DefChainID, req.DefName, req.BindChainID, req.Provider,
		req.RequestHeight, req.RequestIntraTxCounter), bz)

	_, err := k.ck.SendCoins(ctx, req.Consumer, RequestCoinsAccAddr, req.ServiceFee)
	if err != nil {
		return req, err
	}
	k.AddActiveRequest(ctx, req)
	k.AddRequestExpiration(ctx, req)
	k.metrics.ActiveRequests.Add(1)
	return req, nil
}

func (k Keeper) AddActiveRequest(ctx sdk.Context, req SvcRequest) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(req)
	store.Set(GetActiveRequestKey(req.DefChainID, req.DefName, req.BindChainID, req.Provider,
		req.RequestHeight, req.RequestIntraTxCounter), bz)
}

func (k Keeper) DeleteActiveRequest(ctx sdk.Context, req SvcRequest) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetActiveRequestKey(req.DefChainID, req.DefName, req.BindChainID, req.Provider,
		req.RequestHeight, req.RequestIntraTxCounter))
}

func (k Keeper) AddRequestExpiration(ctx sdk.Context, req SvcRequest) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(req)
	store.Set(GetRequestsByExpirationIndexKeyByReq(req), bz)
}

func (k Keeper) DeleteRequestExpiration(ctx sdk.Context, req SvcRequest) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetRequestsByExpirationIndexKeyByReq(req))
}

func (k Keeper) GetActiveRequest(ctx sdk.Context, eHeight, rHeight int64, counter int16) (req SvcRequest, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(GetRequestsByExpirationIndexKey(eHeight, rHeight, counter))
	if value == nil {
		return req, false
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &req)
	return req, true
}

// Returns an iterator for all the request in the Active Queue of specified service binding
func (k Keeper) ActiveBindRequestsIterator(ctx sdk.Context, defChainID, defName, bindChainID string, provider sdk.AccAddress) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, GetSubActiveRequestKey(defChainID, defName, bindChainID, provider))
}

// Returns an iterator for all the request in the Active Queue that expire by block height
func (k Keeper) ActiveRequestQueueIterator(ctx sdk.Context, height int64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, GetRequestsByExpirationPrefix(height))
}

// Returns an iterator for all the request in the Active Queue
func (k Keeper) ActiveAllRequestQueueIterator(store sdk.KVStore) sdk.Iterator {
	return sdk.KVStorePrefixIterator(store, activeRequestKey)
}

//__________________________________________________________________________

func (k Keeper) AddResponse(ctx sdk.Context, resp SvcResponse) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(resp)
	store.Set(GetResponseKey(resp.ReqChainID, resp.ExpirationHeight, resp.RequestHeight, resp.RequestIntraTxCounter), bz)
}

func (k Keeper) GetResponse(ctx sdk.Context, reqChainID string, eHeight, rHeight int64, counter int16) (resp SvcResponse, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(GetResponseKey(reqChainID, eHeight, rHeight, counter))
	if value == nil {
		return resp, false
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &resp)
	return resp, true
}

//__________________________________________________________________________

func (k Keeper) SetReturnFee(ctx sdk.Context, address sdk.AccAddress, coins sdk.Coins) {
	store := ctx.KVStore(k.storeKey)
	fee := NewReturnedFee(address, coins)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(fee)
	store.Set(GetReturnedFeeKey(address), bz)
}

func (k Keeper) GetReturnFee(ctx sdk.Context, address sdk.AccAddress) (fee ReturnedFee, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(GetReturnedFeeKey(address))
	if value == nil {
		return fee, false
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &fee)
	return fee, true
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

// refund fees from a particular consumer, and delete it
func (k Keeper) RefundFee(ctx sdk.Context, address sdk.AccAddress) sdk.Error {
	fee, found := k.GetReturnFee(ctx, address)
	if !found {
		return ErrReturnFeeNotExists(k.Codespace(), address)
	}

	_, err := k.ck.SendCoins(ctx, RequestCoinsAccAddr, address, fee.Coins)
	if err != nil {
		return err
	}
	ctx.Logger().Info("Refund fees", "address", address.String(), "amount", fee.Coins.String())
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetReturnedFeeKey(address))
	return nil
}

func (k Keeper) SetIncomingFee(ctx sdk.Context, address sdk.AccAddress, coins sdk.Coins) {
	store := ctx.KVStore(k.storeKey)
	fee := NewIncomingFee(address, coins)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(fee)
	store.Set(GetIncomingFeeKey(address), bz)
}

func (k Keeper) GetIncomingFee(ctx sdk.Context, address sdk.AccAddress) (fee IncomingFee, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(GetIncomingFeeKey(address))
	if value == nil {
		return fee, false
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &fee)
	return fee, true
}

// Add incoming fee for a particular provider, if it is not existed will create a new
func (k Keeper) AddIncomingFee(ctx sdk.Context, address sdk.AccAddress, coins sdk.Coins) sdk.Error {
	params := k.GetParamSet(ctx)
	feeTax := params.ServiceFeeTax
	taxCoins := sdk.Coins{}
	for _, coin := range coins {
		taxAmount := sdk.NewDecFromInt(coin.Amount).Mul(feeTax).TruncateInt()
		taxCoins = append(taxCoins, sdk.NewCoin(coin.Denom, taxAmount))
	}
	taxCoins = taxCoins.Sort()

	_, err := k.ck.SendCoins(ctx, RequestCoinsAccAddr, TaxCoinsAccAddr, taxCoins)
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

// withdraw fees from a particular provider, and delete it
func (k Keeper) WithdrawFee(ctx sdk.Context, address sdk.AccAddress) sdk.Error {
	fee, found := k.GetIncomingFee(ctx, address)
	if !found {
		return ErrWithdrawFeeNotExists(k.Codespace(), address)
	}
	_, err := k.ck.SendCoins(ctx, RequestCoinsAccAddr, address, fee.Coins)
	if err != nil {
		return err
	}
	ctx.Logger().Info("Withdraw fees", "address", address.String(), "amount", fee.Coins.String())
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetIncomingFeeKey(address))
	return nil
}

func (k Keeper) Slash(ctx sdk.Context, binding SvcBinding, slashCoins sdk.Coins) sdk.Error {
	store := ctx.KVStore(k.storeKey)
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
	svcBindingBytes := k.cdc.MustMarshalBinaryLengthPrefixed(binding)
	store.Set(GetServiceBindingKey(binding.DefChainID, binding.DefName, binding.BindChainID, binding.Provider), svcBindingBytes)
	return nil
}

//__________________________________________________________________________

// get the current in-block request operation counter
func (k Keeper) GetIntraTxCounter(ctx sdk.Context) int16 {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(intraTxCounterKey)
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
	store.Set(intraTxCounterKey, bz)
}

func (k Keeper) getMinDeposit(ctx sdk.Context, prices []sdk.Coin) (sdk.Coins, sdk.Error) {
	params := k.GetParamSet(ctx)
	// min deposit must >= sum(method price) * minDepositMultiple
	minDepositMultiple := sdk.NewInt(params.MinDepositMultiple)
	var minDeposit sdk.Coins
	for _, price := range prices {
		if price.Amount.BigInt().BitLen()+minDepositMultiple.BigInt().BitLen()-1 > 255 {
			return minDeposit, sdk.NewError(DefaultCodespace, CodeIntOverflow, fmt.Sprintf("Int Overflow"))
		}
		minInt := price.Amount.Mul(minDepositMultiple)
		minDeposit = minDeposit.Add(sdk.Coins{sdk.NewCoin(price.Denom, minInt)})
	}
	return minDeposit, nil
}

func (k Keeper) InitMetrics(store sdk.KVStore) {
	activeIterator := k.ActiveAllRequestQueueIterator(store)
	defer activeIterator.Close()
	for ; activeIterator.Valid(); activeIterator.Next() {
		k.metrics.ActiveRequests.Add(1)
	}
}
