package types

import (
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestFarmPool_ExpiredHeight(t *testing.T) {
	type fields struct {
		Name        string
		BeginHeight int64
		Rules       []RewardRule
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{
			name: "test case 1",
			fields: fields{
				Name:        "pool_1",
				BeginHeight: 100,
				Rules: []RewardRule{
					{
						Reward:         sdk.DefaultBondDenom,
						TotalReward:    sdk.NewInt(10000),
						RewardPerBlock: sdk.NewInt(100),
					},
					{
						Reward:         sdk.DefaultBondDenom,
						TotalReward:    sdk.NewInt(100000),
						RewardPerBlock: sdk.NewInt(100),
					},
				},
			},
			want: 200,
		},
		{
			name: "test case 2",
			fields: fields{
				Name:        "pool_2",
				BeginHeight: 10,
				Rules: []RewardRule{
					{
						Reward:         sdk.DefaultBondDenom,
						TotalReward:    sdk.NewInt(10000),
						RewardPerBlock: sdk.NewInt(100),
					},
				},
			},
			want: 110,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := FarmPool{
				Name:        tt.fields.Name,
				StartHeight: tt.fields.BeginHeight,
				Rules:       tt.fields.Rules,
			}
			expiredHeight, _ := fp.ExpiredHeight()
			if got := expiredHeight; got != tt.want {
				t.Errorf("FarmPool.ExpiredHeight() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFarmPool_CaclRewards(t *testing.T) {
	type fields struct {
		Rules []RewardRule
	}
	type args struct {
		farmInfo FarmInfo
		deltaAmt sdk.Int
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantRewards    sdk.Coins
		wantDewardDebt sdk.Coins
	}{
		{
			name: "test case 1",
			fields: fields{
				Rules: []RewardRule{
					{
						Reward:         sdk.DefaultBondDenom,
						RewardPerShare: sdk.NewDec(100),
					},
					{
						Reward:         "uiris",
						RewardPerShare: sdk.NewDecWithPrec(1, 1),
					},
				},
			},
			args: args{
				farmInfo: FarmInfo{
					Locked:     sdk.NewInt(10),
					RewardDebt: []sdk.Coin{},
				},
				deltaAmt: sdk.NewInt(10),
			},
			wantRewards: sdk.NewCoins(
				sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1000)),
				sdk.NewCoin("uiris", sdk.NewInt(1)),
			),
			wantDewardDebt: sdk.NewCoins(
				sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(2000)),
				sdk.NewCoin("uiris", sdk.NewInt(2)),
			),
		},
		{
			name: "test case 2",
			fields: fields{
				Rules: []RewardRule{
					{
						Reward:         sdk.DefaultBondDenom,
						RewardPerShare: sdk.NewDec(100),
					},
					{
						Reward:         "uiris",
						RewardPerShare: sdk.NewDecWithPrec(10, 1),
					},
				},
			},
			args: args{
				farmInfo: FarmInfo{
					Locked: sdk.NewInt(10),
					RewardDebt: sdk.NewCoins(
						sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(500)),
						sdk.NewCoin("uiris", sdk.NewInt(5)),
					),
				},
				deltaAmt: sdk.NewInt(10),
			},
			wantRewards: sdk.NewCoins(
				sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(500)),
				sdk.NewCoin("uiris", sdk.NewInt(5)),
			),
			wantDewardDebt: sdk.NewCoins(
				sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(2000)),
				sdk.NewCoin("uiris", sdk.NewInt(20)),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := FarmPool{
				Rules: tt.fields.Rules,
			}
			gotRewards, gotDewardDebt := fp.CaclRewards(tt.args.farmInfo, tt.args.deltaAmt)
			if !reflect.DeepEqual(gotRewards, tt.wantRewards) {
				t.Errorf("FarmPool.CaclRewards() gotRewards = %v, want %v", gotRewards, tt.wantRewards)
			}
			if !reflect.DeepEqual(gotDewardDebt, tt.wantDewardDebt) {
				t.Errorf("FarmPool.CaclRewards() gotDewardDebt = %v, want %v", gotDewardDebt, tt.wantDewardDebt)
			}
		})
	}
}
