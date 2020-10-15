package keeper

import (
	"bytes"
	"encoding/json"
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
			"qos [%d] must not be greater than maximum request timeout [%d]",
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

	if err := types.ValidatePricing(parsedPricing); err != nil {
		return err
	}

	for i, token := range parsedPricing.Price {
		if total := k.bankKeeper.GetSupply(ctx).GetTotal().AmountOf(token.Denom); !total.IsPositive() {
			return sdkerrors.Wrapf(types.ErrInvalidPricing, "invalid denom: %s", parsedPricing.Price[i].Denom)
		}
	}

	minDeposit := k.getMinDeposit(ctx, parsedPricing)
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
	k.SetServiceBinding(ctx, svcBinding)
	k.SetOwnerServiceBinding(ctx, svcBinding)
	k.SetOwner(ctx, svcBinding.Provider, svcBinding.Owner)
	k.SetOwnerProvider(ctx, svcBinding.Owner, svcBinding.Provider)

	pricing, err := k.ParsePricing(ctx, svcBinding.Pricing)
	if err != nil {
		return err
	}

	k.SetPricing(ctx, svcBinding.ServiceName, svcBinding.Provider, pricing)

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

	if !owner.Equals(binding.Owner) {
		return sdkerrors.Wrap(types.ErrNotAuthorized, "owner not matching")
	}

	updated := false

	if qos != 0 {
		maxReqTimeout := k.MaxRequestTimeout(ctx)
		if qos > uint64(maxReqTimeout) {
			return sdkerrors.Wrapf(
				types.ErrInvalidQoS,
				"qos [%d] must not be greater than maximum request timeout [%d]",
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

		if err := types.ValidatePricing(parsedPricing); err != nil {
			return err
		}

		for i, token := range parsedPricing.Price {
			if total := k.bankKeeper.GetSupply(ctx).GetTotal().AmountOf(token.Denom); !total.IsPositive() {
				return sdkerrors.Wrapf(types.ErrInvalidPricing, "invalid denom: %s", parsedPricing.Price[i].Denom)
			}
		}

		binding.Pricing = pricing
		k.SetPricing(ctx, serviceName, provider, parsedPricing)

		updated = true
	}

	// only check deposit when the binding is available and updated
	if binding.Available && updated {
		minDeposit := k.getMinDeposit(ctx, parsedPricing)
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

	if !owner.Equals(binding.Owner) {
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

	if !owner.Equals(binding.Owner) {
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

	minDeposit := k.getMinDeposit(ctx, k.GetPricing(ctx, serviceName, provider))
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

	if !owner.Equals(binding.Owner) {
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
		ctx, types.DepositAccName, binding.Owner, binding.Deposit,
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

		if err := k.bankKeeper.SendCoinsFromModuleToAccount(
			ctx, types.DepositAccName, binding.Owner, binding.Deposit,
		); err != nil {
			return err
		}
	}

	return nil
}

// SetServiceBinding sets the service binding
func (k Keeper) SetServiceBinding(ctx sdk.Context, svcBinding types.ServiceBinding) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(&svcBinding)
	store.Set(types.GetServiceBindingKey(svcBinding.ServiceName, svcBinding.Provider), bz)
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
	store.Set(types.GetOwnerServiceBindingKey(svcBinding.Owner, svcBinding.ServiceName, svcBinding.Provider), []byte{})
}

// GetOwnerServiceBindings retrieves the service bindings with the specified service name and owner
func (k Keeper) GetOwnerServiceBindings(ctx sdk.Context, owner sdk.AccAddress, serviceName string) []*types.ServiceBinding {
	store := ctx.KVStore(k.storeKey)

	bindings := make([]*types.ServiceBinding, 0)

	iterator := sdk.KVStorePrefixIterator(store, types.GetOwnerBindingsSubspace(owner, serviceName))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		bindingKey := iterator.Key()[sdk.AddrLen+1:]
		sepIndex := bytes.Index(bindingKey, types.EmptyByte)
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

// ParsePricing parses the given string to Pricing
func (k Keeper) ParsePricing(ctx sdk.Context, pricing string) (p types.Pricing, err error) {
	var rawPricing types.RawPricing
	if err := json.Unmarshal([]byte(pricing), &rawPricing); err != nil {
		return p, sdkerrors.Wrapf(types.ErrInvalidPricing, "failed to unmarshal the pricing: %s", err.Error())
	}

	tokenPrice, err := sdk.ParseCoin(rawPricing.Price)
	if err != nil {
		return p, sdkerrors.Wrapf(types.ErrInvalidPricing, "invalid price: %s", err.Error())
	}

	p.Price = sdk.Coins{tokenPrice}
	p.PromotionsByTime = rawPricing.PromotionsByTime
	p.PromotionsByVolume = rawPricing.PromotionsByVolume

	return p, nil
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

// getMinDeposit gets the minimum deposit required for the service binding
func (k Keeper) getMinDeposit(ctx sdk.Context, pricing types.Pricing) sdk.Coins {
	minDepositMultiple := sdk.NewInt(k.MinDepositMultiple(ctx))
	minDepositParam := k.MinDeposit(ctx)
	baseDenom := k.BaseDenom(ctx)

	price := pricing.Price.AmountOf(baseDenom)

	// minimum deposit = max(price * minDepositMultiple, minDepositParam)
	minDeposit := sdk.NewCoins(sdk.NewCoin(baseDenom, price.Mul(minDepositMultiple)))
	if minDeposit.IsAllLT(minDepositParam) {
		minDeposit = minDepositParam
	}

	return minDeposit
}

// validateDeposit validates the given deposit
func (k Keeper) validateDeposit(ctx sdk.Context, deposit sdk.Coins) error {
	baseDenom := k.BaseDenom(ctx)

	if len(deposit) != 1 || deposit[0].Denom != baseDenom {
		return sdkerrors.Wrapf(types.ErrInvalidDeposit, "deposit only accepts %s", baseDenom)
	}

	return nil
}
