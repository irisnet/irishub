package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v2 "mods.irisnet.org/service/migrations/v2"
	"mods.irisnet.org/service/types"
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
	return v2.Migrate(ctx, m.k, m.legacySubspace)
}
