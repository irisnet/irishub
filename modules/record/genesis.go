package record

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/record/keeper"
	"github.com/irisnet/irismod/modules/record/types"
)

// InitGenesis stores genesis data
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	for _, record := range data.Records {
		k.AddRecord(ctx, record)
	}
}

// ExportGenesis outputs genesis data
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	recordsIterator := k.RecordsIterator(ctx)
	defer recordsIterator.Close()

	var records []types.Record
	for ; recordsIterator.Valid(); recordsIterator.Next() {
		var record types.Record
		types.ModuleCdc.MustUnmarshalBinaryBare(recordsIterator.Value(), &record)
		records = append(records, record)
	}

	return types.NewGenesisState(records)
}

// ValidateGenesis validates the provided record genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data types.GenesisState) error {
	for _, record := range data.Records {
		if len(record.Contents) == 0 {
			return fmt.Errorf("contents missing")
		}
		if record.Creator.Empty() {
			return fmt.Errorf("Creator missing")
		}
		for _, content := range record.Contents {
			if len(content.Digest) == 0 {
				return fmt.Errorf("Digest missing")
			}
			if len(content.DigestAlgo) == 0 {
				return fmt.Errorf("DigestAlgo missing")
			}
		}
	}
	return nil
}
