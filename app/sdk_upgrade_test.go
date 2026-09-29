package app_test

import (
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"crypto/sha256"
	"github.com/irisnet/irishub/v5/app/upgrades/sdk0538"
	"github.com/irisnet/irishub/v5/wrapper"
	"google.golang.org/protobuf/encoding/protowire"
	"testing"
	"time"

	nfttypes "github.com/bianjieai/nft-transfer/types"
	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	icahosttypes "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/host/types"
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	solomachine "github.com/cosmos/ibc-go/v10/modules/light-clients/06-solomachine"
	tendermint "github.com/cosmos/ibc-go/v10/modules/light-clients/07-tendermint"
	"github.com/stretchr/testify/require"

	apptestutil "github.com/irisnet/irishub/v5/testutil"
)

// Exercise the real app constructors, genesis, light-client routing and the
// first block: all changed across the SDK 0.53 / IBC v10 upgrade.
func TestSDK053AppWiring(t *testing.T) {
	app := apptestutil.CreateApp(t)
	ctx := app.NewContextLegacy(false, tmproto.Header{Height: 1})
	require.Nil(t, app.GetKey("capability"))
	require.Equal(t, authtypes.NewModuleAddress(govtypes.ModuleName).String(), app.IBCKeeper.GetAuthority())

	for _, clientID := range []string{tendermint.ModuleName + "-0", solomachine.ModuleName + "-0"} {
		client, err := app.IBCKeeper.ClientKeeper.Route(ctx, clientID)
		require.NoError(t, err, clientID)
		require.NotNil(t, client)
	}
	for _, port := range []string{transfertypes.ModuleName, nfttypes.ModuleName, icahosttypes.SubModuleName} {
		require.True(t, app.IBCKeeper.PortKeeper.Router.HasRoute(port), port)
	}
	_, err := app.FinalizeBlock(&abci.RequestFinalizeBlock{Height: 1, Time: time.Now().UTC()})
	require.NoError(t, err)
	_, err = app.Commit()
	require.NoError(t, err)
}

// IBC v8 persists DenomTrace(path=1, base_denom=2); v10 changes that store
// format. Exercise the registered upgrade handler against legacy wire bytes.
func TestSDK053UpgradeMigratesIBCDenom(t *testing.T) {
	app := apptestutil.CreateApp(t)
	ctx := app.NewContextLegacy(false, tmproto.Header{Height: 1})
	path, baseDenom := "transfer/channel-0", "uiris"
	hash := sha256.Sum256([]byte(path + "/" + baseDenom))
	legacy := protowire.AppendString(protowire.AppendTag(nil, 1, protowire.BytesType), path)
	legacy = protowire.AppendString(protowire.AppendTag(legacy, 2, protowire.BytesType), baseDenom)
	key := append(append([]byte{}, transfertypes.DenomTraceKey...), hash[:]...)
	store := ctx.KVStore(app.GetKey(transfertypes.StoreKey))
	store.Set(key, legacy)
	versions, err := app.UpgradeKeeper.GetModuleVersionMap(ctx)
	require.NoError(t, err)
	versions[transfertypes.ModuleName] = 5
	versions["capability"] = 1
	require.NoError(t, app.UpgradeKeeper.SetModuleVersionMap(ctx, versions))
	require.NoError(t, app.UpgradeKeeper.ApplyUpgrade(ctx, upgradetypes.Plan{Name: sdk0538.Upgrade.UpgradeName, Height: 1}))
	denom, found := app.IBCTransferKeeper.GetDenom(ctx, hash[:])
	require.True(t, found)
	require.Equal(t, path+"/"+baseDenom, denom.Path())
	require.True(t, wrapper.NewICS20Keeper(app.IBCTransferKeeper).HasTrace(ctx, denom.IBCDenom()))
	require.False(t, store.Has(key))
	versions, err = app.UpgradeKeeper.GetModuleVersionMap(ctx)
	require.NoError(t, err)
	require.EqualValues(t, 6, versions[transfertypes.ModuleName])
	require.NotContains(t, versions, "capability")
}
