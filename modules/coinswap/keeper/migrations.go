package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v150 "github.com/irisnet/irismod/modules/coinswap/migrations/v150"
	v152 "github.com/irisnet/irismod/modules/coinswap/migrations/v152"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	k Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(k Keeper) Migrator {
	return Migrator{k: k}
}

// Migrate1to2 migrates from version 1 to 2.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return v150.Migrate(ctx, m.k, m.k.bk, m.k.ak)
}

// Migrate2to3 migrates from version 2 to 3.
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	return v152.Migrate(ctx, m.k)
}
