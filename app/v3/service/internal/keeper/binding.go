package keeper

import (
	"fmt"
	"time"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v3/service/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

// AddServiceBinding creates a new service binding
func (k Keeper) AddServiceBinding(
	ctx sdk.Context,
	serviceName string,
	provider sdk.AccAddress,
	deposit sdk.Coins,
	pricing string,
) sdk.Error {
	if _, found := k.GetServiceDefinition(ctx, serviceName); !found {
		return types.ErrUnknownServiceDefinition(k.codespace, serviceName)
	}

	if _, found := k.GetServiceBinding(ctx, serviceName, provider); found {
		return types.ErrServiceBindingExists(k.codespace)
	}

	minDeposit := k.getMinDeposit(ctx, pricing)
	if !deposit.IsAllGTE(minDeposit) {
		return types.ErrInvalidDeposit(k.codespace, fmt.Sprintf("insufficient deposit: minimal deposit %s, %s got", minDeposit, deposit))
	}

	// Send coins from the provider's account to ServiceDepositCoinsAccAddr
	_, err := k.bk.SendCoins(ctx, provider, auth.ServiceDepositCoinsAccAddr, deposit)
	if err != nil {
		return err
	}

	available := true
	disabledTime := time.Time{}

	svcBinding := types.NewServiceBinding(serviceName, provider, deposit, pricing, available, disabledTime)
	k.SetServiceBinding(ctx, svcBinding)

	return nil
}

// UpdateServiceBinding updates the specified service binding
func (k Keeper) UpdateServiceBinding(
	ctx sdk.Context,
	serviceName string,
	provider sdk.AccAddress,
	deposit sdk.Coins,
	pricing string,
) sdk.Error {
	binding, found := k.GetServiceBinding(ctx, serviceName, provider)
	if !found {
		return types.ErrUnknownServiceBinding(k.codespace)
	}

	updated := false

	// add the deposit
	if !deposit.Empty() {
		binding.Deposit = binding.Deposit.Add(deposit)
		updated = true
	}

	// update the pricing
	if len(pricing) != 0 {
		binding.Pricing = pricing
		updated = true
	}

	// only check deposit when the binding is available and updated
	if binding.Available && updated {
		minDeposit := k.getMinDeposit(ctx, binding.Pricing)
		if !binding.Deposit.IsAllGTE(minDeposit) {
			return types.ErrInvalidDeposit(k.codespace, fmt.Sprintf("insufficient deposit: minimal deposit %s, %s got", minDeposit, binding.Deposit))
		}
	}

	if !deposit.Empty() {
		// Send coins from the provider's account to ServiceDepositCoinsAccAddr
		_, err := k.bk.SendCoins(ctx, provider, auth.ServiceDepositCoinsAccAddr, deposit)
		if err != nil {
			return err
		}
	}

	k.SetServiceBinding(ctx, binding)

	return nil
}

// DisableService disables the specified service binding
func (k Keeper) DisableService(ctx sdk.Context, serviceName string, provider sdk.AccAddress) sdk.Error {
	binding, found := k.GetServiceBinding(ctx, serviceName, provider)
	if !found {
		return types.ErrUnknownServiceBinding(k.codespace)
	}

	if !binding.Available {
		return types.ErrServiceBindingUnavailable(k.codespace)
	}

	binding.Available = false
	binding.DisabledTime = ctx.BlockHeader().Time

	k.SetServiceBinding(ctx, binding)

	return nil
}

// EnableService enables the specified service binding
func (k Keeper) EnableService(ctx sdk.Context, serviceName string, provider sdk.AccAddress, deposit sdk.Coins) sdk.Error {
	binding, found := k.GetServiceBinding(ctx, serviceName, provider)
	if !found {
		return types.ErrUnknownServiceBinding(k.codespace)
	}

	if binding.Available {
		return types.ErrServiceBindingAvailable(k.codespace)
	}

	// add the deposit
	if !deposit.Empty() {
		binding.Deposit = binding.Deposit.Add(deposit)
	}

	minDeposit := k.getMinDeposit(ctx, binding.Pricing)
	if !binding.Deposit.IsAllGTE(minDeposit) {
		return types.ErrInvalidDeposit(k.codespace, fmt.Sprintf("insufficient deposit: minimal deposit %s, %s got", minDeposit, binding.Deposit))
	}

	if !deposit.Empty() {
		// Send coins from the provider's account to ServiceDepositCoinsAccAddr
		_, err := k.bk.SendCoins(ctx, provider, auth.ServiceDepositCoinsAccAddr, deposit)
		if err != nil {
			return err
		}
	}

	binding.Available = true
	binding.DisabledTime = time.Time{}

	k.SetServiceBinding(ctx, binding)

	return nil
}

// RefundDeposit refunds the deposit from the specified service binding
func (k Keeper) RefundDeposit(ctx sdk.Context, serviceName string, provider sdk.AccAddress) sdk.Error {
	binding, found := k.GetServiceBinding(ctx, serviceName, provider)
	if !found {
		return types.ErrUnknownServiceBinding(k.codespace)
	}

	if binding.Available {
		return types.ErrServiceBindingAvailable(k.codespace)
	}

	if binding.Deposit.IsZero() {
		return types.ErrInvalidDeposit(k.codespace, "the deposit of the service binding is zero")
	}

	params := k.GetParamSet(ctx)
	refundableTime := binding.DisabledTime.Add(params.ArbitrationTimeLimit).Add(params.ComplaintRetrospect)

	currentTime := ctx.BlockHeader().Time
	if currentTime.Before(refundableTime) {
		return types.ErrIncorrectRefundTime(k.codespace, fmt.Sprintf("%v", refundableTime))
	}

	// Send coins from ServiceDepositCoinsAccAddr to the provider's account
	_, err := k.bk.SendCoins(ctx, auth.ServiceDepositCoinsAccAddr, binding.Provider, binding.Deposit)
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
		var binding types.ServiceBinding
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &binding)

		_, err := k.bk.SendCoins(ctx, auth.ServiceDepositCoinsAccAddr, binding.Provider, binding.Deposit)
		if err != nil {
			return err
		}
	}

	return nil
}

// SetServiceBinding sets the service binding
func (k Keeper) SetServiceBinding(ctx sdk.Context, svcBinding types.ServiceBinding) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(svcBinding)
	store.Set(GetServiceBindingKey(svcBinding.ServiceName, svcBinding.Provider), bz)
}

// GetServiceBinding retrieves the specified service binding
func (k Keeper) GetServiceBinding(ctx sdk.Context, serviceName string, provider sdk.AccAddress) (svcBinding types.ServiceBinding, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(GetServiceBindingKey(serviceName, provider))
	if bz == nil {
		return svcBinding, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &svcBinding)
	return svcBinding, true
}

// SetWithdrawAddress sets the withdrawal address for the specified provider
func (k Keeper) SetWithdrawAddress(ctx sdk.Context, provider, withdrawAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetWithdrawAddrKey(provider), withdrawAddr.Bytes())
}

// GetWithdrawAddress gets the withdrawal address of the specified provider
func (k Keeper) GetWithdrawAddress(ctx sdk.Context, provider sdk.AccAddress) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(GetWithdrawAddrKey(provider))
	if bz == nil {
		return provider
	}

	return sdk.AccAddress(bz)
}

// ServiceBindingsIterator returns an iterator for all bindings of the specified service
func (k Keeper) ServiceBindingsIterator(ctx sdk.Context, serviceName string) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, GetBindingsSubspace(serviceName))
}

// AllServiceBindingsIterator returns an iterator for all bindings
func (k Keeper) AllServiceBindingsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, serviceBindingKey)
}

// GetBasePrice gets the base price of the given service binding
// Note: ensure that the binding is valid
func (k Keeper) GetBasePrice(ctx sdk.Context, binding types.ServiceBinding) sdk.Coins {
	pricing, _ := types.ParsePricing(binding.Pricing)
	return pricing.Price
}

// getMinDeposit gets the minimal deposit required for the service binding
// Note: ensure that the pricing is valid
func (k Keeper) getMinDeposit(ctx sdk.Context, pricing string) sdk.Coins {
	params := k.GetParamSet(ctx)
	minDepositMultiple := sdk.NewInt(params.MinDepositMultiple)
	minDepositParam := params.MinDeposit

	p, _ := types.ParsePricing(pricing)
	price := p.Price.AmountOf(sdk.IrisAtto)

	// minimal deposit = max(price * minDepositMultiple, minDepositParam)
	minDeposit := sdk.NewCoins(sdk.NewCoin(sdk.IrisAtto, price.Mul(minDepositMultiple)))
	if minDeposit.IsAllLT(minDepositParam) {
		minDeposit = minDepositParam
	}

	return minDeposit
}
