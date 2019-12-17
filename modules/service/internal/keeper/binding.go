package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/service/internal/types"
)

func (k Keeper) AddServiceBinding(
	ctx sdk.Context,
	defChainID,
	defName,
	bindChainID string,
	provider sdk.AccAddress,
	bindingType types.BindingType,
	deposit sdk.Coins,
	prices []sdk.Coin,
	level types.Level,
) sdk.Error {
	_, found := k.GetServiceDefinition(ctx, defChainID, defName)
	if !found {
		return types.ErrSvcDefNotExists(k.codespace, defChainID, defName)
	}

	_, found = k.GetServiceBinding(ctx, defChainID, defName, bindChainID, provider)
	if found {
		return types.ErrSvcBindingExists(k.codespace)
	}

	minDeposit, err := k.getMinDeposit(ctx, prices)
	if err != nil {
		return err
	}

	if !deposit.IsAllGTE(minDeposit) {
		return types.ErrLtMinProviderDeposit(k.codespace, minDeposit)
	}

	svcBinding := types.NewSvcBinding(ctx, defChainID, defName, bindChainID, provider, bindingType, deposit, prices, level, true)

	err = k.validateMethodPrices(ctx, svcBinding)
	if err != nil {
		return err
	}

	// Send coins from provider's account to the deposit module account
	err = k.sk.SendCoinsFromAccountToModule(ctx, svcBinding.Provider, types.DepositAccName, svcBinding.Deposit)
	if err != nil {
		return err
	}

	svcBinding.DisableTime = time.Time{}
	k.SetServiceBinding(ctx, svcBinding)

	return nil
}

func (k Keeper) SetServiceBinding(ctx sdk.Context, svcBinding types.SvcBinding) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(svcBinding)
	store.Set(types.GetServiceBindingKey(svcBinding.DefChainID, svcBinding.DefName, svcBinding.BindChainID, svcBinding.Provider), bz)
}

func (k Keeper) GetServiceBinding(ctx sdk.Context, defChainID, defName, bindChainID string, provider sdk.AccAddress) (svcBinding types.SvcBinding, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetServiceBindingKey(defChainID, defName, bindChainID, provider))
	if bz == nil {
		return svcBinding, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &svcBinding)
	return svcBinding, true
}

func (k Keeper) ServiceBindingsIterator(ctx sdk.Context, defChainID, defName string) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.GetBindingsSubspaceKey(defChainID, defName))
}

// AllServiceBindingsIterator returns an iterator for all the binding services
func (k Keeper) AllServiceBindingsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.BindingPropertyKey)
}

func (k Keeper) UpdateServiceBinding(
	ctx sdk.Context,
	defChainID,
	defName,
	bindChainID string,
	provider sdk.AccAddress,
	bindingType types.BindingType,
	deposit sdk.Coins,
	prices []sdk.Coin,
	level types.Level,
) (svcBinding types.SvcBinding, err sdk.Error) {
	oldBinding, found := k.GetServiceBinding(ctx, defChainID, defName, bindChainID, provider)
	if !found {
		return svcBinding, types.ErrSvcBindingNotExists(k.codespace)
	}

	newBinding := types.NewSvcBinding(ctx, defChainID, defName, bindChainID, provider, bindingType,
		deposit, prices, level, false)

	if len(prices) > 0 {
		err = k.validateMethodPrices(ctx, newBinding)
		if err != nil {
			return svcBinding, err
		}

		oldBinding.Prices = newBinding.Prices
	}

	if newBinding.BindingType != 0x00 {
		oldBinding.BindingType = newBinding.BindingType
	}

	// Add coins to svcBinding deposit
	if !newBinding.Deposit.IsAnyNegative() {
		oldBinding.Deposit = oldBinding.Deposit.Add(newBinding.Deposit)
	}

	// Send coins from provider's account to the deposit module account
	err = k.sk.SendCoinsFromAccountToModule(ctx, svcBinding.Provider, types.DepositAccName, newBinding.Deposit)
	if err != nil {
		return svcBinding, err
	}

	if newBinding.Level.UsableTime != 0 {
		oldBinding.Level.UsableTime = newBinding.Level.UsableTime
	}
	if newBinding.Level.AvgRspTime != 0 {
		oldBinding.Level.AvgRspTime = newBinding.Level.AvgRspTime
	}

	// only check deposit if binding is available
	if oldBinding.Available {
		minDeposit, err := k.getMinDeposit(ctx, oldBinding.Prices)
		if err != nil {
			return svcBinding, err
		}

		if !oldBinding.Deposit.IsAllGTE(minDeposit) {
			return svcBinding, types.ErrLtMinProviderDeposit(k.codespace, minDeposit.Sub(oldBinding.Deposit).Add(newBinding.Deposit))
		}
	}

	k.SetServiceBinding(ctx, oldBinding)

	return oldBinding, nil
}

func (k Keeper) Disable(ctx sdk.Context, defChainID, defName, bindChainID string, provider sdk.AccAddress) sdk.Error {
	binding, found := k.GetServiceBinding(ctx, defChainID, defName, bindChainID, provider)
	if !found {
		return types.ErrSvcBindingNotExists(k.codespace)
	}

	if !binding.Available {
		return types.ErrDisable(k.Codespace(), "service binding is unavailable")
	}

	binding.Available = false
	binding.DisableTime = ctx.BlockHeader().Time

	k.SetServiceBinding(ctx, binding)

	return nil
}

func (k Keeper) Enable(ctx sdk.Context, defChainID, defName, bindChainID string, provider sdk.AccAddress, deposit sdk.Coins) sdk.Error {
	binding, found := k.GetServiceBinding(ctx, defChainID, defName, bindChainID, provider)
	if !found {
		return types.ErrSvcBindingNotExists(k.codespace)
	}

	if binding.Available {
		return types.ErrEnable(k.Codespace(), "service binding is available")
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
		return types.ErrLtMinProviderDeposit(k.codespace, minDeposit.Sub(binding.Deposit).Add(deposit))
	}

	// Send coins from provider's account to the deposit module account
	err = k.sk.SendCoinsFromAccountToModule(ctx, binding.Provider, types.DepositAccName, deposit)
	if err != nil {
		return err
	}

	binding.Available = true
	binding.DisableTime = time.Time{}

	k.SetServiceBinding(ctx, binding)

	return nil
}

func (k Keeper) RefundDeposit(ctx sdk.Context, defChainID, defName, bindChainID string, provider sdk.AccAddress) sdk.Error {
	binding, found := k.GetServiceBinding(ctx, defChainID, defName, bindChainID, provider)
	if !found {
		return types.ErrSvcBindingNotExists(k.Codespace())
	}

	if binding.Available {
		return types.ErrRefundDeposit(k.Codespace(), "can't refund from a available service binding")
	}

	if binding.Deposit.IsZero() {
		return types.ErrRefundDeposit(k.Codespace(), "service binding deposit is zero")
	}

	blockTime := ctx.BlockHeader().Time
	params := k.GetParams(ctx)

	refundTime := binding.DisableTime.Add(params.ArbitrationTimeLimit).Add(params.ComplaintRetrospect)
	if blockTime.Before(refundTime) {
		return types.ErrRefundDeposit(k.Codespace(), fmt.Sprintf("can not refund deposit before %s", refundTime.Format("2006-01-02 15:04:05")))
	}

	// Send coins from the deposit module account to the provider's account
	err := k.sk.SendCoinsFromModuleToAccount(ctx, types.DepositAccName, binding.Provider, binding.Deposit)
	if err != nil {
		return err
	}

	binding.Deposit = sdk.Coins{}
	k.SetServiceBinding(ctx, binding)

	return nil
}

// RefundDeposits refunds the deposits of all the binding services
func (k Keeper) RefundDeposits(ctx sdk.Context) sdk.Error {
	iterator := k.AllServiceBindingsIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var binding types.SvcBinding
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &binding)

		err := k.sk.SendCoinsFromModuleToAccount(ctx, types.DepositAccName, binding.Provider, binding.Deposit)
		if err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) getMinDeposit(ctx sdk.Context, prices []sdk.Coin) (sdk.Coins, sdk.Error) {
	params := k.GetParams(ctx)
	// min deposit must >= sum(method price) * minDepositMultiple
	minDepositMultiple := sdk.NewInt(params.MinDepositMultiple)

	var minDeposit sdk.Coins
	for _, price := range prices {
		if price.Amount.BigInt().BitLen()+minDepositMultiple.BigInt().BitLen()-1 > 255 {
			return minDeposit, sdk.NewError(k.codespace, types.CodeIntOverflow, fmt.Sprintf("Int Overflow"))
		}

		minInt := price.Amount.Mul(minDepositMultiple)
		minDeposit = minDeposit.Add(sdk.NewCoins(sdk.NewCoin(price.Denom, minInt)))
	}

	return minDeposit, nil
}

func (k Keeper) validateMethodPrices(ctx sdk.Context, svcBinding types.SvcBinding) sdk.Error {
	iterator := k.GetMethods(ctx, svcBinding.DefChainID, svcBinding.DefName)
	defer iterator.Close()

	var methods []types.MethodProperty
	for ; iterator.Valid(); iterator.Next() {
		var method types.MethodProperty
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &method)
		methods = append(methods, method)
	}

	if len(methods) != len(svcBinding.Prices) {
		return types.ErrInvalidPriceCount(k.Codespace(), len(svcBinding.Prices), len(methods))
	}

	return nil
}
