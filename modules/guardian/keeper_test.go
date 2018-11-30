package guardian

import (
	"testing"
	"github.com/stretchr/testify/require"
)

func TestKeeper_AddProfiler(t *testing.T) {
	ctx, keeper := createTestInput(t)
	profiler := NewProfiler(addrs[0], addrs[1])
	keeper.AddProfiler(ctx, profiler)
	AddedProfiler, found := keeper.GetProfiler(ctx, addrs[0])
	require.True(t, found)
	require.True(t, ProfilerEqual(profiler, AddedProfiler))
}
