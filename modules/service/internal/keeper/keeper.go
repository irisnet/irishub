package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
	"github.com/irisnet/irishub/modules/service/internal/types"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
	sk       types.SupplyKeeper
	gk       types.GuardianKeeper

	// codespace
	codespace sdk.CodespaceType
	// params subspace
	paramSpace params.Subspace
	// metrics
	metrics *types.Metrics
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, sk types.SupplyKeeper, gk types.GuardianKeeper, codespace sdk.CodespaceType, paramSpace params.Subspace, metrics *types.Metrics) Keeper {
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
		codespace:  codespace,
		paramSpace: paramSpace.WithKeyTable(ParamKeyTable()),
		metrics:    metrics,
	}

	return keeper
}

// return the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
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
