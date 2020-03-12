package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeeper(t *testing.T) {
	require.Equal(t, false, isValidVersion(1, 1, 1))
	require.Equal(t, true, isValidVersion(1, 4, 4))
	require.Equal(t, true, isValidVersion(1, 4, 5))
	require.Equal(t, true, isValidVersion(1, 1, 2))
	require.Equal(t, true, isValidVersion(2, 1, 3))
}
