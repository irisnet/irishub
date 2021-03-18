package keeper_test

// import (
// 	"crypto/rand"
// 	"time"

// 	"github.com/tendermint/tendermint/crypto"
// 	tmtime "github.com/tendermint/tendermint/types/time"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

// 	"github.com/irisnet/irismod/modules/htlc/types"
// 	"github.com/irisnet/irismod/simapp"
// )

// var (
// 	DenomMap  = map[int]string{0: "htltbtc", 1: "htlteth", 2: "htltbnb", 3: "htltxrp", 4: "htltdai"}
// 	TestUser1 = sdk.AccAddress(crypto.AddressHash([]byte("KavaTestUser1")))
// 	TestUser2 = sdk.AccAddress(crypto.AddressHash([]byte("KavaTestUser2")))
// )

// func i(in int64) sdk.Int                    { return sdk.NewInt(in) }
// func c(denom string, amount int64) sdk.Coin { return sdk.NewInt64Coin(denom, amount) }
// func cs(coins ...sdk.Coin) sdk.Coins        { return sdk.NewCoins(coins...) }
// func ts(minOffset int) uint64 {
// 	return uint64(tmtime.Now().Add(time.Duration(minOffset) * time.Minute).Unix())
// }

// func NewAuthGenStateFromAccs(accounts ...authtypes.GenesisAccount) simapp.GenesisState {
// 	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), accounts)
// 	return simapp.GenesisState{
// 		authtypes.ModuleName: authtypes.ModuleCdc.MustMarshalJSON(authGenesis),
// 	}
// }

// func NewHTLTGenStateMulti(deputyAddress sdk.AccAddress) simapp.GenesisState {
// 	htlcGenesis := types.GenesisState{
// 		Params: types.Params{
// 			AssetParams: []types.AssetParam{
// 				{
// 					Denom: "htltbnb",
// 					SupplyLimit: types.SupplyLimit{
// 						Limit:          sdk.NewInt(350000000000000),
// 						TimeLimited:    false,
// 						TimeBasedLimit: sdk.ZeroInt(),
// 						TimePeriod:     time.Hour,
// 					},
// 					Active:        true,
// 					DeputyAddress: deputyAddress.String(),
// 					FixedFee:      sdk.NewInt(1000),
// 					MinSwapAmount: sdk.OneInt(),
// 					MaxSwapAmount: sdk.NewInt(1000000000000),
// 					MinBlockLock:  types.DefaultMinBlockLock,
// 					MaxBlockLock:  types.DefaultMaxBlockLock,
// 				},
// 				{
// 					Denom: "htltinc",
// 					SupplyLimit: types.SupplyLimit{
// 						Limit:          sdk.NewInt(100000000000000),
// 						TimeLimited:    true,
// 						TimeBasedLimit: sdk.NewInt(50000000000),
// 						TimePeriod:     time.Hour,
// 					},
// 					Active:        false,
// 					DeputyAddress: deputyAddress.String(),
// 					FixedFee:      sdk.NewInt(1000),
// 					MinSwapAmount: sdk.OneInt(),
// 					MaxSwapAmount: sdk.NewInt(100000000000),
// 					MinBlockLock:  types.DefaultMinBlockLock,
// 					MaxBlockLock:  types.DefaultMaxBlockLock,
// 				},
// 			},
// 		},
// 		Supplies: []types.AssetSupply{
// 			types.NewAssetSupply(
// 				sdk.NewCoin("htltbnb", sdk.ZeroInt()),
// 				sdk.NewCoin("htltbnb", sdk.ZeroInt()),
// 				sdk.NewCoin("htltbnb", sdk.ZeroInt()),
// 				sdk.NewCoin("htltbnb", sdk.ZeroInt()),
// 				time.Duration(0),
// 			),
// 			types.NewAssetSupply(
// 				sdk.NewCoin("htltinc", sdk.ZeroInt()),
// 				sdk.NewCoin("htltinc", sdk.ZeroInt()),
// 				sdk.NewCoin("htltinc", sdk.ZeroInt()),
// 				sdk.NewCoin("htltinc", sdk.ZeroInt()),
// 				time.Duration(0),
// 			),
// 		},
// 		PreviousBlockTime: types.DefaultPreviousBlockTime,
// 	}
// 	return simapp.GenesisState{types.ModuleName: types.ModuleCdc.MustMarshalJSON(&htlcGenesis)}
// }

// func htlts(ctx sdk.Context, count int) []types.HTLC {
// 	var htlts []types.HTLC
// 	for i := 0; i < count; i++ {
// 		htlt := htlt(ctx, i)
// 		htlts = append(htlts, htlt)
// 	}
// 	return htlts
// }

// func htlt(ctx sdk.Context, index int) types.HTLC {
// 	expireOffset := uint64(200)
// 	timestamp := ts(index)
// 	amount := cs(c("htltbnb", 50000))

// 	secret, _ := generateSecret()
// 	hashLock := types.GetHashLock(secret[:], timestamp)
// 	id := types.GetID(TestUser1, TestUser2, amount, hashLock)

// 	return types.NewHTLC(
// 		id,
// 		TestUser1,
// 		TestUser2,
// 		receiverOnOtherChain,
// 		senderOnOtherChain,
// 		amount,
// 		hashLock,
// 		secret,
// 		timestamp,
// 		uint64(ctx.BlockHeight())+expireOffset,
// 		types.Open,
// 		0,
// 		true,
// 		types.Incoming,
// 	)
// }

// func generateSecret() ([]byte, error) {
// 	bytes := make([]byte, 32)
// 	if _, err := rand.Read(bytes); err != nil {
// 		return []byte{}, err
// 	}
// 	return bytes, nil
// }

// func assetSupplies(count int) []types.AssetSupply {
// 	if count > 5 {
// 		return []types.AssetSupply{}
// 	}

// 	var supplies []types.AssetSupply

// 	for i := 0; i < count; i++ {
// 		supply := assetSupply(DenomMap[i])
// 		supplies = append(supplies, supply)
// 	}
// 	return supplies
// }

// func assetSupply(denom string) types.AssetSupply {
// 	return types.NewAssetSupply(c(denom, 0), c(denom, 0), c(denom, 0), c(denom, 0), time.Duration(0))
// }
