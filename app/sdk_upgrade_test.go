package app_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"testing"
	"time"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	nfttypes "github.com/bianjieai/nft-transfer/types"
	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icahosttypes "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/host/types"
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	solomachine "github.com/cosmos/ibc-go/v10/modules/light-clients/06-solomachine"
	tendermint "github.com/cosmos/ibc-go/v10/modules/light-clients/07-tendermint"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protowire"

	v500 "github.com/irisnet/irishub/v5/app/upgrades/v500"
	minttypes "github.com/irisnet/irishub/v5/modules/mint/types"
	apptestutil "github.com/irisnet/irishub/v5/testutil"
	"github.com/irisnet/irishub/v5/wrapper"
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
	require.NoError(t, app.UpgradeKeeper.ApplyUpgrade(ctx, upgradetypes.Plan{Name: v500.Upgrade.UpgradeName, Height: 1}))
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

func TestSDK053UpgradeClearsLSMAndKeepsClaims(t *testing.T) {
	app := apptestutil.CreateApp(t)
	ctx := app.NewContextLegacy(false, tmproto.Header{Height: 1})
	store := ctx.KVStore(app.GetKey(stakingtypes.StoreKey))
	valAddr := sdk.ValAddress(bytes.Repeat([]byte{1}, 20))
	source := sdk.AccAddress(address.Module("lsm", []byte("tokenizeshare_1")))
	holder := sdk.AccAddress(bytes.Repeat([]byte{2}, 20))
	rewardOwner := sdk.AccAddress(bytes.Repeat([]byte{3}, 20))

	validator, err := stakingtypes.NewValidator(valAddr.String(), ed25519.GenPrivKey().PubKey(), stakingtypes.Description{})
	require.NoError(t, err)
	validator.Tokens = math.NewInt(160)
	validator.DelegatorShares = math.LegacyNewDec(160)
	validator.MinSelfDelegation = math.NewInt(100)
	require.NoError(t, app.StakingKeeper.SetValidator(ctx, validator))
	delegation := stakingtypes.NewDelegation(source.String(), valAddr.String(), math.LegacyNewDec(160))
	require.NoError(t, app.StakingKeeper.SetDelegation(ctx, delegation))
	starting := distrtypes.NewDelegatorStartingInfo(1, math.LegacyNewDec(160), 1)
	require.NoError(t, app.DistrKeeper.SetDelegatorStartingInfo(ctx, valAddr, source, starting))
	defaultWithdraw, err := app.DistrKeeper.GetDelegatorWithdrawAddr(ctx, source)
	require.NoError(t, err)
	require.Equal(t, source, defaultWithdraw)

	// Legacy records carry the reward recipient separately from the derivative
	// coin holder, so the migration must move that association before deletion.
	record := protowire.AppendVarint(protowire.AppendTag(nil, 1, protowire.VarintType), 1)
	record = protowire.AppendString(protowire.AppendTag(record, 2, protowire.BytesType), rewardOwner.String())
	record = protowire.AppendString(protowire.AppendTag(record, 3, protowire.BytesType), "tokenizeshare_1")
	record = protowire.AppendString(protowire.AppendTag(record, 4, protowire.BytesType), valAddr.String())
	recordKey := append([]byte{0x81}, make([]byte, 8)...)
	binary.BigEndian.PutUint64(recordKey[1:], 1)
	store.Set(recordKey, record)

	derivative := sdk.NewCoins(sdk.NewInt64Coin(valAddr.String()+"/1", 160))
	require.NoError(t, app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, derivative))
	require.NoError(t, app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, holder, derivative))
	bankBalance := app.BankKeeper.GetAllBalances(ctx, holder)
	bankSupply := app.BankKeeper.GetSupply(ctx, derivative[0].Denom)

	paramRaw := store.Get(stakingtypes.ParamsKey)
	for _, field := range []protowire.Number{7, 8, 9} {
		paramRaw = protowire.AppendBytes(protowire.AppendTag(paramRaw, field, protowire.BytesType), []byte("1.000000000000000000"))
	}
	store.Set(stakingtypes.ParamsKey, paramRaw)
	valKey := stakingtypes.GetValidatorKey(valAddr)
	valRaw := store.Get(valKey)
	for _, field := range []protowire.Number{14, 15} {
		valRaw = protowire.AppendBytes(protowire.AppendTag(valRaw, field, protowire.BytesType), []byte("160.000000000000000000"))
	}
	store.Set(valKey, valRaw)
	delKey := stakingtypes.GetDelegationKey(source, valAddr)
	store.Set(delKey, protowire.AppendVarint(protowire.AppendTag(store.Get(delKey), 4, protowire.VarintType), 1))
	for prefix := byte(0x82); prefix <= 0x87; prefix++ {
		store.Set([]byte{prefix, 1}, []byte{1})
	}

	// HistoricalInfo embeds validators and can therefore retain the removed
	// LSM fields even when the current validator records are clean.
	history := stakingtypes.HistoricalInfo{Header: tmproto.Header{Height: 7}}
	historyRaw, err := app.AppCodec().Marshal(&history)
	require.NoError(t, err)
	historicalValidator, err := app.AppCodec().Marshal(&validator)
	require.NoError(t, err)
	for _, field := range []protowire.Number{14, 15} {
		historicalValidator = protowire.AppendBytes(protowire.AppendTag(historicalValidator, field, protowire.BytesType), []byte("160.000000000000000000"))
	}
	historyRaw = protowire.AppendBytes(protowire.AppendTag(historyRaw, 2, protowire.BytesType), historicalValidator)
	historyKey := stakingtypes.GetHistoricalInfoKey(7)
	store.Set(historyKey, historyRaw)
	require.True(t, hasProtoField(store.Get(stakingtypes.ParamsKey), 7))
	require.True(t, hasProtoField(store.Get(valKey), 15))
	require.True(t, hasProtoField(store.Get(delKey), 4))

	plan := upgradetypes.Plan{Name: v500.Upgrade.UpgradeName, Height: 1}
	require.NoError(t, app.UpgradeKeeper.ApplyUpgrade(ctx, plan))
	for _, field := range []protowire.Number{7, 8, 9} {
		require.False(t, hasProtoField(store.Get(stakingtypes.ParamsKey), field))
	}
	for _, field := range []protowire.Number{14, 15} {
		require.False(t, hasProtoField(store.Get(valKey), field))
	}
	require.False(t, hasProtoField(store.Get(delKey), 4))
	for prefix := byte(0x81); prefix <= 0x87; prefix++ {
		iterator := storetypes.KVStorePrefixIterator(store, []byte{prefix})
		require.False(t, iterator.Valid(), "prefix %x should be empty", prefix)
		require.NoError(t, iterator.Close())
	}
	historyRaw = store.Get(historyKey)
	var gotHistory stakingtypes.HistoricalInfo
	require.NoError(t, app.AppCodec().Unmarshal(historyRaw, &gotHistory))
	require.EqualValues(t, 7, gotHistory.Header.Height)
	require.Len(t, gotHistory.Valset, 1)
	require.Equal(t, validator.Tokens, gotHistory.Valset[0].Tokens)
	require.False(t, hasHistoricalValidatorProtoField(historyRaw, 14))
	require.False(t, hasHistoricalValidatorProtoField(historyRaw, 15))
	gotValidator, err := app.StakingKeeper.GetValidator(ctx, valAddr)
	require.NoError(t, err)
	require.True(t, gotValidator.MinSelfDelegation.IsZero())
	require.Equal(t, validator.Tokens, gotValidator.Tokens)
	require.Equal(t, validator.DelegatorShares, gotValidator.DelegatorShares)
	gotDelegation, err := app.StakingKeeper.GetDelegation(ctx, source, valAddr)
	require.NoError(t, err)
	require.Equal(t, delegation, gotDelegation)
	gotStarting, err := app.DistrKeeper.GetDelegatorStartingInfo(ctx, valAddr, source)
	require.NoError(t, err)
	require.Equal(t, starting, gotStarting)
	gotWithdraw, err := app.DistrKeeper.GetDelegatorWithdrawAddr(ctx, source)
	require.NoError(t, err)
	require.Equal(t, rewardOwner, gotWithdraw)
	require.Equal(t, bankBalance, app.BankKeeper.GetAllBalances(ctx, holder))
	require.Equal(t, bankSupply, app.BankKeeper.GetSupply(ctx, derivative[0].Denom))

	// The cleanup also accepts a state that has already been normalized.
	require.NoError(t, app.UpgradeKeeper.ApplyUpgrade(ctx, plan))
	require.Equal(t, bankBalance, app.BankKeeper.GetAllBalances(ctx, holder))
	require.Equal(t, bankSupply, app.BankKeeper.GetSupply(ctx, derivative[0].Denom))
}

func hasProtoField(data []byte, field protowire.Number) bool {
	for len(data) != 0 {
		number, wireType, tagLength := protowire.ConsumeTag(data)
		if tagLength < 0 {
			return false
		}
		data = data[tagLength:]
		valueLength := protowire.ConsumeFieldValue(number, wireType, data)
		if valueLength < 0 {
			return false
		}
		if number == field {
			return true
		}
		data = data[valueLength:]
	}
	return false
}

func hasHistoricalValidatorProtoField(data []byte, field protowire.Number) bool {
	for len(data) != 0 {
		number, wireType, tagLength := protowire.ConsumeTag(data)
		if tagLength < 0 {
			return false
		}
		data = data[tagLength:]
		valueLength := protowire.ConsumeFieldValue(number, wireType, data)
		if valueLength < 0 {
			return false
		}
		if number == 2 && wireType == protowire.BytesType {
			validatorRaw, length := protowire.ConsumeBytes(data)
			return length >= 0 && hasProtoField(validatorRaw, field)
		}
		data = data[valueLength:]
	}
	return false
}
