package types

import (
	"errors"
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesisState constructs a GenesisState
func NewGenesisState(records []Record) *GenesisState {
	return &GenesisState{
		Records: records,
	}
}

// DefaultGenesisState gets raw genesis raw message for testing
func DefaultGenesisState() *GenesisState {
	return &GenesisState{}
}

// ValidateGenesis validates the provided record genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	for _, record := range data.Records {
		if len(record.Contents) == 0 {
			return errors.New("contents missing")
		}

		_, err := sdk.AccAddressFromBech32(record.Creator)
		if err != nil {
			return fmt.Errorf("invalid record creator address (%s)", err)
		}

		if err := ValidateContents(record.Contents...); err != nil {
			return nil
		}
	}
	return nil
}
