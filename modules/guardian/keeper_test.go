package guardian

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestKeeper(t *testing.T) {
	ctx, keeper := createTestInput(t)
	profiler := NewGuardian("test", Genesis, addrs[0], addrs[1])

	keeper.AddProfiler(ctx, profiler)
	AddedProfiler, found := keeper.GetProfiler(ctx, addrs[0])
	require.True(t, found)
	require.True(t, profiler.Equal(AddedProfiler))

	trustee := NewGuardian("test", Genesis, addrs[0], addrs[1])
	keeper.AddTrustee(ctx, trustee)
	AddedTrustee, found := keeper.GetTrustee(ctx, addrs[0])
	require.True(t, found)
	require.True(t, trustee.Equal(AddedTrustee))

	profilersIterator := keeper.ProfilersIterator(ctx)
	defer profilersIterator.Close()
	var profilers []Guardian
	for ; profilersIterator.Valid(); profilersIterator.Next() {
		var profiler Guardian
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(profilersIterator.Value(), &profiler)
		profilers = append(profilers, profiler)
	}
	require.Equal(t, 1, len(profilers))
	require.True(t, profiler.Equal(profilers[0]))

	trusteesIterator := keeper.TrusteesIterator(ctx)
	defer trusteesIterator.Close()
	var trustees []Guardian
	for ; trusteesIterator.Valid(); trusteesIterator.Next() {
		var trustee Guardian
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(trusteesIterator.Value(), &trustee)
		trustees = append(trustees, trustee)
	}
	require.Equal(t, 1, len(trustees))
	require.True(t, trustee.Equal(trustees[0]))
}
