package simulation

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/irisnet/irishub/baseapp"
	"github.com/irisnet/irishub/simulation/mock/simulation"
)

// TODO Any invariants to check here?
// AllInvariants tests all slashing invariants
func AllInvariants() simulation.Invariant {
	return func(_ *baseapp.BaseApp, _ abci.Header) error {
		return nil
	}
}
