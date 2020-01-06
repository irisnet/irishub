//nolint:bodyclose
package lcdtest

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestCoinswap(t *testing.T) {
	cleanup, _, _, _, err := InitializeLCD(1, []sdk.AccAddress{}, true)
	require.NoError(t, err)
	defer cleanup()
}
