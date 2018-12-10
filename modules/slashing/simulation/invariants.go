package simulation

import (
	"github.com/irisnet/irishub/modules/mock/simulation"
	sdk "github.com/irisnet/irishub/types"
)

// TODO Any invariants to check here?
// AllInvariants tests all slashing invariants
func AllInvariants() simulation.Invariant {
	return func(_ sdk.Context) error {
		return nil
	}
}
