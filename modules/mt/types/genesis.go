package types

// NewGenesisState creates a new genesis state.
func NewGenesisState(collections []Collection) *GenesisState {
	return &GenesisState{
		Collections: collections,
	}
}

// ValidateGenesis performs basic validation of mts genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	for _, c := range data.Collections {
		if err := ValidateDenomID(c.Denom.Name); err != nil {
			return err
		}

		for _, mt := range c.Mts {
			if err := ValidateMTID(mt.GetID()); err != nil {
				return err
			}
		}
	}
	return nil
}
