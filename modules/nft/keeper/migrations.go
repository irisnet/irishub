package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v2 "github.com/irisnet/irismod/modules/nft/migrations/v2"
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
	return v2.Migrate(ctx, m.k.storeKey, m.k.cdc, m.k.Logger(ctx), m.k)
}
