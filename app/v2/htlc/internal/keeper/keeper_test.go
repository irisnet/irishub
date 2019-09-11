package keeper

import (
	"encoding/hex"
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
	outAmount := sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(10))
	inAmount := uint64(100)
	secretStr := "___abcdefghijklmnopqrstuvwxyz___"
	timestamp := uint64(1580000000)
	secretHashLock := sdk.SHA256(append([]byte(secretStr), sdk.Uint64ToBigEndian(timestamp)...))
	timeLock := uint64(50)
	expireHeight := timeLock + uint64(ctx.BlockHeight())
	state := types.StateOpen
	initSecret := make([]byte, 32)

	_, err := keeper.GetHTLC(ctx, secretHashLock)
	require.NotNil(t, err)

	htlc := types.NewHTLC(
		senderAddr,
		receiverAddr,
		receiverOnOtherChain,
		outAmount,
		inAmount,
		initSecret,
		timestamp,
		expireHeight,
		state,
	)

	originSenderAccAmt := ak.GetAccount(ctx, senderAddr).GetCoins().AmountOf(outAmount.Denom)

	htlcAddr := getHTLCAddress(outAmount.Denom)
	require.Nil(t, ak.GetAccount(ctx, htlcAddr))

	_, err = keeper.CreateHTLC(ctx, htlc, secretHashLock)
	require.Nil(t, err)

	htlcAcc := ak.GetAccount(ctx, htlcAddr)
	require.NotNil(t, htlcAcc)

	amountCreatedHTLC := ak.GetAccount(ctx, htlcAddr).GetCoins().AmountOf(outAmount.Denom)
	require.Equal(t, outAmount.Amount.Int64(), amountCreatedHTLC.Int64())

	finalSenderAccAmt := ak.GetAccount(ctx, senderAddr).GetCoins().AmountOf(outAmount.Denom)
	require.Equal(t, originSenderAccAmt.Sub(outAmount.Amount).Int64(), finalSenderAccAmt.Int64())

	htlc, err = keeper.GetHTLC(ctx, secretHashLock)
	require.Nil(t, err)

	require.Equal(t, accs[0].GetAddress(), htlc.Sender)
	require.Equal(t, accs[1].GetAddress(), htlc.Receiver)
	require.Equal(t, receiverOnOtherChain, htlc.ReceiverOnOtherChain)
	require.Equal(t, outAmount, htlc.OutAmount)
	require.Equal(t, inAmount, htlc.InAmount)
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
	outAmount := sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(10))
	inAmount := uint64(100)
	secretStr := "___abcdefghijklmnopqrstuvwxyz___"
	secretHexStr := "5f5f5f6162636465666768696a6b6c6d6e6f707172737475767778797a5f5f5f"
	secret, _ := hex.DecodeString(secretHexStr)
	timestamp := uint64(1580000000)
	secretHashLock := sdk.SHA256(append([]byte(secretStr), sdk.Uint64ToBigEndian(timestamp)...))
	timeLock := uint64(50)
	expireHeight := timeLock + uint64(ctx.BlockHeight())
	state := types.StateOpen
	initSecret := make([]byte, 32)

	htlc := types.NewHTLC(
		senderAddr,
		receiverAddr,
		receiverOnOtherChain,
		outAmount,
		inAmount,
		initSecret,
		timestamp,
		expireHeight,
		state,
	)

	_, err := keeper.CreateHTLC(ctx, htlc, secretHashLock)
	require.Nil(t, err)

	htlc, err = keeper.GetHTLC(ctx, secretHashLock)
	require.Nil(t, err)
	require.Equal(t, types.StateOpen, htlc.State)

	htlcAddr := getHTLCAddress(outAmount.Denom)

	originHTLCAmount := ak.GetAccount(ctx, htlcAddr).GetCoins().AmountOf(outAmount.Denom)
	originReceiverAmount := ak.GetAccount(ctx, receiverAddr).GetCoins().AmountOf(outAmount.Denom)

	_, err = keeper.ClaimHTLC(ctx, secret, secretHashLock)
	require.Nil(t, err)

	htlc, _ = keeper.GetHTLC(ctx, secretHashLock)
	require.Equal(t, types.StateCompleted, htlc.State)

	claimedHTLCAmount := ak.GetAccount(ctx, htlcAddr).GetCoins().AmountOf(outAmount.Denom)
	claimedReceiverAmount := ak.GetAccount(ctx, receiverAddr).GetCoins().AmountOf(outAmount.Denom)

	require.Equal(t, originHTLCAmount.Sub(outAmount.Amount).Int64(), claimedHTLCAmount.Int64())
	require.Equal(t, originReceiverAmount.Add(outAmount.Amount).Int64(), claimedReceiverAmount.Int64())

}

func TestKeeper_RefundHTLC(t *testing.T) {
	ctx, keeper, ak, accs := createTestInput(t, sdk.NewInt(5000000000), 2)

	senderAddr := accs[0].GetAddress().Bytes()
	receiverAddr := accs[1].GetAddress().Bytes()
	receiverOnOtherChain := []byte("receiverOnOtherChain")
	outAmount := sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(10))
	inAmount := uint64(100)
	timestamp := uint64(1580000000)
	secretStr := "___abcdefghijklmnopqrstuvwxyz___"
	secretHashLock := sdk.SHA256(append([]byte(secretStr), sdk.Uint64ToBigEndian(timestamp)...))
	timeLock := uint64(50)
	expireHeight := timeLock + uint64(ctx.BlockHeight())
	state := types.StateExpired
	initSecret := make([]byte, 32)

	htlc := types.NewHTLC(
		senderAddr,
		receiverAddr,
		receiverOnOtherChain,
		outAmount,
		inAmount,
		initSecret,
		timestamp,
		expireHeight,
		state,
	)

	_, err := keeper.CreateHTLC(ctx, htlc, secretHashLock)
	require.Nil(t, err)

	htlc, err = keeper.GetHTLC(ctx, secretHashLock)
	require.Nil(t, err)
	require.Equal(t, types.StateExpired, htlc.State)

	htlcAddr := getHTLCAddress(outAmount.Denom)

	originHTLCAmount := ak.GetAccount(ctx, htlcAddr).GetCoins().AmountOf(outAmount.Denom)
	originSenderAmount := ak.GetAccount(ctx, senderAddr).GetCoins().AmountOf(outAmount.Denom)

	_, err = keeper.RefundHTLC(ctx, secretHashLock)
	require.Nil(t, err)

	htlc, _ = keeper.GetHTLC(ctx, secretHashLock)
	require.Equal(t, types.StateRefunded, htlc.State)

	claimedHTLCAmount := ak.GetAccount(ctx, htlcAddr).GetCoins().AmountOf(outAmount.Denom)
	claimedSenderAmount := ak.GetAccount(ctx, senderAddr).GetCoins().AmountOf(outAmount.Denom)

	require.Equal(t, originHTLCAmount.Sub(outAmount.Amount).Int64(), claimedHTLCAmount.Int64())
	require.Equal(t, originSenderAmount.Add(outAmount.Amount).Int64(), claimedSenderAmount.Int64())
}
