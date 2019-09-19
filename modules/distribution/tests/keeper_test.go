package tests

import (
	"testing"

	"github.com/irisnet/irishub/modules/distribution/types"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
)

func TestSetGetPreviousProposerConsAddr(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInputDefault(t, false, sdk.ZeroInt())

	keeper.SetPreviousProposerConsAddr(ctx, valConsAddr1)
	res := keeper.GetPreviousProposerConsAddr(ctx)
	require.True(t, res.Equals(valConsAddr1), "expected: %v got: %v", valConsAddr1.String(), res.String())
}

func TestSetGetFeePool(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInputDefault(t, false, sdk.ZeroInt())

	fp := types.InitialFeePool()
	fp.TotalValAccum.UpdateHeight = 777

	keeper.SetFeePool(ctx, fp)
	res := keeper.GetFeePool(ctx)
	require.Equal(t, fp.TotalValAccum, res.TotalValAccum)
}
