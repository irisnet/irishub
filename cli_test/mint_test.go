package clitest

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/tests"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/modules/mint"
)

func TestMintGetParams(t *testing.T) {
	t.Parallel()
	f := InitFixtures(t)

	// start iris server
	proc := f.GDStart()
	defer proc.Stop(false)

	params := f.QueryMintParams()
	require.Equal(t, sdk.NewDecWithPrec(4, 2), params.Inflation)
	require.Equal(t, sdk.DefaultBondDenom, params.MintDenom)

	// Cleanup testing directories
	f.Cleanup()
}

// QueryMintParams is iriscli query mint params
func (f *Fixtures) QueryMintParams() mint.Params {
	cmd := fmt.Sprintf("%s query mint params %s", f.IriscliBinary, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var params mint.Params
	err := cdc.UnmarshalJSON([]byte(res), &params)
	require.NoError(f.T, err)
	return params
}
