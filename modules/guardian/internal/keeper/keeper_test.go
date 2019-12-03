package keeper_test

import (
	"github.com/irisnet/irishub/modules/guardian/internal/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestKeeper(t *testing.T) {
	app, ctx := createTestApp(false)
	profiler := types.NewGuardian("test", types.Genesis, addrs[0], addrs[1])

	keeper := app.GuardianKeeper
	cdc := app.Codec()

	keeper.AddProfiler(ctx, profiler)
	AddedProfiler, found := keeper.GetProfiler(ctx, addrs[0])
	require.True(t, found)
	require.True(t, profiler.Equal(AddedProfiler))

	trustee := types.NewGuardian("test", types.Genesis, addrs[0], addrs[1])
	keeper.AddTrustee(ctx, trustee)
	AddedTrustee, found := keeper.GetTrustee(ctx, addrs[0])
	require.True(t, found)
	require.True(t, trustee.Equal(AddedTrustee))

	profilersIterator := keeper.ProfilersIterator(ctx)
	defer profilersIterator.Close()
	var profilers []types.Guardian
	for ; profilersIterator.Valid(); profilersIterator.Next() {
		var profiler types.Guardian
		cdc.MustUnmarshalBinaryLengthPrefixed(profilersIterator.Value(), &profiler)
		profilers = append(profilers, profiler)
	}
	require.Equal(t, 1, len(profilers))
	require.True(t, profiler.Equal(profilers[0]))

	trusteesIterator := keeper.TrusteesIterator(ctx)
	defer trusteesIterator.Close()
	var trustees []types.Guardian
	for ; trusteesIterator.Valid(); trusteesIterator.Next() {
		var trustee types.Guardian
		cdc.MustUnmarshalBinaryLengthPrefixed(trusteesIterator.Value(), &trustee)
		trustees = append(trustees, trustee)
	}
	require.Equal(t, 1, len(trustees))
	require.True(t, trustee.Equal(trustees[0]))
}
