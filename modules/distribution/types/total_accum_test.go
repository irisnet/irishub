package types

import (
	"testing"

	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
)

func TestTotalAccumUpdateForNewHeight(t *testing.T) {

	ta := NewTotalAccum(0)

	ta = ta.UpdateForNewHeight(log.NewNopLogger(), 5, sdk.NewDec(3))
	require.True(sdk.DecEq(t, sdk.NewDec(15), ta.Accum))

	ta = ta.UpdateForNewHeight(log.NewNopLogger(), 8, sdk.NewDec(2))
	require.True(sdk.DecEq(t, sdk.NewDec(21), ta.Accum))
}
