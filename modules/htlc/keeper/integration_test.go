package keeper_test

import (
	"math/rand"
	"time"

	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/htlc/types"
)

var (
	DenomMap = map[int]string{
		0: "htltbtc",
		1: "htlteth",
		2: "htltbnb",
		3: "htltxrp",
		4: "htltdai",
	}
	TestDeputy                  = sdk.AccAddress(crypto.AddressHash([]byte("TestDeputy")))
	TestUser1                   = sdk.AccAddress(crypto.AddressHash([]byte("TestUser1")))
	TestUser2                   = sdk.AccAddress(crypto.AddressHash([]byte("TestUser2")))
	MinTimeLock          uint64 = 220
	MaxTimeLock          uint64 = 270
	ReceiverOnOtherChain        = "ReceiverOnOtherChain"
	SenderOnOtherChain          = "SenderOnOtherChain"
)

func c(denom string, amount int64) sdk.Coin {
	return sdk.NewInt64Coin(denom, amount)
}

func cs(coins ...sdk.Coin) sdk.Coins {
	return sdk.NewCoins(coins...)
}

func ts(minOffset int) uint64 {
	return uint64(time.Now().Add(time.Duration(minOffset) * time.Minute).Unix())
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
					DeputyAddress: deputyAddress.String(),
					FixedFee:      sdk.NewInt(1000),
					MinSwapAmount: sdk.OneInt(),
					MaxSwapAmount: sdk.NewInt(1000000000000),
					MinBlockLock:  MinTimeLock,
					MaxBlockLock:  MaxTimeLock,
				},
				{
					Denom: "htltinc",
					SupplyLimit: types.SupplyLimit{
						Limit:          sdk.NewInt(100000000000000),
						TimeLimited:    true,
						TimeBasedLimit: sdk.NewInt(50000000000),
						TimePeriod:     time.Hour,
					},
					Active:        false,
					DeputyAddress: deputyAddress.String(),
					FixedFee:      sdk.NewInt(1000),
					MinSwapAmount: sdk.OneInt(),
					MaxSwapAmount: sdk.NewInt(100000000000),
					MinBlockLock:  MinTimeLock,
					MaxBlockLock:  MaxTimeLock,
				},
			},
		},
		Htlcs: []types.HTLC{},
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
