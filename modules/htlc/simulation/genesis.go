package simulation

import (
	"encoding/json"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/tendermint/tendermint/crypto"

	"github.com/irisnet/irismod/modules/htlc/types"
)

var (
	MinTimeLock = uint64(50)
	MaxTimeLock = uint64(34560)
	Deputy      = sdk.AccAddress(crypto.AddressHash([]byte("Deputy")))
)

const (
	BNB_DENOM   = "htltbnb"
	OTHER_DENOM = "htltinc"
)

func RandomizedGenState(simState *module.SimulationState) {
	htlcGenesis := &types.GenesisState{
		Params: types.Params{
			AssetParams: []types.AssetParam{
				{
					Denom: BNB_DENOM,
					SupplyLimit: types.SupplyLimit{
						Limit:          sdk.NewInt(350000000000000),
						TimeLimited:    false,
						TimeBasedLimit: sdk.ZeroInt(),
						TimePeriod:     time.Hour,
					},
					Active:        true,
					DeputyAddress: Deputy.String(),
					FixedFee:      sdk.NewInt(1000),
					MinSwapAmount: sdk.NewInt(2000),
					MaxSwapAmount: sdk.NewInt(1000000000000),
					MinBlockLock:  MinTimeLock,
					MaxBlockLock:  MaxTimeLock,
				},
				{
					Denom: OTHER_DENOM,
					SupplyLimit: types.SupplyLimit{
						Limit:          sdk.NewInt(100000000000000),
						TimeLimited:    true,
						TimeBasedLimit: sdk.NewInt(50000000000),
						TimePeriod:     time.Hour,
					},
					Active:        false,
					DeputyAddress: Deputy.String(),
					FixedFee:      sdk.NewInt(1000),
					MinSwapAmount: sdk.NewInt(2000),
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

	bz, err := json.MarshalIndent(htlcGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated %s parameters:\n%s\n", types.ModuleName, bz)

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(htlcGenesis)
}
