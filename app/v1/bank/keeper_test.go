package bank

import (
	"testing"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/codec"
	auth0 "github.com/irisnet/irishub/modules/auth"
	bank0 "github.com/irisnet/irishub/modules/bank"
	"github.com/irisnet/irishub/store"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	authKey := sdk.NewKVStoreKey("authkey")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	return ms, authKey
}

func createV0Account() (auth0.Account,auth0.AccountKeeper,sdk.AccAddress,sdk.Context,*sdk.KVStoreKey){
	cdc0 := codec.New()
	auth0.RegisterBaseAccount(cdc0)

	ms, authKey := setupMultiStore()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	accountKeeper0 := auth0.NewAccountKeeper(cdc0, authKey, auth0.ProtoBaseAccount)
	addr := sdk.AccAddress([]byte("addr1"))
	acc := accountKeeper0.NewAccountWithAddress(ctx, addr)
	accountKeeper0.SetAccount(ctx, acc)

	return acc,accountKeeper0,addr,ctx,authKey
}

func createV1Account(ctx sdk.Context,authKey *sdk.KVStoreKey) (auth.Account,auth.AccountKeeper,sdk.AccAddress){
	cdc1 := codec.New()
	auth.RegisterBaseAccount(cdc1)

	//ms, authKey := setupMultiStore()

	//ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	accountKeeper1 := auth.NewAccountKeeper(cdc1, authKey, auth.ProtoBaseAccount)
	addr := sdk.AccAddress([]byte("addr1"))
	acc := accountKeeper1.NewAccountWithAddress(ctx, addr)
	accountKeeper1.SetAccount(ctx, acc)

	return acc,accountKeeper1,addr
}

func TestAccount(t *testing.T){
	_,accountKeeper0,addr0,ctx,authkey:=createV0Account()
	_,accountKeeper1,addr1:=createV1Account(ctx,authkey)

	bankKeeper0 := bank0.NewBaseKeeper(accountKeeper0)
	//viewKeeper0 := bank0.BaseViewKeeper{accountKeeper0}

	bankKeeper1 :=  NewBaseKeeper(accountKeeper1)
	//viewKeeper1 := NewBaseViewKeeper(accountKeeper1)

	//v0 account add, v1 account get
	//test result: when setting the same ctx and authkey, the account will be same
	bankKeeper0.AddCoins(ctx, addr0, sdk.Coins{sdk.NewInt64Coin("foocoin", 10)})
	//require.True(t, bankKeeper0.GetCoins(ctx0, addr0).IsEqual(sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}))
	//bankKeeper1.AddCoins(ctx1, addr1, sdk.Coins{sdk.NewInt64Coin("foocoin", 10)})
	require.True(t, bankKeeper1.GetCoins(ctx, addr1).IsEqual(sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}))

	//freeze the coins
	bankKeeper0.AddCoins(ctx, addr0, sdk.Coins{sdk.NewInt64Coin("foocoin", 10)})
	bankKeeper1.FreezeCoinFromAddr(ctx,addr1,sdk.NewInt64Coin("foocoin", 10))
	require.True(t, bankKeeper0.GetCoins(ctx, addr0).IsEqual(sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}))
	bankKeeper1.UnfreezeCoinFromAddr(ctx,addr1,sdk.NewInt64Coin("foocoin", 10))
	require.True(t, bankKeeper0.GetCoins(ctx, addr0).IsEqual(sdk.Coins{sdk.NewInt64Coin("foocoin", 20)}))




}

func TestKeeper(t *testing.T) {
	ms, authKey := setupMultiStore()

	cdc := codec.New()
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	accountKeeper := auth.NewAccountKeeper(cdc, authKey, auth.ProtoBaseAccount)
	bankKeeper := NewBaseKeeper(accountKeeper)

	addr := sdk.AccAddress([]byte("addr1"))
	addr2 := sdk.AccAddress([]byte("addr2"))
	addr3 := sdk.AccAddress([]byte("addr3"))
	acc := accountKeeper.NewAccountWithAddress(ctx, addr)

	// Test GetCoins/SetCoins
	accountKeeper.SetAccount(ctx, acc)
	require.True(t, bankKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{}))

	bankKeeper.AddCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 10)})
	require.True(t, bankKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}))

	// Test HasCoins
	require.True(t, bankKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}))
	require.True(t, bankKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 5)}))
	require.False(t, bankKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 15)}))
	require.False(t, bankKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("barcoin", 5)}))

	// Test AddCoins
	bankKeeper.AddCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 15)})
	require.True(t, bankKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("foocoin", 25)}))

	bankKeeper.AddCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("barcoin", 15)})
	require.True(t, bankKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 15), sdk.NewInt64Coin("foocoin", 25)}))

	// Test SubtractCoins
	bankKeeper.SubtractCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 10)})
	bankKeeper.SubtractCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("barcoin", 5)})
	require.True(t, bankKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 10), sdk.NewInt64Coin("foocoin", 15)}))

	bankKeeper.SubtractCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("barcoin", 11)})
	require.True(t, bankKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 10), sdk.NewInt64Coin("foocoin", 15)}))

	bankKeeper.SubtractCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("barcoin", 10)})
	require.True(t, bankKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("foocoin", 15)}))
	require.False(t, bankKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("barcoin", 1)}))

	// Test SendCoins
	bankKeeper.SendCoins(ctx, addr, addr2, sdk.Coins{sdk.NewInt64Coin("foocoin", 5)})
	require.True(t, bankKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}))
	require.True(t, bankKeeper.GetCoins(ctx, addr2).IsEqual(sdk.Coins{sdk.NewInt64Coin("foocoin", 5)}))

	_, err2 := bankKeeper.SendCoins(ctx, addr, addr2, sdk.Coins{sdk.NewInt64Coin("foocoin", 50)})
	assert.Implements(t, (*sdk.Error)(nil), err2)
	require.True(t, bankKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}))
	require.True(t, bankKeeper.GetCoins(ctx, addr2).IsEqual(sdk.Coins{sdk.NewInt64Coin("foocoin", 5)}))

	bankKeeper.AddCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("barcoin", 30)})
	bankKeeper.SendCoins(ctx, addr, addr2, sdk.Coins{sdk.NewInt64Coin("barcoin", 10), sdk.NewInt64Coin("foocoin", 5)})
	require.True(t, bankKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 20), sdk.NewInt64Coin("foocoin", 5)}))
	require.True(t, bankKeeper.GetCoins(ctx, addr2).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 10), sdk.NewInt64Coin("foocoin", 10)}))

	// Test InputOutputCoins
	input1 := NewInput(addr2, sdk.Coins{sdk.NewInt64Coin("foocoin", 2)})
	output1 := NewOutput(addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 2)})
	bankKeeper.InputOutputCoins(ctx, []Input{input1}, []Output{output1})
	require.True(t, bankKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 20), sdk.NewInt64Coin("foocoin", 7)}))
	require.True(t, bankKeeper.GetCoins(ctx, addr2).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 10), sdk.NewInt64Coin("foocoin", 8)}))

	inputs := []Input{
		NewInput(addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 3)}),
		NewInput(addr2, sdk.Coins{sdk.NewInt64Coin("barcoin", 3), sdk.NewInt64Coin("foocoin", 2)}),
	}

	outputs := []Output{
		NewOutput(addr, sdk.Coins{sdk.NewInt64Coin("barcoin", 1)}),
		NewOutput(addr3, sdk.Coins{sdk.NewInt64Coin("barcoin", 2), sdk.NewInt64Coin("foocoin", 5)}),
	}
	bankKeeper.InputOutputCoins(ctx, inputs, outputs)
	require.True(t, bankKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 21), sdk.NewInt64Coin("foocoin", 4)}))
	require.True(t, bankKeeper.GetCoins(ctx, addr2).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 7), sdk.NewInt64Coin("foocoin", 6)}))
	require.True(t, bankKeeper.GetCoins(ctx, addr3).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 2), sdk.NewInt64Coin("foocoin", 5)}))

}

func TestSendKeeper(t *testing.T) {
	ms, authKey := setupMultiStore()

	cdc := codec.New()
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	accountKeeper := auth.NewAccountKeeper(cdc, authKey, auth.ProtoBaseAccount)
	bankKeeper := NewBaseKeeper(accountKeeper)
	sendKeeper := NewBaseSendKeeper(accountKeeper)

	addr := sdk.AccAddress([]byte("addr1"))
	addr2 := sdk.AccAddress([]byte("addr2"))
	addr3 := sdk.AccAddress([]byte("addr3"))
	acc := accountKeeper.NewAccountWithAddress(ctx, addr)

	// Test GetCoins/SetCoins
	accountKeeper.SetAccount(ctx, acc)
	require.True(t, sendKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{}))

	bankKeeper.AddCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 10)})
	bankKeeper.IncreaseLoosenToken(ctx, sdk.Coins{sdk.NewInt64Coin("foocoin", 10)})
	require.True(t, sendKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}))

	// Test HasCoins
	require.True(t, sendKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}))
	require.True(t, sendKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 5)}))
	require.False(t, sendKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 15)}))
	require.False(t, sendKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("barcoin", 5)}))

	bankKeeper.BurnCoinsFromAddr(ctx, addr, bankKeeper.GetCoins(ctx, addr))


	bankKeeper.AddCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 15)})


	bankKeeper.IncreaseLoosenToken(ctx, sdk.Coins{sdk.NewInt64Coin("foocoin", 15)})
	// Test SendCoins
	sendKeeper.SendCoins(ctx, addr, addr2, sdk.Coins{sdk.NewInt64Coin("foocoin", 5)})
	require.True(t, sendKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}))
	require.True(t, sendKeeper.GetCoins(ctx, addr2).IsEqual(sdk.Coins{sdk.NewInt64Coin("foocoin", 5)}))

	_, err2 := sendKeeper.SendCoins(ctx, addr, addr2, sdk.Coins{sdk.NewInt64Coin("foocoin", 50)})
	assert.Implements(t, (*sdk.Error)(nil), err2)
	require.True(t, sendKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}))
	require.True(t, sendKeeper.GetCoins(ctx, addr2).IsEqual(sdk.Coins{sdk.NewInt64Coin("foocoin", 5)}))


	bankKeeper.AddCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("barcoin", 30)})
	bankKeeper.IncreaseLoosenToken(ctx, sdk.Coins{sdk.NewInt64Coin("barcoin", 30)})

	//test get total frozen token when only addr freeze token
	bankKeeper.FreezeCoinFromAddr(ctx, addr, sdk.NewInt64Coin("barcoin", 2))
	require.True(t,bankKeeper.GetFrozenCoin(ctx,accountKeeper,addr,"barcoin").IsEqual(sdk.NewInt64Coin("barcoin", 2)))
	bankKeeper.FreezeCoinFromAddr(ctx, addr, sdk.NewInt64Coin("barcoin", 2))
	require.True(t,bankKeeper.GetFrozenCoin(ctx,accountKeeper,addr,"barcoin").IsEqual(sdk.NewInt64Coin("barcoin", 4)))
	bankKeeper.UnfreezeCoinFromAddr(ctx, addr, sdk.NewInt64Coin("barcoin", 2))
	require.True(t,bankKeeper.GetFrozenCoin(ctx,accountKeeper,addr,"barcoin").IsEqual(sdk.NewInt64Coin("barcoin", 2)))
	bankKeeper.UnfreezeCoinFromAddr(ctx, addr, sdk.NewInt64Coin("barcoin", 2))

	bankKeeper.FreezeCoinFromAddr(ctx, addr, sdk.NewInt64Coin("barcoin", 2))
	//require.True(t,bankKeeper.GetFrozenCoin(ctx,accountKeeper,addr,"barcoin").IsEqual(sdk.NewInt64Coin("barcoin", 2)))
	sendKeeper.SendCoins(ctx, addr, addr2, sdk.Coins{sdk.NewInt64Coin("barcoin", 10), sdk.NewInt64Coin("foocoin", 5)})
	bankKeeper.UnfreezeCoinFromAddr(ctx,addr, sdk.NewInt64Coin("barcoin", 2))

	//test get total frozen token when addr2 freeze token too
	bankKeeper.FreezeCoinFromAddr(ctx, addr2, sdk.NewInt64Coin("barcoin", 2))
	require.True(t,bankKeeper.GetFrozenCoin(ctx,accountKeeper,addr,"barcoin").IsEqual(sdk.NewInt64Coin("barcoin", 2)))
	bankKeeper.FreezeCoinFromAddr(ctx, addr2, sdk.NewInt64Coin("barcoin", 2))
	require.True(t,bankKeeper.GetFrozenCoin(ctx,accountKeeper,addr,"barcoin").IsEqual(sdk.NewInt64Coin("barcoin", 4)))
	bankKeeper.UnfreezeCoinFromAddr(ctx, addr2, sdk.NewInt64Coin("barcoin", 2))
	require.True(t,bankKeeper.GetFrozenCoin(ctx,accountKeeper,addr,"barcoin").IsEqual(sdk.NewInt64Coin("barcoin", 2)))
	bankKeeper.UnfreezeCoinFromAddr(ctx, addr2, sdk.NewInt64Coin("barcoin", 2))

	require.True(t, sendKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 20), sdk.NewInt64Coin("foocoin", 5)}))
	require.True(t, sendKeeper.GetCoins(ctx, addr2).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 10), sdk.NewInt64Coin("foocoin", 10)}))
	// Test InputOutputCoins
	input1 := NewInput(addr2, sdk.Coins{sdk.NewInt64Coin("foocoin", 2)})
	output1 := NewOutput(addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 2)})
	sendKeeper.InputOutputCoins(ctx, []Input{input1}, []Output{output1})

	require.True(t, sendKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 20), sdk.NewInt64Coin("foocoin", 7)}))
	require.True(t, sendKeeper.GetCoins(ctx, addr2).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 10), sdk.NewInt64Coin("foocoin", 8)}))

	inputs := []Input{
		NewInput(addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 3)}),
		NewInput(addr2, sdk.Coins{sdk.NewInt64Coin("barcoin", 3), sdk.NewInt64Coin("foocoin", 2)}),
	}

	outputs := []Output{
		NewOutput(addr, sdk.Coins{sdk.NewInt64Coin("barcoin", 1)}),
		NewOutput(addr3, sdk.Coins{sdk.NewInt64Coin("barcoin", 2), sdk.NewInt64Coin("foocoin", 5)}),
	}
	sendKeeper.InputOutputCoins(ctx, inputs, outputs)

	bankKeeper.FreezeCoinFromAddr(ctx, addr, sdk.NewInt64Coin("barcoin", 2))
	bankKeeper.UnfreezeCoinFromAddr(ctx, addr, sdk.NewInt64Coin("barcoin", 2))
	require.True(t, sendKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 21), sdk.NewInt64Coin("foocoin", 4)}))
	require.True(t, sendKeeper.GetCoins(ctx, addr2).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 7), sdk.NewInt64Coin("foocoin", 6)}))
	require.True(t, sendKeeper.GetCoins(ctx, addr3).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 2), sdk.NewInt64Coin("foocoin", 5)}))


}

func TestViewKeeper(t *testing.T) {
	ms, authKey := setupMultiStore()

	cdc := codec.New()
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	accountKeeper := auth.NewAccountKeeper(cdc, authKey, auth.ProtoBaseAccount)
	bankKeeper := NewBaseKeeper(accountKeeper)
	viewKeeper := NewBaseViewKeeper(accountKeeper)

	addr := sdk.AccAddress([]byte("addr1"))
	acc := accountKeeper.NewAccountWithAddress(ctx, addr)

	// Test GetCoins/SetCoins
	accountKeeper.SetAccount(ctx, acc)
	require.True(t, viewKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{}))

	bankKeeper.AddCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 10)})
	bankKeeper.IncreaseLoosenToken(ctx, sdk.Coins{sdk.NewInt64Coin("foocoin", 10)})
	require.True(t, viewKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}))

	// Test HasCoins
	require.True(t, viewKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}))
	require.True(t, viewKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 5)}))
	require.False(t, viewKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 15)}))
	require.False(t, viewKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("barcoin", 5)}))
}
