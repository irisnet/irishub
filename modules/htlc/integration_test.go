package htlc_test

import (
	"math/rand"
	"time"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmtime "github.com/tendermint/tendermint/types/time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/htlc/types"
)

var (
	DenomMap                    = map[int]string{0: "htltbnb", 1: "htltinc"}
	MinTimeLock          uint64 = 220
	MaxTimeLock          uint64 = 270
	TestDeputy                  = sdk.AccAddress(crypto.AddressHash([]byte("TestDeputy")))
	ReceiverOnOtherChain        = "ReceiverOnOtherChain"
	SenderOnOtherChain          = "SenderOnOtherChain"
)

func i(in int64) sdk.Int {
	return sdk.NewInt(in)
}

func c(denom string, amount int64) sdk.Coin {
	return sdk.NewInt64Coin(denom, amount)
}

func cs(coins ...sdk.Coin) sdk.Coins {
	return sdk.NewCoins(coins...)
}

func ts(minOffset int) uint64 {
	return uint64(tmtime.Now().Add(time.Duration(minOffset) * time.Minute).Unix())
}

func NewHTLTGenesis(deputyAddress sdk.AccAddress) *types.GenesisState {
	return &types.GenesisState{
		Params: types.Params{
			AssetParams: []types.AssetParam{
				{
					Denom: "htltbnb",
					SupplyLimit: types.SupplyLimit{
						Limit:          sdk.NewInt(350000000000000),
						TimeLimited:    false,
						TimeBasedLimit: sdk.ZeroInt(),
						TimePeriod:     time.Hour,
					},
					Active:        true,
					DeputyAddress: TestDeputy.String(),
					FixedFee:      sdk.NewInt(1000),
					MinSwapAmount: sdk.OneInt(),
					MaxSwapAmount: sdk.NewInt(1000000000000),
					MinBlockLock:  MinTimeLock,
					MaxBlockLock:  MaxTimeLock,
				},
				{
					Denom: "htltinc",
					SupplyLimit: types.SupplyLimit{
						Limit:          sdk.NewInt(100000000000),
						TimeLimited:    false,
						TimeBasedLimit: sdk.ZeroInt(),
						TimePeriod:     time.Hour,
					},
					Active:        true,
					DeputyAddress: TestDeputy.String(),
					FixedFee:      sdk.NewInt(1000),
					MinSwapAmount: sdk.OneInt(),
					MaxSwapAmount: sdk.NewInt(1000000000000),
					MinBlockLock:  MinTimeLock,
					MaxBlockLock:  MaxTimeLock,
				},
			},
		},
		Supplies: []types.AssetSupply{
			types.NewAssetSupply(
				sdk.NewCoin("htltbnb", sdk.ZeroInt()),
				sdk.NewCoin("htltbnb", sdk.ZeroInt()),
				sdk.NewCoin("htltbnb", sdk.ZeroInt()),
				sdk.NewCoin("htltbnb", sdk.ZeroInt()),
				time.Duration(0),
			),
			types.NewAssetSupply(
				sdk.NewCoin("htltinc", sdk.ZeroInt()),
				sdk.NewCoin("htltinc", sdk.ZeroInt()),
				sdk.NewCoin("htltinc", sdk.ZeroInt()),
				sdk.NewCoin("htltinc", sdk.ZeroInt()),
				time.Duration(0),
			),
		},
		PreviousBlockTime: types.DefaultPreviousBlockTime,
	}
}

// GeneratePrivKeyAddressPairsFromRand generates (deterministically) a total of n secp256k1 private keys and addresses.
func GeneratePrivKeyAddressPairs(n int) (keys []crypto.PrivKey, addrs []sdk.AccAddress) {
	r := rand.New(rand.NewSource(12345)) // make the generation deterministic
	keys = make([]crypto.PrivKey, n)
	addrs = make([]sdk.AccAddress, n)
	for i := 0; i < n; i++ {
		secret := make([]byte, 32)
		if _, err := r.Read(secret); err != nil {
			panic("Could not read randomness")
		}
		keys[i] = secp256k1.GenPrivKeySecp256k1(secret)
		addrs[i] = sdk.AccAddress(keys[i].PubKey().Address())
	}
	return
}

func GenerateRandomSecret() ([]byte, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return []byte{}, err
	}
	return bytes, nil
}

func loadSwapAndSupply(addr sdk.AccAddress, index int) (types.HTLC, types.AssetSupply) {
	coin := c(DenomMap[index], 50000)
	expireOffset := MinTimeLock
	timestamp := ts(index)
	amount := cs(coin)
	randomSecret, _ := GenerateRandomSecret()
	randomHashLock := types.GetHashLock(randomSecret, timestamp)
	id := types.GetID(addr, addr, amount, randomHashLock)
	htlc := types.NewHTLC(
		id,
		addr,
		addr,
		ReceiverOnOtherChain,
		SenderOnOtherChain,
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

	supply := types.NewAssetSupply(
		coin,
		c(coin.Denom, 0),
		c(coin.Denom, 0),
		c(coin.Denom, 0),
		time.Duration(0),
	)

	return htlc, supply
}
