package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v2 "mods.irisnet.org/modules/coinswap/migrations/v2"
	v3 "mods.irisnet.org/modules/coinswap/migrations/v3"
	v4 "mods.irisnet.org/modules/coinswap/migrations/v4"
	v5 "mods.irisnet.org/modules/coinswap/migrations/v5"
	"mods.irisnet.org/modules/coinswap/types"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	k              Keeper
	legacySubspace types.Subspace
}

// NewMigrator returns a new Migrator.
func NewMigrator(k Keeper, legacySubspace types.Subspace) Migrator {
	return Migrator{k: k, legacySubspace: legacySubspace}
}

// Migrate1to2 migrates from version 1 to 2.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return v2.Migrate(ctx, m.k, m.k.bk, m.k.ak)
}

// Migrate1to2 migrates from version 2 to 3.
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	return v3.Migrate(ctx, m.k, m.legacySubspace)
}

// Migrate1to2 migrates from version 3 to 4.
func (m Migrator) Migrate3to4(ctx sdk.Context) error {
	return v4.Migrate(ctx, m.k, m.legacySubspace)
}

// Migrate1to2 migrates from version 4 to 5.
func (m Migrator) Migrate4to5(ctx sdk.Context) error {
	return v5.Migrate(ctx, m.k, m.legacySubspace)
}
