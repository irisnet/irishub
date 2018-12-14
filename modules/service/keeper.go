package service

import (
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/tools/protoidl"
	"github.com/irisnet/irishub/modules/bank"
	"fmt"
	"github.com/irisnet/irishub/modules/service/params"
	"github.com/irisnet/irishub/modules/arbitration/params"
	"time"
	"github.com/irisnet/irishub/modules/guardian"
	"github.com/tendermint/tendermint/crypto"
)

var DepositedCoinsAccAddr = sdk.AccAddress(crypto.AddressHash([]byte("serviceDepositedCoins")))
var RequestCoinsAccAddr = sdk.AccAddress(crypto.AddressHash([]byte("serviceRequestCoins")))

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
	ck       bank.Keeper
	gk       guardian.Keeper

	// codespace
	codespace sdk.CodespaceType
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, ck bank.Keeper, gk guardian.Keeper, codespace sdk.CodespaceType) Keeper {
	keeper := Keeper{
		storeKey:  key,
		cdc:       cdc,
		ck:        ck,
		gk:        gk,
		codespace: codespace,
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

func (k Keeper) AddServiceBinding(ctx sdk.Context, svcBinding SvcBinding) (sdk.Error, bool) {
	kvStore := ctx.KVStore(k.storeKey)
	_, found := k.GetServiceDefinition(ctx, svcBinding.DefChainID, svcBinding.DefName)
	if !found {
		return ErrSvcDefNotExists(k.Codespace(), svcBinding.DefChainID, svcBinding.DefName), false
	}

	_, found = k.GetServiceBinding(ctx, svcBinding.DefChainID, svcBinding.DefName, svcBinding.BindChainID, svcBinding.Provider)
	if found {
		return ErrSvcBindingExists(k.Codespace()), false
	}

	minDeposit, err := getMinDeposit(ctx, svcBinding.Prices)
	if err != nil {
		return err, false
	}

	if !svcBinding.Deposit.IsAllGTE(minDeposit) {
		return ErrLtMinProviderDeposit(k.Codespace(), minDeposit), false
	}

	err = k.validateMethodPrices(ctx, svcBinding)
	if err != nil {
		return err, false
	}

	// Subtract coins from provider's account
	_, err = k.ck.SendCoins(ctx, svcBinding.Provider, DepositedCoinsAccAddr, svcBinding.Deposit)
	if err != nil {
		return err, false
	}

	svcBinding.DisableTime = time.Time{}
	svcBindingBytes := k.cdc.MustMarshalBinaryLengthPrefixed(svcBinding)
	kvStore.Set(GetServiceBindingKey(svcBinding.DefChainID, svcBinding.DefName, svcBinding.BindChainID, svcBinding.Provider), svcBindingBytes)
	return nil, true
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

func (k Keeper) UpdateServiceBinding(ctx sdk.Context, svcBinding SvcBinding) (sdk.Error, bool) {
	kvStore := ctx.KVStore(k.storeKey)
	oldBinding, found := k.GetServiceBinding(ctx, svcBinding.DefChainID, svcBinding.DefName, svcBinding.BindChainID, svcBinding.Provider)
	if !found {
		return ErrSvcBindingNotExists(k.Codespace()), false
	}

	if len(svcBinding.Prices) > 0 {
		err := k.validateMethodPrices(ctx, svcBinding)
		if err != nil {
			return err, false
		}
		oldBinding.Prices = svcBinding.Prices
	}

	if svcBinding.BindingType != 0x00 {
		oldBinding.BindingType = svcBinding.BindingType
	}

	// Add coins to svcBinding deposit
	if svcBinding.Deposit.IsNotNegative() {
		oldBinding.Deposit = oldBinding.Deposit.Plus(svcBinding.Deposit)
	}

	// Subtract coins from provider's account
	_, err := k.ck.SendCoins(ctx, svcBinding.Provider, DepositedCoinsAccAddr, svcBinding.Deposit)
	if err != nil {
		return err, false
	}

	if svcBinding.Level.UsableTime != 0 {
		oldBinding.Level.UsableTime = svcBinding.Level.UsableTime
	}
	if svcBinding.Level.AvgRspTime != 0 {
		oldBinding.Level.AvgRspTime = svcBinding.Level.AvgRspTime
	}

	// only check deposit if binding is available
	if oldBinding.Available {
		minDeposit, err := getMinDeposit(ctx, oldBinding.Prices)
		if err != nil {
			return err, false
		}

		if !oldBinding.Deposit.IsAllGTE(minDeposit) {
			return ErrLtMinProviderDeposit(k.Codespace(), minDeposit.Minus(oldBinding.Deposit).Plus(svcBinding.Deposit)), false
		}
	}

	svcBindingBytes := k.cdc.MustMarshalBinaryLengthPrefixed(oldBinding)
	kvStore.Set(GetServiceBindingKey(svcBinding.DefChainID, svcBinding.DefName, svcBinding.BindChainID, svcBinding.Provider), svcBindingBytes)
	return nil, true
}

func (k Keeper) Disable(ctx sdk.Context, defChainID, defName, bindChainID string, provider sdk.AccAddress) (sdk.Error, bool) {
	kvStore := ctx.KVStore(k.storeKey)
	binding, found := k.GetServiceBinding(ctx, defChainID, defName, bindChainID, provider)
	if !found {
		return ErrSvcBindingNotExists(k.Codespace()), false
	}

	if !binding.Available {
		return ErrDisable(k.Codespace(), "service binding is unavailable"), false
	}
	binding.Available = false
	binding.DisableTime = ctx.BlockHeader().Time
	svcBindingBytes := k.cdc.MustMarshalBinaryLengthPrefixed(binding)
	kvStore.Set(GetServiceBindingKey(binding.DefChainID, binding.DefName, binding.BindChainID, binding.Provider), svcBindingBytes)
	return nil, true
}

func (k Keeper) Enable(ctx sdk.Context, defChainID, defName, bindChainID string, provider sdk.AccAddress, deposit sdk.Coins) (sdk.Error, bool) {
	kvStore := ctx.KVStore(k.storeKey)
	binding, found := k.GetServiceBinding(ctx, defChainID, defName, bindChainID, provider)
	if !found {
		return ErrSvcBindingNotExists(k.Codespace()), false
	}

	if binding.Available {
		return ErrEnable(k.Codespace(), "service binding is available"), false
	}

	// Add coins to svcBinding deposit
	if deposit.IsNotNegative() {
		binding.Deposit = binding.Deposit.Plus(deposit)
	}

	minDeposit, err := getMinDeposit(ctx, binding.Prices)
	if err != nil {
		return err, false
	}

	if !binding.Deposit.IsAllGTE(minDeposit) {
		return ErrLtMinProviderDeposit(k.Codespace(), minDeposit.Minus(binding.Deposit).Plus(deposit)), false
	}

	// Subtract coins from provider's account
	_, err = k.ck.SendCoins(ctx, binding.Provider, DepositedCoinsAccAddr, deposit)
	if err != nil {
		return err, false
	}

	binding.Available = true
	binding.DisableTime = time.Time{}
	svcBindingBytes := k.cdc.MustMarshalBinaryLengthPrefixed(binding)
	kvStore.Set(GetServiceBindingKey(binding.DefChainID, binding.DefName, binding.BindChainID, binding.Provider), svcBindingBytes)
	return nil, true
}

func (k Keeper) RefundDeposit(ctx sdk.Context, defChainID, defName, bindChainID string, provider sdk.AccAddress) (sdk.Error, bool) {
	kvStore := ctx.KVStore(k.storeKey)
	binding, found := k.GetServiceBinding(ctx, defChainID, defName, bindChainID, provider)
	if !found {
		return ErrSvcBindingNotExists(k.Codespace()), false
	}

	if binding.Available {
		return ErrRefundDeposit(k.Codespace(), "can't refund from a available service binding"), false
	}

	if binding.Deposit.IsZero() {
		return ErrRefundDeposit(k.Codespace(), "service binding deposit is zero"), false
	}

	blockTime := ctx.BlockHeader().Time
	refundTime := binding.DisableTime.Add(arbitrationparams.GetArbitrationTimelimit(ctx)).Add(arbitrationparams.GetComplaintRetrospect(ctx))
	if blockTime.Before(refundTime) {
		return ErrRefundDeposit(k.Codespace(), fmt.Sprintf("can not refund deposit before %s", refundTime.Format("2006-01-02 15:04:05"))), false
	}

	// Add coins to provider's account
	_, err := k.ck.SendCoins(ctx, DepositedCoinsAccAddr, binding.Provider, binding.Deposit)
	if err != nil {
		return err, false
	}

	binding.Deposit = sdk.Coins{}

	svcBindingBytes := k.cdc.MustMarshalBinaryLengthPrefixed(binding)
	kvStore.Set(GetServiceBindingKey(binding.DefChainID, binding.DefName, binding.BindChainID, binding.Provider), svcBindingBytes)
	return nil, true
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

	maxTimeout := serviceparams.GetMaxRequestTimeout(ctx)
	req.ExpirationHeight = req.RequestHeight + maxTimeout

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(req)

	store.Set(GetRequestKey(req.DefChainID, req.DefName, req.BindChainID, req.Provider,
		req.RequestHeight, req.RequestIntraTxCounter), bz)

	_, err := k.ck.SendCoins(ctx, req.Consumer, RequestCoinsAccAddr, req.ServiceFee)
	if err != nil {
		return req, err
	}
	k.AddActiveRequest(ctx, req)
	k.AddRequestExpiration(ctx, req)
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

// Returns an iterator for all the request in the Active Queue that expire by block height
func (k Keeper) ActiveRequestQueueIterator(ctx sdk.Context, height int64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, GetRequestsByExpirationPrefix(height))
}

//__________________________________________________________________________

func (k Keeper) AddResponse(ctx sdk.Context, resp SvcResponse) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(resp)
	store.Set(GetResponseKey(resp.ReqChainID, resp.ExpirationHeight, resp.RequestHeight, resp.RequestIntraTxCounter), bz)
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
	k.SetReturnFee(ctx, address, fee.Coins.Plus(coins))
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
	feeTax := k.GetServiceFeeTax(ctx)
	taxFee := sdk.Coins{}
	for _, coin := range coins {
		taxFee = taxFee.Plus(sdk.Coins{sdk.Coin{Denom: coin.Denom, Amount: sdk.NewDecFromBigInt(coin.Amount.BigInt()).Mul(feeTax).TruncateInt()}})
	}

	taxPool := k.GetServiceFeeTaxPool(ctx)
	taxPool = taxPool.Plus(taxFee)
	k.SetServiceFeeTaxPool(ctx, taxPool)

	incomingFee, hasNeg := coins.SafeMinus(taxFee)
	if hasNeg {
		errMsg := fmt.Sprintf("%s is less than %s", coins, taxFee)
		return sdk.ErrInsufficientFunds(errMsg)
	}
	if !incomingFee.IsNotNegative() {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("%s is less than %s", incomingFee, taxFee))
	}
	fee, found := k.GetIncomingFee(ctx, address)
	if !found {
		k.SetIncomingFee(ctx, address, coins)
	}

	k.SetIncomingFee(ctx, address, fee.Coins.Plus(incomingFee))
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
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetIncomingFeeKey(address))
	return nil
}

//__________________________________________________________________________

func (k Keeper) GetServiceFeeTax(ctx sdk.Context) sdk.Dec {
	var percent sdk.Dec
	store := ctx.KVStore(k.storeKey)
	value := store.Get(serviceFeeTaxKey)
	if value == nil {
		return sdk.Dec{}
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &percent)
	return percent
}

func (k Keeper) SetServiceFeeTax(ctx sdk.Context, percent sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(percent)
	store.Set(serviceFeeTaxKey, bz)
}

func (k Keeper) GetServiceFeeTaxPool(ctx sdk.Context) sdk.Coins {
	var coins sdk.Coins
	store := ctx.KVStore(k.storeKey)
	value := store.Get(serviceFeeTaxPoolKey)
	if value == nil {
		return sdk.Coins{}
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &coins)
	return coins
}

func (k Keeper) SetServiceFeeTaxPool(ctx sdk.Context, coins sdk.Coins) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(coins)
	store.Set(serviceFeeTaxPoolKey, bz)
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

func getMinDeposit(ctx sdk.Context, prices []sdk.Coin) (sdk.Coins, sdk.Error) {
	// min deposit must >= sum(method price) * minDepositMultiple
	minDepositMultiple := sdk.NewInt(serviceparams.GetMinDepositMultiple(ctx))
	var minDeposit sdk.Coins
	for _, price := range prices {
		if price.Amount.BigInt().BitLen()+minDepositMultiple.BigInt().BitLen()-1 > 255 {
			return minDeposit, sdk.NewError(DefaultCodespace, CodeIntOverflow, fmt.Sprintf("Int Overflow"))
		}
		minInt := price.Amount.Mul(minDepositMultiple)
		minDeposit = minDeposit.Plus(sdk.Coins{sdk.Coin{Denom: price.Denom, Amount: minInt}})
	}
	return minDeposit, nil
}
