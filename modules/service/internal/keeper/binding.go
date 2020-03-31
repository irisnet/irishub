package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irishub/modules/service/internal/types"
)

// AddServiceBinding
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
) error {
	if _, found := k.GetServiceDefinition(ctx, defName); !found {
		return sdkerrors.Wrapf(types.ErrUnknownSvcDef, "define chain-id: %s, name: %s", defChainID, defName)
	}

	if _, found := k.GetServiceBinding(ctx, defChainID, defName, bindChainID, provider); found {
		return types.ErrSvcBindingExists
	}

	minDeposit, err := k.getMinDeposit(ctx, prices)
	if err != nil {
		return err
	}

	if !deposit.IsAllGTE(minDeposit) {
		return sdkerrors.Wrapf(types.ErrLtMinProviderDeposit, "mint deposit: %s, deposit: %s", minDeposit.String(), deposit.String())
	}

	svcBinding := types.NewSvcBinding(ctx, defChainID, defName, bindChainID, provider, bindingType, deposit, prices, level, true)

	// Send coins from provider's account to the deposit module account
	if err := k.sk.SendCoinsFromAccountToModule(
		ctx, svcBinding.Provider, types.DepositAccName, svcBinding.Deposit,
	); err != nil {
		return err
	}

	svcBinding.DisableTime = time.Time{}
	k.SetServiceBinding(ctx, svcBinding)

	return nil
}

// SetServiceBinding
func (k Keeper) SetServiceBinding(ctx sdk.Context, svcBinding types.SvcBinding) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(svcBinding)
	store.Set(types.GetServiceBindingKey(svcBinding.DefChainID, svcBinding.DefName, svcBinding.BindChainID, svcBinding.Provider), bz)
}

// GetServiceBinding
func (k Keeper) GetServiceBinding(ctx sdk.Context, defChainID, defName, bindChainID string, provider sdk.AccAddress) (svcBinding types.SvcBinding, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetServiceBindingKey(defChainID, defName, bindChainID, provider))
	if bz == nil {
		return svcBinding, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &svcBinding)
	return svcBinding, true
}

// ServiceBindingsIterator
func (k Keeper) ServiceBindingsIterator(ctx sdk.Context, defChainID, defName string) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.GetBindingsSubspaceKey(defChainID, defName))
}

// AllServiceBindingsIterator returns an iterator for all the binding services
func (k Keeper) AllServiceBindingsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.BindingPropertyKey)
}

// UpdateServiceBinding
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
) (svcBinding types.SvcBinding, err error) {
	oldBinding, found := k.GetServiceBinding(ctx, defChainID, defName, bindChainID, provider)
	if !found {
		return svcBinding, types.ErrUnknownSvcBinding
	}

	newBinding := types.NewSvcBinding(ctx, defChainID, defName, bindChainID, provider, bindingType,
		deposit, prices, level, false)

	// TODO
	oldBinding.Prices = newBinding.Prices

	if newBinding.BindingType != 0x00 {
		oldBinding.BindingType = newBinding.BindingType
	}

	// Add coins to svcBinding deposit
	if !newBinding.Deposit.IsAnyNegative() {
		oldBinding.Deposit = oldBinding.Deposit.Add(newBinding.Deposit...)
	}

	// Send coins from provider's account to the deposit module account
	if err := k.sk.SendCoinsFromAccountToModule(
		ctx, provider, types.DepositAccName, newBinding.Deposit,
	); err != nil {
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
			return svcBinding, sdkerrors.Wrapf(types.ErrLtMinProviderDeposit, "mint deposit: %s, deposit: %s",
				minDeposit.String(), oldBinding.Deposit.String())
		}
	}

	k.SetServiceBinding(ctx, oldBinding)

	return oldBinding, nil
}

// Disable
func (k Keeper) Disable(ctx sdk.Context, defChainID, defName, bindChainID string, provider sdk.AccAddress) error {
	binding, found := k.GetServiceBinding(ctx, defChainID, defName, bindChainID, provider)
	if !found {
		return types.ErrUnknownSvcBinding
	}

	if !binding.Available {
		return types.ErrUnavailable
	}

	binding.Available = false
	binding.DisableTime = ctx.BlockHeader().Time

	k.SetServiceBinding(ctx, binding)

	return nil
}

// Enable
func (k Keeper) Enable(ctx sdk.Context, defChainID, defName, bindChainID string, provider sdk.AccAddress, deposit sdk.Coins) error {
	binding, found := k.GetServiceBinding(ctx, defChainID, defName, bindChainID, provider)
	if !found {
		return types.ErrUnknownSvcBinding
	}

	if binding.Available {
		return types.ErrAvailable
	}

	// Add coins to svcBinding deposit
	if !deposit.IsAnyNegative() {
		binding.Deposit = binding.Deposit.Add(deposit...)
	}

	minDeposit, err := k.getMinDeposit(ctx, binding.Prices)
	if err != nil {
		return err
	}

	if !binding.Deposit.IsAllGTE(minDeposit) {
		return sdkerrors.Wrapf(types.ErrLtMinProviderDeposit, "mint deposit: %s, deposit: %s", minDeposit.String(), binding.Deposit.String())
	}

	// Send coins from provider's account to the deposit module account
	if err := k.sk.SendCoinsFromAccountToModule(
		ctx, binding.Provider, types.DepositAccName, deposit,
	); err != nil {
		return err
	}

	binding.Available = true
	binding.DisableTime = time.Time{}

	k.SetServiceBinding(ctx, binding)

	return nil
}

// RefundDeposit
func (k Keeper) RefundDeposit(ctx sdk.Context, defChainID, defName, bindChainID string, provider sdk.AccAddress) error {
	binding, found := k.GetServiceBinding(ctx, defChainID, defName, bindChainID, provider)
	if !found {
		return types.ErrUnknownSvcBinding
	}

	if binding.Available {
		return sdkerrors.Wrap(types.ErrAvailable, "can't refund from a available service binding")
	}

	if binding.Deposit.IsZero() {
		return sdkerrors.Wrap(types.ErrRefundDeposit, "service binding deposit is zero")
	}

	blockTime := ctx.BlockHeader().Time
	params := k.GetParams(ctx)

	refundTime := binding.DisableTime.Add(params.ArbitrationTimeLimit).Add(params.ComplaintRetrospect)
	if blockTime.Before(refundTime) {
		return sdkerrors.Wrapf(types.ErrRefundDeposit, "can not refund deposit before %s", refundTime.Format("2006-01-02 15:04:05"))
	}

	// Send coins from the deposit module account to the provider's account
	if err := k.sk.SendCoinsFromModuleToAccount(
		ctx, types.DepositAccName, binding.Provider, binding.Deposit,
	); err != nil {
		return err
	}

	binding.Deposit = sdk.Coins{}
	k.SetServiceBinding(ctx, binding)

	return nil
}

// RefundDeposits refunds the deposits of all the binding services
func (k Keeper) RefundDeposits(ctx sdk.Context) error {
	iterator := k.AllServiceBindingsIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var binding types.SvcBinding
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &binding)

		if err := k.sk.SendCoinsFromModuleToAccount(
			ctx, types.DepositAccName, binding.Provider, binding.Deposit,
		); err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) getMinDeposit(ctx sdk.Context, prices []sdk.Coin) (sdk.Coins, error) {
	params := k.GetParams(ctx)
	// min deposit must >= sum(method price) * minDepositMultiple
	minDepositMultiple := sdk.NewInt(params.MinDepositMultiple)

	var minDeposit sdk.Coins
	for _, price := range prices {
		if price.Amount.BigInt().BitLen()+minDepositMultiple.BigInt().BitLen()-1 > 255 {
			return minDeposit, types.ErrIntOverflow
		}

		minInt := price.Amount.Mul(minDepositMultiple)
		minDeposit = minDeposit.Add(sdk.NewCoin(price.Denom, minInt))
	}

	return minDeposit, nil
}
