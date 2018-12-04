package simulation

import (
	"github.com/irisnet/irishub/baseapp"
	"github.com/irisnet/irishub/modules/mock/simulation"
)

// TODO Any invariants to check here?
// AllInvariants tests all slashing invariants
func AllInvariants() simulation.Invariant {
	return func(_ *baseapp.BaseApp) error {
		return nil
	}
}
