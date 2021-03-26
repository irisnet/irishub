package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/htlc/types"
)

type ParamsTestSuite struct {
	suite.Suite
	addr   sdk.AccAddress
	supply []types.SupplyLimit
}

func (suite *ParamsTestSuite) SetupTest() {
	_, addrs := GeneratePrivKeyAddressPairs(1)
	suite.addr = addrs[0]
	supply1 := types.SupplyLimit{
		Limit:          sdk.NewInt(10000000000000),
		TimeLimited:    false,
		TimeBasedLimit: sdk.ZeroInt(),
		TimePeriod:     time.Hour,
	}
	supply2 := types.SupplyLimit{
		Limit:          sdk.NewInt(10000000000000),
		TimeLimited:    true,
		TimeBasedLimit: sdk.NewInt(100000000000),
		TimePeriod:     time.Hour * 24,
	}
	suite.supply = append(suite.supply, supply1, supply2)
}

func (suite *ParamsTestSuite) TestParamValidation() {

	type args struct {
		assetParams []types.AssetParam
	}

	testCases := []struct {
		name       string
		args       args
		expectPass bool
	}{{
		name: "default",
		args: args{
			assetParams: []types.AssetParam{},
		},
		expectPass: true,
	}, {
		name: "valid single asset",
		args: args{
			assetParams: []types.AssetParam{
				types.NewAssetParam(
					"htltbnb",
					714,
					suite.supply[0],
					true,
					suite.addr.String(),
					sdk.NewInt(1000),
					sdk.NewInt(100000000),
					sdk.NewInt(100000000000),
					MinTimeLock,
					MaxTimeLock,
				),
			},
		},
		expectPass: true,
	}, {
		name: "valid single asset time limited",
		args: args{
			assetParams: []types.AssetParam{
				types.NewAssetParam(
					"htltbnb",
					714,
					suite.supply[1],
					true,
					suite.addr.String(),
					sdk.NewInt(1000),
					sdk.NewInt(100000000),
					sdk.NewInt(100000000000),
					MinTimeLock,
					MaxTimeLock,
				),
			},
		},
		expectPass: true,
	}, {
		name: "valid multi asset",
		args: args{
			assetParams: []types.AssetParam{
				types.NewAssetParam(
					"htltbnb",
					714,
					suite.supply[0],
					true,
					suite.addr.String(),
					sdk.NewInt(1000),
					sdk.NewInt(100000000),
					sdk.NewInt(100000000000),
					MinTimeLock,
					MaxTimeLock,
				),
				types.NewAssetParam(
					"htltbtcb",
					0,
					suite.supply[1],
					true,
					suite.addr.String(),
					sdk.NewInt(1000),
					sdk.NewInt(10000000),
					sdk.NewInt(100000000000),
					MinTimeLock,
					MaxTimeLock,
				),
			},
		},
		expectPass: true,
	}, {
		name: "invalid denom - empty",
		args: args{
			assetParams: []types.AssetParam{
				types.NewAssetParam(
					"", 714,
					suite.supply[0],
					true,
					suite.addr.String(),
					sdk.NewInt(1000),
					sdk.NewInt(100000000),
					sdk.NewInt(100000000000),
					MinTimeLock,
					MaxTimeLock,
				),
			},
		},
		expectPass: false,
	}, {
		name: "invalid denom - bad format",
		args: args{
			assetParams: []types.AssetParam{
				types.NewAssetParam(
					"htltBNB",
					714,
					suite.supply[0],
					true,
					suite.addr.String(),
					sdk.NewInt(1000),
					sdk.NewInt(100000000),
					sdk.NewInt(100000000000),
					MinTimeLock,
					MaxTimeLock,
				),
			},
		},
		expectPass: false,
	}, {
		name: "min block lock equal max block lock",
		args: args{
			assetParams: []types.AssetParam{
				types.NewAssetParam(
					"htltbnb",
					714,
					suite.supply[0],
					true,
					suite.addr.String(),
					sdk.NewInt(1000),
					sdk.NewInt(100000000),
					sdk.NewInt(100000000000),
					243,
					243,
				),
			},
		},
		expectPass: true,
	}, {
		name: "min block lock greater max block lock",
		args: args{
			assetParams: []types.AssetParam{
				types.NewAssetParam(
					"htltbnb",
					714,
					suite.supply[0],
					true,
					suite.addr.String(),
					sdk.NewInt(1000),
					sdk.NewInt(100000000),
					sdk.NewInt(100000000000),
					244,
					243,
				),
			},
		},
		expectPass: false,
	}, {
		name: "min swap not positive",
		args: args{
			assetParams: []types.AssetParam{
				types.NewAssetParam(
					"htltbnb",
					714,
					suite.supply[0],
					true,
					suite.addr.String(),
					sdk.NewInt(1000),
					sdk.NewInt(0),
					sdk.NewInt(10000000000),
					MinTimeLock,
					MaxTimeLock,
				),
			},
		},
		expectPass: false,
	}, {
		name: "max swap not positive",
		args: args{
			assetParams: []types.AssetParam{
				types.NewAssetParam(
					"htltbnb",
					714,
					suite.supply[0],
					true,
					suite.addr.String(),
					sdk.NewInt(1000),
					sdk.NewInt(10000),
					sdk.NewInt(0),
					MinTimeLock,
					MaxTimeLock,
				),
			},
		},
		expectPass: false,
	}, {
		name: "min swap greater max swap",
		args: args{
			assetParams: []types.AssetParam{
				types.NewAssetParam(
					"htltbnb",
					714,
					suite.supply[0],
					true,
					suite.addr.String(),
					sdk.NewInt(1000),
					sdk.NewInt(100000000000),
					sdk.NewInt(10000000000),
					MinTimeLock,
					MaxTimeLock,
				),
			},
		},
		expectPass: false,
	}, {
		name: "negative asset limit",
		args: args{
			assetParams: []types.AssetParam{
				types.NewAssetParam(
					"htltbnb", 714,
					types.SupplyLimit{
						sdk.NewInt(-10000000000000),
						false,
						time.Hour,
						sdk.ZeroInt(),
					}, true,
					suite.addr.String(),
					sdk.NewInt(1000),
					sdk.NewInt(100000000),
					sdk.NewInt(100000000000),
					MinTimeLock,
					MaxTimeLock,
				),
			},
		},
		expectPass: false,
	}, {
		name: "negative asset time limit",
		args: args{
			assetParams: []types.AssetParam{
				types.NewAssetParam(
					"htltbnb",
					714,
					types.SupplyLimit{
						sdk.NewInt(10000000000000),
						false, time.Hour,
						sdk.NewInt(-10000000000000),
					},
					true,
					suite.addr.String(),
					sdk.NewInt(1000),
					sdk.NewInt(100000000),
					sdk.NewInt(100000000000),
					MinTimeLock,
					MaxTimeLock,
				),
			},
		},
		expectPass: false,
	}, {
		name: "asset time limit greater than overall limit",
		args: args{
			assetParams: []types.AssetParam{
				types.NewAssetParam(
					"htltbnb",
					714,
					types.SupplyLimit{
						sdk.NewInt(10000000000000),
						true,
						time.Hour,
						sdk.NewInt(100000000000000),
					},
					true,
					suite.addr.String(),
					sdk.NewInt(1000),
					sdk.NewInt(100000000),
					sdk.NewInt(100000000000),
					MinTimeLock,
					MaxTimeLock,
				),
			},
		},
		expectPass: false,
	}, {
		name: "duplicate denom",
		args: args{
			assetParams: []types.AssetParam{
				types.NewAssetParam(
					"htltbnb",
					714,
					suite.supply[0],
					true,
					suite.addr.String(),
					sdk.NewInt(1000),
					sdk.NewInt(100000000),
					sdk.NewInt(100000000000),
					MinTimeLock,
					MaxTimeLock,
				),
				types.NewAssetParam(
					"htltbnb",
					0,
					suite.supply[0],
					true,
					suite.addr.String(),
					sdk.NewInt(1000),
					sdk.NewInt(10000000),
					sdk.NewInt(100000000000),
					MinTimeLock,
					MaxTimeLock,
				),
			},
		},
		expectPass: false,
	}}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			params := types.NewParams(tc.args.assetParams)
			err := params.Validate()
			if tc.expectPass {
				suite.Require().NoError(err, tc.name)
			} else {
				suite.Require().Error(err, tc.name)
			}
		})
	}
}

func TestParamsTestSuite(t *testing.T) {
	suite.Run(t, new(ParamsTestSuite))
}
