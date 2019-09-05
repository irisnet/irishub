package bank

import (
	"testing"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/store"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	authKey := sdk.NewKVStoreKey("authkey")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	return ms, authKey
}

func TestKeeper(t *testing.T) {
	ms, authKey := setupMultiStore()

	cdc := codec.New()
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	accountKeeper := auth.NewAccountKeeper(cdc, authKey, auth.ProtoBaseAccount)
	bankKeeper := NewBaseKeeper(cdc, accountKeeper)

	addr := sdk.AccAddress([]byte("addr1"))
	addr2 := sdk.AccAddress([]byte("addr2"))
	addr3 := sdk.AccAddress([]byte("addr3"))
	acc := accountKeeper.NewAccountWithAddress(ctx, addr)

	// Test GetCoins/SetCoins
	accountKeeper.SetAccount(ctx, acc)
	require.True(t, bankKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{}))

	bankKeeper.AddCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foo-min", 10)})
	require.True(t, bankKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("foo-min", 10)}))

	// Test HasCoins
	require.True(t, bankKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foo-min", 10)}))
	require.True(t, bankKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foo-min", 5)}))
	require.False(t, bankKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foo-min", 15)}))
	require.False(t, bankKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("bar-min", 5)}))

	// Test AddCoins
	bankKeeper.AddCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foo-min", 15)})
	require.True(t, bankKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("foo-min", 25)}))

	bankKeeper.AddCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("bar-min", 15)})
	require.True(t, bankKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("bar-min", 15), sdk.NewInt64Coin("foo-min", 25)}))

	// Test SubtractCoins
	bankKeeper.SubtractCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foo-min", 10)})
	bankKeeper.SubtractCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("bar-min", 5)})
	require.True(t, bankKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("bar-min", 10), sdk.NewInt64Coin("foo-min", 15)}))

	bankKeeper.SubtractCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("bar-min", 11)})
	require.True(t, bankKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("bar-min", 10), sdk.NewInt64Coin("foo-min", 15)}))

	bankKeeper.SubtractCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("bar-min", 10)})
	require.True(t, bankKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("foo-min", 15)}))
	require.False(t, bankKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("bar-min", 1)}))

	// Test SendCoins
	bankKeeper.SendCoins(ctx, addr, addr2, sdk.Coins{sdk.NewInt64Coin("foo-min", 5)})
	require.True(t, bankKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("foo-min", 10)}))
	require.True(t, bankKeeper.GetCoins(ctx, addr2).IsEqual(sdk.Coins{sdk.NewInt64Coin("foo-min", 5)}))

	_, err2 := bankKeeper.SendCoins(ctx, addr, addr2, sdk.Coins{sdk.NewInt64Coin("foo-min", 50)})
	assert.Implements(t, (*sdk.Error)(nil), err2)
	require.True(t, bankKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("foo-min", 10)}))
	require.True(t, bankKeeper.GetCoins(ctx, addr2).IsEqual(sdk.Coins{sdk.NewInt64Coin("foo-min", 5)}))

	bankKeeper.AddCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("bar-min", 30)})
	bankKeeper.SendCoins(ctx, addr, addr2, sdk.Coins{sdk.NewInt64Coin("bar-min", 10), sdk.NewInt64Coin("foo-min", 5)})
	require.True(t, bankKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("bar-min", 20), sdk.NewInt64Coin("foo-min", 5)}))
	require.True(t, bankKeeper.GetCoins(ctx, addr2).IsEqual(sdk.Coins{sdk.NewInt64Coin("bar-min", 10), sdk.NewInt64Coin("foo-min", 10)}))

	// Test InputOutputCoins
	input1 := NewInput(addr2, sdk.Coins{sdk.NewInt64Coin("foo-min", 2)})
	output1 := NewOutput(addr, sdk.Coins{sdk.NewInt64Coin("foo-min", 2)})
	bankKeeper.InputOutputCoins(ctx, []Input{input1}, []Output{output1})
	require.True(t, bankKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("bar-min", 20), sdk.NewInt64Coin("foo-min", 7)}))
	require.True(t, bankKeeper.GetCoins(ctx, addr2).IsEqual(sdk.Coins{sdk.NewInt64Coin("bar-min", 10), sdk.NewInt64Coin("foo-min", 8)}))

	inputs := []Input{
		NewInput(addr, sdk.Coins{sdk.NewInt64Coin("foo-min", 3)}),
		NewInput(addr2, sdk.Coins{sdk.NewInt64Coin("bar-min", 3), sdk.NewInt64Coin("foo-min", 2)}),
	}

	outputs := []Output{
		NewOutput(addr, sdk.Coins{sdk.NewInt64Coin("bar-min", 1)}),
		NewOutput(addr3, sdk.Coins{sdk.NewInt64Coin("bar-min", 2), sdk.NewInt64Coin("foo-min", 5)}),
	}
	bankKeeper.InputOutputCoins(ctx, inputs, outputs)
	require.True(t, bankKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("bar-min", 21), sdk.NewInt64Coin("foo-min", 4)}))
	require.True(t, bankKeeper.GetCoins(ctx, addr2).IsEqual(sdk.Coins{sdk.NewInt64Coin("bar-min", 7), sdk.NewInt64Coin("foo-min", 6)}))
	require.True(t, bankKeeper.GetCoins(ctx, addr3).IsEqual(sdk.Coins{sdk.NewInt64Coin("bar-min", 2), sdk.NewInt64Coin("foo-min", 5)}))

}

func TestSendKeeper(t *testing.T) {
	ms, authKey := setupMultiStore()

	cdc := codec.New()
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	accountKeeper := auth.NewAccountKeeper(cdc, authKey, auth.ProtoBaseAccount)
	bankKeeper := NewBaseKeeper(cdc, accountKeeper)
	sendKeeper := NewBaseSendKeeper(accountKeeper)

	addr := sdk.AccAddress([]byte("addr1"))
	addr2 := sdk.AccAddress([]byte("addr2"))
	addr3 := sdk.AccAddress([]byte("addr3"))
	acc := accountKeeper.NewAccountWithAddress(ctx, addr)

	// Test GetCoins/SetCoins
	accountKeeper.SetAccount(ctx, acc)
	require.True(t, sendKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{}))

	bankKeeper.AddCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foo-min", 10)})
	bankKeeper.IncreaseTotalSupply(ctx, sdk.NewInt64Coin("foo-min", 10))
	require.True(t, sendKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("foo-min", 10)}))

	// Test HasCoins
	require.True(t, sendKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foo-min", 10)}))
	require.True(t, sendKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foo-min", 5)}))
	require.False(t, sendKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foo-min", 15)}))
	require.False(t, sendKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("bar-min", 5)}))

	bankKeeper.BurnCoins(ctx, addr, bankKeeper.GetCoins(ctx, addr))
	bankKeeper.AddCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foo-min", 15)})
	bankKeeper.IncreaseTotalSupply(ctx, sdk.NewInt64Coin("foo-min", 15))
	// Test SendCoins
	sendKeeper.SendCoins(ctx, addr, addr2, sdk.Coins{sdk.NewInt64Coin("foo-min", 5)})
	require.True(t, sendKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("foo-min", 10)}))
	require.True(t, sendKeeper.GetCoins(ctx, addr2).IsEqual(sdk.Coins{sdk.NewInt64Coin("foo-min", 5)}))

	_, err2 := sendKeeper.SendCoins(ctx, addr, addr2, sdk.Coins{sdk.NewInt64Coin("foo-min", 50)})
	assert.Implements(t, (*sdk.Error)(nil), err2)
	require.True(t, sendKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("foo-min", 10)}))
	require.True(t, sendKeeper.GetCoins(ctx, addr2).IsEqual(sdk.Coins{sdk.NewInt64Coin("foo-min", 5)}))

	bankKeeper.AddCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("bar-min", 30)})
	bankKeeper.IncreaseTotalSupply(ctx, sdk.NewInt64Coin("bar-min", 30))

	sendKeeper.SendCoins(ctx, addr, addr2, sdk.Coins{sdk.NewInt64Coin("bar-min", 10), sdk.NewInt64Coin("foo-min", 5)})
	require.True(t, sendKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("bar-min", 20), sdk.NewInt64Coin("foo-min", 5)}))
	require.True(t, sendKeeper.GetCoins(ctx, addr2).IsEqual(sdk.Coins{sdk.NewInt64Coin("bar-min", 10), sdk.NewInt64Coin("foo-min", 10)}))

	// Test InputOutputCoins
	input1 := NewInput(addr2, sdk.Coins{sdk.NewInt64Coin("foo-min", 2)})
	output1 := NewOutput(addr, sdk.Coins{sdk.NewInt64Coin("foo-min", 2)})
	sendKeeper.InputOutputCoins(ctx, []Input{input1}, []Output{output1})
	require.True(t, sendKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("bar-min", 20), sdk.NewInt64Coin("foo-min", 7)}))
	require.True(t, sendKeeper.GetCoins(ctx, addr2).IsEqual(sdk.Coins{sdk.NewInt64Coin("bar-min", 10), sdk.NewInt64Coin("foo-min", 8)}))

	inputs := []Input{
		NewInput(addr, sdk.Coins{sdk.NewInt64Coin("foo-min", 3)}),
		NewInput(addr2, sdk.Coins{sdk.NewInt64Coin("bar-min", 3), sdk.NewInt64Coin("foo-min", 2)}),
	}

	outputs := []Output{
		NewOutput(addr, sdk.Coins{sdk.NewInt64Coin("bar-min", 1)}),
		NewOutput(addr3, sdk.Coins{sdk.NewInt64Coin("bar-min", 2), sdk.NewInt64Coin("foo-min", 5)}),
	}
	sendKeeper.InputOutputCoins(ctx, inputs, outputs)
	require.True(t, sendKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("bar-min", 21), sdk.NewInt64Coin("foo-min", 4)}))
	require.True(t, sendKeeper.GetCoins(ctx, addr2).IsEqual(sdk.Coins{sdk.NewInt64Coin("bar-min", 7), sdk.NewInt64Coin("foo-min", 6)}))
	require.True(t, sendKeeper.GetCoins(ctx, addr3).IsEqual(sdk.Coins{sdk.NewInt64Coin("bar-min", 2), sdk.NewInt64Coin("foo-min", 5)}))

}

func TestViewKeeper(t *testing.T) {
	ms, authKey := setupMultiStore()

	cdc := codec.New()
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	accountKeeper := auth.NewAccountKeeper(cdc, authKey, auth.ProtoBaseAccount)
	bankKeeper := NewBaseKeeper(cdc, accountKeeper)
	viewKeeper := NewBaseViewKeeper(accountKeeper)

	addr := sdk.AccAddress([]byte("addr1"))
	acc := accountKeeper.NewAccountWithAddress(ctx, addr)

	// Test GetCoins/SetCoins
	accountKeeper.SetAccount(ctx, acc)
	require.True(t, viewKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{}))

	bankKeeper.AddCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foo-min", 10)})
	bankKeeper.IncreaseTotalSupply(ctx, sdk.NewInt64Coin("foo-min", 10))
	require.True(t, viewKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("foo-min", 10)}))

	// Test HasCoins
	require.True(t, viewKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foo-min", 10)}))
	require.True(t, viewKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foo-min", 5)}))
	require.False(t, viewKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foo-min", 15)}))
	require.False(t, viewKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("bar-min", 5)}))
}
