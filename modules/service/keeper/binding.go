package keeper

import (
	"bytes"
	"time"

	gogotypes "github.com/gogo/protobuf/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/service/types"
)

// AddServiceBinding creates a new service binding
func (k Keeper) AddServiceBinding(
	ctx sdk.Context,
	serviceName string,
	provider sdk.AccAddress,
	deposit sdk.Coins,
	pricing string,
	qos uint64,
	options string,
	owner sdk.AccAddress,
) error {
	if _, found := k.GetServiceDefinition(ctx, serviceName); !found {
		return sdkerrors.Wrap(types.ErrUnknownServiceDefinition, serviceName)
	}

	if _, found := k.GetServiceBinding(ctx, serviceName, provider); found {
		return sdkerrors.Wrap(types.ErrServiceBindingExists, "")
	}

	currentOwner, found := k.GetOwner(ctx, provider)
	if found && !owner.Equals(currentOwner) {
		return sdkerrors.Wrap(types.ErrNotAuthorized, "owner not matching")
	}

	if err := k.validateDeposit(ctx, deposit); err != nil {
		return err
	}

	maxReqTimeout := k.MaxRequestTimeout(ctx)
	if qos > uint64(maxReqTimeout) {
		return sdkerrors.Wrapf(
			types.ErrInvalidQoS,
			"QoS [%d] must not be greater than maximum request timeout [%d]",
			qos, maxReqTimeout,
		)
	}

	if err := types.ValidateOptions(options); err != nil {
		return err
	}

	parsedPricing, err := k.ParsePricing(ctx, pricing)
	if err != nil {
		return err
	}

	minDeposit, err := k.GetMinDeposit(ctx, parsedPricing)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrInvalidMinDeposit, "%s", err)
	}

	if !deposit.IsAllGTE(minDeposit) {
		return sdkerrors.Wrapf(
			types.ErrInvalidDeposit,
			"insufficient deposit: minimum deposit %s, %s got",
			minDeposit, deposit,
		)
	}

	// Send coins from owner's account to the deposit module account
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, owner, types.DepositAccName, deposit); err != nil {
		return err
	}

	available := true
	disabledTime := time.Time{}

	svcBinding := types.NewServiceBinding(serviceName, provider, deposit, pricing, qos, options, available, disabledTime, owner)

	k.SetServiceBinding(ctx, svcBinding)
	k.SetOwnerServiceBinding(ctx, svcBinding)
	k.SetPricing(ctx, serviceName, provider, parsedPricing)

	if currentOwner.Empty() {
		k.SetOwner(ctx, provider, owner)
		k.SetOwnerProvider(ctx, owner, provider)
	}

	return nil
}

// SetServiceBindingForGenesis sets the service binding for genesis
func (k Keeper) SetServiceBindingForGenesis(
	ctx sdk.Context,
	svcBinding types.ServiceBinding,
) error {
	provider, err := sdk.AccAddressFromBech32(svcBinding.Provider)
	if err != nil {
		return err
	}
	owner, err := sdk.AccAddressFromBech32(svcBinding.Owner)
	if err != nil {
		return err
	}

	pricing, err := types.ParsePricing(svcBinding.Pricing)
	if err != nil {
		return err
	}

	k.SetServiceBinding(ctx, svcBinding)
	k.SetOwnerServiceBinding(ctx, svcBinding)
	k.SetOwner(ctx, provider, owner)
	k.SetOwnerProvider(ctx, owner, provider)

	k.SetPricing(ctx, svcBinding.ServiceName, provider, pricing)

	return nil
}

// UpdateServiceBinding updates the specified service binding
func (k Keeper) UpdateServiceBinding(
	ctx sdk.Context,
	serviceName string,
	provider sdk.AccAddress,
	deposit sdk.Coins,
	pricing string,
	qos uint64,
	options string,
	owner sdk.AccAddress,
) error {
	binding, found := k.GetServiceBinding(ctx, serviceName, provider)
	if !found {
		return sdkerrors.Wrap(types.ErrUnknownServiceBinding, "")
	}

	bindingOwner, err := sdk.AccAddressFromBech32(binding.Owner)
	if err != nil {
		return err
	}

	if !owner.Equals(bindingOwner) {
		return sdkerrors.Wrap(types.ErrNotAuthorized, "owner not matching")
	}

	updated := false

	if qos != 0 {
		maxReqTimeout := k.MaxRequestTimeout(ctx)
		if qos > uint64(maxReqTimeout) {
			return sdkerrors.Wrapf(
				types.ErrInvalidQoS,
				"QoS [%d] must not be greater than maximum request timeout [%d]",
				qos, maxReqTimeout,
			)
		}

		binding.QoS = qos
		updated = true
	}

	// add the deposit
	if !deposit.Empty() {
		if err := k.validateDeposit(ctx, deposit); err != nil {
			return err
		}

		binding.Deposit = binding.Deposit.Add(deposit...)
		updated = true
	}

	parsedPricing := k.GetPricing(ctx, serviceName, provider)

	// update the pricing
	if len(pricing) != 0 {
		parsedPricing, err := k.ParsePricing(ctx, pricing)
		if err != nil {
			return err
		}

		binding.Pricing = pricing
		k.SetPricing(ctx, serviceName, provider, parsedPricing)

		updated = true
	}

	// update options
	if len(options) != 0 {
		if err := types.ValidateOptions(options); err != nil {
			return err
		}
		binding.Options = options
		updated = true
	}

	// only check deposit when the binding is available and updated
	if binding.Available && updated {
		minDeposit, err := k.GetMinDeposit(ctx, parsedPricing)
		if err != nil {
			return sdkerrors.Wrapf(types.ErrInvalidMinDeposit, "%s", err)
		}

		if !binding.Deposit.IsAllGTE(minDeposit) {
			return sdkerrors.Wrapf(
				types.ErrInvalidDeposit,
				"insufficient deposit: minimum deposit %s, %s got",
				minDeposit, binding.Deposit,
			)
		}
	}

	if !deposit.Empty() {
		// Send coins from owner's account to the deposit module account
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, owner, types.DepositAccName, deposit); err != nil {
			return err
		}
	}

	if updated {
		k.SetServiceBinding(ctx, binding)
	}

	return nil
}

// DisableServiceBinding disables the specified service binding
func (k Keeper) DisableServiceBinding(
	ctx sdk.Context,
	serviceName string,
	provider,
	owner sdk.AccAddress,
) error {
	binding, found := k.GetServiceBinding(ctx, serviceName, provider)
	if !found {
		return sdkerrors.Wrap(types.ErrUnknownServiceBinding, "")
	}

	bindingOwner, err := sdk.AccAddressFromBech32(binding.Owner)
	if err != nil {
		return err
	}

	if !owner.Equals(bindingOwner) {
		return sdkerrors.Wrap(types.ErrNotAuthorized, "owner not matching")
	}

	if !binding.Available {
		return sdkerrors.Wrap(types.ErrServiceBindingUnavailable, "")
	}

	binding.Available = false
	binding.DisabledTime = ctx.BlockHeader().Time

	k.SetServiceBinding(ctx, binding)

	return nil
}

// EnableServiceBinding enables the specified service binding
func (k Keeper) EnableServiceBinding(
	ctx sdk.Context,
	serviceName string,
	provider sdk.AccAddress,
	deposit sdk.Coins,
	owner sdk.AccAddress,
) error {
	binding, found := k.GetServiceBinding(ctx, serviceName, provider)
	if !found {
		return sdkerrors.Wrap(types.ErrUnknownServiceBinding, "")
	}

	bindingOwner, err := sdk.AccAddressFromBech32(binding.Owner)
	if err != nil {
		return err
	}

	if !owner.Equals(bindingOwner) {
		return sdkerrors.Wrap(types.ErrNotAuthorized, "owner not matching")
	}

	if binding.Available {
		return sdkerrors.Wrap(types.ErrServiceBindingAvailable, "")
	}

	// add the deposit
	if !deposit.Empty() {
		if err := k.validateDeposit(ctx, deposit); err != nil {
			return err
		}

		binding.Deposit = binding.Deposit.Add(deposit...)
	}

	minDeposit, err := k.GetMinDeposit(ctx, k.GetPricing(ctx, serviceName, provider))
	if err != nil {
		return sdkerrors.Wrapf(types.ErrInvalidMinDeposit, "%s", err)
	}

	if !binding.Deposit.IsAllGTE(minDeposit) {
		return sdkerrors.Wrapf(
			types.ErrInvalidDeposit,
			"insufficient deposit: minimum deposit %s, %s got",
			minDeposit, binding.Deposit,
		)
	}

	if !deposit.Empty() {
		// Send coins from owner's account to the deposit module account
		if err := k.bankKeeper.SendCoinsFromAccountToModule(
			ctx, owner, types.DepositAccName, deposit,
		); err != nil {
			return err
		}
	}

	binding.Available = true
	binding.DisabledTime = time.Time{}

	k.SetServiceBinding(ctx, binding)

	return nil
}

// RefundDeposit refunds the deposit from the specified service binding
func (k Keeper) RefundDeposit(ctx sdk.Context, serviceName string, provider, owner sdk.AccAddress) error {
	binding, found := k.GetServiceBinding(ctx, serviceName, provider)
	if !found {
		return sdkerrors.Wrap(types.ErrUnknownServiceBinding, "")
	}

	bindingOwner, err := sdk.AccAddressFromBech32(binding.Owner)
	if err != nil {
		return err
	}

	if !owner.Equals(bindingOwner) {
		return sdkerrors.Wrap(types.ErrNotAuthorized, "owner not matching")
	}

	if binding.Available {
		return sdkerrors.Wrap(types.ErrServiceBindingAvailable, "")
	}

	if binding.Deposit.IsZero() {
		return sdkerrors.Wrap(types.ErrInvalidDeposit, "the deposit of the service binding is zero")
	}

	refundableTime := binding.DisabledTime.Add(k.ArbitrationTimeLimit(ctx)).Add(k.ComplaintRetrospect(ctx))

	currentTime := ctx.BlockHeader().Time
	if currentTime.Before(refundableTime) {
		return sdkerrors.Wrapf(types.ErrIncorrectRefundTime, "%v", refundableTime)
	}

	// Send coins from the deposit module account to the owner's account
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx, types.DepositAccName, bindingOwner, binding.Deposit,
	); err != nil {
		return err
	}

	binding.Deposit = sdk.Coins{}
	k.SetServiceBinding(ctx, binding)

	return nil
}

// RefundDeposits refunds the deposits of all the service bindings
func (k Keeper) RefundDeposits(ctx sdk.Context) error {
	iterator := k.AllServiceBindingsIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var binding types.ServiceBinding
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &binding)

		bindingOwner, err := sdk.AccAddressFromBech32(binding.Owner)
		if err != nil {
			return err
		}

		if err := k.bankKeeper.SendCoinsFromModuleToAccount(
			ctx, types.DepositAccName, bindingOwner, binding.Deposit,
		); err != nil {
			return err
		}
	}

	return nil
}

// SetServiceBinding sets the service binding
func (k Keeper) SetServiceBinding(ctx sdk.Context, svcBinding types.ServiceBinding) {
	store := ctx.KVStore(k.storeKey)

	provider, _ := sdk.AccAddressFromBech32(svcBinding.Provider)
	bz := k.cdc.MustMarshalBinaryBare(&svcBinding)
	store.Set(types.GetServiceBindingKey(svcBinding.ServiceName, provider), bz)
}

// GetServiceBinding retrieves the specified service binding
func (k Keeper) GetServiceBinding(
	ctx sdk.Context, serviceName string, provider sdk.AccAddress,
) (
	svcBinding types.ServiceBinding, found bool,
) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetServiceBindingKey(serviceName, provider))
	if bz == nil {
		return svcBinding, false
	}

	k.cdc.MustUnmarshalBinaryBare(bz, &svcBinding)
	return svcBinding, true
}

// SetOwnerServiceBinding sets the owner service binding
func (k Keeper) SetOwnerServiceBinding(ctx sdk.Context, svcBinding types.ServiceBinding) {
	store := ctx.KVStore(k.storeKey)
	owner, _ := sdk.AccAddressFromBech32(svcBinding.Owner)
	provider, _ := sdk.AccAddressFromBech32(svcBinding.Provider)
	store.Set(types.GetOwnerServiceBindingKey(owner, svcBinding.ServiceName, provider), []byte{})
}

// GetOwnerServiceBindings retrieves the service bindings with the specified service name and owner
func (k Keeper) GetOwnerServiceBindings(ctx sdk.Context, owner sdk.AccAddress, serviceName string) []*types.ServiceBinding {
	store := ctx.KVStore(k.storeKey)

	bindings := make([]*types.ServiceBinding, 0)

	iterator := sdk.KVStorePrefixIterator(store, types.GetOwnerBindingsSubspace(owner, serviceName))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		bindingKey := iterator.Key()[sdk.AddrLen+1:]
		sepIndex := bytes.Index(bindingKey, types.Delimiter)
		serviceName := string(bindingKey[0:sepIndex])
		provider := sdk.AccAddress(bindingKey[sepIndex+1:])

		if binding, found := k.GetServiceBinding(ctx, serviceName, provider); found {
			bindings = append(bindings, &binding)
		}
	}

	return bindings
}

// SetOwner sets an owner for the specified provider
func (k Keeper) SetOwner(ctx sdk.Context, provider, owner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(&gogotypes.BytesValue{Value: owner})
	store.Set(types.GetOwnerKey(provider), bz)
}

// GetOwner gets the owner for the specified provider
func (k Keeper) GetOwner(ctx sdk.Context, provider sdk.AccAddress) (sdk.AccAddress, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetOwnerKey(provider))
	if bz == nil {
		return nil, false
	}

	addr := gogotypes.BytesValue{}
	k.cdc.MustUnmarshalBinaryBare(bz, &addr)
	return addr.GetValue(), true
}

// SetOwnerProvider sets the provider with the owner
func (k Keeper) SetOwnerProvider(ctx sdk.Context, owner, provider sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetOwnerProviderKey(owner, provider), []byte{})
}

// OwnerProvidersIterator returns an iterator for all providers of the specified owner
func (k Keeper) OwnerProvidersIterator(ctx sdk.Context, owner sdk.AccAddress) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.GetOwnerProvidersSubspace(owner))
}

// SetPricing sets the pricing for the specified service binding
func (k Keeper) SetPricing(
	ctx sdk.Context,
	serviceName string,
	provider sdk.AccAddress,
	pricing types.Pricing,
) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(&pricing)
	store.Set(types.GetPricingKey(serviceName, provider), bz)
}

// GetPricing retrieves the pricing of the specified service binding
func (k Keeper) GetPricing(ctx sdk.Context, serviceName string, provider sdk.AccAddress) (pricing types.Pricing) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetPricingKey(serviceName, provider))
	if bz == nil {
		return
	}

	k.cdc.MustUnmarshalBinaryBare(bz, &pricing)
	return pricing
}

// SetWithdrawAddress sets the withdrawal address for the specified owner
func (k Keeper) SetWithdrawAddress(ctx sdk.Context, owner, withdrawAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetWithdrawAddrKey(owner), withdrawAddr.Bytes())
}

// GetWithdrawAddress gets the withdrawal address of the specified owner
func (k Keeper) GetWithdrawAddress(ctx sdk.Context, owner sdk.AccAddress) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetWithdrawAddrKey(owner))
	if bz == nil {
		return owner
	}

	return sdk.AccAddress(bz)
}

// IterateWithdrawAddresses iterates through all withdrawal addresses
func (k Keeper) IterateWithdrawAddresses(
	ctx sdk.Context,
	op func(owner sdk.AccAddress, withdrawAddress sdk.AccAddress) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.WithdrawAddrKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		ownerAddress := sdk.AccAddress(iterator.Key()[1:])
		withdrawAddress := sdk.AccAddress(iterator.Value())

		if stop := op(ownerAddress, withdrawAddress); stop {
			break
		}
	}
}

// ServiceBindingsIterator returns an iterator for all bindings of the specified service definition
func (k Keeper) ServiceBindingsIterator(ctx sdk.Context, serviceName string) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.GetBindingsSubspace(serviceName))
}

// AllServiceBindingsIterator returns an iterator for all bindings
func (k Keeper) AllServiceBindingsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.ServiceBindingKey)
}

func (k Keeper) IterateServiceBindings(
	ctx sdk.Context,
	op func(binding types.ServiceBinding) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.ServiceBindingKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var binding types.ServiceBinding
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &binding)

		if stop := op(binding); stop {
			break
		}
	}
}

// GetMinDeposit gets the minimum deposit required for the service binding
func (k Keeper) GetMinDeposit(ctx sdk.Context, pricing types.Pricing) (sdk.Coins, error) {
	minDepositMultiple := sdk.NewInt(k.MinDepositMultiple(ctx))
	minDepositParam := k.MinDeposit(ctx)
	baseDenom := k.BaseDenom(ctx)

	priceDenom := pricing.Price.GetDenomByIndex(0)
	price := pricing.Price.AmountOf(priceDenom)

	basePrice := price

	if priceDenom != baseDenom && !price.IsZero() {
		rate, err := k.GetExchangeRate(ctx, priceDenom, baseDenom)
		if err != nil {
			return nil, err
		}

		basePrice = sdk.NewDecFromInt(price).Mul(rate).TruncateInt()
		if basePrice.IsZero() {
			basePrice = sdk.OneInt()
		}
	}

	// minimum deposit = max(price * minDepositMultiple, minDepositParam)
	minDeposit := sdk.NewCoins(sdk.NewCoin(baseDenom, basePrice.Mul(minDepositMultiple)))
	if !minDeposit.IsZero() && minDeposit.IsAllLT(minDepositParam) {
		minDeposit = minDepositParam
	}

	return minDeposit, nil
}

// ParsePricing parses the given pricing
func (k Keeper) ParsePricing(ctx sdk.Context, pricing string) (p types.Pricing, err error) {
	p, err = types.ParsePricing(pricing)
	if err != nil {
		return p, err
	}

	if err := types.CheckPricing(p); err != nil {
		return p, err
	}

	if err := k.validatePricing(ctx, p); err != nil {
		return p, err
	}

	return p, nil
}

// validateDeposit validates the given deposit
func (k Keeper) validateDeposit(ctx sdk.Context, deposit sdk.Coins) error {
	baseDenom := k.BaseDenom(ctx)

	if len(deposit) != 1 || deposit[0].Denom != baseDenom {
		return sdkerrors.Wrapf(types.ErrInvalidDeposit, "deposit only accepts %s", baseDenom)
	}

	return nil
}

// validatePricing validates the given pricing
func (k Keeper) validatePricing(ctx sdk.Context, pricing types.Pricing) error {
	priceDenom := pricing.Price.GetDenomByIndex(0)

	if k.RestrictedServiceFeeDenom(ctx) {
		baseDenom := k.BaseDenom(ctx)

		if priceDenom != baseDenom {
			return sdkerrors.Wrapf(types.ErrInvalidPricing, "invalid denom: %s, service fee only accepts %s", priceDenom, baseDenom)
		}
	}

	if supply := k.bankKeeper.GetSupply(ctx).GetTotal().AmountOf(priceDenom); !supply.IsPositive() {
		return sdkerrors.Wrapf(types.ErrInvalidPricing, "invalid denom: %s", priceDenom)
	}

	return nil
}
