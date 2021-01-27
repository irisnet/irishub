package types

import (
	"fmt"
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	token = Token{
		Symbol:        "iris",
		Name:          "irisnet",
		Scale:         18,
		MinUnit:       "atto",
		InitialSupply: 1000000,
		MaxSupply:     10000000,
		Mintable:      true,
		Owner:         "",
	}
)

func TestToken_ToMinCoin(t *testing.T) {
	type args struct {
		coin sdk.DecCoin
	}

	for i := uint32(0); i <= MaximumScale; i++ {
		token.Scale = i
		tests := []struct {
			name    string
			args    args
			want    sdk.Coin
			wantErr bool
			success bool
		}{
			{
				name:    fmt.Sprintf("Main Coin to Min Coin,scale=%d", i),
				wantErr: false,
				args:    args{coin: sdk.NewDecCoin(token.Symbol, sdk.NewInt(10))},
				want:    sdk.NewCoin(token.MinUnit, sdk.NewIntWithDecimal(10, int(token.Scale))),
				success: true,
			},
			{
				name:    fmt.Sprintf("Main Coin to Min Coin Failed,scale=%d", i),
				wantErr: false,
				args:    args{coin: sdk.NewDecCoin(token.Symbol, sdk.NewInt(10))},
				want:    sdk.NewCoin(token.MinUnit, sdk.NewInt(10)),
				success: (i == 0),
			},
			{
				name:    fmt.Sprintf("Min Coin to Min Coin Success,scale=%d", i),
				wantErr: false,
				args:    args{coin: sdk.NewDecCoin(token.MinUnit, sdk.NewInt(10))},
				want:    sdk.NewCoin(token.MinUnit, sdk.NewInt(10)),
				success: true,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tr := Token{
					Symbol:  token.Symbol,
					Scale:   token.Scale,
					MinUnit: token.MinUnit,
				}
				got, err := tr.ToMinCoin(tt.args.coin)
				if (err != nil) != tt.wantErr {
					t.Errorf("Token.ToMainCoin() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if tt.success != reflect.DeepEqual(got, tt.want) {
					t.Errorf("Token.ToMainCoin() = %v, want %v", got, tt.want)
				}
			})
		}
	}
}

func TestToken_ToMainCoin(t *testing.T) {
	type args struct {
		coin sdk.Coin
	}

	for i := uint32(0); i <= MaximumScale; i++ {
		token.Scale = i
		tests := []struct {
			name    string
			args    args
			want    sdk.DecCoin
			wantErr bool
			success bool
		}{
			{
				name:    "Main Coin to Main Coin",
				wantErr: false,
				args:    args{coin: sdk.NewCoin(token.Symbol, sdk.NewInt(10))},
				want:    sdk.NewInt64DecCoin(token.Symbol, 10),
				success: true,
			},
			{
				name:    "Min Coin to Main Coin Failed",
				wantErr: false,
				args:    args{coin: sdk.NewCoin(token.MinUnit, sdk.NewInt(10))},
				want:    sdk.NewInt64DecCoin(token.Symbol, 10),
				success: (i == 0),
			},
			{
				name:    "Min Coin to Main Coin Success",
				wantErr: false,
				args:    args{coin: sdk.NewCoin(token.MinUnit, sdk.NewInt(10))},
				want: sdk.NewDecCoinFromDec(token.Symbol,
					sdk.NewDecWithPrec(1, int64(token.Scale)).MulInt64(10)),
				success: true,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tr := Token{
					Symbol:  token.Symbol,
					Scale:   token.Scale,
					MinUnit: token.MinUnit,
				}
				got, err := tr.ToMainCoin(tt.args.coin)
				if (err != nil) != tt.wantErr {
					t.Errorf("Token.ToMainCoin() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if tt.success != reflect.DeepEqual(got, tt.want) {
					t.Errorf("Token.ToMainCoin() = %v, want %v", got, tt.want)
				}
			})
		}
	}
}
