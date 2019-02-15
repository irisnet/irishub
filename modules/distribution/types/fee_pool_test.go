package types

import (
	"testing"

	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
)

func TestUpdateTotalValAccum(t *testing.T) {

	fp := InitialFeePool()

	fp = fp.UpdateTotalValAccum(log.NewNopLogger(), 5, sdk.NewDec(3))
	require.True(sdk.DecEq(t, sdk.NewDec(15), fp.TotalValAccum.Accum))

	fp = fp.UpdateTotalValAccum(log.NewNopLogger(), 8, sdk.NewDec(2))
	require.True(sdk.DecEq(t, sdk.NewDec(21), fp.TotalValAccum.Accum))
}
