package types_test

import (
	"math/rand"
	time "time"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmtime "github.com/tendermint/tendermint/types/time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/htlc/types"
)

const (
	MinTimeLock uint64 = 220
	MaxTimeLock uint64 = 270
)

func c(denom string, amount int64) sdk.Coin {
	return sdk.NewInt64Coin(denom, amount)
}

func cs(coins ...sdk.Coin) sdk.Coins {
	return sdk.NewCoins(coins...)
}

func htlcs(count int) []types.HTLC {
	var htlcs []types.HTLC
	for i := 0; i < count; i++ {
		htlc := htlc(i)
		htlcs = append(htlcs, htlc)
	}
	return htlcs
}

func htlc(index int) types.HTLC {
	expireOffset := uint64((index * 15) + 360)
	timestamp := uint64(tmtime.Now().Add(time.Duration(index) * time.Minute).Unix())
	randomSecret, _ := GenerateRandomSecret()
	randomHashLock := types.GetHashLock(randomSecret, timestamp)
	amount := cs(c("htltbnb", 50000))
	id := types.GetID(sender, recipient, amount, randomHashLock)

	htlc := types.NewHTLC(
		id,
		sender,
		recipient,
		receiverOnOtherChain,
		senderOnOtherChain,
		amount,
		randomHashLock,
		[]byte{},
		timestamp,
		expireOffset,
		types.Open,
		1,
		true,
		types.Incoming,
	)

	return htlc
}

func GenerateRandomSecret() ([]byte, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return []byte{}, err
	}
	return bytes, nil
}

func GeneratePrivKeyAddressPairs(n int) (keys []crypto.PrivKey, addrs []sdk.AccAddress) {
	r := rand.New(rand.NewSource(12345)) // make the generation deterministic
	keys = make([]crypto.PrivKey, n)
	addrs = make([]sdk.AccAddress, n)
	for i := 0; i < n; i++ {
		secret := make([]byte, 32)
		_, err := r.Read(secret)
		if err != nil {
			panic("Could not read randomness")
		}
		keys[i] = secp256k1.GenPrivKeySecp256k1(secret)
		addrs[i] = sdk.AccAddress(keys[i].PubKey().Address())
	}
	return
}
