package config

import (
	"testing"

	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	require.True(t, cfg.MinimumFees().IsZero())
}

func TestSetMinimumFees(t *testing.T) {
	cfg := DefaultConfig()
	cfg.SetMinimumFees(sdk.Coins{sdk.NewCoin("foo", sdk.NewInt(100))})
	require.Equal(t, "100foo", cfg.MinFees)

	cfg.SetInvariant("error")
	require.Equal(t, "error", cfg.InvariantLevel)


	cfg.SetInvariant("panic")
	require.Equal(t, "panic", cfg.InvariantLevel)

	cfg.SetInvariant("no1")
	require.Equal(t, "", cfg.InvariantLevel)
}
