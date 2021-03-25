package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/htlc/types"
)

type GenesisTestSuite struct {
	suite.Suite
	htlcs    []types.HTLC
	supplies []types.AssetSupply
}

func (suite *GenesisTestSuite) SetupTest() {

	coin := sdk.NewCoin("htltbnb", sdk.OneInt())
	suite.htlcs = htlcs(10)

	supply := types.NewAssetSupply(coin, coin, coin, coin, time.Duration(0))
	suite.supplies = []types.AssetSupply{supply}
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

func (suite *GenesisTestSuite) TestValidate() {
	type args struct {
		swaps             []types.HTLC
		supplies          []types.AssetSupply
		previousBlockTime time.Time
	}
	testCases := []struct {
		name       string
		args       args
		expectPass bool
	}{{
		"default",
		args{
			swaps:             []types.HTLC{},
			previousBlockTime: types.DefaultPreviousBlockTime,
		},
		true,
	}, {
		"with swaps",
		args{
			swaps:             suite.htlcs,
			previousBlockTime: types.DefaultPreviousBlockTime,
		},
		true,
	}, {
		"with supplies",
		args{
			swaps:             []types.HTLC{},
			supplies:          suite.supplies,
			previousBlockTime: types.DefaultPreviousBlockTime,
		},
		true,
	}, {
		"invalid supply",
		args{
			swaps:             []types.HTLC{},
			supplies:          []types.AssetSupply{{IncomingSupply: sdk.Coin{Denom: "Invalid", Amount: sdk.ZeroInt()}}},
			previousBlockTime: types.DefaultPreviousBlockTime,
		},
		false,
	}, {
		"duplicate swaps",
		args{
			swaps:             []types.HTLC{suite.htlcs[2], suite.htlcs[2]},
			previousBlockTime: types.DefaultPreviousBlockTime,
		},
		false,
	}, {
		"invalid swap",
		args{
			swaps:             []types.HTLC{{Amount: sdk.Coins{sdk.Coin{Denom: "Invalid Denom", Amount: sdk.NewInt(-1)}}}},
			previousBlockTime: types.DefaultPreviousBlockTime,
		},
		false,
	}, {
		"blocktime not set",
		args{
			swaps: []types.HTLC{},
		},
		true,
	}}

	for _, tc := range testCases {
		suite.Run(
			tc.name,
			func() {
				var gs *types.GenesisState
				if tc.name == "default" {
					gs = types.DefaultGenesisState()
				} else {
					gs = types.NewGenesisState(types.DefaultParams(), tc.args.swaps, tc.args.supplies, tc.args.previousBlockTime)
				}

				err := types.ValidateGenesis(*gs)
				if tc.expectPass {
					suite.Require().NoError(err, tc.name)
				} else {
					suite.Require().Error(err, tc.name)
				}
			},
		)
	}
}
