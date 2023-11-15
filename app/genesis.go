package app

import (
	"github.com/irisnet/irishub/v2/types"
)

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState() types.GenesisState {
	encCfg := MakeEncodingConfig()
	return ModuleBasics.DefaultGenesis(encCfg.Marshaler)
}
