package record

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/record/keeper"
	"github.com/irisnet/irismod/modules/record/types"
)

// InitGenesis stores genesis data
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	if err := types.ValidateGenesis(data); err != nil {
		panic(err.Error())
	}

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
