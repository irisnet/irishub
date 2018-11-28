package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/upgrade"
)

type Hook interface{}

const stakeTrigger = "stake"

type hooks []Hook
type hooksVersion map[int64]hooks

type HookHub struct {
	upgradeKeeper upgrade.Keeper
	triggeredHook map[string]hooksVersion
}

var _ sdk.StakingHooks = HookHub{}
var _ Hook = HookHub{}

func NewHooksHub(keeper upgrade.Keeper) HookHub {
	return HookHub{
		upgradeKeeper: keeper,
		triggeredHook: make(map[string]hooksVersion),
	}
}

func (hkhub HookHub) AddHook(trigger string, version int64, hk Hook) (hkh HookHub) {
	hkversion, ok := hkhub.triggeredHook[trigger]
	if !ok {
		hkversion = make(map[int64]hooks)
	}

	hkversion[version] = []Hook(append(hkversion[version], hk))
	hkhub.triggeredHook[trigger] = hkversion

	return hkhub
}

func (hkhub HookHub) GetHooks(trigger string, version int64) (hk hooks) {
	hkversion, ok := hkhub.triggeredHook[trigger]
	if !ok {
		return nil
	}

	hks, ok := hkversion[version]
	if !ok {
		return nil
	}

	return hks
}

func (h HookHub) GetCurrentVersionHooks(ctx sdk.Context, trigger string) hooks {
	version := h.upgradeKeeper.GetCurrentVersion(ctx)

	hookversion, ok := h.triggeredHook[trigger]
	if !ok {
		panic("The stakeTrigger of the hookHub doesn't existed!")
	}

	return hookversion[version.Id]
}

//______________________________________________________________________________________________
// functions for StakingHooks
func (h HookHub) OnValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) {
	hks := h.GetCurrentVersionHooks(ctx, stakeTrigger)

	for _, hook := range hks {
		hook.(sdk.StakingHooks).OnValidatorCreated(ctx, valAddr)
	}
}
func (h HookHub) OnValidatorModified(ctx sdk.Context, valAddr sdk.ValAddress) {
	hks := h.GetCurrentVersionHooks(ctx, stakeTrigger)

	for _, hook := range hks {
		hook.(sdk.StakingHooks).OnValidatorModified(ctx, valAddr)
	}
}

func (h HookHub) OnValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	hks := h.GetCurrentVersionHooks(ctx, stakeTrigger)

	for _, hook := range hks {
		hook.(sdk.StakingHooks).OnValidatorRemoved(ctx, consAddr, valAddr)
	}
}

func (h HookHub) OnValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	hks := h.GetCurrentVersionHooks(ctx, stakeTrigger)

	for _, hook := range hks {
		hook.(sdk.StakingHooks).OnValidatorBonded(ctx, consAddr, valAddr)
	}
}

func (h HookHub) OnValidatorPowerDidChange(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	hks := h.GetCurrentVersionHooks(ctx, stakeTrigger)

	for _, hook := range hks {
		hook.(sdk.StakingHooks).OnValidatorPowerDidChange(ctx, consAddr, valAddr)
	}
}

func (h HookHub) OnValidatorBeginUnbonding(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	hks := h.GetCurrentVersionHooks(ctx, stakeTrigger)

	for _, hook := range hks {
		hook.(sdk.StakingHooks).OnValidatorBeginUnbonding(ctx, consAddr, valAddr)
	}
}

func (h HookHub) OnDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	hks := h.GetCurrentVersionHooks(ctx, stakeTrigger)

	for _, hook := range hks {
		hook.(sdk.StakingHooks).OnDelegationCreated(ctx, delAddr, valAddr)
	}
}

func (h HookHub) OnDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	hks := h.GetCurrentVersionHooks(ctx, stakeTrigger)

	for _, hook := range hks {
		hook.(sdk.StakingHooks).OnDelegationSharesModified(ctx, delAddr, valAddr)
	}
}

func (h HookHub) OnDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	hks := h.GetCurrentVersionHooks(ctx, stakeTrigger)

	for _, hook := range hks {
		hook.(sdk.StakingHooks).OnDelegationRemoved(ctx, delAddr, valAddr)
	}
}
