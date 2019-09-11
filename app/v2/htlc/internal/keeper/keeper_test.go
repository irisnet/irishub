package keeper

import (
	"testing"

	"github.com/irisnet/irishub/app/v2/htlc/internal/types"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_CreateHTLC(t *testing.T) {
	ctx, keeper, ak, accs := createTestInput(t, sdk.NewInt(5000000000), 2)

	senderAddr := accs[0].GetAddress().Bytes()
	receiverAddr := accs[1].GetAddress().Bytes()
	receiverOnOtherChain := []byte("receiverOnOtherChain")
	amount := sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(10))
	secret := []byte("___abcdefghijklmnopqrstuvwxyz___")
	timestamp := uint64(1580000000)
	hashLock := sdk.SHA256(append(secret, sdk.Uint64ToBigEndian(timestamp)...))
	timeLock := uint64(50)
	expireHeight := timeLock + uint64(ctx.BlockHeight())
	state := types.StateOpen
	initSecret := make([]byte, 32)

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

	originSenderAccAmt := ak.GetAccount(ctx, senderAddr).GetCoins().AmountOf(amount.Denom)

	htlcAddr := getHTLCAddress(amount.Denom)
	require.Nil(t, ak.GetAccount(ctx, htlcAddr))

	_, err = keeper.CreateHTLC(ctx, htlc, hashLock)
	require.Nil(t, err)

	htlcAcc := ak.GetAccount(ctx, htlcAddr)
	require.NotNil(t, htlcAcc)

	amountCreatedHTLC := ak.GetAccount(ctx, htlcAddr).GetCoins().AmountOf(amount.Denom)
	require.Equal(t, amount.Amount.Int64(), amountCreatedHTLC.Int64())

	finalSenderAccAmt := ak.GetAccount(ctx, senderAddr).GetCoins().AmountOf(amount.Denom)
	require.Equal(t, originSenderAccAmt.Sub(amount.Amount).Int64(), finalSenderAccAmt.Int64())

	htlc, err = keeper.GetHTLC(ctx, hashLock)
	require.Nil(t, err)

	require.Equal(t, accs[0].GetAddress(), htlc.Sender)
	require.Equal(t, accs[1].GetAddress(), htlc.Receiver)
	require.Equal(t, receiverOnOtherChain, htlc.ReceiverOnOtherChain)
	require.Equal(t, amount, htlc.Amount)
	require.Equal(t, initSecret, htlc.Secret)
	require.Equal(t, timestamp, htlc.Timestamp)
	require.Equal(t, expireHeight, htlc.ExpireHeight)
	require.Equal(t, state, htlc.State)
}

func TestKeeper_ClaimHTLC(t *testing.T) {
	ctx, keeper, ak, accs := createTestInput(t, sdk.NewInt(5000000000), 2)

	senderAddr := accs[0].GetAddress().Bytes()
	receiverAddr := accs[1].GetAddress().Bytes()
	receiverOnOtherChain := []byte("receiverOnOtherChain")
	amount := sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(10))
	secret := []byte("___abcdefghijklmnopqrstuvwxyz___")
	timestamp := uint64(1580000000)
	hashLock := sdk.SHA256(append(secret, sdk.Uint64ToBigEndian(timestamp)...))
	timeLock := uint64(50)
	expireHeight := timeLock + uint64(ctx.BlockHeight())
	state := types.StateOpen
	initSecret := make([]byte, 32)

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
	require.Equal(t, types.StateOpen, htlc.State)

	htlcAddr := getHTLCAddress(amount.Denom)

	originHTLCAmount := ak.GetAccount(ctx, htlcAddr).GetCoins().AmountOf(amount.Denom)
	originReceiverAmount := ak.GetAccount(ctx, receiverAddr).GetCoins().AmountOf(amount.Denom)

	_, err = keeper.ClaimHTLC(ctx, secret, hashLock)
	require.Nil(t, err)

	htlc, _ = keeper.GetHTLC(ctx, hashLock)
	require.Equal(t, types.StateCompleted, htlc.State)

	claimedHTLCAmount := ak.GetAccount(ctx, htlcAddr).GetCoins().AmountOf(amount.Denom)
	claimedReceiverAmount := ak.GetAccount(ctx, receiverAddr).GetCoins().AmountOf(amount.Denom)

	require.Equal(t, originHTLCAmount.Sub(amount.Amount).Int64(), claimedHTLCAmount.Int64())
	require.Equal(t, originReceiverAmount.Add(amount.Amount).Int64(), claimedReceiverAmount.Int64())
}

func TestKeeper_RefundHTLC(t *testing.T) {
	ctx, keeper, ak, accs := createTestInput(t, sdk.NewInt(5000000000), 2)

	senderAddr := accs[0].GetAddress().Bytes()
	receiverAddr := accs[1].GetAddress().Bytes()
	receiverOnOtherChain := []byte("receiverOnOtherChain")
	amount := sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(10))
	timestamp := uint64(1580000000)
	secret := []byte("___abcdefghijklmnopqrstuvwxyz___")
	hashLock := sdk.SHA256(append(secret, sdk.Uint64ToBigEndian(timestamp)...))
	timeLock := uint64(50)
	expireHeight := timeLock + uint64(ctx.BlockHeight())
	state := types.StateExpired
	initSecret := make([]byte, 32)

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
	require.Equal(t, types.StateExpired, htlc.State)

	htlcAddr := getHTLCAddress(amount.Denom)

	originHTLCAmount := ak.GetAccount(ctx, htlcAddr).GetCoins().AmountOf(amount.Denom)
	originSenderAmount := ak.GetAccount(ctx, senderAddr).GetCoins().AmountOf(amount.Denom)

	_, err = keeper.RefundHTLC(ctx, hashLock)
	require.Nil(t, err)

	htlc, _ = keeper.GetHTLC(ctx, hashLock)
	require.Equal(t, types.StateRefunded, htlc.State)

	claimedHTLCAmount := ak.GetAccount(ctx, htlcAddr).GetCoins().AmountOf(amount.Denom)
	claimedSenderAmount := ak.GetAccount(ctx, senderAddr).GetCoins().AmountOf(amount.Denom)

	require.Equal(t, originHTLCAmount.Sub(amount.Amount).Int64(), claimedHTLCAmount.Int64())
	require.Equal(t, originSenderAmount.Add(amount.Amount).Int64(), claimedSenderAmount.Int64())
}
