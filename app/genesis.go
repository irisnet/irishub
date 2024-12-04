package app

import (
	"github.com/irisnet/irishub/v4/types"
)

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState() types.GenesisState {
	encCfg := RegisterEncodingConfig()
	return ModuleBasics.DefaultGenesis(encCfg.Marshaler)
}
