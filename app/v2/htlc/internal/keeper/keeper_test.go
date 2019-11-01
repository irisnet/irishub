package keeper

import (
	"testing"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v2/htlc/internal/types"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_CreateHTLC(t *testing.T) {
	ctx, keeper, ak, accs := createTestInput(t, sdk.NewInt(5000000000), 2)

	senderAddr := accs[0].GetAddress().Bytes()
	receiverAddr := accs[1].GetAddress().Bytes()
	receiverOnOtherChain := "receiverOnOtherChain"
	amount := sdk.NewCoins(sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(10)))
	secret := []byte("___abcdefghijklmnopqrstuvwxyz___")
	timestamp := uint64(1580000000)
	hashLock := sdk.SHA256(append(secret, sdk.Uint64ToBigEndian(timestamp)...))
	timeLock := uint64(50)
	expireHeight := timeLock + uint64(ctx.BlockHeight())
	state := types.OPEN
	initSecret := make([]byte, 0)

	_, err := keeper.GetHTLC(ctx, hashLock)
	require.NotNil(t, err)

	htlc := types.NewHTLC(
		senderAddr,
		receiverAddr,
		receiverOnOtherChain,
		amount,
		initSecret,
		timestamp,
		expireHeight,
		state,
	)

	originSenderAccAmt := ak.GetAccount(ctx, senderAddr).GetCoins()

	htlcAddr := auth.HTLCLockedCoinsAccAddr
	require.Nil(t, ak.GetAccount(ctx, htlcAddr))

	_, err = keeper.CreateHTLC(ctx, htlc, hashLock)
	require.Nil(t, err)

	htlcAcc := ak.GetAccount(ctx, htlcAddr)
	require.NotNil(t, htlcAcc)

	amountCreatedHTLC := ak.GetAccount(ctx, htlcAddr).GetCoins()
	require.True(t, amount.IsEqual(amountCreatedHTLC))

	finalSenderAccAmt := ak.GetAccount(ctx, senderAddr).GetCoins()
	require.True(t, originSenderAccAmt.Sub(amount).IsEqual(finalSenderAccAmt))

	htlc, err = keeper.GetHTLC(ctx, hashLock)
	require.Nil(t, err)

	require.Equal(t, accs[0].GetAddress(), htlc.Sender)
	require.Equal(t, accs[1].GetAddress(), htlc.To)
	require.Equal(t, receiverOnOtherChain, htlc.ReceiverOnOtherChain)
	require.Equal(t, amount, htlc.Amount)
	require.Equal(t, []byte(nil), htlc.Secret)
	require.Equal(t, timestamp, htlc.Timestamp)
	require.Equal(t, expireHeight, htlc.ExpireHeight)
	require.Equal(t, state, htlc.State)

	store := ctx.KVStore(keeper.storeKey)
	require.True(t, store.Has(KeyHTLCExpireQueue(htlc.ExpireHeight, hashLock)))
}

func newHashLock(secret []byte, timestamp uint64) []byte {
	if timestamp > 0 {
		return sdk.SHA256(append(secret, sdk.Uint64ToBigEndian(timestamp)...))
	}
	return sdk.SHA256(secret)
}

func TestKeeper_ClaimHTLC(t *testing.T) {
	ctx, keeper, ak, accs := createTestInput(t, sdk.NewInt(5000000000), 2)

	senderAddr := accs[0].GetAddress().Bytes()
	receiverAddr := accs[1].GetAddress().Bytes()
	receiverOnOtherChain := "receiverOnOtherChain"
	amount := sdk.NewCoins(sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(10)))
	secret1 := []byte("___abcdefghijklmnopqrstuvwxyz___")
	secret2 := []byte("___00000000000000000000000000___")
	timestamp := uint64(1580000000)
	timestampNil := uint64(0)
	timeLock := uint64(50)
	expireHeight := timeLock + uint64(ctx.BlockHeight())
	state := types.OPEN
	initSecret := make([]byte, 0)

	testData := []struct {
		expectPass           bool
		senderAddr           []byte
		toAddr               []byte
		receiverOnOtherChain string
		amount               sdk.Coins
		secret               []byte
		timestamp            uint64
		hashLock             []byte
		timeLock             uint64
		expireHeight         uint64
		state                types.HTLCState
		initSecret           []byte
	}{
		// timestamp > 0
		{true, senderAddr, receiverAddr, receiverOnOtherChain, amount, secret1, timestamp, newHashLock(secret1, timestamp), timeLock, expireHeight, state, initSecret},
		// timestamp = 0
		{true, senderAddr, receiverAddr, receiverOnOtherChain, amount, secret1, timestampNil, newHashLock(secret1, timestampNil), timeLock, expireHeight, state, initSecret},
		// invalid secret
		{false, senderAddr, receiverAddr, receiverOnOtherChain, amount, secret1, timestampNil, newHashLock(secret2, timestampNil), timeLock, expireHeight, state, initSecret},
	}

	for i, td := range testData {
		if td.expectPass {
			htlc := types.NewHTLC(
				td.senderAddr,
				td.toAddr,
				td.receiverOnOtherChain,
				td.amount,
				td.initSecret,
				td.timestamp,
				td.expireHeight,
				td.state,
			)

			_, err := keeper.CreateHTLC(ctx, htlc, td.hashLock)
			require.Nil(t, err, "TestData: %d", i)

			htlc, err = keeper.GetHTLC(ctx, td.hashLock)
			require.Nil(t, err, "TestData: %d", i)
			require.Equal(t, types.OPEN, htlc.State, "TestData: %d", i)

			htlcAddr := auth.HTLCLockedCoinsAccAddr

			originHTLCAmount := ak.GetAccount(ctx, htlcAddr).GetCoins()
			originReceiverAmount := ak.GetAccount(ctx, receiverAddr).GetCoins()

			_, err = keeper.ClaimHTLC(ctx, td.hashLock, td.secret)
			require.Nil(t, err, "TestData: %d", i)

			htlc, _ = keeper.GetHTLC(ctx, td.hashLock)
			require.Equal(t, types.COMPLETED, htlc.State, "TestData: %d", i)

			store := ctx.KVStore(keeper.storeKey)
			require.True(t, !store.Has(KeyHTLCExpireQueue(htlc.ExpireHeight, td.hashLock)))

			claimedHTLCAmount := ak.GetAccount(ctx, htlcAddr).GetCoins()
			claimedReceiverAmount := ak.GetAccount(ctx, receiverAddr).GetCoins()

			require.True(t, originHTLCAmount.Sub(amount).IsEqual(claimedHTLCAmount), "TestData: %d", i)
			require.True(t, originReceiverAmount.Add(amount).IsEqual(claimedReceiverAmount), "TestData: %d", i)

		} else {
			htlc := types.NewHTLC(
				td.senderAddr,
				td.toAddr,
				td.receiverOnOtherChain,
				td.amount,
				td.initSecret,
				td.timestamp,
				td.expireHeight,
				td.state,
			)

			_, err := keeper.CreateHTLC(ctx, htlc, td.hashLock)
			require.Nil(t, err, "TestData: %d", i)

			htlc, err = keeper.GetHTLC(ctx, td.hashLock)
			require.Nil(t, err, "TestData: %d", i)
			require.Equal(t, types.OPEN, htlc.State, "TestData: %d", i)

			htlcAddr := auth.HTLCLockedCoinsAccAddr

			originHTLCAmount := ak.GetAccount(ctx, htlcAddr).GetCoins()
			originReceiverAmount := ak.GetAccount(ctx, receiverAddr).GetCoins()

			_, err = keeper.ClaimHTLC(ctx, td.hashLock, td.secret)
			require.NotNil(t, err, "TestData: %d", i)

			htlc, _ = keeper.GetHTLC(ctx, td.hashLock)
			require.Equal(t, types.OPEN, htlc.State, "TestData: %d", i)

			claimedHTLCAmount := ak.GetAccount(ctx, htlcAddr).GetCoins()
			claimedReceiverAmount := ak.GetAccount(ctx, receiverAddr).GetCoins()

			require.True(t, originHTLCAmount.IsEqual(claimedHTLCAmount), "TestData: %d", i)
			require.True(t, originReceiverAmount.IsEqual(claimedReceiverAmount), "TestData: %d", i)
		}
	}
}

func TestKeeper_RefundHTLC(t *testing.T) {
	ctx, keeper, ak, accs := createTestInput(t, sdk.NewInt(5000000000), 2)

	senderAddr := accs[0].GetAddress().Bytes()
	receiverAddr := accs[1].GetAddress().Bytes()
	receiverOnOtherChain := "receiverOnOtherChain"
	amount := sdk.NewCoins(sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(10)))
	timestamp := uint64(1580000000)
	secret := []byte("___abcdefghijklmnopqrstuvwxyz___")
	hashLock := sdk.SHA256(append(secret, sdk.Uint64ToBigEndian(timestamp)...))
	timeLock := uint64(50)
	expireHeight := timeLock + uint64(ctx.BlockHeight())
	state := types.EXPIRED
	initSecret := make([]byte, 0)

	htlc := types.NewHTLC(
		senderAddr,
		receiverAddr,
		receiverOnOtherChain,
		amount,
		initSecret,
		timestamp,
		expireHeight,
		state,
	)

	_, err := keeper.CreateHTLC(ctx, htlc, hashLock)
	require.Nil(t, err)

	htlc, err = keeper.GetHTLC(ctx, hashLock)
	require.Nil(t, err)
	require.Equal(t, types.EXPIRED, htlc.State)

	htlcAddr := auth.HTLCLockedCoinsAccAddr

	originHTLCAmount := ak.GetAccount(ctx, htlcAddr).GetCoins()
	originSenderAmount := ak.GetAccount(ctx, senderAddr).GetCoins()

	_, err = keeper.RefundHTLC(ctx, hashLock)
	require.Nil(t, err)

	htlc, _ = keeper.GetHTLC(ctx, hashLock)
	require.Equal(t, types.REFUNDED, htlc.State)

	claimedHTLCAmount := ak.GetAccount(ctx, htlcAddr).GetCoins()
	claimedSenderAmount := ak.GetAccount(ctx, senderAddr).GetCoins()

	require.True(t, originHTLCAmount.Sub(amount).IsEqual(claimedHTLCAmount))
	require.True(t, originSenderAmount.Add(amount).IsEqual(claimedSenderAmount))
}
