package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"

	"github.com/irisnet/irishub/modules/service/internal/types"
)

// Keeper
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
	sk       types.SupplyKeeper
	gk       types.GuardianKeeper

	// params subspace
	paramSpace params.Subspace
}

// NewKeeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, sk types.SupplyKeeper, gk types.GuardianKeeper, paramSpace params.Subspace) Keeper {
	// ensure service module accounts are set
	if addr := sk.GetModuleAddress(types.DepositAccName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.DepositAccName))
	}

	if addr := sk.GetModuleAddress(types.RequestAccName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.RequestAccName))
	}

	if addr := sk.GetModuleAddress(types.TaxAccName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.TaxAccName))
	}

	keeper := Keeper{
		storeKey:   key,
		cdc:        cdc,
		sk:         sk,
		gk:         gk,
		paramSpace: paramSpace.WithKeyTable(ParamKeyTable()),
	}

	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("%s", types.ModuleName))
}

// GetCdc returns the cdc
func (k Keeper) GetCdc() *codec.Codec {
	return k.cdc
}

// GetServiceDepositAccount returns the service depost ModuleAccount
func (k Keeper) GetServiceDepositAccount(ctx sdk.Context) exported.ModuleAccountI {
	return k.sk.GetModuleAccount(ctx, types.DepositAccName)
}

// GetServiceRequestAccount returns the service request ModuleAccount
func (k Keeper) GetServiceRequestAccount(ctx sdk.Context) exported.ModuleAccountI {
	return k.sk.GetModuleAccount(ctx, types.RequestAccName)
}

// GetServiceTaxAccount returns the service tax ModuleAccount
func (k Keeper) GetServiceTaxAccount(ctx sdk.Context) exported.ModuleAccountI {
	return k.sk.GetModuleAccount(ctx, types.TaxAccName)
}
