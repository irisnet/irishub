package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/upgrade"
)

const trigger = "stake"
var _ sdk.StakingHooks = HooksHub{}

type hooks []sdk.StakingHooks
type hooksVersion map[int64]hooks

type HooksHub struct {
	upgradeKeeper upgrade.Keeper
	triggeredHook map[string]hooksVersion
}

func NewHooksHub(keeper upgrade.Keeper) HooksHub {
	return HooksHub{
		upgradeKeeper: keeper,
		triggeredHook: make(map[string]hooksVersion),
	}
}

func (hkhub *HooksHub) AddHook(trigger string, version int64, shk sdk.StakingHooks) (hkh *HooksHub)  {
	hkversion, ok := hkhub.triggeredHook[trigger]
	if !ok {
		hkversion = make(map[int64]hooks)
	}

	hkversion[version] = []sdk.StakingHooks(append(hkversion[version], shk))
	hkhub.triggeredHook[trigger] = hkversion

	return hkhub
}

func (hkhub *HooksHub) GetHooks(trigger string, version int64) (hk *hooks) {
	hkversion, ok := hkhub.triggeredHook[trigger]
	if !ok {
		return nil
	}

	hks, ok := hkversion[version]
	if !ok {
		return nil
	}

	return &hks
}

func (h HooksHub) GetCurrentVersionHooks(ctx sdk.Context) hooks {
	version := h.upgradeKeeper.GetCurrentVersion(ctx)

	hookversion, ok := h.triggeredHook[trigger]
	if !ok {
		panic("The trigger of the hooks doesn't existed!")
	}

	return hookversion[version.Id]
}

func (h HooksHub) OnValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) {
	hks := h.GetCurrentVersionHooks(ctx)

	for _, hook := range hks {
		hook.OnValidatorCreated(ctx, valAddr)
	}
}
func (h HooksHub) OnValidatorModified(ctx sdk.Context, valAddr sdk.ValAddress) {
	hks := h.GetCurrentVersionHooks(ctx)

	for _, hook := range hks {
		hook.OnValidatorModified(ctx, valAddr)
	}
}

func (h HooksHub) OnValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	hks := h.GetCurrentVersionHooks(ctx)

	for _, hook := range hks {
		hook.OnValidatorRemoved(ctx, consAddr, valAddr)
	}
}

func (h HooksHub) OnValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	hks := h.GetCurrentVersionHooks(ctx)

	for _, hook := range hks {
		hook.OnValidatorBonded(ctx, consAddr, valAddr)
	}
}

func (h HooksHub) OnValidatorPowerDidChange(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	hks := h.GetCurrentVersionHooks(ctx)

	for _, hook := range hks {
		hook.OnValidatorPowerDidChange(ctx, consAddr, valAddr)
	}
}

func (h HooksHub) OnValidatorBeginUnbonding(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	hks := h.GetCurrentVersionHooks(ctx)

	for _, hook := range hks {
		hook.OnValidatorBeginUnbonding(ctx, consAddr, valAddr)
	}
}

func (h HooksHub) OnDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	hks := h.GetCurrentVersionHooks(ctx)

	for _, hook := range hks {
		hook.OnDelegationCreated(ctx, delAddr, valAddr)
	}
}

func (h HooksHub) OnDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	hks := h.GetCurrentVersionHooks(ctx)

	for _, hook := range hks {
		hook.OnDelegationSharesModified(ctx, delAddr, valAddr)
	}
}

func (h HooksHub) OnDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	hks := h.GetCurrentVersionHooks(ctx)

	for _, hook := range hks {
		hook.OnDelegationRemoved(ctx, delAddr, valAddr)
	}
}

