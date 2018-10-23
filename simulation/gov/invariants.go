package simulation

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/irisnet/irishub/baseapp"
	"github.com/irisnet/irishub/simulation/mock/simulation"
)

// AllInvariants tests all governance invariants
func AllInvariants() simulation.Invariant {
	return func(t *testing.T, app *baseapp.BaseApp, log string) {
		// TODO Add some invariants!
		// Checking proposal queues, no passed-but-unexecuted proposals, etc.
		require.Nil(t, nil)
	}
}
